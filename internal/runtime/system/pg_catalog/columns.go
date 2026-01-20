package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/storage"
)

func ExecutePgColumns(stmt *statements.SelectStatement, store *storage.DistributedStorageVClock) ([]map[string]interface{}, error) {
	// Return columns for system tables
	allRows := []map[string]interface{}{
		{
			"table_schema":             "pg_catalog",
			"table_name":               "pg_database",
			"column_name":              "datname",
			"ordinal_position":         1,
			"column_default":           nil,
			"is_nullable":              "NO",
			"data_type":                "name",
			"character_maximum_length": nil,
			"numeric_precision":        nil,
		},
		{
			"table_schema":             "pg_catalog",
			"table_name":               "pg_database",
			"column_name":              "datdba",
			"ordinal_position":         2,
			"column_default":           nil,
			"is_nullable":              "NO",
			"data_type":                "oid",
			"character_maximum_length": nil,
			"numeric_precision":        nil,
		},
		{
			"table_schema":             "pg_catalog",
			"table_name":               "pg_database",
			"column_name":              "encoding",
			"ordinal_position":         3,
			"column_default":           nil,
			"is_nullable":              "NO",
			"data_type":                "integer",
			"character_maximum_length": nil,
			"numeric_precision":        32,
		},
		{
			"table_schema":             "pg_catalog",
			"table_name":               "pg_database",
			"column_name":              "datcollate",
			"ordinal_position":         4,
			"column_default":           nil,
			"is_nullable":              "NO",
			"data_type":                "name",
			"character_maximum_length": nil,
			"numeric_precision":        nil,
		},
		{
			"table_schema":             "pg_catalog",
			"table_name":               "pg_database",
			"column_name":              "datctype",
			"ordinal_position":         5,
			"column_default":           nil,
			"is_nullable":              "NO",
			"data_type":                "name",
			"character_maximum_length": nil,
			"numeric_precision":        nil,
		},
	}
	return rhelpers.FilterColumns(allRows, stmt), nil
}
