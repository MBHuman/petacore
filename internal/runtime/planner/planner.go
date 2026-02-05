package planner

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
	"strings"
)

// PlannerContext содержит контекст для создания плана
type PlannerContext struct {
	Database string
	Schema   string
}

// CreateQueryPlan создает план выполнения из SELECT statement
func CreateQueryPlan(stmt *statements.SelectStatement, ctx PlannerContext) (*QueryPlan, error) {
	rootNode, err := buildPlanFromSelectStatement(stmt, ctx)
	if err != nil {
		return nil, err
	}

	return &QueryPlan{
		Root:       rootNode,
		Statement:  stmt,
		IsReadOnly: true,
	}, nil
}

// buildPlanFromSelectStatement рекурсивно строит план из SELECT statement
func buildPlanFromSelectStatement(stmt *statements.SelectStatement, ctx PlannerContext) (PlanNode, error) {
	if stmt.IsCombined() {
		return buildCombinedPlan(stmt.Combined, ctx)
	} else if stmt.IsPrimary() {
		return buildPrimaryPlan(stmt.Primary, ctx)
	}

	return nil, fmt.Errorf("invalid SELECT statement: neither Primary nor Combined")
}

// buildCombinedPlan строит план для UNION/INTERSECT/EXCEPT
func buildCombinedPlan(combined *statements.CombinedSelectStatement, ctx PlannerContext) (PlanNode, error) {
	// Рекурсивно строим левую и правую части
	leftPlan, err := buildPlanFromSelectStatement(combined.Left, ctx)
	if err != nil {
		return nil, fmt.Errorf("error building left side of %s: %w", combined.OperationType, err)
	}

	rightPlan, err := buildPlanFromSelectStatement(combined.Right, ctx)
	if err != nil {
		return nil, fmt.Errorf("error building right side of %s: %w", combined.OperationType, err)
	}

	// Создаем соответствующий узел плана
	switch combined.OperationType {
	case statements.UnionOperation:
		return &UnionPlanNode{
			Left:  leftPlan,
			Right: rightPlan,
			All:   combined.All,
		}, nil
	case statements.IntersectOperation:
		return &IntersectPlanNode{
			Left:  leftPlan,
			Right: rightPlan,
			All:   combined.All,
		}, nil
	case statements.ExceptOperation:
		return &ExceptPlanNode{
			Left:  leftPlan,
			Right: rightPlan,
			All:   combined.All,
		}, nil
	default:
		return nil, fmt.Errorf("unknown set operation: %s", combined.OperationType)
	}
}

// buildPrimaryPlan строит план для одного PRIMARY SELECT
func buildPrimaryPlan(primary *statements.PrimarySelectStatement, ctx PlannerContext) (PlanNode, error) {
	var currentPlan PlanNode

	// 1. Если нет FROM - это VALUES узел (SELECT без таблицы)
	if primary.From == nil {
		currentPlan = &ValuesPlanNode{
			Columns: primary.Columns,
		}
	} else {
		// 2. Начинаем с FROM clause
		fromPlan, err := buildFromPlan(primary.From, ctx)
		if err != nil {
			return nil, err
		}
		currentPlan = fromPlan
	}

	// --- SCALAR SUBQUERY IN WHERE ---
	nextParam := 0
	if primary.Where != nil {
		// Переписываем выражение WHERE, собираем subplans
		newWhereExpr, subplans, err := rewriteExprAndCollectSubplans(primary.Where.ExpressionContext, ctx, &nextParam)
		if err != nil {
			return nil, err
		}
		// Оборачиваем Expr в *items.WhereClause
		newWhere := &items.WhereClause{ExpressionContext: exprToParserContext(newWhereExpr)}
		currentPlan = &FilterPlanNode{
			Input:     currentPlan,
			Condition: newWhere,
		}
		// Если есть scalar subplans — оборачиваем в InitPlanNode
		if len(subplans) > 0 {
			currentPlan = &InitPlanNode{
				Input:    currentPlan,
				Subplans: subplans,
			}
		}
	}

	// 3. Добавляем WHERE filter
	// (WHERE уже обработано выше)

	// 4. Проверяем наличие агрегатов или GROUP BY
	hasAggregates := false
	for _, col := range primary.Columns {
		if col.Function != nil && col.Function.IsAggregate {
			hasAggregates = true
			break
		}
	}

	if len(primary.GroupBy) > 0 || hasAggregates {
		currentPlan = &AggregatePlanNode{
			Input:      currentPlan,
			GroupBy:    primary.GroupBy,
			Aggregates: primary.Columns,
		}
	} else {
		// 5. Добавляем проекцию (выбор колонок) только если нет агрегатов
		// AggregatePlanNode уже обрабатывает проекцию колонок
		currentPlan = &ProjectPlanNode{
			Input:   currentPlan,
			Columns: primary.Columns,
		}
	}

	// 6. Добавляем ORDER BY
	if len(primary.OrderBy) > 0 {
		currentPlan = &SortPlanNode{
			Input:   currentPlan,
			OrderBy: primary.OrderBy,
		}
	}

	// 7. Добавляем LIMIT/OFFSET
	if primary.Limit > 0 || primary.Offset > 0 {
		currentPlan = &LimitPlanNode{
			Input:  currentPlan,
			Limit:  primary.Limit,
			Offset: primary.Offset,
		}
	}

	return currentPlan, nil
}

