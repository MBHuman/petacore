package rhelpers

import "fmt"

// Helper functions for arithmetic operations
func addValues(a, b interface{}) interface{} {
	// Convert to float64 for simplicity
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af + bf
	}
	return fmt.Sprintf("%v%v", a, b) // fallback to string concatenation
}

func subtractValues(a, b interface{}) interface{} {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af - bf
	}
	return 0 // fallback
}

func multiplyValues(a, b interface{}) interface{} {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af * bf
	}
	return 0 // fallback
}

func divideValues(a, b interface{}) interface{} {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok && bf != 0 {
		return af / bf
	}
	return 0 // fallback
}
