package storage

import (
	"context"
	"fmt"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"time"
)

// DistributedStorageVClock представляет распределенное хранилище с Vector Clock
// Реализует CP модель с quorum-based чтением
type DistributedStorageVClock struct {
	synchronizer   *distributed.SynchronizerVClock
	logicalClock   *core.LClock
	mvccVClock     *core.MVCCWithVClock
	isolationLevel core.IsolationLevel
	nodeID         string

	// Конфигурация кворума
	totalNodes int
	minAcks    int // Минимум подтверждений для безопасного чтения (обычно N/2 + 1)

	// Таймауты
	writeTimeout time.Duration
	readTimeout  time.Duration
}

// NewDistributedStorageVClock создает новое распределенное хранилище с VClock
// minAcks:
//
//	0 или не указан - используется N/2 + 1 (строгий quorum, по умолчанию)
//	-1 - используется N (все узлы, максимальная консистентность)
//	> 0 - конкретное значение
func NewDistributedStorageVClock(kvStore distributed.KVStore, nodeID string, totalNodes int, isolationLevel core.IsolationLevel, minAcks int) *DistributedStorageVClock {
	logicalClock := core.NewLClock()
	mvccVClock := core.NewMVCCWithVClock()

	// Обработка специальных значений minAcks
	if minAcks == 0 {
		// По умолчанию - простое большинство (quorum)
		minAcks = totalNodes/2 + 1
	} else if minAcks == -1 {
		// Все узлы - максимальная консистентность
		minAcks = totalNodes
	} else if minAcks < 0 {
		// Некорректное отрицательное значение - используем по умолчанию
		minAcks = totalNodes/2 + 1
	}

	// Проверка корректности minAcks
	if minAcks > totalNodes {
		minAcks = totalNodes
	}
	if minAcks < 1 {
		minAcks = 1
	}

	synchronizer := distributed.NewSynchronizerVClock(kvStore, mvccVClock, logicalClock, nodeID)

	return &DistributedStorageVClock{
		synchronizer:   synchronizer,
		logicalClock:   logicalClock,
		mvccVClock:     mvccVClock,
		isolationLevel: isolationLevel,
		nodeID:         nodeID,
		totalNodes:     totalNodes,
		minAcks:        minAcks,
		writeTimeout:   5 * time.Second,
		readTimeout:    1 * time.Second,
	}
}

// Start запускает синхронизацию
func (ds *DistributedStorageVClock) Start() error {
	return ds.synchronizer.Start()
}

// Stop останавливает синхронизацию
func (ds *DistributedStorageVClock) Stop() {
	ds.synchronizer.Stop()
}

// IsSynced проверяет, синхронизирован ли узел
func (ds *DistributedStorageVClock) IsSynced() bool {
	return ds.synchronizer.IsSynced()
}

// GetMinAcks возвращает минимум подтверждений для quorum
func (ds *DistributedStorageVClock) GetMinAcks() int {
	return ds.minAcks
}

// GetTotalNodes возвращает общее количество узлов
func (ds *DistributedStorageVClock) GetTotalNodes() int {
	return ds.totalNodes
}

// SetMinAcks устанавливает минимальное количество подтверждений
// minAcks:
//
//	0 - использовать N/2 + 1 (строгий quorum)
//	-1 - использовать N (все узлы)
//	> 0 - конкретное значение
func (ds *DistributedStorageVClock) SetMinAcks(minAcks int) {
	if minAcks == 0 {
		ds.minAcks = ds.totalNodes/2 + 1
	} else if minAcks == -1 {
		ds.minAcks = ds.totalNodes
	} else if minAcks < 0 {
		ds.minAcks = ds.totalNodes/2 + 1
	} else {
		ds.minAcks = minAcks
	}

	// Проверка корректности
	if ds.minAcks > ds.totalNodes {
		ds.minAcks = ds.totalNodes
	}
	if ds.minAcks < 1 {
		ds.minAcks = 1
	}
}

