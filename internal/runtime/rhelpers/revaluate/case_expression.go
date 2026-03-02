package revaluate

import (
	"context"
	"fmt"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

// evaluateCaseExpression evaluates a CASE WHEN THEN ELSE END expression
func EvaluateCaseExpression(
	allocator pmem.Allocator,
	goCtx context.Context,
	caseExpr *rmodels.CaseExpression,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
	runtimeParams map[int]interface{},
) (rmodels.Expression, error) {
	ctx := caseExpr.Context
	// In new grammar: CASE (WHEN expression THEN expression)+ (ELSE expression)? END
	// AllExpression returns all expressions: WHEN1, THEN1, WHEN2, THEN2, ..., ELSE (if present)
	allExpressions := ctx.AllExpression()
	numWhen := len(ctx.AllWHEN())

	if numWhen == 0 {
		return nil, fmt.Errorf("[EvaluateCaseExpression] invalid CASE expression")
	}

	// Expressions alternate: WHEN1, THEN1, WHEN2, THEN2, ...
	for i := 0; i < numWhen; i++ {
		whenIdx := i * 2
		thenIdx := i*2 + 1

		if whenIdx >= len(allExpressions) {
			break
		}

		// Evaluate WHEN condition
		condition, err := EvaluateExpressionContext(allocator, goCtx, allExpressions[whenIdx], row, subExec, runtimeParams)
		if err != nil {
			return nil, err
		}

		// Check if true - нужно правильно извлечь булево значение
		isTrue := false
		if boolExpr, ok := condition.(*rmodels.BoolExpression); ok {
			isTrue = boolExpr.Value
		} else if resultExpr, ok := condition.(*rmodels.ResultRowsExpression); ok {
			// Извлекаем значение из ResultRowsExpression
			if len(resultExpr.Row.Rows) > 0 {
				val, oid, err := resultExpr.Row.Schema.GetField(resultExpr.Row.Rows[0], 0)
				if err != nil {
					return nil, err
				}
				value, err := serializers.DeserializeGeneric(val, oid)
				if err != nil {
					return nil, err
				}

				valueBool, ok := ptypes.TryIntoBool(value)
				if ok {
					isTrue = valueBool.IntoGo()
				} else {
					// TODO возможно тут падать будет
					return nil, fmt.Errorf("[EvaluateCaseExpression] CASE WHEN condition does not evaluate to boolean")
				}
			}
		}

		if isTrue {
			if thenIdx < len(allExpressions) {
				// Evaluate THEN result
				result, err := EvaluateExpressionContext(allocator, goCtx, allExpressions[thenIdx], row, subExec, runtimeParams)
				if err != nil {
					return nil, err
				}
				return result, nil
			}
		}
	}

	// Check for ELSE
	elseIdx := numWhen * 2
	if elseIdx < len(allExpressions) {
		result, err := EvaluateExpressionContext(allocator, goCtx, allExpressions[elseIdx], row, subExec, runtimeParams)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	// No match and no ELSE
	return nil, fmt.Errorf("[EvaluateCaseExpression] CASE expression has no matching WHEN condition and no ELSE clause")
}
