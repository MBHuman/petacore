package rparser

import (
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
)

// parseExtractFunction handles EXTRACT expressions
func ParseExtractFunction(extractExpr parser.IExtractFunctionContext, row *table.ResultRow) (rmodels.Expression, error) {
	if extractExpr == nil {
		return nil, nil
	}

	field := extractExpr.IDENTIFIER().GetText()
	sourceExpr := extractExpr.Expression()
	if sourceExpr == nil {
		return nil, nil
	}

	source, err := ParseExpression(sourceExpr, row)
	if err != nil {
		return nil, err
	}
	args := []interface{}{field, source}
	value, _ := functions.ExecuteFunction("EXTRACT", args)
	return &rmodels.ResultRowsExpression{Row: value}, nil
}
