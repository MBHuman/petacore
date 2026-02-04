package rparser

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
)

// ParseConcatExpression handles string concatenation with ||
func ParseConcatExpression(concatExpr parser.IConcatExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	// logger.Debug("ParseConcatExpression")
	if concatExpr == nil {
		return nil, nil
	}

	// Get all additive expressions to concatenate
	exprs := concatExpr.AllAdditiveExpression()
	if len(exprs) == 0 {
		return nil, nil
	}

	// If only one expression and no CONCAT operators, return it directly
	if len(exprs) == 1 && len(concatExpr.AllCONCAT()) == 0 {
		return ParseAdditiveExpression(exprs[0], row)
	}

	// Multiple expressions with CONCAT, concatenate as strings
	result := &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows:    [][]interface{}{{""}},
			Columns: []table.TableColumn{{Type: table.ColTypeString}},
		},
	}
	for _, e := range exprs {
		val, err := ParseAdditiveExpression(e, row)
		if err != nil {
			return nil, err
		}
		if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok {
			err = checkConcatExpr(result, valExpr)
			if err != nil {
				return nil, err
			}
			if val != nil {
				result.Row.Rows[0][0] = fmt.Sprintf("%v", result.Row.Rows[0][0]) + fmt.Sprintf("%v", valExpr.Row.Rows[0][0])
			}
		} else {
			return nil, fmt.Errorf("concatenation supports only result row expressions")
		}

	}
	logger.Debugf("Concatenated result: %s\n", result)
	return result, nil
}

func checkConcatExpr(a, b *rmodels.ResultRowsExpression) error {
	if a == nil || b == nil {
		return fmt.Errorf("nil operand in concatenation")
	}
	if len(a.Row.Rows) == 0 || len(b.Row.Rows) == 0 {
		return fmt.Errorf("empty rows in concatenation")
	}
	if len(a.Row.Columns) == 0 || len(b.Row.Columns) == 0 {
		return fmt.Errorf("empty columns in concatenation")
	}
	if len(a.Row.Rows) > 1 || len(b.Row.Rows) > 1 {
		return fmt.Errorf("concatenation supports only single-row expressions")
	}
	if a.Row.Columns[0].Type != table.ColTypeString || b.Row.Columns[0].Type != table.ColTypeString {
		return fmt.Errorf("concatenation supports only string types")
	}
	return nil
}
