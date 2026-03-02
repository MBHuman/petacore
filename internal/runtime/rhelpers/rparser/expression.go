package rparser

import (
	"context"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
)

// parseExpression evaluates an ANTLR expression context and returns the value
func ParseExpression(allocator pmem.Allocator, ctx context.Context, expr parser.IExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	if expr == nil {
		return nil, nil
	}

	if orExpr := expr.OrExpression(); orExpr != nil {
		return ParseOrExpressionWithContext(allocator, ctx, orExpr, row, subExec)
	}

	return nil, nil
}
