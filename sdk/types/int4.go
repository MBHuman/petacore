// types/int4.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
)

type TypeInt4 struct {
	BufferPtr []byte
}

var _ BaseType[int32] = (*TypeInt4)(nil)
var _ NumericType[int32] = (*TypeInt4)(nil)
var _ OrderedType[int32] = (*TypeInt4)(nil)
var _ NullableType[int32] = (*TypeInt4)(nil)
var _ BitwiseType[int32] = (*TypeInt4)(nil)

func (t TypeInt4) GetType() OID { return PTypeInt4 }

func (t TypeInt4) Compare(other BaseType[int32]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeInt4) GetBuffer() []byte { return t.BufferPtr }

func (t TypeInt4) IntoGo() int32 {
	if len(t.BufferPtr) < 4 {
		return 0
	}
	return int32(binary.BigEndian.Uint32(t.BufferPtr) ^ 0x80000000)
}

// NullableType

func (t TypeInt4) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeInt4) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeInt4) LessThan(other BaseType[int32]) bool       { return t.Compare(other) < 0 }
func (t TypeInt4) GreaterThan(other BaseType[int32]) bool    { return t.Compare(other) > 0 }
func (t TypeInt4) LessOrEqual(other BaseType[int32]) bool    { return t.Compare(other) <= 0 }
func (t TypeInt4) GreaterOrEqual(other BaseType[int32]) bool { return t.Compare(other) >= 0 }
func (t TypeInt4) Between(low, high BaseType[int32]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func int4FromVal(allocator pmem.Allocator, v int32) (TypeInt4, error) {
	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return TypeInt4{}, fmt.Errorf("int4: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint32(buf, uint32(v)^0x80000000)
	return TypeInt4{BufferPtr: buf}, nil
}

// NumericType

func (t TypeInt4) Add(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	return int4FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeInt4) Sub(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	return int4FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeInt4) Mul(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	return int4FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeInt4) Div(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int4: division by zero")
	}
	return int4FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeInt4) Mod(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int4: modulo by zero")
	}
	return int4FromVal(allocator, t.IntoGo()%other.IntoGo())
}

func (t TypeInt4) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeInt4) Neg(allocator pmem.Allocator) NumericType[int32] {
	v, _ := int4FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeInt4) Abs(allocator pmem.Allocator) NumericType[int32] {
	v := t.IntoGo()
	if v < 0 {
		v = -v
	}
	result, _ := int4FromVal(allocator, v)
	return result
}

// BitwiseType

func (t TypeInt4) And(allocator pmem.Allocator, other BitwiseType[int32]) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()&other.IntoGo())
	return v
}

func (t TypeInt4) Or(allocator pmem.Allocator, other BitwiseType[int32]) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()|other.IntoGo())
	return v
}

func (t TypeInt4) Xor(allocator pmem.Allocator, other BitwiseType[int32]) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()^other.IntoGo())
	return v
}

func (t TypeInt4) Not(allocator pmem.Allocator) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, ^t.IntoGo())
	return v
}

func (t TypeInt4) ShiftLeft(allocator pmem.Allocator, n uint) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()<<n)
	return v
}

func (t TypeInt4) ShiftRight(allocator pmem.Allocator, n uint) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()>>n)
	return v
}
