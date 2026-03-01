package rparser

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/utils"
	"strconv"
	"strings"
	"time"
)

// parsePrimaryExpression handles the basic expressions
// Always returns ResultRowExpression
func ParsePrimaryExpression(ctx context.Context, primExpr parser.IPrimaryExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	// logger.Debug("ParsePrimaryExpression")

	// Check for parenthesized expression
	if primExpr.Expression() != nil {
		return ParseExpression(ctx, primExpr.Expression(), row, subExec)
	}

	// Check for CASE expression
	if primExpr.CaseExpression() != nil {
		return ParseCaseExpression(ctx, primExpr.CaseExpression())
	}

	// Check for subquery expression
	if primExpr.SubqueryExpression() != nil {
		sqCtx := primExpr.SubqueryExpression()
		selCtx := sqCtx.SelectStatement()
		if selCtx == nil {
			return nil, fmt.Errorf("invalid subquery context")
		}
		selectStmt, err := ParseSelectStatement(selCtx)
		if err != nil {
			return nil, fmt.Errorf("error parsing subquery: %v", err)
		}
		if selectStmt == nil {
			return nil, fmt.Errorf("failed to build select statement from context")
		}
		// Возвращаем выражение-подзапрос
		return &rmodels.SubqueryExpression{Select: selectStmt}, nil
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

			val, err := ParseExpression(ctx, argExpr.Expression(), row, subExec)
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

		value, err := functions.ExecuteFunctionWithContext(ctx, funcName, args)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, fmt.Errorf("function %s returned nil result", funcName)
		}
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
		return ParseExtractFunction(ctx, primExpr.ExtractFunction(), row, subExec)
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
					Rows: [][]interface{}{{time.Now().UnixMicro()}},
					Columns: []table.TableColumn{
						{Idx: 0, Name: "?column?", Type: table.ColTypeTimestampTz},
					},
				},
			}, nil
		}

		// Check for typed literals (DATE, TIMESTAMP, INTERVAL)
		if valueCtx.TypedLiteral() != nil {
			typedLit := valueCtx.TypedLiteral()

			// Get the string literal value
			strLit := typedLit.STRING_LITERAL()
			if strLit == nil {
				return nil, fmt.Errorf("typed literal missing string value")
			}

			str := strLit.GetText()
			// Remove quotes
			if len(str) >= 2 && str[0] == '\'' && str[len(str)-1] == '\'' {
				str = str[1 : len(str)-1]
			}

			// Check which type of literal
			if typedLit.INTERVAL_TYPE() != nil {
				// Parse interval string
				intervalMicros, err := utils.ParseInterval(str)
				if err != nil {
					return nil, fmt.Errorf("invalid interval: %v", err)
				}
				return &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows: [][]interface{}{{intervalMicros}},
						Columns: []table.TableColumn{
							{Idx: 0, Name: "?column?", Type: table.ColTypeInterval},
						},
					},
				}, nil
			} else if typedLit.DATE_TYPE() != nil || typedLit.TIMESTAMP_TYPE() != nil {
				// Parse date/timestamp string
				// TODO: Implement proper date/timestamp parsing
				// For now, try parsing ISO 8601 format
				t, err := time.Parse("2006-01-02", str)
				if err != nil {
					t, err = time.Parse("2006-01-02 15:04:05", str)
					if err != nil {
						return nil, fmt.Errorf("invalid date/timestamp format: %s", str)
					}
				}
				colType := table.ColTypeTimestamp
				if typedLit.TIMESTAMP_TYPE() != nil {
					colType = table.ColTypeTimestampTz
				}
				return &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows: [][]interface{}{{t.UnixMicro()}},
						Columns: []table.TableColumn{
							{Idx: 0, Name: "?column?", Type: colType},
						},
					},
				}, nil
			}
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
			// Сначала ищем по полному qualified имени (например, "c.id" или "o.user_id")
			// todo убрать паразитный цикл на колонки
			for i, col := range row.Columns {
				// Проверяем прямое совпадение с col.Name
				if col.Name == fullQualifiedName {
					if i < len(row.Row) {
						return &rmodels.ResultRowsExpression{
							Row: &table.ExecuteResult{
								Rows: [][]interface{}{{row.Row[i]}},
								Columns: []table.TableColumn{
									{Idx: 0, Name: col.Name, Type: col.Type, TableIdentifier: col.TableIdentifier},
								},
							},
						}, nil
					}
				}

				// Проверяем составное имя TableIdentifier.Name
				if col.TableIdentifier != "" {
					qualifiedColName := col.TableIdentifier + "." + col.Name
					if qualifiedColName == fullQualifiedName {
						if i < len(row.Row) {
							return &rmodels.ResultRowsExpression{
								Row: &table.ExecuteResult{
									Rows: [][]interface{}{{row.Row[i]}},
									Columns: []table.TableColumn{
										{Idx: 0, Name: col.Name, Type: col.Type, TableIdentifier: col.TableIdentifier},
									},
								},
							}, nil
						}
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
						// Проверяем прямое совпадение с именем колонки
						if col.Name == columnName {
							matchedCols = append(matchedCols, i)
							if col.TableIdentifier != "" {
								matchedColNames = append(matchedColNames, col.TableIdentifier+"."+col.Name)
							} else {
								matchedColNames = append(matchedColNames, col.Name)
							}
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
										{Idx: 0, Name: row.Columns[i].Name, Type: row.Columns[i].Type, TableIdentifier: row.Columns[i].TableIdentifier},
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

		// Column not found - return error
		logger.Debugf("Column %s (qualified: %s) not found in row context", actualColName, fullQualifiedName)

		// Если это qualified имя (например, "nsp.nspname"), выдаем ошибку о неизвестной таблице
		if strings.Contains(fullQualifiedName, ".") {
			parts := strings.Split(fullQualifiedName, ".")
			if len(parts) == 2 {
				return nil, fmt.Errorf("missing FROM-clause entry for table \"%s\"", parts[0])
			}
		}

		// Для простых имен выдаем ошибку о несуществующей колонке
		return nil, fmt.Errorf("column \"%s\" does not exist", actualColName)
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

// ExtractPrimaryExpression tries to extract a simple primaryExpression from an expression tree
// Returns nil if the expression is complex (has operators, etc.)
func ExtractPrimaryExpression(expr parser.IExpressionContext) parser.IPrimaryExpressionContext {
	if expr == nil {
		return nil
	}

	// expression -> orExpression
	orExpr := expr.OrExpression()
	if orExpr == nil {
		return nil
	}

	// orExpression -> andExpression (if only one)
	andExprs := orExpr.AllAndExpression()
	if len(andExprs) != 1 {
		return nil // multiple AND expressions
	}

	// andExpression -> notExpression (if only one)
	notExprs := andExprs[0].AllNotExpression()
	if len(notExprs) != 1 {
		return nil // multiple NOT expressions
	}

	notExpr := notExprs[0]
	if notExpr.NOT() != nil {
		return nil // has NOT operator
	}

	// notExpression -> comparisonExpression
	compExpr := notExpr.ComparisonExpression()
	if compExpr == nil {
		return nil
	}

	// comparisonExpression should have no operators
	if compExpr.Operator() != nil || compExpr.IN() != nil || compExpr.LIKE() != nil || compExpr.IS() != nil {
		return nil // has comparison operators
	}

	// comparisonExpression -> concatExpression (if only one)
	concatExprs := compExpr.AllConcatExpression()
	if len(concatExprs) != 1 {
		return nil
	}

	concatExpr := concatExprs[0]
	if len(concatExpr.AllCONCAT()) > 0 {
		return nil // has concatenation
	}

	// concatExpression -> additiveExpression (if only one)
	addExprs := concatExpr.AllAdditiveExpression()
	if len(addExprs) != 1 {
		return nil
	}

	addExpr := addExprs[0]
	if len(addExpr.AllPLUS()) > 0 || len(addExpr.AllMINUS()) > 0 {
		return nil // has addition/subtraction
	}

	// additiveExpression -> multiplicativeExpression (if only one)
	multExprs := addExpr.AllMultiplicativeExpression()
	if len(multExprs) != 1 {
		return nil
	}

	multExpr := multExprs[0]
	if len(multExpr.AllSTAR()) > 0 || len(multExpr.AllSLASH()) > 0 {
		return nil // has multiplication/division
	}

	// multiplicativeExpression -> unaryExpression (if only one)
	unaryExprs := multExpr.AllUnaryExpression()
	if len(unaryExprs) != 1 {
		return nil
	}

	unaryExpr := unaryExprs[0]
	if unaryExpr.PLUS() != nil || unaryExpr.MINUS() != nil {
		return nil // has unary operator
	}

	// unaryExpression -> castExpression
	castExpr := unaryExpr.CastExpression()
	if castExpr == nil {
		return nil
	}

	if len(castExpr.AllPostfix()) > 0 {
		return nil // has postfix (cast, collate, at time zone)
	}

	// castExpression -> primaryExpression
	primExpr := castExpr.PrimaryExpression()

	// atTimeZoneExpression -> primaryExpression
	return primExpr
}
