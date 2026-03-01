// types/int8.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
)

type TypeInt8 struct {
	BufferPtr []byte
}

var _ BaseType[int64] = (*TypeInt8)(nil)
var _ NumericType[int64] = (*TypeInt8)(nil)
var _ OrderedType[int64] = (*TypeInt8)(nil)
var _ NullableType[int64] = (*TypeInt8)(nil)
var _ BitwiseType[int64] = (*TypeInt8)(nil)

func (t TypeInt8) GetType() OID { return PTypeInt8 }

func (t TypeInt8) Compare(other BaseType[int64]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeInt8) GetBuffer() []byte { return t.BufferPtr }

func (t TypeInt8) IntoGo() int64 {
	if len(t.BufferPtr) < 8 {
		return 0
	}
	return int64(binary.BigEndian.Uint64(t.BufferPtr) ^ 0x8000000000000000)
}

// NullableType

func (t TypeInt8) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeInt8) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeInt8) LessThan(other BaseType[int64]) bool       { return t.Compare(other) < 0 }
func (t TypeInt8) GreaterThan(other BaseType[int64]) bool    { return t.Compare(other) > 0 }
func (t TypeInt8) LessOrEqual(other BaseType[int64]) bool    { return t.Compare(other) <= 0 }
func (t TypeInt8) GreaterOrEqual(other BaseType[int64]) bool { return t.Compare(other) >= 0 }
func (t TypeInt8) Between(low, high BaseType[int64]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func int8FromVal(allocator pmem.Allocator, v int64) (TypeInt8, error) {
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return TypeInt8{}, fmt.Errorf("int8: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint64(buf, uint64(v)^0x8000000000000000)
	return TypeInt8{BufferPtr: buf}, nil
}

// NumericType

func (t TypeInt8) Add(allocator pmem.Allocator, other NumericType[int64]) (NumericType[int64], error) {
	return int8FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeInt8) Sub(allocator pmem.Allocator, other NumericType[int64]) (NumericType[int64], error) {
	return int8FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeInt8) Mul(allocator pmem.Allocator, other NumericType[int64]) (NumericType[int64], error) {
	return int8FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeInt8) Div(allocator pmem.Allocator, other NumericType[int64]) (NumericType[int64], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int8: division by zero")
	}
	return int8FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeInt8) Mod(allocator pmem.Allocator, other NumericType[int64]) (NumericType[int64], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int8: modulo by zero")
	}
	return int8FromVal(allocator, t.IntoGo()%other.IntoGo())
}

func (t TypeInt8) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeInt8) Neg(allocator pmem.Allocator) NumericType[int64] {
	v, _ := int8FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeInt8) Abs(allocator pmem.Allocator) NumericType[int64] {
	v := t.IntoGo()
	if v < 0 {
		v = -v
	}
	result, _ := int8FromVal(allocator, v)
	return result
}

// BitwiseType

func (t TypeInt8) And(allocator pmem.Allocator, other BitwiseType[int64]) BitwiseType[int64] {
	v, _ := int8FromVal(allocator, t.IntoGo()&other.IntoGo())
	return v
}

func (t TypeInt8) Or(allocator pmem.Allocator, other BitwiseType[int64]) BitwiseType[int64] {
	v, _ := int8FromVal(allocator, t.IntoGo()|other.IntoGo())
	return v
}

func (t TypeInt8) Xor(allocator pmem.Allocator, other BitwiseType[int64]) BitwiseType[int64] {
	v, _ := int8FromVal(allocator, t.IntoGo()^other.IntoGo())
	return v
}

func (t TypeInt8) Not(allocator pmem.Allocator) BitwiseType[int64] {
	v, _ := int8FromVal(allocator, ^t.IntoGo())
	return v
}

func (t TypeInt8) ShiftLeft(allocator pmem.Allocator, n uint) BitwiseType[int64] {
	v, _ := int8FromVal(allocator, t.IntoGo()<<n)
	return v
}

func (t TypeInt8) ShiftRight(allocator pmem.Allocator, n uint) BitwiseType[int64] {
	v, _ := int8FromVal(allocator, t.IntoGo()>>n)
	return v
}
