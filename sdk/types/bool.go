// types/bool.go
package ptypes

import (
	"bytes"
	"fmt"
	"petacore/sdk/pmem"
)

type TypeBool struct {
	BufferPtr []byte
}

var _ BaseType[bool] = (*TypeBool)(nil)
var _ OrderedType[bool] = (*TypeBool)(nil)
var _ NullableType[bool] = (*TypeBool)(nil)

func (t TypeBool) GetType() OID {
	return PTypeBool
}

func (t TypeBool) Compare(other BaseType[bool]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeBool) GetBuffer() []byte {
	return t.BufferPtr
}

func (t TypeBool) IntoGo() bool {
	if len(t.BufferPtr) == 0 {
		return false
	}
	return t.BufferPtr[0] == 1
}

// OrderedType

func (t TypeBool) LessThan(other BaseType[bool]) bool       { return t.Compare(other) < 0 }
func (t TypeBool) GreaterThan(other BaseType[bool]) bool    { return t.Compare(other) > 0 }
func (t TypeBool) LessOrEqual(other BaseType[bool]) bool    { return t.Compare(other) <= 0 }
func (t TypeBool) GreaterOrEqual(other BaseType[bool]) bool { return t.Compare(other) >= 0 }
func (t TypeBool) Between(low, high BaseType[bool]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// NullableType

func (t TypeBool) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeBool) IsNotNull() bool { return t.BufferPtr != nil }

// логические операции через аллокатор

func (t TypeBool) And(allocator pmem.Allocator, other TypeBool) (TypeBool, error) {
	buf, err := allocator.Alloc(1)
	if err != nil {
		return TypeBool{}, fmt.Errorf("bool and: %w", err)
	}
	if t.IntoGo() && other.IntoGo() {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return TypeBool{BufferPtr: buf}, nil
}

func (t TypeBool) Or(allocator pmem.Allocator, other TypeBool) (TypeBool, error) {
	buf, err := allocator.Alloc(1)
	if err != nil {
		return TypeBool{}, fmt.Errorf("bool or: %w", err)
	}
	if t.IntoGo() || other.IntoGo() {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return TypeBool{BufferPtr: buf}, nil
}

func (t TypeBool) Not(allocator pmem.Allocator) (TypeBool, error) {
	buf, err := allocator.Alloc(1)
	if err != nil {
		return TypeBool{}, fmt.Errorf("bool not: %w", err)
	}
	if t.IntoGo() {
		buf[0] = 0
	} else {
		buf[0] = 1
	}
	return TypeBool{BufferPtr: buf}, nil
}

func (t TypeBool) String() string {
	return "bool(" + fmt.Sprintf("%v", t.IntoGo()) + ")"
}
