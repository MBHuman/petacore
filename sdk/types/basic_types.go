package ptypes

import (
	"bytes"
	"fmt"
	"petacore/sdk/pmem"
	"reflect"
	"regexp"
)

type BaseType[T any] interface {
	GetType() OID
	Compare(other BaseType[T]) int
	GetBuffer() []byte
	IntoGo() T
}

func TryIntoBool(val BaseType[any]) (BaseType[bool], bool) {
	w, ok := val.(AnyWrapper[bool])
	if !ok {
		return nil, false
	}
	inner, ok := w.Inner().(BaseType[bool])
	return inner, ok
}

func TryIntoNumeric[T any](val BaseType[any]) (NumericType[T], bool) {
	w, ok := val.(AnyWrapper[T])
	if !ok {
		return nil, false
	}
	n, ok := w.Inner().(NumericType[T])
	return n, ok
}

func TryIntoOrdered[T any](val BaseType[any]) (OrderedType[T], bool) {
	w, ok := val.(AnyWrapper[T])
	if !ok {
		return nil, false
	}
	o, ok := w.Inner().(OrderedType[T])
	return o, ok
}

func TryIntoBitwise[T any](val BaseType[any]) (BitwiseType[T], bool) {
	w, ok := val.(AnyWrapper[T])
	if !ok {
		return nil, false
	}
	b, ok := w.Inner().(BitwiseType[T])
	return b, ok
}

func TryIntoText[T any](val BaseType[any]) (TextType[T], bool) {
	w, ok := val.(AnyWrapper[T])
	if !ok {
		return nil, false
	}
	t, ok := w.Inner().(TextType[T])
	return t, ok
}

// NumericType — Add/Sub/Mul/Div/Mod создают новое значение — нужен аллокатор
// Neg/Abs тоже создают новое значение
// IsZero — только читает, аллокатор не нужен
type NumericType[T any] interface {
	BaseType[T]

	Add(allocator pmem.Allocator, other NumericType[T]) (NumericType[T], error)
	Sub(allocator pmem.Allocator, other NumericType[T]) (NumericType[T], error)
	Mul(allocator pmem.Allocator, other NumericType[T]) (NumericType[T], error)
	Div(allocator pmem.Allocator, other NumericType[T]) (NumericType[T], error)
	Mod(allocator pmem.Allocator, other NumericType[T]) (NumericType[T], error)

	IsZero() bool
	Neg(allocator pmem.Allocator) NumericType[T]
	Abs(allocator pmem.Allocator) NumericType[T]
}

// OrderedType — только сравнения, ничего не создаёт — аллокатор не нужен
type OrderedType[T any] interface {
	BaseType[T]

	LessThan(other BaseType[T]) bool
	GreaterThan(other BaseType[T]) bool
	LessOrEqual(other BaseType[T]) bool
	GreaterOrEqual(other BaseType[T]) bool
	Between(low, high BaseType[T]) bool
}

// BitwiseType — все операции создают новое значение — нужен аллокатор
type BitwiseType[T any] interface {
	BaseType[T]

	And(allocator pmem.Allocator, other BitwiseType[T]) BitwiseType[T]
	Or(allocator pmem.Allocator, other BitwiseType[T]) BitwiseType[T]
	Xor(allocator pmem.Allocator, other BitwiseType[T]) BitwiseType[T]
	Not(allocator pmem.Allocator) BitwiseType[T]
	ShiftLeft(allocator pmem.Allocator, n uint) BitwiseType[T]
	ShiftRight(allocator pmem.Allocator, n uint) BitwiseType[T]
}

// TextType — все операции создающие новую строку требуют аллокатор
// Like/ILike/StartsWith/Contains/Length — только читают, аллокатор не нужен
type TextType[T any] interface {
	BaseType[T]

	Concat(allocator pmem.Allocator, other TextType[T]) (TextType[T], error)
	Length() int
	Like(pattern string) bool
	ILike(pattern string) bool
	StartsWith(prefix string) bool
	Contains(substr string) bool
	ToUpper(allocator pmem.Allocator) (TextType[T], error)
	ToLower(allocator pmem.Allocator) (TextType[T], error)
	Trim(allocator pmem.Allocator) (TextType[T], error)
	Substring(allocator pmem.Allocator, start, length int) (TextType[T], error)
	AsStr() string

	RegexpMatch(pattern string) (bool, error)
	RegexpMatchCompiled(re *regexp.Regexp) bool
	Replace(allocator pmem.Allocator, old, new string) (TextType[T], error)
	RegexpReplace(allocator pmem.Allocator, pattern, replacement string) (TextType[T], error)
	RegexpReplaceCompiled(allocator pmem.Allocator, re *regexp.Regexp, replacement string) (TextType[T], error)
}

