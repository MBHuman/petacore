package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rops"
	"petacore/internal/runtime/rsql/table"
	"sort"
)

// parseMultiplicativeExpression handles multiplication and division
func ParseMultiplicativeExpression(multExpr parser.IMultiplicativeExpressionContext, row *table.ResultRow) (result rmodels.Expression, err error) {
	// logger.Debug("ParseMultiplicativeExpression")
	if multExpr == nil {
		return nil, nil
	}

	// Get all unary expressions
	unaryExprs := multExpr.AllUnaryExpression()
	if len(unaryExprs) == 0 {
		return nil, nil
	}

	// Get the first unary expression
	result, err = ParseUnaryExpression(unaryExprs[0], row)
	if err != nil {
		return nil, err
	}

	// Handle additional terms with operators
	stars := multExpr.AllSTAR()
	slashes := multExpr.AllSLASH()

	// Merge operators by token index
	type opInfo struct {
		op    string
		index int
	}
	var ops []opInfo
	for _, s := range stars {
		ops = append(ops, opInfo{"*", s.GetSymbol().GetTokenIndex()})
	}
	for _, s := range slashes {
		ops = append(ops, opInfo{"/", s.GetSymbol().GetTokenIndex()})
	}
	// Sort by token index
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].index < ops[j].index
	})

	// Skip the first one since we already processed it
	for i, op := range ops {
		if i+1 >= len(unaryExprs) {
			break
		}
		nextValue, err := ParseUnaryExpression(unaryExprs[i+1], row)
		if err != nil {
			return nil, err
		}

		switch op.op {
		case "*":
			if val, ok := result.(*rmodels.ResultRowsExpression); ok {
				if nextVal, ok := nextValue.(*rmodels.ResultRowsExpression); ok {
					result, err = rops.MultiplyValues(val, nextVal)
				} else {
					return nil, fmt.Errorf("multiplication is only supported for result row expressions")
				}
			} else {
				return nil, fmt.Errorf("multiplication is only supported for result row expressions")
			}
		case "/":
			if val, ok := result.(*rmodels.ResultRowsExpression); ok {
				if nextVal, ok := nextValue.(*rmodels.ResultRowsExpression); ok {
					result, err = rops.DivideValues(val, nextVal)
				} else {
					return nil, fmt.Errorf("division is only supported for result row expressions")
				}
			} else {
				return nil, fmt.Errorf("division is only supported for result row expressions")
			}
		}
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
