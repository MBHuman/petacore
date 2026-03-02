package table

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/storage"
	"petacore/internal/utils"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"sort"
	"strconv"
)

type ExecuteResult struct {
	Rows   []*ptypes.Row
	Schema *serializers.BaseSchema
}

type ResultRow struct {
	Row     *ptypes.Row
	Schema  *serializers.BaseSchema
	Context context.Context // For passing queryStartTime and other runtime context
}

type SelectColumn struct {
	Name  string
	IsAll bool
}

type ITable interface {
	CreateTable(name string, columns []ColumnDef, ifNotExists bool) error
	Insert(tableName string, values []*ptypes.Row) error
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
	Type ptypes.OID
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
func (t *Table) validateValueType(value []byte, expectedType ptypes.OID) error {
	return serializers.ValidateGeneric(value, expectedType)
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
	if found && len(seqValueStr) > 0 {
		if parsed, err := strconv.ParseUint(string(seqValueStr), 10, 64); err == nil {
			seqValue = parsed
		}
	}

	logger.Debugf("Using sequence prefix key %s with value %d", sequencePrefixKey, seqValue)
	tx.Write([]byte(sequencePrefixKey), []byte(strconv.FormatUint(seqValue+1, 10)))
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
	if !found || len(metaDataStr) == 0 {
		return nil, fmt.Errorf("table %s does not exist", t.Name)
	}

	var meta TableMetadata
	err := json.Unmarshal([]byte(metaDataStr), &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func (t *Table) GetSerializerSchema(tx *storage.DistributedTransactionVClock) (*serializers.BaseSchema, error) {
	meta, err := t.GetTableMetadataInTx(tx)
	if err != nil {
		return nil, err
	}

	type indexedCol struct {
		name string
		meta ColumnMetadata
	}
	cols := make([]indexedCol, 0, len(meta.Columns))
	for colName, colMeta := range meta.Columns {
		cols = append(cols, indexedCol{name: colName, meta: colMeta})
	}
	sort.Slice(cols, func(i, j int) bool {
		return cols[i].meta.Idx < cols[j].meta.Idx
	})

	schemaFields := make([]serializers.FieldDef, 0, len(cols))
	for _, col := range cols {
		schemaFields = append(schemaFields, serializers.FieldDef{
			Name: col.name,
			OID:  col.meta.Type,
		})
	}

	return &serializers.BaseSchema{
		Fields: schemaFields,
	}, nil
}
