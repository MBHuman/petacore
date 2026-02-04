package executor

import (
	"fmt"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

// ExecuteInsert вставляет данные
func ExecuteInsert(stmt *statements.InsertStatement, store *storage.DistributedStorageVClock, exCtx ExecutorContext) error {
	// Резолвим схему и имя таблицы
	schema, tableName := ComputeSchemaAndTableName(stmt.TableName, &exCtx)
	tbl := table.NewTable(store, exCtx.Database, schema, tableName)

	insertValues := make([][]interface{}, 0, len(stmt.Values))

	for _, rowValues := range stmt.Values {
		if len(rowValues) != len(stmt.Columns) {
			return fmt.Errorf("number of values does not match number of columns")
		}
		// Создаем map значений для этой строки
		// values := make(map[string]interface{})
		// for i, col := range stmt.Columns {
		// 	values[col] = rowValues[i]
		// }

		values := make([]interface{}, len(stmt.Columns))
		copy(values, rowValues)

		insertValues = append(insertValues, values)
	}

	columnNames := make([]string, 0, len(stmt.Columns))
	for _, col := range stmt.Columns {
		columnNames = append(columnNames, col)
	}

	if err := tbl.Insert(stmt.TableName, insertValues, columnNames); err != nil {
		return err
	}

	return nil
}
