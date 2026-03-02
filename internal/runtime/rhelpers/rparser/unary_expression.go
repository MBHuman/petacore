package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
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
				fieldVal, oid, err := val.Row.Schema.GetField(val.Row.Rows[0], 0)
				if err != nil {
					return nil, err
				}
				desVal, err := serializers.DeserializeGeneric(fieldVal, oid)
				if err != nil {
					return nil, err
				}
				negated, err := ptypes.ApplyNeg(allocator, desVal, oid)
				if err != nil {
					return nil, err
				}
				resultRow, err := val.Row.Schema.Pack(allocator, [][]byte{negated.GetBuffer()})
				if err != nil {
					return nil, err
				}
				result = &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows:   []*ptypes.Row{resultRow},
						Schema: val.Row.Schema,
					},
				}
			} else {
				return nil, fmt.Errorf("unary minus is only supported for single-row single-column result expressions")
			}
		}
	}

	return result, nil
}
