package psdk

import (
	"petacore/internal/runtime/rsql/table"
	"reflect"
)

const (
	// OIDs for basic types
	PTypeBool    OID = 16
	PTypeBytea   OID = 17
	PTypeChar    OID = 18
	PTypeInt8    OID = 20
	PTypeInt2    OID = 21
	PTypeInt4    OID = 23
	PTypeText    OID = 25
	PTypeFloat4  OID = 700
	PTypeFloat8  OID = 701
	PTypeVarchar OID = 1043
	PTypeNumeric OID = 1700
)

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
	case PTypeText, PTypeVarchar:
		return reflect.TypeFor[string]()
	case PTypeChar:
		return reflect.TypeFor[rune]()
	case PTypeBytea:
		return reflect.TypeFor[[]byte]()
	case PTypeNumeric: // десятичная арифметика произвольной точности
		return reflect.TypeFor[float64]()
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
	case PTypeText, PTypeVarchar, PTypeChar:
		return table.ColTypeString
	case PTypeBytea:
		return table.ColTypeString // bytea можно представить как строку в кодировке base64
	default:
		return table.ColTypeString // по умолчанию считаем строкой
	}
}
