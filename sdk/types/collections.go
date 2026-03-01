// types/array.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
)

// Comparator определяет порядок сравнения элементов массива
type Comparator[K any, T BaseType[K]] func(a, b T) int

type ElementFactory[K any, T BaseType[K]] func(buf []byte) T

type CollectionType[K any, T BaseType[K]] interface {
	BaseType[[]T]

	GetPos(idx int) BaseType[T]
	DeletePos(allocator pmem.Allocator, idx int) CollectionType[K, T]
	SetPos(allocator pmem.Allocator, idx int, value BaseType[T]) CollectionType[K, T]
	Append(allocator pmem.Allocator, value BaseType[T]) (CollectionType[K, T], error)
	Slice(allocator pmem.Allocator, start, length int) (CollectionType[K, T], error)
	Len() int

	IsNull() bool
	IsNotNull() bool

	LessThan(other BaseType[[]T]) bool
	GreaterThan(other BaseType[[]T]) bool
	LessOrEqual(other BaseType[[]T]) bool
	GreaterOrEqual(other BaseType[[]T]) bool
	Between(low, high BaseType[[]T]) bool

	Contains(value BaseType[T]) bool
}

type TypeArray[K any, T BaseType[K]] struct {
	BufferPtr  []byte
	Factory    ElementFactory[K, T]
	Comparator Comparator[K, T]
}

var _ CollectionType[bool, TypeBool] = (*TypeArray[bool, TypeBool])(nil)

func (t *TypeArray[K, T]) GetType() OID {
	if len(t.BufferPtr) < 4 {
		panic(fmt.Sprintf("array: buffer too short to read inner type OID, expected at least 4 bytes, got %d", len(t.BufferPtr)))
	}
	return OID(binary.BigEndian.Uint32(t.BufferPtr[:4]))
}

func (t *TypeArray[K, T]) GetBuffer() []byte { return t.BufferPtr }

// compare выполняет семантическое сравнение двух массивов
// если Comparator задан — элемент за элементом
// иначе — лексикографически по буферу
func (t *TypeArray[K, T]) compare(other BaseType[[]T]) int {
	if t.Comparator == nil {
		return bytes.Compare(t.BufferPtr, other.GetBuffer())
	}

	otherArr, ok := other.(*TypeArray[K, T])
	if !ok {
		return bytes.Compare(t.BufferPtr, other.GetBuffer())
	}

	aLen := t.Len()
	bLen := otherArr.Len()
	minLen := aLen
	if bLen < minLen {
		minLen = bLen
	}

	// сравниваем поэлементно
	for i := 0; i < minLen; i++ {
		aBuf, err := t.GetPosBuffer(i)
		if err != nil {
			break
		}
		bBuf, err := otherArr.GetPosBuffer(i)
		if err != nil {
			break
		}
		aElem := t.Factory(aBuf)
		bElem := t.Factory(bBuf)
		if cmp := t.Comparator(aElem, bElem); cmp != 0 {
			return cmp
		}
	}

	// все общие элементы равны — короткий массив меньше
	if aLen < bLen {
		return -1
	}
	if aLen > bLen {
		return 1
	}
	return 0
}

func (t *TypeArray[K, T]) Compare(other BaseType[[]T]) int {
	return t.compare(other)
}

func (t *TypeArray[K, T]) Len() int {
	if len(t.BufferPtr) < 8 {
		return 0
	}
	return int(binary.BigEndian.Uint32(t.BufferPtr[4:8]))
}

func (t *TypeArray[K, T]) IntoGo() []T {
	count := t.Len()
	result := make([]T, count)
	for i := 0; i < count; i++ {
		buf, err := t.GetPosBuffer(i)
		if err != nil {
			continue
		}
		result[i] = t.Factory(buf)
	}
	return result
}

func (t *TypeArray[K, T]) IsNull() bool    { return t.BufferPtr == nil }
func (t *TypeArray[K, T]) IsNotNull() bool { return t.BufferPtr != nil }

