package table

import (
	"encoding/json"
	"fmt"
	"log"
	"petacore/internal/core"
	"petacore/internal/storage"
	"time"
)

// Insert вставляет строку в таблицу
func (t *Table) Insert(tableName string, values []map[string]interface{}) error {
	log.Printf("DEBUG: Insert into %s: %+v\n", tableName, values)
	return t.Storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		// Получаем метаданные таблицы
		metaPrefixKey := t.getMetadataPrefixKey()
		metaStr, found := tx.Read([]byte(metaPrefixKey))
		if !found || metaStr == "" {
			return fmt.Errorf("table %s does not exist", tableName)
		}

		var meta TableMetadata
		if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
			return err
		}

		log.Printf("DEBUG: Table metadata: %+v\n", meta)

		for _, value := range values {
			// Применяем defaults и генерируем SERIAL значения
			for colName, colMeta := range meta.Columns {
				if _, exists := value[colName]; !exists {
					if colMeta.IsSerial {
						// // Генерируем следующий ID из sequence
						seqValue := t.genSequenceKey(tx, colName)
						value[colName] = seqValue
					} else if colMeta.DefaultValue != nil {
						if colMeta.DefaultValue == "CURRENT_TIMESTAMP" {
							value[colName] = time.Now().Format("2006-01-02 15:04:05")
						} else {
							value[colName] = colMeta.DefaultValue
						}
					}
				}
			}

			// Проверяем NOT NULL constraints
			for colName, colMeta := range meta.Columns {
				if !colMeta.IsNullable {
					if value, exists := value[colName]; !exists || value == nil {
						return fmt.Errorf("null value in column %s violates not-null constraint", colName)
					}
				}
			}

			// Проверяем UNIQUE constraints
			// TODO поменять на фильтр блума в метаданных, для быстрой проверки
			for colName, colMeta := range meta.Columns {
				if colMeta.IsUnique || colMeta.IsPrimaryKey {
					if value, exists := value[colName]; exists {
						prefix := t.getRowPrefixKey()
						kvMap, err := tx.Scan([]byte(prefix), core.IteratorTypeAll, -1)
						if err != nil {
							return err
						}
						for _, val := range kvMap {
							var rowData map[string]interface{}
							if err := json.Unmarshal([]byte(val), &rowData); err != nil {
								continue
							}
							if rowData[colName] == value {
								return fmt.Errorf("duplicate key value violates unique constraint \"%s\"", colName)
							}
						}
					}
				}
			}

			// Генерируем уникальный ID для строки
			var rowID string
			// Если есть PRIMARY KEY, используем его значение как rowID
			for colName, colMeta := range meta.Columns {
				if colMeta.IsPrimaryKey {
					if val, exists := value[colName]; exists {
						rowID = fmt.Sprintf("%v", val)
					}
					break
				}
			}
			// Обязательно должен быть rowID
			if rowID == "" {
				return fmt.Errorf("cannot determine row ID for table %s", tableName)
			}
			// Сохраняем строку
			rowKey := t.getRowKey(rowID)
			rowData, err := json.Marshal(value)
			if err != nil {
				return err
			}
			tx.Write([]byte(rowKey), string(rowData))
			log.Printf("DEBUG: Saved row %s: %s\n", rowKey, string(rowData))
		}

		return nil
	})
}