// RunTransaction выполняет транзакцию с quorum-based чтением
func (ds *DistributedStorageVClock) RunTransaction(txFunc func(tx *DistributedTransactionVClock) error) error {
	tx := NewDistributedTransactionVClock(
		ds.mvccVClock,
		ds.logicalClock,
		ds.synchronizer,
		ds.isolationLevel,
		ds.nodeID,
		ds.minAcks,
		ds.totalNodes,
	)
	defer tx.Release()

	tx.Begin()

	if err := txFunc(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// BeginTransaction начинает долгоживущую транзакцию
func (ds *DistributedStorageVClock) BeginTransaction() *DistributedTransactionVClock {
	tx := NewDistributedTransactionVClock(
		ds.mvccVClock,
		ds.logicalClock,
		ds.synchronizer,
		ds.isolationLevel,
		ds.nodeID,
		ds.minAcks,
		ds.totalNodes,
	)
	tx.Begin()
	return tx
}

// CommitTransaction коммитит долгоживущую транзакцию
func (ds *DistributedStorageVClock) CommitTransaction(tx *DistributedTransactionVClock) error {
	defer tx.Release()
	return tx.Commit()
}

// DistributedTransactionVClock транзакция с Vector Clock
type DistributedTransactionVClock struct {
	mvccVClock     *core.MVCCWithVClock
	logicalClock   *core.LClock
	synchronizer   *distributed.SynchronizerVClock
	isolationLevel core.IsolationLevel
	nodeID         string
	minAcks        int
	totalNodes     int

	// Локальные изменения в рамках транзакции
	localWrites map[string]string

	// Snapshot для Snapshot Isolation
	snapshotVClock *core.VectorClock
}

// NewDistributedTransactionVClock создает новую транзакцию
func NewDistributedTransactionVClock(
	mvccVClock *core.MVCCWithVClock,
	logicalClock *core.LClock,
	synchronizer *distributed.SynchronizerVClock,
	isolationLevel core.IsolationLevel,
	nodeID string,
	minAcks int,
	totalNodes int,
) *DistributedTransactionVClock {
	return &DistributedTransactionVClock{
		mvccVClock:     mvccVClock,
		logicalClock:   logicalClock,
		synchronizer:   synchronizer,
		isolationLevel: isolationLevel,
		nodeID:         nodeID,
		minAcks:        minAcks,
		totalNodes:     totalNodes,
		localWrites:    make(map[string]string),
	}
}

// Begin начинает транзакцию
func (dtx *DistributedTransactionVClock) Begin() {
	if dtx.isolationLevel == core.SnapshotIsolation {
		// Для SI фиксируем текущий VectorClock как snapshot
		dtx.snapshotVClock = dtx.synchronizer.GetGlobalVectorClock().Clone()
	}
}

// Read читает с quorum-based проверкой
// Возвращает последнюю БЕЗОПАСНУЮ версию (с подтвержденным кворумом)
// Если запись не синхронизирована на достаточное количество узлов - возвращает старую версию
func (dtx *DistributedTransactionVClock) Read(key string) (string, bool) {
	// Сначала проверяем локальные изменения
	if value, ok := dtx.localWrites[key]; ok {
		return value, true
	}

	// Определяем snapshot VClock в зависимости от уровня изоляции
	var snapshotVC *core.VectorClock
	if dtx.isolationLevel == core.SnapshotIsolation && dtx.snapshotVClock != nil {
		// Для Snapshot Isolation используем фиксированный snapshot
		snapshotVC = dtx.snapshotVClock
	} else {
		// Для ReadCommitted используем текущий global VClock
		snapshotVC = dtx.synchronizer.GetGlobalVectorClock()
	}

	// Читаем из MVCC с snapshot и quorum проверкой
	value, _, ok := dtx.mvccVClock.ReadWithSnapshot(key, snapshotVC, dtx.minAcks, dtx.totalNodes, dtx.nodeID)
	if !ok {
		return "", false
	}
	return value, true
}

// Write записывает в локальный буфер
func (dtx *DistributedTransactionVClock) Write(key string, value string) {
	dtx.localWrites[key] = value
}

// Commit фиксирует транзакцию
// 1. Синхронно пишет в ETCD с Vector Clock
// 2. Локально пишет в MVCC с Vector Clock (nodeID инкрементирован)
// 3. Фоновая синхронизация обновит Vector Clock на других узлах
func (dtx *DistributedTransactionVClock) Commit() error {
	if len(dtx.localWrites) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for key, value := range dtx.localWrites {
		// WriteThroughVClock:
		// 1. Создает Vector Clock с инкрементом для текущего узла
		// 2. Пишет в ETCD (синхронно, гарантия персистентности)
		// 3. Пишет в локальный MVCC с VClock
		// 4. Watch на других узлах увидит изменение и обновит их VClock
		if err := dtx.synchronizer.WriteThroughVClock(ctx, key, value); err != nil {
			return fmt.Errorf("failed to write key %s: %w", key, err)
		}
	}

	return nil
}

// GetLocalWrites возвращает локальные изменения
func (dtx *DistributedTransactionVClock) GetLocalWrites() map[string]string {
	return dtx.localWrites
}

// Release освобождает ресурсы транзакции
func (dtx *DistributedTransactionVClock) Release() {
	// Очищаем локальные записи
	for k := range dtx.localWrites {
		delete(dtx.localWrites, k)
	}
	dtx.snapshotVClock = nil
}
