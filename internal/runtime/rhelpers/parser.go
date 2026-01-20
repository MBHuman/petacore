package rhelpers

import (
	"fmt"
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
func ParseExpression(expr parser.IExpressionContext) interface{} {
	if expr == nil {
		return nil
	}

	// Check for castExpression
	if castExpr := expr.CastExpression(); castExpr != nil {
		return parseCastExpression(castExpr)
	}

	// Check for atTimeZoneExpression
	if atExpr := expr.AtTimeZoneExpression(); atExpr != nil {
		return ParseAtTimeZoneExpression(atExpr)
	}

	// Check for caseExpression
	if caseExpr := expr.CaseExpression(); caseExpr != nil {
		return ParseCaseExpression(caseExpr)
	}

	// Get the additive expression (top level)
	addExpr := expr.AdditiveExpression(0)
	if addExpr == nil {
		return nil
	}

	return ParseAdditiveExpression(addExpr)
}

// parseAdditiveExpression handles addition and subtraction
func ParseAdditiveExpression(addExpr parser.IAdditiveExpressionContext) interface{} {
	if addExpr == nil {
		return nil
	}

	// Get the first multiplicative expression
	multExpr := addExpr.MultiplicativeExpression(0)
	if multExpr == nil {
		return nil
	}

	result := parseMultiplicativeExpression(multExpr)

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
		nextValue := parseMultiplicativeExpression(multExprs[i+1])

		if op.op == "+" {
			result = addValues(result, nextValue)
		} else if op.op == "-" {
			result = subtractValues(result, nextValue)
		}
	}

	return result
}

// parseAtTimeZoneExpression handles AT TIME ZONE expressions
func ParseAtTimeZoneExpression(atExpr parser.IAtTimeZoneExpressionContext) interface{} {
	if atExpr == nil {
		return nil
	}

	// Get the primary expression (timestamp)
	primExpr := atExpr.PrimaryExpression()
	if primExpr == nil {
		return nil
	}

	// For simplicity, just return the primary expression value, ignore timezone
	return parsePrimaryExpression(primExpr)
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

// parseExtractFunction handles EXTRACT expressions
func ParseExtractFunction(extractExpr parser.IExtractFunctionContext) interface{} {
	if extractExpr == nil {
		return nil
	}

	field := extractExpr.IDENTIFIER().GetText()
	sourceExpr := extractExpr.Expression()
	if sourceExpr == nil {
		return nil
	}

	source := ParseExpression(sourceExpr)
	args := []interface{}{field, source}
	value, _ := functions.ExecuteFunction("EXTRACT", args)
	return value
}

// parsePrimaryExpression handles the basic expressions
func parsePrimaryExpression(primExpr parser.IPrimaryExpressionContext) interface{} {
	if primExpr == nil {
		return nil
	}

	// Check for parenthesized expression
	if primExpr.Expression() != nil {
		return ParseExpression(primExpr.Expression())
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
			args = append(args, ParseExpression(argExpr))
		}
		value, _ := functions.ExecuteFunction(funcName, args)
		return value
	}

	// Check for extract function
	if primExpr.ExtractFunction() != nil {
		return ParseExtractFunction(primExpr.ExtractFunction())
	}

	// Check for column name
	if primExpr.ColumnName() != nil {
		// For now, return the column name as string
		// In a real implementation, this would look up the column value
		return primExpr.ColumnName().GetText()
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
func parseMultiplicativeExpression(multExpr parser.IMultiplicativeExpressionContext) interface{} {
	if multExpr == nil {
		return nil
	}

	// Get the first primary expression
	primExpr := multExpr.PrimaryExpression(0)
	if primExpr == nil {
		return nil
	}

	result := parsePrimaryExpression(primExpr)

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

	primExprs := multExpr.AllPrimaryExpression()
	// Skip the first one since we already processed it
	for i, op := range ops {
		if i+1 >= len(primExprs) {
			break
		}
		nextValue := parsePrimaryExpression(primExprs[i+1])

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
func parseCastExpression(castExpr parser.ICastExpressionContext) interface{} {
	if castExpr == nil {
		return nil
	}

	// Get the primary expression to cast
	primExpr := castExpr.PrimaryExpression()
	if primExpr == nil {
		return nil
	}

	// Get the target type
	typeName := strings.ToLower(castExpr.IDENTIFIER().GetText())

	// Parse the value
	value := parsePrimaryExpression(primExpr)

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
