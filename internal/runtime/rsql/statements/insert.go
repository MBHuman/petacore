package statements

// InsertStatement представляет INSERT
type InsertStatement struct {
	TableName string
	Columns   []string
	Values    [][]interface{} // Multiple rows
}

func (i *InsertStatement) Type() string { return "INSERT" }
