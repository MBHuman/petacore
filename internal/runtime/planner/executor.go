package planner

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/revaluate"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ExecutorContext содержит контекст выполнения
type ExecutorContext struct {
	Database         string
	Schema           string
	Storage          *storage.DistributedStorageVClock
	SubqueryExecutor subquery.SubqueryExecutor
	// QueryStartTime   int64 // timestamp in microseconds for NOW() function
	GoCtx context.Context
}

// ExecutePlan выполняет план запроса
func ExecutePlan(plan *QueryPlan, ctx ExecutorContext) (*table.ExecuteResult, error) {
	var result *table.ExecuteResult
	var err error
	runtimeParams := make(map[int]interface{})

	// // Set query start time if not already set
	// if ctx.QueryStartTime == 0 {
	// 	ctx.QueryStartTime = time.Now().UnixMicro()
	// }
	if ctx.GoCtx == nil {
		ctx.GoCtx = context.Background()
	}
	ctx.GoCtx = context.WithValue(ctx.GoCtx, "queryStartTime", time.Now().UnixMicro())

	// Выполняем в транзакции
	err = ctx.Storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		// Attach subquery executor to the plan so expressions can execute subqueries
		plan.SubqueryExecutor = func(stmt *statements.SelectStatement) (*table.ExecuteResult, error) {
			subPlan, err := CreateQueryPlan(stmt, PlannerContext{Database: ctx.Database, Schema: ctx.Schema})
			if err != nil {
				return nil, err
			}
			return executePlanNode(subPlan.Root, plan, ctx, tx, runtimeParams)
		}
		ctx.SubqueryExecutor = plan.SubqueryExecutor

		result, err = executePlanNode(plan.Root, plan, ctx, tx, runtimeParams)
		return err
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// executePlanNode рекурсивно выполняет узел плана
func executePlanNode(
	node PlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	logger.Debug("executing plan node", zap.String("type", node.NodeType()))

	switch n := node.(type) {
	case *InitPlanNode:
		return executeInitPlanNode(n, plan, ctx, tx, runtimeParams)
	case *ScalarSubqueryPlanNode:
		return executeScalarSubqueryPlanNode(n, plan, ctx, tx, runtimeParams)
	case *ScanPlanNode:
		return executeScan(n, plan, ctx, tx, runtimeParams)
	case *ValuesPlanNode:
		result, err := executeValues(n, plan, ctx, tx, runtimeParams)
		logger.Debug("ValuePlanNode result", zap.Any("result", result))
		return result, err
	case *ProjectPlanNode:
		return executeProject(n, plan, ctx, tx, runtimeParams)
	case *FilterPlanNode:
		return executeFilter(n, plan, ctx, tx, runtimeParams)
	case *JoinPlanNode:
		return executeJoin(n, plan, ctx, tx, runtimeParams)
	case *AggregatePlanNode:
		return executeAggregate(n, plan, ctx, tx, runtimeParams)
	case *SortPlanNode:
		return executeSort(n, plan, ctx, tx, runtimeParams)
	case *LimitPlanNode:
		return executeLimit(n, plan, ctx, tx, runtimeParams)
	case *UnionPlanNode:
		return executeUnion(n, plan, ctx, tx, runtimeParams)
	case *IntersectPlanNode:
		return executeIntersect(n, plan, ctx, tx, runtimeParams)
	case *ExceptPlanNode:
		return executeExcept(n, plan, ctx, tx, runtimeParams)
	case *SubqueryPlanNode:
		return executeSubquery(n, plan, ctx, tx, runtimeParams)
	default:
		return nil, fmt.Errorf("unknown plan node type: %s", node.NodeType())
	}
}

// executeScan выполняет сканирование таблицы
func executeScan(
	node *ScanPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	tbl := &table.Table{
		Storage:  ctx.Storage,
		Database: node.Database,
		Schema:   node.Schema,
		Name:     node.TableName,
	}

	// Читаем все строки и все колонки
	selectColumns := []table.SelectColumn{{IsAll: true}}
	result, err := tbl.Select(tx, node.TableName, selectColumns, nil, 0)
	if err != nil {
		return nil, err
	}

	// Устанавливаем TableIdentifier для всех колонок (алиас или имя таблицы)
	tableIdentifier := node.Alias
	if tableIdentifier == "" {
		tableIdentifier = node.TableName
	}
	for i := range result.Columns {
		result.Columns[i].TableIdentifier = tableIdentifier
		result.Columns[i].OriginalTableName = node.TableName
	}

	logger.Debug("scan result", zap.Int("rows", len(result.Rows)), zap.Int("cols", len(result.Columns)))
	return result, nil
}

// executeValues выполняет VALUES узел (SELECT без таблицы)
// TODO тут могут быть вложенные запросы на expression с подзапросами, надо это учесть
func executeValues(
	node *ValuesPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	resultColumns := make([]table.TableColumn, 0, len(node.Columns))
	row := make([]interface{}, 0, len(node.Columns))

	for i, col := range node.Columns {
		var value *table.ExecuteResult

		if col.Function != nil {
			// Вычисляем аргументы функции
			args := make([]interface{}, 0, len(col.Function.Args))
			for _, argExpr := range col.Function.Args {
				if argExpr.STAR() != nil {
					args = append(args, 1)
					continue
				}

				if argExpr.Expression() == nil {
					return nil, fmt.Errorf("invalid function argument")
				}
				val, err := revaluate.EvaluateExpressionContext(ctx.GoCtx, argExpr.Expression(), nil, ctx.SubqueryExecutor, runtimeParams)
				if err != nil {
					return nil, fmt.Errorf("error evaluating function arg: %w", err)
				}
				if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
					args = append(args, valExpr.Row.Rows[0][0])
				} else {
					return nil, fmt.Errorf("invalid function arg")
				}
			}
			v, err := functions.ExecuteFunctionWithContext(
				ctx.GoCtx,
				col.Function.Name,
				args,
			)
			if err != nil {
				return nil, err
			}
			value = v
		} else if col.ExpressionContext != nil {
			v, err := revaluate.EvaluateExpressionContext(ctx.GoCtx, col.ExpressionContext, nil, ctx.SubqueryExecutor, runtimeParams)
			if err != nil {
				return nil, err
			}
			if execRes, ok := v.(*rmodels.ResultRowsExpression); ok && len(execRes.Row.Rows) > 0 && len(execRes.Row.Rows[0]) > 0 {
				value = execRes.Row
			} else if boolExpr, ok := v.(*rmodels.BoolExpression); ok {
				// Convert BoolExpression to ExecuteResult
				value = &table.ExecuteResult{
					Rows:    [][]interface{}{{boolExpr.Value}},
					Columns: []table.TableColumn{{Idx: 0, Name: "?column?", Type: table.ColTypeBool}},
				}
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

		logger.Debugf("Values node - processed column %d/%d", i+1, len(node.Columns))
	}

	return &table.ExecuteResult{
		Rows:    [][]interface{}{row},
		Columns: resultColumns,
	}, nil
}

// executeProject выполняет проекцию (выбор колонок)
func executeProject(
	node *ProjectPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	// Сначала выполняем дочерний узел
	inputResult, err := executePlanNode(node.Input, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	// Если все колонки - SELECT *, возвращаем как есть
	allSelectAll := true
	for _, col := range node.Columns {
		if !col.IsSelectAll {
			allSelectAll = false
			break
		}
	}
	if allSelectAll {
		return inputResult, nil
	}

	// Обрабатываем каждую колонку
	newColumns := make([]table.TableColumn, 0, len(node.Columns))
	newRows := make([][]interface{}, 0, len(inputResult.Rows))
	logger.Debug("Executing Project: ", zap.Any("columns", inputResult.Columns))
	// Определяем колонки результата
	for _, col := range node.Columns {
		if col.IsSelectAll {
			// Добавляем все колонки с qualified именами
			for _, inputCol := range inputResult.Columns {
				// Создаем копию колонки
				newCol := table.TableColumn{
					Idx:               inputCol.Idx,
					Name:              inputCol.Name,
					Type:              inputCol.Type,
					TableIdentifier:   inputCol.TableIdentifier,
					OriginalTableName: inputCol.OriginalTableName,
				}
				// Если у колонки есть TableIdentifier, добавляем его к имени для вывода
				if inputCol.TableIdentifier != "" && !strings.Contains(inputCol.Name, ".") {
					newCol.Name = inputCol.TableIdentifier + "." + inputCol.Name
				}
				newColumns = append(newColumns, newCol)
			}
		} else if col.ColumnName != "" && col.Function == nil && col.ExpressionContext == nil {
			// Простая колонка по имени
			colIdx, err := findColumnIndex(inputResult.Columns, col.ColumnName, col.TableAlias)
			if err != nil {
				return nil, err
			}

			inputCol := inputResult.Columns[colIdx]

			// Создаем копию колонки
			newCol := table.TableColumn{
				Idx:               inputCol.Idx,
				Type:              inputCol.Type,
				TableIdentifier:   inputCol.TableIdentifier,
				OriginalTableName: inputCol.OriginalTableName,
			}

			// Определяем имя колонки для вывода
			if col.Alias != "" {
				// Если задан явный алиас в SELECT (например, SELECT id AS user_id)
				newCol.Name = col.Alias
				newCol.TableIdentifier = "" // Убираем TableIdentifier для явного алиаса
			} else if col.TableAlias != "" && inputCol.TableIdentifier != "" {
				// Если колонка выбрана с qualified именем (например, SELECT o.id)
				// Выводим с qualified именем
				newCol.Name = inputCol.TableIdentifier + "." + inputCol.Name
			} else {
				// Простое имя без алиаса таблицы (например, SELECT id без JOIN)
				newCol.Name = col.ColumnName
			}

			newColumns = append(newColumns, newCol)
		}
		// } else {
		// 	// Выражение или функция - создаем новую колонку
		// 	colName := col.Alias
		// 	if colName == "" {
		// 		if col.Function != nil {
		// 			colName = col.Function.Name
		// 		} else {
		// 			colName = "?column?"
		// 		}
		// 	}
		// 	inputCol := inputResult.Columns[0] // TODO: тут может быть несколько колонок из выражения, надо это учесть
		// 	var colType table.ColType = inputCol.Type
		// 	if col.Function != nil {
		// 		if fn, ok := functions.GetRegisteredFunction(col.Function.Name); ok {
		// 			colType = fn.GetFunction().ProRetType.ToColType()
		// 		}
		// 	}
		// 	newColumns = append(newColumns, table.TableColumn{
		// 		Name: colName,
		// 		Type: colType,
		// 	})
		// }
	}

	// Обрабатываем каждую строку
	for _, row := range inputResult.Rows {
		newRow := make([]interface{}, 0, len(node.Columns))

		for _, col := range node.Columns {
			if col.IsSelectAll {
				// Добавляем все значения
				newRow = append(newRow, row...)
			} else if col.ColumnName != "" && col.Function == nil && col.ExpressionContext == nil {
				// Простая колонка
				colIdx, err := findColumnIndex(inputResult.Columns, col.ColumnName, col.TableAlias)
				if err != nil {
					return nil, err
				}
				newRow = append(newRow, row[colIdx])
			} else {
				// Выражение или функция
				val, err := evaluateColumnExpression(&col, inputResult, row, ctx)
				if err != nil {
					return nil, err
				}
				var colType table.ColType = val.Columns[0].Type
				if col.Function != nil {
					if fn, ok := functions.GetRegisteredFunction(col.Function.Name); ok {
						colType = fn.GetFunction().ProRetType.ToColType()
					}
				}
				newColumns = append(newColumns, table.TableColumn{
					Name: col.Alias,
					Type: colType,
				})
				value := val.Rows[0][0]

				newRow = append(newRow, value)
			}
		}

		newRows = append(newRows, newRow)
	}

	logger.Debug("Project result", zap.Any("rows", newRows), zap.Any("cols", newColumns))

	return &table.ExecuteResult{
		Rows:    newRows,
		Columns: newColumns,
	}, nil
}

// executeFilter выполняет фильтрацию (WHERE)
func executeFilter(
	node *FilterPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	inputResult, err := executePlanNode(node.Input, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	// Получаем текущий statement из плана
	var statement *statements.SelectStatement
	if plan != nil && plan.Statement != nil {
		statement = plan.Statement
	}
	return revaluate.EvaluateFilterRowsByWhere(ctx.GoCtx, inputResult, node.Condition, statement, ctx.SubqueryExecutor, runtimeParams), nil
}

// executeJoin выполняет соединение таблиц
func executeJoin(
	node *JoinPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	leftResult, err := executePlanNode(node.Left, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, fmt.Errorf("error executing left side of join: %w", err)
	}

	rightResult, err := executePlanNode(node.Right, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, fmt.Errorf("error executing right side of join: %w", err)
	}

	// Объединяем колонки
	combinedColumns := append(leftResult.Columns, rightResult.Columns...)

	// Выполняем join в зависимости от типа
	var resultRows [][]interface{}
	switch strings.ToUpper(node.JoinType) {
	case "INNER", "":
		resultRows = performInnerJoin(leftResult, rightResult, node.Condition, ctx)
	case "LEFT":
		resultRows = performLeftJoin(leftResult, rightResult, node.Condition, ctx)
	case "CROSS":
		resultRows = performCrossJoin(leftResult, rightResult)
	default:
		return nil, fmt.Errorf("unsupported join type: %s", node.JoinType)
	}

	return &table.ExecuteResult{
		Rows:    resultRows,
		Columns: combinedColumns,
	}, nil
}

// executeAggregate выполняет агрегацию (GROUP BY)
func executeAggregate(
	node *AggregatePlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	inputResult, err := executePlanNode(node.Input, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	// Используем существующую реализацию из executor
	// TODO: Перенести логику сюда для чистоты
	return processGroupByAndAggregates(node.GroupBy, node.Aggregates, inputResult, ctx)
}

// executeSort выполняет сортировку (ORDER BY)
func executeSort(
	node *SortPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	inputResult, err := executePlanNode(node.Input, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	revaluate.EvaluateSortRows(inputResult, node.OrderBy)
	return inputResult, nil
}

// executeLimit выполняет ограничение количества строк (LIMIT/OFFSET)
func executeLimit(
	node *LimitPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	inputResult, err := executePlanNode(node.Input, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	// Применяем OFFSET
	startIdx := node.Offset
	if startIdx > len(inputResult.Rows) {
		startIdx = len(inputResult.Rows)
	}

	// Применяем LIMIT
	endIdx := len(inputResult.Rows)
	if node.Limit > 0 && startIdx+node.Limit < endIdx {
		endIdx = startIdx + node.Limit
	}

	inputResult.Rows = inputResult.Rows[startIdx:endIdx]
	return inputResult, nil
}

// executeUnion выполняет UNION
func executeUnion(
	node *UnionPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	leftResult, err := executePlanNode(node.Left, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	rightResult, err := executePlanNode(node.Right, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	return applyUnion(leftResult, rightResult, node.All)
}

// executeIntersect выполняет INTERSECT
func executeIntersect(
	node *IntersectPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	leftResult, err := executePlanNode(node.Left, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	rightResult, err := executePlanNode(node.Right, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	return applyIntersect(leftResult, rightResult, node.All)
}

// executeExcept выполняет EXCEPT
func executeExcept(
	node *ExceptPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	leftResult, err := executePlanNode(node.Left, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	rightResult, err := executePlanNode(node.Right, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	return applyExcept(leftResult, rightResult, node.All)
}

// executeSubquery выполняет подзапрос
func executeSubquery(
	node *SubqueryPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	result, err := executePlanNode(node.Subquery, plan, ctx, tx, runtimeParams)
	if err != nil {
		return nil, err
	}

	// Применяем алиас к колонкам
	if node.Alias != "" {
		for i := range result.Columns {
			result.Columns[i].TableIdentifier = node.Alias
		}
	}

	return result, nil
}

// Вспомогательные функции

func findColumnIndex(columns []table.TableColumn, name string, tableAlias string) (int, error) {
	if tableAlias != "" {
		// Если указан алиас таблицы, проверяем его
		for i, col := range columns {
			if col.TableIdentifier == tableAlias && col.Name == name {
				return i, nil
			}
		}
		return -1, fmt.Errorf("column \"%s.%s\" does not exist", tableAlias, name)
	}

	// Ищем по простому имени и проверяем уникальность
	var matchedIndices []int
	var matchedNames []string

	for i, col := range columns {
		if col.Name == name {
			matchedIndices = append(matchedIndices, i)
			if col.TableIdentifier != "" {
				matchedNames = append(matchedNames, col.TableIdentifier+"."+col.Name)
			} else {
				matchedNames = append(matchedNames, col.Name)
			}
		}
	}

	if len(matchedIndices) == 0 {
		return -1, fmt.Errorf("column \"%s\" does not exist", name)
	}

	if len(matchedIndices) > 1 {
		return -1, fmt.Errorf("column reference \"%s\" is ambiguous", name)
	}

	return matchedIndices[0], nil
}

func evaluateColumnExpression(col *items.SelectItem, inputResult *table.ExecuteResult, row []interface{}, ctx ExecutorContext) (*table.ExecuteResult, error) {
	if col.Function != nil {
		// Проверяем, что это не агрегатная функция
		// Агрегатные функции должны обрабатываться в AggregatePlanNode
		if col.Function.IsAggregate {
			return nil, fmt.Errorf("aggregate function %s cannot be used in non-aggregate context", col.Function.Name)
		}

		// Вычисляем аргументы функции
		args := make([]interface{}, 0, len(col.Function.Args))
		for _, argExpr := range col.Function.Args {
			if argExpr.STAR() != nil {
				args = append(args, 1)
				continue
			}

			if argExpr.Expression() == nil {
				return nil, fmt.Errorf("invalid function argument")
			}

			// Создаем временный ResultRow для вычисления выражения
			tempRow := &table.ResultRow{
				Row:     row,
				Columns: inputResult.Columns,
			}

			val, err := revaluate.EvaluateExpressionContext(ctx.GoCtx, argExpr.Expression(), tempRow, ctx.SubqueryExecutor, nil)
			if err != nil {
				return nil, fmt.Errorf("error evaluating function arg: %w", err)
			}
			if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
				args = append(args, valExpr.Row.Rows[0][0])
			} else {
				return nil, fmt.Errorf("invalid function arg")
			}
		}

		// Вызываем обычную (неагрегатную) функцию
		result, err := functions.ExecuteFunctionWithContext(
			ctx.GoCtx,
			col.Function.Name,
			args,
		)
		if err != nil {
			return nil, err
		}
		if len(result.Rows) > 0 && len(result.Rows[0]) > 0 {
			return &table.ExecuteResult{
				Rows:    result.Rows,
				Columns: result.Columns,
			}, nil
		}
		return nil, nil
	} else if col.ExpressionContext != nil {
		// Вычисляем выражение
		tempRow := &table.ResultRow{
			Row:     row,
			Columns: inputResult.Columns,
		}

		val, err := revaluate.EvaluateExpressionContext(ctx.GoCtx, col.ExpressionContext, tempRow, ctx.SubqueryExecutor, nil)
		if err != nil {
			return nil, err
		}
		if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
			return &table.ExecuteResult{
				Rows:    valExpr.Row.Rows,
				Columns: valExpr.Row.Columns,
			}, nil
		}
		if boolExpr, ok := val.(*rmodels.BoolExpression); ok {
			return &table.ExecuteResult{
				Rows:    [][]interface{}{{boolExpr.Value}},
				Columns: []table.TableColumn{{Idx: 0, Name: "?column?", Type: table.ColTypeBool}},
			}, nil
		}
		return nil, fmt.Errorf("invalid expression result")
	}

	return nil, fmt.Errorf("unsupported column expression")
}

func performInnerJoin(left, right *table.ExecuteResult, condition parser.IExpressionContext, ctx ExecutorContext) [][]interface{} {
	// Если нет условия, выполняем CROSS JOIN
	if condition == nil {
		return performCrossJoin(left, right)
	}

	// Иначе используем nested loop join
	return performNestedLoopJoin(left, right, condition, ctx)
}

// findColumnInSet ищет колонку по имени или qualified имени в наборе колонок
func findColumnInSet(name string, cols []table.TableColumn) int {
	for i, col := range cols {
		// Проверяем полное совпадение с именем
		if col.Name == name {
			return i
		}
		// Проверяем совпадение с qualified именем (table.column)
		if col.TableIdentifier != "" {
			qualifiedName := col.TableIdentifier + "." + col.Name
			if qualifiedName == name {
				return i
			}
		}
	}
	return -1
}

// performNestedLoopJoin выполняет nested loop join (fallback для сложных условий)
func performNestedLoopJoin(left, right *table.ExecuteResult, condition parser.IExpressionContext, ctx ExecutorContext) [][]interface{} {
	var result [][]interface{}

	// Объединяем колонки для вычисления условия
	combinedColumns := append(left.Columns, right.Columns...)

	errorCount := 0
	for rowIdx, leftRow := range left.Rows {
		for rightIdx, rightRow := range right.Rows {
			combinedRow := append(append([]interface{}{}, leftRow...), rightRow...)

			// Проверяем условие
			if condition != nil {
				tempRow := &table.ResultRow{
					Row:     combinedRow,
					Columns: combinedColumns,
				}

				// if rowIdx == 0 && rightIdx == 0 {

				val, err := revaluate.EvaluateExpressionContext(ctx.GoCtx, condition, tempRow, ctx.SubqueryExecutor, nil)
				if err != nil {
					logger.Debugf("Error evaluating JOIN condition for row %d (left) and row %d (right): %v", rowIdx, rightIdx, err)
					// Если ошибка при вычислении условия, пропускаем строку
					if errorCount < 3 {
						logger.Debug("JOIN condition error",
							zap.Error(err),
							zap.Any("leftRow", leftRow),
							zap.Any("rightRow", rightRow))
						errorCount++
					}
					continue
				}

				matches := evaluateConditionResult(val)
				if !matches {
					continue
				}
			}

			result = append(result, combinedRow)
		}
	}

	return result
}

// evaluateConditionResult извлекает булево значение из результата выражения
func evaluateConditionResult(val rmodels.Expression) bool {
	if boolExpr, ok := val.(*rmodels.BoolExpression); ok {
		return boolExpr.Value
	} else if resultExpr, ok := val.(*rmodels.ResultRowsExpression); ok {
		// Извлекаем значение из ResultRowsExpression
		if len(resultExpr.Row.Rows) > 0 && len(resultExpr.Row.Rows[0]) > 0 {
			rowVal := resultExpr.Row.Rows[0][0]
			if boolVal, ok := rowVal.(bool); ok {
				return boolVal
			}
		}
	}
	return false
}

// sortByColumn сортирует строки по указанной колонке
func sortByColumn(result *table.ExecuteResult, colIdx int) [][]interface{} {
	sorted := make([][]interface{}, len(result.Rows))
	copy(sorted, result.Rows)

	// Простая сортировка пузырьком (можно заменить на более эффективную)
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if compareValues(sorted[j][colIdx], sorted[j+1][colIdx]) > 0 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

// compareValues сравнивает два значения, возвращает -1, 0, 1
func compareValues(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	switch av := a.(type) {
	case int:
		if bv, ok := b.(int); ok {
			if av < bv {
				return -1
			} else if av > bv {
				return 1
			}
			return 0
		}
	case float64:
		if bv, ok := b.(float64); ok {
			if av < bv {
				return -1
			} else if av > bv {
				return 1
			}
			return 0
		}
	case string:
		if bv, ok := b.(string); ok {
			if av < bv {
				return -1
			} else if av > bv {
				return 1
			}
			return 0
		}
	case bool:
		if bv, ok := b.(bool); ok {
			if !av && bv {
				return -1
			} else if av && !bv {
				return 1
			}
			return 0
		}
	}

	// Если типы не совпадают или неизвестны, считаем их равными
	return 0
}

func performLeftJoin(left, right *table.ExecuteResult, condition parser.IExpressionContext, ctx ExecutorContext) [][]interface{} {
	var result [][]interface{}

	// Объединяем колонки для вычисления условия
	combinedColumns := append(left.Columns, right.Columns...)

	for _, leftRow := range left.Rows {
		hasMatch := false

		for _, rightRow := range right.Rows {
			combinedRow := append(append([]interface{}{}, leftRow...), rightRow...)

			// Проверяем условие, если оно задано
			matches := true
			if condition != nil {
				tempRow := &table.ResultRow{
					Row:     combinedRow,
					Columns: combinedColumns,
				}

				val, err := revaluate.EvaluateExpressionContext(ctx.GoCtx, condition, tempRow, ctx.SubqueryExecutor, nil)
				if err != nil {
					matches = false
				} else {
					matches = false
					if boolExpr, ok := val.(*rmodels.BoolExpression); ok {
						matches = boolExpr.Value
					} else if resultExpr, ok := val.(*rmodels.ResultRowsExpression); ok {
						// Извлекаем значение из ResultRowsExpression
						if len(resultExpr.Row.Rows) > 0 && len(resultExpr.Row.Rows[0]) > 0 {
							rowVal := resultExpr.Row.Rows[0][0]
							if boolVal, ok := rowVal.(bool); ok {
								matches = boolVal
							}
						}
					}
				}
			}

			if matches {
				result = append(result, combinedRow)
				hasMatch = true
			}
		}

		// Если не было совпадений, добавляем левую строку с NULL значениями для правой части
		if !hasMatch {
			nullRow := make([]interface{}, len(right.Columns))
			for i := range nullRow {
				nullRow[i] = nil
			}
			combinedRow := append(append([]interface{}{}, leftRow...), nullRow...)
			result = append(result, combinedRow)
		}
	}
	return result
}

func performCrossJoin(left, right *table.ExecuteResult) [][]interface{} {
	var result [][]interface{}
	for _, leftRow := range left.Rows {
		for _, rightRow := range right.Rows {
			combinedRow := append(append([]interface{}{}, leftRow...), rightRow...)
			result = append(result, combinedRow)
		}
	}
	return result
}

func processGroupByAndAggregates(groupBy []items.SelectItem, aggregates []items.SelectItem, inputResult *table.ExecuteResult, ctx ExecutorContext) (*table.ExecuteResult, error) {
	// Check for invalid combinations
	hasAggregates := false
	for _, col := range aggregates {
		if col.Function != nil && col.Function.IsAggregate {
			hasAggregates = true
			break
		}
	}

	if len(groupBy) == 0 && hasAggregates {
		// If aggregates without GROUP BY, no non-aggregate columns allowed
		for _, col := range aggregates {
			if col.ColumnName != "" && col.Function == nil {
				return nil, fmt.Errorf("column must appear in the GROUP BY clause or be used in an aggregate function")
			}
		}
	}

	var groups map[string][]*table.ResultRow

	// If GROUP BY exists, validate non-aggregate columns are in GROUP BY
	if len(groupBy) > 0 {
		for _, col := range aggregates {
			if col.ColumnName != "" && col.Function == nil {
				found := false
				for _, gb := range groupBy {
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

		// Group rows by GROUP BY keys
		groups = make(map[string][]*table.ResultRow)
		for _, row := range inputResult.Rows {
			resultRow := &table.ResultRow{
				Row:     row,
				Columns: inputResult.Columns,
			}
			groupKey, err := computeGroupKey(groupBy, resultRow, ctx)
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
	newColumns := make([]table.TableColumn, 0, len(aggregates))
	newRows := make([][]interface{}, 0, len(groups))

	// Define columns
	for colIdx, col := range aggregates {
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
			if len(groupBy) > 0 {
				found := false
				for _, gb := range groupBy {
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
		newRow := make([]interface{}, 0, len(aggregates))

		for idx, col := range aggregates {
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
							val, err := revaluate.EvaluateExpressionContext(ctx.GoCtx, argExpr.Expression(), groupRow, ctx.SubqueryExecutor, nil)
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

func computeGroupKey(groupBy []items.SelectItem, row *table.ResultRow, ctx ExecutorContext) (string, error) {
	keyParts := make([]string, 0, len(groupBy))
	for _, gb := range groupBy {
		var val interface{}
		var err error
		if gb.ExpressionContext != nil {
			val, err = revaluate.EvaluateExpressionContext(ctx.GoCtx, gb.ExpressionContext, row, ctx.SubqueryExecutor, nil)
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
