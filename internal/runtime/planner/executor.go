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
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
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
	Allocator        pmem.Allocator
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
	err = ctx.Storage.RunTransactionWithAllocator(ctx.Allocator, func(tx *storage.DistributedTransactionVClock) error {
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
	result, err := tbl.Select(ctx.Allocator, tx, node.TableName, selectColumns, nil, 0)
	if err != nil {
		return nil, err
	}

	// Устанавливаем TableIdentifier для всех колонок (алиас или имя таблицы)
	tableIdentifier := node.Alias
	if tableIdentifier == "" {
		tableIdentifier = node.TableName
	}
	// Store table identifier in TableAlias — keeps Name bare for display
	if tableIdentifier != "" {
		for i := range result.Schema.Fields {
			result.Schema.Fields[i].TableAlias = tableIdentifier
		}
		result.Schema.RebuildIndex()
	}

	logger.Debug("scan result", zap.Int("rows", len(result.Rows)), zap.Int("cols", len(result.Schema.Fields)))
	return result, nil
}

// executeValues выполняет VALUES узел (SELECT без таблицы)
func executeValues(
	node *ValuesPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	resultColumns := make([]serializers.FieldDef, 0, len(node.Columns))
	rowValues := make([][]byte, 0, len(node.Columns))

	for _, col := range node.Columns {
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
				val, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, argExpr.Expression(), nil, ctx.SubqueryExecutor, runtimeParams)
				if err != nil {
					return nil, fmt.Errorf("error evaluating function arg: %w", err)
				}
				if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Schema.Fields) > 0 {
					buf, oid, err := valExpr.Row.Schema.GetField(valExpr.Row.Rows[0], 0)
					if err != nil {
						return nil, fmt.Errorf("failed to get field: %w", err)
					}
					desVal, err := serializers.DeserializeGeneric(buf, oid)
					if err != nil {
						return nil, fmt.Errorf("failed to deserialize: %w", err)
					}
					args = append(args, desVal)
				} else {
					return nil, fmt.Errorf("invalid function arg")
				}
			}
			v, err := functions.ExecuteFunctionWithContext(
				ctx.Allocator,
				ctx.GoCtx,
				col.Function.Name,
				args,
			)
			if err != nil {
				return nil, err
			}
			value = v
		} else if col.ExpressionContext != nil {
			v, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, col.ExpressionContext, nil, ctx.SubqueryExecutor, runtimeParams)
			if err != nil {
				return nil, err
			}
			if execRes, ok := v.(*rmodels.ResultRowsExpression); ok && len(execRes.Row.Rows) > 0 && len(execRes.Row.Schema.Fields) > 0 {
				value = execRes.Row
			} else if boolExpr, ok := v.(*rmodels.BoolExpression); ok {
				// Convert BoolExpression to ExecuteResult
				fields := []serializers.FieldDef{{
					Name: "?column?",
					OID:  ptypes.PTypeBool,
				}}
				schema := serializers.NewBaseSchema(fields)
				buf, err := serializers.BoolSerializerInstance.Serialize(ctx.Allocator, boolExpr.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to serialize bool: %w", err)
				}
				row, err := schema.Pack(ctx.Allocator, [][]byte{buf})
				if err != nil {
					return nil, fmt.Errorf("failed to pack bool: %w", err)
				}
				value = &table.ExecuteResult{
					Rows:   []*ptypes.Row{row},
					Schema: schema,
				}
			} else {
				return nil, fmt.Errorf("expected ExecuteResult from expression, got %T", v)
			}
		} else {
			return nil, fmt.Errorf("unsupported select column without table")
		}

		if col.Alias != "" {
			value.Schema.Fields[0].Name = col.Alias
		}

		buf, _, err := value.Schema.GetField(value.Rows[0], 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get field: %w", err)
		}

		rowValues = append(rowValues, buf)
		resultColumns = append(resultColumns, value.Schema.Fields[0])
	}

	// Создаем схему и упаковываем результат
	resultSchema := serializers.NewBaseSchema(resultColumns)
	packedRow, err := resultSchema.Pack(ctx.Allocator, rowValues)
	if err != nil {
		return nil, fmt.Errorf("failed to pack result row: %w", err)
	}

	return &table.ExecuteResult{
		Rows:   []*ptypes.Row{packedRow},
		Schema: resultSchema,
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

	// SELECT * — strip table prefix and return
	allSelectAll := true
	for _, col := range node.Columns {
		if !col.IsSelectAll {
			allSelectAll = false
			break
		}
	}
	if allSelectAll {
		for i := range inputResult.Schema.Fields {
			inputResult.Schema.Fields[i].TableAlias = "" // clear qualifier for output
		}
		return inputResult, nil
	}

	// Build a flat list of column specs so we can handle all column kinds uniformly.
	type colKind int
	const (
		colKindSimple     colKind = iota // plain column reference
		colKindSelectAll                 // SELECT *  (expanded inline)
		colKindFunction                  // function call
		colKindExpression                // arbitrary expression
	)
	type colSpec struct {
		kind   colKind
		field  serializers.FieldDef
		colIdx int              // for colKindSimple / colKindSelectAll
		col    items.SelectItem // original column (for func/expr)
	}

	specs := make([]colSpec, 0, len(node.Columns))
	for _, col := range node.Columns {
		if col.IsSelectAll {
			for i, field := range inputResult.Schema.Fields {
				f := field
				f.TableAlias = "" // clear qualifier for output
				specs = append(specs, colSpec{kind: colKindSelectAll, field: f, colIdx: i})
			}
		} else if col.ColumnName != "" && col.Function == nil && col.ExpressionContext == nil {
			colIdx, err := findColumnIndexByFieldName(inputResult.Schema.Fields, col.ColumnName, col.TableAlias)
			if err != nil {
				return nil, err
			}
			field := inputResult.Schema.Fields[colIdx]
			if col.Alias != "" {
				field.Name = col.Alias
			}
			field.TableAlias = "" // clear qualifier for output
			specs = append(specs, colSpec{kind: colKindSimple, field: field, colIdx: colIdx, col: col})
		} else if col.Function != nil {
			name := col.Function.Name
			if col.Alias != "" {
				name = col.Alias
			}
			oid := ptypes.PTypeText
			if fn, ok := functions.GetRegisteredFunction(col.Function.Name); ok {
				oid = fn.GetFunction().ProRetType
			}
			specs = append(specs, colSpec{kind: colKindFunction, field: serializers.FieldDef{Name: name, OID: oid}, col: col})
		} else if col.ExpressionContext != nil {
			name := "?column?"
			if col.Alias != "" {
				name = col.Alias
			}
			specs = append(specs, colSpec{kind: colKindExpression, field: serializers.FieldDef{Name: name, OID: ptypes.PTypeText}, col: col})
		} else {
			return nil, fmt.Errorf("unsupported column kind in projection")
		}
	}

	// Build schema from specs (OIDs for func/expr cols may be refined on first row)
	newFields := make([]serializers.FieldDef, len(specs))
	for i, sp := range specs {
		newFields[i] = sp.field
	}
	resultSchema := serializers.NewBaseSchema(newFields)

	resultRows := make([]*ptypes.Row, 0, len(inputResult.Rows))
	newRowValues := make([][]byte, 0, len(specs))

	for rowNum, inputRow := range inputResult.Rows {
		newRowValues = newRowValues[:0]
		currentRow := &table.ResultRow{Row: inputRow, Schema: inputResult.Schema}

		for si, sp := range specs {
			switch sp.kind {
			case colKindSimple, colKindSelectAll:
				buf, _, err := inputResult.Schema.GetField(inputRow, sp.colIdx)
				if err != nil {
					return nil, fmt.Errorf("failed to get field %d: %w", sp.colIdx, err)
				}
				newRowValues = append(newRowValues, buf)

			case colKindFunction:
				col := sp.col
				args := make([]interface{}, 0, len(col.Function.Args))
				for _, argExpr := range col.Function.Args {
					if argExpr.STAR() != nil {
						args = append(args, 1)
						continue
					}
					if argExpr.Expression() == nil {
						return nil, fmt.Errorf("invalid function argument")
					}
					val, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, argExpr.Expression(), currentRow, ctx.SubqueryExecutor, runtimeParams)
					if err != nil {
						return nil, fmt.Errorf("error evaluating function arg: %w", err)
					}
					if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 {
						buf, oid, err := valExpr.Row.Schema.GetField(valExpr.Row.Rows[0], 0)
						if err != nil {
							return nil, fmt.Errorf("failed to get function arg field: %w", err)
						}
						desVal, err := serializers.DeserializeGeneric(buf, oid)
						if err != nil {
							return nil, fmt.Errorf("failed to deserialize function arg: %w", err)
						}
						args = append(args, desVal)
					} else {
						args = append(args, nil)
					}
				}
				v, err := functions.ExecuteFunctionWithContext(ctx.Allocator, ctx.GoCtx, col.Function.Name, args)
				if err != nil {
					return nil, err
				}
				if len(v.Rows) > 0 && v.Schema != nil && len(v.Schema.Fields) > 0 {
					buf, oid, err := v.Schema.GetField(v.Rows[0], 0)
					if err == nil {
						if rowNum == 0 {
							resultSchema.Fields[si].OID = oid
						}
						newRowValues = append(newRowValues, buf)
						continue
					}
				}
				newRowValues = append(newRowValues, nil)

			case colKindExpression:
				col := sp.col
				val, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, col.ExpressionContext, currentRow, ctx.SubqueryExecutor, runtimeParams)
				if err != nil {
					return nil, err
				}
				switch v := val.(type) {
				case *rmodels.ResultRowsExpression:
					if len(v.Row.Rows) > 0 && v.Row.Schema != nil && len(v.Row.Schema.Fields) > 0 {
						buf, oid, err := v.Row.Schema.GetField(v.Row.Rows[0], 0)
						if err == nil {
							if rowNum == 0 {
								resultSchema.Fields[si].OID = oid
							}
							newRowValues = append(newRowValues, buf)
							continue
						}
					}
					newRowValues = append(newRowValues, nil)
				case *rmodels.BoolExpression:
					buf, err := serializers.BoolSerializerInstance.Serialize(ctx.Allocator, v.Value)
					if err != nil {
						return nil, fmt.Errorf("failed to serialize bool expression: %w", err)
					}
					if rowNum == 0 {
						resultSchema.Fields[si].OID = ptypes.PTypeBool
					}
					newRowValues = append(newRowValues, buf)
				default:
					newRowValues = append(newRowValues, nil)
				}
			}
		}

		newRow, err := resultSchema.Pack(ctx.Allocator, newRowValues)
		if err != nil {
			return nil, fmt.Errorf("failed to pack row: %w", err)
		}
		resultRows = append(resultRows, newRow)
	}

	logger.Debug("Project result", zap.Int("rows", len(resultRows)), zap.Int("cols", len(newFields)))

	return &table.ExecuteResult{
		Rows:   resultRows,
		Schema: resultSchema,
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
	return revaluate.EvaluateFilterRowsByWhere(ctx.Allocator, ctx.GoCtx, inputResult, node.Condition, statement, ctx.SubqueryExecutor, runtimeParams), nil
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

	// Объединяем схемы
	combinedFields := make([]serializers.FieldDef, 0, len(leftResult.Schema.Fields)+len(rightResult.Schema.Fields))
	combinedFields = append(combinedFields, leftResult.Schema.Fields...)
	combinedFields = append(combinedFields, rightResult.Schema.Fields...)
	combinedSchema := serializers.NewBaseSchema(combinedFields)

	// Выполняем join в зависимости от типа
	var resultRows []*ptypes.Row
	switch strings.ToUpper(node.JoinType) {
	case "INNER", "":
		resultRows = performInnerJoin(leftResult, rightResult, node.Condition, ctx, combinedSchema, runtimeParams)
	case "LEFT":
		resultRows = performLeftJoin(leftResult, rightResult, node.Condition, ctx, combinedSchema, runtimeParams)
	case "CROSS":
		resultRows = performCrossJoin(leftResult, rightResult, ctx, combinedSchema)
	default:
		return nil, fmt.Errorf("unsupported join type: %s", node.JoinType)
	}

	return &table.ExecuteResult{
		Rows:   resultRows,
		Schema: combinedSchema,
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

	err = revaluate.EvaluateSortRows(ctx.Allocator, inputResult, node.OrderBy)
	if err != nil {
		return nil, err
	}
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

	// Store alias in TableAlias — keeps Name bare
	if node.Alias != "" {
		for i := range result.Schema.Fields {
			result.Schema.Fields[i].TableAlias = node.Alias
		}
		result.Schema.RebuildIndex()
	}

	return result, nil
}

// Вспомогательные функции

// findColumnIndexByFieldName ищет индекс колонки в схеме по имени
func findColumnIndexByFieldName(fields []serializers.FieldDef, name string, tableAlias string) (int, error) {
	nameLower := strings.ToLower(name)
	aliasLower := strings.ToLower(tableAlias)

	if tableAlias != "" {
		// Qualified lookup: match alias + bare name
		for i, field := range fields {
			if strings.ToLower(field.Name) == nameLower && strings.ToLower(field.TableAlias) == aliasLower {
				return i, nil
			}
		}
		return -1, fmt.Errorf("column \"%s.%s\" does not exist", tableAlias, name)
	}

	// Unqualified: match by bare Name, ensure uniqueness
	var matchedIndices []int
	for i, field := range fields {
		if strings.ToLower(field.Name) == nameLower {
			matchedIndices = append(matchedIndices, i)
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

func performInnerJoin(left, right *table.ExecuteResult, condition parser.IExpressionContext, ctx ExecutorContext, combinedSchema *serializers.BaseSchema, runtimeParams map[int]interface{}) []*ptypes.Row {
	// Если нет условия, выполняем CROSS JOIN
	if condition == nil {
		return performCrossJoin(left, right, ctx, combinedSchema)
	}

	// Иначе используем nested loop join
	return performNestedLoopJoin(left, right, condition, ctx, combinedSchema, runtimeParams)
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
func performNestedLoopJoin(left, right *table.ExecuteResult, condition parser.IExpressionContext, ctx ExecutorContext, combinedSchema *serializers.BaseSchema, runtimeParams map[int]interface{}) []*ptypes.Row {
	var result []*ptypes.Row
	scratch := make([][]byte, len(left.Schema.Fields)+len(right.Schema.Fields)) // scratch для объединения строк

	errorCount := 0
	for rowIdx, leftRow := range left.Rows {
		for rightIdx, rightRow := range right.Rows {
			// Объединяем строки
			combinedRow, err := combineRows(leftRow, rightRow, left.Schema, right.Schema, combinedSchema, ctx.Allocator, scratch)
			if err != nil {
				logger.Debugf("Error combining rows: %v", err)
				continue
			}

			// Проверяем условие
			if condition != nil {
				tempRow := &table.ResultRow{
					Row:    combinedRow,
					Schema: combinedSchema,
				}

				val, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, condition, tempRow, ctx.SubqueryExecutor, runtimeParams)
				if err != nil {
					logger.Debugf("Error evaluating JOIN condition for row %d (left) and row %d (right): %v", rowIdx, rightIdx, err)
					// Если ошибка при вычислении условия, пропускаем строку
					if errorCount < 3 {
						logger.Debug("JOIN condition error",
							zap.Error(err))
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
		if len(resultExpr.Row.Rows) > 0 {
			// Пытаемся извлечь первое поле первой строки
			firstRow := resultExpr.Row.Rows[0]
			if resultExpr.Row.Schema != nil && len(resultExpr.Row.Schema.Fields) > 0 {
				buf, oid, err := resultExpr.Row.Schema.GetField(firstRow, 0)
				if err == nil {
					// Десериализуем значение
					val, err := serializers.DeserializeGeneric(buf, oid)
					if err == nil {
						boolVal, ok := ptypes.TryIntoBool(val)
						if ok {
							return boolVal.IntoGo()
						}
					}
				}
			}
		}
	}
	return false
}

// combineRows объединяет две бинарные строки в одну
func combineRows(leftRow, rightRow *ptypes.Row, leftSchema, rightSchema, combinedSchema *serializers.BaseSchema, allocator pmem.Allocator, scratch [][]byte) (*ptypes.Row, error) {
	// scratch передаётся снаружи и переиспользуется
	nL := len(leftSchema.Fields)
	nR := len(rightSchema.Fields)
	scratch = scratch[:nL+nR]

	for i := range leftSchema.Fields {
		buf, _, err := leftSchema.GetField(leftRow, i)
		if err != nil {
			return nil, err
		}
		scratch[i] = buf
	}
	for i := range rightSchema.Fields {
		buf, _, err := rightSchema.GetField(rightRow, i)
		if err != nil {
			return nil, err
		}
		scratch[nL+i] = buf
	}
	return combinedSchema.Pack(allocator, scratch)
}

// createNullRow создает строку со всеми NULL значениями для заданной схемы
func createNullRow(schema *serializers.BaseSchema, allocator pmem.Allocator) (*ptypes.Row, error) {
	nullValues := make([][]byte, len(schema.Fields))
	for i := range schema.Fields {
		// NULL представлен как nil в байтовом массиве
		nullValues[i] = nil
	}

	row, err := schema.Pack(allocator, nullValues)
	if err != nil {
		return nil, fmt.Errorf("error packing null row: %w", err)
	}

	return row, nil
}

func performLeftJoin(left, right *table.ExecuteResult, condition parser.IExpressionContext, ctx ExecutorContext, combinedSchema *serializers.BaseSchema, runtimeParams map[int]interface{}) []*ptypes.Row {
	var result []*ptypes.Row

	for _, leftRow := range left.Rows {
		hasMatch := false

		scratch := make([][]byte, len(left.Schema.Fields)+len(right.Schema.Fields)) // scratch для объединения строк

		for _, rightRow := range right.Rows {
			// Объединяем строки
			combinedRow, err := combineRows(leftRow, rightRow, left.Schema, right.Schema, combinedSchema, ctx.Allocator, scratch)
			if err != nil {
				logger.Debugf("Error combining rows: %v", err)
				continue
			}

			// Проверяем условие, если оно задано
			matches := true
			if condition != nil {
				tempRow := &table.ResultRow{
					Row:    combinedRow,
					Schema: combinedSchema,
				}

				val, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, condition, tempRow, ctx.SubqueryExecutor, runtimeParams)
				if err != nil {
					matches = false
				} else {
					matches = evaluateConditionResult(val)
				}
			}

			if matches {
				result = append(result, combinedRow)
				hasMatch = true
			}
		}

		// Если не было совпадений, добавляем левую строку с NULL значениями для правой части
		if !hasMatch {
			// Создаем пустую правую строку (NULL значения)
			nullRightRow, err := createNullRow(right.Schema, ctx.Allocator)
			if err != nil {
				logger.Debugf("Error creating null row: %v", err)
				continue
			}
			combinedRow, err := combineRows(leftRow, nullRightRow, left.Schema, right.Schema, combinedSchema, ctx.Allocator, scratch)
			if err != nil {
				logger.Debugf("Error combining with null row: %v", err)
				continue
			}
			result = append(result, combinedRow)
		}
	}
	return result
}

func performCrossJoin(left, right *table.ExecuteResult, ctx ExecutorContext, combinedSchema *serializers.BaseSchema) []*ptypes.Row {
	var result []*ptypes.Row
	scratch := make([][]byte, len(left.Schema.Fields)+len(right.Schema.Fields)) // scratch для объединения строк
	for _, leftRow := range left.Rows {
		for _, rightRow := range right.Rows {
			combinedRow, err := combineRows(leftRow, rightRow, left.Schema, right.Schema, combinedSchema, ctx.Allocator, scratch)
			if err != nil {
				logger.Debugf("Error combining rows in CROSS JOIN: %v", err)
				continue
			}
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
				Row:    row,
				Schema: inputResult.Schema,
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
				Row:    row,
				Schema: inputResult.Schema,
			}
		}
		groups = map[string][]*table.ResultRow{"": groupRows}
	}

	// Now process each group
	newFields := make([]serializers.FieldDef, 0, len(aggregates))
	newRows := make([]*ptypes.Row, 0, len(groups))

	// Define columns (fields)
	for _, col := range aggregates {
		var colName string
		var colOID ptypes.OID

		if col.IsSelectAll {
			return nil, fmt.Errorf("SELECT * not allowed with GROUP BY or aggregates")
		}

		if col.ExpressionContext != nil {
			colOID = ptypes.PTypeText // default
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = "?column?"
			}
		} else if col.Function != nil {
			colOID = ptypes.PTypeText // default, will be updated by aggregate function
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
			// Get actual type from input schema fields
			colOID = ptypes.PTypeText // default
			for _, inputField := range inputResult.Schema.Fields {
				if strings.EqualFold(inputField.Name, col.ColumnName) &&
					(col.TableAlias == "" || strings.EqualFold(inputField.TableAlias, col.TableAlias)) {
					colOID = inputField.OID
					break
				}
			}
			if col.Alias != "" {
				colName = col.Alias
			} else {
				colName = col.ColumnName
			}
		}

		newFields = append(newFields, serializers.FieldDef{
			Name: colName,
			OID:  colOID,
		})
	}

	newSchema := serializers.NewBaseSchema(newFields)

	// Process each group
	for _, groupRows := range groups {
		newRowValues := make([][]byte, 0, len(aggregates))

		for idx, col := range aggregates {
			var colValueBuf []byte

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
							runtimeParams := make(map[int]interface{})
							val, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, argExpr.Expression(), groupRow, ctx.SubqueryExecutor, runtimeParams)
							if err != nil {
								return nil, fmt.Errorf("error evaluating aggregate arg: %w", err)
							}
							if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok && len(valExpr.Row.Rows) > 0 {
								// Extract first field from first row
								firstRow := valExpr.Row.Rows[0]
								if valExpr.Row.Schema != nil && len(valExpr.Row.Schema.Fields) > 0 {
									buf, oid, err := valExpr.Row.Schema.GetField(firstRow, 0)
									if err == nil {
										deserVal, err := serializers.DeserializeGeneric(buf, oid)
										if err == nil {
											groupValues = append(groupValues, deserVal)
										} else {
											groupValues = append(groupValues, nil)
										}
									} else {
										groupValues = append(groupValues, nil)
									}
								} else {
									groupValues = append(groupValues, nil)
								}
							} else {
								groupValues = append(groupValues, nil)
							}
						}
						args = append(args, groupValues)
					}
					v, err := functions.ExecuteAggregateFunction(ctx.Allocator, col.Function.Name, args)
					if err != nil {
						return nil, fmt.Errorf("error executing aggregate function %s: %w", col.Function.Name, err)
					}
					// Extract result value and serialize
					if len(v.Rows) > 0 {
						firstRow := v.Rows[0]
						if v.Schema != nil && len(v.Schema.Fields) > 0 {
							buf, oid, err := v.Schema.GetField(firstRow, 0)
							if err == nil {
								colValueBuf = buf
								// Update field OID based on actual result
								newFields[idx].OID = oid
							} else {
								colValueBuf = nil
							}
						} else {
							colValueBuf = nil
						}
					} else {
						colValueBuf = nil
					}
				} else {
					// Non-aggregate function, evaluate on first row or something? For now error
					return nil, fmt.Errorf("non-aggregate functions not supported in GROUP BY")
				}
			} else if col.ColumnName != "" {
				// Non-aggregate column, take from first row in group
				if len(groupRows) > 0 {
					resultRow := groupRows[0]
					found := false
					for i, origField := range inputResult.Schema.Fields {
						if strings.EqualFold(origField.Name, col.ColumnName) &&
							(col.TableAlias == "" || strings.EqualFold(origField.TableAlias, col.TableAlias)) {
							buf, _, err := resultRow.Schema.GetField(resultRow.Row, i)
							if err == nil {
								colValueBuf = buf
							} else {
								colValueBuf = nil
							}
							found = true
							break
						}
					}
					if !found {
						return nil, fmt.Errorf("column %s not found", col.ColumnName)
					}
				}
			}

			newRowValues = append(newRowValues, colValueBuf)
		}

		// Pack the row
		newRow, err := newSchema.Pack(ctx.Allocator, newRowValues)
		if err != nil {
			return nil, fmt.Errorf("error packing aggregated row: %w", err)
		}
		newRows = append(newRows, newRow)
	}

	return &table.ExecuteResult{
		Schema: newSchema,
		Rows:   newRows,
	}, nil
}

