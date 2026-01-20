package visitor

import "petacore/internal/runtime/parser"

func (l *sqlListener) EnterQuery(ctx *parser.QueryContext) {
	// Инициализируем в зависимости от типа
}
