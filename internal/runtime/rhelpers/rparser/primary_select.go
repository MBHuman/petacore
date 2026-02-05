package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"strconv"
)

// parsePrimarySelectStatement парсит один PRIMARY SELECT или скобки
func ParsePrimarySelectStatement(ctx parser.IPrimarySelectStatementContext) (stmt *statements.SelectStatement, err error) {
	if ctx == nil {
		return nil, nil
	}

	// Проверяем, это скобки или реальный SELECT
	if ctx.LPAREN() != nil && ctx.SelectStatement() != nil {
		// Это (SELECT ...), рекурсивно парсим
		if selectStmt, ok := ctx.SelectStatement().(*parser.SelectStatementContext); ok {
			return ParseSelectStatement(selectStmt)
		}
	}

	// Это реальный PRIMARY SELECT
	primary := &statements.PrimarySelectStatement{}

	// FROM clause
	if ctx.FromClause() != nil {
		primary.From, err = ParseFromClause(ctx.FromClause())
		if err != nil {
			return nil, fmt.Errorf("error parsing Primary SELECT FROM clause: %v", err)
		}
	}

	// SELECT list (columns)
	if ctx.SelectList() != nil {
		primary.Columns = ParseSelectList(ctx.SelectList())
	}

	// WHERE clause
	if ctx.WhereClause() != nil {
		primary.Where = &items.WhereClause{
			ExpressionContext: ctx.WhereClause().Expression(),
		}
	}

	// GROUP BY clause
	if ctx.GroupByClause() != nil {
		primary.GroupBy, err = ParseGroupByClause(ctx.GroupByClause())
		if err != nil {
			return nil, fmt.Errorf("error parsing GROUP BY clause: %v", err)
		}
	}

	// ORDER BY clause
	if ctx.OrderByClause() != nil {
		primary.OrderBy = ParseOrderByClause(ctx.OrderByClause())
	}

	// LIMIT
	if ctx.LimitValue() != nil {
		limitStr := ctx.LimitValue().GetText()
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid LIMIT value: %v", err)
		}
		primary.Limit = limit
	}

	// OFFSET
	if ctx.OffsetValue() != nil {
		offsetStr := ctx.OffsetValue().GetText()
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid OFFSET value: %v", err)
		}
		primary.Offset = offset
	}

	return &statements.SelectStatement{
		Primary:       primary,
		SubqueryCache: make(map[*statements.SelectStatement]interface{}),
	}, nil
}
