package visitor

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rparser"
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
			value, err := rparser.ParseExpression(expr, nil)
			if err != nil {
				l.err = err
				return
			}
			if value == nil {
				l.err = fmt.Errorf("invalid expression in INSERT")
				return
			}
			if val, ok := value.(*rmodels.ResultRowsExpression); ok {
				rowValues = append(rowValues, val.Row.Rows[0][0])
			} else if val, ok := value.(*rmodels.BoolExpression); ok {
				rowValues = append(rowValues, val.Value)
			}
		}
		stmt.Values = append(stmt.Values, rowValues)
	}

	l.stmt = stmt
}
