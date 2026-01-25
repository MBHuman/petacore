package rparser

import (
	"petacore/internal/runtime/rsql/table"
	"strings"
)

func ParseDataType(typeStr string) table.ColType {
	switch strings.ToUpper(typeStr) {
	case "STRING", "TEXT":
		return table.ColTypeString
	case "INT":
		return table.ColTypeInt
	case "FLOAT":
		return table.ColTypeFloat
	case "BOOL":
		return table.ColTypeBool
	default:
		return table.ColTypeString
	}
}
