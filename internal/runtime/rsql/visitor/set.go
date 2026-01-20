package visitor

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
	"strconv"
	"strings"
)

func (l *sqlListener) EnterSetStatement(ctx *parser.SetStatementContext) {
	stmt := &statements.SetStatement{}
	if ctx.IDENTIFIER() != nil {
		stmt.Variable = ctx.IDENTIFIER().GetText()
	}
	if ctx.Value() != nil {
		valueStr := ctx.Value().GetText()
		var value interface{}
		if strings.HasPrefix(valueStr, "'") && strings.HasSuffix(valueStr, "'") {
			value = valueStr[1 : len(valueStr)-1]
		} else if numVal, err := strconv.Atoi(valueStr); err == nil {
			value = numVal
		} else if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
			value = floatVal
		} else {
			value = valueStr
		}
		stmt.Value = value
	}
	l.stmt = stmt
}
