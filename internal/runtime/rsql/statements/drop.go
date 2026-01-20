package statements

// DropTableStatement представляет DROP TABLE
type DropTableStatement struct {
	TableName string
}

func (d *DropTableStatement) Type() string { return "DROP_TABLE" }
