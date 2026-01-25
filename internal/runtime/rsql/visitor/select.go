package visitor

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"strconv"

	"github.com/antlr4-go/antlr/v4"
)

// TODO пересмотреть, возможно можно проще сделать
func (l *sqlListener) EnterSelectStatement(ctx *parser.SelectStatementContext) {
	stmt := &statements.SelectStatement{}

	// From clause
	if ctx.FromClause() != nil {
		from := &statements.FromClause{}

		if tableFactor := ctx.FromClause().TableFactor(); tableFactor != nil {
			if tableName := tableFactor.TableName(); tableName != nil {
				tableNameText := tableName.GetText()

				from.TableName = tableNameText
				if tableFactor.Alias() != nil {
					from.Alias = tableFactor.Alias().GetText()
				}
			}

			if subSeleqct := tableFactor.SelectStatement(); subSeleqct != nil {
				logger.Debugf("DEBUG: parsing FROM subquery")
				subStmtListener := &sqlListener{}
				antlr.ParseTreeWalkerDefault.Walk(subStmtListener, subSeleqct)
				from.SelectStatement = subStmtListener.stmt.(*statements.SelectStatement)
				if tableFactor.Alias() != nil {
					from.Alias = tableFactor.Alias().GetText()
				}
			}
		}

		// Parse joins
		for _, jc := range ctx.FromClause().AllJoinClause() {
			join := statements.JoinClause{}
			if jc.INNER() != nil {
				join.Type = "INNER"
			} else if jc.LEFT() != nil {
				join.Type = "LEFT"
			} else if jc.RIGHT() != nil {
				join.Type = "RIGHT"
			} else if jc.FULL() != nil {
				join.Type = "FULL"
			} else if jc.CROSS() != nil {
				join.Type = "CROSS"
			} else {
				join.Type = "INNER" // default
			}
			if jc.QualifiedName() != nil {
				join.TableName = jc.QualifiedName().GetText()
			}
			if alias := jc.Alias(); alias != nil {
				if id := alias.IDENTIFIER(); id != nil {
					join.Alias = id.GetText()
				}
			}
			if jc.ON() != nil && jc.Expression() != nil {
				join.OnCondition = jc.Expression()
			}
			from.Joins = append(from.Joins, join)
		}
		stmt.From = from
	}

	// Columns
	if ctx.SelectList() != nil {
		for _, item := range ctx.SelectList().AllSelectItem() {
			logger.Debugf("Parsing select item: %s", item.GetText())
			selectItem := items.SelectItem{}

			if selectAll := item.SelectAll(); selectAll != nil {
				selectItem.IsSelectAll = true
				stmt.Columns = append(stmt.Columns, selectItem)
			} else if expr := item.Expression(); expr != nil {
				// Try to extract simple column or function names for optimization
				if primExpr := extractPrimaryExpression(expr); primExpr != nil {
					if primExpr.ColumnName() != nil {
						fullName := primExpr.ColumnName().GetText()
						// Parse qualified name like "c.name" or "customers.name"
						if qn := primExpr.ColumnName().QualifiedName(); qn != nil {
							parts := qn.AllNamePart()
							if len(parts) == 2 {
								// table_alias.column_name
								selectItem.TableAlias = parts[0].GetText()
								selectItem.ColumnName = parts[1].GetText()
							} else if len(parts) == 1 {
								// just column_name
								selectItem.ColumnName = parts[0].GetText()
							} else {
								// schema.table.column or more complex
								selectItem.ColumnName = fullName
							}
						} else {
							selectItem.ColumnName = fullName
						}
					} else if primExpr.FunctionCall() != nil {
						fc := primExpr.FunctionCall()
						funcCall := &items.FunctionCall{}
						if qn := fc.QualifiedName(); qn != nil {
							nameParts := qn.AllNamePart()
							parts := []string{}
							for _, np := range nameParts {
								parts = append(parts, np.GetText())
							}
							funcCall.Name = parts[len(parts)-1]
						}

						// Args - expression contexts
						for _, argExpr := range fc.AllExpression() {
							funcCall.Args = append(funcCall.Args, argExpr)
						}
						selectItem.Function = funcCall
					} else {
						// Complex expression - save ExpressionContext
						selectItem.ExpressionContext = expr
					}
				} else {
					// Complex expression - save ExpressionContext
					selectItem.ExpressionContext = expr
				}

				// Alias
				if alias := item.Alias(); alias != nil {
					if id := alias.IDENTIFIER(); id != nil {
						selectItem.Alias = id.GetText()
					}
				}
				stmt.Columns = append(stmt.Columns, selectItem)
			}
		}

		// Order By
		if ctx.OrderByClause() != nil {
			for _, item := range ctx.OrderByClause().AllOrderByItem() {
				orderItem := items.OrderByItem{}

				// Извлекаем имя колонки или индекс из expression
				if expr := item.Expression(); expr != nil {
					// Пробуем извлечь простое выражение
					if primExpr := extractPrimaryExpression(expr); primExpr != nil {
						// Проверяем, это число (индекс колонки) или имя колонки
						if primExpr.Value() != nil {
							valueText := primExpr.Value().GetText()
							// Пробуем распарсить как число
							if idx, err := strconv.Atoi(valueText); err == nil && idx > 0 {
								orderItem.ColumnIndex = idx
							} else {
								// Не число - используем как имя
								orderItem.ColumnName = valueText
							}
						} else if primExpr.ColumnName() != nil {
							orderItem.ColumnName = primExpr.ColumnName().GetText()
						}
					} else {
						// Сложное выражение - берем как текст
						orderItem.ColumnName = expr.GetText()
					}
				}

				if item.ASC() != nil {
					orderItem.Direction = "ASC"
				} else if item.DESC() != nil {
					orderItem.Direction = "DESC"
				} else {
					orderItem.Direction = "ASC" // default
				}

				stmt.OrderBy = append(stmt.OrderBy, orderItem)
			}
		}

		// Limit
		if ctx.LimitValue() != nil {
			limitStr := ctx.LimitValue().GetText()
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				l.err = fmt.Errorf("invalid LIMIT value: %v", err)
				return
			}
			stmt.Limit = limit
		}

		l.stmt = stmt
	}
}

// extractPrimaryExpression tries to extract a simple primaryExpression from an expression tree
// Returns nil if the expression is complex (has operators, etc.)
func extractPrimaryExpression(expr parser.IExpressionContext) parser.IPrimaryExpressionContext {
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