// buildFromPlan строит план для FROM clause (может включать JOIN)
func buildFromPlan(from *statements.FromClause, ctx PlannerContext) (PlanNode, error) {
	var currentPlan PlanNode

	// Основная таблица или подзапрос
	if from.SelectStatement != nil {
		// Подзапрос в FROM
		subqueryPlan, err := buildPlanFromSelectStatement(from.SelectStatement, ctx)
		if err != nil {
			return nil, fmt.Errorf("error building subquery: %w", err)
		}
		currentPlan = &SubqueryPlanNode{
			Subquery: subqueryPlan,
			Alias:    from.Alias,
		}
	} else {
		// Обычная таблица
		schema, tableName := resolveSchemaAndTable(from.TableName, ctx.Schema)
		currentPlan = &ScanPlanNode{
			Database:      ctx.Database,
			Schema:        schema,
			TableName:     tableName,
			Alias:         from.Alias,
			IsSystemTable: isPgCatalogTable(tableName),
		}
	}

	// Добавляем JOIN'ы
	for _, join := range from.Joins {
		schema, tableName := resolveSchemaAndTable(join.TableName, ctx.Schema)
		rightPlan := &ScanPlanNode{
			Database:      ctx.Database,
			Schema:        schema,
			TableName:     tableName,
			Alias:         join.Alias,
			IsSystemTable: isPgCatalogTable(tableName),
		}

		currentPlan = &JoinPlanNode{
			Left:      currentPlan,
			Right:     rightPlan,
			JoinType:  join.Type,
			Condition: join.OnCondition,
		}
	}

	return currentPlan, nil
}

// resolveSchemaAndTable определяет схему и имя таблицы
// Правила:
// 1. Если имя таблицы содержит schema.table, используем явную схему
// 2. Если таблица системная (pg_catalog), используем схему pg_catalog
// 3. Иначе используем схему из контекста
func resolveSchemaAndTable(tableName string, defaultSchema string) (string, string) {
	parts := strings.SplitN(tableName, ".", 2)
	if len(parts) == 2 {
		// Явно указана схема (например, pg_catalog.pg_type)
		return parts[0], parts[1]
	}

	// Проверяем, является ли это системной таблицей
	if isPgCatalogTable(tableName) {
		return "pg_catalog", tableName
	}

	// Используем схему по умолчанию
	return defaultSchema, tableName
}

// isPgCatalogTable проверяет, является ли таблица системной таблицей pg_catalog
func isPgCatalogTable(tableName string) bool {
	switch tableName {
	case "pg_tables", "pg_columns", "pg_class",
		"pg_attribute", "pg_proc", "pg_type", "pg_namespace",
		"pg_database", "pg_tablespace", "pg_roles",
		"pg_stat_ssl", "pg_shdescription", "pg_am",
		"pg_index", "pg_constraint":
		return true
	default:
		return false
	}
}

// --- SCALAR SUBQUERY SUPPORT ---
type InitPlanNode struct {
	Input    PlanNode
	Subplans []Subplan
}

func (n *InitPlanNode) NodeType() string { return "Init" }
func (n *InitPlanNode) OutputColumns() []table.TableColumn {
	return n.Input.OutputColumns()
}

type Subplan struct {
	Plan  PlanNode
	Param int
	Kind  string // "scalar"
}

type ScalarSubqueryPlanNode struct {
	Input PlanNode
}

func (n *ScalarSubqueryPlanNode) NodeType() string { return "ScalarSubquery" }
func (n *ScalarSubqueryPlanNode) OutputColumns() []table.TableColumn {
	return n.Input.OutputColumns()
}

