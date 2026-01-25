package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
)

// parseAndExpression handles AND expressions
func ParseAndExpression(andExpr parser.IAndExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	// logger.Debug("ParseAndExpression")
	if andExpr == nil {
		return nil, nil
	}

	notExprs := andExpr.AllNotExpression()
	if len(notExprs) == 0 {
		return nil, nil
	}

	// Evaluate first NOT expression
	result, err := ParseNotExpression(notExprs[0], row)
	if err != nil {
		return nil, err
	}
	if leftVal, ok := result.(*rmodels.BoolExpression); ok {
		// If multiple NOT expressions connected by AND
		for i := 1; i < len(notExprs); i++ {
			rightVal, err := ParseNotExpression(notExprs[i], row)
			if err != nil {
				return nil, err
			}
			if rv, ok := rightVal.(*rmodels.BoolExpression); ok {
				leftVal.Value = leftVal.Value && rv.Value
			} else {
				return nil, fmt.Errorf("expected BoolExpression, got %T", rightVal)
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
