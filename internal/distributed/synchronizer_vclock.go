package distributed

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"petacore/internal/core"
	"petacore/internal/logger"
	"sync"
	"time"

	"github.com/hamba/avro/v2"

	"go.uber.org/zap"
)

var ErrKeyNotFound = errors.New("key not found")

// SynchronizerVClock синхронизатор с Vector Clock для quorum-based чтения
type SynchronizerVClock struct {
	kvStore      KVStore
	mvccVClock   *core.MVCCWithVClock
	logicalClock *core.LClock
	nodeID       string

	// Глобальный Vector Clock отслеживает все узлы
	globalVClock *core.VectorClock
	vclockMu     sync.RWMutex

	status   SyncStatus
	statusMu sync.RWMutex

	lastSyncRevision int64
	revisionMu       sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// buffer pool to reduce json marshal allocations
var vclockBufPool sync.Pool

// NewSynchronizerVClock создает новый синхронизатор с VClock
func NewSynchronizerVClock(kvStore KVStore, mvccVClock *core.MVCCWithVClock, logicalClock *core.LClock, nodeID string) *SynchronizerVClock {
	ctx, cancel := context.WithCancel(context.Background())

	return &SynchronizerVClock{
		kvStore:      kvStore,
		mvccVClock:   mvccVClock,
		logicalClock: logicalClock,
		nodeID:       nodeID,
		globalVClock: core.NewVectorClock(),
		status:       SyncStatusSyncing,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start запускает синхронизацию
func (s *SynchronizerVClock) Start() error {
	// Запускаем sync для непрерывной синхронизации с initial load
	s.wg.Add(1)
	go s.syncLoopVClock()

	return nil
}

// Stop останавливает синхронизацию
func (s *SynchronizerVClock) Stop() {
	s.cancel()
	s.wg.Wait()
}

// GetStatus возвращает текущий статус синхронизации
func (s *SynchronizerVClock) GetStatus() SyncStatus {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}

// setStatus устанавливает статус синхронизации
func (s *SynchronizerVClock) setStatus(status SyncStatus) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status = status
}

// IsSynced проверяет, синхронизирован ли узел
func (s *SynchronizerVClock) IsSynced() bool {
	return s.GetStatus() == SyncStatusSynced
}

// GetGlobalVectorClock возвращает копию глобального Vector Clock
func (s *SynchronizerVClock) GetGlobalVectorClock() *core.VectorClock {
	s.vclockMu.RLock()
	defer s.vclockMu.RUnlock()
	return s.globalVClock.Clone()
}

// ScanPrefix сканирует ключи с префиксом из ETCD
func (s *SynchronizerVClock) ScanPrefix(ctx context.Context, prefix []byte) (map[string]string, error) {
	entries, err := s.kvStore.ScanPrefix(ctx, prefix)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, entry := range entries {
		result[string(entry.Key)] = string(entry.Value)
	}
	return result, nil
}

// VClockEntry структура для хранения в ETCD
// VClockEntryMeta хранит метаданные записи, сериализуемые через Avro
type VClockEntryMeta struct {
	Timestamp   int64            `avro:"timestamp"`
	VectorClock map[string]int64 `avro:"vclock"`
}

var vclockAvroSchema = `{"type":"record","name":"VClockEntryMeta","fields":[{"name":"timestamp","type":"long"},{"name":"vclock","type":{"type":"map","values":"long"}}]}`
var vclockAvroParsed avro.Schema

func init() {
	vclockBufPool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
	vclockAvroParsed = avro.MustParse(vclockAvroSchema)
}

// syncLoopVClock непрерывно синхронизирует изменения через SyncIterator
func (s *SynchronizerVClock) syncLoopVClock() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		if err := s.syncOnceVClock(); err != nil {
			s.setStatus(SyncStatusError)

			select {
			case <-time.After(5 * time.Second):
				continue
			case <-s.ctx.Done():
				return
			}
		}
	}
}