func (t *TypeArray[K, T]) LessThan(other BaseType[[]T]) bool       { return t.compare(other) < 0 }
func (t *TypeArray[K, T]) GreaterThan(other BaseType[[]T]) bool    { return t.compare(other) > 0 }
func (t *TypeArray[K, T]) LessOrEqual(other BaseType[[]T]) bool    { return t.compare(other) <= 0 }
func (t *TypeArray[K, T]) GreaterOrEqual(other BaseType[[]T]) bool { return t.compare(other) >= 0 }
func (t *TypeArray[K, T]) Between(low, high BaseType[[]T]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

func (t *TypeArray[K, T]) Contains(value BaseType[T]) bool {
	count := t.Len()
	needle := value.GetBuffer()
	for i := 0; i < count; i++ {
		buf, err := t.GetPosBuffer(i)
		if err != nil {
			continue
		}
		if bytes.Equal(buf, needle) {
			return true
		}
	}
	return false
}

func (t *TypeArray[K, T]) Append(allocator pmem.Allocator, value BaseType[T]) (CollectionType[K, T], error) {
	elements := t.readAllBuffers()
	elements = append(elements, value.GetBuffer())
	buf, err := SerializeArrayElements(allocator, t.GetType(), elements)
	if err != nil {
		return nil, fmt.Errorf("array append: %w", err)
	}
	return &TypeArray[K, T]{BufferPtr: buf, Factory: t.Factory, Comparator: t.Comparator}, nil
}

func (t *TypeArray[K, T]) Slice(allocator pmem.Allocator, start, length int) (CollectionType[K, T], error) {
	count := t.Len()
	if start < 0 || start+length > count {
		return nil, fmt.Errorf("array slice: [%d:%d] out of bounds %d", start, start+length, count)
	}
	elements := make([][]byte, length)
	for i := 0; i < length; i++ {
		buf, err := t.GetPosBuffer(start + i)
		if err != nil {
			return nil, fmt.Errorf("array slice: %w", err)
		}
		elements[i] = buf
	}
	buf, err := SerializeArrayElements(allocator, t.GetType(), elements)
	if err != nil {
		return nil, fmt.Errorf("array slice: %w", err)
	}
	return &TypeArray[K, T]{BufferPtr: buf, Factory: t.Factory, Comparator: t.Comparator}, nil
}

func (t *TypeArray[K, T]) GetPos(idx int) BaseType[T] {
	buf, err := t.GetPosBuffer(idx)
	if err != nil {
		return nil
	}
	return &arrayElemWrapper[K, T]{val: t.Factory(buf)}
}

func (t *TypeArray[K, T]) DeletePos(allocator pmem.Allocator, idx int) CollectionType[K, T] {
	count := t.Len()
	if idx < 0 || idx >= count {
		return t
	}
	elements := t.readAllBuffers()
	newElements := make([][]byte, 0, count-1)
	for i, e := range elements {
		if i != idx {
			newElements = append(newElements, e)
		}
	}
	buf, err := SerializeArrayElements(allocator, t.GetType(), newElements)
	if err != nil {
		return t
	}
	return &TypeArray[K, T]{BufferPtr: buf, Factory: t.Factory, Comparator: t.Comparator}
}

func (t *TypeArray[K, T]) SetPos(allocator pmem.Allocator, idx int, value BaseType[T]) CollectionType[K, T] {
	count := t.Len()
	if idx < 0 || idx >= count {
		return t
	}
	elements := t.readAllBuffers()
	elements[idx] = value.GetBuffer()
	buf, err := SerializeArrayElements(allocator, t.GetType(), elements)
	if err != nil {
		return t
	}
	return &TypeArray[K, T]{BufferPtr: buf, Factory: t.Factory, Comparator: t.Comparator}
}

func (t *TypeArray[K, T]) GetPosBuffer(idx int) ([]byte, error) {
	count := t.Len()
	if idx < 0 || idx >= count {
		return nil, fmt.Errorf("array: index %d out of bounds [0, %d)", idx, count)
	}
	offsetPos := 8 + idx*4
	lengthPos := 8 + count*4 + idx*4
	if len(t.BufferPtr) < lengthPos+4 {
		return nil, fmt.Errorf("array: buffer too short")
	}
	offset := int(binary.BigEndian.Uint32(t.BufferPtr[offsetPos : offsetPos+4]))
	length := int(binary.BigEndian.Uint32(t.BufferPtr[lengthPos : lengthPos+4]))
	if offset+length > len(t.BufferPtr) {
		return nil, fmt.Errorf("array: element %d out of buffer bounds", idx)
	}
	return t.BufferPtr[offset : offset+length], nil
}

func (t *TypeArray[K, T]) readAllBuffers() [][]byte {
	count := t.Len()
	result := make([][]byte, count)
	for i := 0; i < count; i++ {
		buf, _ := t.GetPosBuffer(i)
		result[i] = buf
	}
	return result
}

func SerializeArrayElements(allocator pmem.Allocator, innerType OID, elements [][]byte) ([]byte, error) {
	count := len(elements)
	dataSize := 0
	for _, e := range elements {
		dataSize += len(e)
	}
	headerSize := 4 + 4 + count*4 + count*4
	totalSize := headerSize + dataSize
	buf, err := allocator.Alloc(totalSize)
	if err != nil {
		return nil, fmt.Errorf("array: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint32(buf[0:4], uint32(innerType))
	binary.BigEndian.PutUint32(buf[4:8], uint32(count))
	dataOffset := headerSize
	for i, elem := range elements {
		offsetPos := 8 + i*4
		lengthPos := 8 + count*4 + i*4
		binary.BigEndian.PutUint32(buf[offsetPos:offsetPos+4], uint32(dataOffset))
		binary.BigEndian.PutUint32(buf[lengthPos:lengthPos+4], uint32(len(elem)))
		copy(buf[dataOffset:], elem)
		dataOffset += len(elem)
	}
	return buf, nil
}

type arrayElemWrapper[K any, T BaseType[K]] struct{ val T }

func (w *arrayElemWrapper[K, T]) GetType() OID      { return w.val.GetType() }
func (w *arrayElemWrapper[K, T]) GetBuffer() []byte { return w.val.GetBuffer() }
func (w *arrayElemWrapper[K, T]) Compare(other BaseType[T]) int {
	return bytes.Compare(w.val.GetBuffer(), other.GetBuffer())
}
func (w *arrayElemWrapper[K, T]) IntoGo() T { return w.val }
