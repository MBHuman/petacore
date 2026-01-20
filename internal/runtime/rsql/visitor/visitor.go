package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
)

// sqlListener реализует ANTLR listener для разбора SQL
type sqlListener struct {
	*parser.BasesqlListener
	stmt statements.SQLStatement
	err  error
}
