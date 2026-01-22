package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
)

func ExecutePgAmExpanded(stmt *statements.SelectStatement) ([]map[string]interface{}, error) {
	allRows := []map[string]interface{}{
		{
			"oid":       2,
			"amname":    "heap",
			"amhandler": 0,
			"amtype":    "t",
		},
		{
			"oid":       403,
			"amname":    "btree",
			"amhandler": 0,
			"amtype":    "i",
		},
		{
			"oid":       405,
			"amname":    "hash",
			"amhandler": 0,
			"amtype":    "i",
		},
		{
			"oid":       783,
			"amname":    "gist",
			"amhandler": 0,
			"amtype":    "i",
		},
		{
			"oid":       2742,
			"amname":    "gin",
			"amhandler": 0,
			"amtype":    "i",
		},
		{
			"oid":       4000,
			"amname":    "spgist",
			"amhandler": 0,
			"amtype":    "i",
		},
		{
			"oid":       3580,
			"amname":    "brin",
			"amhandler": 0,
			"amtype":    "i",
		},
	}

	return rhelpers.FilterColumns(allRows, stmt), nil
}
