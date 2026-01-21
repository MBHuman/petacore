package rhelpers

import (
	"fmt"
	"log"
	"sort"

	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/items"
	"petacore/internal/runtime/rsql/table"
	"strconv"
	"strings"
	"time"
)

func ParseDataType(typeStr string) table.ColType {
	switch strings.ToUpper(typeStr) {
	case "STRING", "TEXT":
		return table.ColTypeString
	case "INT":
		return table.ColTypeInt
	case "FLOAT":
		return table.ColTypeFloat
	case "BOOL":
		return table.ColTypeBool
	default:
		return table.ColTypeString
	}
}

// parseExpression evaluates an ANTLR expression context and returns the value
func ParseExpression(expr parser.IExpressionContext, row map[string]interface{}) interface{} {
	if expr == nil {
		return nil
	}

	// New grammar: expression -> orExpression -> andExpression -> ... -> primaryExpression
	// Just parse the orExpression which will handle the entire tree
	if orExpr := expr.OrExpression(); orExpr != nil {
		return parseOrExpression(orExpr, row)
	}

	return nil
}

// parseAdditiveExpression handles addition and subtraction
func ParseAdditiveExpression(addExpr parser.IAdditiveExpressionContext, row map[string]interface{}) interface{} {
	if addExpr == nil {
		return nil
	}

	// Get the first multiplicative expression
	multExpr := addExpr.MultiplicativeExpression(0)
	if multExpr == nil {
		return nil
	}

	result := parseMultiplicativeExpression(multExpr, row)

	// Handle additional terms with operators
	plusOps := addExpr.AllPLUS()
	minusOps := addExpr.AllMINUS()

	// Merge operators by token index
	type opInfo struct {
		op    string
		index int
	}
	var ops []opInfo
	for _, p := range plusOps {
		ops = append(ops, opInfo{"+", p.GetSymbol().GetTokenIndex()})
	}
	for _, m := range minusOps {
		ops = append(ops, opInfo{"-", m.GetSymbol().GetTokenIndex()})
	}
	// Sort by token index
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].index < ops[j].index
	})

	multExprs := addExpr.AllMultiplicativeExpression()
	// Skip the first one since we already processed it
	for i, op := range ops {
		if i+1 >= len(multExprs) {
			break
		}
		nextValue := parseMultiplicativeExpression(multExprs[i+1], row)

		if op.op == "+" {
			result = addValues(result, nextValue)
		} else if op.op == "-" {
			result = subtractValues(result, nextValue)
		}
	}

	return result
}

// parseCaseExpression handles CASE WHEN THEN ELSE END expressions
func ParseCaseExpression(caseExpr parser.ICaseExpressionContext) interface{} {
	if caseExpr == nil {
		return nil
	}

	// For now, return a placeholder - actual evaluation will be done during execution
	// We need to store the case expression for later evaluation
	return &items.CaseExpression{
		Context: caseExpr,
	}
}

// parseConcatExpression handles string concatenation with ||
func parseConcatExpression(concatExpr parser.IConcatExpressionContext, row map[string]interface{}) interface{} {
	if concatExpr == nil {
		return nil
	}

	// Get all additive expressions to concatenate
	exprs := concatExpr.AllAdditiveExpression()
	if len(exprs) == 0 {
		return nil
	}

	// If only one expression and no CONCAT operators, return it directly
	if len(exprs) == 1 && len(concatExpr.AllCONCAT()) == 0 {
		return ParseAdditiveExpression(exprs[0], row)
	}

	// Multiple expressions with CONCAT, concatenate as strings
	result := ""
	for _, e := range exprs {
		val := ParseAdditiveExpression(e, row)
		if val != nil {
			result += fmt.Sprintf("%v", val)
		}
	}
	return result
}

// parseExtractFunction handles EXTRACT expressions
func ParseExtractFunction(extractExpr parser.IExtractFunctionContext, row map[string]interface{}) interface{} {
	if extractExpr == nil {
		return nil
	}

	field := extractExpr.IDENTIFIER().GetText()
	sourceExpr := extractExpr.Expression()
	if sourceExpr == nil {
		return nil
	}

	source := ParseExpression(sourceExpr, row)
	args := []interface{}{field, source}
	value, _ := functions.ExecuteFunction("EXTRACT", args)
	return value
}

