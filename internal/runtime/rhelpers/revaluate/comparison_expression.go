package revaluate

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rparser"
	"petacore/internal/runtime/rsql/table"
)

func EvaluateComparisonExpressionContext(ctx parser.IComparisonExpressionContext, row *table.ResultRow) (interface{}, error) {
	parsed, err := rparser.ParseComparisonExpression(ctx, row)
	return parsed, err
}
