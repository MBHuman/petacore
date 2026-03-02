package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

// parsePrimaryExpression handles the basic expressions
// Always returns ResultRowExpression
func ParsePrimaryExpression(
	allocator pmem.Allocator,
	ctx context.Context,
	primExpr parser.IPrimaryExpressionContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (rmodels.Expression, error) {
	// Check for parenthesized expression
	if primExpr.Expression() != nil {
		return ParseExpression(allocator, ctx, primExpr.Expression(), row, subExec)
	}

	// Check for CASE expression
	if primExpr.CaseExpression() != nil {
		return ParseCaseExpression(allocator, ctx, primExpr.CaseExpression())
	}

	// Check for subquery expression
	if primExpr.SubqueryExpression() != nil {
		sqCtx := primExpr.SubqueryExpression()
		selCtx := sqCtx.SelectStatement()
		if selCtx == nil {
			return nil, fmt.Errorf("[ParsePrimaryExpression] invalid subquery context")
		}
		selectStmt, err := ParseSelectStatement(selCtx)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryExpression] error parsing subquery: %v", err)
		}
		if selectStmt == nil {
			return nil, fmt.Errorf("[ParsePrimaryExpression] failed to build select statement from context")
		}
		// Возвращаем выражение-подзапрос
		return &rmodels.SubqueryExpression{Select: selectStmt}, nil
	}

	// Check for function call
	if funcCall := primExpr.FunctionCall(); funcCall != nil {
		return ParsePrimaryExpressionFunctionCall(allocator, ctx, funcCall, row, subExec)
	}

	// Check for extract function
	if primExpr.ExtractFunction() != nil {
		return ParseExtractFunction(allocator, ctx, primExpr.ExtractFunction(), row, subExec)
	}

	if valueLit := primExpr.Value(); valueLit != nil {
		res, err := ParsePrimaryLiteral(allocator, valueLit)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryExpression] error parsing primary literal: %v", err)
		}
		if res != nil {
			return res, nil
		}
	}

	// Check for column name
	if colName := primExpr.ColumnName(); colName != nil {
		return ParsePrimaryColName(allocator, colName, row)
	}

	// If nothing matches, return the text
	text := primExpr.GetText()
	fields := []serializers.FieldDef{{
		Name: "?column?",
		OID:  ptypes.PTypeText,
	}}
	resultSchema := serializers.NewBaseSchema(fields)
	serVal, err := serializers.TextSerializerInstance.Serialize(allocator, text)
	if err != nil {
		return nil, fmt.Errorf("[ParsePrimaryExpression] failed to serialize text value: %v", err)
	}
	resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
	if err != nil {
		return nil, fmt.Errorf("[ParsePrimaryExpression] failed to pack text value: %v", err)
	}
	return &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows:   []*ptypes.Row{resultRow},
			Schema: resultSchema,
		},
	}, nil
}
