package psdk

import (
	"petacore/internal/runtime/rsql/table"
	"reflect"
)

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
