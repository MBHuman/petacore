package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
)

// parseExtractFunction handles EXTRACT expressions
func ParseExtractFunction(allocator pmem.Allocator, ctx context.Context, extractExpr parser.IExtractFunctionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	if extractExpr == nil {
		return nil, nil
	}

	field := extractExpr.IDENTIFIER().GetText()
	sourceExpr := extractExpr.Expression()
	if sourceExpr == nil {
		return nil, nil
	}

	source, err := ParseExpression(allocator, ctx, sourceExpr, row, subExec)
	if err != nil {
		return nil, err
	}
	args := []interface{}{field, source}
	value, err := functions.ExecuteFunction(allocator, "EXTRACT", args)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("extract function returned nil")
	}
	return &rmodels.ResultRowsExpression{Row: value}, nil
}
