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

// parseNotExpression handles NOT expression
func ParseNotExpression(allocator pmem.Allocator, ctx context.Context, notExpr parser.INotExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	if notExpr == nil {
		return nil, nil
	}

	if notExpr.SubqueryExpression() != nil {
		subqCtx := notExpr.SubqueryExpression()
		selCtx := subqCtx.SelectStatement()
		if selCtx == nil {
			return nil, fmt.Errorf("[ParseNotExpression] expected select statement in NOT EXISTS subquery")
		}
		selectStmt, err := ParseSelectStatement(selCtx)
		if err != nil {
			return nil, fmt.Errorf("[ParseNotExpression] error parsing subquery in NOT EXISTS operator: %w", err)
		}
		res, err := subExec(selectStmt)
		if err != nil {
			return nil, err
		}
		exists := len(res.Rows) > 0
		if notExpr.NOT() != nil {
			exists = !exists
		}
		return &rmodels.BoolExpression{Value: exists}, nil
	}

	compExpr := notExpr.ComparisonExpression()
	if compExpr == nil {
		return nil, nil
	}

	result, err := ParseComparisonExpression(allocator, ctx, compExpr, row, subExec)
	if err != nil {
		return nil, err
	}

	// Apply NOT if present
	if notExpr.NOT() != nil {
		if resultBool, ok := result.(*rmodels.BoolExpression); ok {
			return &rmodels.BoolExpression{Value: !resultBool.Value}, nil
		} else {
			return nil, fmt.Errorf("[ParseNotExpression] expected BoolExpression, got %T", result)
		}
	}

	return result, nil
}
