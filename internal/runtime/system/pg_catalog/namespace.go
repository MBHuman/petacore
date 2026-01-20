package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
)

func ExecutePgNamespace(stmt *statements.SelectStatement) ([]map[string]interface{}, error) {
	allRows := []map[string]interface{}{
		{
			"oid":      11,
			"nspname":  "pg_catalog",
			"nspowner": 10,
			"nspacl":   nil,
		},
		{
			"oid":      99,
			"nspname":  "pg_toast",
			"nspowner": 10,
			"nspacl":   nil,
		},
		{
			"oid":      2200,
			"nspname":  "public",
			"nspowner": 10,
			"nspacl":   nil,
		},
	}
	return rhelpers.FilterColumns(allRows, stmt), nil
}
