// types/bool.go
package ptypes

import (
	"bytes"
	"encoding/binary"
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

func (t TypeBool) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	switch targetType {
	case PTypeInt2:
		buf, err := allocator.AllocAligned(2, 2)
		if err != nil {
			return nil, fmt.Errorf("bool cast to int2: %w", err)
		}
		var v uint16
		if t.IntoGo() {
			v = uint16(1) ^ 0x8000
		} else {
			v = uint16(0) ^ 0x8000
		}
		binary.BigEndian.PutUint16(buf, v)
		return anyWrapper[int16]{TypeInt2{BufferPtr: buf}}, nil

	case PTypeInt4:
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("bool cast to int4: %w", err)
		}
		var v uint32
		if t.IntoGo() {
			v = uint32(1) ^ 0x80000000
		} else {
			v = uint32(0) ^ 0x80000000
		}
		binary.BigEndian.PutUint32(buf, v)
		return anyWrapper[int32]{TypeInt4{BufferPtr: buf}}, nil

	case PTypeInt8:
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("bool cast to int8: %w", err)
		}
		var v uint64
		if t.IntoGo() {
			v = uint64(1) ^ 0x8000000000000000
		} else {
			v = uint64(0) ^ 0x8000000000000000
		}
		binary.BigEndian.PutUint64(buf, v)
		return anyWrapper[int64]{TypeInt8{BufferPtr: buf}}, nil

	case PTypeText, PTypeVarchar:
		s := "false"
		if t.IntoGo() {
			s = "true"
		}
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("bool cast to text: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("bool: unsupported cast to OID %d", targetType)
	}
}
