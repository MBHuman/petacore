package revaluate

import (
	"context"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rparser"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
)

func EvaluateComparisonExpressionContext(goCtx context.Context, ctx parser.IComparisonExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (interface{}, error) {
	parsed, err := rparser.ParseComparisonExpression(goCtx, ctx, row, subExec)
	return parsed, err
}
