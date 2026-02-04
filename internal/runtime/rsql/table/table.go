package table

import (
	"bytes"
	"encoding/json"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/storage"
	"petacore/internal/utils"
	"strconv"
)

type ExecuteResult struct {
	Rows    [][]interface{}
	Columns []TableColumn
}

type ResultRow struct {
	Row     []interface{}
	Columns []TableColumn
}

type SelectColumn struct {
	Name  string
	IsAll bool
}

type ITable interface {
	CreateTable(name string, columns []ColumnDef, ifNotExists bool) error
	Insert(tableName string, values map[string]interface{}) error
	Select(
		tableName string,
		columns []SelectColumn,
		where map[string]interface{},
		limit int,
	) (*ExecuteResult, error)
	DropTable(name string) error
}

type Table struct {
	Storage  *storage.DistributedStorageVClock
	Database string
	Schema   string
	Name     string
}

type TableMetadata struct {
	Name        string
	Columns     map[string]ColumnMetadata
	PrimaryKeys []int
}

type ColumnMetadata struct {
	Type ColType
	// IsPrimaryKey bool
	IsNullable   bool
	DefaultValue interface{}
	IsSerial     bool
	IsUnique     bool
	Idx          int
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
	case ColTypeBigInt:
		switch value.(type) {
		case int, int32, int64, float64:
			// Разрешаем числовые типы для bigint
		default:
			return fmt.Errorf("expected bigint, got %T", value)
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

	seqValueStr, found := tx.Read([]byte(sequencePrefixKey))
	seqValue := uint64(1)
	if found && seqValueStr != "" {
		if parsed, err := strconv.ParseUint(seqValueStr, 10, 64); err == nil {
			seqValue = parsed
		}
	}

	logger.Debugf("Using sequence prefix key %s with value %d", sequencePrefixKey, seqValue)
	tx.Write([]byte(sequencePrefixKey), strconv.FormatUint(seqValue+1, 10))
	return seqValue
}

func (t *Table) getRowKey(primaryKeys []interface{}) []byte {
	tablePrefixKey := t.getTablePrefixKey()
	var buf bytes.Buffer
	for _, pk := range primaryKeys {
		buf.WriteString(fmt.Sprintf("%v", pk))
	}
	rowID := buf.Bytes()
	rowKey := utils.GenTableRowKey(tablePrefixKey, string(rowID))
	return []byte(rowKey)
}

func (t *Table) getRowPrefixKey() string {
	tablePrefixKey := t.getTablePrefixKey()
	rowPrefixKey := utils.GenTableRowPrefix(tablePrefixKey)
	return rowPrefixKey
}

// TableExists проверяет, существует ли таблица
func (t *Table) TableExists(tx *storage.DistributedTransactionVClock) bool {
	var exists bool
	metaPrefixKey := t.getMetadataPrefixKey()
	_, found := tx.Read([]byte(metaPrefixKey))
	exists = found
	return exists
}

// GetTableMetadataInTx получает метаданные таблицы в транзакции
func (t *Table) GetTableMetadataInTx(tx *storage.DistributedTransactionVClock) (*TableMetadata, error) {
	metaPrefixKey := t.getMetadataPrefixKey()
	metaDataStr, found := tx.Read([]byte(metaPrefixKey))
	if !found || metaDataStr == "" {
		return nil, fmt.Errorf("table %s does not exist", t.Name)
	}

	var meta TableMetadata
	err := json.Unmarshal([]byte(metaDataStr), &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}
