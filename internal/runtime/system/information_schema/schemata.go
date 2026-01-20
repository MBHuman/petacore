package informationschema

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
)

func ExecuteSchemata(stmt *statements.SelectStatement) ([]map[string]interface{}, error) {
	rows := []map[string]interface{}{
		{
			"catalog_name":                  "testdb",
			"schema_name":                   "pg_catalog",
			"schema_owner":                  "postgres",
			"default_character_set_catalog": nil,
			"default_character_set_schema":  nil,
			"default_character_set_name":    nil,
			"sql_path":                      nil,
		},
		{
			"catalog_name":                  "testdb",
			"schema_name":                   "public",
			"schema_owner":                  "postgres",
			"default_character_set_catalog": nil,
			"default_character_set_schema":  nil,
			"default_character_set_name":    nil,
			"sql_path":                      nil,
		},
	}
	return rhelpers.FilterColumns(rows, stmt), nil
}
