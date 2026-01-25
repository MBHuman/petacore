package rhelpers

func IsTrue(value interface{}) bool {
	if value == nil {
		return false
	}
	if boolVal, ok := value.(bool); ok {
		return boolVal
	}
	// For other types, treat as true if not zero/false
	if intVal, ok := value.(int); ok {
		return intVal != 0
	}
	if floatVal, ok := value.(float64); ok {
		return floatVal != 0
	}
	if strVal, ok := value.(string); ok {
		return strVal != ""
	}
	return true
}
