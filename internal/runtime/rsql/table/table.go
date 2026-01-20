package table

import (
	"fmt"
	"petacore/internal/storage"
	"petacore/internal/utils"
	"sort"
	"strconv"
)

type ITable interface {
	CreateTable(name string, columns []ColumnDef) error
	Insert(tableName string, values map[string]interface{}) error
	Select(tableName string, columns []string, where map[string]interface{}, limit int) ([]map[string]interface{}, error)
	DropTable(name string) error
}

type Table struct {
	Storage  *storage.DistributedStorageVClock
	Database string
	Schema   string
	Name     string
}

type TableMetadata struct {
	Name    string
	Columns map[string]ColumnMetadata
}

type ColumnMetadata struct {
	Type         ColType
	IsPrimaryKey bool
	IsNullable   bool
	DefaultValue interface{}
	IsSerial     bool
	IsUnique     bool
}

type Row map[string]interface{}

// NewTable создает новый экземпляр Table
func NewTable(
	storage *storage.DistributedStorageVClock,
	database, schema, name string,
) *Table {
	return &Table{
		Storage:  storage,
		Database: database,
		Schema:   schema,
		Name:     name,
	}
}

// validateValueType проверяет тип значения
func (t *Table) validateValueType(value interface{}, expectedType ColType) error {
	switch expectedType {
	case ColTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case ColTypeInt:
		switch value.(type) {
		case int, int32, int64, float64:
			// Разрешаем числовые типы для int
		default:
			return fmt.Errorf("expected int, got %T", value)
		}
	case ColTypeFloat:
		switch value.(type) {
		case float32, float64, int, int32, int64:
			// Разрешаем числовые типы для float
		default:
			return fmt.Errorf("expected float, got %T", value)
		}
	case ColTypeBool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected bool, got %T", value)
		}
	}
	return nil
}

// getNextRowID получает следующий ID для строки
// getColumnNames возвращает имена колонок в порядке определения
func getColumnNames(meta TableMetadata) []string {
	names := make([]string, 0, len(meta.Columns))
	for name := range meta.Columns {
		names = append(names, name)
	}
	// Сортируем для консистентности
	sort.Strings(names)
	return names
}

// matchesWhere проверяет, соответствует ли строка условию WHERE
func (t *Table) matchesWhere(row map[string]interface{}, where map[string]interface{}) bool {
	if len(where) == 0 {
		return true
	}

	for key, expectedValue := range where {
		actualValue, exists := row[key]
		if !exists {
			return false
		}

		// Простое сравнение (в реальности нужно поддерживать операторы)
		if actualValue != expectedValue {
			return false
		}
	}

	return true
}

func (t *Table) getTablePrefixKey() string {
	key := utils.GenTablePrefix(
		&utils.TableKey{
			Database: t.Database,
			Schema:   t.Schema,
			Table:    t.Name,
		},
	)
	return key
}

func (t *Table) getMetadataPrefixKey() string {
	tableKeyPrefix := t.getTablePrefixKey()
	metadataPrefixKey := utils.GenMetaPrefix(tableKeyPrefix)
	return metadataPrefixKey
}

func (t *Table) getSequencePrefixKey(colName string) string {
	sequenceKeyPrefix := utils.GenSequenceKey(
		&utils.SequenceKey{
			Database: t.Database,
			Schema:   t.Schema,
			Table:    t.Name,
			Column:   colName,
		},
	)
	sequenceKeyPrefix = utils.GenSequencePrefix(sequenceKeyPrefix)
	return sequenceKeyPrefix
}

func (t *Table) getAllSequencePrefixKey() string {
	sequencePrefixKey := utils.GenSequencePrefixKey(
		&utils.SequenceKey{
			Database: t.Database,
			Schema:   t.Schema,
			Table:    t.Name,
		},
	)
	sequencePrefixKey = utils.GenSequencePrefix(sequencePrefixKey)
	return sequencePrefixKey
}

func (t *Table) genSequenceKey(tx *storage.DistributedTransactionVClock, colName string) uint64 {
	sequencePrefixKey := t.getSequencePrefixKey(colName)
	sequencePrefixKey = utils.GenSequencePrefix(sequencePrefixKey)

	seqValueStr, found := tx.Read([]byte(sequencePrefixKey))
	seqValue := uint64(1)
	if found && seqValueStr != "" {
		if parsed, err := strconv.ParseUint(seqValueStr, 10, 64); err == nil {
			seqValue = parsed
		}
	}

	tx.Write([]byte(sequencePrefixKey), strconv.FormatUint(seqValue+1, 10))
	return seqValue
}

func (t *Table) getRowKey(rowID string) string {
	tablePrefixKey := t.getTablePrefixKey()
	rowKey := utils.GenTableRowKey(tablePrefixKey, rowID)
	return rowKey
}

func (t *Table) getRowPrefixKey() string {
	tablePrefixKey := t.getTablePrefixKey()
	rowPrefixKey := utils.GenTableRowPrefix(tablePrefixKey)
	return rowPrefixKey
}
