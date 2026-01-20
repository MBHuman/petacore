package informationschema

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/storage"
)

func ExecuteColumns(stmt *statements.SelectStatement, store *storage.DistributedStorageVClock) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}

	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		// For now, return hardcoded columns for system tables
		// In a full implementation, we would scan all schema keys and extract column info
		systemTableColumns := map[string][]map[string]interface{}{
			"pg_database": {
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "oid", "ordinal_position": 1, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datname", "ordinal_position": 2, "column_default": nil, "is_nullable": "NO", "data_type": "name"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datdba", "ordinal_position": 3, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "encoding", "ordinal_position": 4, "column_default": nil, "is_nullable": "NO", "data_type": "int4"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datcollate", "ordinal_position": 5, "column_default": nil, "is_nullable": "NO", "data_type": "name"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datctype", "ordinal_position": 6, "column_default": nil, "is_nullable": "NO", "data_type": "name"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datistemplate", "ordinal_position": 7, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datallowconn", "ordinal_position": 8, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datconnlimit", "ordinal_position": 9, "column_default": nil, "is_nullable": "NO", "data_type": "int4"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datlastsysoid", "ordinal_position": 10, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datfrozenxid", "ordinal_position": 11, "column_default": nil, "is_nullable": "NO", "data_type": "xid"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datminmxid", "ordinal_position": 12, "column_default": nil, "is_nullable": "NO", "data_type": "xid"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "dattablespace", "ordinal_position": 13, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_database", "column_name": "datacl", "ordinal_position": 14, "column_default": nil, "is_nullable": "YES", "data_type": "aclitem[]"},
			},
			"pg_namespace": {
				{"table_schema": "pg_catalog", "table_name": "pg_namespace", "column_name": "oid", "ordinal_position": 1, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_namespace", "column_name": "nspname", "ordinal_position": 2, "column_default": nil, "is_nullable": "NO", "data_type": "name"},
				{"table_schema": "pg_catalog", "table_name": "pg_namespace", "column_name": "nspowner", "ordinal_position": 3, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_namespace", "column_name": "nspacl", "ordinal_position": 4, "column_default": nil, "is_nullable": "YES", "data_type": "aclitem[]"},
			},
			"pg_class": {
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "oid", "ordinal_position": 1, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relname", "ordinal_position": 2, "column_default": nil, "is_nullable": "NO", "data_type": "name"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relnamespace", "ordinal_position": 3, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "reltype", "ordinal_position": 4, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "reloftype", "ordinal_position": 5, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relowner", "ordinal_position": 6, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relam", "ordinal_position": 7, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relfilenode", "ordinal_position": 8, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "reltablespace", "ordinal_position": 9, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relpages", "ordinal_position": 10, "column_default": nil, "is_nullable": "NO", "data_type": "int4"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "reltuples", "ordinal_position": 11, "column_default": nil, "is_nullable": "NO", "data_type": "float4"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relallvisible", "ordinal_position": 12, "column_default": nil, "is_nullable": "NO", "data_type": "int4"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "reltoastrelid", "ordinal_position": 13, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "reltoastidxid", "ordinal_position": 14, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relhasindex", "ordinal_position": 15, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relisshared", "ordinal_position": 16, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relpersistence", "ordinal_position": 17, "column_default": nil, "is_nullable": "NO", "data_type": "char"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relkind", "ordinal_position": 18, "column_default": nil, "is_nullable": "NO", "data_type": "char"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relnatts", "ordinal_position": 19, "column_default": nil, "is_nullable": "NO", "data_type": "int2"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relchecks", "ordinal_position": 20, "column_default": nil, "is_nullable": "NO", "data_type": "int2"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relhasrules", "ordinal_position": 21, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relhastriggers", "ordinal_position": 22, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relhassubclass", "ordinal_position": 23, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relrowsecurity", "ordinal_position": 24, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relforcerowsecurity", "ordinal_position": 25, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relispopulated", "ordinal_position": 26, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relreplident", "ordinal_position": 27, "column_default": nil, "is_nullable": "NO", "data_type": "char"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relispartition", "ordinal_position": 28, "column_default": nil, "is_nullable": "NO", "data_type": "bool"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relrewrite", "ordinal_position": 29, "column_default": nil, "is_nullable": "NO", "data_type": "oid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relfrozenxid", "ordinal_position": 30, "column_default": nil, "is_nullable": "NO", "data_type": "xid"},
				{"table_schema": "pg_catalog", "table_name": "pg_class", "column_name": "relminmxid", "ordinal_position": 31, "column_default": nil, "is_nullable": "NO", "data_type": "xid"},
			},
		}

		for _, columns := range systemTableColumns {
			rows = append(rows, columns...)
		}
		return nil
	})

	return rhelpers.FilterColumns(rows, stmt), err
}