type ParamRefExpr struct {
	Index int
}

// Expr interface for rewriting (adapt to your AST)
type Expr interface{}

type SubqueryExpr struct {
	Select *statements.SelectStatement
}

// Example: rewriteExprAndCollectSubplans
func rewriteExprAndCollectSubplans(
	expr Expr,
	ctx PlannerContext,
	nextParam *int,
) (Expr, []Subplan, error) {
	var subplans []Subplan
	var walk func(e Expr) (Expr, error)
	walk = func(e Expr) (Expr, error) {
		switch t := e.(type) {
		case *SubqueryExpr:
			p, err := buildPlanFromSelectStatement(t.Select, ctx)
			if err != nil {
				return nil, err
			}
			scalar := &ScalarSubqueryPlanNode{Input: p}
			param := *nextParam
			*nextParam++
			subplans = append(subplans, Subplan{Plan: scalar, Param: param, Kind: "scalar"})
			return &ParamRefExpr{Index: param}, nil
		// ... handle other expression types ...
		default:
			return e, nil
		}
	}
	newExpr, err := walk(expr)
	if err != nil {
		return nil, nil, err
	}
	return newExpr, subplans, nil
}

// Исполнение InitPlanNode
func executeInitPlanNode(
	node *InitPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	// Вычисляем все scalar subplans
	logger.Debugf("InitPlanNode: executing %d scalar subplans", len(node.Subplans))
	for i, sub := range node.Subplans {
		// sub.Plan всегда ScalarSubqueryPlanNode
		scalarNode, ok := sub.Plan.(*ScalarSubqueryPlanNode)
		if !ok {
			return nil, fmt.Errorf("subplan is not ScalarSubqueryPlanNode")
		}
		logger.Debugf("InitPlanNode: executing subplan %d param=%d nodeType=%s", i, sub.Param, scalarNode.Input.NodeType())
		res, err := executeScalarSubqueryPlanNode(scalarNode, plan, ctx, tx, runtimeParams)
		if err != nil {
			logger.Errorf("InitPlanNode: error executing subplan param=%d: %v", sub.Param, err)
			return nil, err
		}
		// Проверяем, что результат scalar: максимум 1 строка, 1 колонка
		val := interface{}(nil)
		rows := 0
		cols := 0
		if res != nil {
			rows = len(res.Rows)
			if rows > 0 {
				cols = len(res.Rows[0])
			}
		}
		if rows == 0 {
			val = nil // NULL
		} else if rows == 1 && cols == 1 {
			val = res.Rows[0][0]
		} else {
			logger.Errorf("InitPlanNode: scalar subquery param=%d returned rows=%d cols=%d", sub.Param, rows, cols)
			return nil, fmt.Errorf("scalar subquery returned more than one row or column")
		}
		logger.Debugf("InitPlanNode: subplan param=%d result rows=%d cols=%d value=%v", sub.Param, rows, cols, val)
		runtimeParams[sub.Param] = val
	}
	// Запускаем основной Input
	return executePlanNode(node.Input, plan, ctx, tx, runtimeParams)
}

// Исполнение ScalarSubqueryPlanNode
func executeScalarSubqueryPlanNode(
	node *ScalarSubqueryPlanNode,
	plan *QueryPlan,
	ctx ExecutorContext,
	tx *storage.DistributedTransactionVClock,
	runtimeParams map[int]interface{},
) (*table.ExecuteResult, error) {
	logger.Debugf("Executing ScalarSubqueryPlanNode input type=%s", node.Input.NodeType())
	res, err := executePlanNode(node.Input, plan, ctx, tx, runtimeParams)
	if err != nil {
		logger.Errorf("Error executing scalar subquery node type=%s: %v", node.Input.NodeType(), err)
		return nil, err
	}
	if res == nil {
		logger.Debugf("ScalarSubqueryPlanNode returned nil result for node type=%s", node.Input.NodeType())
	} else {
		logger.Debugf("ScalarSubqueryPlanNode result rows=%d cols=%d for node type=%s", len(res.Rows), func() int {
			if len(res.Rows) > 0 {
				return len(res.Rows[0])
			} else {
				return 0
			}
		}(), node.Input.NodeType())
	}
	return res, nil
}

// Конвертер Expr → parser.IExpressionContext
func exprToParserContext(expr Expr) parser.IExpressionContext {
	switch t := expr.(type) {
	case parser.IExpressionContext:
		return t
	// ... другие Expr ...
	default:
		return nil
	}
}