// // parseConcatExpression handles concatenation
// func parseConcatExpression(concatExpr parser.IConcatExpressionContext, row map[string]interface{}) interface{} {
// 	if concatExpr == nil {
// 		return nil
// 	}

// 	result := ParseAdditiveExpression(concatExpr.AdditiveExpression(0), row)

// 	// Handle additional concatenations
// 	addExprs := concatExpr.AllAdditiveExpression()
// 	// Skip the first one
// 	for i := 1; i < len(addExprs); i++ {
// 		next := ParseAdditiveExpression(addExprs[i], row)
// 		result = fmt.Sprintf("%v%v", result, next) // simple concatenation
// 	}

// 	return result
// }

// parseOrExpression handles OR expressions
func parseOrExpression(orExpr parser.IOrExpressionContext, row map[string]interface{}) interface{} {
	if orExpr == nil {
		return nil
	}

	andExprs := orExpr.AllAndExpression()
	if len(andExprs) == 0 {
		return nil
	}

	// Evaluate first AND expression
	result := parseAndExpression(andExprs[0], row)

	// If multiple AND expressions connected by OR
	for i := 1; i < len(andExprs); i++ {
		rightVal := parseAndExpression(andExprs[i], row)
		// OR logic: true if either is true
		if boolLeft, ok := result.(bool); ok {
			if boolRight, ok := rightVal.(bool); ok {
				result = boolLeft || boolRight
				continue
			}
		}
		// If not bools, treat as truthy values
		if isTrueHelper(result) || isTrueHelper(rightVal) {
			result = true
		} else {
			result = false
		}
	}

	return result
}

// parseAndExpression handles AND expressions
func parseAndExpression(andExpr parser.IAndExpressionContext, row map[string]interface{}) interface{} {
	if andExpr == nil {
		return nil
	}

	notExprs := andExpr.AllNotExpression()
	if len(notExprs) == 0 {
		return nil
	}

	// Evaluate first NOT expression
	result := parseNotExpression(notExprs[0], row)

	// If multiple NOT expressions connected by AND
	for i := 1; i < len(notExprs); i++ {
		rightVal := parseNotExpression(notExprs[i], row)
		// AND logic: true only if both are true
		if boolLeft, ok := result.(bool); ok {
			if boolRight, ok := rightVal.(bool); ok {
				result = boolLeft && boolRight
				continue
			}
		}
		// If not bools, treat as truthy values
		if isTrueHelper(result) && isTrueHelper(rightVal) {
			result = true
		} else {
			result = false
		}
	}

	return result
}

// parseNotExpression handles NOT expression
func parseNotExpression(notExpr parser.INotExpressionContext, row map[string]interface{}) interface{} {
	if notExpr == nil {
		return nil
	}

	compExpr := notExpr.ComparisonExpression()
	if compExpr == nil {
		return nil
	}

	result := parseComparisonExpression(compExpr, row)

	// Apply NOT if present
	if notExpr.NOT() != nil {
		if boolVal, ok := result.(bool); ok {
			return !boolVal
		}
		// If not a bool, treat as truthy value
		return !isTrueHelper(result)
	}

	return result
}

