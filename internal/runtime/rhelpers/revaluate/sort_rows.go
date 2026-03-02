package revaluate

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"sort"
	"strings"

	"go.uber.org/zap"
)

// sortRows sorts the rows based on OrderBy items
func EvaluateSortRows(allocator pmem.Allocator, execResult *table.ExecuteResult, orderBy []items.OrderByItem) error {
	if len(orderBy) == 0 {
		return nil
	}

	// Log the schema fields and order-by items for debugging
	fields := make([]string, len(execResult.Schema.Fields))
	for i, f := range execResult.Schema.Fields {
		fields[i] = f.Name
	}
	logger.Debug("EvaluateSortRows",
		zap.Strings("schemaFields", fields),
		zap.Any("orderBy", orderBy),
	)

	resolveIdx := func(ob items.OrderByItem) (int, error) {
		if ob.ColumnIndex > 0 {
			idx := ob.ColumnIndex - 1
			if idx < len(execResult.Schema.Fields) {
				logger.Debug("EvaluateSortRows: resolved by index", zap.Int("colIndex", ob.ColumnIndex), zap.Int("fieldIdx", idx))
				return idx, nil
			}
			logger.Debug("EvaluateSortRows: column index out of range", zap.Int("colIndex", ob.ColumnIndex), zap.Int("numFields", len(execResult.Schema.Fields)))
			return -1, fmt.Errorf("[EvaluateSortRows] ORDER BY position %d is not in select list", ob.ColumnIndex)
		}
		if ob.ColumnName != "" {
			lower := strings.ToLower(ob.ColumnName)
			// Qualified name like "t.oid" or "pg_type.oid": match by TableAlias.Name first.
			// If no match (e.g. Project already cleared TableAlias), fall back to bare-name lookup.
			if dotIdx := strings.LastIndex(lower, "."); dotIdx >= 0 {
				tAlias := lower[:dotIdx]
				colName := lower[dotIdx+1:]
				for i, field := range execResult.Schema.Fields {
					if strings.ToLower(field.Name) == colName &&
						strings.ToLower(field.TableAlias) == tAlias {
						logger.Debug("EvaluateSortRows: resolved qualified", zap.String("colName", ob.ColumnName), zap.Int("fieldIdx", i))
						return i, nil
					}
				}
				// Qualified match failed — fallback to bare-name (TableAlias may have been
				// cleared by the Project node). Error if the bare name is ambiguous.
				var fallback []int
				for i, field := range execResult.Schema.Fields {
					if strings.ToLower(field.Name) == colName {
						fallback = append(fallback, i)
					}
				}
				if len(fallback) == 1 {
					logger.Debug("EvaluateSortRows: resolved qualified via bare-name fallback", zap.String("colName", ob.ColumnName), zap.Int("fieldIdx", fallback[0]))
					return fallback[0], nil
				}
				if len(fallback) > 1 {
					return -1, fmt.Errorf("[EvaluateSortRows] column reference \"%s\" is ambiguous", ob.ColumnName)
				}
				logger.Debug("EvaluateSortRows: qualified name not found", zap.String("colName", ob.ColumnName))
				return -1, fmt.Errorf("[EvaluateSortRows] column \"%s\" does not exist", ob.ColumnName)
			}
			// Unqualified: match by bare Name — reject if ambiguous.
			var matched []int
			for i, field := range execResult.Schema.Fields {
				if strings.ToLower(field.Name) == lower {
					matched = append(matched, i)
				}
			}
			logger.Debug("EvaluateSortRows: resolved by name", zap.String("colName", ob.ColumnName), zap.Ints("matchedIndices", matched))
			if len(matched) == 1 {
				return matched[0], nil
			}
			if len(matched) > 1 {
				return -1, fmt.Errorf("[EvaluateSortRows] column reference \"%s\" is ambiguous", ob.ColumnName)
			}
		}
		logger.Debug("EvaluateSortRows: could not resolve column", zap.Any("ob", ob))
		return -1, fmt.Errorf("[EvaluateSortRows] column \"%s\" does not exist", ob.ColumnName)
	}

	// Pre-compute sort keys once per row (O(n*k)), then sort a combined slice
	// so that keys and rows stay in sync as sort.Slice reorders them.
	type sortKey = []ptypes.BaseType[any]
	type indexedRow struct {
		row  *ptypes.Row
		keys sortKey
	}

	// Pre-resolve all ORDER BY column indices before touching any rows.
	// This way a bad column reference returns an error immediately.
	colIndices := make([]int, len(orderBy))
	for k, ob := range orderBy {
		idx, resolveErr := resolveIdx(ob)
		if resolveErr != nil {
			return fmt.Errorf("[EvaluateSortRows] %w", resolveErr)
		}
		colIndices[k] = idx
	}

	indexed := make([]indexedRow, len(execResult.Rows))
	for i, row := range execResult.Rows {
		sk := make(sortKey, len(orderBy))
		for j := range orderBy {
			colIdx := colIndices[j]
			if colIdx < 0 {
				continue
			}
			buf, oid, getErr := execResult.Schema.GetField(row, colIdx)
			if getErr != nil {
				continue
			}
			sk[j], _ = serializers.DeserializeGeneric(buf, oid)
		}
		indexed[i] = indexedRow{row: row, keys: sk}
	}

	sort.Slice(indexed, func(i, j int) bool {
		for k, ob := range orderBy {
			vi, vj := indexed[i].keys[k], indexed[j].keys[k]
			if vi == nil || vj == nil {
				continue
			}
			if cmp := vi.Compare(vj); cmp != 0 {
				return cmp < 0 != (ob.Direction == "DESC")
			}
		}
		return false
	})

	// Write sorted rows back
	for i, ir := range indexed {
		execResult.Rows[i] = ir.row
	}

	return nil
}
