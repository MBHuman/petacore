package visitor

import (
	"errors"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rparser"
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

	primaryKeys := make([]int, 0)
	primaryKeysMap := make(map[int]struct{})
	// Columns
	for idx, colDef := range ctx.AllColumnDefinition() {
		col := table.ColumnDef{}
		if colDef.ColumnName() != nil {
			col.Name = colDef.ColumnName().GetText()
		}
		if colDef.DataType() != nil {
			dataTypeText := strings.ToUpper(colDef.DataType().GetText())
			col.Type = rparser.ParseDataType(dataTypeText)
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
			primaryKeys = append(primaryKeys, idx+1)
			primaryKeysMap[idx+1] = struct{}{}
		}

		// PRIMARY KEY implies NOT NULL
		if _, ok := primaryKeysMap[idx+1]; ok {
			col.IsNullable = false
		}

		stmt.Columns = append(stmt.Columns, col)
	}

	allColumnNames := ctx.AllColumnName()
	if len(primaryKeys) > 0 && len(allColumnNames) > 0 {
		l.err = errors.New("PRIMARY KEY already defined in column definitions")
		return // PRIMARY KEY already defined in column definitions
	}

	for _, colCtx := range allColumnNames {
		colName := colCtx.GetText()
		// Find column index
		for i, col := range stmt.Columns {
			if col.Name == colName {
				pkIdx := i + 1
				// Avoid duplicates
				if _, exists := primaryKeysMap[pkIdx]; !exists {
					primaryKeys = append(primaryKeys, pkIdx)
					primaryKeysMap[pkIdx] = struct{}{}
					// PRIMARY KEY implies NOT NULL
					stmt.Columns[i].IsNullable = false
				}
				break
			}
		}
	}

	stmt.PrimaryKeys = primaryKeys

	l.stmt = stmt
}
