package executor

import (
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteDropTable удаляет таблицу
func ExecuteDropTable(stmt *statements.DropTableStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	// Резолвим схему и имя таблицы
	schema, tableName := ComputeSchemaAndTableName(stmt.TableName, &exCtx)
	tbl := table.NewTable(store, exCtx.Database, schema, tableName)
	return tbl.DropTable(stmt.TableName)
}
