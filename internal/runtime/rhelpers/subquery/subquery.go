package subquery

import (
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
)

// SubqueryExecutor executes a parsed SELECT statement and returns its result.
type SubqueryExecutor func(*statements.SelectStatement) (*table.ExecuteResult, error)
