package visitor

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rparser"
)

// EnterSelectStatement - входная точка для парсинга SELECT
func (l *sqlListener) EnterSelectStatement(ctx *parser.SelectStatementContext) {
	// Если после парсинга stmt не создан, явно записываем ошибку
	stmt, err := rparser.ParseSelectStatement(ctx)
	if err != nil {
		l.err = fmt.Errorf("error parsing SELECT statement: %v", err)
		return
	}
	// Only set top-level statement once; don't overwrite if nested SELECTs are parsed later
	if stmt != nil && l.stmt == nil {
		l.stmt = stmt
	}
	if stmt == nil && l.err == nil {
		l.err = fmt.Errorf("SELECT: unsupported or invalid syntax")
	}
}
