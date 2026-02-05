package statements

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
)

type FromClause struct {
	Alias           string
	Joins           []JoinClause
	TableName       string
	SelectStatement *SelectStatement
}

type JoinClause struct {
	// TODO перевести на enum
	Type        string // INNER, LEFT, etc.
	TableName   string
	Alias       string
	OnCondition parser.IExpressionContext
}

type SetOperationType int

const (
	UnionOperation SetOperationType = iota
	IntersectOperation
	ExceptOperation
)

func (s SetOperationType) String() string {
	switch s {
	case UnionOperation:
		return "UNION"
	case IntersectOperation:
		return "INTERSECT"
	case ExceptOperation:
		return "EXCEPT"
	default:
		return "UNKNOWN"
	}
}

// SelectStatement - главный узел для любого SELECT запроса
// Может быть либо простым (Primary) либо комбинированным (Combined)
type SelectStatement struct {
	Primary      *PrimarySelectStatement
	Combined     *CombinedSelectStatement
	Subqueries   []*SelectStatement // вложенные SELECT-запросы
	SubqueryCache map[*SelectStatement]interface{} // кэш для скалярных подзапросов
}

func (s *SelectStatement) Type() string { return "SELECT" }

func (s *SelectStatement) IsPrimary() bool {
	return s.Primary != nil
}

func (s *SelectStatement) IsCombined() bool {
	return s.Combined != nil
}

// PrimarySelectStatement - один SELECT запрос со всеми клаузами
type PrimarySelectStatement struct {
	Columns []items.SelectItem
	From    *FromClause
	Where   *items.WhereClause
	GroupBy []items.SelectItem
	OrderBy []items.OrderByItem
	Limit   int
	Offset  int
}

// CombinedSelectStatement - объединение двух SELECT через UNION/INTERSECT/EXCEPT
type CombinedSelectStatement struct {
	OperationType SetOperationType
	All           bool // true для UNION ALL, INTERSECT ALL, EXCEPT ALL
	Left          *SelectStatement
	Right         *SelectStatement
}
