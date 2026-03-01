package planner

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
)

// PlanNode - базовый интерфейс для узлов плана выполнения
type PlanNode interface {
	NodeType() string
	OutputColumns() []table.TableColumn
}

// ScanPlanNode - сканирование таблицы
type ScanPlanNode struct {
	Database      string
	Schema        string
	TableName     string
	Alias         string
	IsSystemTable bool // true если это системная таблица pg_catalog
}

func (n *ScanPlanNode) NodeType() string { return "Scan" }
func (n *ScanPlanNode) OutputColumns() []table.TableColumn {
	return nil
}

// ProjectPlanNode - проекция (выбор колонок)
type ProjectPlanNode struct {
	Input   PlanNode
	Columns []items.SelectItem
}

func (n *ProjectPlanNode) NodeType() string { return "Project" }
func (n *ProjectPlanNode) OutputColumns() []table.TableColumn {
	return nil
}

// FilterPlanNode - фильтрация (WHERE)
type FilterPlanNode struct {
	Input     PlanNode
	Condition *items.WhereClause
}

func (n *FilterPlanNode) NodeType() string { return "Filter" }
func (n *FilterPlanNode) OutputColumns() []table.TableColumn {
	return n.Input.OutputColumns()
}

// JoinPlanNode - соединение таблиц
type JoinPlanNode struct {
	Left      PlanNode
	Right     PlanNode
	JoinType  string
	Condition parser.IExpressionContext
}

func (n *JoinPlanNode) NodeType() string { return "Join" }
func (n *JoinPlanNode) OutputColumns() []table.TableColumn {
	return nil
}

// AggregatePlanNode - агрегация (GROUP BY)
type AggregatePlanNode struct {
	Input      PlanNode
	GroupBy    []items.SelectItem
	Aggregates []items.SelectItem
}

func (n *AggregatePlanNode) NodeType() string { return "Aggregate" }
func (n *AggregatePlanNode) OutputColumns() []table.TableColumn {
	return nil
}

// SortPlanNode - сортировка (ORDER BY)
type SortPlanNode struct {
	Input   PlanNode
	OrderBy []items.OrderByItem
}

func (n *SortPlanNode) NodeType() string { return "Sort" }
func (n *SortPlanNode) OutputColumns() []table.TableColumn {
	return n.Input.OutputColumns()
}

// LimitPlanNode - ограничение количества строк
type LimitPlanNode struct {
	Input  PlanNode
	Limit  int
	Offset int
}

func (n *LimitPlanNode) NodeType() string { return "Limit" }
func (n *LimitPlanNode) OutputColumns() []table.TableColumn {
	return n.Input.OutputColumns()
}

// UnionPlanNode - объединение результатов (UNION)
type UnionPlanNode struct {
	Left  PlanNode
	Right PlanNode
	All   bool
}

func (n *UnionPlanNode) NodeType() string { return "Union" }
func (n *UnionPlanNode) OutputColumns() []table.TableColumn {
	return n.Left.OutputColumns()
}

// IntersectPlanNode - пересечение результатов (INTERSECT)
type IntersectPlanNode struct {
	Left  PlanNode
	Right PlanNode
	All   bool
}

func (n *IntersectPlanNode) NodeType() string { return "Intersect" }
func (n *IntersectPlanNode) OutputColumns() []table.TableColumn {
	return n.Left.OutputColumns()
}

// ExceptPlanNode - разность результатов (EXCEPT)
type ExceptPlanNode struct {
	Left  PlanNode
	Right PlanNode
	All   bool
}

func (n *ExceptPlanNode) NodeType() string { return "Except" }
func (n *ExceptPlanNode) OutputColumns() []table.TableColumn {
	return n.Left.OutputColumns()
}

// ValuesPlanNode - значения без таблицы (SELECT 1, 2, 3)
type ValuesPlanNode struct {
	Columns []items.SelectItem
}

func (n *ValuesPlanNode) NodeType() string { return "Values" }
func (n *ValuesPlanNode) OutputColumns() []table.TableColumn {
	return nil
}

// SubqueryPlanNode - подзапрос в FROM
type SubqueryPlanNode struct {
	Subquery PlanNode
	Alias    string
}

func (n *SubqueryPlanNode) NodeType() string { return "Subquery" }
func (n *SubqueryPlanNode) OutputColumns() []table.TableColumn {
	return n.Subquery.OutputColumns()
}

// QueryPlan - корневой узел плана запроса
type QueryPlan struct {
	Root             PlanNode
	Statement        *statements.SelectStatement
	IsReadOnly       bool
	SubqueryExecutor subquery.SubqueryExecutor
}

func (p *QueryPlan) String() string {
	return formatPlan(p.Root, 0)
}

func formatPlan(node PlanNode, indent int) string {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	var result string

	switch n := node.(type) {
	case *ScanPlanNode:
		tableName := n.TableName
		if n.Schema != "" && n.Schema != "public" {
			tableName = n.Schema + "." + tableName
		}
		if n.Alias != "" {
			tableName += " AS " + n.Alias
		}
		result = prefix + node.NodeType() + " (" + tableName + ")\n"
	case *ValuesPlanNode:
		result = prefix + node.NodeType() + " (constant values)\n"
	case *SubqueryPlanNode:
		alias := ""
		if n.Alias != "" {
			alias = " AS " + n.Alias
		}
		result = prefix + node.NodeType() + alias + "\n"
		result += formatPlan(n.Subquery, indent+1)
	case *LimitPlanNode:
		limitInfo := ""
		if n.Limit > 0 {
			limitInfo = fmt.Sprintf(" (limit=%d", n.Limit)
			if n.Offset > 0 {
				limitInfo += fmt.Sprintf(", offset=%d", n.Offset)
			}
			limitInfo += ")"
		} else if n.Offset > 0 {
			limitInfo = fmt.Sprintf(" (offset=%d)", n.Offset)
		}
		result = prefix + node.NodeType() + limitInfo + "\n"
	default:
		result = prefix + node.NodeType() + "\n"
	}

	switch n := node.(type) {
	case *ProjectPlanNode:
		result += formatPlan(n.Input, indent+1)
	case *FilterPlanNode:
		result += formatPlan(n.Input, indent+1)
	case *JoinPlanNode:
		result += formatPlan(n.Left, indent+1)
		result += formatPlan(n.Right, indent+1)
	case *AggregatePlanNode:
		result += formatPlan(n.Input, indent+1)
	case *SortPlanNode:
		result += formatPlan(n.Input, indent+1)
	case *LimitPlanNode:
		result += formatPlan(n.Input, indent+1)
	case *UnionPlanNode:
		result += formatPlan(n.Left, indent+1)
		result += formatPlan(n.Right, indent+1)
	case *IntersectPlanNode:
		result += formatPlan(n.Left, indent+1)
		result += formatPlan(n.Right, indent+1)
	case *ExceptPlanNode:
		result += formatPlan(n.Left, indent+1)
		result += formatPlan(n.Right, indent+1)
	}

	return result
}
