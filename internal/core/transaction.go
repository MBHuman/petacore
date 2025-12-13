package core

import (
	"errors"
	"sync"
)

// IsolationLevel определяет уровень изоляции транзакции
type IsolationLevel int

const (
	// ReadCommitted - транзакция видит только committed данные,
	// при каждом чтении получает последнюю закоммиченную версию
	ReadCommitted IsolationLevel = iota
	// SnapshotIsolation - транзакция видит фиксированный snapshot данных
	// на момент начала транзакции
	SnapshotIsolation
)

// Transaction pool для переиспользования объектов транзакций
var transactionPool = sync.Pool{
	New: func() interface{} {
		return &Transaction{
			localWrites: make(map[string]string, 16),
		}
	},
}

type Transaction struct {
	mvcc         *MVCC
	logicalClock *LClock

	isolationLevel  IsolationLevel
	snapshotVersion *uint64
	localWrites     map[string]string
}

func NewTransaction(mvcc *MVCC, logicalClock *LClock, isolationLevel IsolationLevel) *Transaction {
	tx := transactionPool.Get().(*Transaction)
	tx.mvcc = mvcc
	tx.logicalClock = logicalClock
	tx.isolationLevel = isolationLevel
	tx.snapshotVersion = nil
	// Очищаем localWrites, но сохраняем capacity
	for k := range tx.localWrites {
		delete(tx.localWrites, k)
	}
	return tx
}

func (tx *Transaction) Begin() {
	// Для SnapshotIsolation фиксируем версию snapshot
	// Для ReadCommitted snapshotVersion остается nil
	if tx.isolationLevel == SnapshotIsolation {
		version := tx.logicalClock.Get()
		tx.snapshotVersion = &version
	}
}

func (tx *Transaction) Read(key string) (string, bool) {
	// Сначала проверяем локальные записи
	if value, ok := tx.localWrites[key]; ok {
		return value, true
	}

	// Для ReadCommitted читаем последнюю committed версию
	if tx.isolationLevel == ReadCommitted {
		currentVersion := int64(tx.logicalClock.Get())
		return tx.mvcc.Read(key, currentVersion)
	}

	// Для SnapshotIsolation используем фиксированную версию snapshot
	if tx.snapshotVersion == nil {
		return "", false
	}
	return tx.mvcc.Read(key, int64(*tx.snapshotVersion))
}

func (tx *Transaction) Write(key string, value string) {
	if tx.localWrites == nil {
		tx.localWrites = make(map[string]string)
	}
	tx.localWrites[key] = value
}

func (tx *Transaction) Commit() error {
	// Для SnapshotIsolation проверяем, что транзакция была начата
	if tx.isolationLevel == SnapshotIsolation && tx.snapshotVersion == nil {
		return errors.New("transaction has not been started")
	}

	// Записываем все локальные изменения в MVCC с новыми версиями
	for key, value := range tx.localWrites {
		newVersion := int64(tx.logicalClock.SendOrLocal())
		tx.mvcc.Write(key, value, newVersion)
	}

	return nil
}

// Release возвращает транзакцию в пул для переиспользования
func (tx *Transaction) Release() {
	tx.mvcc = nil
	tx.logicalClock = nil
	tx.snapshotVersion = nil
	// Не очищаем localWrites здесь - сделаем при следующем Get из пула
	transactionPool.Put(tx)
}

// GetSnapshotVersion возвращает версию snapshot для тестирования
func (tx *Transaction) GetSnapshotVersion() *uint64 {
	return tx.snapshotVersion
}

// GetLocalWrites возвращает локальные записи транзакции
func (tx *Transaction) GetLocalWrites() map[string]string {
	return tx.localWrites
}
