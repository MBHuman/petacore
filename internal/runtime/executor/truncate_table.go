package executor

import (
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteTruncateTable удаляет все строки из таблицы
func ExecuteTruncateTable(stmt *statements.TruncateTableStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	tbl := table.NewTable(store, exCtx.Database, exCtx.Schema, stmt.TableName)
	return tbl.TruncateTable(stmt.TableName)
}
