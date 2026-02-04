package executor

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/revaluate"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
	"strings"

	"go.uber.org/zap"
)

func ExecuteSelect(
	stmt *statements.SelectStatement,
	store *storage.DistributedStorageVClock,
	exCtx ExecutorContext,
) (*table.ExecuteResult, error) {
	logger.Debug("execute select", zap.Any("stmt", stmt))

	// --- 1) SELECT без таблицы (функции/выражения) ---
	if stmt.From == nil {
		return ExecuteSelectWithoutTable(stmt, store, exCtx)
	} else {
		// --- 3) Normal table select ---
		return ExecuteNormalTable(stmt, store, exCtx)
	}
}

func ExecuteSelectWithoutTable(
	stmt *statements.SelectStatement,
	store *storage.DistributedStorageVClock,
	exCtx ExecutorContext,
) (*table.ExecuteResult, error) {
	resultColumns := make([]table.TableColumn, 0, len(stmt.Columns))
	row := make([]interface{}, 0, len(stmt.Columns))
	size := len(stmt.Columns)
	for i, col := range stmt.Columns {
		var value *table.ExecuteResult

		logger.Debug("Select without table - processing column", zap.Int("index", i), zap.Any("column", col))

		if col.Function != nil {
			// Evaluate function args
			args := make([]interface{}, 0, len(col.Function.Args))
			for _, argExpr := range col.Function.Args {
				if argExpr.STAR() != nil {
					// Special case for COUNT(*)
					args = append(args, 1)
					continue
				}

				if argExpr.Expression() == nil {
					return nil, fmt.Errorf("invalid function argument")
				}
				val, err := revaluate.EvaluateExpressionContext(argExpr.Expression(), nil)
				if err != nil {
					return nil, fmt.Errorf("error evaluating function arg: %w", err)
				}
				if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
					args = append(args, valExpr.Row.Rows[0][0])
				} else {
					return nil, fmt.Errorf("invalid function arg")
				}
			}
			v, err := functions.ExecuteFunction(col.Function.Name, args)
			if err != nil {
				return nil, err
			}
			if len(v.Columns) > 1 && size > 1 {
				return nil, fmt.Errorf("function %s returns multiple columns, cannot be used in multi-column select", col.Function.Name)
			} else if len(v.Columns) > 1 && size == 1 {
				return v, nil
			}
			value = v
		} else if col.ExpressionContext != nil {
			v, err := revaluate.EvaluateExpressionContext(col.ExpressionContext, nil)
			if err != nil {
				return nil, err
			}
			if execRes, ok := v.(*rmodels.ResultRowsExpression); ok && len(execRes.Row.Rows) > 0 && len(execRes.Row.Rows[0]) > 0 {
				value = execRes.Row
			} else {
				return nil, fmt.Errorf("expected ExecuteResult from expression, got %T", v)
			}
		} else {
			return nil, fmt.Errorf("unsupported select column without table")
		}

		if col.Alias != "" {
			value.Columns[0].Name = col.Alias
		}

		row = append(row, value.Rows[0][0])
		resultColumns = append(resultColumns, value.Columns[0])

		logger.Debugf("Select without table - processed column %d/%d: %v", i+1, size, row)
	}
	if len(row) != len(stmt.Columns) {
		return nil, fmt.Errorf("mismatch in number of columns and row values")
	}

	execResult := &table.ExecuteResult{
		Rows:    [][]interface{}{row},
		Columns: resultColumns,
	}

	return execResult, nil
}

