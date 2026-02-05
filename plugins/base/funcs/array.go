package basefuncs

import (
	"context"
	"fmt"
	psdk "petacore/sdk"
	"strings"
)

// ArrayToStringFunction implements PostgreSQL array_to_string(anyarray, text) -> text
// Concatenates array elements using supplied delimiter
type ArrayToStringFunction struct {
	*psdk.BaseFunction
}

func (f *ArrayToStringFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         395,
		ProName:     "ARRAY_TO_STRING",
		ProArgTypes: []psdk.OID{psdk.PTypeTextArray, psdk.PTypeText},
		ProRetType:  psdk.PTypeText,
	}
}

func (f *ArrayToStringFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("ARRAY_TO_STRING requires at least 2 arguments")
	}

	delimiter, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ARRAY_TO_STRING delimiter must be text")
	}

	// Handle different array types
	switch arr := args[0].(type) {
	case []string:
		return strings.Join(arr, delimiter), nil
	case []int:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(strs, delimiter), nil
	case []int32:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(strs, delimiter), nil
	case []int64:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(strs, delimiter), nil
	case []float32:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%g", v)
		}
		return strings.Join(strs, delimiter), nil
	case []float64:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%g", v)
		}
		return strings.Join(strs, delimiter), nil
	case []bool:
		strs := make([]string, len(arr))
		for i, v := range arr {
			if v {
				strs[i] = "t"
			} else {
				strs[i] = "f"
			}
		}
		return strings.Join(strs, delimiter), nil
	case []interface{}:
		// Generic slice handling
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%v", v)
		}
		return strings.Join(strs, delimiter), nil
	default:
		return nil, fmt.Errorf("ARRAY_TO_STRING: unsupported array type %T", args[0])
	}
}

// ArrayToStringIntFunction for integer arrays
type ArrayToStringIntFunction struct {
	*psdk.BaseFunction
}

func (f *ArrayToStringIntFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         376,
		ProName:     "ARRAY_TO_STRING",
		ProArgTypes: []psdk.OID{psdk.PTypeInt4Array, psdk.PTypeText},
		ProRetType:  psdk.PTypeText,
	}
}

func (f *ArrayToStringIntFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("ARRAY_TO_STRING requires at least 2 arguments")
	}

	delimiter, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("ARRAY_TO_STRING delimiter must be text")
	}

	switch arr := args[0].(type) {
	case []int:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(strs, delimiter), nil
	case []int32:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(strs, delimiter), nil
	case []int64:
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = fmt.Sprintf("%d", v)
		}
		return strings.Join(strs, delimiter), nil
	default:
		return nil, fmt.Errorf("ARRAY_TO_STRING: unsupported array type %T", args[0])
	}
}

// ArrayLengthFunction implements PostgreSQL array_length(anyarray, int) -> int
// Returns the length of the requested array dimension
type ArrayLengthFunction struct {
	*psdk.BaseFunction
}

func (f *ArrayLengthFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2176,
		ProName:     "ARRAY_LENGTH",
		ProArgTypes: []psdk.OID{psdk.PTypeTextArray, psdk.PTypeInt4},
		ProRetType:  psdk.PTypeInt4,
	}
}

func (f *ArrayLengthFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ARRAY_LENGTH requires exactly 2 arguments")
	}

	dimension, ok := args[1].(int32)
	if !ok {
		// Try int
		if dimInt, ok := args[1].(int); ok {
			dimension = int32(dimInt)
		} else {
			return nil, fmt.Errorf("ARRAY_LENGTH dimension must be integer")
		}
	}

	// PostgreSQL arrays are 1-indexed, and we only support 1-dimensional arrays
	if dimension != 1 {
		return nil, nil // Return NULL for invalid dimensions
	}

	// Handle different array types
	switch arr := args[0].(type) {
	case []string:
		return int32(len(arr)), nil
	case []int:
		return int32(len(arr)), nil
	case []int32:
		return int32(len(arr)), nil
	case []int64:
		return int32(len(arr)), nil
	case []float32:
		return int32(len(arr)), nil
	case []float64:
		return int32(len(arr)), nil
	case []bool:
		return int32(len(arr)), nil
	case []interface{}:
		return int32(len(arr)), nil
	default:
		return nil, fmt.Errorf("ARRAY_LENGTH: unsupported array type %T", args[0])
	}
}

// ArrayLengthIntFunction for integer arrays
type ArrayLengthIntFunction struct {
	*psdk.BaseFunction
}

func (f *ArrayLengthIntFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2177,
		ProName:     "ARRAY_LENGTH",
		ProArgTypes: []psdk.OID{psdk.PTypeInt4Array, psdk.PTypeInt4},
		ProRetType:  psdk.PTypeInt4,
	}
}

func (f *ArrayLengthIntFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ARRAY_LENGTH requires exactly 2 arguments")
	}

	dimension, ok := args[1].(int32)
	if !ok {
		if dimInt, ok := args[1].(int); ok {
			dimension = int32(dimInt)
		} else {
			return nil, fmt.Errorf("ARRAY_LENGTH dimension must be integer")
		}
	}

	if dimension != 1 {
		return nil, nil
	}

	switch arr := args[0].(type) {
	case []int:
		return int32(len(arr)), nil
	case []int32:
		return int32(len(arr)), nil
	case []int64:
		return int32(len(arr)), nil
	default:
		return nil, fmt.Errorf("ARRAY_LENGTH: unsupported array type %T", args[0])
	}
}
