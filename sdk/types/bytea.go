// types/bytea.go
package ptypes

import (
	"bytes"
	"fmt"
	"petacore/sdk/pmem"
)

type TypeBytea struct {
	BufferPtr []byte
}

var _ BaseType[[]byte] = (*TypeBytea)(nil)
var _ OrderedType[[]byte] = (*TypeBytea)(nil)
var _ NullableType[[]byte] = (*TypeBytea)(nil)

func NewTypeBytea(val []byte) TypeBytea {
	return TypeBytea{BufferPtr: val}
}

func (t TypeBytea) GetType() OID {
	return PTypeBytea
}

func (t TypeBytea) Compare(other BaseType[[]byte]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeBytea) GetBuffer() []byte {
	return t.BufferPtr
}

func (t TypeBytea) IntoGo() []byte {
	return t.BufferPtr
}

// OrderedType

func (t TypeBytea) LessThan(other BaseType[[]byte]) bool       { return t.Compare(other) < 0 }
func (t TypeBytea) GreaterThan(other BaseType[[]byte]) bool    { return t.Compare(other) > 0 }
func (t TypeBytea) LessOrEqual(other BaseType[[]byte]) bool    { return t.Compare(other) <= 0 }
func (t TypeBytea) GreaterOrEqual(other BaseType[[]byte]) bool { return t.Compare(other) >= 0 }
func (t TypeBytea) Between(low, high BaseType[[]byte]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// NullableType

func (t TypeBytea) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeBytea) IsNotNull() bool { return t.BufferPtr != nil }

// операции специфичные для bytea — через аллокатор

// Concat объединяет два буфера
func (t TypeBytea) Concat(allocator pmem.Allocator, other TypeBytea) (TypeBytea, error) {
	size := len(t.BufferPtr) + len(other.BufferPtr)
	buf, err := allocator.Alloc(size)
	if err != nil {
		return TypeBytea{}, fmt.Errorf("bytea concat: %w", err)
	}
	copy(buf, t.BufferPtr)
	copy(buf[len(t.BufferPtr):], other.BufferPtr)
	return TypeBytea{BufferPtr: buf}, nil
}

// Slice возвращает подмассив — zero copy, аллокатор не нужен
func (t TypeBytea) Slice(start, length int) (TypeBytea, error) {
	if start < 0 || start+length > len(t.BufferPtr) {
		return TypeBytea{}, fmt.Errorf("bytea slice: [%d:%d] out of bounds %d", start, start+length, len(t.BufferPtr))
	}
	return TypeBytea{BufferPtr: t.BufferPtr[start : start+length]}, nil
}

// Length — только читает, аллокатор не нужен
func (t TypeBytea) Length() int {
	return len(t.BufferPtr)
}

// Contains — только читает
func (t TypeBytea) Contains(sub []byte) bool {
	return bytes.Contains(t.BufferPtr, sub)
}

// StartsWith — только читает
func (t TypeBytea) StartsWith(prefix []byte) bool {
	return bytes.HasPrefix(t.BufferPtr, prefix)
}

// EndsWith — только читает
func (t TypeBytea) EndsWith(suffix []byte) bool {
	return bytes.HasSuffix(t.BufferPtr, suffix)
}

// Overlay заменяет участок буфера — создаёт новый буфер
func (t TypeBytea) Overlay(allocator pmem.Allocator, replacement []byte, start, length int) (TypeBytea, error) {
	if start < 0 || start+length > len(t.BufferPtr) {
		return TypeBytea{}, fmt.Errorf("bytea overlay: [%d:%d] out of bounds %d", start, start+length, len(t.BufferPtr))
	}

	newSize := len(t.BufferPtr) - length + len(replacement)
	buf, err := allocator.Alloc(newSize)
	if err != nil {
		return TypeBytea{}, fmt.Errorf("bytea overlay: %w", err)
	}

	copy(buf, t.BufferPtr[:start])
	copy(buf[start:], replacement)
	copy(buf[start+len(replacement):], t.BufferPtr[start+length:])

	return TypeBytea{BufferPtr: buf}, nil
}

func (t TypeBytea) String() string {
	return "bytea(" + fmt.Sprintf("%v", t.BufferPtr) + ")"
}
