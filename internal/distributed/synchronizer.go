package distributed

import (
	"context"
	"fmt"
	"petacore/internal/core"
	"sync"
	"time"
)

// SyncStatus статус синхронизации узла
type SyncStatus int

const (
	// SyncStatusSyncing узел синхронизируется
	SyncStatusSyncing SyncStatus = iota
	// SyncStatusSynced узел синхронизирован
	SyncStatusSynced
	// SyncStatusError ошибка синхронизации
	SyncStatusError
)

// Synchronizer синхронизирует данные между ETCD и локальным MVCC
type Synchronizer struct {
	kvStore      KVStore
	mvcc         *core.MVCC
	logicalClock *core.LClock

	status   SyncStatus
	statusMu sync.RWMutex

	lastSyncRevision int64
	revisionMu       sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewSynchronizer создает новый синхронизатор
func NewSynchronizer(kvStore KVStore, mvcc *core.MVCC, logicalClock *core.LClock) *Synchronizer {
	ctx, cancel := context.WithCancel(context.Background())

	return &Synchronizer{
		kvStore:      kvStore,
		mvcc:         mvcc,
		logicalClock: logicalClock,
		status:       SyncStatusSyncing,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start запускает синхронизацию
func (s *Synchronizer) Start() error {
	// Запускаем sync для непрерывной синхронизации с initial load
	s.wg.Add(1)
	go s.syncLoop()

	return nil
}

// Stop останавливает синхронизацию
func (s *Synchronizer) Stop() {
	s.cancel()
	s.wg.Wait()
}

// GetStatus возвращает текущий статус синхронизации
func (s *Synchronizer) GetStatus() SyncStatus {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}

// setStatus устанавливает статус синхронизации
func (s *Synchronizer) setStatus(status SyncStatus) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.status = status
}

// IsSynced проверяет, синхронизирован ли узел
func (s *Synchronizer) IsSynced() bool {
	return s.GetStatus() == SyncStatusSynced
}

// GetLastSyncRevision возвращает последнюю синхронизированную ревизию
func (s *Synchronizer) GetLastSyncRevision() int64 {
	s.revisionMu.RLock()
	defer s.revisionMu.RUnlock()
	return s.lastSyncRevision
}

// setLastSyncRevision устанавливает последнюю синхронизированную ревизию
func (s *Synchronizer) setLastSyncRevision(revision int64) {
	s.revisionMu.Lock()
	defer s.revisionMu.Unlock()
	s.lastSyncRevision = revision
}

// syncLoop непрерывно синхронизирует изменения через SyncIterator
func (s *Synchronizer) syncLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			// log.Println("[Synchronizer] Sync loop stopped")
			return
		default:
		}

		if err := s.syncOnce(); err != nil {
			// log.Printf("[Synchronizer] Sync error: %v, retrying in 5s...", err)
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

// syncOnce выполняет один цикл синхронизации через SyncIterator
func (s *Synchronizer) syncOnce() error {
	syncCtx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	// log.Println("[Synchronizer] Starting sync...")
	eventChan, err := s.kvStore.SyncIterator(syncCtx, "")
	if err != nil {
		return fmt.Errorf("failed to start sync iterator: %w", err)
	}

	s.setStatus(SyncStatusSynced)

	for {
		select {
		case <-s.ctx.Done():
			return nil

		case event, ok := <-eventChan:
			if !ok {
				return fmt.Errorf("sync iterator channel closed")
			}

			if err := s.handleWatchEvent(event); err != nil {
				// log.Printf("[Synchronizer] Failed to handle event: %v", err)
				// Продолжаем работу, не возвращаем ошибку
			}
		}
	}
}

// handleWatchEvent обрабатывает событие изменения из ETCD
func (s *Synchronizer) handleWatchEvent(event *WatchEvent) error {
	switch event.Type {
	case EventTypePut:
		if event.Entry == nil {
			return fmt.Errorf("put event without entry")
		}

		// log.Printf("[Synchronizer] Received PUT: key=%s, version=%d, revision=%d",
		// event.Entry.Key, event.Entry.Version, event.Entry.Revision)

		// Синхронизируем HLC
		s.logicalClock.Recv(uint64(event.Entry.Version))

		// Записываем в локальный MVCC
		s.mvcc.Write(event.Entry.Key, event.Entry.Value, event.Entry.Version)

		// Обновляем ревизию
		s.setLastSyncRevision(event.Entry.Revision)

	case EventTypeDelete:
		if event.Entry == nil {
			return fmt.Errorf("delete event without entry")
		}

		// log.Printf("[Synchronizer] Received DELETE: key=%s", event.Entry.Key)

		// Удаляем из локального MVCC
		s.mvcc.Delete(event.Entry.Key)
	}

	return nil
}

// WriteThrough записывает данные в ETCD и локальный MVCC
// Это метод для записи, который гарантирует консистентность
func (s *Synchronizer) WriteThrough(ctx context.Context, key string, value string) error {
	// Генерируем новую версию с помощью HLC
	version := int64(s.logicalClock.SendOrLocal())

	// Сначала пишем в ETCD (source of truth)
	if err := s.kvStore.Put(ctx, key, value, version); err != nil {
		return fmt.Errorf("failed to write to etcd: %w", err)
	}

	// log.Printf("[Synchronizer] WriteThrough: key=%s, version=%d", key, version)

	// ETCD Watch автоматически синхронизирует на все узлы
	// Но мы также записываем локально для немедленного доступа
	s.mvcc.Write(key, value, version)

	return nil
}

// ReadLocal читает из локального MVCC кеша
func (s *Synchronizer) ReadLocal(key string, version int64) (string, bool) {
	return s.mvcc.Read(key, version)
}
