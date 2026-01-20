package executor

import (
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteDropTable удаляет таблицу
func ExecuteDropTable(stmt *statements.DropTableStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	tbl := table.NewTable(store, exCtx.Database, exCtx.Schema, stmt.TableName)
	return tbl.DropTable(stmt.TableName)
}
