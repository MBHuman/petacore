package executor

import (
	"fmt"

	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/runtime/system"
	"petacore/internal/storage"
)

func ExecuteSelect(
	stmt *statements.SelectStatement,
	store *storage.DistributedStorageVClock,
	exCtx ExecutorContext,
) (map[string]interface{}, error) {
	fmt.Printf("DEBUG: execute select %v\n", stmt)

	// --- 1) SELECT без таблицы (функции/выражения) ---
	if stmt.TableName == "" {
		row := make(map[string]interface{})
		for _, col := range stmt.Columns {
			var colName string
			var value interface{}

			if col.Function != nil {
				v, err := functions.ExecuteFunction(col.Function.Name, col.Function.Args)
				if err != nil {
					return nil, err
				}
				value = v
				colName = col.Function.Name
			} else if col.ExpressionContext != nil {
				v, err := rhelpers.EvaluateExpressionContext(col.ExpressionContext)
				if err != nil {
					return nil, err
				}
				value = v
				colName = "?column?"
			} else if col.ColumnName != "" {
				colName = col.ColumnName
				value = nil
			}

			if col.Alias != "" {
				colName = col.Alias
			} else if col.ExpressionContext != nil {
				colName = "?column?"
			}

			row[colName] = value
		}

		cols := rhelpers.GetColumnNamesFromRow(row)
		types := rhelpers.GetColumnTypesFromRow(row)

		return map[string]interface{}{
			"rows":        []map[string]interface{}{row},
			"columns":     cols,
			"columnTypes": types,
		}, nil
	}

	// --- 2) System table select ---
	if system.IsSystemTable(stmt.TableName) {
		fullRows, err := system.ExecuteSystemTableSelect(stmt, store)
		if err != nil {
			return nil, err
		}

		// Filter to selected columns
		var filteredRows []map[string]interface{}
		for _, row := range fullRows {
			filteredRow := make(map[string]interface{})

			for _, col := range stmt.Columns {
				// SELECT *
				if (col.ExpressionContext != nil && col.ExpressionContext.GetText() == "*") || col.ColumnName == "*" {
					for k, v := range row {
						filteredRow[k] = v
					}
					continue
				}

				var colName string
				var value interface{}

				if col.ColumnName != "" {
					colName = col.ColumnName
					if v, ok := row[col.ColumnName]; ok {
						value = v
					} else {
						value = nil
					}
				} else if col.Function != nil {
					colName = col.Function.Name
					v, err := functions.ExecuteFunction(col.Function.Name, col.Function.Args)
					if err != nil {
						return nil, err
					}
					value = v
				} else if col.ExpressionContext != nil {
					colName = "?column?"
					v, err := rhelpers.EvaluateExpressionContext(col.ExpressionContext)
					if err != nil {
						return nil, err
					}
					value = v
				}

				if col.Alias != "" {
					colName = col.Alias
				} else if col.ExpressionContext != nil {
					colName = "?column?"
				}

				filteredRow[colName] = value
			}

			filteredRows = append(filteredRows, filteredRow)
		}

		// WHERE
		if stmt.Where != nil {
			var filtered []map[string]interface{}
			for _, row := range filteredRows {
				if rhelpers.EvaluateWhereCondition(stmt.Where, row) {
					filtered = append(filtered, row)
				}
			}
			filteredRows = filtered
		}

		rhelpers.SortRows(filteredRows, stmt.OrderBy)

		// columns/types: если нет строк — попытаться вывести из stmt (иначе будет пустой RowDescription)
		var cols []string
		var types []table.ColType
		if len(filteredRows) > 0 {
			cols = rhelpers.GetColumnNamesFromRow(filteredRows[0])
			types = rhelpers.GetColumnTypesFromRow(filteredRows[0])
		} else {
			// минимальный фоллбек: из SELECT списка (без типов)
			for _, item := range stmt.Columns {
				if item.ColumnName == "*" {
					// неизвестно какие именно колонки → оставим пусто
					cols = nil
					break
				}
				if item.ColumnName != "" {
					cols = append(cols, item.ColumnName)
				}
			}
			if len(cols) > 0 {
				types = make([]table.ColType, len(cols))
				for i := range types {
					types[i] = table.ColTypeString
				}
			}
		}

		return map[string]interface{}{
			"rows":        filteredRows,
			"columns":     cols,
			"columnTypes": types,
		}, nil
	}

	// --- 3) Normal table select ---
	tbl := &table.Table{
		Storage:  store,
		Database: exCtx.Database,
		Schema:   exCtx.Schema,
		Name:     stmt.TableName,
	}

	// columns from stmt
	var selectCols []string
	for _, item := range stmt.Columns {
		if item.ColumnName != "" {
			selectCols = append(selectCols, item.ColumnName)
		} else {
			return nil, fmt.Errorf("functions in table select not supported yet")
		}
	}

	var whereMap map[string]interface{}
	if stmt.Where != nil {
		whereMap = map[string]interface{}{
			stmt.Where.Field: stmt.Where.Value,
		}
	}

	var rows []map[string]interface{}
	var finalColumns []string
	var columnTypes []table.ColType

	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		var err error
		rows, finalColumns, columnTypes, err = tbl.Select(tx, stmt.TableName, selectCols, whereMap, stmt.Limit)
		return err
	})
	if err != nil {
		return nil, err
	}

	rhelpers.SortRows(rows, stmt.OrderBy)

	// ВАЖНО: даже если rows пустой — finalColumns/columnTypes должны быть заполнены tbl.Select()
	// чтобы wire мог послать RowDescription.
	return map[string]interface{}{
		"rows":        rows,
		"columns":     finalColumns,
		"columnTypes": columnTypes,
	}, nil
}
