// types/float.go
package ptypes

import (
	"math"
)

// OrderableFloat32bits конвертирует float32 в order-preserving uint32
func OrderableFloat32bits(f float32) uint32 {
	bits := math.Float32bits(f)
	if bits>>31 == 1 {
		// отрицательное — инвертируем все биты
		return ^bits
	}
	// положительное — инвертируем только знаковый бит
	return bits ^ 0x80000000
}

// Float32fromOrderableBits — обратная операция
func Float32fromOrderableBits(bits uint32) float32 {
	if bits>>31 == 0 {
		// был отрицательным
		return math.Float32frombits(^bits)
	}
	return math.Float32frombits(bits ^ 0x80000000)
}

// OrderableFloat64bits конвертирует float64 в order-preserving uint64
func OrderableFloat64bits(f float64) uint64 {
	bits := math.Float64bits(f)
	if bits>>63 == 1 {
		return ^bits
	}
	return bits ^ 0x8000000000000000
}

func Float64fromOrderableBits(bits uint64) float64 {
	if bits>>63 == 0 {
		return math.Float64frombits(^bits)
	}
	return math.Float64frombits(bits ^ 0x8000000000000000)
}
