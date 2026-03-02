package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"strings"
)

func ParsePrimaryColName(
	allocator pmem.Allocator,
	colName parser.IColumnNameContext,
	row *table.ResultRow,
) (rmodels.Expression, error) {
	colNameText := colName.GetText()

	// Parse qualified name like "c.name" or "table.column"
	var actualColName string
	var fullQualifiedName string
	if qn := colName.QualifiedName(); qn != nil {
		parts := qn.AllNamePart()
		if len(parts) >= 2 {
			// table_alias.column_name
			tablePrefix := parts[0].GetText()
			columnName := parts[len(parts)-1].GetText()
			actualColName = columnName
			fullQualifiedName = tablePrefix + "." + columnName
		} else if len(parts) == 1 {
			// just column_name
			actualColName = parts[0].GetText()
			fullQualifiedName = actualColName
		} else {
			actualColName = colNameText
			fullQualifiedName = colNameText
		}
	} else {
		actualColName = colNameText
		fullQualifiedName = colNameText
	}

	// Look up column in row - try full qualified name first, then just column name
	if row != nil {
		// Fast path: use the schema nameIdx which already has both bare names and
		// qualified names (tableAlias.colName) pre-indexed by RebuildIndex().
		// This correctly resolves e.g. "t.typelem" when TableAlias="t", Name="typelem".
		if i, ok := row.Schema.FieldIndex(fullQualifiedName); ok {
			fields := []serializers.FieldDef{{
				Name: row.Schema.Fields[i].Name,
				OID:  row.Schema.Fields[i].OID,
			}}
			resultSchema := serializers.NewBaseSchema(fields)
			buf, _, err := row.Schema.GetField(row.Row, i)
			if err != nil {
				return nil, fmt.Errorf("[ParsePrimaryColName] failed to get column value: %v", err)
			}
			resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
			if err != nil {
				return nil, fmt.Errorf("[ParsePrimaryColName] failed to pack column value: %v", err)
			}
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:   []*ptypes.Row{resultRow},
					Schema: resultSchema,
				},
			}, nil
		}

		// Сначала ищем по полному qualified имени (например, "c.id" или "o.user_id")
		for i, col := range row.Schema.Fields {
			// Проверяем прямое совпадение с col.Name
			if col.Name == fullQualifiedName {
				if i < len(row.Schema.Fields) {
					fields := []serializers.FieldDef{{
						Name: col.Name,
						OID:  col.OID,
					}}
					resultSchema := serializers.NewBaseSchema(fields)
					buf, _, err := row.Schema.GetField(row.Row, i)
					if err != nil {
						return nil, fmt.Errorf("[ParsePrimaryColName] failed to get column value: %v", err)
					}
					resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
					if err != nil {
						return nil, fmt.Errorf("[ParsePrimaryColName] failed to pack column value: %v", err)
					}
					return &rmodels.ResultRowsExpression{
						Row: &table.ExecuteResult{
							Rows:   []*ptypes.Row{resultRow},
							Schema: resultSchema,
						},
					}, nil
				}
			}

			// Проверяем составное имя с префиксом (например table.column)
			// В FieldDef нет TableIdentifier, так что просто ищем колонку с полным именем
			if strings.Contains(col.Name, ".") {
				if col.Name == fullQualifiedName {
					if i < len(row.Schema.Fields) {
						fields := []serializers.FieldDef{{
							Name: col.Name,
							OID:  col.OID,
						}}
						resultSchema := serializers.NewBaseSchema(fields)
						buf, _, err := row.Schema.GetField(row.Row, i)
						if err != nil {
							return nil, fmt.Errorf("[ParsePrimaryColName] failed to get column value: %v", err)
						}
						resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
						if err != nil {
							return nil, fmt.Errorf("[ParsePrimaryColName] failed to pack column value: %v", err)
						}
						return &rmodels.ResultRowsExpression{
							Row: &table.ExecuteResult{
								Rows:   []*ptypes.Row{resultRow},
								Schema: resultSchema,
							},
						}, nil
					}
				}
			}
		}

		// Проверяем, не использует ли пользователь оригинальное имя таблицы вместо алиаса
		// В новой структуре нет OriginalTableName и TableIdentifier в FieldDef
		// Эта проверка на алиасы пока пропускается
		// TODO: добавить метаданные о таблицах и алиасах в схему, если требуется

		// Если не нашли точное совпадение и имя не квалифицированное (без префикса),
		// ищем по имени колонки без префикса (например, "amount" найдет "o.amount")
		if qn := colName.QualifiedName(); qn != nil {
			parts := qn.AllNamePart()
			if len(parts) == 1 {
				// Простое имя без префикса - ищем среди всех колонок
				columnName := parts[0].GetText()
				var matchedCols []int
				var matchedColNames []string

				for i, col := range row.Schema.Fields {
					// Проверяем прямое совпадение с именем колонки или с последней частью
					if col.Name == columnName {
						matchedCols = append(matchedCols, i)
						matchedColNames = append(matchedColNames, col.Name)
					} else if strings.Contains(col.Name, ".") {
						// Check if last part matches (e.g., "o.amount" matches "amount")
						colParts := strings.Split(col.Name, ".")
						if colParts[len(colParts)-1] == columnName {
							matchedCols = append(matchedCols, i)
							matchedColNames = append(matchedColNames, col.Name)
						}
					}
				}

				if len(matchedCols) == 1 {
					// Найдена ровно одна колонка - используем её
					i := matchedCols[0]
					if i < len(row.Schema.Fields) {
						fields := []serializers.FieldDef{{
							Name: row.Schema.Fields[i].Name,
							OID:  row.Schema.Fields[i].OID,
						}}
						resultSchema := serializers.NewBaseSchema(fields)
						buf, _, err := row.Schema.GetField(row.Row, i)
						if err != nil {
							return nil, fmt.Errorf("[ParsePrimaryColName] failed to get column value: %v", err)
						}
						resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
						if err != nil {
							return nil, fmt.Errorf("[ParsePrimaryColName] failed to pack column value: %v", err)
						}
						return &rmodels.ResultRowsExpression{
							Row: &table.ExecuteResult{
								Rows:   []*ptypes.Row{resultRow},
								Schema: resultSchema,
							},
						}, nil
					}
				} else if len(matchedCols) > 1 {
					// Найдено несколько колонок - неоднозначность
					return nil, fmt.Errorf(
						"[ParsePrimaryColName] column reference \"%s\" is ambiguous\nHINT: Could refer to: %s",
						columnName,
						strings.Join(matchedColNames, ", "),
					)
				}
			}
		}
	}

	// Если это qualified имя (например, "nsp.nspname"), выдаем ошибку о неизвестной таблице
	if strings.Contains(fullQualifiedName, ".") {
		parts := strings.Split(fullQualifiedName, ".")
		if len(parts) == 2 {
			return nil, fmt.Errorf("[ParsePrimaryColName] missing FROM-clause entry for table \"%s\"", parts[0])
		}
	}

	// Для простых имен выдаем ошибку о несуществующей колонке
	return nil, fmt.Errorf("[ParsePrimaryColName] column \"%s\" does not exist", actualColName)
}
