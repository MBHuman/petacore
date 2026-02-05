package baseplugin

import (
	"context"
	"errors"
	"fmt"
	psdk "petacore/sdk"
)

// CountFunction - агрегатная функция COUNT
type CountFunction struct {
	*psdk.BaseFunction
}

func (f *CountFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2147,
		ProName:     "COUNT",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeInt8},
		ProRetType:  psdk.PTypeFloat8,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *CountFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("count requires at least one argument")
	}

	// Для агрегатных функций args[0] должен быть []interface{}
	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("count argument must be a slice of values")
	}

	count := 0
	for _, val := range values {
		if val != nil {
			count++
		}
	}
	return float64(count), nil
}

// SumFunction - агрегатная функция SUM
type SumFunction struct {
	*psdk.BaseFunction
}

func (f *SumFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2108,
		ProName:     "SUM",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeFloat8},
		ProRetType:  psdk.PTypeFloat8,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *SumFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("sum requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("sum argument must be a slice of values")
	}

	var sum float64
	for _, val := range values {
		if val == nil {
			continue
		}
		switch v := val.(type) {
		case int:
			sum += float64(v)
		case int64:
			sum += float64(v)
		case float64:
			sum += v
		case float32:
			sum += float64(v)
		default:
			return nil, fmt.Errorf("sum can only operate on numeric values, got %T", v)
		}
	}
	return sum, nil
}

// AvgFunction - агрегатная функция AVG
type AvgFunction struct {
	*psdk.BaseFunction
}

func (f *AvgFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2101,
		ProName:     "AVG",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeFloat8},
		ProRetType:  psdk.PTypeFloat8,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *AvgFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("avg requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("avg argument must be a slice of values")
	}

	var sum float64
	count := 0
	for _, val := range values {
		if val == nil {
			continue
		}
		switch v := val.(type) {
		case int:
			sum += float64(v)
		case int64:
			sum += float64(v)
		case float64:
			sum += v
		case float32:
			sum += float64(v)
		default:
			return nil, fmt.Errorf("avg can only operate on numeric values, got %T", v)
		}
		count++
	}

	if count == 0 {
		return nil, nil
	}

	return sum / float64(count), nil
}

// MaxFunction - агрегатная функция MAX
type MaxFunction struct {
	*psdk.BaseFunction
}

func (f *MaxFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2116,
		ProName:     "MAX",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeFloat8},
		ProRetType:  psdk.PTypeFloat8,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *MaxFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("max requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("max argument must be a slice of values")
	}

	if len(values) == 0 {
		return nil, nil
	}

	var maxVal interface{}
	for _, val := range values {
		if val == nil {
			continue
		}

		if maxVal == nil {
			maxVal = val
			continue
		}

		switch v := val.(type) {
		case int:
			if mv, ok := maxVal.(int); ok && v > mv {
				maxVal = v
			} else if mv, ok := maxVal.(int64); ok && int64(v) > mv {
				maxVal = v
			} else if mv, ok := maxVal.(float64); ok && float64(v) > mv {
				maxVal = v
			}
		case int64:
			if mv, ok := maxVal.(int64); ok && v > mv {
				maxVal = v
			} else if mv, ok := maxVal.(int); ok && v > int64(mv) {
				maxVal = v
			} else if mv, ok := maxVal.(float64); ok && float64(v) > mv {
				maxVal = v
			}
		case float64:
			if mv, ok := maxVal.(float64); ok && v > mv {
				maxVal = v
			} else if mv, ok := maxVal.(int); ok && v > float64(mv) {
				maxVal = v
			} else if mv, ok := maxVal.(int64); ok && v > float64(mv) {
				maxVal = v
			}
		case string:
			if mv, ok := maxVal.(string); ok && v > mv {
				maxVal = v
			}
		}
	}

	return maxVal, nil
}

// MinFunction - агрегатная функция MIN
type MinFunction struct {
	*psdk.BaseFunction
}

