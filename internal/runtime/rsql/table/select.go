package table

import (
	"encoding/json"
	"fmt"
	"log"
	"petacore/internal/core"
	"petacore/internal/storage"
	"sort"
)

// Select выполняет SELECT запрос
// TODO переписать, чтобы Select генерировал план запроса, а не просто что-то сканировал. Всё должно проходить в транзакции на этапе executor
func (t *Table) Select(tx *storage.DistributedTransactionVClock, tableName string, columns []string, where map[string]interface{}, limit int) ([]map[string]interface{}, []string, []ColType, error) {
	// fmt.Printf("DEBUG: Select from %s, columns: %v, where: %+v, limit: %d\n", tableName, columns, where, limit)
	var results []map[string]interface{}

	// Получаем метаданные таблицы
	tableKey := t.getMetadataPrefixKey()
	metaStr, found := tx.Read([]byte(tableKey))
	if !found || metaStr == "" {
		return nil, nil, nil, fmt.Errorf("table %s does not exist", tableName)
	}

	var meta TableMetadata
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		return nil, nil, nil, err
	}

	log.Printf("DEBUG: Table metadata: %+v\n", meta)

	// Определяем columns и columnTypes
	var finalColumns []string
	var columnTypes []ColType
	if columns == nil || (len(columns) == 1 && columns[0] == "*") {
		finalColumns = make([]string, 0, len(meta.Columns))
		for name := range meta.Columns {
			finalColumns = append(finalColumns, name)
		}
		sort.Strings(finalColumns)
		columnTypes = make([]ColType, len(finalColumns))
		for i, name := range finalColumns {
			columnTypes[i] = meta.Columns[name].Type
		}
	} else {
		finalColumns = columns
		columnTypes = make([]ColType, len(columns))
		for i, name := range columns {
			if colMeta, ok := meta.Columns[name]; ok {
				columnTypes[i] = colMeta.Type
			} else {
				return nil, nil, nil, fmt.Errorf("column %s does not exist", name)
			}
		}
	}

	// Сканируем все строки
	prefix := t.getRowPrefixKey()
	// log.Printf("DEBUG: Scanning rows with prefix: %s\n", prefix)
	kvMap, err := tx.Scan([]byte(prefix), core.IteratorTypeAll, limit)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, value := range kvMap {
		var rowData []map[string]interface{}
		if err := json.Unmarshal([]byte(value), &rowData); err != nil {
			// Try as single row
			var singleRow map[string]interface{}
			if err2 := json.Unmarshal([]byte(value), &singleRow); err2 != nil {
				// log.Printf("DEBUG: Failed to unmarshal row data: %v\n", err)
				continue
			}
			rowData = []map[string]interface{}{singleRow}
		}

		for _, row := range rowData {
			if len(row) == 0 {
				continue
			}
			// No WHERE filtering here, done in executor
			if columns == nil || (len(columns) == 1 && columns[0] == "*") {
				results = append(results, row)
			} else {
				filtered := make(map[string]interface{})
				for _, col := range columns {
					if val, ok := row[col]; ok {
						filtered[col] = val
					}
				}

				results = append(results, filtered)
			}
		}
	}

	// Применяем LIMIT
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, finalColumns, columnTypes, nil
}
