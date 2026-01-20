package executor

import (
	"fmt"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteInsert вставляет данные
func ExecuteInsert(stmt *statements.InsertStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	tbl := table.NewTable(store, exCtx.Database, exCtx.Schema, stmt.TableName)

	insertValues := make([]map[string]interface{}, 0, len(stmt.Values))

	for _, rowValues := range stmt.Values {
		if len(rowValues) != len(stmt.Columns) {
			return fmt.Errorf("number of values does not match number of columns")
		}
		// Создаем map значений для этой строки
		values := make(map[string]interface{})
		for i, col := range stmt.Columns {
			values[col] = rowValues[i]
		}

		insertValues = append(insertValues, values)
	}
	if err := tbl.Insert(stmt.TableName, insertValues); err != nil {
		return err
	}

	return nil
}
