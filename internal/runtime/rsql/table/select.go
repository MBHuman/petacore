package table

import (
	"encoding/json"
	"fmt"
	"petacore/internal/core"
	"petacore/internal/logger"
	"petacore/internal/storage"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"sort"

	"go.uber.org/zap"
)

// Select выполняет SELECT запрос
func (t *Table) Select(
	allocator pmem.Allocator,
	tx *storage.DistributedTransactionVClock,
	tableName string,
	columns []SelectColumn,
	where map[string]interface{},
	limit int,
) (*ExecuteResult, error) {

	// Получаем метаданные таблицы
	tableKey := t.getMetadataPrefixKey()
	metaStr, found := tx.Read([]byte(tableKey))
	if !found || len(metaStr) == 0 {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	var meta TableMetadata
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		return nil, err
	}

	logger.Debug("Table metadata: ", zap.Any("meta", meta))

	tableColumns := make([]TableColumn, 0, len(meta.Columns))

	for colName, colMeta := range meta.Columns {
		tableColumns = append(tableColumns, TableColumn{
			Idx:  colMeta.Idx,
			Name: colName,
			Type: colMeta.Type,
		})
	}

	tableColumnsMap := make(map[string]TableColumn)
	for _, col := range tableColumns {
		tableColumnsMap[col.Name] = col
	}

	sort.Slice(tableColumns, func(i, j int) bool {
		return tableColumns[i].Idx < tableColumns[j].Idx
	})

	// Build schemaFields in the same Idx order as the stored rows
	schemaFields := make([]serializers.FieldDef, 0, len(tableColumns))
	for _, col := range tableColumns {
		schemaFields = append(schemaFields, serializers.FieldDef{
			Name: col.Name,
			OID:  col.Type,
		})
	}

	logger.Debug("Table columns", zap.Any("columns", tableColumns))
	logger.Debug("Select columns", zap.Any("columns", columns))

	finalColumns := make([]TableColumn, 0, 5)

	idxCnt := 1

	for _, selectColumn := range columns {
		if selectColumn.IsAll {
			for _, col := range tableColumns {
				finalColumns = append(finalColumns, col)
			}
		} else {
			if colMeta, ok := meta.Columns[selectColumn.Name]; ok {
				finalColumns = append(finalColumns, TableColumn{
					Idx:  idxCnt,
					Name: selectColumn.Name,
					Type: colMeta.Type,
				})
			} else {
				return nil, fmt.Errorf("column %s does not exist", selectColumn.Name)
			}
		}
		idxCnt++
	}
	schema := serializers.NewBaseSchema(schemaFields)

	// Сканируем все строки
	prefix := t.getRowPrefixKey()
	kvMap, err := tx.Scan([]byte(prefix), core.IteratorTypeAll, limit)
	if err != nil {
		return nil, err
	}

	results := make([]*ptypes.Row, 0) // TODO переделать на итератор

	for _, value := range kvMap {
		results = append(results, ptypes.RowFactory(value))
	}

	execResult := &ExecuteResult{
		Rows:   results,
		Schema: schema,
	}

	return execResult, nil
}
