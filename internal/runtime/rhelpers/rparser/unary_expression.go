package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rops"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
)

// ParseUnaryExpression handles unary operators (+ and -)
func ParseUnaryExpression(
	allocator pmem.Allocator,
	ctx context.Context,
	unaryExpr parser.IUnaryExpressionContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (result rmodels.Expression, err error) {
	if unaryExpr == nil {
		return nil, nil
	}

	// Get the cast expression
	castExpr := unaryExpr.CastExpression()
	if castExpr == nil {
		return nil, nil
	}

	// Parse the cast expression
	result, err = ParseCastExpression(allocator, ctx, castExpr, row, subExec)
	if err != nil {
		return nil, err
	}

	// Check for unary operators
	hasMinus := unaryExpr.MINUS() != nil

	if hasMinus {
		// Apply unary minus
		if val, ok := result.(*rmodels.ResultRowsExpression); ok {
			// TODO перевести на inplace если есть возможность, сейчас это будет создавать новый результат, что не очень эффективно
			if len(val.Row.Rows) == 1 && len(val.Row.Schema.Fields) == 1 {
				result, err = rops.NegateValue(allocator, val.Row.Rows[0], val.Row.Schema)
				if err != nil {
					return nil, fmt.Errorf("[ParseUnaryExpression] NegateValue error: %w", err)
				}
				return result, nil
			} else {
				return nil, fmt.Errorf("[ParseUnaryExpression] unary minus is only supported for single-row single-column result expressions")
			}
		}
	}

	return result, nil
}
