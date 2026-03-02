package table

import (
	"encoding/json"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/storage"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"sort"
	"time"

	"go.uber.org/zap"
)

// Insert вставляет строку в таблицу
// TODO убрать хардкодинг, сделать поддержку всех типов данных и ограничений
func (t *Table) Insert(allocator pmem.Allocator, tableName string, values [][]ptypes.BaseType[any], columnNames []string) error {
	logger.Debugf("DEBUG: Insert into %s: %+v\n", tableName, values)
	return t.Storage.RunTransactionWithAllocator(allocator, func(tx *storage.DistributedTransactionVClock) error {
		// Получаем метаданные таблицы
		metaPrefixKey := t.getMetadataPrefixKey()
		metaStr, found := tx.Read([]byte(metaPrefixKey))
		if !found || len(metaStr) == 0 {
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

		// Создаем Schema на основе метаданных таблицы
		// Сортируем колонки по индексу
		type indexedColumn struct {
			name string
			meta ColumnMetadata
		}
		sortedCols := make([]indexedColumn, 0, len(meta.Columns))
		for colName, colMeta := range meta.Columns {
			sortedCols = append(sortedCols, indexedColumn{name: colName, meta: colMeta})
		}
		sort.Slice(sortedCols, func(i, j int) bool {
			return sortedCols[i].meta.Idx < sortedCols[j].meta.Idx
		})

		// Создаем Schema
		schemaFields := make([]serializers.FieldDef, 0, len(sortedCols))
		for _, col := range sortedCols {
			schemaFields = append(schemaFields, serializers.FieldDef{
				Name: col.name,
				OID:  col.meta.Type,
			})
		}
		schema := serializers.NewBaseSchema(schemaFields)

		for _, value := range values {
			// Подготавливаем данные для вставки в правильном порядке (по Idx)
			insertRowBuffers := make([][]byte, len(meta.Columns))

			// Применяем defaults и генерируем SERIAL значения
			for colName, colMeta := range meta.Columns {
				var fieldBuffer []byte

				if idx, exists := columnsMap[colName]; exists {
					// Значение предоставлено пользователем
					if value[idx] != nil {
						fieldBuffer = value[idx].GetBuffer()
					}
				} else {
					// Применяем default значения
					if colMeta.IsSerial {
						// Генерируем следующий ID из sequence
						seqValue := t.genSequenceKey(tx, colName)
						// Сериализуем sequence value (int32)
						// serialized, err := serializers.SerializeInt4(allocator, seqValue)
						serialized, err := serializers.Int8SerializerInstance.Serialize(allocator, int64(seqValue))
						if err != nil {
							return fmt.Errorf("failed to serialize serial value: %w", err)
						}
						fieldBuffer = serialized
					} else if colMeta.DefaultValue != nil {
						if colMeta.DefaultValue == "CURRENT_TIMESTAMP" {
							// Store current timestamp as int64 microseconds
							ts := time.Now()
							serialized, err := serializers.TimestampSerializerInstance.Serialize(allocator, &ts)
							if err != nil {
								return fmt.Errorf("failed to serialize current timestamp: %w", err)
							}
							fieldBuffer = serialized
						} else {
							// TODO: handle other default values
							return fmt.Errorf("unsupported default value type")
						}
					}
				}

				insertRowBuffers[colMeta.Idx-1] = fieldBuffer
			}

			// Проверяем NOT NULL constraints
			for colName, colMeta := range meta.Columns {
				if !colMeta.IsNullable {
					if insertRowBuffers[colMeta.Idx-1] == nil {
						return fmt.Errorf("null value in column %s violates not-null constraint", colName)
					}
				}
			}

			// Генерируем primary key для rowKey
			primaryKeyBuffers := make([][]byte, 0, len(meta.PrimaryKeys))
			for _, pkIdx := range meta.PrimaryKeys {
				if insertRowBuffers[pkIdx-1] == nil {
					return fmt.Errorf("primary key column index %d cannot be null", pkIdx)
				}
				primaryKeyBuffers = append(primaryKeyBuffers, insertRowBuffers[pkIdx-1])
			}

			// Обязательно должен быть rowID
			if len(primaryKeyBuffers) == 0 {
				return fmt.Errorf("cannot determine row ID for table %s", tableName)
			}

			// Создаем ключи из буферов для getRowKey
			// Нужно десериализовать primaryKeyBuffers в interface{} для совместимости
			primaryKeys := make([]interface{}, 0, len(primaryKeyBuffers))
			for i, pkIdx := range meta.PrimaryKeys {
				colMeta := sortedCols[pkIdx-1].meta
				buf := primaryKeyBuffers[i]
				deserialized, err := serializers.DeserializeGeneric(buf, colMeta.Type)
				if err != nil {
					return fmt.Errorf("failed to deserialize primary key: %w", err)
				}
				primaryKeys = append(primaryKeys, deserialized.IntoGo())
			}

			// Сохраняем строку
			rowKey := t.getRowKey(primaryKeys)

			if _, ok := tx.Read(rowKey); ok {
				return fmt.Errorf("duplicate key value violates primary key constraint")
			}

			// Упаковываем row в бинарный формат
			packedRow, err := schema.Pack(allocator, insertRowBuffers)
			if err != nil {
				return fmt.Errorf("failed to pack row: %w", err)
			}

			tx.Write(rowKey, packedRow.BufferPtr)
			logger.Debug("DEBUG: Saved row (binary): ",
				zap.String("rowKey", string(rowKey)),
				zap.Int("rowSize", len(packedRow.BufferPtr)),
			)
		}

		return nil
	})
}