func isTrueHelper(value interface{}) bool {
	if value == nil {
		return false
	}
	if boolVal, ok := value.(bool); ok {
		return boolVal
	}
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

// parseComparisonExpression handles comparison expressions including IN, LIKE, IS NULL
func parseComparisonExpression(compExpr parser.IComparisonExpressionContext, row map[string]interface{}) interface{} {
	if compExpr == nil {
		return nil
	}

	// Get the left side (concatExpression)
	concatExprs := compExpr.AllConcatExpression()
	if len(concatExprs) == 0 {
		return nil
	}

	left := parseConcatExpression(concatExprs[0], row)

	// Check for operator
	if compExpr.Operator() != nil {
		if len(concatExprs) < 2 {
			return left
		}
		op := compExpr.Operator().GetText()
		right := parseConcatExpression(concatExprs[1], row)
		return evaluateComparison(left, right, op)
	}

	// Check for IN
	if compExpr.IN() != nil {
		not := compExpr.NOT() != nil
		var values []interface{}
		for _, expr := range compExpr.AllExpression() {
			values = append(values, ParseExpression(expr, row))
		}
		in := contains(values, left)
		if not {
			in = !in
		}
		return in
	}

	// Check for LIKE
	if compExpr.LIKE() != nil {
		if len(concatExprs) < 2 {
			return left
		}
		not := compExpr.NOT() != nil
		right := parseConcatExpression(concatExprs[1], row)
		// Simple LIKE implementation
		leftStr := fmt.Sprintf("%v", left)
		rightStr := fmt.Sprintf("%v", right)
		match := strings.Contains(leftStr, strings.ReplaceAll(strings.ReplaceAll(rightStr, "%", ""), "_", ""))
		if not {
			return !match
		}
		return match
	}

	// Check for IS NULL
	if compExpr.IS() != nil {
		not := compExpr.NOT() != nil
		isNull := left == nil
		if not {
			return !isNull
		}
		return isNull
	}

	// No operator, just return the left value
	return left
}

// Old parseComparisonExpression - keeping for reference
func parseComparisonExpressionOLD(compExpr parser.IComparisonExpressionContext, row map[string]interface{}) interface{} {
	if compExpr == nil {
		return nil
	}

	left := parseConcatExpression(compExpr.ConcatExpression(0), row)

	if compExpr.Operator() != nil {
		// Regular comparison
		op := compExpr.Operator().GetText()
		right := parseConcatExpression(compExpr.ConcatExpression(1), row)
		return evaluateComparison(left, right, op)
	}

	if compExpr.IN() != nil {
		// IN expression
		not := compExpr.NOT() != nil
		var values []interface{}
		for _, expr := range compExpr.AllExpression() {
			values = append(values, ParseExpression(expr, row))
		}
		in := contains(values, left)
		if not {
			in = !in
		}
		return in
	}

	// No operator or IN, just the left value
	return left
}

// evaluateComparison evaluates a comparison operator
func evaluateComparison(left, right interface{}, op string) bool {
	log.Printf("DEBUG evaluateComparison: left=%v (type %T), right=%v (type %T), op=%s\n", left, left, right, right, op)
	switch op {
	case "=":
		result := compareEquals(left, right)
		log.Printf("DEBUG evaluateComparison: result=%v\n", result)
		return result
	case "!=":
		return left != right
	case "<":
		if la, ok := left.(int); ok {
			if ra, ok := right.(int); ok {
				return la < ra
			}
		}
		if la, ok := left.(float64); ok {
			if ra, ok := right.(float64); ok {
				return la < ra
			}
		}
		if la, ok := left.(string); ok {
			if ra, ok := right.(string); ok {
				return la < ra
			}
		}
	case ">":
		if la, ok := left.(int); ok {
			if ra, ok := right.(int); ok {
				return la > ra
			}
		}
		if la, ok := left.(float64); ok {
			if ra, ok := right.(float64); ok {
				return la > ra
			}
		}
		if la, ok := left.(string); ok {
			if ra, ok := right.(string); ok {
				return la > ra
			}
		}
	case "<=":
		if la, ok := left.(int); ok {
			if ra, ok := right.(int); ok {
				return la <= ra
			}
		}
		if la, ok := left.(float64); ok {
			if ra, ok := right.(float64); ok {
				return la <= ra
			}
		}
		if la, ok := left.(string); ok {
			if ra, ok := right.(string); ok {
				return la <= ra
			}
		}
	case ">=":
		if la, ok := left.(int); ok {
			if ra, ok := right.(int); ok {
				return la >= ra
			}
		}
		if la, ok := left.(float64); ok {
			if ra, ok := right.(float64); ok {
				return la >= ra
			}
		}
		if la, ok := left.(string); ok {
			if ra, ok := right.(string); ok {
				return la >= ra
			}
		}
	case "LIKE":
		if la, ok := left.(string); ok {
			if ra, ok := right.(string); ok {
				return strings.Contains(la, ra) // simple like
			}
		}
	}
	return false
}

// compareEquals properly compares two values for equality, handling type conversions
func compareEquals(left, right interface{}) bool {
	// Handle nil cases
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}

	// Try direct comparison first (same types)
	if left == right {
		return true
	}

	// Handle numeric comparisons with type conversion
	leftNum, leftIsNum := toNumber(left)
	rightNum, rightIsNum := toNumber(right)

	if leftIsNum && rightIsNum {
		return leftNum == rightNum
	}

	// String comparison
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)
	return leftStr == rightStr
}

