package rparser

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
)

// ParseUnaryExpression handles unary operators (+ and -)
func ParseUnaryExpression(unaryExpr parser.IUnaryExpressionContext, row *table.ResultRow) (result rmodels.Expression, err error) {
	if unaryExpr == nil {
		return nil, nil
	}

	// Get the cast expression
	castExpr := unaryExpr.CastExpression()
	if castExpr == nil {
		return nil, nil
	}

	// Parse the cast expression
	result, err = ParseCastExpression(castExpr, row)
	if err != nil {
		return nil, err
	}

	// Check for unary operators
	hasMinus := unaryExpr.MINUS() != nil
	// hasPlus := unaryExpr.PLUS() != nil

	if hasMinus {
		// Apply unary minus
		if val, ok := result.(*rmodels.ResultRowsExpression); ok {
			if len(val.Row.Rows) > 0 && len(val.Row.Rows[0]) > 0 {
				value := val.Row.Rows[0][0]
				switch v := value.(type) {
				case int:
					val.Row.Rows[0][0] = -v
				case int32:
					val.Row.Rows[0][0] = -v
				case int64:
					val.Row.Rows[0][0] = -v
				case float32:
					val.Row.Rows[0][0] = -v
				case float64:
					val.Row.Rows[0][0] = -v
				}
			}
		}
	}
	// For unary plus, we don't need to do anything

	return result, nil
}
