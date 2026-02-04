package executor

import (
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteTruncateTable удаляет все строки из таблицы
func ExecuteTruncateTable(stmt *statements.TruncateTableStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	// Резолвим схему и имя таблицы
	schema, tableName := ComputeSchemaAndTableName(stmt.TableName, &exCtx)
	tbl := table.NewTable(store, exCtx.Database, schema, tableName)
	return tbl.TruncateTable(stmt.TableName)
}
