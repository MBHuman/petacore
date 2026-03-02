// serializers/array.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

type ArraySerializer[K any, T ptypes.BaseType[K]] struct {
	InnerType ptypes.OID
	Factory   ptypes.ElementFactory[K, T]
}

func NewArraySerializer[K any, T ptypes.BaseType[K]](innerType ptypes.OID, factory ptypes.ElementFactory[K, T]) *ArraySerializer[K, T] {
	return &ArraySerializer[K, T]{InnerType: innerType, Factory: factory}
}

func (s *ArraySerializer[K, T]) Serialize(allocator pmem.Allocator, elements [][]byte) ([]byte, error) {
	return ptypes.SerializeArrayElements(allocator, s.InnerType, elements)
}

func (s *ArraySerializer[K, T]) Deserialize(data []byte) (*ptypes.TypeArray[K, T], error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("array deserialize: buffer too short")
	}
	return &ptypes.TypeArray[K, T]{
		BufferPtr: data,
		Factory:   s.Factory,
	}, nil
}

func (s *ArraySerializer[K, T]) Validate(value *ptypes.TypeArray[K, T]) error {
	if len(value.BufferPtr) < 8 {
		return fmt.Errorf("array validate: buffer too short")
	}
	count := int(binary.BigEndian.Uint32(value.BufferPtr[4:8]))
	minSize := 8 + count*8
	if len(value.BufferPtr) < minSize {
		return fmt.Errorf("array validate: buffer too short for %d elements", count)
	}
	// проверяем что OID в буфере совпадает с ожидаемым
	storedOID := ptypes.OID(binary.BigEndian.Uint32(value.BufferPtr[:4]))
	if storedOID != s.InnerType {
		return fmt.Errorf("array validate: OID mismatch, expected %d got %d", s.InnerType, storedOID)
	}
	return nil
}

func (s *ArraySerializer[K, T]) Append(allocator pmem.Allocator, arr *ptypes.TypeArray[K, T], elemBuf []byte) (*ptypes.TypeArray[K, T], error) {
	count := arr.Len()
	elements := make([][]byte, count+1)
	for i := 0; i < count; i++ {
		buf, err := arr.GetPosBuffer(i)
		if err != nil {
			return nil, fmt.Errorf("array append: %w", err)
		}
		elements[i] = buf
	}
	elements[count] = elemBuf

	buf, err := ptypes.SerializeArrayElements(allocator, s.InnerType, elements)
	if err != nil {
		return nil, err
	}
	return &ptypes.TypeArray[K, T]{BufferPtr: buf, Factory: s.Factory}, nil
}
