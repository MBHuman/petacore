package statements

// DescribeStatement представляет DESCRIBE TABLE
type DescribeStatement struct {
	TableName string
}

func (d *DescribeStatement) Type() string { return "DESCRIBE" }
