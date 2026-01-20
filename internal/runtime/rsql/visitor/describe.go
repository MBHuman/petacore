package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
)

func (l *sqlListener) EnterDescribeStatement(ctx *parser.DescribeStatementContext) {
	stmt := &statements.DescribeStatement{}
	if ctx.TableName() != nil {
		stmt.TableName = ctx.TableName().GetText()
	}
	l.stmt = stmt
}
