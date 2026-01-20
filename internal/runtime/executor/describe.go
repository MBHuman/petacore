package executor

import (
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteDescribe возвращает метаданные таблицы
func ExecuteDescribe(stmt *statements.DescribeStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) ([]map[string]interface{}, error) {
	// Create table instance
	t := &table.Table{
		Storage:  store,
		Database: exCtx.Database,
		Schema:   exCtx.Schema,
		Name:     stmt.TableName,
	}

	var result []map[string]interface{}
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		var err error
		result, err = t.Describe(tx)
		return err
	})
	return result, err
}
