package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
)

func (l *sqlListener) EnterWhereClause(ctx *parser.WhereClauseContext) {
	if ctx.Expression() != nil {
		where := &items.WhereClause{
			ExpressionContext: ctx.Expression(),
		}

		if selectStmt, ok := l.stmt.(*statements.SelectStatement); ok {
			selectStmt.Where = where
		}
	}
}
