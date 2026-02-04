package statements

// TruncateTableStatement представляет TRUNCATE TABLE
type TruncateTableStatement struct {
	TableName string
}

func (t *TruncateTableStatement) Type() string { return "TRUNCATE_TABLE" }
