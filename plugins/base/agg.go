package baseplugin

import (
	"context"
	"errors"
	"fmt"
	psdk "petacore/sdk"
	ptypes "petacore/sdk/types"
)

// toFloat64Val extracts a float64 from a value, unwrapping AnyWrapper[T] if needed.
func toFloat64Val(v interface{}) (float64, bool) {
	if v == nil {
		return 0, false
	}
	if w, ok := v.(interface{ IntoGo() any }); ok {
		v = w.IntoGo()
	}
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	}
	return 0, false
}

// toStringVal extracts a string from a value, unwrapping AnyWrapper[T] if needed.
func toStringVal(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}
	if w, ok := v.(interface{ IntoGo() any }); ok {
		v = w.IntoGo()
	}
	if s, ok := v.(string); ok {
		return s, true
	}
	return "", false
}

// CountFunction - агрегатная функция COUNT
type CountFunction struct {
	*psdk.BaseFunction
}

func (f *CountFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2147,
		ProName:     "COUNT",
		IsAggregate: true,
		ProArgTypes: []ptypes.OID{ptypes.PTypeInt8},
		ProRetType:  ptypes.PTypeInt8,
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
	return int64(count), nil
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeFloat8},
		ProRetType:  ptypes.PTypeFloat8,
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
		f64, ok := toFloat64Val(val)
		if !ok {
			return nil, fmt.Errorf("sum can only operate on numeric values, got %T", val)
		}
		sum += f64
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeFloat8},
		ProRetType:  ptypes.PTypeFloat8,
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
		f64, ok := toFloat64Val(val)
		if !ok {
			return nil, fmt.Errorf("avg can only operate on numeric values, got %T", val)
		}
		sum += f64
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeFloat8},
		ProRetType:  ptypes.PTypeFloat8,
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

	var maxFloat float64
	haveVal := false
	for _, val := range values {
		if val == nil {
			continue
		}
		f64, ok := toFloat64Val(val)
		if !ok {
			continue
		}
		if !haveVal || f64 > maxFloat {
			maxFloat = f64
			haveVal = true
		}
	}
	if !haveVal {
		return nil, nil
	}
	return maxFloat, nil
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeFloat8},
		ProRetType:  ptypes.PTypeFloat8,
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

	var minFloat float64
	haveVal := false
	for _, val := range values {
		if val == nil {
			continue
		}
		f64, ok := toFloat64Val(val)
		if !ok {
			continue
		}
		if !haveVal || f64 < minFloat {
			minFloat = f64
			haveVal = true
		}
	}
	if !haveVal {
		return nil, nil
	}
	return minFloat, nil
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeInt4},
		ProRetType:  ptypes.PTypeInt4,
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

	var maxVal int64
	haveVal := false
	for _, val := range values {
		if val == nil {
			continue
		}
		f64, ok := toFloat64Val(val)
		if !ok {
			continue
		}
		v := int64(f64)
		if !haveVal || v > maxVal {
			maxVal = v
			haveVal = true
		}
	}
	if !haveVal {
		return nil, nil
	}
	return int32(maxVal), nil
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeInt4},
		ProRetType:  ptypes.PTypeInt4,
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

	var minVal int64
	haveVal := false
	for _, val := range values {
		if val == nil {
			continue
		}
		f64, ok := toFloat64Val(val)
		if !ok {
			continue
		}
		v := int64(f64)
		if !haveVal || v < minVal {
			minVal = v
			haveVal = true
		}
	}
	if !haveVal {
		return nil, nil
	}
	return int32(minVal), nil
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeText},
		ProRetType:  ptypes.PTypeText,
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

	var maxStr string
	first := true
	for _, val := range values {
		if val == nil {
			continue
		}
		v, ok := toStringVal(val)
		if !ok {
			continue
		}
		if first || v > maxStr {
			maxStr = v
			first = false
		}
	}

	if first {
		return nil, nil
	}

	return maxStr, nil
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
		ProArgTypes: []ptypes.OID{ptypes.PTypeText},
		ProRetType:  ptypes.PTypeText,
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

	var minStr string
	first := true
	for _, val := range values {
		if val == nil {
			continue
		}
		v, ok := toStringVal(val)
		if !ok {
			continue
		}
		if first || v < minStr {
			minStr = v
			first = false
		}
	}

	if first {
		return nil, nil
	}

	return minStr, nil
}
