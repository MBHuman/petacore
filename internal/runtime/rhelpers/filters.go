package rhelpers

import (
	"log"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
)

// filterColumns filters the rows to only include the requested columns
func FilterColumns(rows []map[string]interface{}, stmt *statements.SelectStatement) []map[string]interface{} {
	if len(stmt.Columns) == 0 {
		return rows
	}

	// Special case for DataGrip query: SELECT N FROM pg_catalog.pg_database N
	if stmt.TableName == "pg_catalog.pg_database" && len(stmt.Columns) == 1 && stmt.Columns[0].ColumnName == "N" {
		return rows
	}

	// Check if selecting table alias (e.g., SELECT N FROM table N)
	if len(stmt.Columns) == 1 && stmt.Columns[0].ColumnName == stmt.TableAlias && stmt.TableAlias != "" {
		return rows
	}

	// Check if SELECT * is used
	if len(stmt.Columns) == 1 && stmt.Columns[0].ColumnName == "*" {
		return rows
	}

	// Get the list of requested columns
	var requestedColumns []string
	for _, col := range stmt.Columns {
		if col.ColumnName != "" && col.ColumnName != "*" {
			requestedColumns = append(requestedColumns, col.ColumnName)
		}
	}

	// If no specific columns requested, return all
	if len(requestedColumns) == 0 {
		return rows
	}

	// Filter each row to only include requested columns
	var filteredRows []map[string]interface{}
	for _, row := range rows {
		filteredRow := make(map[string]interface{})
		for _, col := range requestedColumns {
			if value, exists := row[col]; exists {
				filteredRow[col] = value
			}
		}
		filteredRows = append(filteredRows, filteredRow)
	}

	return filteredRows
}

// filterRowsByWhere filters rows based on the WHERE clause
func FilterRowsByWhere(rows []map[string]interface{}, where *items.WhereClause) []map[string]interface{} {
	if where == nil {
		return rows
	}

	log.Printf("DEBUG FilterRowsByWhere: filtering %d rows with WHERE: %+v", len(rows), where)
	var filteredRows []map[string]interface{}
	for _, row := range rows {
		matches := EvaluateWhereCondition(where, row)
		log.Printf("DEBUG FilterRowsByWhere: row %v matches: %v", row, matches)
		if matches {
			filteredRows = append(filteredRows, row)
		}
	}
	log.Printf("DEBUG FilterRowsByWhere: filtered to %d rows", len(filteredRows))
	return filteredRows
}
