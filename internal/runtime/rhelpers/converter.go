package rhelpers

import "petacore/internal/runtime/rsql/table"

// toFloat64 converts a value to float64 if possible
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

// GetColumnNamesFromRow extracts column names from a row map in insertion order
func GetColumnNamesFromRow(row map[string]interface{}) []string {
	var cols []string
	for k := range row {
		cols = append(cols, k)
	}
	return cols
}

// GetColumnTypesFromRow infers column types from a row map values in insertion order
func GetColumnTypesFromRow(row map[string]interface{}) []table.ColType {
	var types []table.ColType
	for _, v := range row {
		if v == nil {
			types = append(types, table.ColTypeString) // default for nil
		} else {
			switch v.(type) {
			case int, int32, int64:
				types = append(types, table.ColTypeInt)
			case float32, float64:
				types = append(types, table.ColTypeFloat)
			case bool:
				types = append(types, table.ColTypeBool)
			default:
				types = append(types, table.ColTypeString)
			}
		}
	}
	return types
}
