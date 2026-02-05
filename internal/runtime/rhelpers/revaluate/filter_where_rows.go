package revaluate

import (
	"context"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
)

// filterRowsByWhere filters rows based on the WHERE clause
func EvaluateFilterRowsByWhere(goCtx context.Context, execResult *table.ExecuteResult, where *items.WhereClause, statement *statements.SelectStatement, subExec subquery.SubqueryExecutor, runtimeParams map[int]interface{}) *table.ExecuteResult {
	if where == nil {
		return execResult
	}

	var filteredRows [][]interface{}
	for _, row := range execResult.Rows {
		matches := EvaluateWhereCondition(goCtx, where, &table.ResultRow{
			Row:     row,
			Columns: execResult.Columns,
		}, statement, subExec, runtimeParams)
		if matches {
			filteredRows = append(filteredRows, row)
		}
	}

	filteredResult := &table.ExecuteResult{
		Columns: execResult.Columns,
		Rows:    filteredRows,
	}
	return filteredResult
}
