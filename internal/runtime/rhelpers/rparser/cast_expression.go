package rparser

import (
	"context"
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/rops"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"strings"
)

func ParseCastExpression(
	allocator pmem.Allocator,
	ctx context.Context,
	castExpr parser.ICastExpressionContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (rmodels.Expression, error) {
	if castExpr == nil {
		return nil, nil
	}

	primExpr := castExpr.PrimaryExpression()
	if primExpr == nil {
		return nil, nil
	}

	value, err := ParsePrimaryExpression(allocator, ctx, primExpr, row, subExec)
	if err != nil {
		return nil, err
	}

	for _, postfix := range castExpr.AllPostfix() {
		switch {
		case postfix.AT() != nil && postfix.TIME() != nil && postfix.ZONE() != nil:
			val, ok := value.(*rmodels.ResultRowsExpression)
			if !ok {
				return nil, fmt.Errorf("[ParseCastExpression] AT TIME ZONE: expected result expression")
			}
			value = ApplyTimeZone(val, postfix.STRING_LITERAL().GetText())

		case postfix.COLLATE() != nil:
			// игнорируем COLLATE

		default:
			for _, castingOp := range postfix.AllTypeName() {
				typeName := strings.ToLower(castingOp.QualifiedName().GetText())
				targetOID := ptypes.ColTypeFromString(typeName)

				val, ok := value.(*rmodels.ResultRowsExpression)
				if !ok {
					return nil, fmt.Errorf("[ParseCastExpression] cast: expected result expression, got %T", value)
				}

				value, err = rops.CastValue(allocator, val, targetOID)
				if err != nil {
					return nil, fmt.Errorf("[ParseCastExpression] cast to %s: %w", typeName, err)
				}
			}
		}
	}

	return value, nil
}

func ApplyTimeZone(value *rmodels.ResultRowsExpression, tzStr string) rmodels.Expression {
	// TODO: реализовать корректную работу с часовыми поясами
	return value
}
