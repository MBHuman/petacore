package revaluate

import (
	"context"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

// filterRowsByWhere filters rows based on the WHERE clause
func EvaluateFilterRowsByWhere(
	allocator pmem.Allocator,
	goCtx context.Context,
	execResult *table.ExecuteResult,
	where *items.WhereClause,
	statement *statements.SelectStatement,
	subExec subquery.SubqueryExecutor,
	runtimeParams map[int]interface{},
) *table.ExecuteResult {
	if where == nil {
		return execResult
	}

	var filteredRows []*ptypes.Row
	for _, row := range execResult.Rows {
		matches := EvaluateWhereCondition(
			allocator,
			goCtx, where, &table.ResultRow{
				Row:    row,
				Schema: execResult.Schema,
			}, statement, subExec, runtimeParams,
		)
		if matches {
			filteredRows = append(filteredRows, row)
		}
	}

	filteredResult := &table.ExecuteResult{
		Schema: execResult.Schema,
		Rows:   filteredRows,
	}
	return filteredResult
}
