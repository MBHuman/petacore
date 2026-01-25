package revaluate

import (
	"petacore/internal/logger"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/table"
)

func EvaluateWhereCondition(where *items.WhereClause, row *table.ResultRow) bool {
	if where == nil {
		return true
	}
	result, err := EvaluateExpressionContext(where.ExpressionContext, row)
	if err != nil {
		logger.Errorf("Error evaluating WHERE condition: %v", err)
		return false
	}
	logger.Debugf("WHERE condition result: %v (type: %T)", result, result)

	if boolVal, ok := result.(*rmodels.BoolExpression); ok {
		logger.Debugf("BoolExpression value: %v", boolVal.Value)
		return boolVal.Value
	}

	// Если это ResultRowsExpression, пробуем извлечь булево значение
	if resultVal, ok := result.(*rmodels.ResultRowsExpression); ok {
		if len(resultVal.Row.Rows) > 0 && len(resultVal.Row.Rows[0]) > 0 {
			val := resultVal.Row.Rows[0][0]
			logger.Debugf("ResultRowsExpression value: %v (type: %T)", val, val)
			if boolVal, ok := val.(bool); ok {
				return boolVal
			}
		}
	}

	logger.Debugf("WHERE condition returned non-bool result, treating as false")
	return false
}
