package table

import (
	"encoding/json"
	"fmt"
	"petacore/internal/core"
	"petacore/internal/logger"
	"petacore/internal/storage"
	"sort"

	"go.uber.org/zap"
)

// Select выполняет SELECT запрос
func (t *Table) Select(
	tx *storage.DistributedTransactionVClock,
	tableName string,
	columns []SelectColumn,
	where map[string]interface{},
	limit int,
) (*ExecuteResult, error) {

	// Получаем метаданные таблицы
	tableKey := t.getMetadataPrefixKey()
	metaStr, found := tx.Read([]byte(tableKey))
	if !found || metaStr == "" {
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

	// Сканируем все строки
	prefix := t.getRowPrefixKey()
	// TODO сюда можно захардкодить системные таблицы, чтобы много не переписывать
	kvMap, err := tx.Scan([]byte(prefix), core.IteratorTypeAll, limit)
	if err != nil {
		return nil, err
	}

	results := make([][]interface{}, 0, 10)

	for _, value := range kvMap {
		// TODO пересмотреть концепцию хранения строк, нужна унификация
		var singleRow []interface{}
		if err := json.Unmarshal([]byte(value), &singleRow); err != nil {
			logger.Error("Failed to unmarshal", zap.Error(err))
			continue
		}
		logger.Debugf("Single row data: %+v\n", singleRow)

		resultRow := make([]interface{}, 0, len(finalColumns))

		for _, finalColumn := range finalColumns {
			if tableColumn, ok := tableColumnsMap[finalColumn.Name]; ok {
				resultRow = append(resultRow, singleRow[tableColumn.Idx-1])
				continue
			}
			return nil, fmt.Errorf("column %s not found in table columns map", finalColumn.Name)
		}

		results = append(results, resultRow)
	}

	// Применяем LIMIT
	// TODO убрать на уровень executor
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	execResult := &ExecuteResult{
		Rows:    results,
		Columns: finalColumns,
	}
	logger.Debug("Select result: ", zap.Any("result", execResult))

	return execResult, nil
}
