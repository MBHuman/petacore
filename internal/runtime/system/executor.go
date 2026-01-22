package system

import (
	"fmt"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	informationschema "petacore/internal/runtime/system/information_schema"
	pgcatalog "petacore/internal/runtime/system/pg_catalog"
	"petacore/internal/storage"
	"strings"
)

// TODO перевести логику выполнения системных таблиц на общий рантайм через планировщик
func ExecuteSystemTableSelect(stmt *statements.SelectStatement, store *storage.DistributedStorageVClock) ([]map[string]interface{}, error) {
	originalTableName := strings.ToLower(stmt.TableName)
	tableName := originalTableName
	// Strip pg_catalog prefix if present
	tableName = strings.TrimPrefix(tableName, "pg_catalog.")

	// Handle special cases
	if originalTableName == "pg_catalog" {
		return nil, fmt.Errorf("\"pg_catalog\" is not a table")
	}

	var rows []map[string]interface{}
	var err error

	// Handle information_schema tables
	if strings.HasPrefix(originalTableName, "information_schema.") {
		tableName = strings.TrimPrefix(originalTableName, "information_schema.")
		switch tableName {
		case "tables":
			rows, err = informationschema.ExecuteTables(stmt, store)
		case "columns":
			rows, err = informationschema.ExecuteColumns(stmt, store)
		case "schemata":
			rows, err = informationschema.ExecuteSchemata(stmt)
		default:
			return []map[string]interface{}{}, nil
		}
	} else {
		// Handle pg_catalog tables
		switch tableName {
		case "pg_tables":
			rows, err = pgcatalog.ExecutePgTables(stmt, store)
		case "pg_columns":
			rows, err = pgcatalog.ExecutePgColumns(stmt, store)
		case "pg_class":
			rows, err = pgcatalog.ExecutePgClass(stmt)
		case "pg_attribute":
			rows, err = pgcatalog.ExecutePgAttribute(stmt)
		case "pg_proc":
			rows, err = pgcatalog.ExecutePgProc(stmt)
		case "pg_type":
			rows, err = pgcatalog.ExecutePgTypeExpanded(stmt)
		case "pg_namespace":
			rows, err = pgcatalog.ExecutePgNamespace(stmt)
		case "pg_database":
			rows, err = pgcatalog.ExecutePgDatabase(stmt)
		case "pg_tablespace":
			rows, err = pgcatalog.ExecutePgTablespace(stmt)
		case "pg_roles":
			rows, err = pgcatalog.ExecutePgRoles(stmt)
		case "pg_stat_ssl":
			rows, err = pgcatalog.ExecutePgStatSsl(stmt)
		case "pg_shdescription":
			rows, err = pgcatalog.ExecutePgShdescription(stmt)
		case "pg_am":
			rows, err = pgcatalog.ExecutePgAmExpanded(stmt)
		default:
			return nil, fmt.Errorf("unknown system table: %s", tableName)
		}
	}

	if err != nil {
		return nil, err
	}

	// Apply WHERE clause filtering
	if stmt.Where != nil {
		rows = rhelpers.FilterRowsByWhere(rows, stmt.Where)
	}

	// Apply LIMIT
	if stmt.Limit > 0 && len(rows) > stmt.Limit {
		rows = rows[:stmt.Limit]
	}

	return rows, nil
}
