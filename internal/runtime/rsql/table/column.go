package table

import (
	"fmt"
	"petacore/internal/utils"
)

type ColType int

const (
	ColTypeString ColType = iota
	ColTypeInt
	ColTypeBigInt
	ColTypeFloat
	ColTypeBool
	ColTypeTimestamp
	ColTypeTimestampTz
	ColTypeInterval
	ColTypeDate
	// Array types
	ColTypeStringArray
	ColTypeIntArray
	ColTypeBigIntArray
	ColTypeFloatArray
	ColTypeBoolArray
)

func (c ColType) String() string {
	switch c {
	case ColTypeString:
		return "text"
	case ColTypeInt:
		return "integer"
	case ColTypeBigInt:
		return "bigint"
	case ColTypeFloat:
		return "real"
	case ColTypeBool:
		return "boolean"
	case ColTypeTimestamp:
		return "timestamp"
	case ColTypeTimestampTz:
		return "timestamp with time zone"
	case ColTypeInterval:
		return "interval"
	case ColTypeDate:
		return "date"
	case ColTypeStringArray:
		return "text[]"
	case ColTypeIntArray:
		return "integer[]"
	case ColTypeBigIntArray:
		return "bigint[]"
	case ColTypeFloatArray:
		return "real[]"
	case ColTypeBoolArray:
		return "boolean[]"
	default:
		return "text"
	}
}

func (c ColType) TypeOps() TypeOps {
	switch c {
	case ColTypeString:
		return &StringOps{}
	case ColTypeInt:
		return &IntOps{}
	case ColTypeBigInt:
		return &IntOps{} // BigInt uses same ops as Int
	case ColTypeFloat:
		return &FloatOps{}
	case ColTypeBool:
		return &BoolOps{}
	case ColTypeTimestamp, ColTypeTimestampTz:
		return &IntOps{} // Treat timestamp as int64
	case ColTypeInterval:
		return &IntOps{} // Treat interval as int64 (microseconds)
	case ColTypeStringArray:
		return &ArrayOps{ElementType: ColTypeString}
	case ColTypeIntArray:
		return &ArrayOps{ElementType: ColTypeInt}
	case ColTypeBigIntArray:
		return &ArrayOps{ElementType: ColTypeBigInt}
	case ColTypeFloatArray:
		return &ArrayOps{ElementType: ColTypeFloat}
	case ColTypeBoolArray:
		return &ArrayOps{ElementType: ColTypeBool}
	default:
		panic("unexpected ColType")
	}
}

type TypeOps interface {
	CastTo(value interface{}, targetType ColType) (interface{}, error)
	Compare(a, b interface{}, rightTyp ColType) (int, error)
}

type StringOps struct{}

func (s *StringOps) CastTo(value interface{}, targetType ColType) (interface{}, error) {
	switch targetType {
	case ColTypeString:
		return fmt.Sprintf("%v", value), nil
	case ColTypeInt:
		if i, ok := utils.ToInt(value); ok {
			return i, nil
		}
		return 0, fmt.Errorf("cannot cast to int")
	case ColTypeFloat:
		if f, ok := utils.ToFloat64(value); ok {
			return f, nil
		}
		return 0.0, fmt.Errorf("cannot cast to float")
	case ColTypeBool:
		if b, ok := utils.ToBool(value); ok {
			return b, nil
		}
		return false, fmt.Errorf("cannot cast to bool")
	default:
		panic("unexpected target ColType")
	}
}

