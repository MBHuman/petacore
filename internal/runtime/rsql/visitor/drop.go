package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
)

func (l *sqlListener) EnterDropTableStatement(ctx *parser.DropTableStatementContext) {
	stmt := &statements.DropTableStatement{}
	if ctx.TableName() != nil {
		stmt.TableName = ctx.TableName().GetText()
	}
	l.stmt = stmt
}
