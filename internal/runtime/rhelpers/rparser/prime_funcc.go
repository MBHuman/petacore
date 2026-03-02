// PrimaryExpression FunctionCall
package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/functions"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func ParsePrimaryExpressionFunctionCall(
	allocator pmem.Allocator,
	ctx context.Context,
	funcCall parser.IFunctionCallContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (rmodels.Expression, error) {
	fc := funcCall
	funcName := ""
	if qn := fc.QualifiedName(); qn != nil {
		nameParts := qn.AllNamePart()
		if len(nameParts) > 0 {
			parts := []string{}
			for _, np := range nameParts {
				parts = append(parts, np.GetText())
			}
			funcName = parts[len(parts)-1]
		}
	}
	var args []interface{}
	for _, argExpr := range fc.AllFunctionArg() {
		if argExpr.STAR() != nil {
			// Handle COUNT(*) case - pass a special marker value
			args = append(args, "*")
			continue
		}

		if argExpr.Expression() == nil {
			return nil, fmt.Errorf("[ParsePrimaryExpressionFunctionCall] unsupported function argument: %s", argExpr.GetText())
		}

		val, err := ParseExpression(allocator, ctx, argExpr.Expression(), row, subExec)
		if err != nil {
			return nil, err
		}
		// TODO переделать передачу аргументов функций
		// на строго типизированную передачу Expression
		if valExpr, ok := val.(*rmodels.ResultRowsExpression); ok {
			if len(valExpr.Row.Rows) == 1 && len(valExpr.Row.Schema.Fields) == 1 {
				buf, oid, err := valExpr.Row.Schema.GetField(valExpr.Row.Rows[0], 0)
				if err != nil {
					return nil, fmt.Errorf("[ParsePrimaryExpressionFunctionCall] failed to get field: %w", err)
				}
				desVal, err := serializers.DeserializeGeneric(buf, oid)
				if err != nil {
					return nil, fmt.Errorf("[ParsePrimaryExpressionFunctionCall] failed to deserialize arg: %w", err)
				}
				args = append(args, desVal)
			} else {
				return nil, fmt.Errorf("[ParsePrimaryExpressionFunctionCall] function arguments must be single-row single-column expressions")
			}
			// if len(valExpr.Row.Rows) > 0 && len(valExpr.Row.Rows[0]) > 0 {
			// 	args = append(args, valExpr.Row.Rows[0][0])
			// } else {
			// 	args = append(args, nil)
			// }
		}
		// args = append(args, ParseExpression(argExpr, row))
	}

	value, err := functions.ExecuteFunctionWithContext(allocator, ctx, funcName, args)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("function %s returned nil result", funcName)
	}
	// If function returns a single boolean value, return BoolExpression for WHERE conditions
	if len(value.Rows) == 1 && len(value.Schema.Fields) == 1 && value.Schema.Fields[0].OID == ptypes.PTypeBool {
		field, _, err := value.Schema.GetField(value.Rows[0], 0)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryExpressionFunctionCall] error getting field from function result: %v", err)
		}
		boolVal, err := serializers.BoolSerializerInstance.Deserialize(field)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryExpressionFunctionCall] error deserializing boolean value from function result: %v", err)
		}
		return &rmodels.BoolExpression{Value: boolVal.IntoGo()}, nil
	}
	return &rmodels.ResultRowsExpression{Row: value}, nil
}
