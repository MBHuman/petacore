package statements

import "petacore/internal/runtime/rsql/table"

// CreateTableStatement представляет CREATE TABLE
type CreateTableStatement struct {
	TableName   string
	Columns     []table.ColumnDef
	PrimaryKeys []int
	IfNotExists bool
}

func (c *CreateTableStatement) Type() string { return "CREATE_TABLE" }