func (f *MinFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2131,
		ProName:     "MIN",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeFloat8},
		ProRetType:  psdk.PTypeFloat8,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *MinFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("min requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("min argument must be a slice of values")
	}

	if len(values) == 0 {
		return nil, nil
	}

	var minVal interface{}
	for _, val := range values {
		if val == nil {
			continue
		}

		if minVal == nil {
			minVal = val
			continue
		}

		switch v := val.(type) {
		case int:
			if mv, ok := minVal.(int); ok && v < mv {
				minVal = v
			} else if mv, ok := minVal.(int64); ok && int64(v) < mv {
				minVal = v
			} else if mv, ok := minVal.(float64); ok && float64(v) < mv {
				minVal = v
			}
		case int64:
			if mv, ok := minVal.(int64); ok && v < mv {
				minVal = v
			} else if mv, ok := minVal.(int); ok && v < int64(mv) {
				minVal = v
			} else if mv, ok := minVal.(float64); ok && float64(v) < mv {
				minVal = v
			}
		case float64:
			if mv, ok := minVal.(float64); ok && v < mv {
				minVal = v
			} else if mv, ok := minVal.(int); ok && v < float64(mv) {
				minVal = v
			} else if mv, ok := minVal.(int64); ok && v < float64(mv) {
				minVal = v
			}
		case string:
			if mv, ok := minVal.(string); ok && v < mv {
				minVal = v
			}
		}
	}

	return minVal, nil
}

// MaxFunctionInt - MAX для INT4
type MaxFunctionInt struct {
	*psdk.BaseFunction
}

func (f *MaxFunctionInt) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2117,
		ProName:     "MAX",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeInt4},
		ProRetType:  psdk.PTypeInt4,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *MaxFunctionInt) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("max requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("max argument must be a slice of values")
	}

	if len(values) == 0 {
		return nil, nil
	}

	var maxVal interface{}
	for _, val := range values {
		if val == nil {
			continue
		}

		if maxVal == nil {
			maxVal = val
			continue
		}

		// Поддержка int, int32, int64
		switch v := val.(type) {
		case int:
			if mv, ok := maxVal.(int); ok && v > mv {
				maxVal = v
			}
		case int32:
			if mv, ok := maxVal.(int32); ok && v > mv {
				maxVal = v
			}
		case int64:
			if mv, ok := maxVal.(int64); ok && v > mv {
				maxVal = v
			}
		}
	}

	return maxVal, nil
}

// MinFunctionInt - MIN для INT4
type MinFunctionInt struct {
	*psdk.BaseFunction
}

func (f *MinFunctionInt) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2132,
		ProName:     "MIN",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeInt4},
		ProRetType:  psdk.PTypeInt4,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *MinFunctionInt) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("min requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("min argument must be a slice of values")
	}

	if len(values) == 0 {
		return nil, nil
	}

	var minVal interface{}
	for _, val := range values {
		if val == nil {
			continue
		}

		if minVal == nil {
			minVal = val
			continue
		}

		// Поддержка int, int32, int64
		switch v := val.(type) {
		case int:
			if mv, ok := minVal.(int); ok && v < mv {
				minVal = v
			}
		case int32:
			if mv, ok := minVal.(int32); ok && v < mv {
				minVal = v
			}
		case int64:
			if mv, ok := minVal.(int64); ok && v < mv {
				minVal = v
			}
		}
	}

	return minVal, nil
}

// MaxFunctionText - MAX для TEXT
type MaxFunctionText struct {
	*psdk.BaseFunction
}

func (f *MaxFunctionText) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2129,
		ProName:     "MAX",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeText},
		ProRetType:  psdk.PTypeText,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *MaxFunctionText) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("max requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("max argument must be a slice of values")
	}

	if len(values) == 0 {
		return nil, nil
	}

	var maxVal string
	first := true
	for _, val := range values {
		if val == nil {
			continue
		}

		if v, ok := val.(string); ok {
			if first || v > maxVal {
				maxVal = v
				first = false
			}
		}
	}

	if first {
		return nil, nil
	}

	return maxVal, nil
}

// MinFunctionText - MIN для TEXT
type MinFunctionText struct {
	*psdk.BaseFunction
}

func (f *MinFunctionText) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2145,
		ProName:     "MIN",
		IsAggregate: true,
		ProArgTypes: []psdk.OID{psdk.PTypeText},
		ProRetType:  psdk.PTypeText,
		Meta: psdk.FunctionMeta{
			ProVariadic: 1,
		},
	}
}

func (f *MinFunctionText) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) == 0 {
		return nil, errors.New("min requires at least one argument")
	}

	values, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("min argument must be a slice of values")
	}

	if len(values) == 0 {
		return nil, nil
	}

	var minVal string
	first := true
	for _, val := range values {
		if val == nil {
			continue
		}

		if v, ok := val.(string); ok {
			if first || v < minVal {
				minVal = v
				first = false
			}
		}
	}

	if first {
		return nil, nil
	}

	return minVal, nil
}
