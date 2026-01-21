package visitor

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"strconv"
)

func (l *sqlListener) EnterSelectStatement(ctx *parser.SelectStatementContext) {
	stmt := &statements.SelectStatement{}

	// From clause
	if ctx.FromClause() != nil {
		from := &statements.FromClause{}
		if ctx.FromClause().TableName() != nil {
			from.TableName = ctx.FromClause().TableName().GetText()
			// Table alias
			if alias := ctx.FromClause().Alias(); alias != nil {
				if id := alias.IDENTIFIER(); id != nil {
					from.Alias = id.GetText()
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
			if jc.TableName() != nil {
				join.TableName = jc.TableName().GetText()
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
		// For backward compatibility
		if from.TableName != "" {
			stmt.TableName = from.TableName
			stmt.TableAlias = from.Alias
		}
	}

	// Columns
	if ctx.SelectList() != nil {
		if ctx.SelectList().STAR() != nil {
			// SELECT *
			stmt.Columns = []items.SelectItem{{ColumnName: "*"}}
		} else {
			// Parse select items
			for _, item := range ctx.SelectList().AllSelectItem() {
				selectItem := items.SelectItem{}
				expr := item.Expression()
				if expr != nil {
					// Try to extract simple column or function names for optimization
					if primExpr := extractPrimaryExpression(expr); primExpr != nil {
						if primExpr.ColumnName() != nil {
							selectItem.ColumnName = primExpr.ColumnName().GetText()
						} else if primExpr.FunctionCall() != nil {
							fc := primExpr.FunctionCall()
							funcCall := &items.FunctionCall{}
							if qn := fc.QualifiedName(); qn != nil {
								ids := qn.AllIDENTIFIER()
								parts := []string{}
								for _, id := range ids {
									parts = append(parts, id.GetText())
								}
								funcCall.Name = parts[len(parts)-1]
							}

							// Args - now expressions
							for _, argExpr := range fc.AllExpression() {
								value := rhelpers.ParseExpression(argExpr, nil)
								funcCall.Args = append(funcCall.Args, value)
							}
							selectItem.Function = funcCall
						} else {
							// Complex expression
							selectItem.ExpressionContext = expr
						}
					} else {
						// Complex expression
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
		}

		// Order By
		if ctx.OrderByClause() != nil {
			for _, item := range ctx.OrderByClause().AllOrderByItem() {
				orderItem := items.OrderByItem{}
				if item.Expression() != nil {
					orderItem.ExpressionContext = item.Expression()
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

	// multiplicativeExpression -> castExpression (if only one)
	castExprs := multExpr.AllCastExpression()
	if len(castExprs) != 1 {
		return nil
	}

	castExpr := castExprs[0]
	if castExpr.COLON_COLON() != nil {
		return nil // has cast operator
	}

	// castExpression -> atTimeZoneExpression
	atExpr := castExpr.AtTimeZoneExpression()
	if atExpr == nil {
		return nil
	}

	if atExpr.AT() != nil {
		return nil // has AT TIME ZONE
	}

	// atTimeZoneExpression -> primaryExpression
	return atExpr.PrimaryExpression()
}
