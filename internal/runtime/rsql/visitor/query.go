package visitor

import "petacore/internal/runtime/parser"

// TODO пересмотреть
func (l *sqlListener) EnterQuery(ctx *parser.QueryContext) {
	// Инициализируем в зависимости от типа
}
