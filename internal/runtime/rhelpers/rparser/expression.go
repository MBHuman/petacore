package rparser

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
)

// parseExpression evaluates an ANTLR expression context and returns the value
func ParseExpression(expr parser.IExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	// logger.Debug("Parsing expression")
	if expr == nil {
		return nil, nil
	}

	// New grammar: expression -> orExpression -> andExpression -> ... -> primaryExpression
	// Just parse the orExpression which will handle the entire tree
	if orExpr := expr.OrExpression(); orExpr != nil {
		return ParseOrExpression(orExpr, row)
	}

	return nil, nil
}
