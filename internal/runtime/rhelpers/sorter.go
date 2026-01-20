package rhelpers

import (
	"petacore/internal/runtime/rsql/items"
	"sort"
	"strings"
)

// sortRows sorts the rows based on OrderBy items
func SortRows(rows []map[string]interface{}, orderBy []items.OrderByItem) {
	if len(orderBy) == 0 {
		return
	}

	sort.Slice(rows, func(i, j int) bool {
		for _, ob := range orderBy {
			// Evaluate the expression for both rows
			// For simplicity, assume it's a column name
			exprText := ob.ExpressionContext.GetText()
			var valI, valJ interface{}
			if strings.Contains(exprText, ".") {
				// Column name like N.oid
				parts := strings.Split(exprText, ".")
				if len(parts) == 2 {
					colName := parts[1]
					// For now, ignore table alias
					valI = rows[i][colName]
					valJ = rows[j][colName]
				}
			} else {
				valI = rows[i][exprText]
				valJ = rows[j][exprText]
			}

			// Compare values
			cmp := compareValues(valI, valJ)
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
