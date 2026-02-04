package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
	"strings"
)

func (l *sqlListener) EnterShowStatement(ctx *parser.ShowStatementContext) {
	stmt := &statements.ShowStatement{}
	// Collect all identifiers
	var parts []string
	for _, id := range ctx.AllIDENTIFIER() {
		parts = append(parts, id.GetText())
	}
	stmt.Parameter = strings.Join(parts, " ")
	l.stmt = stmt
}
