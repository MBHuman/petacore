package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
)

// parseOrExpression handles OR expressions
func ParseOrExpression(orExpr parser.IOrExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	// logger.Debug("ParseOrExpression")
	if orExpr == nil {
		return nil, nil
	}

	andExprs := orExpr.AllAndExpression()
	if len(andExprs) == 0 {
		return nil, nil
	}

	// Evaluate first AND expression
	result, err := ParseAndExpression(andExprs[0], row)
	if err != nil {
		return nil, err
	}

	if leftVal, ok := result.(*rmodels.BoolExpression); ok {
		result = leftVal

		// If multiple AND expressions connected by OR
		for i := 1; i < len(andExprs); i++ {
			rightVal, err := ParseAndExpression(andExprs[i], row)
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
	} else {
		return nil, fmt.Errorf("expected BoolExpression, got %T", result)
	}
}
