package visitor

import (
	"petacore/internal/runtime/parser"
)

// EnterWhereClause больше не используется
// WHERE обрабатывается в parsePrimarySelectStatement в select.go
func (l *sqlListener) EnterWhereClause(ctx *parser.WhereClauseContext) {
	// No-op: WHERE clause is now handled in parsePrimarySelectStatement
}
