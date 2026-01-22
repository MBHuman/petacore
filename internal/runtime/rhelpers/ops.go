package rhelpers

import "fmt"

// TODO, реализовать логику работы с операциями через типы данных,
// всё переводить в float64 для арифметики плохая идея, нужно поддерживать разные типы

// Helper functions for arithmetic operations
func AddValues(a, b interface{}) interface{} {
	// Convert to float64 for simplicity
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af + bf
	}
	return fmt.Sprintf("%v%v", a, b) // fallback to string concatenation
}

func SubtractValues(a, b interface{}) interface{} {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af - bf
	}
	return 0 // fallback
}

func MultiplyValues(a, b interface{}) interface{} {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af * bf
	}
	return 0 // fallback
}

func DivideValues(a, b interface{}) interface{} {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok && bf != 0 {
		return af / bf
	}
	return 0 // fallback
}
