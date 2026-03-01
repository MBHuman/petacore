package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rops"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"sort"
)

// parseAdditiveExpression handles addition and subtraction
func ParseAdditiveExpression(ctx context.Context, addExpr parser.IAdditiveExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (result rmodels.Expression, err error) {
	if addExpr == nil {
		return nil, nil
	}

	// Get the first multiplicative expression
	multExpr := addExpr.MultiplicativeExpression(0)
	if multExpr == nil {
		return nil, nil
	}

	result, err = ParseMultiplicativeExpression(ctx, multExpr, row, subExec)
	if err != nil {
		return nil, err
	}

	// Handle subsequent multiplicative expressions with operators
	plusOps := addExpr.AllPLUS()
	minusOps := addExpr.AllMINUS()

	// Build operator order by token index
	type opInfo struct {
		op    string
		index int
	}
	var ops []opInfo
	for _, p := range plusOps {
		ops = append(ops, opInfo{"+", p.GetSymbol().GetTokenIndex()})
	}
	for _, m := range minusOps {
		ops = append(ops, opInfo{"-", m.GetSymbol().GetTokenIndex()})
	}
	// Sort by token index
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].index < ops[j].index
	})

	// Simple stable order since grammar preserves order of terms
	multExprs := addExpr.AllMultiplicativeExpression()
	// Skip the first one since we already processed it
	for i, op := range ops {
		if i+1 >= len(multExprs) {
			break
		}
		nextValue, err := ParseMultiplicativeExpression(ctx, multExprs[i+1], row, subExec)
		if err != nil {
			return nil, err
		}

		switch op.op {
		case "+":
			if val, ok := result.(*rmodels.ResultRowsExpression); ok {
				if nextVal, ok := nextValue.(*rmodels.ResultRowsExpression); ok {
					result, err = rops.AddValues(val, nextVal)
				} else {
					return nil, fmt.Errorf("addition is only supported for result row expressions")
				}
			} else {
				return nil, fmt.Errorf("addition is only supported for result row expressions")
			}
		case "-":
			if val, ok := result.(*rmodels.ResultRowsExpression); ok {
				if nextVal, ok := nextValue.(*rmodels.ResultRowsExpression); ok {
					result, err = rops.SubtractValues(val, nextVal)
				} else {
					return nil, fmt.Errorf("subtraction is only supported for result row expressions")
				}
			} else {
				return nil, fmt.Errorf("subtraction is only supported for result row expressions")
			}
		}
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
