package revaluate

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/table"
	"sort"
)

// sortRows sorts the rows based on OrderBy items
// TODO сделать нормальную поддержку сортировки с учётом типов данных колонок
// не просто через interface{}
func EvaluateSortRows(execResult *table.ExecuteResult, orderBy []items.OrderByItem) {
	if len(orderBy) == 0 {
		return
	}

	columnIndexMap := make(map[string]int)
	for idx, col := range execResult.Columns {
		columnIndexMap[col.Name] = idx
	}

	sort.Slice(execResult.Rows, func(i, j int) bool {
		for _, ob := range orderBy {

			// Evaluate the expression for both rows
			// For simplicity, assume it's a column name
			exprText := ob.ColumnName
			var valI, valJ interface{}
			if ob.ColumnIndex > 0 {
				// SQL uses 1-based indexing, arrays use 0-based
				idx := ob.ColumnIndex - 1
				if idx >= len(execResult.Rows[i]) || idx >= len(execResult.Rows[j]) {
					continue // Skip invalid index
				}
				valI = execResult.Rows[i][idx]
				valJ = execResult.Rows[j][idx]
			} else if ob.ColumnName == "" {
				continue // Skip if neither index nor name specified
			} else {
				if colIdx, ok := columnIndexMap[exprText]; ok {
					valI = execResult.Rows[i][colIdx]
					valJ = execResult.Rows[j][colIdx]
				} else {
					continue
				}
			}

			// Compare values
			cmp := rhelpers.СompareValues(valI, valJ)
			if cmp != 0 {
				if ob.Direction == "DESC" {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return false
	})
}
