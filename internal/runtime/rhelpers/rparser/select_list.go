package rparser

import (
	"petacore/internal/logger"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
)

// parseSelectList парсит список колонок в SELECT
func ParseSelectList(ctx parser.ISelectListContext) []items.SelectItem {
	if ctx == nil {
		return nil
	}

	var columns []items.SelectItem

	for _, item := range ctx.AllSelectItem() {
		logger.Debugf("Parsing select item: %s", item.GetText())
		selectItem := items.SelectItem{}

		if selectAll := item.SelectAll(); selectAll != nil {
			selectItem.IsSelectAll = true
			columns = append(columns, selectItem)
		} else if expr := item.Expression(); expr != nil {
			// Try to extract simple column or function names for optimization
			if primExpr := ExtractPrimaryExpression(expr); primExpr != nil {
				if primExpr.ColumnName() != nil {
					fullName := primExpr.ColumnName().GetText()
					// Parse qualified name like "c.name" or "customers.name"
					if qn := primExpr.ColumnName().QualifiedName(); qn != nil {
						parts := qn.AllNamePart()
						if len(parts) == 2 {
							selectItem.TableAlias = parts[0].GetText()
							selectItem.ColumnName = parts[1].GetText()
						} else if len(parts) == 1 {
							selectItem.ColumnName = parts[0].GetText()
						} else {
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

					// Check if aggregate function через SDK registry
					funcCall.IsAggregate = functions.IsAggregateFunction(funcCall.Name)

					// Args - expression contexts
					for _, argExpr := range fc.AllFunctionArg() {
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
			columns = append(columns, selectItem)
		}
	}

	return columns
}
