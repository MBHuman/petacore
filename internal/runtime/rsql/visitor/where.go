package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/statements"
	"strconv"
	"strings"
)

func (l *sqlListener) EnterWhereClause(ctx *parser.WhereClauseContext) {
	if ctx.Expression() != nil {
		expr := ctx.Expression()
		// Check if it's a simple comparison: expr operator expr
		if expr.AdditiveExpression(0) != nil && len(expr.AllOperator()) == 1 && len(expr.AllAdditiveExpression()) == 1 {
			leftExpr := expr.AdditiveExpression(0)
			operator := expr.AllOperator()[0].GetText()
			rightExpr := expr.AllAdditiveExpression()[0]

			// Check if left is column name and right is value
			if leftExpr.MultiplicativeExpression(0) != nil &&
				leftExpr.MultiplicativeExpression(0).PrimaryExpression(0) != nil &&
				leftExpr.MultiplicativeExpression(0).PrimaryExpression(0).ColumnName() != nil &&
				rightExpr.MultiplicativeExpression(0) != nil &&
				rightExpr.MultiplicativeExpression(0).PrimaryExpression(0) != nil &&
				rightExpr.MultiplicativeExpression(0).PrimaryExpression(0).Value() != nil {

				field := leftExpr.MultiplicativeExpression(0).PrimaryExpression(0).ColumnName().GetText()
				valueCtx := rightExpr.MultiplicativeExpression(0).PrimaryExpression(0).Value()
				valueStr := valueCtx.GetText()

				var value interface{}
				if strings.HasPrefix(valueStr, "'") && strings.HasSuffix(valueStr, "'") {
					value = valueStr[1 : len(valueStr)-1]
				} else if numVal, err := strconv.Atoi(valueStr); err == nil {
					value = numVal
				} else if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
					value = floatVal
				} else {
					value = valueStr // fallback
				}

				where := &items.WhereClause{
					Field:    field,
					Operator: operator,
					Value:    value,
				}

				if selectStmt, ok := l.stmt.(*statements.SelectStatement); ok {
					selectStmt.Where = where
				}
			}
		}
	}
}
