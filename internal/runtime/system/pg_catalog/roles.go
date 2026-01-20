package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
)

func ExecutePgRoles(stmt *statements.SelectStatement) ([]map[string]interface{}, error) {
	allRows := []map[string]interface{}{
		{
			"oid":            10,
			"rolname":        "postgres",
			"rolsuper":       true,
			"rolinherit":     true,
			"rolcreaterole":  true,
			"rolcreatedb":    true,
			"rolcanlogin":    true,
			"rolreplication": true,
			"rolconnlimit":   -1,
			"rolpassword":    nil,
			"rolvaliduntil":  nil,
			"rolbypassrls":   true,
			"rolconfig":      nil,
		},
	}
	return rhelpers.FilterColumns(allRows, stmt), nil
}
