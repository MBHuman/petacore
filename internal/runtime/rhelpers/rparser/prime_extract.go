package rparser

import "petacore/internal/runtime/parser"

// ExtractPrimaryExpression tries to extract a simple primaryExpression from an expression tree
// Returns nil if the expression is complex (has operators, etc.)
func ExtractPrimaryExpression(expr parser.IExpressionContext) parser.IPrimaryExpressionContext {
	if expr == nil {
		return nil
	}

	// expression -> orExpression
	orExpr := expr.OrExpression()
	if orExpr == nil {
		return nil
	}

	// orExpression -> andExpression (if only one)
	andExprs := orExpr.AllAndExpression()
	if len(andExprs) != 1 {
		return nil // multiple AND expressions
	}

	// andExpression -> notExpression (if only one)
	notExprs := andExprs[0].AllNotExpression()
	if len(notExprs) != 1 {
		return nil // multiple NOT expressions
	}

	notExpr := notExprs[0]
	if notExpr.NOT() != nil {
		return nil // has NOT operator
	}

	// notExpression -> comparisonExpression
	compExpr := notExpr.ComparisonExpression()
	if compExpr == nil {
		return nil
	}

	// comparisonExpression should have no operators
	if compExpr.Operator() != nil || compExpr.IN() != nil || compExpr.LIKE() != nil || compExpr.IS() != nil {
		return nil // has comparison operators
	}

	// comparisonExpression -> concatExpression (if only one)
	concatExprs := compExpr.AllConcatExpression()
	if len(concatExprs) != 1 {
		return nil
	}

	concatExpr := concatExprs[0]
	if len(concatExpr.AllCONCAT()) > 0 {
		return nil // has concatenation
	}

	// concatExpression -> additiveExpression (if only one)
	addExprs := concatExpr.AllAdditiveExpression()
	if len(addExprs) != 1 {
		return nil
	}

	addExpr := addExprs[0]
	if len(addExpr.AllPLUS()) > 0 || len(addExpr.AllMINUS()) > 0 {
		return nil // has addition/subtraction
	}

	// additiveExpression -> multiplicativeExpression (if only one)
	multExprs := addExpr.AllMultiplicativeExpression()
	if len(multExprs) != 1 {
		return nil
	}

	multExpr := multExprs[0]
	if len(multExpr.AllSTAR()) > 0 || len(multExpr.AllSLASH()) > 0 {
		return nil // has multiplication/division
	}

	// multiplicativeExpression -> unaryExpression (if only one)
	unaryExprs := multExpr.AllUnaryExpression()
	if len(unaryExprs) != 1 {
		return nil
	}

	unaryExpr := unaryExprs[0]
	if unaryExpr.PLUS() != nil || unaryExpr.MINUS() != nil {
		return nil // has unary operator
	}

	// unaryExpression -> castExpression
	castExpr := unaryExpr.CastExpression()
	if castExpr == nil {
		return nil
	}

	if len(castExpr.AllPostfix()) > 0 {
		return nil // has postfix (cast, collate, at time zone)
	}

	// castExpression -> primaryExpression
	primExpr := castExpr.PrimaryExpression()

	// atTimeZoneExpression -> primaryExpression
	return primExpr
}
