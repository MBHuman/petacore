package informationschema

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/storage"
)

func ExecuteTables(stmt *statements.SelectStatement, store *storage.DistributedStorageVClock) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}

	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		// For now, return hardcoded system tables since we don't have a way to enumerate all user tables
		// In a full implementation, we would scan all keys with prefix "schema:"
		systemTables := []string{
			"pg_database", "pg_tables", "pg_columns", "pg_class", "pg_namespace",
			"pg_type", "pg_attribute", "pg_proc", "pg_roles", "pg_tablespace",
		}

		for _, tableName := range systemTables {
			rows = append(rows, map[string]interface{}{
				"table_schema": "pg_catalog",
				"table_name":   tableName,
				"table_type":   "BASE TABLE",
			})
		}

		return nil
	})

	return rhelpers.FilterColumns(rows, stmt), err
}