func ExecuteNormalTable(
	stmt *statements.SelectStatement,
	store *storage.DistributedStorageVClock,
	exCtx ExecutorContext,
) (*table.ExecuteResult, error) {

	tableName := stmt.From.TableName

	// Проверяем есть ли JOIN
	if stmt.From != nil && len(stmt.From.Joins) > 0 {
		return ExecuteSelectWithJoins(stmt, store, exCtx)
	}

	// --- 3) Normal table select ---
	tbl := &table.Table{
		Storage:  store,
		Database: exCtx.Database,
		Schema:   exCtx.Schema,
		Name:     tableName,
	}

	logger.Debug("Executing normal SELECT on table", zap.String("table", tableName))

	var result *table.ExecuteResult

	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		var err error
		// Проверяем есть ли только простые колонки без выражений
		hasExpressions := false
		for _, col := range stmt.Columns {
			if col.ExpressionContext != nil || col.Function != nil {
				hasExpressions = true
				break
			}
		}

		// Если есть WHERE или выражения, получаем ВСЕ колонки
		// потому что WHERE может ссылаться на колонки не в SELECT
		if stmt.Where != nil || hasExpressions {
			selectColumns := []table.SelectColumn{{IsAll: true}}
			logger.Debugf("Executing SELECT with WHERE/expressions on table %s - fetching all columns", tableName)
			result, err = tbl.Select(tx, tableName, selectColumns, nil, 0)
		} else {
			// Простой SELECT без WHERE - собираем список колонок для выборки
			var selectColumns []table.SelectColumn
			for _, col := range stmt.Columns {
				if col.IsSelectAll {
					selectColumns = append(selectColumns, table.SelectColumn{
						IsAll: true,
					})
				} else {
					selectColumns = append(selectColumns, table.SelectColumn{
						Name: col.ColumnName,
					})
				}
			}
			logger.Debugf("Executing SELECT on table %s with columns: %+v", tableName, selectColumns)
			result, err = tbl.Select(tx, tableName, selectColumns, nil, 0)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	logger.Debugf("Select fetched %d rows: %v from table %s", len(result.Rows), result.Rows, tableName)
	logger.Debug("Select fetched", zap.Any("columns", result.Columns), zap.Any("rows", result.Rows))

	// Apply WHERE filtering after getting rows
	if stmt.Where != nil {
		result = revaluate.EvaluateFilterRowsByWhere(result, stmt.Where)
	}

	// Check for aggregates
	hasAggregates := false
	for _, col := range stmt.Columns {
		if col.Function != nil && col.Function.IsAggregate {
			hasAggregates = true
			break
		}
	}

	if len(stmt.GroupBy) > 0 || hasAggregates {
		// Handle GROUP BY and aggregates
		result, err = processGroupByAndAggregates(stmt, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	// Normal processing
	// Проверяем нужно ли обрабатывать выражения или выбирать конкретные колонки
	hasExpressions := false
	needColumnFiltering := false

	for _, col := range stmt.Columns {
		if col.ExpressionContext != nil || col.Function != nil {
			hasExpressions = true
			break
		}
		if !col.IsSelectAll {
			needColumnFiltering = true
		}
	}

	if hasExpressions {
		// Если есть выражения, обрабатываем их
		result, err = processSelectExpressions(stmt, result)
		if err != nil {
			return nil, err
		}
	} else if needColumnFiltering && stmt.Where != nil {
		// Если был WHERE и нужны только конкретные колонки, фильтруем
		result, err = filterSelectColumns(stmt, result)
		if err != nil {
			return nil, err
		}
	} else {
		// Просто применяем алиасы к обычным колонкам
		for i, col := range stmt.Columns {
			if col.Alias != "" && i < len(result.Columns) {
				result.Columns[i].Name = col.Alias
			}
		}
	}

	revaluate.EvaluateSortRows(result, stmt.OrderBy)

	if stmt.Limit > 0 && len(result.Rows) > stmt.Limit {
		result.Rows = result.Rows[:stmt.Limit]
	}

	return result, nil
}

// filterSelectColumns фильтрует колонки из полного набора данных
func filterSelectColumns(stmt *statements.SelectStatement, inputResult *table.ExecuteResult) (*table.ExecuteResult, error) {
	newColumns := make([]table.TableColumn, 0)
	columnIndices := make([]int, 0)

	for _, col := range stmt.Columns {
		if col.IsSelectAll {
			// Для * берем все колонки
			return inputResult, nil
		}

		// Ищем колонку по имени с учётом table alias
		found := false
		searchName := col.ColumnName
		if col.TableAlias != "" {
			searchName = col.TableAlias + "." + col.ColumnName
		}

		for i, origCol := range inputResult.Columns {
			// Проверяем точное совпадение или совпадение без префикса
			if origCol.Name == searchName || origCol.Name == col.ColumnName {
				newCol := origCol
				if col.Alias != "" {
					newCol.Name = col.Alias
				} else {
					newCol.Name = col.ColumnName
				}
				newCol.Idx = len(newColumns)
				newColumns = append(newColumns, newCol)
				columnIndices = append(columnIndices, i)
				found = true
				break
			}
			// Также проверяем если origCol имеет префикс, а мы ищем без него
			if col.TableAlias == "" && strings.HasSuffix(origCol.Name, "."+col.ColumnName) {
				newCol := origCol
				if col.Alias != "" {
					newCol.Name = col.Alias
				} else {
					newCol.Name = col.ColumnName
				}
				newCol.Idx = len(newColumns)
				newColumns = append(newColumns, newCol)
				columnIndices = append(columnIndices, i)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("column %s not found", col.ColumnName)
		}
	}

	// Фильтруем строки, оставляя только нужные колонки
	newRows := make([][]interface{}, len(inputResult.Rows))
	for i, row := range inputResult.Rows {
		newRow := make([]interface{}, len(columnIndices))
		for j, idx := range columnIndices {
			if idx < len(row) {
				newRow[j] = row[idx]
			}
		}
		newRows[i] = newRow
	}

	return &table.ExecuteResult{
		Rows:    newRows,
		Columns: newColumns,
	}, nil
}

// processSelectExpressions обрабатывает выражения в SELECT (включая CASE, функции, алиасы)
func processSelectExpressions(stmt *statements.SelectStatement, inputResult *table.ExecuteResult) (*table.ExecuteResult, error) {
	newColumns := make([]table.TableColumn, 0, len(stmt.Columns))
	newRows := make([][]interface{}, 0, len(inputResult.Rows))

	// Сначала определяем колонки
	for colIdx, col := range stmt.Columns {
		var colName string
		var colType table.ColType

		if col.IsSelectAll {
			// Для * добавляем все колонки из исходного результата
			for _, origCol := range inputResult.Columns {
				newColumns = append(newColumns, origCol)
			}
			break // * заменяет все
		}

		// Обработка выражений
		if col.ExpressionContext != nil {
			// Для выражений, тип определим позже, но пока string
			colType = table.ColTypeString
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = "?column?"
			}
		} else if col.Function != nil {
			// Для функций, тип определим позже
			colType = table.ColTypeString
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = "?column?"
			}
		} else if col.ColumnName != "" {
			// Простая колонка
			searchName := col.ColumnName
			if col.TableAlias != "" {
				searchName = col.TableAlias + "." + col.ColumnName
			}
			found := false
			for _, origCol := range inputResult.Columns {
				if origCol.Name == searchName || origCol.Name == col.ColumnName {
					colType = origCol.Type
					found = true
					break
				}
				if col.TableAlias == "" && strings.HasSuffix(origCol.Name, "."+col.ColumnName) {
					colType = origCol.Type
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("column %s not found", col.ColumnName)
			}
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = searchName
			}
		}

		newColumns = append(newColumns, table.TableColumn{
			Idx:  colIdx,
			Name: colName,
			Type: colType,
		})
	}

	// Теперь обрабатываем строки
	for _, row := range inputResult.Rows {
		newRow := make([]interface{}, 0, len(stmt.Columns))

		resultRow := &table.ResultRow{
			Row:     row,
			Columns: inputResult.Columns,
		}

		for _, col := range stmt.Columns {
			var colValue interface{}

			if col.IsSelectAll {
				// Для * добавляем все значения
				newRow = append(newRow, row...)
				break
			}

			// Обработка выражений
			if col.ExpressionContext != nil {
				val, err := revaluate.EvaluateExpressionContext(col.ExpressionContext, resultRow)
				if err != nil {
					return nil, fmt.Errorf("error evaluating expression: %w", err)
				}
				if boolVal, ok := val.(*rmodels.BoolExpression); ok {
					colValue = boolVal.Value
				} else if resultVal, ok := val.(*rmodels.ResultRowsExpression); ok {
					if len(resultVal.Row.Rows) > 0 && len(resultVal.Row.Rows[0]) > 0 {
						colValue = resultVal.Row.Rows[0][0]
					} else {
						colValue = nil
					}
				} else {
					colValue = nil
				}
			} else if col.Function != nil {
				// Обработка функций
				args := make([]interface{}, 0, len(col.Function.Args))
				for _, argExpr := range col.Function.Args {
					if argExpr.STAR() != nil {
						// Special case for COUNT(*)
						args = append(args, 1)
						continue
					}
					if argExpr.Expression() == nil {
						return nil, fmt.Errorf("invalid function argument")
					}
					val, err := revaluate.EvaluateExpressionContext(argExpr.Expression(), resultRow)
					if err != nil {
						return nil, fmt.Errorf("error evaluating function arg: %w", err)
					}
					if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
						args = append(args, valExpr.Row.Rows[0][0])
					} else {
						return nil, fmt.Errorf("invalid function arg")
					}
				}
				v, err := functions.ExecuteFunction(col.Function.Name, args)
				if err != nil {
					return nil, fmt.Errorf("error executing function %s: %w", col.Function.Name, err)
				}
				if len(v.Rows) > 0 && len(v.Rows[0]) > 0 {
					colValue = v.Rows[0][0]
				} else {
					colValue = nil
				}
			} else if col.ColumnName != "" {
				// Простая колонка
				searchName := col.ColumnName
				if col.TableAlias != "" {
					searchName = col.TableAlias + "." + col.ColumnName
				}
				found := false
				for i, origCol := range inputResult.Columns {
					if origCol.Name == searchName || origCol.Name == col.ColumnName {
						colValue = row[i]
						found = true
						break
					}
					if col.TableAlias == "" && strings.HasSuffix(origCol.Name, "."+col.ColumnName) {
						colValue = row[i]
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("column %s not found", col.ColumnName)
				}
			}

			newRow = append(newRow, colValue)
		}

		newRows = append(newRows, newRow)
	}

	return &table.ExecuteResult{
		Rows:    newRows,
		Columns: newColumns,
	}, nil
}

// inferTypeFromValue определяет тип колонки по значению
func inferTypeFromValue(val interface{}) table.ColType {
	switch val.(type) {
	case int, int32, int64:
		return table.ColTypeInt
	case float32, float64:
		return table.ColTypeFloat
	case bool:
		return table.ColTypeBool
	default:
		return table.ColTypeString
	}
}

// ExecuteSelectWithJoins выполняет SELECT с JOIN операциями
func ExecuteSelectWithJoins(
	stmt *statements.SelectStatement,
	store *storage.DistributedStorageVClock,
	exCtx ExecutorContext,
) (*table.ExecuteResult, error) {
	var mainResult *table.ExecuteResult

	// Получаем данные из главной таблицы
	mainTable := &table.Table{
		Storage:  store,
		Database: exCtx.Database,
		Schema:   exCtx.Schema,
		Name:     stmt.From.TableName,
	}

	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		var err error
		selectColumns := []table.SelectColumn{{IsAll: true}}
		mainResult, err = mainTable.Select(tx, stmt.From.TableName, selectColumns, nil, 0)
		return err
	})
	if err != nil {
		return nil, err
	}

	// Добавляем префикс к колонкам главной таблицы
	mainAlias := stmt.From.Alias
	if mainAlias == "" {
		mainAlias = stmt.From.TableName
	}
	for i := range mainResult.Columns {
		originalName := mainResult.Columns[i].Name
		mainResult.Columns[i].Name = mainAlias + "." + originalName
		mainResult.Columns[i].TableIdentifier = mainAlias
		mainResult.Columns[i].OriginalTableName = stmt.From.TableName
	}

	currentResult := mainResult

	// Выполняем JOIN-ы последовательно
	for _, join := range stmt.From.Joins {
		var rightResult *table.ExecuteResult

		rightTable := &table.Table{
			Storage:  store,
			Database: exCtx.Database,
			Schema:   exCtx.Schema,
			Name:     join.TableName,
		}

		err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			var err error
			selectColumns := []table.SelectColumn{{IsAll: true}}
			rightResult, err = rightTable.Select(tx, join.TableName, selectColumns, nil, 0)
			return err
		})
		if err != nil {
			return nil, err
		}

		// Добавляем префикс к колонкам правой таблицы
		rightAlias := join.Alias
		if rightAlias == "" {
			rightAlias = join.TableName
		}
		for i := range rightResult.Columns {
			originalName := rightResult.Columns[i].Name
			rightResult.Columns[i].Name = rightAlias + "." + originalName
			rightResult.Columns[i].TableIdentifier = rightAlias
			rightResult.Columns[i].OriginalTableName = join.TableName
		}

		// Выполняем JOIN
		currentResult, err = performJoin(currentResult, rightResult, join)
		if err != nil {
			return nil, err
		}
	}

	// Применяем WHERE
	if stmt.Where != nil {
		currentResult = revaluate.EvaluateFilterRowsByWhere(currentResult, stmt.Where)
	}

	// Check for aggregates
	hasAggregates := false
	for _, col := range stmt.Columns {
		if col.Function != nil && col.Function.IsAggregate {
			hasAggregates = true
			break
		}
	}

	if len(stmt.GroupBy) > 0 || hasAggregates {
		// Handle GROUP BY and aggregates
		currentResult, err = processGroupByAndAggregates(stmt, currentResult)
		if err != nil {
			return nil, err
		}
		// Применяем ORDER BY после GROUP BY
		revaluate.EvaluateSortRows(currentResult, stmt.OrderBy)
	} else {
		// Применяем ORDER BY перед SELECT
		revaluate.EvaluateSortRows(currentResult, stmt.OrderBy)
		// Обрабатываем SELECT колонки
		currentResult, err = processSelectExpressions(stmt, currentResult)
		if err != nil {
			return nil, err
		}
	}

	// Применяем LIMIT
	if stmt.Limit > 0 && len(currentResult.Rows) > stmt.Limit {
		currentResult.Rows = currentResult.Rows[:stmt.Limit]
	}

	return currentResult, nil
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

