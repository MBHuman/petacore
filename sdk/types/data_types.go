package ptypes

import (
	"strings"
)

func ParseDataType(typeStr string) OID {
	upper := strings.ToUpper(typeStr)

	// Strip size/precision modifiers, e.g. CHAR(1) -> CHAR, VARCHAR(255) -> VARCHAR
	if idx := strings.IndexByte(upper, '('); idx != -1 {
		upper = strings.TrimSpace(upper[:idx])
	}

	switch upper {
	case "BOOL", "BOOLEAN":
		return PTypeBool
	case "BYTEA":
		return PTypeBytea
	case "CHAR", "\"CHAR\"":
		return PTypeChar
	case "NAME":
		return PTypeName
	case "INT8", "BIGINT":
		return PTypeInt8
	case "INT2", "SMALLINT":
		return PTypeInt2
	case "INT", "INT4", "INTEGER":
		return PTypeInt4
	case "STRING", "TEXT":
		return PTypeText
	case "FLOAT4", "REAL":
		return PTypeFloat4
	case "FLOAT", "FLOAT8", "DOUBLE", "DOUBLEPRECISION":
		return PTypeFloat8
	case "VARCHAR", "CHARACTERVARYING":
		return PTypeVarchar
	case "NUMERIC", "DECIMAL":
		return PTypeNumeric
	case "DATE":
		return PTypeDate
	case "TIME", "TIMEWITHOUTtimezone":
		return PTypeTime
	case "TIMESTAMP", "TIMESTAMPWITHOUTTIMEZONE":
		return PTypeTimestamp
	case "TIMESTAMPTZ", "TIMESTAMPWITHTIMEZONE":
		return PTypeTimestampz
	default:
		return PTypeText
	}
}
