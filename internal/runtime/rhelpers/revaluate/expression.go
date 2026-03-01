package revaluate

import (
	"context"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rparser"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
)

// evaluateExpressionContext evaluates an expression using the ANTLR parsed context
// Может возвращать bool для условий и *table.ExecuteResult для других выражений
func EvaluateExpressionContext(goCtx context.Context, ctx parser.IExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor, runtimeParams map[int]interface{}) (rmodels.Expression, error) {
	parsed, err := rparser.ParseExpression(goCtx, ctx, row, subExec)
	if err != nil {
		return nil, err
	}
	// Если это ParamRefExpression — возвращаем значение параметра
	if paramRef, ok := parsed.(*rmodels.ParamRefExpression); ok {
		var val interface{} = nil
		if runtimeParams != nil {
			val = runtimeParams[paramRef.Index]
		}
		// Оборачиваем в ResultRowsExpression для совместимости
		return &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:    [][]interface{}{{val}},
				Columns: []table.TableColumn{{Idx: 0, Name: "?param?", Type: table.ColTypeString}},
			},
		}, nil
	}
	if caseExpr, ok := parsed.(*rmodels.CaseExpression); ok {
		return EvaluateCaseExpression(goCtx, caseExpr, row, subExec, runtimeParams)
	}

	// Если это SubqueryExpression — выполняем подзапрос и возвращаем ExecuteResult
	if subq, ok := parsed.(*rmodels.SubqueryExpression); ok {
		res, err := subExec(subq.Select)
		if err != nil {
			return nil, err
		}
		return &rmodels.ResultRowsExpression{Row: res}, nil
	}
	return parsed, nil
}
