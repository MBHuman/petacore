package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
)

func ExecutePgTablespace(stmt *statements.SelectStatement) ([]map[string]interface{}, error) {
	allRows := []map[string]interface{}{
		{
			"oid":        1663,
			"spcname":    "pg_default",
			"spcowner":   10,
			"spcacl":     nil,
			"spcoptions": nil,
		},
		{
			"oid":        1664,
			"spcname":    "pg_global",
			"spcowner":   10,
			"spcacl":     nil,
			"spcoptions": nil,
		},
	}
	return rhelpers.FilterColumns(allRows, stmt), nil
}
