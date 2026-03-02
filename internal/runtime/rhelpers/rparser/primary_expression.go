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
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"strconv"
	"strings"
	"time"
)

// parsePrimaryExpression handles the basic expressions
// Always returns ResultRowExpression
func ParsePrimaryExpression(
	allocator pmem.Allocator,
	ctx context.Context,
	primExpr parser.IPrimaryExpressionContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (rmodels.Expression, error) {
	// Check for parenthesized expression
	if primExpr.Expression() != nil {
		return ParseExpression(allocator, ctx, primExpr.Expression(), row, subExec)
	}

	// Check for CASE expression
	if primExpr.CaseExpression() != nil {
		return ParseCaseExpression(allocator, ctx, primExpr.CaseExpression())
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

			val, err := ParseExpression(allocator, ctx, argExpr.Expression(), row, subExec)
			if err != nil {
				return nil, err
			}
			// TODO переделать передачу аргументов функций
			// на строго типизированную передачу Expression
			if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok {
				if len(valExpr.Row.Rows) == 1 && len(valExpr.Row.Schema.Fields) == 1 {
					buf, oid, err := valExpr.Row.Schema.GetField(valExpr.Row.Rows[0], 0)
					if err != nil {
						return nil, fmt.Errorf("failed to get field: %w", err)
					}
					desVal, err := serializers.DeserializeGeneric(buf, oid)
					if err != nil {
						return nil, fmt.Errorf("failed to deserialize arg: %w", err)
					}
					args = append(args, desVal)
				} else {
					return nil, fmt.Errorf("function arguments must be single-row single-column expressions")
				}
				// if len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
				// 	args = append(args, valExpr.Row.Rows[0][0])
				// } else {
				// 	args = append(args, nil)
				// }
			}
			// args = append(args, ParseExpression(argExpr, row))
		}

		value, err := functions.ExecuteFunctionWithContext(allocator, ctx, funcName, args)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, fmt.Errorf("function %s returned nil result", funcName)
		}
		// If function returns a single boolean value, return BoolExpression for WHERE conditions
		if len(value.Rows) == 1 && len(value.Schema.Fields) == 1 && value.Schema.Fields[0].OID == ptypes.PTypeBool {
			field, _, err := value.Schema.GetField(value.Rows[0], 0)
			if err != nil {
				return nil, fmt.Errorf("error getting field from function result: %v", err)
			}
			boolVal, err := serializers.BoolSerializerInstance.Deserialize(field)
			if err != nil {
				return nil, fmt.Errorf("error deserializing boolean value from function result: %v", err)
			}
			return &rmodels.BoolExpression{Value: boolVal.IntoGo()}, nil
		}
		return &rmodels.ResultRowsExpression{Row: value}, nil
	}

	// Check for extract function
	if primExpr.ExtractFunction() != nil {
		return ParseExtractFunction(allocator, ctx, primExpr.ExtractFunction(), row, subExec)
	}

	// Check for literal values
	if primExpr.Value() != nil {
		valueCtx := primExpr.Value()
		if valueCtx.NUMBER() != nil {
			numStr := valueCtx.NUMBER().GetText()
			if strings.Contains(numStr, ".") {
				if val, err := strconv.ParseFloat(numStr, 64); err == nil {
					fields := []serializers.FieldDef{{
						Name: "?column?",
						OID:  ptypes.PTypeFloat8,
					}}
					resultSchema := serializers.NewBaseSchema(fields)
					serVal, err := serializers.Float8SerializerInstance.Serialize(allocator, val)
					if err != nil {
						return nil, fmt.Errorf("failed to serialize float value: %v", err)
					}
					resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
					if err != nil {
						return nil, fmt.Errorf("failed to pack float value: %v", err)
					}
					return &rmodels.ResultRowsExpression{
						Row: &table.ExecuteResult{
							Rows:   []*ptypes.Row{resultRow},
							Schema: resultSchema,
						},
					}, nil
				}
			} else {
				if val, err := strconv.Atoi(numStr); err == nil {
					fields := []serializers.FieldDef{{
						Name: "?column?",
						OID:  ptypes.PTypeInt4,
					}}
					resultSchema := serializers.NewBaseSchema(fields)
					serVal, err := serializers.Int4SerializerInstance.Serialize(allocator, int32(val))
					if err != nil {
						return nil, fmt.Errorf("failed to serialize integer value: %v", err)
					}
					resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
					if err != nil {
						return nil, fmt.Errorf("failed to pack integer value: %v", err)
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

		if valueCtx.STRING_LITERAL() != nil {
			str := valueCtx.STRING_LITERAL().GetText()
			// Remove quotes
			if len(str) >= 2 && str[0] == '\'' && str[len(str)-1] == '\'' {
				str = str[1 : len(str)-1]
			}
			fields := []serializers.FieldDef{{
				Name: "?column?",
				OID:  ptypes.PTypeText,
			}}
			resultSchema := serializers.NewBaseSchema(fields)
			serVal, err := serializers.TextSerializerInstance.Serialize(allocator, str)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize string value: %v", err)
			}
			resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
			if err != nil {
				return nil, fmt.Errorf("failed to pack string value: %v", err)
			}
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:   []*ptypes.Row{resultRow},
					Schema: resultSchema,
				},
			}, nil
		}

		if valueCtx.TRUE() != nil {
			fields := []serializers.FieldDef{{
				Name: "?column?",
				OID:  ptypes.PTypeBool,
			}}
			resultSchema := serializers.NewBaseSchema(fields)
			val, err := serializers.BoolSerializerInstance.Serialize(allocator, true)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize boolean value: %v", err)
			}
			resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
			if err != nil {
				return nil, fmt.Errorf("failed to pack boolean value: %v", err)
			}
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:   []*ptypes.Row{resultRow},
					Schema: resultSchema,
				},
			}, nil
		}

		if valueCtx.FALSE() != nil {
			fields := []serializers.FieldDef{{
				Name: "?column?",
				OID:  ptypes.PTypeBool,
			}}
			resultSchema := serializers.NewBaseSchema(fields)
			val, err := serializers.BoolSerializerInstance.Serialize(allocator, false)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize boolean value: %v", err)
			}
			resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
			if err != nil {
				return nil, fmt.Errorf("failed to pack boolean value: %v", err)
			}
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:   []*ptypes.Row{resultRow},
					Schema: resultSchema,
				},
			}, nil
		}

		if valueCtx.CURRENT_TIMESTAMP() != nil {
			fields := []serializers.FieldDef{{
				Name: "?column?",
				OID:  ptypes.PTypeTimestamp,
			}}
			resultSchema := serializers.NewBaseSchema(fields)
			timeVal := time.Now()
			val, err := serializers.TimestampSerializerInstance.Serialize(allocator, &timeVal)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize timestamp value: %v", err)
			}
			resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
			if err != nil {
				return nil, fmt.Errorf("failed to pack timestamp value: %v", err)
			}
			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:   []*ptypes.Row{resultRow},
					Schema: resultSchema,
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
			// TODO not supported in typesystem
			// if typedLit.INTERVAL_TYPE() != nil {
			// 	// Parse interval string
			// 	intervalMicros, err := utils.ParseInterval(str)
			// 	if err != nil {
			// 		return nil, fmt.Errorf("invalid interval: %v", err)
			// 	}
			// 	return &rmodels.ResultRowsExpression{
			// 		Row: &table.ExecuteResult{
			// 			Rows: [][]interface{}{{intervalMicros}},
			// 			Columns: []table.TableColumn{
			// 				{Idx: 0, Name: "?column?", Type: table.ColTypeInterval},
			// 			},
			// 		},
			// 	}, nil
			// } else
			if typedLit.DATE_TYPE() != nil || typedLit.TIMESTAMP_TYPE() != nil {
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
				colType := ptypes.PTypeTimestamp
				if typedLit.TIMESTAMP_TYPE() != nil {
					colType = ptypes.PTypeTimestampz
				}
				fields := []serializers.FieldDef{{
					Name: "?column?",
					OID:  colType,
				}}
				resultSchema := serializers.NewBaseSchema(fields)
				val, err := serializers.TimestampSerializerInstance.Serialize(allocator, &t)
				if err != nil {
					return nil, fmt.Errorf("failed to serialize timestamp value: %v", err)
				}
				resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
				if err != nil {
					return nil, fmt.Errorf("failed to pack timestamp value: %v", err)
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
					return nil, fmt.Errorf("failed to get column value: %v", err)
				}
				resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
				if err != nil {
					return nil, fmt.Errorf("failed to pack column value: %v", err)
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
							return nil, fmt.Errorf("failed to get column value: %v", err)
						}
						resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
						if err != nil {
							return nil, fmt.Errorf("failed to pack column value: %v", err)
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
								return nil, fmt.Errorf("failed to get column value: %v", err)
							}
							resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
							if err != nil {
								return nil, fmt.Errorf("failed to pack column value: %v", err)
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
			if qn := primExpr.ColumnName().QualifiedName(); qn != nil {
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
								return nil, fmt.Errorf("failed to get column value: %v", err)
							}
							resultRow, err := resultSchema.Pack(allocator, [][]byte{buf})
							if err != nil {
								return nil, fmt.Errorf("failed to pack column value: %v", err)
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
	fields := []serializers.FieldDef{{
		Name: "?column?",
		OID:  ptypes.PTypeText,
	}}
	resultSchema := serializers.NewBaseSchema(fields)
	serVal, err := serializers.TextSerializerInstance.Serialize(allocator, text)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize text value: %v", err)
	}
	resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
	if err != nil {
		return nil, fmt.Errorf("failed to pack text value: %v", err)
	}
	return &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows:   []*ptypes.Row{resultRow},
			Schema: resultSchema,
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
