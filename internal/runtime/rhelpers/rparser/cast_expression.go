package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"strings"
)

// parseCastExpression handles type casting with ::<type>, COLLATE, AT TIME ZONE
func ParseCastExpression(castExpr parser.ICastExpressionContext, row *table.ResultRow) (result rmodels.Expression, err error) {
	// logger.Debug("ParseCastExpression")
	if castExpr == nil {
		return nil, nil
	}

	// Get the primary expression
	primExpr := castExpr.PrimaryExpression()
	if primExpr == nil {
		return nil, nil
	}

	// Parse the primary expression
	value, err := ParsePrimaryExpression(primExpr, row)
	if err != nil {
		return nil, err
	}

	// Apply postfixes
	postfixes := castExpr.AllPostfix()
	for _, postfix := range postfixes {
		if postfix.AT() != nil && postfix.TIME() != nil && postfix.ZONE() != nil && postfix.STRING_LITERAL() != nil {
			// AT TIME ZONE
			if val, ok := value.(*rmodels.ResultRowsExpression); ok {
				value = ApplyTimeZone(val, postfix.STRING_LITERAL().GetText())
			} else {
				return nil, fmt.Errorf("AT TIME ZONE can only be applied to timestamp expressions")
			}
		} else if postfix.COLLATE() != nil && postfix.QualifiedName() != nil {
			// COLLATE - for now, ignore
			// value = applyCollate(value, postfix.QualifiedName().GetText())
		} else {
			castingTypes := postfix.AllTypeName()
			for _, castingOp := range castingTypes {
				typeName := strings.ToLower(castingOp.QualifiedName().GetText())
				colType := table.ColTypeFromString(typeName)
				if val, ok := value.(*rmodels.ResultRowsExpression); ok {
					value, err = CastValue(val, colType)
				} else if val, ok := value.(*rmodels.BoolExpression); ok {
					newVal := &rmodels.ResultRowsExpression{
						Row: &table.ExecuteResult{
							Rows:    [][]interface{}{{val.Value}},
							Columns: []table.TableColumn{{Type: table.ColTypeBool}},
						},
					}
					value, err = CastValue(newVal, colType)
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return value, nil
}

// applyTimeZone applies time zone to a timestamp
// TODO реализовать корректную работу с часовыми поясами
func ApplyTimeZone(value *rmodels.ResultRowsExpression, tzStr string) rmodels.Expression {
	// For simplicity, just return the value, ignore timezone
	return value
}

// castValue performs type casting to the specified type
// TODO перевести cast в отдельный модуль + пересмотреть работу с типами
func CastValue(expression *rmodels.ResultRowsExpression, colType table.ColType) (rmodels.Expression, error) {
	err := checkCastExpr(expression, colType)
	if err != nil {
		return nil, err
	}

	rows := expression.Row.Rows
	val := rows[0][0]

	typeOps := expression.Row.Columns[0].Type.TypeOps()
	castedVal, err := typeOps.CastTo(val, colType)
	if err != nil {
		return nil, err
	}

	// Сохраняем имя колонки из исходного выражения
	colName := "?column?"
	if len(expression.Row.Columns) > 0 {
		colName = expression.Row.Columns[0].Name
	}

	return &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows:    [][]interface{}{{castedVal}},
			Columns: []table.TableColumn{{Name: colName, Type: colType}},
		},
	}, nil
}

func checkCastExpr(a *rmodels.ResultRowsExpression, colType table.ColType) error {
	if a == nil {
		return fmt.Errorf("nil operand in cast")
	}
	if len(a.Row.Rows) == 0 {
		return fmt.Errorf("empty rows in cast")
	}
	if len(a.Row.Columns) == 0 {
		return fmt.Errorf("empty columns in cast")
	}
	if len(a.Row.Columns) > 1 {
		return fmt.Errorf("cast supports only single-column expressions")
	}
	if len(a.Row.Rows) > 1 {
		return fmt.Errorf("cast supports only single-row expressions")
	}
	return nil
}