// toNumber converts a value to float64 if it's numeric
func toNumber(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case string:
		// Try to parse string as number
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, true
		}
		if num, err := strconv.Atoi(v); err == nil {
			return float64(num), true
		}
	}
	return 0, false
}

// contains checks if value is in the list
func contains(list []interface{}, value interface{}) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// parsePrimaryExpression handles the basic expressions
func parsePrimaryExpression(primExpr parser.IPrimaryExpressionContext, row map[string]interface{}) interface{} {
	if primExpr == nil {
		return nil
	}

	// Check for parenthesized expression
	if primExpr.Expression() != nil {
		return ParseExpression(primExpr.Expression(), row)
	}

	// Check for CASE expression
	if primExpr.CaseExpression() != nil {
		return ParseCaseExpression(primExpr.CaseExpression())
	}

	// Check for function call
	if primExpr.FunctionCall() != nil {
		fc := primExpr.FunctionCall()
		funcName := ""
		if qn := fc.QualifiedName(); qn != nil {
			ids := qn.AllIDENTIFIER()
			if len(ids) > 0 {
				parts := []string{}
				for _, id := range ids {
					parts = append(parts, id.GetText())
				}
				funcName = parts[len(parts)-1]
			}
		}
		var args []interface{}
		for _, argExpr := range fc.AllExpression() {
			args = append(args, ParseExpression(argExpr, row))
		}
		value, _ := functions.ExecuteFunction(funcName, args)
		return value
	}

	// Check for extract function
	if primExpr.ExtractFunction() != nil {
		return ParseExtractFunction(primExpr.ExtractFunction(), row)
	}

	// Check for column name
	if primExpr.ColumnName() != nil {
		columnName := primExpr.ColumnName().GetText()
		// Handle qualified names like table.column or alias.column
		parts := strings.Split(columnName, ".")

		// Try full qualified name first (alias.column or table.column)
		if value, ok := row[columnName]; ok {
			return value
		}

		// Try just the column name (last part)
		actualColumn := parts[len(parts)-1]
		if value, ok := row[actualColumn]; ok {
			return value
		}

		// If column not found, return nil or error
		return nil
	}

	// Check for literal values
	if primExpr.Value() != nil {
		valueCtx := primExpr.Value()
		if valueCtx.NUMBER() != nil {
			numStr := valueCtx.NUMBER().GetText()
			if strings.Contains(numStr, ".") {
				if val, err := strconv.ParseFloat(numStr, 64); err == nil {
					return val
				}
			} else {
				if val, err := strconv.Atoi(numStr); err == nil {
					return val
				}
			}
		}

		if valueCtx.STRING_LITERAL() != nil {
			str := valueCtx.STRING_LITERAL().GetText()
			// Remove quotes
			if len(str) >= 2 && str[0] == '\'' && str[len(str)-1] == '\'' {
				return str[1 : len(str)-1]
			}
			return str
		}

		if valueCtx.TRUE() != nil {
			return true
		}

		if valueCtx.FALSE() != nil {
			return false
		}

		if valueCtx.CURRENT_TIMESTAMP() != nil {
			return time.Now().Format("2006-01-02 15:04:05")
		}
	}

	// If nothing matches, return the text
	return primExpr.GetText()
}

