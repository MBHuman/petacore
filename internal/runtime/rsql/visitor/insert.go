package visitor

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rparser"
	"petacore/internal/runtime/rsql/statements"
	"petacore/sdk/serializers"
	"strings"
)

func (l *sqlListener) EnterInsertStatement(ctx *parser.InsertStatementContext) {
	stmt := &statements.InsertStatement{}
	if ctx.TableName() != nil {
		stmt.TableName = ctx.TableName().GetText()
	}
	goCtx := context.Background()

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
			value, err := rparser.ParseExpression(l.allocator, goCtx, expr, nil, l.subExec)
			if err != nil {
				l.err = err
				return
			}
			if value == nil {
				l.err = fmt.Errorf("[InsertStatement] invalid expression in INSERT")
				return
			}
			if val, ok := value.(*rmodels.ResultRowsExpression); ok {
				if len(val.Row.Rows) > 0 && len(val.Row.Schema.Fields) > 0 {
					buf, oid, err := val.Row.Schema.GetField(val.Row.Rows[0], 0)
					if err != nil {
						l.err = fmt.Errorf("[InsertStatement] failed to get field: %w", err)
						return
					}
					desVal, err := serializers.DeserializeGeneric(buf, oid)
					if err != nil {
						l.err = fmt.Errorf("[InsertStatement] failed to deserialize: %w", err)
						return
					}
					rowValues = append(rowValues, desVal)
				} else {
					rowValues = append(rowValues, nil)
				}
			} else if val, ok := value.(*rmodels.BoolExpression); ok {
				rowValues = append(rowValues, val.Value)
			}
		}
		stmt.Values = append(stmt.Values, rowValues)
	}

	l.stmt = stmt
}
