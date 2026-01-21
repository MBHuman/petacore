package items

import "petacore/internal/runtime/parser"

type SelectItem struct {
	TableAlias        string
	ColumnName        string
	Function          *FunctionCall
	ExpressionContext parser.IExpressionContext
	Alias             string
}

// WhereClause представляет WHERE условие
type WhereClause struct {
	ExpressionContext parser.IExpressionContext
}

type OrderByItem struct {
	ExpressionContext parser.IExpressionContext
	Direction         string // "ASC" or "DESC"
}

type FunctionCall struct {
	Name string
	Args []interface{}
}

type CaseExpression struct {
	Context parser.ICaseExpressionContext
}
