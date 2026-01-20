package table

import (
	"encoding/json"
	"fmt"
	"petacore/internal/storage"
	"sort"
)

// Describe возвращает метаданные таблицы в формате для DESCRIBE
func (t *Table) Describe(tx *storage.DistributedTransactionVClock) ([]map[string]interface{}, error) {
	// Получаем метаданные таблицы
	tableKey := t.getMetadataPrefixKey()
	metaStr, found := tx.Read([]byte(tableKey))
	if !found || metaStr == "" {
		return nil, fmt.Errorf("table %s does not exist", t.Name)
	}

	var meta TableMetadata
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal table metadata: %v", err)
	}

	// Форматируем результат как []map[string]interface{}
	var rows []map[string]interface{}
	for colName, colMeta := range meta.Columns {
		row := map[string]interface{}{
			"Field":   colName,
			"Type":    colMeta.Type.String(),
			"Null":    colMeta.IsNullable,
			"Key":     "",
			"Default": colMeta.DefaultValue,
			"Extra":   "",
		}
		if colMeta.IsPrimaryKey {
			row["Key"] = "PRI"
		}
		if colMeta.IsUnique {
			if row["Key"] == "" {
				row["Key"] = "UNI"
			} else {
				row["Key"] = row["Key"].(string) + ",UNI"
			}
		}
		if colMeta.IsSerial {
			row["Extra"] = "auto_increment"
		}
		rows = append(rows, row)
	}

	// Сортируем по имени колонки
	sort.Slice(rows, func(i, j int) bool {
		return rows[i]["Field"].(string) < rows[j]["Field"].(string)
	})

	return rows, nil
}
