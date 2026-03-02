package rparser

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"regexp"

	"go.uber.org/zap"
)

func getSubqueryCache(statement *statements.SelectStatement) map[*statements.SelectStatement]*ptypes.Row {
	if statement == nil {
		return nil
	}
	return statement.SubqueryCache
}

// parseComparisonExpression handles comparison expressions including IN, LIKE, IS NULL
func ParseComparisonExpression(
	allocator pmem.Allocator,
	ctx context.Context,
	compExpr parser.IComparisonExpressionContext,
	row *table.ResultRow,
	subExec subquery.SubqueryExecutor,
) (rmodels.Expression, error) {
	// logger.Debug("ParseComparisonExpression")
	if compExpr == nil {
		return nil, nil
	}

	concatExprs := compExpr.AllConcatExpression()
	if len(concatExprs) == 0 {
		return nil, nil
	}

	left, err := ParseConcatExpression(allocator, ctx, concatExprs[0], row, subExec)
	if err != nil {
		return nil, err
	}

	// Если правый операнд есть
	var right rmodels.Expression
	if len(concatExprs) > 1 {
		right, err = ParseConcatExpression(allocator, ctx, concatExprs[1], row, subExec)
		if err != nil {
			return nil, err
		}
	}

	// Если один из операндов — SubqueryExpression, извлекаем скалярное значение и используем кэш
	var scalarLeft, scalarRight *ptypes.Row
	if subq, ok := left.(*rmodels.SubqueryExpression); ok {
		cache := getSubqueryCache(subq.Select)
		if cache != nil {
			if val, ok := cache[subq.Select]; ok {
				scalarLeft = val
			} else {
				res, err := subExec(subq.Select)
				if err != nil {
					return nil, err
				}
				if res != nil && len(res.Rows) > 0 {
					scalarLeft = res.Rows[0]
					cache[subq.Select] = scalarLeft
				}
			}
		}
	}

	// Подставляем скалярные значения вместо выражений
	if scalarLeft != nil {
		// Создаем ResultRowsExpression с правильной структурой
		// Используем первое поле из scalarLeft для определения типа
		fields := []serializers.FieldDef{{
			Name: "?column?",
			OID:  ptypes.PTypeText, // default to text, может быть улучшено
		}}
		resultSchema := serializers.NewBaseSchema(fields)
		left = &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:   []*ptypes.Row{scalarLeft},
				Schema: resultSchema,
			},
		}
	}

	// OPERATOR: =, !=, <, >, <=, >=, LIKE, ~ etc (смотря что в грамматике)
	if compExpr.Operator() != nil {
		op := compExpr.Operator().GetText()

		// Случай: есть оператор в дереве, но нет правой части (у тебя так бывает)
		// Тогда ожидаем, что left — булево выражение (1x1 bool)
		if len(concatExprs) < 2 {
			switch l := left.(type) {
			case *rmodels.BoolExpression:
				return l, nil
			case *rmodels.ResultRowsExpression:
				// Проверим, что это 1x1 и bool
				// TODO добавить проверки на 1x1 они постоянно идут, надо это учитывать
				val, oid, err := l.Row.Schema.GetField(l.Row.Rows[0], 0)
				if err != nil {
					return nil, fmt.Errorf("error getting field from result row: %w", err)
				}
				desVal, err := serializers.DeserializeGeneric(val, oid)
				if err != nil {
					return nil, fmt.Errorf("error deserializing value: %w", err)
				}
				if oid == ptypes.PTypeBool {
					bv, ok := ptypes.TryIntoBool(desVal)
					if !ok {
						return nil, fmt.Errorf("expected bool value, got %T", desVal)
					}
					return &rmodels.BoolExpression{Value: bv.IntoGo()}, nil
				}
			default:
				return nil, fmt.Errorf("left expression is not boolean (got %T)", left)
			}
		}

		// Нормальный бинарный case: left OP right
		right, err = ParseConcatExpression(allocator, ctx, concatExprs[1], row, subExec)
		if err != nil {
			return nil, err
		}

		if subq, ok := right.(*rmodels.SubqueryExpression); ok {
			logger.Debug("Entering into subquery expression on right side of comparison", zap.Any("subquery", subq.Select))
			// cache := getSubqueryCache(subq.Select)
			// if cache != nil {
			// 	if val, ok := cache[subq.Select]; ok {
			// 		logger.Debug("Getting cached value for subquery", zap.Any("value", val))
			// 		scalarRight = val
			// 	} else {
			logger.Debug("Executing subquery for right side of comparison", zap.Any("subquery", subq.Select))
			res, err := subExec(subq.Select)
			if err != nil {
				return nil, err
			}
			logger.Debug("Getting value from subquery result", zap.Any("result", res))
			scalarRight = res.Rows[0]
			// 	cache[subq.Select] = scalarRight
			// }
			// }
			// }
		}

		if scalarRight != nil {
			// Создаем ResultRowsExpression с правильной структурой
			fields := []serializers.FieldDef{{
				Name: "?column?",
				OID:  ptypes.PTypeText, // default
			}}
			resultSchema := serializers.NewBaseSchema(fields)
			right = &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:   []*ptypes.Row{scalarRight},
					Schema: resultSchema,
				},
			}
		}

		lvv, okL := left.(*rmodels.ResultRowsExpression)
		rrr, okR := right.(*rmodels.ResultRowsExpression)
		if !okL || !okR {
			return nil, fmt.Errorf("comparison operands must be ResultRowsExpression, got %T and %T", left, right)
		}

		// Simplified approach: just use the expressions directly
		// Both should have Schema with at least one field
		if lvv == nil || lvv.Row == nil || lvv.Row.Schema == nil || len(lvv.Row.Schema.Fields) == 0 {
			return nil, fmt.Errorf("left operand has no result row or schema")
		}
		if rrr == nil || rrr.Row == nil || rrr.Row.Schema == nil || len(rrr.Row.Schema.Fields) == 0 {
			return nil, fmt.Errorf("right operand has no result row or schema")
		}

		// Get types from schema
		lType := lvv.Row.Schema.Fields[0].OID
		rType := rrr.Row.Schema.Fields[0].OID

		// Get values
		if len(lvv.Row.Rows) == 0 || len(rrr.Row.Rows) == 0 {
			return nil, fmt.Errorf("comparison operands have no data rows")
		}

		leftBuf, lOID, err := lvv.Row.Schema.GetField(lvv.Row.Rows[0], 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get left field: %w", err)
		}
		rightBuf, rOID, err := rrr.Row.Schema.GetField(rrr.Row.Rows[0], 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get right field: %w", err)
		}

		// Deserialize values
		leftVal, err := serializers.DeserializeGeneric(leftBuf, lOID)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize left value: %w", err)
		}
		rightVal, err := serializers.DeserializeGeneric(rightBuf, rOID)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize right value: %w", err)
		}

		// Type cast if needed: always cast right to left type
		if lType != rType {
			castable, ok := rightVal.(ptypes.CastableType[any])
			if !ok {
				return nil, fmt.Errorf("cannot compare types %d and %d: right operand (%T) does not support casting", lType, rType, rightVal)
			}
			converted, err := castable.CastTo(allocator, lType)
			if err != nil {
				return nil, fmt.Errorf("cannot cast right operand from type %d to %d: %w", rType, lType, err)
			}
			rightVal = converted
			rType = lType
		}

		// Compare values
		compareResult := leftVal.Compare(rightVal)
		switch op {
		case "=":
			return &rmodels.BoolExpression{Value: compareResult == 0}, nil
		case "!=", "<>":
			return &rmodels.BoolExpression{Value: compareResult != 0}, nil
		case "<":
			return &rmodels.BoolExpression{Value: compareResult < 0}, nil
		case "<=":
			return &rmodels.BoolExpression{Value: compareResult <= 0}, nil
		case ">":
			return &rmodels.BoolExpression{Value: compareResult > 0}, nil
		case ">=":
			return &rmodels.BoolExpression{Value: compareResult >= 0}, nil
		case "~":
			if lType == ptypes.PTypeText || lType == ptypes.PTypeVarchar {
				ls, ok := ptypes.TryIntoText[string](leftVal)
				if !ok {
					return nil, fmt.Errorf("left value is not text")
				}
				rs, ok := ptypes.TryIntoText[string](rightVal)
				if !ok {
					return nil, fmt.Errorf("right value is not text")
				}
				matched, _ := regexp.MatchString(rs.IntoGo(), ls.IntoGo())
				return &rmodels.BoolExpression{Value: matched}, nil
			}
			return nil, fmt.Errorf("regex operator ~ requires string operands")
		case "!~":
			if lType == ptypes.PTypeText || lType == ptypes.PTypeVarchar {
				ls, ok := ptypes.TryIntoText[string](leftVal)
				if !ok {
					return nil, fmt.Errorf("left value is not text")
				}
				rs, ok := ptypes.TryIntoText[string](rightVal)
				if !ok {
					return nil, fmt.Errorf("right value is not text")
				}
				matched, _ := regexp.MatchString(rs.IntoGo(), ls.IntoGo())
				return &rmodels.BoolExpression{Value: !matched}, nil
			}
			return nil, fmt.Errorf("regex operator !~ requires string operands")
		case "~*":
			if lType == ptypes.PTypeText || lType == ptypes.PTypeVarchar {
				ls, ok := ptypes.TryIntoText[string](leftVal)
				if !ok {
					return nil, fmt.Errorf("left value is not text")
				}
				rs, ok := ptypes.TryIntoText[string](rightVal)
				if !ok {
					return nil, fmt.Errorf("right value is not text")
				}
				matched, _ := regexp.MatchString("(?i)"+rs.IntoGo(), ls.IntoGo())
				return &rmodels.BoolExpression{Value: matched}, nil
			}
			return nil, fmt.Errorf("regex operator ~* requires string operands")
		case "!~*":
			if lType == ptypes.PTypeText || lType == ptypes.PTypeVarchar {
				ls, ok := ptypes.TryIntoText[string](leftVal)
				if !ok {
					return nil, fmt.Errorf("left value is not text")
				}
				rs, ok := ptypes.TryIntoText[string](rightVal)
				if !ok {
					return nil, fmt.Errorf("right value is not text")
				}
				matched, _ := regexp.MatchString("(?i)"+rs.IntoGo(), ls.IntoGo())
				return &rmodels.BoolExpression{Value: !matched}, nil
			}
			return nil, fmt.Errorf("regex operator !~* requires string operands")
		default:
			return nil, fmt.Errorf("unsupported comparison operator: %s", op)
		}
	}

	// OPERATOR(expr) (например OPERATOR(pg_catalog.~))
	if compExpr.OperatorExpr() != nil {
		if len(concatExprs) < 2 {
			return left, nil
		}
		opExpr := compExpr.OperatorExpr()
		right, err := ParseConcatExpression(allocator, ctx, concatExprs[1], row, subExec)
		if err != nil {
			return nil, err
		}

		resBool := rhelpers.EvaluateOperator(left, right, opExpr)
		return &rmodels.BoolExpression{Value: resBool}, nil
	}

	// IN (...)
	if compExpr.IN() != nil {
		not := compExpr.NOT() != nil

		// Получаем левое значение как ResultRowsExpression для доступа к типу
		leftExpr, ok := left.(*rmodels.ResultRowsExpression)
		if !ok {
			return nil, fmt.Errorf("IN operator requires ResultRowsExpression on left, got %T", left)
		}

		if leftExpr.Row == nil || leftExpr.Row.Schema == nil || len(leftExpr.Row.Schema.Fields) == 0 {
			return nil, fmt.Errorf("left expression has no schema")
		}
		if len(leftExpr.Row.Rows) == 0 {
			return nil, fmt.Errorf("left expression has no rows")
		}

		leftType := leftExpr.Row.Schema.Fields[0].OID
		leftBuf, _, err := leftExpr.Row.Schema.GetField(leftExpr.Row.Rows[0], 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get left field: %w", err)
		}
		leftValue, err := serializers.DeserializeGeneric(leftBuf, leftType)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize left value: %w", err)
		}

		// Проверяем каждое значение из IN списка
		found := false

		if sqCtx := compExpr.SubqueryExpression(); sqCtx != nil {
			selCtx := sqCtx.SelectStatement()
			if selCtx == nil {
				return nil, fmt.Errorf("expected subquery in IN operator")
			}
			selectStmt, err := ParseSelectStatement(selCtx)
			if err != nil {
				return nil, fmt.Errorf("error parsing subquery in IN operator: %w", err)
			}
			res, err := subExec(selectStmt)
			if err != nil {
				return nil, err
			}
			for _, r := range res.Rows {
				if res.Schema == nil || len(res.Schema.Fields) == 0 {
					continue
				}
				rightBuf, rOID, err := res.Schema.GetField(r, 0)
				if err != nil {
					continue
				}
				rightValue, err := serializers.DeserializeGeneric(rightBuf, rOID)
				if err != nil {
					continue
				}

				// Compare
				if lComp, ok := leftValue.(ptypes.BaseType[any]); ok {
					if rComp, ok := rightValue.(ptypes.BaseType[any]); ok {
						if lComp.Compare(rComp) == 0 {
							found = true
							break
						}
					}
				}
			}
		} else if compExpr.AllExpression() != nil {
			for _, e := range compExpr.AllExpression() {
				v, err := ParseExpression(allocator, ctx, e, row, subExec)
				if err != nil {
					return nil, err
				}

				rightExpr, ok := v.(*rmodels.ResultRowsExpression)
				if !ok {
					return nil, fmt.Errorf("IN value must be ResultRowsExpression, got %T", v)
				}

				if len(rightExpr.Row.Rows) == 0 || rightExpr.Row.Schema == nil {
					continue
				}
				rightBuf, _, err := rightExpr.Row.Schema.GetField(rightExpr.Row.Rows[0], 0)
				if err != nil {
					continue
				}
				rightType := rightExpr.Row.Schema.Fields[0].OID
				rightValue, err := serializers.DeserializeGeneric(rightBuf, rightType)
				if err != nil {
					continue
				}

				// Compare
				if lComp, ok := leftValue.(ptypes.BaseType[any]); ok {
					if rComp, ok := rightValue.(ptypes.BaseType[any]); ok {
						if lComp.Compare(rComp) == 0 {
							found = true
							break
						}
					}
				}
			}
		}

		if not {
			found = !found
		}
		return &rmodels.BoolExpression{Value: found}, nil
	}

	// LIKE
	if compExpr.LIKE() != nil {
		if len(concatExprs) < 2 {
			return left, nil
		}
		not := compExpr.NOT() != nil

		right, err := ParseConcatExpression(allocator, ctx, concatExprs[1], row, subExec)
		if err != nil {
			return nil, err
		}

		ls, err := exprToString(left)
		if err != nil {
			return nil, err
		}
		rs, err := exprToString(right)
		if err != nil {
			return nil, err
		}

		// используем TypeText для корректного SQL LIKE
		text := ptypes.TypeText{BufferPtr: []byte(ls)}
		match := text.Like(rs)
		if not {
			match = !match
		}
		return &rmodels.BoolExpression{Value: match}, nil
	}

	// IS NULL
	if compExpr.IS() != nil {
		not := compExpr.NOT() != nil

		// IS NULL для Expression: если это 1x1 и value==nil -> true
		isNull, err := exprIsNull(left)
		if err != nil {
			return nil, err
		}
		if not {
			isNull = !isNull
		}
		return &rmodels.BoolExpression{Value: isNull}, nil
	}

	// logger.Debug("Comparison expression result", zap.Any("result", left))
	return left, nil
}