// NullableType — только читает, аллокатор не нужен
type NullableType[T any] interface {
	BaseType[T]

	IsNull() bool
	IsNotNull() bool
}

type CastableType[T any] interface {
	BaseType[T]

	CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error)
}

func ApplyNeg(allocator pmem.Allocator, val BaseType[any], oid OID) (BaseType[any], error) {
	switch oid {
	case PTypeInt2:
		n, ok := TryIntoNumeric[int16](val)
		if !ok {
			return nil, fmt.Errorf("unary minus: failed to extract int2")
		}
		return NewAnyWrapper(n.Neg(allocator)), nil

	case PTypeInt4:
		n, ok := TryIntoNumeric[int32](val)
		if !ok {
			return nil, fmt.Errorf("unary minus: failed to extract int4")
		}
		return NewAnyWrapper(n.Neg(allocator)), nil

	case PTypeInt8:
		n, ok := TryIntoNumeric[int64](val)
		if !ok {
			return nil, fmt.Errorf("unary minus: failed to extract int8")
		}
		return NewAnyWrapper(n.Neg(allocator)), nil

	case PTypeFloat4:
		n, ok := TryIntoNumeric[float32](val)
		if !ok {
			return nil, fmt.Errorf("unary minus: failed to extract float4")
		}
		return NewAnyWrapper(n.Neg(allocator)), nil

	case PTypeFloat8:
		n, ok := TryIntoNumeric[float64](val)
		if !ok {
			return nil, fmt.Errorf("unary minus: failed to extract float8")
		}
		return NewAnyWrapper(n.Neg(allocator)), nil

	case PTypeNumeric:
		n, ok := TryIntoNumeric[[]byte](val)
		if !ok {
			return nil, fmt.Errorf("unary minus: failed to extract numeric")
		}
		return NewAnyWrapper(n.Neg(allocator)), nil

	default:
		return nil, fmt.Errorf("unary minus can only be applied to numeric values, got OID %d", oid)
	}
}

// AnyWrapper оборачивает конкретный BaseType[T] в BaseType[any]
// нужен для реализации CastableType — CastTo возвращает BaseType[any]
type AnyWrapper[T any] struct {
	inner BaseType[T]
}

func NewAnyWrapper[T any](inner BaseType[T]) AnyWrapper[T] {
	return AnyWrapper[T]{inner: inner}
}

func (w AnyWrapper[T]) GetType() OID      { return w.inner.GetType() }
func (w AnyWrapper[T]) GetBuffer() []byte { return w.inner.GetBuffer() }
func (w AnyWrapper[T]) Compare(other BaseType[any]) int {
	return bytes.Compare(w.inner.GetBuffer(), other.GetBuffer())
}
func (w AnyWrapper[T]) IntoGo() any        { return w.inner.IntoGo() }
func (w AnyWrapper[T]) Inner() BaseType[T] { return w.inner }

// CastTo delegates to the inner type's CastTo if it implements CastableType[T].
// This allows AnyWrapper[T] to satisfy CastableType[any], enabling the
// comparison code to upcast mismatched integer types (e.g. Int4 → Int8)
// before doing a byte-level compare.
func (w AnyWrapper[T]) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	if castable, ok := w.inner.(CastableType[T]); ok {
		return castable.CastTo(allocator, targetType)
	}
	return nil, fmt.Errorf("type %T does not support casting to OID %d", w.inner, targetType)
}

// NumericType[any] methods - forward to inner if it's numeric
// These allow AnyWrapper[T] to act as NumericType[any] when T is numeric
func (w AnyWrapper[T]) Add(allocator pmem.Allocator, other NumericType[any]) (NumericType[any], error) {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		// Try to extract other's inner value
		var otherT NumericType[T]
		if otherWrapper, ok := other.(AnyWrapper[T]); ok {
			if otherNumeric, ok := otherWrapper.inner.(NumericType[T]); ok {
				otherT = otherNumeric
			} else {
				return nil, fmt.Errorf("incompatible numeric types for Add")
			}
		} else {
			return nil, fmt.Errorf("incompatible numeric types for Add")
		}
		result, err := numeric.Add(allocator, otherT)
		if err != nil {
			return nil, err
		}
		return NewAnyWrapper[T](result), nil
	}
	return nil, fmt.Errorf("inner type is not numeric")
}

