package statements

import "petacore/internal/runtime/rsql/items"

type SelectStatement struct {
	TableName  string
	TableAlias string
	Columns    []items.SelectItem
	Where      *items.WhereClause
	OrderBy    []items.OrderByItem
	Limit      int
}

func (s *SelectStatement) Type() string { return "SELECT" }
