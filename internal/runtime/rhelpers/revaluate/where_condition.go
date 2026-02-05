package revaluate

import (
	"context"
	"petacore/internal/logger"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
)

func EvaluateWhereCondition(goCtx context.Context, where *items.WhereClause, row *table.ResultRow, statement *statements.SelectStatement, subExec subquery.SubqueryExecutor, runtimeParams map[int]interface{}) bool {
	if where == nil {
		return true
	}
	result, err := EvaluateExpressionContext(goCtx, where.ExpressionContext, row, subExec, runtimeParams)
	if err != nil {
		logger.Errorf("Error evaluating WHERE condition: %v", err)
		return false
	}

	if boolVal, ok := result.(*rmodels.BoolExpression); ok {
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

	return false
}

// getSubqueryCache извлекает кэш подзапросов из row, если возможно
func getSubqueryCache(statement *statements.SelectStatement) map[*statements.SelectStatement]interface{} {
	if statement == nil {
		return nil
	}
	return statement.SubqueryCache
}
