package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
)

func ExecutePgShdescription(stmt *statements.SelectStatement) ([]map[string]interface{}, error) {
	allItems := []map[string]interface{}{
		{
			"objoid":      1,
			"classoid":    1262, // pg_database
			"objsubid":    0,
			"description": "default database",
		},
	}

	return rhelpers.FilterColumns(allItems, stmt), nil
}
