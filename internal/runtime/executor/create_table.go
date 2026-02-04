package executor

import (
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteCreateTable создает таблицу
func ExecuteCreateTable(stmt *statements.CreateTableStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	// Резолвим схему и имя таблицы
	schema, tableName := ComputeSchemaAndTableName(stmt.TableName, &exCtx)
	tbl := table.NewTable(store, exCtx.Database, schema, tableName)

	columns := make([]table.ColumnDef, len(stmt.Columns))
	for i, col := range stmt.Columns {
		columns[i] = table.ColumnDef{
			Idx:          i + 1,
			Name:         col.Name,
			Type:         col.Type,
			IsNullable:   col.IsNullable,
			IsUnique:     col.IsUnique,
			IsSerial:     col.IsSerial,
			DefaultValue: col.DefaultValue,
		}
	}

	return tbl.CreateTable(
		stmt.TableName,
		columns,
		stmt.PrimaryKeys,
		stmt.IfNotExists,
		exCtx.IsInformationSchemaInit,
	)
}
