package rhelpers

import (
	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

// EvaluateOperator evaluates an operator expression like OPERATOR(pg_catalog.~)
// TODO: Review operator handling and remove hardcoding
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
	var leftVal, rightVal ptypes.BaseType[any]
	if l, ok := left.(*rmodels.ResultRowsExpression); ok {
		if len(l.Row.Rows) == 1 && len(l.Row.Schema.Fields) == 1 {
			fieldVal, oid, err := l.Row.Schema.GetField(l.Row.Rows[0], 0)
			if err != nil {
				logger.Errorf("EvaluateOperator: failed to get field from left expression: %v", err)
				return false
			}
			desVal, err := serializers.DeserializeGeneric(fieldVal, oid)
			if err != nil {
				logger.Errorf("EvaluateOperator: failed to deserialize left value: %v", err)
				return false
			}
			leftVal = desVal
		} else {
			logger.Errorf("EvaluateOperator: left expression has multiple rows or columns, which is not supported")
			return false
		}
	}
	if r, ok := right.(*rmodels.ResultRowsExpression); ok {
		if len(r.Row.Rows) == 1 && len(r.Row.Schema.Fields) == 1 {
			fieldVal, oid, err := r.Row.Schema.GetField(r.Row.Rows[0], 0)
			if err != nil {
				logger.Errorf("EvaluateOperator: failed to get field from right expression: %v", err)
				return false
			}
			desVal, err := serializers.DeserializeGeneric(fieldVal, oid)
			if err != nil {
				logger.Errorf("EvaluateOperator: failed to deserialize right value: %v", err)
				return false
			}
			rightVal = desVal
		} else {
			logger.Errorf("EvaluateOperator: right expression has multiple rows or columns, which is not supported")
			return false
		}
	}
	lText, ok := ptypes.TryIntoText[string](leftVal)
	rText, ok2 := ptypes.TryIntoText[string](rightVal)
	if !ok || !ok2 {
		logger.Errorf("EvaluateOperator: both left and right values must be text for regex operators, got %T and %T", leftVal, rightVal)
		return false
	}

	logger.Debugf("EvaluateOperator: opName=%s, leftVal=%v (%T), rightVal=%v (%T)", opName, leftVal, leftVal, rightVal, rightVal)

	// For now, handle specific operators
	switch opName {
	case "pg_catalog.~":
		matched, err := lText.RegexpMatch(rText.AsStr())
		if err != nil {
			logger.Errorf("EvaluateOperator: failed to match regex: %v", err)
			return false
		}
		return matched
	case "pg_catalog.!~":
		matched, err := lText.RegexpMatch(rText.AsStr())
		if err != nil {
			logger.Errorf("EvaluateOperator: failed to match regex: %v", err)
			return false
		}
		return !matched
	case "pg_catalog.~*":
		matched, err := lText.RegexpMatch("(?i)" + rText.AsStr())
		if err != nil {
			logger.Errorf("EvaluateOperator: failed to match regex: %v", err)
			return false
		}
		return matched
	case "pg_catalog.!~*":
		matched, err := lText.RegexpMatch("(?i)" + rText.AsStr())
		if err != nil {
			logger.Errorf("EvaluateOperator: failed to match regex: %v", err)
			return false
		}
		return !matched
	}

	return false
}
