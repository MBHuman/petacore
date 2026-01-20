package rhelpers

import (
	"fmt"
	"log"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
)

// evaluateExpressionContext evaluates an expression using the ANTLR parsed context
func EvaluateExpressionContext(ctx parser.IExpressionContext) (interface{}, error) {
	parsed := ParseExpression(ctx)
	log.Printf("Evaluating expression: %v", ctx)
	log.Printf("Parsed expression result: %v", parsed)
	if caseExpr, ok := parsed.(*items.CaseExpression); ok {
		return EvaluateCaseExpression(caseExpr)
	}
	return parsed, nil
}

// evaluateCaseExpression evaluates a CASE WHEN THEN ELSE END expression
func EvaluateCaseExpression(caseExpr *items.CaseExpression) (interface{}, error) {
	ctx := caseExpr.Context
	// Get all WHEN THEN pairs
	whenThens := ctx.AllWHEN()
	expressions := ctx.AllExpression()

	// CASE WHEN expr THEN expr [WHEN expr THEN expr]* [ELSE expr] END
	// expressions[0] is first WHEN condition, expressions[1] is first THEN result, etc.
	// If ELSE, last expression is ELSE result

	numWhen := len(whenThens)
	if numWhen == 0 {
		return nil, fmt.Errorf("invalid CASE expression")
	}

	exprIndex := 0
	for i := 0; i < numWhen; i++ {
		// Evaluate WHEN condition
		conditionExpr := expressions[exprIndex]
		conditionValue := ParseExpression(conditionExpr)
		// For simplicity, treat non-nil as true
		if conditionValue != nil && conditionValue != false && conditionValue != 0 {
			// Evaluate THEN result
			resultExpr := expressions[exprIndex+1]
			return ParseExpression(resultExpr), nil
		}
		exprIndex += 2
	}

	// Check for ELSE
	if exprIndex < len(expressions) {
		elseExpr := expressions[exprIndex]
		return ParseExpression(elseExpr), nil
	}

	// No match and no ELSE
	return nil, nil
}

func EvaluateWhereCondition(where *items.WhereClause, row map[string]interface{}) bool {
	if where == nil {
		return true
	}
	fieldValue, ok := row[where.Field]
	if !ok {
		return false
	}
	// Simple comparison, assume =
	return fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", where.Value)
}
