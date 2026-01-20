package table

import (
	"encoding/json"
	"fmt"
	"petacore/internal/storage"
)

// CreateTable создает новую таблицу
func (t *Table) CreateTable(name string, columns []ColumnDef) error {
	return t.Storage.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		// Проверяем, существует ли таблица
		metaPrefixKey := t.getMetadataPrefixKey()
		existing, found := tx.Read([]byte(metaPrefixKey))
		if found && existing != "" {
			return fmt.Errorf("table %s already exists", name)
		}

		// Создаем метаданные таблицы
		meta := TableMetadata{
			Name:    name,
			Columns: make(map[string]ColumnMetadata),
		}

		for _, col := range columns {
			meta.Columns[col.Name] = ColumnMetadata{
				Type:         col.Type,
				IsPrimaryKey: col.IsPrimaryKey,
				IsNullable:   col.IsNullable,
				DefaultValue: col.DefaultValue,
				IsSerial:     col.IsSerial,
			}
		}

		// Сохраняем метаданные
		metaData, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		tx.Write([]byte(metaPrefixKey), string(metaData))

		// Инициализируем sequences для SERIAL колонок
		for _, col := range columns {
			if col.IsSerial {
				seqKey := t.getSequencePrefixKey(col.Name)
				tx.Write([]byte(seqKey), "1")
			}
		}

		return nil
	})
}
