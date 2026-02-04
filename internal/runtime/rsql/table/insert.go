package table

import (
	"encoding/json"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/storage"
	"time"

	"go.uber.org/zap"
)

// Insert вставляет строку в таблицу
// TODO убрать хардкодинг, сделать поддержку всех типов данных и ограничений
func (t *Table) Insert(tableName string, values [][]interface{}, columnNames []string) error {
	logger.Debugf("DEBUG: Insert into %s: %+v\n", tableName, values)
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

		columnsMap := make(map[string]int)
		for idx, colName := range columnNames {
			columnsMap[colName] = idx
		}

		logger.Debug("DEBUG: Table metadata:", zap.Any("meta", meta))

		for _, value := range values {
			insertRow := make([]interface{}, len(meta.Columns))
			// Применяем defaults и генерируем SERIAL значения
			for colName, colMeta := range meta.Columns {
				if idx, exists := columnsMap[colName]; exists {
					insertRow[colMeta.Idx-1] = value[idx]
				} else {
					if colMeta.IsSerial {
						// // Генерируем следующий ID из sequence
						seqValue := t.genSequenceKey(tx, colName)
						insertRow[colMeta.Idx-1] = seqValue
					} else if colMeta.DefaultValue != nil {
						if colMeta.DefaultValue == "CURRENT_TIMESTAMP" {
							insertRow[colMeta.Idx-1] = time.Now().Format("2006-01-02 15:04:05")
						} else {
							insertRow[colMeta.Idx-1] = colMeta.DefaultValue
						}
					}
				}
			}

			// Проверяем NOT NULL constraints
			for colName, colMeta := range meta.Columns {
				if !colMeta.IsNullable {
					if val := insertRow[colMeta.Idx-1]; val == nil {
						return fmt.Errorf("null value in column %s violates not-null constraint", colName)
					}
				}
			}

			// Проверяем UNIQUE constraints
			// TODO сделать нормальный UNIQUE для быстрой проверки на стороне индекса,
			// возможно через фильтр блума

			// for colName, colMeta := range meta.Columns {
			// 	if colMeta.IsUnique || colMeta.IsPrimaryKey {
			// 		if value, exists := value[colName]; exists {
			// 			prefix := t.getRowPrefixKey()
			// 			kvMap, err := tx.Scan([]byte(prefix), core.IteratorTypeAll, -1)
			// 			if err != nil {
			// 				return err
			// 			}
			// 			for _, val := range kvMap {
			// 				var rowData map[string]interface{}
			// 				if err := json.Unmarshal([]byte(val), &rowData); err != nil {
			// 					continue
			// 				}
			// 				if rowData[colName] == value {
			// 					return fmt.Errorf("duplicate key value violates unique constraint \"%s\"", colName)
			// 				}
			// 			}
			// 		}
			// 	}
			// }

			// Генерируем уникальный ID для строки

			primaryKeys := make([]interface{}, 0, len(meta.PrimaryKeys))
			for _, pkIdx := range meta.PrimaryKeys {
				if insertRow[pkIdx-1] == nil {
					return fmt.Errorf("primary key column index %d cannot be null", pkIdx)
				}
				primaryKeys = append(primaryKeys, insertRow[pkIdx-1])
			}

			// Обязательно должен быть rowID
			if len(primaryKeys) == 0 {
				return fmt.Errorf("cannot determine row ID for table %s", tableName)
			}
			// Сохраняем строку
			rowKey := t.getRowKey(primaryKeys)

			if _, ok := tx.Read(rowKey); ok {
				return fmt.Errorf("duplicate key value violates primary key constraint")
			}

			rowData, err := json.Marshal(insertRow)
			if err != nil {
				return err
			}
			tx.Write(rowKey, string(rowData))
			logger.Debug("DEBUG: Saved row: ",
				zap.String("rowKey", string(rowKey)),
				zap.String("rowData", string(rowData)),
			)
		}

		return nil
	})
}
