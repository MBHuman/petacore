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

	// Table name
	if ctx.FromClause() != nil && ctx.FromClause().TableName() != nil {
		stmt.TableName = ctx.FromClause().TableName().GetText()
	}

	// Columns
	if ctx.ColumnList() != nil {
		if ctx.ColumnList().STAR() != nil {
			// SELECT *
			stmt.Columns = []items.SelectItem{{ColumnName: "*"}}
		} else {
			// Parse select items
			for _, item := range ctx.ColumnList().AllSelectItem() {
				selectItem := items.SelectItem{}
				expr := item.Expression()
				if expr != nil {
					// Check if it's a simple expression
					if expr.AdditiveExpression(0) != nil && len(expr.AllOperator()) == 0 {
						// Simple expression, check the primary
						addExpr := expr.AdditiveExpression(0)
						if addExpr.MultiplicativeExpression(0) != nil && len(addExpr.AllPLUS()) == 0 && len(addExpr.AllMINUS()) == 0 {
							multExpr := addExpr.MultiplicativeExpression(0)
							if multExpr.PrimaryExpression(0) != nil && len(multExpr.AllSTAR()) == 0 && len(multExpr.AllSLASH()) == 0 {
								primExpr := multExpr.PrimaryExpression(0)
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
										value := rhelpers.ParseExpression(argExpr)
										funcCall.Args = append(funcCall.Args, value)
									}
									selectItem.Function = funcCall
								} else {
									// Other primary expressions
									selectItem.ExpressionContext = expr
								}
							} else {
								// Complex multiplicative expression
								selectItem.ExpressionContext = expr
							}
						} else {
							// Complex additive expression
							selectItem.ExpressionContext = expr
						}
					} else {
						// Complex expression with operators
						selectItem.ExpressionContext = expr
					}
				}
				// Alias
				if item.IDENTIFIER() != nil {
					selectItem.Alias = item.IDENTIFIER().GetText()
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
