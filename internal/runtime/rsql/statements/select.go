package statements

import "petacore/internal/runtime/rsql/items"

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
	OnCondition interface{} // *parser.IExpressionContext
}

type SelectStatement struct {
	TableAlias string
	From       *FromClause
	Columns    []items.SelectItem
	GroupBy    []items.SelectItem
	Where      *items.WhereClause
	OrderBy    []items.OrderByItem
	Limit      int
}

func (s *SelectStatement) Type() string { return "SELECT" }