func (s *StringOps) Compare(a, b interface{}, rightTyp ColType) (int, error) {
	switch rightTyp {
	case ColTypeString:
		as := fmt.Sprintf("%v", a)
		bs := fmt.Sprintf("%v", b)
		if as < bs {
			return -1, nil
		} else if as > bs {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot compare string with %v", rightTyp)
	}
}

type IntOps struct{}

func (i *IntOps) CastTo(value interface{}, targetType ColType) (interface{}, error) {
	switch targetType {
	case ColTypeString:
		return fmt.Sprintf("%v", value), nil
	case ColTypeInt:
		if iv, ok := utils.ToInt(value); ok {
			return iv, nil
		}
		return 0, fmt.Errorf("cannot cast to int")
	case ColTypeBigInt:
		if iv, ok := utils.ToInt(value); ok {
			return int64(iv), nil
		}
		return int64(0), fmt.Errorf("cannot cast to bigint")
	case ColTypeFloat:
		if f, ok := utils.ToFloat64(value); ok {
			return f, nil
		}
		return 0.0, fmt.Errorf("cannot cast to float")
	case ColTypeBool:
		if b, ok := utils.ToBool(value); ok {
			return b, nil
		}
		return false, fmt.Errorf("cannot cast to bool")
	case ColTypeTimestamp, ColTypeTimestampTz, ColTypeInterval, ColTypeDate:
		if iv, ok := utils.ToInt(value); ok {
			return iv, nil
		}
		return 0, fmt.Errorf("cannot cast to timestamp/interval/date")
	default:
		panic("unexpected target ColType")
	}
}

func (i *IntOps) Compare(a, b interface{}, rightTyp ColType) (int, error) {
	switch rightTyp {
	case ColTypeInt, ColTypeBigInt, ColTypeTimestamp, ColTypeTimestampTz, ColTypeInterval, ColTypeDate:
		// All these types are stored as int64, so we can compare them directly
		ai64, aok := utils.ToInt64(a)
		bi64, bok := utils.ToInt64(b)
		if aok && bok {
			if ai64 < bi64 {
				return -1, nil
			} else if ai64 > bi64 {
				return 1, nil
			}
			return 0, nil
		}
		// Fallback to int comparison for smaller values
		ai, aok := utils.ToInt(a)
		bi, bok := utils.ToInt(b)
		if aok && bok {
			if ai < bi {
				return -1, nil
			} else if ai > bi {
				return 1, nil
			}
			return 0, nil
		}
		return 0, fmt.Errorf("cannot convert to int for comparison")
	case ColTypeFloat:
		af, aok := utils.ToFloat64(a)
		bf, bok := utils.ToFloat64(b)
		if aok && bok {
			if af < bf {
				return -1, nil
			} else if af > bf {
				return 1, nil
			}
			return 0, nil
		}
		return 0, fmt.Errorf("cannot convert to float64 for comparison")
	default:
		return 0, fmt.Errorf("cannot compare int with %v", rightTyp)
	}
}

type FloatOps struct{}

func (f *FloatOps) CastTo(value interface{}, targetType ColType) (interface{}, error) {
	switch targetType {
	case ColTypeString:
		return fmt.Sprintf("%v", value), nil
	case ColTypeInt:
		if iv, ok := utils.ToInt(value); ok {
			return iv, nil
		}
		return 0, fmt.Errorf("cannot cast to int")
	case ColTypeFloat:
		if fv, ok := utils.ToFloat64(value); ok {
			return fv, nil
		}
		return 0.0, fmt.Errorf("cannot cast to float")
	case ColTypeBool:
		if b, ok := utils.ToBool(value); ok {
			return b, nil
		}
		return false, fmt.Errorf("cannot cast to bool")
	default:
		panic("unexpected target ColType")
	}
}

func (f *FloatOps) Compare(a, b interface{}, rightTyp ColType) (int, error) {
	switch rightTyp {
	case ColTypeFloat:
		af, aok := utils.ToFloat64(a)
		bf, bok := utils.ToFloat64(b)
		if aok && bok {
			if af < bf {
				return -1, nil
			} else if af > bf {
				return 1, nil
			}
			return 0, nil
		}
		return 0, fmt.Errorf("cannot convert to float64 for comparison")
	case ColTypeInt:
		af, aok := utils.ToFloat64(a)
		bi, bok := utils.ToInt(b)
		if aok && bok {
			bf := float64(bi)
			if af < bf {
				return -1, nil
			} else if af > bf {
				return 1, nil
			}
			return 0, nil
		}
		return 0, fmt.Errorf("cannot convert to int for comparison")
	default:
		return 0, fmt.Errorf("cannot compare float with %v", rightTyp)
	}
}

type BoolOps struct{}

func (b *BoolOps) CastTo(value interface{}, targetType ColType) (interface{}, error) {
	switch targetType {
	case ColTypeString:
		return fmt.Sprintf("%v", value), nil
	case ColTypeInt:
		if iv, ok := utils.ToInt(value); ok {
			return iv, nil
		}
		return 0, fmt.Errorf("cannot cast to int")
	case ColTypeFloat:
		if fv, ok := utils.ToFloat64(value); ok {
			return fv, nil
		}
		return 0.0, fmt.Errorf("cannot cast to float")
	case ColTypeBool:
		if bv, ok := utils.ToBool(value); ok {
			return bv, nil
		}
		return false, fmt.Errorf("cannot cast to bool")
	default:
		panic("unexpected target ColType")
	}
}

func (b *BoolOps) Compare(a, bval interface{}, rightTyp ColType) (int, error) {
	switch rightTyp {
	case ColTypeBool:
		ab, aok := utils.ToBool(a)
		bb, bok := utils.ToBool(bval)
		if aok && bok {
			if !ab && bb {
				return -1, nil
			} else if ab && !bb {
				return 1, nil
			}
			return 0, nil
		}
		return 0, fmt.Errorf("cannot convert to bool for comparison")
	default:
		return 0, fmt.Errorf("cannot compare bool with %v", rightTyp)
	}
}

// ArrayOps handles array type operations
type ArrayOps struct {
	ElementType ColType
}

func (ao *ArrayOps) CastTo(value interface{}, targetType ColType) (interface{}, error) {
	// Arrays can be cast to string representation
	if targetType == ColTypeString {
		return fmt.Sprintf("%v", value), nil
	}
	// For now, we don't support casting arrays to other types
	return nil, fmt.Errorf("cannot cast array to %v", targetType)
}

func (ao *ArrayOps) Compare(a, b interface{}, rightTyp ColType) (int, error) {
	// Array comparison is not commonly needed, but we provide basic support
	return 0, fmt.Errorf("array comparison not implemented")
}

type TableColumn struct {
	Idx               int
	Name              string
	Type              ColType
	TableIdentifier   string // table alias or name (used in query)
	OriginalTableName string // original table name (for error messages)
}

func ColTypeFromString(typeStr string) ColType {
	switch typeStr {
	case "text", "varchar", "character varying", "name":
		return ColTypeString
	case "int", "int4", "integer":
		return ColTypeInt
	case "bigint", "int8":
		return ColTypeBigInt
	case "float", "float8", "double precision":
		return ColTypeFloat
	case "bool", "boolean":
		return ColTypeBool
	case "timestamp", "timestamp without time zone":
		return ColTypeTimestamp
	case "timestamptz", "timestamp with time zone":
		return ColTypeTimestampTz
	case "interval":
		return ColTypeInterval
	case "text[]", "varchar[]", "name[]":
		return ColTypeStringArray
	case "int[]", "int4[]", "integer[]":
		return ColTypeIntArray
	case "bigint[]", "int8[]":
		return ColTypeBigIntArray
	case "float[]", "float8[]", "double precision[]":
		return ColTypeFloatArray
	case "bool[]", "boolean[]":
		return ColTypeBoolArray
	default:
		return ColTypeString
	}
}

// ColumnDef определяет колонку
type ColumnDef struct {
	Idx  int
	Name string
	Type ColType
	// IsPrimaryKey bool
	IsNullable   bool
	IsUnique     bool
	IsSerial     bool
	DefaultValue interface{}
}