// parseMultiplicativeExpression handles multiplication and division
func parseMultiplicativeExpression(multExpr parser.IMultiplicativeExpressionContext, row map[string]interface{}) interface{} {
	if multExpr == nil {
		return nil
	}

	// Get all cast expressions
	castExprs := multExpr.AllCastExpression()
	if len(castExprs) == 0 {
		return nil
	}

	// Get the first cast expression
	result := parseCastExpression(castExprs[0], row)

	// Handle additional terms with operators
	stars := multExpr.AllSTAR()
	slashes := multExpr.AllSLASH()

	// Merge operators by token index
	type opInfo struct {
		op    string
		index int
	}
	var ops []opInfo
	for _, s := range stars {
		ops = append(ops, opInfo{"*", s.GetSymbol().GetTokenIndex()})
	}
	for _, s := range slashes {
		ops = append(ops, opInfo{"/", s.GetSymbol().GetTokenIndex()})
	}
	// Sort by token index
	sort.Slice(ops, func(i, j int) bool {
		return ops[i].index < ops[j].index
	})

	// Skip the first one since we already processed it
	for i, op := range ops {
		if i+1 >= len(castExprs) {
			break
		}
		nextValue := parseCastExpression(castExprs[i+1], row)

		switch op.op {
		case "*":
			result = multiplyValues(result, nextValue)
		case "/":
			result = divideValues(result, nextValue)
		}
	}

	return result
}

// parseCastExpression handles type casting with ::<type>
func parseCastExpression(castExpr parser.ICastExpressionContext, row map[string]interface{}) interface{} {
	if castExpr == nil {
		return nil
	}

	// Get the atTimeZoneExpression
	atExpr := castExpr.AtTimeZoneExpression()
	if atExpr == nil {
		return nil
	}

	// Get the value
	value := parseAtTimeZoneExpression(atExpr, row)

	// Check if there's a cast operator
	if castExpr.COLON_COLON() != nil && castExpr.IDENTIFIER() != nil {
		// Get the target type
		typeName := strings.ToLower(castExpr.IDENTIFIER().GetText())

		// Perform type casting based on target type
		switch typeName {
		case "int", "int4", "integer":
			if f, ok := toFloat64(value); ok {
				return int(f)
			}
			if s, ok := value.(string); ok {
				if i, err := strconv.Atoi(s); err == nil {
					return i
				}
			}
			return 0 // fallback
		case "bigint", "int8":
			if f, ok := toFloat64(value); ok {
				return int64(f)
			}
			if s, ok := value.(string); ok {
				if i, err := strconv.ParseInt(s, 10, 64); err == nil {
					return i
				}
			}
			return int64(0) // fallback
		case "float", "float8", "double precision":
			if f, ok := toFloat64(value); ok {
				return f
			}
			if s, ok := value.(string); ok {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					return f
				}
			}
			return 0.0 // fallback
		case "text", "varchar", "character varying":
			return fmt.Sprintf("%v", value)
		case "bool", "boolean":
			if b, ok := value.(bool); ok {
				return b
			}
			if s, ok := value.(string); ok {
				if strings.ToLower(s) == "true" {
					return true
				} else if strings.ToLower(s) == "false" {
					return false
				}
			}
			if i, ok := toFloat64(value); ok && i != 0 {
				return true
			}
			return false // fallback
		case "oid":
			// OID is typically int4
			if f, ok := toFloat64(value); ok {
				return int(f)
			}
			return 0
		case "name":
			// name is like text but limited
			return fmt.Sprintf("%v", value)
		default:
			// For unknown types, return as is
			return value
		}
	}

	// No cast, return the value as is
	return value
}

// parseAtTimeZoneExpression handles AT TIME ZONE expressions
func parseAtTimeZoneExpression(atExpr parser.IAtTimeZoneExpressionContext, row map[string]interface{}) interface{} {
	if atExpr == nil {
		return nil
	}

	// Get the primary expression (timestamp)
	primExpr := atExpr.PrimaryExpression()
	if primExpr == nil {
		return nil
	}

	value := parsePrimaryExpression(primExpr, row)

	// If there's an AT TIME ZONE clause, handle it
	if atExpr.AT() != nil && atExpr.STRING_LITERAL() != nil {
		// For simplicity, just return the value - proper timezone conversion would require more logic
		// In a real implementation, you'd convert the timestamp to the specified timezone
		return value
	}

	return value
}

// applyTimeZone applies time zone to a timestamp
func applyTimeZone(value interface{}, tzStr string) interface{} {
	// For simplicity, just return the value, ignore timezone
	return value
}
