package revaluate

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
)

// evaluateCaseExpression evaluates a CASE WHEN THEN ELSE END expression
func EvaluateCaseExpression(goCtx context.Context, caseExpr *rmodels.CaseExpression, row *table.ResultRow, subExec subquery.SubqueryExecutor, runtimeParams map[int]interface{}) (rmodels.Expression, error) {
	ctx := caseExpr.Context
	// In new grammar: CASE (WHEN expression THEN expression)+ (ELSE expression)? END
	// AllExpression returns all expressions: WHEN1, THEN1, WHEN2, THEN2, ..., ELSE (if present)
	allExpressions := ctx.AllExpression()
	numWhen := len(ctx.AllWHEN())

	if numWhen == 0 {
		return nil, fmt.Errorf("invalid CASE expression")
	}

	// Expressions alternate: WHEN1, THEN1, WHEN2, THEN2, ...
	for i := 0; i < numWhen; i++ {
		whenIdx := i * 2
		thenIdx := i*2 + 1

		if whenIdx >= len(allExpressions) {
			break
		}

		// Evaluate WHEN condition
		condition, err := EvaluateExpressionContext(goCtx, allExpressions[whenIdx], row, subExec, runtimeParams)
		if err != nil {
			return nil, err
		}

		// Check if true - нужно правильно извлечь булево значение
		isTrue := false
		if boolExpr, ok := condition.(*rmodels.BoolExpression); ok {
			isTrue = boolExpr.Value
		} else if resultExpr, ok := condition.(*rmodels.ResultRowsExpression); ok {
			// Извлекаем значение из ResultRowsExpression
			if len(resultExpr.Row.Rows) > 0 && len(resultExpr.Row.Rows[0]) > 0 {
				val := resultExpr.Row.Rows[0][0]
				if boolVal, ok := val.(bool); ok {
					isTrue = boolVal
				} else {
					// Для других типов используем IsTrue
					isTrue = rhelpers.IsTrue(val)
				}
			}
		} else {
			// Fallback - используем IsTrue для других типов
			isTrue = rhelpers.IsTrue(condition)
		}

		if isTrue {
			if thenIdx < len(allExpressions) {
				// Evaluate THEN result
				result, err := EvaluateExpressionContext(goCtx, allExpressions[thenIdx], row, subExec, runtimeParams)
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
		result, err := EvaluateExpressionContext(goCtx, allExpressions[elseIdx], row, subExec, runtimeParams)
		if err != nil {
			return nil, err
		}
		logger.Debugf("CASE ELSE result: %v", result)
		return result, nil
	}

	// No match and no ELSE
	return nil, nil
}
