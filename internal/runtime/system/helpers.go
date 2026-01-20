package system

import (
	"strings"
)

func IsSystemTable(tableName string) bool {
	switch strings.ToLower(tableName) {
	case "pg_database", "pg_tables", "pg_columns", "pg_class", "pg_namespace", "pg_type", "pg_attribute", "pg_index", "pg_constraint", "pg_roles", "pg_catalog", "pg_stat_ssl":
		return true
	default:
		return strings.HasPrefix(strings.ToLower(tableName), "information_schema.") || strings.HasPrefix(strings.ToLower(tableName), "pg_catalog.")
	}
}
