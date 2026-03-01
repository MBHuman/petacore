package ptypes

import (
	"bytes"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"reflect"
)

type BaseType[T any] interface {
	GetType() OID
	Compare(other BaseType[T]) int
	GetBuffer() []byte
	IntoGo() T
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

// anyWrapper оборачивает конкретный BaseType[T] в BaseType[any]
// нужен для реализации CastableType — CastTo возвращает BaseType[any]
type anyWrapper[T any] struct {
	inner BaseType[T]
}

func (w anyWrapper[T]) GetType() OID      { return w.inner.GetType() }
func (w anyWrapper[T]) GetBuffer() []byte { return w.inner.GetBuffer() }
func (w anyWrapper[T]) Compare(other BaseType[any]) int {
	return bytes.Compare(w.inner.GetBuffer(), other.GetBuffer())
}
func (w anyWrapper[T]) IntoGo() any { return w.inner.IntoGo() }

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

func FromColType(colType table.ColType) OID {
	switch colType {
	case table.ColTypeString:
		return PTypeText
	case table.ColTypeInt:
		return PTypeInt4
	case table.ColTypeBigInt:
		return PTypeInt8
	case table.ColTypeFloat:
		return PTypeFloat8
	case table.ColTypeBool:
		return PTypeBool
	case table.ColTypeTimestamp:
		return PTypeTimestamp
	case table.ColTypeTimestampTz:
		return PTypeTimestampz
	case table.ColTypeInterval:
		return PTypeInterval
	case table.ColTypeStringArray:
		return PTypeTextArray
	case table.ColTypeIntArray:
		return PTypeInt4Array
	case table.ColTypeBigIntArray:
		return PTypeInt8Array
	case table.ColTypeFloatArray:
		return PTypeFloat8Array
	case table.ColTypeBoolArray:
		return PTypeBoolArray
	case table.ColTypeDate:
		return PTypeDate
	default:
		return PTypeText
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

// For backward compatibility

func (oid OID) ToColType() table.ColType {
	switch oid {
	case PTypeBool:
		return table.ColTypeBool
	case PTypeInt2, PTypeInt4, PTypeInt8:
		return table.ColTypeInt
	case PTypeFloat4, PTypeFloat8, PTypeNumeric:
		return table.ColTypeFloat
	case PTypeText, PTypeVarchar, PTypeChar, PTypeName:
		return table.ColTypeString
	case PTypeBytea:
		return table.ColTypeString // bytea можно представить как строку в кодировке base64
	case PTypeTimestamp:
		return table.ColTypeTimestamp
	case PTypeTimestampz:
		return table.ColTypeTimestampTz
	case PTypeInterval:
		return table.ColTypeInterval
	case PTypeDate:
		return table.ColTypeDate
	case PTypeBoolArray:
		return table.ColTypeBoolArray
	case PTypeInt2Array, PTypeInt4Array, PTypeInt8Array:
		return table.ColTypeIntArray
	case PTypeFloat4Array, PTypeFloat8Array:
		return table.ColTypeFloatArray
	case PTypeTextArray, PTypeVarcharArray, PTypeNameArray:
		return table.ColTypeStringArray
	default:
		return table.ColTypeString // по умолчанию считаем строкой
	}
}
