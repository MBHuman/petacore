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
	aType := a.Row.Columns[0].Type
	bType := b.Row.Columns[0].Type

	// Handle timestamp + interval
	if (aType == table.ColTypeTimestamp || aType == table.ColTypeTimestampTz) && bType == table.ColTypeInterval {
		// timestamp + interval = timestamp
		tsMicros, tsOk := utils.ToInt64(aVal)
		intervalMicros, intOk := utils.ToInt64(bVal)
		if tsOk && intOk {
			resultRow := &table.ExecuteResult{
				Rows: [][]interface{}{{tsMicros + intervalMicros}},
				Columns: []table.TableColumn{
					{Idx: 0, Name: "?column?", Type: aType}, // preserve timestamp type
				},
			}
			return &rmodels.ResultRowsExpression{Row: resultRow}, nil
		}
		return nil, fmt.Errorf("cannot add timestamp and interval")
	}

	// Handle interval + timestamp
	if aType == table.ColTypeInterval && (bType == table.ColTypeTimestamp || bType == table.ColTypeTimestampTz) {
		// interval + timestamp = timestamp
		intervalMicros, intOk := utils.ToInt64(aVal)
		tsMicros, tsOk := utils.ToInt64(bVal)
		if intOk && tsOk {
			resultRow := &table.ExecuteResult{
				Rows: [][]interface{}{{tsMicros + intervalMicros}},
				Columns: []table.TableColumn{
					{Idx: 0, Name: "?column?", Type: bType}, // preserve timestamp type
				},
			}
			return &rmodels.ResultRowsExpression{Row: resultRow}, nil
		}
		return nil, fmt.Errorf("cannot add interval and timestamp")
	}

	// Если типы не совпадают, пробуем привести к типу левой части
	if aType != bType {
		aOps := aType.TypeOps()
		convertedValue, err := aOps.CastTo(bVal, aType)
		if err == nil {
			bVal = convertedValue
			bType = aType
		} else {
			// Если не получилось привести к типу левой, пробуем привести левую к типу правой
			bOps := bType.TypeOps()
			convertedValue, err := bOps.CastTo(aVal, bType)
			if err == nil {
				aVal = convertedValue
				aType = bType
			} else {
				return nil, fmt.Errorf("AddValues: cannot cast types %s and %s", aType.String(), bType.String())
			}
		}
	}

	val, err := applyAdd(aVal, bVal, aType)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: aType},
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
	aType := a.Row.Columns[0].Type
	bType := b.Row.Columns[0].Type

	// Handle timestamp - interval
	if (aType == table.ColTypeTimestamp || aType == table.ColTypeTimestampTz) && bType == table.ColTypeInterval {
		// timestamp - interval = timestamp
		tsMicros, tsOk := utils.ToInt64(aVal)
		intervalMicros, intOk := utils.ToInt64(bVal)
		if tsOk && intOk {
			resultRow := &table.ExecuteResult{
				Rows: [][]interface{}{{tsMicros - intervalMicros}},
				Columns: []table.TableColumn{
					{Idx: 0, Name: "?column?", Type: aType}, // preserve timestamp type
				},
			}
			return &rmodels.ResultRowsExpression{Row: resultRow}, nil
		}
		return nil, fmt.Errorf("cannot subtract interval from timestamp")
	}

	// Если типы не совпадают, пробуем привести к типу левой части
	if aType != bType {
		aOps := aType.TypeOps()
		convertedValue, err := aOps.CastTo(bVal, aType)
		if err == nil {
			bVal = convertedValue
			bType = aType
		} else {
			// Если не получилось привести к типу левой, пробуем привести левую к типу правой
			bOps := bType.TypeOps()
			convertedValue, err := bOps.CastTo(aVal, bType)
			if err == nil {
				aVal = convertedValue
				aType = bType
			} else {
				return nil, fmt.Errorf("SubtractValues: cannot cast types %s and %s", aType.String(), bType.String())
			}
		}
	}

	val, err := applySubtract(aVal, bVal, aType)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: aType},
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
	aType := a.Row.Columns[0].Type
	bType := b.Row.Columns[0].Type

	// Если типы не совпадают, пробуем привести к типу левой части
	if aType != bType {
		aOps := aType.TypeOps()
		convertedValue, err := aOps.CastTo(bVal, aType)
		if err == nil {
			bVal = convertedValue
			bType = aType
		} else {
			// Если не получилось привести к типу левой, пробуем привести левую к типу правой
			bOps := bType.TypeOps()
			convertedValue, err := bOps.CastTo(aVal, bType)
			if err == nil {
				aVal = convertedValue
				aType = bType
			} else {
				return nil, fmt.Errorf("MultiplyValues: cannot cast types %s and %s", aType.String(), bType.String())
			}
		}
	}

	val, err := applyMultiply(aVal, bVal, aType)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: aType},
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
	aType := a.Row.Columns[0].Type
	bType := b.Row.Columns[0].Type

	// Если типы не совпадают, пробуем привести к типу левой части
	if aType != bType {
		aOps := aType.TypeOps()
		convertedValue, err := aOps.CastTo(bVal, aType)
		if err == nil {
			bVal = convertedValue
			bType = aType
		} else {
			// Если не получилось привести к типу левой, пробуем привести левую к типу правой
			bOps := bType.TypeOps()
			convertedValue, err := bOps.CastTo(aVal, bType)
			if err == nil {
				aVal = convertedValue
				aType = bType
			} else {
				return nil, fmt.Errorf("DivideValues: cannot cast types %s and %s", aType.String(), bType.String())
			}
		}
	}

	val, err := applyDivide(aVal, bVal, aType)
	if err != nil {
		return nil, err
	}

	resultRow := &table.ExecuteResult{
		Rows: [][]interface{}{{val}},
		Columns: []table.TableColumn{
			{Idx: 0, Name: "?column?", Type: aType},
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

	aType := a.Row.Columns[0].Type
	bType := b.Row.Columns[0].Type

	// Allow timestamp +/- interval operations
	if (aType == table.ColTypeTimestamp || aType == table.ColTypeTimestampTz) && bType == table.ColTypeInterval {
		return nil
	}
	if aType == table.ColTypeInterval && (bType == table.ColTypeTimestamp || bType == table.ColTypeTimestampTz) {
		return nil
	}

	// Проверяем, что хотя бы один из типов поддерживает арифметику
	if !(aType == table.ColTypeInt || aType == table.ColTypeFloat ||
		bType == table.ColTypeInt || bType == table.ColTypeFloat) {
		return fmt.Errorf("AddValues: types %s and %s not supported for arithmetic", aType.String(), bType.String())
	}

	// Проверяем, что типы можно привести друг к другу
	if aType != bType {
		aOps := aType.TypeOps()
		_, errA := aOps.CastTo(bRows[0][0], aType)
		if errA != nil {
			bOps := bType.TypeOps()
			_, errB := bOps.CastTo(aRows[0][0], bType)
			if errB != nil {
				return fmt.Errorf("AddValues: cannot cast between types %s and %s", aType.String(), bType.String())
			}
		}
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
	case table.ColTypeBigInt, table.ColTypeTimestamp, table.ColTypeTimestampTz, table.ColTypeInterval:
		ai, aok := utils.ToInt64(a)
		bi, bok := utils.ToInt64(b)
		if aok && bok {
			return ai + bi, nil
		}
		return nil, fmt.Errorf("applyAdd: cannot convert to int64")
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
	case table.ColTypeBigInt, table.ColTypeTimestamp, table.ColTypeTimestampTz, table.ColTypeInterval:
		ai, aok := utils.ToInt64(a)
		bi, bok := utils.ToInt64(b)
		if aok && bok {
			return ai - bi, nil
		}
		return nil, fmt.Errorf("applySubtract: cannot convert to int64")
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