func (w AnyWrapper[T]) Sub(allocator pmem.Allocator, other NumericType[any]) (NumericType[any], error) {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		var otherT NumericType[T]
		if otherWrapper, ok := other.(AnyWrapper[T]); ok {
			if otherNumeric, ok := otherWrapper.inner.(NumericType[T]); ok {
				otherT = otherNumeric
			} else {
				return nil, fmt.Errorf("incompatible numeric types for Sub")
			}
		} else {
			return nil, fmt.Errorf("incompatible numeric types for Sub")
		}
		result, err := numeric.Sub(allocator, otherT)
		if err != nil {
			return nil, err
		}
		return NewAnyWrapper[T](result), nil
	}
	return nil, fmt.Errorf("inner type is not numeric")
}

func (w AnyWrapper[T]) Mul(allocator pmem.Allocator, other NumericType[any]) (NumericType[any], error) {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		var otherT NumericType[T]
		if otherWrapper, ok := other.(AnyWrapper[T]); ok {
			if otherNumeric, ok := otherWrapper.inner.(NumericType[T]); ok {
				otherT = otherNumeric
			} else {
				return nil, fmt.Errorf("incompatible numeric types for Mul")
			}
		} else {
			return nil, fmt.Errorf("incompatible numeric types for Mul")
		}
		result, err := numeric.Mul(allocator, otherT)
		if err != nil {
			return nil, err
		}
		return NewAnyWrapper[T](result), nil
	}
	return nil, fmt.Errorf("inner type is not numeric")
}

func (w AnyWrapper[T]) Div(allocator pmem.Allocator, other NumericType[any]) (NumericType[any], error) {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		var otherT NumericType[T]
		if otherWrapper, ok := other.(AnyWrapper[T]); ok {
			if otherNumeric, ok := otherWrapper.inner.(NumericType[T]); ok {
				otherT = otherNumeric
			} else {
				return nil, fmt.Errorf("incompatible numeric types for Div")
			}
		} else {
			return nil, fmt.Errorf("incompatible numeric types for Div")
		}
		result, err := numeric.Div(allocator, otherT)
		if err != nil {
			return nil, err
		}
		return NewAnyWrapper[T](result), nil
	}
	return nil, fmt.Errorf("inner type is not numeric")
}

func (w AnyWrapper[T]) Mod(allocator pmem.Allocator, other NumericType[any]) (NumericType[any], error) {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		var otherT NumericType[T]
		if otherWrapper, ok := other.(AnyWrapper[T]); ok {
			if otherNumeric, ok := otherWrapper.inner.(NumericType[T]); ok {
				otherT = otherNumeric
			} else {
				return nil, fmt.Errorf("incompatible numeric types for Mod")
			}
		} else {
			return nil, fmt.Errorf("incompatible numeric types for Mod")
		}
		result, err := numeric.Mod(allocator, otherT)
		if err != nil {
			return nil, err
		}
		return NewAnyWrapper[T](result), nil
	}
	return nil, fmt.Errorf("inner type is not numeric")
}

func (w AnyWrapper[T]) IsZero() bool {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		return numeric.IsZero()
	}
	return false
}

func (w AnyWrapper[T]) Neg(allocator pmem.Allocator) NumericType[any] {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		result := numeric.Neg(allocator)
		return NewAnyWrapper[T](result)
	}
	return nil
}

func (w AnyWrapper[T]) Abs(allocator pmem.Allocator) NumericType[any] {
	if numeric, ok := w.inner.(NumericType[T]); ok {
		result := numeric.Abs(allocator)
		return NewAnyWrapper[T](result)
	}
	return nil
}

const (
	// OIDs for basic types
	PTypeBool       OID = 16
	PTypeBytea      OID = 17
	PTypeChar       OID = 18
	PTypeName       OID = 19
	PTypeInt8       OID = 20
	PTypeInt2       OID = 21
	PTypeInt4       OID = 23
	PTypeText       OID = 25
	PTypeFloat4     OID = 700
	PTypeFloat8     OID = 701
	PTypeVarchar    OID = 1043
	PTypeNumeric    OID = 1700
	PTypeTimestamp  OID = 1114
	PTypeTimestampz OID = 1184 // timestamp with time zone
	PTypeInterval   OID = 1186 // interval
	PTypeDate       OID = 1082
	PTypeTime       OID = 1083

	// OIDs for array types
	PTypeBoolArray    OID = 1000
	PTypeNameArray    OID = 1003
	PTypeInt2Array    OID = 1005
	PTypeInt4Array    OID = 1007
	PTypeTextArray    OID = 1009
	PTypeInt8Array    OID = 1016
	PTypeVarcharArray OID = 1015
	PTypeFloat4Array  OID = 1021
	PTypeFloat8Array  OID = 1022
)

