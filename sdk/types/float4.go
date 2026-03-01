// types/float.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"petacore/sdk/pmem"
)

type TypeFloat4 struct {
	BufferPtr []byte
}

var _ BaseType[float32] = (*TypeFloat4)(nil)
var _ NumericType[float32] = (*TypeFloat4)(nil)
var _ OrderedType[float32] = (*TypeFloat4)(nil)
var _ NullableType[float32] = (*TypeFloat4)(nil)

func (t TypeFloat4) GetType() OID { return PTypeFloat4 }

func (t TypeFloat4) Compare(other BaseType[float32]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeFloat4) GetBuffer() []byte { return t.BufferPtr }

func (t TypeFloat4) IntoGo() float32 {
	if len(t.BufferPtr) < 4 {
		return 0
	}
	return Float32fromOrderableBits(binary.BigEndian.Uint32(t.BufferPtr))
}

// NullableType

func (t TypeFloat4) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeFloat4) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeFloat4) LessThan(other BaseType[float32]) bool       { return t.Compare(other) < 0 }
func (t TypeFloat4) GreaterThan(other BaseType[float32]) bool    { return t.Compare(other) > 0 }
func (t TypeFloat4) LessOrEqual(other BaseType[float32]) bool    { return t.Compare(other) <= 0 }
func (t TypeFloat4) GreaterOrEqual(other BaseType[float32]) bool { return t.Compare(other) >= 0 }
func (t TypeFloat4) Between(low, high BaseType[float32]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func float4FromVal(allocator pmem.Allocator, v float32) (TypeFloat4, error) {
	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return TypeFloat4{}, fmt.Errorf("float4: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint32(buf, OrderableFloat32bits(v))
	return TypeFloat4{BufferPtr: buf}, nil
}

// NumericType

func (t TypeFloat4) Add(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	return float4FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeFloat4) Sub(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	return float4FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeFloat4) Mul(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	return float4FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeFloat4) Div(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float4: division by zero")
	}
	return float4FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeFloat4) Mod(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float4: modulo by zero")
	}
	return float4FromVal(allocator, float32(math.Mod(float64(t.IntoGo()), float64(other.IntoGo()))))
}

func (t TypeFloat4) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeFloat4) Neg(allocator pmem.Allocator) NumericType[float32] {
	v, _ := float4FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeFloat4) Abs(allocator pmem.Allocator) NumericType[float32] {
	v, _ := float4FromVal(allocator, float32(math.Abs(float64(t.IntoGo()))))
	return v
}
