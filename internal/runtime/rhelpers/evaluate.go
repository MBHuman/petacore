package rhelpers

import (
	"fmt"
	"log"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
)

// evaluateExpressionContext evaluates an expression using the ANTLR parsed context
func EvaluateExpressionContext(ctx parser.IExpressionContext, row map[string]interface{}) (interface{}, error) {
	parsed := ParseExpression(ctx, row)
	// log.Printf("Evaluating expression: %v", ctx)
	// log.Printf("Parsed expression result: %v", parsed)
	if caseExpr, ok := parsed.(*items.CaseExpression); ok {
		return EvaluateCaseExpression(caseExpr, row)
	}
	return parsed, nil
}

func EvaluateComparisonExpressionContext(ctx parser.IComparisonExpressionContext, row map[string]interface{}) (interface{}, error) {
	parsed := parseComparisonExpression(ctx, row)
	// log.Printf("Evaluating comparison expression: %v", ctx.GetText())
	// log.Printf("Parsed comparison expression result: %v", parsed)
	return parsed, nil
}

// evaluateCaseExpression evaluates a CASE WHEN THEN ELSE END expression
func EvaluateCaseExpression(caseExpr *items.CaseExpression, row map[string]interface{}) (interface{}, error) {
	ctx := caseExpr.Context
	// In new grammar: CASE (WHEN expression THEN expression)+ (ELSE expression)? END
	// AllExpression returns all expressions: WHEN1, THEN1, WHEN2, THEN2, ..., ELSE (if present)
	allExpressions := ctx.AllExpression()
	numWhen := len(ctx.AllWHEN())

	if numWhen == 0 {
		return nil, fmt.Errorf("invalid CASE expression")
	}

	// Expressions alternate: WHEN1, THEN1, WHEN2, THEN2, ...
	for i := 0; i < numWhen; i++ {
		whenIdx := i * 2
		thenIdx := i*2 + 1

		if whenIdx >= len(allExpressions) {
			break
		}

		// Evaluate WHEN condition
		condition, err := EvaluateExpressionContext(allExpressions[whenIdx], row)
		if err != nil {
			return nil, err
		}
		log.Printf("CASE condition: %v, value: %v", allExpressions[whenIdx].GetText(), condition)

		// Check if true
		if IsTrue(condition) {
			if thenIdx < len(allExpressions) {
				// Evaluate THEN result
				result, _ := EvaluateExpressionContext(allExpressions[thenIdx], row)
				log.Printf("CASE THEN result: %v", result)
				return result, nil
			}
		}
	}

	// Check for ELSE
	elseIdx := numWhen * 2
	if elseIdx < len(allExpressions) {
		result, _ := EvaluateExpressionContext(allExpressions[elseIdx], row)
		log.Printf("CASE ELSE result: %v", result)
		return result, nil
	}

	// No match and no ELSE
	return nil, nil
}

func IsTrue(value interface{}) bool {
	if value == nil {
		return false
	}
	if boolVal, ok := value.(bool); ok {
		return boolVal
	}
	// For other types, treat as true if not zero/false
	if intVal, ok := value.(int); ok {
		return intVal != 0
	}
	if floatVal, ok := value.(float64); ok {
		return floatVal != 0
	}
	if strVal, ok := value.(string); ok {
		return strVal != ""
	}
	return true
}

func EvaluateWhereCondition(where *items.WhereClause, row map[string]interface{}) bool {
	if where == nil {
		return true
	}
	result, err := EvaluateExpressionContext(where.ExpressionContext, row)
	if err != nil {
		log.Printf("DEBUG EvaluateWhereCondition: error evaluating expression: %v", err)
		return false
	}
	log.Printf("DEBUG EvaluateWhereCondition: expression result = %v (type %T)", result, result)
	// Convert to bool
	if boolVal, ok := result.(bool); ok {
		return boolVal
	}
	log.Printf("DEBUG EvaluateWhereCondition: result is not bool, returning false")
	return false
}
