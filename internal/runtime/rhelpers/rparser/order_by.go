package rparser

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
	"strconv"
)

// parseOrderByClause парсит ORDER BY
func ParseOrderByClause(ctx parser.IOrderByClauseContext) []items.OrderByItem {
	if ctx == nil {
		return nil
	}

	var orderBy []items.OrderByItem

	for _, item := range ctx.AllOrderByItem() {
		orderItem := items.OrderByItem{}

		if expr := item.Expression(); expr != nil {
			if primExpr := ExtractPrimaryExpression(expr); primExpr != nil {
				if primExpr.Value() != nil {
					valueText := primExpr.Value().GetText()
					if idx, err := strconv.Atoi(valueText); err == nil && idx > 0 {
						orderItem.ColumnIndex = idx
					} else {
						orderItem.ColumnName = valueText
					}
				} else if primExpr.ColumnName() != nil {
					orderItem.ColumnName = primExpr.ColumnName().GetText()
				}
			} else {
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

		orderBy = append(orderBy, orderItem)
	}

	return orderBy
}
