package rhelpers

import (
	"fmt"
	"reflect"
	"regexp"

	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"strconv"
	"strings"
)

// TODO пропал concatenation в evaluate - нужно добавить обратно

func IsTrueHelper(value interface{}) bool {
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

func EvaluateComparison(left, right *rmodels.ResultRowsExpression, op string) bool {
	logger.Debugf("Evaluating comparison: left=%v, right=%v, op=%s\n", left, right, op)
	lv, rv, ok := oneValueSameType(left, right)
	if !ok {
		return false
	}

	switch op {
	case "=", "!=":
		eq := equal(lv, rv)
		if op == "!=" {
			return !eq
		}
		return eq

	case "<", ">", "<=", ">=":
		return orderCompare(lv, rv, op)

	case "LIKE":
		ls, ok := lv.(string)
		if !ok {
			return false
		}
		return strings.Contains(ls, rv.(string))

	case "~", "!~", "~*", "!~*":
		ls, ok := lv.(string)
		if !ok {
			return false
		}
		pat := rv.(string)
		if op == "~*" || op == "!~*" {
			pat = "(?i)" + pat
		}
		matched, err := regexp.MatchString(pat, ls)
		if err != nil {
			return false
		}
		if op == "!~" || op == "!~*" {
			return !matched
		}
		return matched
	}

	return false
}

func oneValueSameType(l, r *rmodels.ResultRowsExpression) (interface{}, interface{}, bool) {
	if l == nil || r == nil || len(l.Row.Rows) != 1 || len(r.Row.Rows) != 1 {
		return nil, nil, false
	}
	lv, rv := l.Row.Rows[0][0], r.Row.Rows[0][0]
	if lv == nil || rv == nil || reflect.TypeOf(lv) != reflect.TypeOf(rv) {
		return nil, nil, false
	}
	return lv, rv, true
}

func equal(a, b interface{}) bool {
	switch av := a.(type) {
	case int:
		return av == b.(int)
	case float64:
		return av == b.(float64)
	case string:
		return av == b.(string)
	case bool:
		return av == b.(bool)
	default:
		return false
	}
}

func orderCompare(a, b interface{}, op string) bool {
	switch av := a.(type) {
	case int:
		bv := b.(int)
		switch op {
		case "<":
			return av < bv
		case ">":
			return av > bv
		case "<=":
			return av <= bv
		case ">=":
			return av >= bv
		}
	case float64:
		bv := b.(float64)
		switch op {
		case "<":
			return av < bv
		case ">":
			return av > bv
		case "<=":
			return av <= bv
		case ">=":
			return av >= bv
		}
	case string:
		bv := b.(string)
		switch op {
		case "<":
			return av < bv
		case ">":
			return av > bv
		case "<=":
			return av <= bv
		case ">=":
			return av >= bv
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
func Contains(list []interface{}, value interface{}) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// evaluateOperator evaluates an operator expression like OPERATOR(pg_catalog.~)
// TODO пересмотреть работу с операторами, убрать хардкодинг
func EvaluateOperator(left, right rmodels.Expression, opExpr parser.IOperatorExprContext) bool {
	if opExpr == nil {
		return false
	}

	// Get the qualified name of the operator
	qualName := opExpr.QualifiedName()
	if qualName == nil {
		return false
	}

	// Get the operator name (e.g., "pg_catalog.~")
	opName := qualName.GetText()

	// Extract values
	var leftVal, rightVal interface{}
	if l, ok := left.(*rmodels.ResultRowsExpression); ok && len(l.Row.Rows) > 0 && len(l.Row.Rows[0]) > 0 {
		leftVal = l.Row.Rows[0][0]
	}
	if r, ok := right.(*rmodels.ResultRowsExpression); ok && len(r.Row.Rows) > 0 && len(r.Row.Rows[0]) > 0 {
		rightVal = r.Row.Rows[0][0]
	}

	logger.Debugf("EvaluateOperator: opName=%s, leftVal=%v (%T), rightVal=%v (%T)", opName, leftVal, leftVal, rightVal, rightVal)

	// For now, handle specific operators
	switch opName {
	case "pg_catalog.~":
		if la, ok := leftVal.(string); ok {
			if ra, ok := rightVal.(string); ok {
				matched, err := regexp.MatchString(ra, la)
				logger.Debugf("EvaluateOperator ~: pattern=%s, str=%s, matched=%v, err=%v", ra, la, matched, err)
				return matched
			}
		}
	case "pg_catalog.!~":
		if la, ok := leftVal.(string); ok {
			if ra, ok := rightVal.(string); ok {
				matched, err := regexp.MatchString(ra, la)
				logger.Debugf("EvaluateOperator !~: pattern=%s, str=%s, matched=%v, err=%v", ra, la, matched, err)
				return !matched
			}
		}
	case "pg_catalog.~*":
		if la, ok := leftVal.(string); ok {
			if ra, ok := rightVal.(string); ok {
				matched, err := regexp.MatchString("(?i)"+ra, la)
				logger.Debugf("EvaluateOperator ~*: pattern=%s, str=%s, matched=%v, err=%v", "(?i)"+ra, la, matched, err)
				return matched
			}
		}
	case "pg_catalog.!~*":
		if la, ok := leftVal.(string); ok {
			if ra, ok := rightVal.(string); ok {
				matched, err := regexp.MatchString("(?i)"+ra, la)
				logger.Debugf("EvaluateOperator !~*: pattern=%s, str=%s, matched=%v, err=%v", "(?i)"+ra, la, matched, err)
				return !matched
			}
		}
	}

	return false
}
