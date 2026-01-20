package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/storage"
)

func ExecutePgTables(stmt *statements.SelectStatement, store *storage.DistributedStorageVClock) ([]map[string]interface{}, error) {
	// Return system tables
	allItems := []map[string]interface{}{
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_database",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_tables",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_columns",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_class",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_namespace",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_type",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_attribute",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_index",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
		{
			"schemaname":  "pg_catalog",
			"tablename":   "pg_constraint",
			"tableowner":  "postgres",
			"tablespace":  nil,
			"hasindexes":  false,
			"hasrules":    false,
			"hastriggers": false,
			"rowsecurity": false,
		},
	}

	return rhelpers.FilterColumns(allItems, stmt), nil
}