// syncOnceVClock выполняет один цикл синхронизации через SyncIterator
func (s *SynchronizerVClock) syncOnceVClock() error {
	syncCtx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	eventChan, err := s.kvStore.SyncIterator(syncCtx, []byte{})
	if err != nil {
		return fmt.Errorf("failed to start sync iterator: %w", err)
	}

	for {
		select {
		case <-s.ctx.Done():
			return nil

		case event, ok := <-eventChan:
			if !ok {
				return fmt.Errorf("sync iterator channel closed")
			}

			// Обрабатываем событие
			s.handleWatchEventVClock(event)
		}
	}
}

// handleWatchEventVClock обрабатывает событие от watch
func (s *SynchronizerVClock) handleWatchEventVClock(event *WatchEvent) {
	// Обработка завершения синхронизации
	if event.Type == EventTypeSyncComplete {
		s.setStatus(SyncStatusSynced)
		logger.Info("[SynchronizerVClock] Sync complete", zap.String("nodeID", s.nodeID))
		return
	}

	// Пропускаем удаления
	if event.Type == EventTypeDelete {
		return
	}

	// Parse binary format: [4-byte metaLen][avro(meta)][value bytes]
	data := []byte(event.Entry.Value)
	if len(data) < 4 {
		logger.Warn("[SynchronizerVClock] Warning: invalid entry size", zap.String("key", string(event.Entry.Key)))
		return
	}
	metaLen := binary.BigEndian.Uint32(data[:4])
	if int(metaLen) > len(data)-4 {
		logger.Warn("[SynchronizerVClock] Warning: invalid meta length", zap.String("key", string(event.Entry.Key)))
		return
	}
	metaBytes := data[4 : 4+metaLen]
	valueBytes := data[4+metaLen:]

	var meta VClockEntryMeta
	if err := avro.Unmarshal(vclockAvroParsed, metaBytes, &meta); err != nil {
		logger.Warn("[SynchronizerVClock] Warning: failed to unmarshal avro meta",
			zap.String("key", string(event.Entry.Key)),
			zap.Error(err),
		)
		return
	}

	logger.Info("[SynchronizerVClock] Loaded key",
		zap.String("key", string(event.Entry.Key)),
		zap.Int64("timestamp", meta.Timestamp),
		zap.Any("vclock", meta.VectorClock),
	)
	// Создаем Vector Clock
	vclock := core.NewVectorClock()
	// Convert map[int64] to map[uint64]
	vmap := make(map[string]uint64, len(meta.VectorClock))
	for k, v := range meta.VectorClock {
		vmap[k] = uint64(v)
	}
	vclock.UpdateFromMap(vmap)

	// Проверяем: если это наша собственная запись, не инкрементируем повторно
	// Наша запись уже содержит наш nodeID в VClock
	currentNodeValue := vclock.Get(s.nodeID)
	if currentNodeValue == 0 {
		// Это запись с другого узла - инкрементируем для текущего узла
		// Это показывает, что данный узел видел и применил это изменение
		vclock.Increment(s.nodeID)
	}
	// Если currentNodeValue > 0, это наша собственная запись, уже содержит наш nodeID

	// Обновляем глобальный Vector Clock
	s.vclockMu.Lock()
	s.globalVClock.Update(vclock)
	s.vclockMu.Unlock()

	// Синхронизируем HLC
	s.logicalClock.Recv(uint64(meta.Timestamp))

	// Записываем в локальный MVCC с VClock (не модифицированным для своих записей)
	s.mvccVClock.WriteWithVClock(event.Entry.Key, valueBytes, uint64(meta.Timestamp), vclock)

	// Обновляем ревизию
	s.setLastSyncRevision(event.Entry.Revision)

	// logger.Info("[SynchronizerVClock] Node applied", zap.String("nodeID", s.nodeID), zap.String("key", string(event.Entry.Key)), zap.Int64("timestamp", vclockEntry.Timestamp), zap.Any("vclock", vclock), zap.Int64("revision", event.Entry.Revision))
}

