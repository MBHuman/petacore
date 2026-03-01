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
	"regexp"
	"strings"

	"go.uber.org/zap"
)

func getSubqueryCache(statement *statements.SelectStatement) map[*statements.SelectStatement]interface{} {
	if statement == nil {
		return nil
	}
	return statement.SubqueryCache
}

// parseComparisonExpression handles comparison expressions including IN, LIKE, IS NULL
func ParseComparisonExpression(ctx context.Context, compExpr parser.IComparisonExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	// getSubqueryCache извлекает кэш подзапросов из row, если возможно

	// logger.Debug("ParseComparisonExpression")
	if compExpr == nil {
		return nil, nil
	}

	concatExprs := compExpr.AllConcatExpression()
	if len(concatExprs) == 0 {
		return nil, nil
	}

	left, err := ParseConcatExpression(ctx, concatExprs[0], row, subExec)
	if err != nil {
		return nil, err
	}

	// Если правый операнд есть
	var right rmodels.Expression
	if len(concatExprs) > 1 {
		right, err = ParseConcatExpression(ctx, concatExprs[1], row, subExec)
		if err != nil {
			return nil, err
		}
	}

	// Если один из операндов — SubqueryExpression, извлекаем скалярное значение и используем кэш
	var scalarLeft, scalarRight interface{}
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
				if res != nil && len(res.Rows) > 0 && len(res.Rows[0]) > 0 {
					scalarLeft = res.Rows[0][0]
					cache[subq.Select] = scalarLeft
				}
			}
		}
	}

	// Подставляем скалярные значения вместо выражений
	if scalarLeft != nil {
		left = &rmodels.ResultRowsExpression{Row: &table.ExecuteResult{Rows: [][]interface{}{{scalarLeft}}, Columns: []table.TableColumn{{Type: table.ColTypeString}}}}
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
				if err := checkBoolSingleCell(l); err != nil {
					return nil, err
				}
				bv, ok := l.Row.Rows[0][0].(bool)
				if !ok {
					return nil, fmt.Errorf("expected bool value, got %T", l.Row.Rows[0][0])
				}
				return &rmodels.BoolExpression{Value: bv}, nil

			default:
				return nil, fmt.Errorf("left expression is not boolean (got %T)", left)
			}
		}

		// Нормальный бинарный case: left OP right
		right, err = ParseConcatExpression(ctx, concatExprs[1], row, subExec)
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
			// if res != nil && len(res.Rows) > 0 && len(res.Rows[0]) > 0 {
			scalarRight = res.Rows[0][0]
			// 	cache[subq.Select] = scalarRight
			// }
			// }
			// }
		}

		if scalarRight != nil {
			right = &rmodels.ResultRowsExpression{Row: &table.ExecuteResult{Rows: [][]interface{}{{scalarRight}}, Columns: []table.TableColumn{{Type: table.ColTypeString}}}}
		}

		lvv, okL := left.(*rmodels.ResultRowsExpression)
		rrr, okR := right.(*rmodels.ResultRowsExpression)
		if !okL || !okR {
			return nil, fmt.Errorf("comparison operands must be ResultRowsExpression, got %T and %T", left, right)
		}

		// If left expression represents a computed scalar (function result, literal,
		// or otherwise not bound to a table column), use its scalar value directly
		// instead of trying to map it to a column from the input row.
		var lrr *rmodels.ResultRowsExpression
		if lvv == nil || lvv.Row == nil || len(lvv.Row.Columns) == 0 {
			return nil, fmt.Errorf("left operand has no result row")
		}

		leftCol := lvv.Row.Columns[0]
		// If left is not a table-bound column (no TableIdentifier or placeholder name),
		// create a ResultRowsExpression from its scalar value.
		if leftCol.TableIdentifier == "" || leftCol.Name == "?column?" {
			leftVal := interface{}(nil)
			if len(lvv.Row.Rows) > 0 && len(lvv.Row.Rows[0]) > 0 {
				leftVal = lvv.Row.Rows[0][0]
			}
			lrr = &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:    [][]interface{}{{leftVal}},
					Columns: []table.TableColumn{{Type: leftCol.Type}},
				},
			}
		} else {
			// Try to find matching column in the input row by TableIdentifier/Name.
			idxLeft := -1
			for idx, col := range row.Columns {
				if col.TableIdentifier == leftCol.TableIdentifier && col.Name == leftCol.Name {
					idxLeft = idx
					break
				}
			}
			if idxLeft >= 0 && idxLeft < len(row.Row) {
				lrr = &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows:    [][]interface{}{{row.Row[idxLeft]}},
						Columns: []table.TableColumn{row.Columns[idxLeft]},
					},
				}
			} else {
				// Fallback: use computed scalar value
				leftVal := interface{}(nil)
				if len(lvv.Row.Rows) > 0 && len(lvv.Row.Rows[0]) > 0 {
					leftVal = lvv.Row.Rows[0][0]
				}
				lrr = &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows:    [][]interface{}{{leftVal}},
						Columns: []table.TableColumn{{Type: leftCol.Type}},
					},
				}
			}
		}

		if err := checkComparisonExpr(lrr, rrr); err != nil {
			return nil, err
		}

		lType := lrr.Row.Columns[0].Type
		rType := rrr.Row.Columns[0].Type

		// Если типы не совпадают, пробуем привести правую часть к типу левой
		rightValue := rrr.Row.Rows[0][0]
		if lType != rType {
			lOps := lType.TypeOps()
			convertedValue, err := lOps.CastTo(rightValue, lType)
			if err == nil {
				rightValue = convertedValue
				// Обновляем тип правой части
				rrr = &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows:    [][]interface{}{{rightValue}},
						Columns: []table.TableColumn{{Type: lType}},
					},
				}
			}
		}

		lOps := lType.TypeOps()

		compareResult, err := lOps.Compare(lrr.Row.Rows[0][0], rightValue, lType)
		if err != nil {
			return nil, fmt.Errorf("comparison error: %w", err)
		}
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
			if lType == table.ColTypeString {
				ls := lrr.Row.Rows[0][0].(string)
				rs := rightValue.(string)
				matched, _ := regexp.MatchString(rs, ls)
				return &rmodels.BoolExpression{Value: matched}, nil
			}
			return nil, fmt.Errorf("regex operator ~ requires string operands")
		case "!~":
			if lType == table.ColTypeString {
				ls := lrr.Row.Rows[0][0].(string)
				rs := rightValue.(string)
				matched, _ := regexp.MatchString(rs, ls)
				return &rmodels.BoolExpression{Value: !matched}, nil
			}
			return nil, fmt.Errorf("regex operator !~ requires string operands")
		case "~*":
			if lType == table.ColTypeString {
				ls := lrr.Row.Rows[0][0].(string)
				rs := rightValue.(string)
				matched, _ := regexp.MatchString("(?i)"+rs, ls)
				return &rmodels.BoolExpression{Value: matched}, nil
			}
			return nil, fmt.Errorf("regex operator ~* requires string operands")
		case "!~*":
			if lType == table.ColTypeString {
				ls := lrr.Row.Rows[0][0].(string)
				rs := rightValue.(string)
				matched, _ := regexp.MatchString("(?i)"+rs, ls)
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
		right, err := ParseConcatExpression(ctx, concatExprs[1], row, subExec)
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

		if err := checkComparisonExpr(leftExpr, leftExpr); err != nil {
			return nil, err
		}

		leftType := leftExpr.Row.Columns[0].Type
		leftValue := leftExpr.Row.Rows[0][0]
		leftOps := leftType.TypeOps()

		// Проверяем каждое значение из IN списка
		found := false

		if sqCtx := compExpr.SubqueryExpression(); sqCtx != nil {
			sqCtx := compExpr.SubqueryExpression()
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
				if len(r) == 0 {
					continue
				}
				rightValue := r[0]
				rightType := res.Columns[0].Type

				// Приводим правое значение к типу левого, если типы не совпадают
				if leftType != rightType {
					convertedValue, err := leftOps.CastTo(rightValue, leftType)
					if err == nil {
						rightValue = convertedValue
					}
				}

				// Используем Compare для проверки равенства
				compareResult, err := leftOps.Compare(leftValue, rightValue, leftType)
				if err == nil && compareResult == 0 {
					found = true
					break
				}
			}
		} else if compExpr.AllExpression() != nil {
			for _, e := range compExpr.AllExpression() {
				v, err := ParseExpression(ctx, e, row, subExec)
				if err != nil {
					return nil, err
				}

				rightExpr, ok := v.(*rmodels.ResultRowsExpression)
				if !ok {
					return nil, fmt.Errorf("IN value must be ResultRowsExpression, got %T", v)
				}

				rightValue := rightExpr.Row.Rows[0][0]
				rightType := rightExpr.Row.Columns[0].Type

				// Приводим правое значение к типу левого, если типы не совпадают
				if leftType != rightType {
					convertedValue, err := leftOps.CastTo(rightValue, leftType)
					if err == nil {
						rightValue = convertedValue
					}
				}

				// Используем Compare для проверки равенства
				compareResult, err := leftOps.Compare(leftValue, rightValue, leftType)
				if err == nil && compareResult == 0 {
					found = true
					break
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

		right, err := ParseConcatExpression(ctx, concatExprs[1], row, subExec)
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

		match := strings.Contains(ls, strings.ReplaceAll(strings.ReplaceAll(rs, "%", ""), "_", ""))
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

	logger.Debug("Comparison expression result", zap.Any("result", left))
	return left, nil
}

// ---- helpers ----

func checkBoolSingleCell(a *rmodels.ResultRowsExpression) error {
	if a == nil || a.Row == nil {
		return fmt.Errorf("nil operand")
	}
	if len(a.Row.Columns) != 1 {
		return fmt.Errorf("expected 1 column, got %d", len(a.Row.Columns))
	}
	if len(a.Row.Rows) != 1 || len(a.Row.Rows[0]) != 1 {
		return fmt.Errorf("expected 1x1 result")
	}
	if a.Row.Columns[0].Type != table.ColTypeBool {
		return fmt.Errorf("expected bool column, got %v", a.Row.Columns[0].Type)
	}
	return nil
}

func checkComparisonExpr(a, b *rmodels.ResultRowsExpression) error {
	if a == nil || b == nil || a.Row == nil || b.Row == nil {
		return fmt.Errorf("nil operand in comparison")
	}
	if len(a.Row.Columns) != 1 || len(b.Row.Columns) != 1 {
		return fmt.Errorf("multiple columns in comparison not supported")
	}
	if len(a.Row.Rows) != 1 || len(b.Row.Rows) != 1 || len(a.Row.Rows[0]) != 1 || len(b.Row.Rows[0]) != 1 {
		return fmt.Errorf("comparison supports only 1x1 values")
	}
	// Типы не обязаны совпадать - будем пытаться приводить при сравнении
	return nil
}

func exprToScalar(e rmodels.Expression) (interface{}, error) {
	switch v := e.(type) {
	case *rmodels.BoolExpression:
		return v.Value, nil
	case *rmodels.ResultRowsExpression:
		if v.Row == nil || len(v.Row.Rows) != 1 || len(v.Row.Rows[0]) != 1 {
			return nil, fmt.Errorf("expected 1x1 result for scalar, got %v", e.Type())
		}
		return v.Row.Rows[0][0], nil
	default:
		return nil, fmt.Errorf("unsupported expression type for scalar: %T", e)
	}
}

func exprToString(e rmodels.Expression) (string, error) {
	val, err := exprToScalar(e)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	return fmt.Sprintf("%v", val), nil
}

func exprIsNull(e rmodels.Expression) (bool, error) {
	switch v := e.(type) {
	case *rmodels.BoolExpression:
		return false, nil // bool никогда не NULL в твоей модели
	case *rmodels.ResultRowsExpression:
		if v.Row == nil || len(v.Row.Rows) != 1 || len(v.Row.Rows[0]) != 1 {
			return false, fmt.Errorf("expected 1x1 result for IS NULL")
		}
		return v.Row.Rows[0][0] == nil, nil
	default:
		return false, fmt.Errorf("unsupported expression type for IS NULL: %T", e)
	}
}
