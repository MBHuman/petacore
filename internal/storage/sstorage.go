package storage

import "petacore/internal/core"

type SimpleStorage struct {
	logicalClock   *core.LClock
	mvcc           *core.MVCC
	isolationLevel core.IsolationLevel
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		logicalClock:   core.NewLClock(),
		mvcc:           core.NewMVCC(),
		isolationLevel: core.ReadCommitted, // По умолчанию Read Committed
	}
}

func NewSimpleStorageWithIsolation(isolationLevel core.IsolationLevel) *SimpleStorage {
	return &SimpleStorage{
		logicalClock:   core.NewLClock(),
		mvcc:           core.NewMVCC(),
		isolationLevel: isolationLevel,
	}
}

func (ss *SimpleStorage) RunTransaction(txFunc func(tx *core.Transaction) error) error {
	tx := core.NewTransaction(ss.mvcc, ss.logicalClock, ss.isolationLevel)
	defer tx.Release() // Возвращаем транзакцию в пул

	tx.Begin()
	res := txFunc(tx)
	if res != nil {
		return res
	}
	return tx.Commit()
}
