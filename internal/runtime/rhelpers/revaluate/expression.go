package revaluate

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rparser"
	"petacore/internal/runtime/rsql/table"
)

// evaluateExpressionContext evaluates an expression using the ANTLR parsed context
// Может возвращать bool для условий и *table.ExecuteResult для других выражений
func EvaluateExpressionContext(ctx parser.IExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	parsed, err := rparser.ParseExpression(ctx, row)
	if err != nil {
		return nil, err
	}
	if caseExpr, ok := parsed.(*rmodels.CaseExpression); ok {
		return EvaluateCaseExpression(caseExpr, row)
	}
	return parsed, nil
}
