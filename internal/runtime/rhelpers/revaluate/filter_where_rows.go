package revaluate

import (
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/table"
)

// filterRowsByWhere filters rows based on the WHERE clause
func EvaluateFilterRowsByWhere(execResult *table.ExecuteResult, where *items.WhereClause) *table.ExecuteResult {
	if where == nil {
		return execResult
	}

	var filteredRows [][]interface{}
	for _, row := range execResult.Rows {
		matches := EvaluateWhereCondition(where, &table.ResultRow{
			Row:     row,
			Columns: execResult.Columns,
		})
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
