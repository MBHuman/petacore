package revaluate

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rparser"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

// evaluateExpressionContext evaluates an expression using the ANTLR parsed context
// Может возвращать bool для условий и *table.ExecuteResult для других выражений
func EvaluateExpressionContext(allocator pmem.Allocator, goCtx context.Context, ctx parser.IExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor, runtimeParams map[int]interface{}) (rmodels.Expression, error) {
	parsed, err := rparser.ParseExpression(allocator, goCtx, ctx, row, subExec)
	if err != nil {
		return nil, err
	}
	// Если это ParamRefExpression — возвращаем значение параметра
	if paramRef, ok := parsed.(*rmodels.ParamRefExpression); ok {
		var val interface{} = nil
		if runtimeParams != nil {
			val = runtimeParams[paramRef.Index]
		}
		fields := []serializers.FieldDef{{
			Name: "?param?",
			OID:  ptypes.PTypeText,
		}}
		schema := serializers.NewBaseSchema(fields)
		outRow, err := schema.Pack(allocator, [][]byte{[]byte(fmt.Sprintf("%v", val))})
		if err != nil {
			return nil, err
		}
		return &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:   []*ptypes.Row{outRow},
				Schema: schema,
			},
		}, nil
	}
	if caseExpr, ok := parsed.(*rmodels.CaseExpression); ok {
		return EvaluateCaseExpression(allocator, goCtx, caseExpr, row, subExec, runtimeParams)
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
