package distributed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"petacore/internal/core"
	"sync"
	"time"
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
		result[string(entry.Key)] = entry.Value
	}
	return result, nil
}

// VClockEntry структура для хранения в ETCD
type VClockEntry struct {
	Value       string            `json:"value"`
	Timestamp   uint64            `json:"timestamp"`
	VectorClock map[string]uint64 `json:"vclock"`
}

// syncLoopVClock непрерывно синхронизирует изменения через SyncIterator
func (s *SynchronizerVClock) syncLoopVClock() {
	defer s.wg.Done()

	// log.Printf("[SynchronizerVClock] Starting sync for node %s...", s.nodeID)

	for {
		select {
		case <-s.ctx.Done():
			// log.Println("[SynchronizerVClock] Sync loop stopped")
			return
		default:
		}

		if err := s.syncOnceVClock(); err != nil {
			// log.Printf("[SynchronizerVClock] Sync error: %v, retrying in 5s...", err)
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

	// log.Printf("[SynchronizerVClock] Starting sync iterator for node %s...", s.nodeID)
	eventChan, err := s.kvStore.SyncIterator(syncCtx, []byte{})
	if err != nil {
		return fmt.Errorf("failed to start sync iterator: %w", err)
	}

	// s.setStatus(SyncStatusSynced) // Убрано, теперь устанавливается по EventTypeSyncComplete

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
		log.Printf("[SynchronizerVClock] Sync complete for node %s\n", s.nodeID)
		log.Printf("[SynchronizerVClock] Node mvccVclock: %v is now synced", s.mvccVClock)
		val, _, _ := s.mvccVClock.ReadLatest([]byte("a"))
		log.Printf("[SynchronizerVClock] Node mvccVclock: key 'a' has value: %s", val)
		return
	}

	// Пропускаем удаления
	if event.Type == EventTypeDelete {
		return
	}

	// Парсим VClockEntry
	var vclockEntry VClockEntry
	if err := json.Unmarshal([]byte(event.Entry.Value), &vclockEntry); err != nil {
		log.Printf("[SynchronizerVClock] Warning: failed to parse VClock entry for key %s: %v\n", event.Entry.Key, err)
		return
	}

	log.Printf("[SynchronizerVClock] Loaded key %s, value %s, vclock %v\n", event.Entry.Key, vclockEntry.Value, vclockEntry.VectorClock)

	// Создаем Vector Clock
	vclock := core.NewVectorClock()
	vclock.UpdateFromMap(vclockEntry.VectorClock)

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
	s.logicalClock.Recv(vclockEntry.Timestamp)

	// Записываем в локальный MVCC с VClock (не модифицированным для своих записей)
	s.mvccVClock.WriteWithVClock(event.Entry.Key, vclockEntry.Value, vclockEntry.Timestamp, vclock)

	// Обновляем ревизию
	s.setLastSyncRevision(event.Entry.Revision)

	// log.Printf("[SynchronizerVClock] Node %s applied: key=%s, timestamp=%d, vclock=%v, revision=%d",
	// s.nodeID, event.Entry.Key, vclockEntry.Timestamp, vclock, event.Entry.Revision)
}

// WriteThroughVClock записывает в ETCD и локальный MVCC с Vector Clock
// Ключевой метод для CP модели:
// 1. Инкрементируем Vector Clock для текущего узла
// 2. Пишем в ETCD (синхронно, блокирующая операция)
// 3. Пишем в локальный MVCC
// 4. Другие узлы получат через watch и обновят свои VClock
func (s *SynchronizerVClock) WriteThroughVClock(ctx context.Context, key []byte, value string) error {
	// Инкрементируем логическое время
	timestamp := s.logicalClock.SendOrLocal()

	// Создаем новый Vector Clock с инкрементом для текущего узла
	s.vclockMu.Lock()
	vclock := s.globalVClock.Clone()
	vclock.Increment(s.nodeID)
	s.globalVClock.Update(vclock)
	s.vclockMu.Unlock()

	// Сериализуем в VClockEntry
	vclockEntry := VClockEntry{
		Value:       value,
		Timestamp:   timestamp,
		VectorClock: vclock.ToMap(),
	}

	entryJSON, err := json.Marshal(vclockEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal VClock entry: %w", err)
	}

	// Записываем в ETCD (синхронно, CP гарантия)
	if err := s.kvStore.Put(ctx, key, string(entryJSON), int64(timestamp)); err != nil {
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

	var vclockEntry VClockEntry
	if err := json.Unmarshal([]byte(entry.Value), &vclockEntry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VClock entry: %w", err)
	}

	vclock := core.NewVectorClock()
	vclock.UpdateFromMap(vclockEntry.VectorClock)
	return vclock, nil
}

// setLastSyncRevision устанавливает последнюю синхронизированную ревизию
func (s *SynchronizerVClock) setLastSyncRevision(revision int64) {
	s.revisionMu.Lock()
	defer s.revisionMu.Unlock()
	s.lastSyncRevision = revision
}
