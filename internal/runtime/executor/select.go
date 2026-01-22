package executor

import (
	"log"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/runtime/system"
	"petacore/internal/storage"
)

// TODO рефакторинг: сильно раздутый ExecuteSelect разбить на подфункции по типам выборок
func ExecuteSelect(
	stmt *statements.SelectStatement,
	store *storage.DistributedStorageVClock,
	exCtx ExecutorContext,
) (map[string]interface{}, error) {
	log.Printf("DEBUG: execute select %v\n", stmt)

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
				v, err := rhelpers.EvaluateExpressionContext(col.ExpressionContext, row)
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

		// Build columns list from SELECT statement to preserve order
		var cols []string
		var types []table.ColType
		for _, col := range stmt.Columns {
			var colName string
			if col.Alias != "" {
				colName = col.Alias
			} else if col.ColumnName != "" {
				colName = col.ColumnName
			} else {
				colName = "?column?"
			}
			cols = append(cols, colName)

			// Get type from row value
			if v, ok := row[colName]; ok {
				switch v.(type) {
				case int, int32, int64:
					types = append(types, table.ColTypeInt)
				case float32, float64:
					types = append(types, table.ColTypeFloat)
				case bool:
					types = append(types, table.ColTypeBool)
				default:
					types = append(types, table.ColTypeString)
				}
			} else {
				types = append(types, table.ColTypeString)
			}
		}

		return map[string]interface{}{
			"rows":        []map[string]interface{}{row},
			"columns":     cols,
			"columnTypes": types,
		}, nil
	}

	// --- 2) System table select ---
	if system.IsSystemTable(stmt.TableName) {
		// Check if there are joins - if so, use executeFromClause
		if stmt.From != nil && len(stmt.From.Joins) > 0 {
			log.Printf("DEBUG: Executing FROM clause with joins for system table %s", stmt.TableName)
			fullRows, err := executeFromClause(stmt.From, store, stmt.Limit)
			if err != nil {
				return nil, err
			}
			log.Printf("DEBUG: Got %d rows from executeFromClause", len(fullRows))

			// Apply WHERE filtering after getting rows
			if stmt.Where != nil {
				fullRows = rhelpers.FilterRowsByWhere(fullRows, stmt.Where)
			}

			// Filter to selected columns
			var filteredRows []map[string]interface{}
			for _, row := range fullRows {
				filteredRow := make(map[string]interface{})

				log.Printf("DEBUG: Processing row with keys: %v", getKeys(row))

				for _, col := range stmt.Columns {
					log.Printf("DEBUG: Processing column - ColumnName: %q, ExpressionContext: %v", col.ColumnName, col.ExpressionContext != nil)
					if col.ColumnName == "*" {
						// copy all columns
						for k, v := range row {
							filteredRow[k] = v
						}
					} else {
						var colName string
						var value interface{}

						if col.ColumnName != "" {
							if val, ok := row[col.ColumnName]; ok {
								value = val
							} else {
								value = nil
							}
						} else if col.ExpressionContext != nil {
							v, err := rhelpers.EvaluateExpressionContext(col.ExpressionContext, row)
							if err != nil {
								return nil, err
							}
							value = v
						}

						if col.Alias != "" {
							colName = col.Alias
						} else if col.ColumnName != "" {
							colName = col.ColumnName
						} else {
							colName = "?column?"
						}

						filteredRow[colName] = value
					}
				}

				filteredRows = append(filteredRows, filteredRow)
			}

			rhelpers.SortRows(filteredRows, stmt.OrderBy)

			if stmt.Limit > 0 && len(filteredRows) > stmt.Limit {
				filteredRows = filteredRows[:stmt.Limit]
			}

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

		// No joins, execute normally
		fullRows, err := system.ExecuteSystemTableSelect(stmt, store)
		if err != nil {
			return nil, err
		}

		// Filter to selected columns
		var filteredRows []map[string]interface{}
		for _, row := range fullRows {
			filteredRow := make(map[string]interface{})

			log.Printf("DEBUG: Processing row with keys: %v", getKeys(row))

			for _, col := range stmt.Columns {
				log.Printf("DEBUG: Processing column - ColumnName: %q, ExpressionContext: %v", col.ColumnName, col.ExpressionContext != nil)
				if col.ColumnName == "*" {
					// copy all columns
					for k, v := range row {
						filteredRow[k] = v
					}
				} else {
					var colName string
					var value interface{}

					if col.ColumnName != "" {
						if val, ok := row[col.ColumnName]; ok {
							value = val
						} else {
							value = nil
						}
					} else if col.ExpressionContext != nil {
						v, err := rhelpers.EvaluateExpressionContext(col.ExpressionContext, row)
						if err != nil {
							return nil, err
						}
						value = v
					}

					if col.Alias != "" {
						colName = col.Alias
					} else if col.ColumnName != "" {
						colName = col.ColumnName
					} else {
						colName = "?column?"
					}

					filteredRow[colName] = value
				}
			}

			filteredRows = append(filteredRows, filteredRow)
		}

		rhelpers.SortRows(filteredRows, stmt.OrderBy)

		if stmt.Limit > 0 && len(filteredRows) > stmt.Limit {
			filteredRows = filteredRows[:stmt.Limit]
		}

		// columns/types: если нет строк — попытаться вывести из stmt (иначе будет пустой RowDescription)
		var cols []string
		var types []table.ColType

		// Build columns list from SELECT statement to preserve order
		for _, item := range stmt.Columns {
			if item.ColumnName == "*" {
				// For *, get all columns from first row
				if len(filteredRows) > 0 {
					for k := range filteredRows[0] {
						cols = append(cols, k)
					}
				}
			} else {
				// Use alias if available, otherwise column name or ?column?
				var colName string
				if item.Alias != "" {
					colName = item.Alias
				} else if item.ColumnName != "" {
					colName = item.ColumnName
				} else {
					colName = "?column?"
				}
				cols = append(cols, colName)
			}
		}

		// Get types from first row using the column order
		if len(filteredRows) > 0 {
			log.Printf("DEBUG: Determining types for cols: %v", cols)
			for _, colName := range cols {
				if v, ok := filteredRows[0][colName]; ok {
					// Fallback to inferring from value
					log.Printf("DEBUG: Column %q has value %v (type %T)", colName, v, v)
					switch v.(type) {
					case int, int32, int64:
						types = append(types, table.ColTypeInt)
					case float32, float64:
						types = append(types, table.ColTypeFloat)
					case bool:
						types = append(types, table.ColTypeBool)
					default:
						types = append(types, table.ColTypeString)
					}
				} else {
					log.Printf("DEBUG: Column %q not found in row", colName)
					types = append(types, table.ColTypeString)
				}
			}
			log.Printf("DEBUG: First filtered row: %+v", filteredRows[0])
			log.Printf("DEBUG: Column names: %v", cols)
			log.Printf("DEBUG: Column types: %v", types)
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

	var rows []map[string]interface{}
	var finalColumns []string
	var columnTypes []table.ColType

	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		var err error
		// Собираем список колонок для выборки
		var selectColumns []string
		if len(stmt.Columns) == 1 && stmt.Columns[0].ColumnName == "*" {
			selectColumns = nil // все колонки
		} else {
			for _, col := range stmt.Columns {
				if col.ColumnName != "" {
					selectColumns = append(selectColumns, col.ColumnName)
				}
			}
		}
		rows, finalColumns, columnTypes, err = tbl.Select(tx, stmt.TableName, selectColumns, nil, 0)
		return err
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Select fetched %d rows: %v from table %s", len(rows), rows, stmt.TableName)

	// Apply WHERE filtering after getting rows
	if stmt.Where != nil {
		rows = rhelpers.FilterRowsByWhere(rows, stmt.Where)
	}

	// Filter to selected columns
	var filteredRows []map[string]interface{}
	for _, row := range rows {
		filteredRow := make(map[string]interface{})

		// Только явно выбранные поля, даже если row содержит больше
		for _, col := range stmt.Columns {
			var colName string
			var value interface{}
			if col.ColumnName == "*" {
				// SELECT * — копируем все поля
				for k, v := range row {
					filteredRow[k] = v
				}
			} else {
				if col.ColumnName != "" {
					if val, ok := row[col.ColumnName]; ok {
						value = val
					} else {
						value = nil
					}
				} else if col.ExpressionContext != nil {
					v, err := rhelpers.EvaluateExpressionContext(col.ExpressionContext, row)
					if err != nil {
						return nil, err
					}
					value = v
				}
				if col.Alias != "" {
					colName = col.Alias
				} else if col.ColumnName != "" {
					colName = col.ColumnName
				} else {
					colName = "?column?"
				}
				filteredRow[colName] = value
			}
		}

		filteredRows = append(filteredRows, filteredRow)
	}

	rhelpers.SortRows(filteredRows, stmt.OrderBy)

	if stmt.Limit > 0 && len(filteredRows) > stmt.Limit {
		filteredRows = filteredRows[:stmt.Limit]
	}

	return map[string]interface{}{
		"rows":        filteredRows,
		"columns":     finalColumns,
		"columnTypes": columnTypes,
	}, nil
}

func executeFromClause(from *statements.FromClause, store *storage.DistributedStorageVClock, limit int) ([]map[string]interface{}, error) {
	// Get rows for main table
	stmt := &statements.SelectStatement{TableName: from.TableName}
	stmt.Limit = limit
	mainRows, err := system.ExecuteSystemTableSelect(stmt, store)
	if err != nil {
		return nil, err
	}

	log.Printf("DEBUG JOIN: fetched %d rows from main table %s", len(mainRows), from.TableName)
	if len(mainRows) > 0 {
		log.Printf("DEBUG JOIN: first main row keys: %v", getKeys(mainRows[0]))
	}

	// Add table alias prefix to main table columns
	mainAlias := from.Alias
	if mainAlias == "" {
		mainAlias = from.TableName
	}
	prefixedMainRows := prefixRowKeys(mainRows, mainAlias)

	log.Printf("DEBUG JOIN: after prefix - main table %s, alias %s, rows: %d", from.TableName, mainAlias, len(prefixedMainRows))
	if len(prefixedMainRows) > 0 {
		log.Printf("DEBUG JOIN: first prefixed row keys: %v", getKeys(prefixedMainRows[0]))
	}

	// For each join, perform hash join
	currentRows := prefixedMainRows
	for _, join := range from.Joins {
		log.Printf("DEBUG JOIN: processing join with %s (alias: %s)", join.TableName, join.Alias)
		rightStmt := &statements.SelectStatement{TableName: join.TableName}
		rightRows, err := system.ExecuteSystemTableSelect(rightStmt, store)
		if err != nil {
			return nil, err
		}

		log.Printf("DEBUG JOIN: fetched %d rows from right table %s", len(rightRows), join.TableName)

		// Add table alias prefix to right table columns
		rightAlias := join.Alias
		if rightAlias == "" {
			rightAlias = join.TableName
		}
		prefixedRightRows := prefixRowKeys(rightRows, rightAlias)

		currentRows, err = performJoin(currentRows, prefixedRightRows, join)
		if err != nil {
			return nil, err
		}

		log.Printf("DEBUG JOIN: after join, currentRows: %d", len(currentRows))
		if len(currentRows) > 0 {
			log.Printf("DEBUG JOIN: first joined row keys: %v", getKeys(currentRows[0]))
		}
	}

	return currentRows, nil
}

// prefixRowKeys adds table/alias prefix to all column names in rows
func prefixRowKeys(rows []map[string]interface{}, prefix string) []map[string]interface{} {
	prefixedRows := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		prefixedRow := make(map[string]interface{})
		for key, value := range row {
			// Add both prefixed and unprefixed versions
			prefixedRow[prefix+"."+key] = value
			prefixedRow[key] = value // Keep unprefixed for backward compatibility
		}
		prefixedRows[i] = prefixedRow
	}
	return prefixedRows
}

func performJoin(leftRows, rightRows []map[string]interface{}, join statements.JoinClause) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	switch join.Type {
	case "LEFT":
		for _, leftRow := range leftRows {
			found := false
			for _, rightRow := range rightRows {
				combined := make(map[string]interface{})
				// Merge left and right rows (both already have prefixed keys)
				for k, v := range leftRow {
					combined[k] = v
				}
				for k, v := range rightRow {
					combined[k] = v
				}
				// Check ON condition
				if join.OnCondition != nil {
					if expr, ok := join.OnCondition.(*parser.ExpressionContext); ok {
						val, err := rhelpers.EvaluateExpressionContext(expr, combined)
						if err != nil {
							return nil, err
						}
						if val == nil || val == false || val == 0 {
							continue
						}
					}
				}
				result = append(result, combined)
				found = true
			}
			if !found {
				// Add leftRow with nulls for right columns
				combined := make(map[string]interface{})
				for k, v := range leftRow {
					combined[k] = v
				}
				// For right columns, set to nil
				for _, rightRow := range rightRows {
					for k := range rightRow {
						if _, ok := combined[k]; !ok {
							combined[k] = nil
						}
					}
					break // Only need column names from one row
				}
				result = append(result, combined)
			}
		}
	case "INNER", "":
		// Default to INNER
		for _, leftRow := range leftRows {
			for _, rightRow := range rightRows {
				combined := make(map[string]interface{})
				// Merge left and right rows (both already have prefixed keys)
				for k, v := range leftRow {
					combined[k] = v
				}
				for k, v := range rightRow {
					combined[k] = v
				}
				// Check ON condition
				if join.OnCondition != nil {
					if expr, ok := join.OnCondition.(*parser.ExpressionContext); ok {
						val, err := rhelpers.EvaluateExpressionContext(expr, combined)
						if err != nil {
							return nil, err
						}
						if val == nil || val == false || val == 0 {
							continue
						}
					}
				}
				result = append(result, combined)
			}
		}
	default:
		// For now, treat as INNER
		return performJoin(leftRows, rightRows, statements.JoinClause{Type: "INNER", TableName: join.TableName, Alias: join.Alias, OnCondition: join.OnCondition})
	}
	return result, nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
