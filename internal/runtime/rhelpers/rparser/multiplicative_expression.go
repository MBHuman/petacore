package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rops"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"sort"
)

func ParseMultiplicativeExpression(
	allocator pmem.Allocator,
	ctx context.Context,
	multExpr parser.IMultiplicativeExpressionContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (rmodels.Expression, error) {
	if multExpr == nil {
		return nil, nil
	}

	unaryExprs := multExpr.AllUnaryExpression()
	if len(unaryExprs) == 0 {
		return nil, nil
	}

	result, err := ParseUnaryExpression(allocator, ctx, unaryExprs[0], row, subExec)
	if err != nil {
		return nil, err
	}

	if len(unaryExprs) == 1 {
		return result, nil
	}

	type opInfo struct {
		op    string
		index int
	}
	var ops []opInfo
	for _, s := range multExpr.AllSTAR() {
		ops = append(ops, opInfo{"*", s.GetSymbol().GetTokenIndex()})
	}
	for _, s := range multExpr.AllSLASH() {
		ops = append(ops, opInfo{"/", s.GetSymbol().GetTokenIndex()})
	}
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].index < ops[j].index
	})

	if len(ops) != len(unaryExprs)-1 {
		return nil, fmt.Errorf("multiplicative: operator count mismatch: %d ops for %d expressions",
			len(ops), len(unaryExprs))
	}

	for i, op := range ops {
		nextResult, err := ParseUnaryExpression(allocator, ctx, unaryExprs[i+1], row, subExec)
		if err != nil {
			return nil, err
		}

		a, aOk := result.(*rmodels.ResultRowsExpression)
		b, bOk := nextResult.(*rmodels.ResultRowsExpression)
		if !aOk || !bOk {
			return nil, fmt.Errorf("multiplicative: only result row expressions supported")
		}
		if len(a.Row.Rows) != 1 || len(b.Row.Rows) != 1 {
			return nil, fmt.Errorf("multiplicative: only single-row expressions supported")
		}

		switch op.op {
		case "*":
			result, err = rops.MultiplyValues(allocator, a.Row.Rows[0], b.Row.Rows[0], a.Row.Schema, b.Row.Schema)
		case "/":
			result, err = rops.DivideValues(allocator, a.Row.Rows[0], b.Row.Rows[0], a.Row.Schema, b.Row.Schema)
		default:
			return nil, fmt.Errorf("multiplicative: unknown operator %q", op.op)
		}
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
