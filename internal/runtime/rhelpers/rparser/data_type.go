package rparser

import (
	ptypes "petacore/sdk/types"
	"strings"
)

func ParseDataType(typeStr string) ptypes.OID {
	upper := strings.ToUpper(typeStr)

	// Strip size/precision modifiers, e.g. CHAR(1) -> CHAR, VARCHAR(255) -> VARCHAR
	if idx := strings.IndexByte(upper, '('); idx != -1 {
		upper = strings.TrimSpace(upper[:idx])
	}

	switch upper {
	case "BOOL", "BOOLEAN":
		return ptypes.PTypeBool
	case "BYTEA":
		return ptypes.PTypeBytea
	case "CHAR", "\"CHAR\"":
		return ptypes.PTypeChar
	case "NAME":
		return ptypes.PTypeName
	case "INT8", "BIGINT":
		return ptypes.PTypeInt8
	case "INT2", "SMALLINT":
		return ptypes.PTypeInt2
	case "INT", "INT4", "INTEGER":
		return ptypes.PTypeInt4
	case "STRING", "TEXT":
		return ptypes.PTypeText
	case "FLOAT4", "REAL":
		return ptypes.PTypeFloat4
	case "FLOAT", "FLOAT8", "DOUBLE", "DOUBLEPRECISION":
		return ptypes.PTypeFloat8
	case "VARCHAR", "CHARACTERVARYING":
		return ptypes.PTypeVarchar
	case "NUMERIC", "DECIMAL":
		return ptypes.PTypeNumeric
	case "DATE":
		return ptypes.PTypeDate
	case "TIME", "TIMEWITHOUTtimezone":
		return ptypes.PTypeTime
	case "TIMESTAMP", "TIMESTAMPWITHOUTTIMEZONE":
		return ptypes.PTypeTimestamp
	case "TIMESTAMPTZ", "TIMESTAMPWITHTIMEZONE":
		return ptypes.PTypeTimestampz
	default:
		return ptypes.PTypeText
	}
}
