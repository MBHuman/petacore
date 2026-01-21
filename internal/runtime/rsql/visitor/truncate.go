package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
)

func (l *sqlListener) EnterTruncateTableStatement(ctx *parser.TruncateTableStatementContext) {
	stmt := &statements.TruncateTableStatement{}
	if ctx.TableName() != nil {
		stmt.TableName = ctx.TableName().GetText()
	}
	l.stmt = stmt
}
