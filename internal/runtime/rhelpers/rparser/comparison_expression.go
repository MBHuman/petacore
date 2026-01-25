package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"regexp"
	"strings"
)

// parseComparisonExpression handles comparison expressions including IN, LIKE, IS NULL
func ParseComparisonExpression(compExpr parser.IComparisonExpressionContext, row *table.ResultRow) (rmodels.Expression, error) {
	// logger.Debug("ParseComparisonExpression")
	if compExpr == nil {
		return nil, nil
	}

	concatExprs := compExpr.AllConcatExpression()
	if len(concatExprs) == 0 {
		return nil, nil
	}

	left, err := ParseConcatExpression(concatExprs[0], row)
	if err != nil {
		return nil, err
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
		right, err := ParseConcatExpression(concatExprs[1], row)
		if err != nil {
			return nil, err
		}

		lvv, okL := left.(*rmodels.ResultRowsExpression)
		idxLeft := 0

		for idx, col := range row.Columns {
			if col.Name == lvv.Row.Columns[0].Name {
				idxLeft = idx
				break
			}
		}

		lrr := &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:    [][]interface{}{{row.Row[idxLeft]}},
				Columns: []table.TableColumn{row.Columns[idxLeft]},
			},
		}

		rrr, okR := right.(*rmodels.ResultRowsExpression)
		if !okL || !okR {
			return nil, fmt.Errorf("comparison operands must be ResultRowsExpression, got %T and %T", left, right)
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
		right, err := ParseConcatExpression(concatExprs[1], row)
		if err != nil {
			return nil, err
		}

		resBool := rhelpers.EvaluateOperator(left, right, opExpr)
		return &rmodels.BoolExpression{Value: resBool}, nil
	}

	// IN (...)
	if compExpr.IN() != nil {
		// Тут у тебя старая логика была на interface{} — если хочешь IN для ResultRowsExpression,
		// нужно отдельно доставать single cell из left.
		not := compExpr.NOT() != nil

		var values []interface{}
		for _, e := range compExpr.AllExpression() {
			v, err := ParseExpression(e, row) // если у тебя ParseExpression возвращает interface{}
			if err != nil {
				return nil, err
			}
			values = append(values, v)
		}

		// left может быть Expression — под IN разумно поддержать 1x1 scalar
		lv, err := exprToScalar(left)
		if err != nil {
			return nil, err
		}

		in := rhelpers.Contains(values, lv)
		if not {
			in = !in
		}
		return &rmodels.BoolExpression{Value: in}, nil
	}

	// LIKE
	if compExpr.LIKE() != nil {
		if len(concatExprs) < 2 {
			return left, nil
		}
		not := compExpr.NOT() != nil

		right, err := ParseConcatExpression(concatExprs[1], row)
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
