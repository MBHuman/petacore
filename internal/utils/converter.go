package utils

import "fmt"

// toFloat64 converts a value to float64 if possible
func ToFloat64(v interface{}) (float64, bool) {
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

func ToBool(v interface{}) (bool, bool) {
	switch val := v.(type) {
	case bool:
		return val, true
	case string:
		switch val {
		case "true":
			return true, true
		case "false":
			return false, true
		}
		return false, false
	case int:
		return val != 0, true
	case int32:
		return val != 0, true
	case int64:
		return val != 0, true
	case float32:
		return val != 0.0, true
	case float64:
		return val != 0.0, true
	default:
		return false, false
	}
}

func ToInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int32:
		return int(val), true
	case int64:
		return int(val), true
	case float32:
		return int(val), true
	case float64:
		return int(val), true
	case string:
		var i int
		n, err := fmt.Sscanf(val, "%d", &i)
		if err == nil && n == 1 {
			return i, true
		}
		return 0, false
	default:
		return 0, false
	}
}