// WriteThroughVClock записывает в ETCD и локальный MVCC с Vector Clock
// Ключевой метод для CP модели:
// 1. Инкрементируем Vector Clock для текущего узла
// 2. Пишем в ETCD (синхронно, блокирующая операция)
// 3. Пишем в локальный MVCC
// 4. Другие узлы получат через watch и обновят свои VClock
func (s *SynchronizerVClock) WriteThroughVClock(ctx context.Context, key []byte, value []byte) error {
	// Инкрементируем логическое время
	timestamp := s.logicalClock.SendOrLocal()

	// Создаем новый Vector Clock с инкрементом для текущего узла
	s.vclockMu.Lock()
	vclock := s.globalVClock.Clone()
	vclock.Increment(s.nodeID)
	s.globalVClock.Update(vclock)
	s.vclockMu.Unlock()

	// Сериализуем метаданные через Avro
	meta := VClockEntryMeta{
		Timestamp:   int64(timestamp),
		VectorClock: make(map[string]int64),
	}
	for k, v := range vclock.ToMap() {
		meta.VectorClock[k] = int64(v)
	}

	metaBytes, err := avro.Marshal(vclockAvroParsed, meta)
	if err != nil {
		return fmt.Errorf("failed to marshal avro meta: %w", err)
	}

	// Final payload: [4-byte metaLen][metaBytes][value bytes]
	final := make([]byte, 4+len(metaBytes)+len(value))
	binary.BigEndian.PutUint32(final[:4], uint32(len(metaBytes)))
	copy(final[4:4+len(metaBytes)], metaBytes)
	copy(final[4+len(metaBytes):], value)

	// Записываем в ETCD (синхронно, CP гарантия)
	if err := s.kvStore.Put(ctx, key, final, int64(timestamp)); err != nil {
		return fmt.Errorf("failed to write to ETCD: %w", err)
	}

	// Записываем в локальный MVCC
	s.mvccVClock.WriteWithVClock(key, value, timestamp, vclock)

	// log.Printf("[SynchronizerVClock] Node %s WriteThrough: key=%s, timestamp=%d, vclock=%v",
	// s.nodeID, key, timestamp, vclock)

	return nil
}

// GetKVStore возвращает KVStore для прямого доступа
func (s *SynchronizerVClock) GetKVStore() KVStore {
	return s.kvStore
}

// GetCurrentVersion получает текущую версию ключа из KV с Vector Clock
func (s *SynchronizerVClock) GetCurrentVersion(ctx context.Context, key []byte) (*core.VectorClock, error) {
	entry, err := s.kvStore.Get(ctx, key)
	if err != nil {
		if err == ErrKeyNotFound {
			return core.NewVectorClock(), nil
		}
		return nil, fmt.Errorf("failed to get current version: %w", err)
	}

	// Try to parse binary payload: [4-byte metaLen][avro(meta)][value bytes]
	data := []byte(entry.Value)
	if len(data) < 4 {
		// unknown format, return empty vclock
		return core.NewVectorClock(), nil
	}
	metaLen := binary.BigEndian.Uint32(data[:4])
	if int(metaLen) > len(data)-4 {
		return core.NewVectorClock(), nil
	}
	metaBytes := data[4 : 4+metaLen]

	var meta VClockEntryMeta
	if err := avro.Unmarshal(vclockAvroParsed, metaBytes, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal avro meta in GetCurrentVersion: %w", err)
	}

	vclock := core.NewVectorClock()
	vmap := make(map[string]uint64, len(meta.VectorClock))
	for k, v := range meta.VectorClock {
		vmap[k] = uint64(v)
	}
	vclock.UpdateFromMap(vmap)
	return vclock, nil
}

// setLastSyncRevision устанавливает последнюю синхронизированную ревизию
func (s *SynchronizerVClock) setLastSyncRevision(revision int64) {
	s.revisionMu.Lock()
	defer s.revisionMu.Unlock()
	s.lastSyncRevision = revision
}