func computeGroupKey(groupBy []items.SelectItem, row *table.ResultRow, ctx ExecutorContext) (string, error) {
	keyParts := make([]string, 0, len(groupBy))
	for _, gb := range groupBy {
		var val interface{}
		if gb.ExpressionContext != nil {
			runtimeParams := make(map[int]interface{})
			valExpr, err := revaluate.EvaluateExpressionContext(ctx.Allocator, ctx.GoCtx, gb.ExpressionContext, row, ctx.SubqueryExecutor, runtimeParams)
			if err != nil {
				return "", err
			}
			if resultExpr, ok := valExpr.(*rmodels.ResultRowsExpression); ok && len(resultExpr.Row.Rows) > 0 {
				// Extract first field from first row
				firstRow := resultExpr.Row.Rows[0]
				if resultExpr.Row.Schema != nil && len(resultExpr.Row.Schema.Fields) > 0 {
					buf, oid, err := resultExpr.Row.Schema.GetField(firstRow, 0)
					if err == nil {
						val, _ = serializers.DeserializeGeneric(buf, oid)
					} else {
						val = nil
					}
				} else {
					val = nil
				}
			} else {
				val = nil
			}
		} else if gb.ColumnName != "" {
			found := false
			for i, field := range row.Schema.Fields {
				if strings.EqualFold(field.Name, gb.ColumnName) &&
					(gb.TableAlias == "" || strings.EqualFold(field.TableAlias, gb.TableAlias)) {
					buf, oid, err := row.Schema.GetField(row.Row, i)
					if err == nil {
						val, _ = serializers.DeserializeGeneric(buf, oid)
					} else {
						val = nil
					}
					found = true
					break
				}
			}
			if !found {
				return "", fmt.Errorf("group by column %s not found", gb.ColumnName)
			}
		} else {
			return "", fmt.Errorf("invalid group by item")
		}
		keyParts = append(keyParts, fmt.Sprintf("%v", val))
	}
	return strings.Join(keyParts, "|"), nil
}
