package items

import "petacore/internal/runtime/parser"

type SelectItem struct {
	IsSelectAll       bool
	TableAlias        string
	ColumnName        string
	Function          *FunctionCall
	ExpressionContext parser.IExpressionContext
	Alias             string

	// Можно использовать для приведения типов подряд в SELECT
	// Например: SELECT col1::int::text
	TypeCasting []string
}

// WhereClause представляет WHERE условие
type WhereClause struct {
	ExpressionContext parser.IExpressionContext
}

type OrderDirection int

const (
	OrderAsc OrderDirection = iota
	OrderDesc
)

type OrderByItem struct {
	ColumnName  string
	ColumnIndex int
	Direction   OrderDirection
}

type FunctionCall struct {
	Name        string
	Args        []parser.IFunctionArgContext
	IsAggregate bool
}
