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

func ParseAdditiveExpression(
	allocator pmem.Allocator,
	ctx context.Context,
	addExpr parser.IAdditiveExpressionContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (rmodels.Expression, error) {
	if addExpr == nil {
		return nil, nil
	}

	multExprs := addExpr.AllMultiplicativeExpression()
	if len(multExprs) == 0 {
		return nil, nil
	}

	result, err := ParseMultiplicativeExpression(allocator, ctx, multExprs[0], row, subExec)
	if err != nil {
		return nil, err
	}

	if len(multExprs) == 1 {
		return result, nil
	}

	// собираем операторы в порядке появления
	type opInfo struct {
		op    string
		index int
	}
	var ops []opInfo
	for _, p := range addExpr.AllPLUS() {
		ops = append(ops, opInfo{"+", p.GetSymbol().GetTokenIndex()})
	}
	for _, m := range addExpr.AllMINUS() {
		ops = append(ops, opInfo{"-", m.GetSymbol().GetTokenIndex()})
	}
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].index < ops[j].index
	})

	if len(ops) != len(multExprs)-1 {
		return nil, fmt.Errorf("additive: operator count mismatch: %d ops for %d expressions",
			len(ops), len(multExprs))
	}

	for i, op := range ops {
		nextResult, err := ParseMultiplicativeExpression(allocator, ctx, multExprs[i+1], row, subExec)
		if err != nil {
			return nil, err
		}

		a, aOk := result.(*rmodels.ResultRowsExpression)
		b, bOk := nextResult.(*rmodels.ResultRowsExpression)
		if !aOk || !bOk {
			return nil, fmt.Errorf("additive: only result row expressions supported")
		}
		if len(a.Row.Rows) != 1 || len(b.Row.Rows) != 1 {
			return nil, fmt.Errorf("additive: only single-row expressions supported")
		}

		switch op.op {
		case "+":
			result, err = rops.AddValues(allocator, a.Row.Rows[0], b.Row.Rows[0], a.Row.Schema, b.Row.Schema)
		case "-":
			result, err = rops.SubtractValues(allocator, a.Row.Rows[0], b.Row.Rows[0], a.Row.Schema, b.Row.Schema)
		default:
			return nil, fmt.Errorf("additive: unknown operator %q", op.op)
		}
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
