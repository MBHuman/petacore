// types/float8.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"petacore/sdk/pmem"
)

type TypeFloat8 struct {
	BufferPtr []byte
}

var _ BaseType[float64] = (*TypeFloat8)(nil)
var _ NumericType[float64] = (*TypeFloat8)(nil)
var _ OrderedType[float64] = (*TypeFloat8)(nil)
var _ NullableType[float64] = (*TypeFloat8)(nil)

func (t TypeFloat8) GetType() OID { return PTypeFloat8 }

func (t TypeFloat8) Compare(other BaseType[float64]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeFloat8) GetBuffer() []byte { return t.BufferPtr }

func (t TypeFloat8) IntoGo() float64 {
	if len(t.BufferPtr) < 8 {
		return 0
	}
	return Float64fromOrderableBits(binary.BigEndian.Uint64(t.BufferPtr))
}

// NullableType

func (t TypeFloat8) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeFloat8) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeFloat8) LessThan(other BaseType[float64]) bool       { return t.Compare(other) < 0 }
func (t TypeFloat8) GreaterThan(other BaseType[float64]) bool    { return t.Compare(other) > 0 }
func (t TypeFloat8) LessOrEqual(other BaseType[float64]) bool    { return t.Compare(other) <= 0 }
func (t TypeFloat8) GreaterOrEqual(other BaseType[float64]) bool { return t.Compare(other) >= 0 }
func (t TypeFloat8) Between(low, high BaseType[float64]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func float8FromVal(allocator pmem.Allocator, v float64) (TypeFloat8, error) {
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return TypeFloat8{}, fmt.Errorf("float8: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint64(buf, OrderableFloat64bits(v))
	return TypeFloat8{BufferPtr: buf}, nil
}

// NumericType

func (t TypeFloat8) Add(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	return float8FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeFloat8) Sub(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	return float8FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeFloat8) Mul(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	return float8FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeFloat8) Div(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float8: division by zero")
	}
	return float8FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeFloat8) Mod(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float8: modulo by zero")
	}
	return float8FromVal(allocator, math.Mod(t.IntoGo(), other.IntoGo()))
}

func (t TypeFloat8) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeFloat8) Neg(allocator pmem.Allocator) NumericType[float64] {
	v, _ := float8FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeFloat8) Abs(allocator pmem.Allocator) NumericType[float64] {
	v, _ := float8FromVal(allocator, math.Abs(t.IntoGo()))
	return v
}

func (t TypeFloat8) String() string {
	return "float8(" + fmt.Sprintf("%v", t.IntoGo()) + ")"
}
