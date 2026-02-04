package rparser

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"strconv"
	"strings"
	"time"
)

// parsePrimaryExpression handles the basic expressions
// Always returns ResultRowExpression
func ParsePrimaryExpression(primExpr parser.IPrimaryExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	// logger.Debug("ParsePrimaryExpression")
	if primExpr == nil {
		return nil, nil
	}

	// Check for parenthesized expression
	if primExpr.Expression() != nil {
		return ParseExpression(primExpr.Expression(), row)
	}

	// Check for CASE expression
	if primExpr.CaseExpression() != nil {
		return ParseCaseExpression(primExpr.CaseExpression())
	}

	// Check for function call
	if primExpr.FunctionCall() != nil {
		fc := primExpr.FunctionCall()
		funcName := ""
		if qn := fc.QualifiedName(); qn != nil {
			nameParts := qn.AllNamePart()
			if len(nameParts) > 0 {
				parts := []string{}
				for _, np := range nameParts {
					parts = append(parts, np.GetText())
				}
				funcName = parts[len(parts)-1]
			}
		}
		var args []interface{}
		for _, argExpr := range fc.AllFunctionArg() {
			if argExpr.STAR() != nil {
				// Handle COUNT(*) case - pass a special marker value
				args = append(args, "*")
				continue
			}

			if argExpr.Expression() == nil {
				return nil, fmt.Errorf("unsupported function argument: %s", argExpr.GetText())
			}

			val, err := ParseExpression(argExpr.Expression(), row)
			if err != nil {
				return nil, err
			}
			// TODO переделать передачу аргументов функций
			// на строго типизированную передачу Expression
			if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok {
				if len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
					args = append(args, valExpr.Row.Rows[0][0])
				} else {
					args = append(args, nil)
				}
			}
			// args = append(args, ParseExpression(argExpr, row))
		}
		value, _ := functions.ExecuteFunction(funcName, args)
		// If function returns a single boolean value, return BoolExpression for WHERE conditions
		if len(value.Rows) == 1 && len(value.Rows[0]) == 1 && len(value.Columns) == 1 && value.Columns[0].Type == table.ColTypeBool {
			if boolVal, ok := value.Rows[0][0].(bool); ok {
				return &rmodels.BoolExpression{Value: boolVal}, nil
			}
		}
		return &rmodels.ResultRowsExpression{Row: value}, nil
	}

	// Check for extract function
	if primExpr.ExtractFunction() != nil {
		return ParseExtractFunction(primExpr.ExtractFunction(), row)
	}

	// Check for literal values
	if primExpr.Value() != nil {
		valueCtx := primExpr.Value()
		if valueCtx.NUMBER() != nil {
			numStr := valueCtx.NUMBER().GetText()
			if strings.Contains(numStr, ".") {
				if val, err := strconv.ParseFloat(numStr, 64); err == nil {
					return &rmodels.ResultRowsExpression{
						Row: &table.ExecuteResult{
							Rows: [][]interface{}{{val}},
							Columns: []table.TableColumn{
								{Idx: 0, Name: "?column?", Type: table.ColTypeFloat},
							},
						},
					}, nil
				}
			} else {
				if val, err := strconv.Atoi(numStr); err == nil {
					return &rmodels.ResultRowsExpression{
						Row: &table.ExecuteResult{
							Rows: [][]interface{}{{val}},
							Columns: []table.TableColumn{
								{Idx: 0, Name: "?column?", Type: table.ColTypeInt},
							},
						},
					}, nil
				}
			}
		}

		if valueCtx.STRING_LITERAL() != nil {
			str := valueCtx.STRING_LITERAL().GetText()
			// Remove quotes
			if len(str) >= 2 && str[0] == '\'' && str[len(str)-1] == '\'' {
				// return str[1 : len(str)-1]
				str = str[1 : len(str)-1]
				return &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows: [][]interface{}{{str}},
						Columns: []table.TableColumn{
							{Idx: 0, Name: "?column?", Type: table.ColTypeString},
						},
					},
				}, nil
			}
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows: [][]interface{}{{str}},
					Columns: []table.TableColumn{
						{Idx: 0, Name: "?column?", Type: table.ColTypeString},
					},
				},
			}, nil
		}

		if valueCtx.TRUE() != nil {
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows: [][]interface{}{{true}},
					Columns: []table.TableColumn{
						{Idx: 0, Name: "?column?", Type: table.ColTypeBool},
					},
				},
			}, nil
		}

		if valueCtx.FALSE() != nil {
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows: [][]interface{}{{false}},
					Columns: []table.TableColumn{
						{Idx: 0, Name: "?column?", Type: table.ColTypeBool},
					},
				},
			}, nil
		}

		if valueCtx.CURRENT_TIMESTAMP() != nil {
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows: [][]interface{}{{time.Now().Format("2006-01-02 15:04:05")}},
					Columns: []table.TableColumn{
						{Idx: 0, Name: "?column?", Type: table.ColTypeString},
					},
				},
			}, nil
		}
	}

	// Check for column name
	if primExpr.ColumnName() != nil {
		colNameText := primExpr.ColumnName().GetText()

		// Parse qualified name like "c.name" or "table.column"
		var actualColName string
		var fullQualifiedName string
		if qn := primExpr.ColumnName().QualifiedName(); qn != nil {
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
			// Сначала ищем по полному qualified имени (например, "c.id")
			for i, col := range row.Columns {
				if col.Name == fullQualifiedName {
					if i < len(row.Row) {
						return &rmodels.ResultRowsExpression{
							Row: &table.ExecuteResult{
								Rows: [][]interface{}{{row.Row[i]}},
								Columns: []table.TableColumn{
									{Idx: 0, Name: col.Name, Type: col.Type},
								},
							},
						}, nil
					}
				}
			}

			// Проверяем, не использует ли пользователь оригинальное имя таблицы вместо алиаса
			if qn := primExpr.ColumnName().QualifiedName(); qn != nil {
				parts := qn.AllNamePart()
				if len(parts) >= 2 {
					tablePrefix := parts[0].GetText()
					columnName := parts[len(parts)-1].GetText()

					// Ищем колонку с таким именем, но с другим TableIdentifier
					for _, col := range row.Columns {
						if col.OriginalTableName == tablePrefix && col.TableIdentifier != tablePrefix {
							// Проверяем, что имя колонки совпадает
							colParts := strings.Split(col.Name, ".")
							if len(colParts) >= 2 && colParts[len(colParts)-1] == columnName {
								return nil, fmt.Errorf(
									"invalid reference to FROM-clause entry for table \"%s\"\nHINT: Perhaps you meant to reference the table alias \"%s\"",
									tablePrefix,
									col.TableIdentifier,
								)
							}
						}
					}
				}
			}

			// Если не нашли точное совпадение и имя не квалифицированное (без префикса),
			// ищем по имени колонки без префикса (например, "amount" найдет "o.amount")
			if qn := primExpr.ColumnName().QualifiedName(); qn != nil {
				parts := qn.AllNamePart()
				if len(parts) == 1 {
					// Простое имя без префикса - ищем среди всех колонок
					columnName := parts[0].GetText()
					var matchedCols []int
					var matchedColNames []string

					for i, col := range row.Columns {
						// Проверяем, заканчивается ли имя колонки на искомое имя
						colParts := strings.Split(col.Name, ".")
						if len(colParts) >= 2 && colParts[len(colParts)-1] == columnName {
							matchedCols = append(matchedCols, i)
							matchedColNames = append(matchedColNames, col.Name)
						} else if col.Name == columnName {
							// Колонка без префикса (может быть в случае без JOIN)
							matchedCols = append(matchedCols, i)
							matchedColNames = append(matchedColNames, col.Name)
						}
					}

					if len(matchedCols) == 1 {
						// Найдена ровно одна колонка - используем её
						i := matchedCols[0]
						if i < len(row.Row) {
							return &rmodels.ResultRowsExpression{
								Row: &table.ExecuteResult{
									Rows: [][]interface{}{{row.Row[i]}},
									Columns: []table.TableColumn{
										{Idx: 0, Name: row.Columns[i].Name, Type: row.Columns[i].Type},
									},
								},
							}, nil
						}
					} else if len(matchedCols) > 1 {
						// Найдено несколько колонок - неоднозначность
						return nil, fmt.Errorf(
							"column reference \"%s\" is ambiguous\nHINT: Could refer to: %s",
							columnName,
							strings.Join(matchedColNames, ", "),
						)
					}
				}
			}
		}

		// Column not found - return nil or error
		logger.Debugf("Column %s (qualified: %s) not found in row context", actualColName, fullQualifiedName)
		return &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows: [][]interface{}{{nil}},
				Columns: []table.TableColumn{
					{Idx: 0, Name: actualColName, Type: table.ColTypeString},
				},
			},
		}, nil
	}

	// If nothing matches, return the text
	text := primExpr.GetText()
	return &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows: [][]interface{}{{text}},
			Columns: []table.TableColumn{
				{Idx: 0, Name: "?column?", Type: table.ColTypeString},
			},
		},
	}, nil
}
