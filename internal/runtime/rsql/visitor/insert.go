package visitor

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"strings"
)

func (l *sqlListener) EnterInsertStatement(ctx *parser.InsertStatementContext) {
	stmt := &statements.InsertStatement{}
	if ctx.TableName() != nil {
		stmt.TableName = ctx.TableName().GetText()
	}

	// Columns
	if ctx.ColumnList() != nil {
		columnsText := ctx.ColumnList().GetText()
		if columnsText != "" {
			// Remove parentheses and split
			columnsText = strings.Trim(columnsText, "()")
			stmt.Columns = strings.Split(columnsText, ",")
			for i, col := range stmt.Columns {
				stmt.Columns[i] = strings.TrimSpace(col)
			}
		}
	}

	// Values - multiple value lists
	for _, vl := range ctx.AllValueList() {
		var rowValues []interface{}
		for _, expr := range vl.AllExpression() {
			value := rhelpers.ParseExpression(expr, nil)
			if value == nil {
				l.err = fmt.Errorf("invalid expression in INSERT")
				return
			}
			rowValues = append(rowValues, value)
		}
		stmt.Values = append(stmt.Values, rowValues)
	}

	l.stmt = stmt
}
