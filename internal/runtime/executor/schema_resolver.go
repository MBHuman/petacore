package executor

import "strings"

// ComputeSchemaAndTableName определяет схему и имя таблицы на основе имени таблицы и контекста выполнения.
// Реализует логику search_path для системных таблиц pg_catalog.
//
// Правила резолва:
// 1. Если имя таблицы содержит явное указание схемы (schema.table), используется явная схема
// 2. Если имя таблицы - известная системная таблица pg_catalog, используется схема pg_catalog
// 3. Иначе используется схема из контекста выполнения (обычно public)
func ComputeSchemaAndTableName(tableName string, exCtx *ExecutorContext) (string, string) {
	parts := strings.SplitN(tableName, ".", 2)
	if len(parts) == 2 {
		// Явно указана схема
		return parts[0], parts[1]
	}

	// Проверяем, является ли это системной таблицей pg_catalog
	if IsPgCatalogTable(tableName) {
		return "pg_catalog", tableName
	}

	// Используем схему из контекста
	return exCtx.Schema, tableName
}

// IsPgCatalogTable проверяет, является ли таблица системной таблицей pg_catalog
func IsPgCatalogTable(tableName string) bool {
	switch tableName {
	case "pg_tables", "pg_columns", "pg_class",
		"pg_attribute", "pg_proc", "pg_type", "pg_namespace",
		"pg_database", "pg_tablespace", "pg_roles",
		"pg_stat_ssl", "pg_shdescription", "pg_am":
		return true
	default:
		return false
	}
}
