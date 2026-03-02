package executor

import (
	"fmt"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
	ptypes "petacore/sdk/types"
)

// ExecuteInsert вставляет данные
func ExecuteInsert(stmt *statements.InsertStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	// Резолвим схему и имя таблицы
	// TODO перевести на новые типы данных и убрать лишнее копирование
	schemaName, tableName := ComputeSchemaAndTableName(stmt.TableName, &exCtx)
	tbl := table.NewTable(store, exCtx.Database, schemaName, tableName)

	insertValues := make([][]ptypes.BaseType[any], 0, len(stmt.Values))

	for _, rowValues := range stmt.Values {
		if len(rowValues) != len(stmt.Columns) {
			return fmt.Errorf("number of values does not match number of columns")
		}
		values := make([]ptypes.BaseType[any], len(stmt.Columns))
		// Convert interface{} values to BaseType[any]
		for i, val := range rowValues {
			// Use sdk/types conversion function
			bt, err := ptypes.ToBaseTypeAny(val)
			if err != nil {
				return fmt.Errorf("failed to convert value at index %d: %w", i, err)
			}
			if bt == nil {
				return fmt.Errorf("expected ptypes.BaseType[any] at index %d, got %T", i, val)
			}
			values[i] = bt
		}

		insertValues = append(insertValues, values)
	}

	columnNames := make([]string, 0, len(stmt.Columns))
	for _, col := range stmt.Columns {
		columnNames = append(columnNames, col)
	}

	if err := tbl.Insert(exCtx.Allocator, stmt.TableName, insertValues, columnNames); err != nil {
		return err
	}

	return nil
}
