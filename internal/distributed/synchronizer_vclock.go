package distributed

import (
	"context"
	"encoding/json"
	"fmt"
	"petacore/internal/core"
	"sync"
	"time"
)

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
	// Инициализируем локальный MVCC данными из ETCD
	if err := s.initialSyncVClock(); err != nil {
		s.setStatus(SyncStatusError)
		return fmt.Errorf("initial sync failed: %w", err)
	}

	s.setStatus(SyncStatusSynced)

	// Запускаем watch для непрерывной синхронизации
	s.wg.Add(1)
	go s.watchLoopVClock()

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

// VClockEntry структура для хранения в ETCD
type VClockEntry struct {
	Value       string            `json:"value"`
	Timestamp   uint64            `json:"timestamp"`
	VectorClock map[string]uint64 `json:"vclock"`
}

// initialSyncVClock загружает все данные из ETCD в локальный MVCC с VClock
func (s *SynchronizerVClock) initialSyncVClock() error {
	// log.Printf("[SynchronizerVClock] Starting initial sync for node %s...", s.nodeID)

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()

	// Получаем все данные из ETCD
	entries, err := s.kvStore.GetAll(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get all entries: %w", err)
	}

	// log.Printf("[SynchronizerVClock] Loaded %d entries from ETCD", len(entries))

	// Загружаем в локальный MVCC
	maxRevision := int64(0)
	for _, entry := range entries {
		// Парсим VClockEntry из value
		var vclockEntry VClockEntry
		if err := json.Unmarshal([]byte(entry.Value), &vclockEntry); err != nil {
			// log.Printf("[SynchronizerVClock] Warning: failed to parse VClock entry for key %s: %v", entry.Key, err)
			// Пропускаем некорректные записи
			continue
		}

		// Синхронизируем HLC
		s.logicalClock.Recv(vclockEntry.Timestamp)

		// Создаем Vector Clock из map
		vclock := core.NewVectorClock()
		vclock.UpdateFromMap(vclockEntry.VectorClock)

		// Обновляем глобальный Vector Clock
		s.vclockMu.Lock()
		s.globalVClock.Update(vclock)
		s.vclockMu.Unlock()

		// Записываем в локальный MVCC
		s.mvccVClock.WriteWithVClock(entry.Key, vclockEntry.Value, vclockEntry.Timestamp, vclock)

		if entry.Revision > maxRevision {
			maxRevision = entry.Revision
		}
	}

	s.setLastSyncRevision(maxRevision)
	// log.Printf("[SynchronizerVClock] Initial sync completed, last revision: %d", maxRevision)

	return nil
}

// watchLoopVClock непрерывно следит за изменениями в ETCD
func (s *SynchronizerVClock) watchLoopVClock() {
	defer s.wg.Done()

	// log.Printf("[SynchronizerVClock] Starting watch for node %s...", s.nodeID)

	// Начинаем watch
	watchChan, err := s.kvStore.Watch(s.ctx, "")
	if err != nil {
		// log.Printf("[SynchronizerVClock] Failed to start watch: %v", err)
		return
	}

	for {
		select {
		case <-s.ctx.Done():
			// log.Println("[SynchronizerVClock] Watch loop stopped")
			return

		case event, ok := <-watchChan:
			if !ok {
				// log.Println("[SynchronizerVClock] Watch channel closed, reconnecting...")
				time.Sleep(1 * time.Second)
				watchChan, err = s.kvStore.Watch(s.ctx, "")
				if err != nil {
					// log.Printf("[SynchronizerVClock] Failed to reconnect watch: %v", err)
					return
				}
				continue
			}

			// Обрабатываем событие
			s.handleWatchEventVClock(event)
		}
	}
}

// handleWatchEventVClock обрабатывает событие от watch
func (s *SynchronizerVClock) handleWatchEventVClock(event *WatchEvent) {
	// Пропускаем удаления
	if event.Type == EventTypeDelete {
		return
	}

	// Парсим VClockEntry
	var vclockEntry VClockEntry
	if err := json.Unmarshal([]byte(event.Entry.Value), &vclockEntry); err != nil {
		// log.Printf("[SynchronizerVClock] Warning: failed to parse VClock entry for key %s: %v", event.Entry.Key, err)
		return
	}

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
func (s *SynchronizerVClock) WriteThroughVClock(ctx context.Context, key string, value string) error {
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

// GetLastSyncRevision возвращает последнюю синхронизированную ревизию
func (s *SynchronizerVClock) GetLastSyncRevision() int64 {
	s.revisionMu.RLock()
	defer s.revisionMu.RUnlock()
	return s.lastSyncRevision
}

// setLastSyncRevision устанавливает последнюю синхронизированную ревизию
func (s *SynchronizerVClock) setLastSyncRevision(revision int64) {
	s.revisionMu.Lock()
	defer s.revisionMu.Unlock()
	s.lastSyncRevision = revision
}
