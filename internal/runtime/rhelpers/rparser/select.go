package rparser

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
)

// parseSelectStatement парсит весь SELECT statement (может быть primary или combined)
func ParseSelectStatement(ctx parser.ISelectStatementContext) (*statements.SelectStatement, error) {
	if ctx == nil {
		return nil, nil
	}

	unionExceptCtx := ctx.UnionExceptStatement()
	if unionExceptCtx == nil {
		return nil, nil
	}

	return ParseUnionExceptStatement(unionExceptCtx)
}