func exprToScalar(e rmodels.Expression) (interface{}, error) {
	switch v := e.(type) {
	case *rmodels.BoolExpression:
		return v.Value, nil
	case *rmodels.ResultRowsExpression:
		if v.Row == nil || len(v.Row.Rows) != 1 || v.Row.Schema == nil || len(v.Row.Schema.Fields) != 1 {
			return nil, fmt.Errorf("expected 1x1 result for scalar, got %v", e.Type())
		}
		buf, oid, err := v.Row.Schema.GetField(v.Row.Rows[0], 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get field: %w", err)
		}
		retval, err := serializers.DeserializeGeneric(buf, oid)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize: %w", err)
		}
		return retval, nil
	default:
		return nil, fmt.Errorf("unsupported expression type for scalar: %T", e)
	}
}

func exprToString(e rmodels.Expression) (string, error) {
	val, ok := e.(*rmodels.ResultRowsExpression)
	if !ok {
		return "", fmt.Errorf("exprToString: expected ResultRowsExpression, got %T", e)
	}
	if len(val.Row.Rows) == 0 || len(val.Row.Schema.Fields) == 0 {
		return "", nil
	}

	buf, oid, err := val.Row.Schema.GetField(val.Row.Rows[0], 0)
	if err != nil {
		return "", fmt.Errorf("exprToString: get field: %w", err)
	}

	switch oid {
	case ptypes.PTypeText, ptypes.PTypeVarchar:
		return string(buf), nil
	default:
		// для остальных типов — через coerceToString
		return serializers.CoerceToString(buf, oid)
	}
}

func exprIsNull(e rmodels.Expression) (bool, error) {
	switch v := e.(type) {
	case *rmodels.BoolExpression:
		return false, nil // bool никогда не NULL в твоей модели
	case *rmodels.ResultRowsExpression:
		if v.Row == nil || len(v.Row.Rows) != 1 || v.Row.Schema == nil || len(v.Row.Schema.Fields) != 1 {
			return false, fmt.Errorf("expected 1x1 result for IS NULL")
		}
		buf, _, err := v.Row.Schema.GetField(v.Row.Rows[0], 0)
		if err != nil {
			return false, err
		}
		// Check if buf is empty (representing NULL)
		return len(buf) == 0, nil
	default:
		return false, fmt.Errorf("unsupported expression type for IS NULL: %T", e)
	}
}
