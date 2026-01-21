package statements

import "petacore/internal/runtime/rsql/items"

type FromClause struct {
	TableName string
	Alias     string
	Joins     []JoinClause
}

type JoinClause struct {
	Type        string // INNER, LEFT, etc.
	TableName   string
	Alias       string
	OnCondition interface{} // *parser.IExpressionContext
}

type SelectStatement struct {
	TableName  string
	TableAlias string
	From       *FromClause
	Columns    []items.SelectItem
	Where      *items.WhereClause
	OrderBy    []items.OrderByItem
	Limit      int
}

func (s *SelectStatement) Type() string { return "SELECT" }