func performJoin(leftRows, rightRows *table.ExecuteResult, join statements.JoinClause) (*table.ExecuteResult, error) {
	// Собираем итоговые колонки: left + right, с корректными Idx
	combColumns := make([]table.TableColumn, 0, len(leftRows.Columns)+len(rightRows.Columns))

	for i, c := range leftRows.Columns {
		cc := c
		cc.Idx = i
		combColumns = append(combColumns, cc)
	}
	for i, c := range rightRows.Columns {
		cc := c
		cc.Idx = len(leftRows.Columns) + i
		combColumns = append(combColumns, cc)
	}

	rightPad := make([]interface{}, len(rightRows.Columns)) // nil-ы по умолчанию

	evalOn := func(combRow []interface{}) (bool, error) {
		if join.OnCondition == nil {
			return true, nil
		}
		expr, ok := join.OnCondition.(*parser.ExpressionContext)
		if !ok {
			return false, fmt.Errorf("unsupported OnCondition type: %T", join.OnCondition)
		}
		val, err := revaluate.EvaluateExpressionContext(expr, &table.ResultRow{
			Row:     combRow,
			Columns: combColumns,
		})
		if err != nil {
			logger.Errorf("Error evaluating JOIN ON condition: %v", err)
			return false, err
		}

		if valBool, ok := val.(*rmodels.BoolExpression); ok {
			return valBool.Value, nil
		} else if resultVal, ok := val.(*rmodels.ResultRowsExpression); ok {
			// Пробуем извлечь булево значение из результата
			if len(resultVal.Row.Rows) > 0 && len(resultVal.Row.Rows[0]) > 0 {
				if boolVal, ok := resultVal.Row.Rows[0][0].(bool); ok {
					return boolVal, nil
				}
			}
			logger.Errorf("ON condition returned ResultRowsExpression with non-bool value: %v", resultVal)
			return false, fmt.Errorf("ON condition did not evaluate to bool, got %T with value %v", val, resultVal.Row.Rows)
		} else {
			logger.Errorf("ON condition returned unexpected type: %T", val)
			return false, fmt.Errorf("OnCondition did not evaluate to bool, got %T", val)
		}
	}

	resultRows := make([][]interface{}, 0)

	switch join.Type {
	case "LEFT":
		for _, lrow := range leftRows.Rows {
			found := false

			for _, rrow := range rightRows.Rows {
				combRow := make([]interface{}, 0, len(lrow)+len(rrow))
				combRow = append(combRow, lrow...)
				combRow = append(combRow, rrow...)

				ok, err := evalOn(combRow)
				if err != nil {
					return nil, err
				}
				if !ok {
					continue
				}

				resultRows = append(resultRows, combRow)
				found = true
			}

			if !found {
				// Левую строку оставляем, правую часть заполняем nil-ами
				combRow := make([]interface{}, 0, len(lrow)+len(rightPad))
				combRow = append(combRow, lrow...)
				combRow = append(combRow, rightPad...)
				resultRows = append(resultRows, combRow)
			}
		}

	case "INNER", "":
		fallthrough
	default:
		// Любой другой тип пока трактуем как INNER
		for _, lrow := range leftRows.Rows {
			for _, rrow := range rightRows.Rows {
				combRow := make([]interface{}, 0, len(lrow)+len(rrow))
				combRow = append(combRow, lrow...)
				combRow = append(combRow, rrow...)

				ok, err := evalOn(combRow)
				if err != nil {
					logger.Errorf("JOIN ON error: %v. Skipping this row combination", err)
					continue
				}
				if !ok {
					continue
				}

				resultRows = append(resultRows, combRow)
			}
		}
	}

	return &table.ExecuteResult{
		Rows:    resultRows,
		Columns: combColumns,
	}, nil
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func processGroupByAndAggregates(stmt *statements.SelectStatement, inputResult *table.ExecuteResult) (*table.ExecuteResult, error) {
	// Check for invalid combinations
	hasAggregates := false
	for _, col := range stmt.Columns {
		if col.Function != nil && col.Function.IsAggregate {
			hasAggregates = true
			break
		}
	}

	if len(stmt.GroupBy) == 0 && hasAggregates {
		// If aggregates without GROUP BY, no non-aggregate columns allowed
		for _, col := range stmt.Columns {
			if col.ColumnName != "" || col.ExpressionContext != nil {
				return nil, fmt.Errorf("column must appear in the GROUP BY clause or be used in an aggregate function")
			}
		}
	}

	var groups map[string][]*table.ResultRow

	// If GROUP BY exists, validate non-aggregate columns are in GROUP BY
	if len(stmt.GroupBy) > 0 {
		for _, col := range stmt.Columns {
			if col.ColumnName != "" {
				found := false
				for _, gb := range stmt.GroupBy {
					if gb.ColumnName == col.ColumnName && gb.TableAlias == col.TableAlias {
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("column %s must appear in GROUP BY clause", col.ColumnName)
				}
			}
		}

		// If GROUP BY exists, group rows
		groups = make(map[string][]*table.ResultRow)
		for _, row := range inputResult.Rows {
			resultRow := &table.ResultRow{
				Row:     row,
				Columns: inputResult.Columns,
			}
			groupKey, err := computeGroupKey(stmt.GroupBy, resultRow)
			if err != nil {
				return nil, err
			}
			groups[groupKey] = append(groups[groupKey], resultRow)
		}
	} else {
		// No GROUP BY, treat all rows as one group
		groupRows := make([]*table.ResultRow, len(inputResult.Rows))
		for i, row := range inputResult.Rows {
			groupRows[i] = &table.ResultRow{
				Row:     row,
				Columns: inputResult.Columns,
			}
		}
		groups = map[string][]*table.ResultRow{"": groupRows}
	}

	// Now process each group
	newColumns := make([]table.TableColumn, 0, len(stmt.Columns))
	newRows := make([][]interface{}, 0, len(groups))

	// Define columns
	for colIdx, col := range stmt.Columns {
		var colName string
		var colType table.ColType

		if col.IsSelectAll {
			return nil, fmt.Errorf("SELECT * not allowed with GROUP BY or aggregates")
		}

		if col.ExpressionContext != nil {
			colType = table.ColTypeString
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = "?column?"
			}
		} else if col.Function != nil {
			colType = table.ColTypeString
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = col.Function.Name
			}
		} else if col.ColumnName != "" {
			// For non-aggregate columns, they must be in GROUP BY
			if len(stmt.GroupBy) > 0 {
				found := false
				for _, gb := range stmt.GroupBy {
					if gb.ColumnName == col.ColumnName {
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("column %s must appear in GROUP BY clause", col.ColumnName)
				}
			}
			// Get actual type from input columns
			colType = table.ColTypeString // default
			for _, inputCol := range inputResult.Columns {
				if inputCol.Name == col.ColumnName || (col.TableAlias != "" && inputCol.Name == col.TableAlias+"."+col.ColumnName) {
					colType = inputCol.Type
					break
				}
			}
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = col.ColumnName
			}
		}

		newColumns = append(newColumns, table.TableColumn{
			Idx:  colIdx,
			Name: colName,
			Type: colType,
		})
	}

	// Process each group
	for _, groupRows := range groups {
		newRow := make([]interface{}, 0, len(stmt.Columns))

		for idx, col := range stmt.Columns {
			var colValue interface{}

			if col.ExpressionContext != nil {
				// For now, assume no expressions in aggregates
				return nil, fmt.Errorf("expressions not supported in GROUP BY yet")
			} else if col.Function != nil {
				if col.Function.IsAggregate {
					// Compute aggregate over group
					args := make([]interface{}, 0, len(col.Function.Args))
					for _, argExpr := range col.Function.Args {
						// For aggregates, collect values from all rows in group
						groupValues := make([]interface{}, 0, len(groupRows))
						for _, groupRow := range groupRows {
							if argExpr.STAR() != nil {
								// Special case for COUNT(*)
								groupValues = append(groupValues, 1)
								continue
							}

							if argExpr.Expression() == nil {
								return nil, fmt.Errorf("invalid aggregate argument")
							}
							val, err := revaluate.EvaluateExpressionContext(argExpr.Expression(), groupRow)
							if err != nil {
								return nil, fmt.Errorf("error evaluating aggregate arg: %w", err)
							}
							if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
								groupValues = append(groupValues, valExpr.Row.Rows[0][0])
							} else {
								groupValues = append(groupValues, nil)
							}
						}
						args = append(args, groupValues)
					}
					v, err := functions.ExecuteAggregateFunction(col.Function.Name, args)
					if err != nil {
						return nil, fmt.Errorf("error executing aggregate function %s: %w", col.Function.Name, err)
					}
					// TODO посмотреть, может можно лучше сделать
					colValue = v.Rows[0][0]
					newColumns[idx].Type = v.Columns[0].Type
				} else {
					// Non-aggregate function, evaluate on first row or something? For now error
					return nil, fmt.Errorf("non-aggregate functions not supported in GROUP BY")
				}
			} else if col.ColumnName != "" {
				// Non-aggregate column, take from first row in group
				if len(groupRows) > 0 {
					resultRow := groupRows[0]
					searchName := col.ColumnName
					if col.TableAlias != "" {
						searchName = col.TableAlias + "." + col.ColumnName
					}
					found := false
					for i, origCol := range inputResult.Columns {
						if origCol.Name == searchName || origCol.Name == col.ColumnName {
							colValue = resultRow.Row[i]
							found = true
							break
						}
						if col.TableAlias == "" && strings.HasSuffix(origCol.Name, "."+col.ColumnName) {
							colValue = resultRow.Row[i]
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("column %s not found", col.ColumnName)
					}
				}
			}

			newRow = append(newRow, colValue)
		}

		newRows = append(newRows, newRow)
	}

	return &table.ExecuteResult{
		Columns: newColumns,
		Rows:    newRows,
	}, nil
}
func computeGroupKey(groupBy []items.SelectItem, row *table.ResultRow) (string, error) {
	keyParts := make([]string, 0, len(groupBy))
	for _, gb := range groupBy {
		var val interface{}
		var err error
		if gb.ExpressionContext != nil {
			val, err = revaluate.EvaluateExpressionContext(gb.ExpressionContext, row)
			if err != nil {
				return "", err
			}
			if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
				val = valExpr.Row.Rows[0][0]
			} else {
				val = nil
			}
		} else if gb.ColumnName != "" {
			// Find column value
			searchName := gb.ColumnName
			if gb.TableAlias != "" {
				searchName = gb.TableAlias + "." + gb.ColumnName
			}
			found := false
			for i, col := range row.Columns {
				if col.Name == searchName || col.Name == gb.ColumnName {
					val = row.Row[i]
					found = true
					break
				} else if gb.TableAlias == "" && strings.HasSuffix(col.Name, "."+gb.ColumnName) {
					val = row.Row[i]
					found = true
					break
				}
			}
			if !found {
				return "", fmt.Errorf("group by column %s not found", searchName)
			}
		} else {
			return "", fmt.Errorf("invalid group by item")
		}
		keyParts = append(keyParts, fmt.Sprintf("%v", val))
	}
	return strings.Join(keyParts, "|"), nil
}
