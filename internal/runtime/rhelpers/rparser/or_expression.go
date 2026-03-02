package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
)

// parseOrExpression handles OR expressions
func ParseOrExpression(allocator pmem.Allocator, orExpr parser.IOrExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	return ParseOrExpressionWithContext(allocator, context.Background(), orExpr, row, subExec)
}

// ParseOrExpressionWithContext парсит OR выражение с контекстом
func ParseOrExpressionWithContext(allocator pmem.Allocator, ctx context.Context, orExpr parser.IOrExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	if orExpr == nil {
		return nil, nil
	}

	andExprs := orExpr.AllAndExpression()
	if len(andExprs) == 0 {
		return nil, nil
	}

	// Evaluate first AND expression
	result, err := ParseAndExpression(allocator, ctx, andExprs[0], row, subExec)
	if err != nil {
		return nil, err
	}

	if leftVal, ok := result.(*rmodels.BoolExpression); ok {
		result = leftVal

		// If multiple AND expressions connected by OR
		for i := 1; i < len(andExprs); i++ {
			rightVal, err := ParseAndExpression(allocator, ctx, andExprs[i], row, subExec)
			if err != nil {
				return nil, err
			}

			if rv, ok := rightVal.(*rmodels.BoolExpression); ok {
				leftVal.Value = leftVal.Value || rv.Value
			}
		}
		return leftVal, nil
	} else if leftVal, ok := result.(*rmodels.ResultRowsExpression); ok {
		return leftVal, nil
	} else if leftVal, ok := result.(*rmodels.CaseExpression); ok {
		return leftVal, nil
	} else if leftVal, ok := result.(*rmodels.SubqueryExpression); ok {
		return leftVal, nil
	} else {
		return nil, fmt.Errorf("expected BoolExpression, got %T", result)
	}
}
