package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
)

// parseNotExpression handles NOT expression
func ParseNotExpression(notExpr parser.INotExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	// logger.Debug("ParseNotExpression")
	if notExpr == nil {
		return nil, nil
	}

	compExpr := notExpr.ComparisonExpression()
	if compExpr == nil {
		return nil, nil
	}

	result, err := ParseComparisonExpression(compExpr, row)
	if err != nil {
		return nil, err
	}

	// Apply NOT if present
	if notExpr.NOT() != nil {
		if resultBool, ok := result.(*rmodels.BoolExpression); ok {
			return &rmodels.BoolExpression{Value: !resultBool.Value}, nil
		} else {
			return nil, fmt.Errorf("expected BoolExpression, got %T", result)
		}
	}

	return result, nil
}
