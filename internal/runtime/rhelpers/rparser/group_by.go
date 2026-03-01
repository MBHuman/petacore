package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
)

// ParseGroupByClause парсит GROUP BY
func ParseGroupByClause(ctx parser.IGroupByClauseContext) ([]items.SelectItem, error) {
	if ctx == nil {
		return nil, nil
	}

	var groupBy []items.SelectItem

	for _, expr := range ctx.AllExpression() {
		groupByItem := items.SelectItem{}
		if primExpr := ExtractPrimaryExpression(expr); primExpr != nil {
			if primExpr.ColumnName() != nil {
				fullName := primExpr.ColumnName().GetText()
				if qn := primExpr.ColumnName().QualifiedName(); qn != nil {
					parts := qn.AllNamePart()
					if len(parts) == 2 {
						groupByItem.TableAlias = parts[0].GetText()
						groupByItem.ColumnName = parts[1].GetText()
					} else if len(parts) == 1 {
						groupByItem.ColumnName = parts[0].GetText()
					} else {
						groupByItem.ColumnName = fullName
					}
				} else {
					groupByItem.ColumnName = fullName
				}
			}
		} else {
			return nil, fmt.Errorf("complex GROUP BY expressions are not supported yet")
		}
		groupBy = append(groupBy, groupByItem)
	}

	return groupBy, nil
}
