package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/statements"
	"petacore/sdk/pmem"
)

// sqlListener реализует ANTLR listener для разбора SQL
type sqlListener struct {
	*parser.BasesqlListener
	allocator pmem.Allocator
	stmt      statements.SQLStatement
	subExec   subquery.SubqueryExecutor
	err       error
	// subqueries collects parsed sub-select statements keyed by their SQL text
	subqueries map[string]*statements.SelectStatement
}