func ColTypeFromString(typeStr string) OID {
	switch typeStr {
	case "bool":
		return PTypeBool
	case "bytea":
		return PTypeBytea
	case "char":
		return PTypeChar
	case "name":
		return PTypeName
	case "int8":
		return PTypeInt8
	case "int2":
		return PTypeInt2
	case "int4":
		return PTypeInt4
	case "text":
		return PTypeText
	case "float4":
		return PTypeFloat4
	case "float8":
		return PTypeFloat8
	case "varchar":
		return PTypeVarchar
	case "numeric":
		return PTypeNumeric
	case "timestamp":
		return PTypeTimestamp
	case "timestamptz":
		return PTypeTimestampz
	case "interval":
		return PTypeInterval
	case "date":
		return PTypeDate
	case "time":
		return PTypeTime
	default:
		return 0 // unknown type
	}
}

func TypeToGoType(oid OID) reflect.Type {
	switch oid {
	case PTypeBool:
		return reflect.TypeFor[bool]()
	case PTypeInt2:
		return reflect.TypeFor[int16]()
	case PTypeInt4:
		return reflect.TypeFor[int32]()
	case PTypeInt8:
		return reflect.TypeFor[int64]()
	case PTypeFloat4:
		return reflect.TypeFor[float32]()
	case PTypeFloat8: // IEEE-754 double
		return reflect.TypeFor[float64]()
	case PTypeText, PTypeVarchar, PTypeName:
		return reflect.TypeFor[string]()
	case PTypeChar:
		return reflect.TypeFor[rune]()
	case PTypeBytea:
		return reflect.TypeFor[[]byte]()
	case PTypeNumeric: // десятичная арифметика произвольной точности
		return reflect.TypeFor[float64]()
	case PTypeTimestamp, PTypeTimestampz:
		return reflect.TypeFor[int64]() // timestamp as int64 (microseconds since epoch)
	case PTypeInterval:
		return reflect.TypeFor[int64]() // interval as int64 (microseconds)
	case PTypeDate:
		return reflect.TypeFor[int32]() // date as int32 (days since epoch)
	// Array types
	case PTypeBoolArray:
		return reflect.TypeFor[[]bool]()
	case PTypeInt2Array:
		return reflect.TypeFor[[]int16]()
	case PTypeInt4Array:
		return reflect.TypeFor[[]int32]()
	case PTypeInt8Array:
		return reflect.TypeFor[[]int64]()
	case PTypeFloat4Array:
		return reflect.TypeFor[[]float32]()
	case PTypeFloat8Array:
		return reflect.TypeFor[[]float64]()
	case PTypeTextArray, PTypeVarcharArray, PTypeNameArray:
		return reflect.TypeFor[[]string]()
	default:
		return reflect.TypeOf(nil)
	}
}

// ToBaseTypeAny converts various types to BaseType[any]
// Handles: BaseType[any], AnyWrapper[T] - returns as BaseType[any]
// AnyWrapper[T] already implements BaseType[any] interface
func ToBaseTypeAny(val any) (BaseType[any], error) {
	// Already BaseType[any]
	if bt, ok := val.(BaseType[any]); ok {
		return bt, nil
	}

	// AnyWrapper[T] implements BaseType[any], so return as-is
	// Note: We check specific types to ensure they implement the interface
	switch v := val.(type) {
	case AnyWrapper[bool]:
		return v, nil
	case AnyWrapper[int16]:
		return v, nil
	case AnyWrapper[int32]:
		return v, nil
	case AnyWrapper[int64]:
		return v, nil
	case AnyWrapper[float32]:
		return v, nil
	case AnyWrapper[float64]:
		return v, nil
	case AnyWrapper[string]:
		return v, nil
	case AnyWrapper[[]byte]:
		return v, nil
	case AnyWrapper[any]:
		return v, nil
	default:
		return nil, nil
	}
}

// ToNumericAny converts various types to NumericType[any]
// Handles: BaseType[any], AnyWrapper[T], and checks if it's numeric
func ToNumericAny(val any) (NumericType[any], error) {
	bt, err := ToBaseTypeAny(val)
	if err != nil {
		return nil, err
	}
	if bt == nil {
		return nil, nil
	}

	numeric, ok := TryIntoNumeric[any](bt)
	if !ok {
		return nil, nil
	}
	return numeric, nil
}
