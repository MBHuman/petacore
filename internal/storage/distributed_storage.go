package storage

import (
	"context"
	"fmt"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"time"
)

// DistributedStorage представляет распределенное хранилище с MVCC кешем
// Реализует CP (Consistency + Partition tolerance) модель
type DistributedStorage struct {
	synchronizer   *distributed.Synchronizer
	logicalClock   *core.LClock
	mvcc           *core.MVCC
	isolationLevel core.IsolationLevel

	// Таймауты для операций
	writeTimeout time.Duration
	readTimeout  time.Duration
}

// NewDistributedStorage создает новое распределенное хранилище
func NewDistributedStorage(kvStore distributed.KVStore, isolationLevel core.IsolationLevel) *DistributedStorage {
	logicalClock := core.NewLClock()
	mvcc := core.NewMVCC()

	synchronizer := distributed.NewSynchronizer(kvStore, mvcc, logicalClock)

	return &DistributedStorage{
		synchronizer:   synchronizer,
		logicalClock:   logicalClock,
		mvcc:           mvcc,
		isolationLevel: isolationLevel,
		writeTimeout:   5 * time.Second,
		readTimeout:    1 * time.Second,
	}
}

// Start запускает синхронизацию с распределенным хранилищем
func (ds *DistributedStorage) Start() error {
	return ds.synchronizer.Start()
}

// Stop останавливает синхронизацию
func (ds *DistributedStorage) Stop() {
	ds.synchronizer.Stop()
}

// IsSynced проверяет, синхронизирован ли узел
func (ds *DistributedStorage) IsSynced() bool {
	return ds.synchronizer.IsSynced()
}

// GetSyncStatus возвращает статус синхронизации
func (ds *DistributedStorage) GetSyncStatus() distributed.SyncStatus {
	return ds.synchronizer.GetStatus()
}

// RunTransaction выполняет транзакцию
// Для CP модели:
// - Чтение: из локального MVCC кеша (быстро)
// - Запись: через ETCD с синхронизацией на все узлы (медленно, но консистентно)
func (ds *DistributedStorage) RunTransaction(txFunc func(tx *DistributedTransaction) error) error {
	tx := NewDistributedTransaction(ds.mvcc, ds.logicalClock, ds.synchronizer, ds.isolationLevel)
	defer tx.Release()

	tx.Begin()

	if err := txFunc(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// DistributedTransaction транзакция для распределенного хранилища
type DistributedTransaction struct {
	*core.Transaction
	synchronizer *distributed.Synchronizer
}

// NewDistributedTransaction создает новую распределенную транзакцию
func NewDistributedTransaction(mvcc *core.MVCC, logicalClock *core.LClock,
	synchronizer *distributed.Synchronizer, isolationLevel core.IsolationLevel) *DistributedTransaction {

	coreTx := core.NewTransaction(mvcc, logicalClock, isolationLevel)

	return &DistributedTransaction{
		Transaction:  coreTx,
		synchronizer: synchronizer,
	}
}

// Read читает из локального MVCC кеша
// CP модель: если узел не синхронизирован, возвращаем старые данные
func (dtx *DistributedTransaction) Read(key string) (string, bool) {
	// Просто делегируем базовой транзакции - она читает из локального MVCC
	return dtx.Transaction.Read(key)
}

// Write записывает в локальный буфер (как обычно)
func (dtx *DistributedTransaction) Write(key string, value string) {
	dtx.Transaction.Write(key, value)
}

// Commit фиксирует транзакцию
// Записывает изменения в ETCD, который затем синхронизируется на все узлы
func (dtx *DistributedTransaction) Commit() error {
	// Получаем локальные изменения
	localWrites := dtx.GetLocalWrites()

	if len(localWrites) == 0 {
		// Нет изменений - просто завершаем
		return nil
	}

	// Записываем каждое изменение в ETCD
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for key, value := range localWrites {
		// WriteThrough записывает в ETCD и локальный MVCC
		if err := dtx.synchronizer.WriteThrough(ctx, key, value); err != nil {
			return fmt.Errorf("failed to write key %s: %w", key, err)
		}
	}

	return nil
}

// GetLocalWrites возвращает локальные изменения (для внутреннего использования)
func (dtx *DistributedTransaction) GetLocalWrites() map[string]string {
	// Это требует доступа к приватному полю localWrites из Transaction
	// Можно добавить метод в Transaction или использовать reflection
	// Для простоты добавим публичный метод в Transaction
	return dtx.Transaction.GetLocalWrites()
}

// Release возвращает транзакцию в пул
func (dtx *DistributedTransaction) Release() {
	dtx.Transaction.Release()
}
