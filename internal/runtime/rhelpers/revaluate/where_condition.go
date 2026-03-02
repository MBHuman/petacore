package revaluate

import (
	"context"
	"petacore/internal/logger"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func EvaluateWhereCondition(
	allocator pmem.Allocator,
	goCtx context.Context,
	where *items.WhereClause,
	row *table.ResultRow,
	statement *statements.SelectStatement,
	subExec subquery.SubqueryExecutor,
	runtimeParams map[int]interface{},
) bool {
	if where == nil {
		return true
	}
	result, err := EvaluateExpressionContext(allocator, goCtx, where.ExpressionContext, row, subExec, runtimeParams)
	if err != nil {
		logger.Errorf("[EvaluateWhereCondition] Error evaluating WHERE condition: %v", err)
		return false
	}

	if boolVal, ok := result.(*rmodels.BoolExpression); ok {
		return boolVal.Value
	}

	// Если это ResultRowsExpression, пробуем извлечь булево значение
	if resultVal, ok := result.(*rmodels.ResultRowsExpression); ok {
		var err error
		val, oid, err := resultVal.Row.Schema.GetField(resultVal.Row.Rows[0], 0)
		if err != nil {
			logger.Errorf("[EvaluateWhereCondition] Error getting field from ResultRowsExpression: %v", err)
			return false
		}
		valDes, err := serializers.DeserializeGeneric(val, oid)
		if err != nil {
			logger.Errorf("[EvaluateWhereCondition] Error deserializing field from ResultRowsExpression: %v", err)
			return false
		}
		if oid == ptypes.PTypeBool {
			valueBool, ok := ptypes.TryIntoBool(valDes)
			if ok {
				return valueBool.IntoGo()
			}
		}
	}

	return false
}

// getSubqueryCache извлекает кэш подзапросов из row, если возможно
func getSubqueryCache(statement *statements.SelectStatement) map[*statements.SelectStatement]*ptypes.Row {
	if statement == nil {
		return nil
	}
	return statement.SubqueryCache
}
