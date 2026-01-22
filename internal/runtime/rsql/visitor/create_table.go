package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// TODO убрать хардкодинг типов данных и ограничений
func (l *sqlListener) EnterCreateTableStatement(ctx *parser.CreateTableStatementContext) {
	stmt := &statements.CreateTableStatement{}
	if ctx.TableName() != nil {
		stmt.TableName = ctx.TableName().GetText()
	}
	// Check for IF NOT EXISTS
	if ctx.IF() != nil && ctx.NOT() != nil && ctx.EXISTS() != nil {
		stmt.IfNotExists = true
	}

	// Columns
	for _, colDef := range ctx.AllColumnDefinition() {
		col := table.ColumnDef{}
		if colDef.ColumnName() != nil {
			col.Name = colDef.ColumnName().GetText()
		}
		if colDef.DataType() != nil {
			dataTypeText := strings.ToUpper(colDef.DataType().GetText())
			col.Type = rhelpers.ParseDataType(dataTypeText)
			if strings.Contains(dataTypeText, "SERIAL") {
				col.IsSerial = true
				col.Type = table.ColTypeInt // SERIAL is integer
			}
		}

		// Parse constraints
		col.IsNullable = true // Default
		if colDef.ColumnConstraints() != nil {
			children := colDef.ColumnConstraints().GetChildren()
			i := 0
			for i < len(children) {
				child := children[i]
				var text string
				if pt, ok := child.(antlr.ParseTree); ok {
					text = strings.ToUpper(pt.GetText())
				} else {
					i++
					continue
				}
				if text == "NOT" && i+1 < len(children) {
					if pt2, ok := children[i+1].(antlr.ParseTree); ok && strings.ToUpper(pt2.GetText()) == "NULL" {
						col.IsNullable = false
						i += 2
					} else {
						i++
					}
				} else if text == "UNIQUE" {
					col.IsUnique = true
					i++
				} else if text == "DEFAULT" && i+1 < len(children) {
					if pt2, ok := children[i+1].(antlr.ParseTree); ok {
						defaultValue := pt2.GetText()
						if defaultValue == "CURRENT_TIMESTAMP" {
							col.DefaultValue = "CURRENT_TIMESTAMP"
						} else if strings.HasPrefix(defaultValue, "'") && strings.HasSuffix(defaultValue, "'") {
							col.DefaultValue = defaultValue[1 : len(defaultValue)-1]
						} else if numVal, err := strconv.Atoi(defaultValue); err == nil {
							col.DefaultValue = numVal
						} else {
							col.DefaultValue = defaultValue
						}
					}
					i += 2
				} else {
					i++
				}
			}
		}

		// Special handling for PRIMARY KEY column
		if colDef.PRIMARY() != nil && colDef.KEY() != nil {
			col.IsPrimaryKey = true
		}

		// PRIMARY KEY implies NOT NULL
		if col.IsPrimaryKey {
			col.IsNullable = false
		}

		stmt.Columns = append(stmt.Columns, col)
	}

	l.stmt = stmt
}
