package rops

import (
	"fmt"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/utils"
)

// Helper functions for arithmetic operations
func AddValues(a, b *rmodels.ResultRowsExpression) (*rmodels.ResultRowsExpression, error) {

	if err := checkOps(a, b); err != nil {
		return nil, err
	}

	aRows := a.Row.Rows
	bRows := b.Row.Rows

	aVal := aRows[0][0]
	bVal := bRows[0][0]
	typ := a.Row.Columns[0].Type

	val, err := applyAdd(aVal, bVal, typ)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: typ},
		},
	}
	return &rmodels.ResultRowsExpression{Row: resultRow}, nil
}

func SubtractValues(a, b *rmodels.ResultRowsExpression) (*rmodels.ResultRowsExpression, error) {

	if err := checkOps(a, b); err != nil {
		return nil, err
	}

	aRows := a.Row.Rows
	bRows := b.Row.Rows

	aVal := aRows[0][0]
	bVal := bRows[0][0]
	typ := a.Row.Columns[0].Type

	val, err := applySubtract(aVal, bVal, typ)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: typ},
		},
	}
	return &rmodels.ResultRowsExpression{Row: resultRow}, nil
}

func MultiplyValues(a, b *rmodels.ResultRowsExpression) (*rmodels.ResultRowsExpression, error) {

	if err := checkOps(a, b); err != nil {
		return nil, err
	}

	aRows := a.Row.Rows
	bRows := b.Row.Rows

	aVal := aRows[0][0]
	bVal := bRows[0][0]
	typ := a.Row.Columns[0].Type

	val, err := applyMultiply(aVal, bVal, typ)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: typ},
		},
	}
	return &rmodels.ResultRowsExpression{Row: resultRow}, nil
}

func DivideValues(a, b *rmodels.ResultRowsExpression) (*rmodels.ResultRowsExpression, error) {

	if err := checkOps(a, b); err != nil {
		return nil, err
	}

	aRows := a.Row.Rows
	bRows := b.Row.Rows

	aVal := aRows[0][0]
	bVal := bRows[0][0]
	typ := a.Row.Columns[0].Type

	val, err := applyDivide(aVal, bVal, typ)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: typ},
		},
	}
	return &rmodels.ResultRowsExpression{Row: resultRow}, nil
}

func checkOps(a, b *rmodels.ResultRowsExpression) error {
	aRows := a.Row.Rows
	bRows := b.Row.Rows

	if len(aRows) > 1 || len(bRows) > 1 {
		return fmt.Errorf("AddValues: multiple rows not supported")
	}

	if len(aRows) == 0 || len(bRows) == 0 {
		return fmt.Errorf("AddValues: no rows to operate on")
	}

	if a.Row.Columns[0].Type != b.Row.Columns[0].Type {
		return fmt.Errorf("AddValues: column types do not match")
	}

	if !(a.Row.Columns[0].Type == table.ColTypeInt ||
		a.Row.Columns[0].Type == table.ColTypeFloat) {
		return fmt.Errorf("AddValues: use type: %s", a.Row.Columns[0].Type.String())
	}
	return nil
}

func applyAdd(a, b interface{}, colType table.ColType) (interface{}, error) {
	switch colType {
	case table.ColTypeFloat:
		af, aok := utils.ToFloat64(a)
		bf, bok := utils.ToFloat64(b)
		if aok && bok {
			return af + bf, nil
		}
		return nil, fmt.Errorf("applyAdd: cannot convert to float64")
	case table.ColTypeInt:
		ai, aok := utils.ToInt(a)
		bi, bok := utils.ToInt(b)
		if aok && bok {
			return ai + bi, nil
		}
		return nil, fmt.Errorf("applyAdd: cannot convert to int")
	default:
		panic(fmt.Sprintf("unexpected table.ColType: %#v", colType))
	}
}

func applySubtract(a, b interface{}, colType table.ColType) (interface{}, error) {
	switch colType {
	case table.ColTypeFloat:
		af, aok := utils.ToFloat64(a)
		bf, bok := utils.ToFloat64(b)
		if aok && bok {
			return af - bf, nil
		}
		return nil, fmt.Errorf("applySubtract: cannot convert to float64")
	case table.ColTypeInt:
		ai, aok := utils.ToInt(a)
		bi, bok := utils.ToInt(b)
		if aok && bok {
			return ai - bi, nil
		}
		return nil, fmt.Errorf("applySubtract: cannot convert to int")
	default:
		panic(fmt.Sprintf("unexpected table.ColType: %#v", colType))
	}
}

func applyMultiply(a, b interface{}, colType table.ColType) (interface{}, error) {
	switch colType {
	case table.ColTypeFloat:
		af, aok := utils.ToFloat64(a)
		bf, bok := utils.ToFloat64(b)
		if aok && bok {
			return af * bf, nil
		}
		return nil, fmt.Errorf("applyMultiply: cannot convert to float64")
	case table.ColTypeInt:
		ai, aok := utils.ToInt(a)
		bi, bok := utils.ToInt(b)
		if aok && bok {
			return ai * bi, nil
		}
		return nil, fmt.Errorf("applyMultiply: cannot convert to int")
	default:
		panic(fmt.Sprintf("unexpected table.ColType: %#v", colType))
	}
}

func applyDivide(a, b interface{}, colType table.ColType) (interface{}, error) {
	switch colType {
	case table.ColTypeFloat:
		af, aok := utils.ToFloat64(a)
		bf, bok := utils.ToFloat64(b)
		if aok && bok {
			if bf == 0 {
				return nil, fmt.Errorf("applyDivide: division by zero")
			}
			return af / bf, nil
		}
		return nil, fmt.Errorf("applyDivide: cannot convert to float64")
	case table.ColTypeInt:
		ai, aok := utils.ToInt(a)
		bi, bok := utils.ToInt(b)
		if aok && bok {
			if bi == 0 {
				return nil, fmt.Errorf("applyDivide: division by zero")
			}
			return ai / bi, nil
		}
		return nil, fmt.Errorf("applyDivide: cannot convert to int")
	default:
		panic(fmt.Sprintf("unexpected table.ColType: %#v", colType))
	}
}
