// serializers/int.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

type Int2Serializer struct{}

var Int2SerializerInstance BaseSerializer[int16, ptypes.TypeInt2] = &Int2Serializer{}

func (s *Int2Serializer) Serialize(allocator pmem.Allocator, value int16) ([]byte, error) {
	buf, err := allocator.AllocAligned(2, 2)
	if err != nil {
		return nil, fmt.Errorf("int2 serialize: %w", err)
	}
	binary.BigEndian.PutUint16(buf, uint16(value)^0x8000)
	return buf, nil
}

func (s *Int2Serializer) Deserialize(data []byte) (ptypes.TypeInt2, error) {
	if len(data) < 2 {
		return ptypes.TypeInt2{}, fmt.Errorf("int2 deserialize: expected 2 bytes, got %d", len(data))
	}
	return ptypes.TypeInt2{BufferPtr: data[:2]}, nil
}

func (s *Int2Serializer) Validate(value ptypes.TypeInt2) error {
	if len(value.BufferPtr) < 2 {
		return fmt.Errorf("int2 validate: buffer too short, expected 2, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *Int2Serializer) GetType() ptypes.OID { return ptypes.PTypeInt2 }

// ----

type Int4Serializer struct{}

var Int4SerializerInstance BaseSerializer[int32, ptypes.TypeInt4] = &Int4Serializer{}

func (s *Int4Serializer) Serialize(allocator pmem.Allocator, value int32) ([]byte, error) {
	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return nil, fmt.Errorf("int4 serialize: %w", err)
	}
	binary.BigEndian.PutUint32(buf, uint32(value)^0x80000000)
	return buf, nil
}

func (s *Int4Serializer) Deserialize(data []byte) (ptypes.TypeInt4, error) {
	if len(data) < 4 {
		return ptypes.TypeInt4{}, fmt.Errorf("int4 deserialize: expected 4 bytes, got %d", len(data))
	}
	return ptypes.TypeInt4{BufferPtr: data[:4]}, nil
}

func (s *Int4Serializer) Validate(value ptypes.TypeInt4) error {
	if len(value.BufferPtr) < 4 {
		return fmt.Errorf("int4 validate: buffer too short, expected 4, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *Int4Serializer) GetType() ptypes.OID { return ptypes.PTypeInt4 }

// ----

type Int8Serializer struct{}

var Int8SerializerInstance BaseSerializer[int64, ptypes.TypeInt8] = &Int8Serializer{}

func (s *Int8Serializer) Serialize(allocator pmem.Allocator, value int64) ([]byte, error) {
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return nil, fmt.Errorf("int8 serialize: %w", err)
	}
	binary.BigEndian.PutUint64(buf, uint64(value)^0x8000000000000000)
	return buf, nil
}

func (s *Int8Serializer) Deserialize(data []byte) (ptypes.TypeInt8, error) {
	if len(data) < 8 {
		return ptypes.TypeInt8{}, fmt.Errorf("int8 deserialize: expected 8 bytes, got %d", len(data))
	}
	return ptypes.TypeInt8{BufferPtr: data[:8]}, nil
}

func (s *Int8Serializer) Validate(value ptypes.TypeInt8) error {
	if len(value.BufferPtr) < 8 {
		return fmt.Errorf("int8 validate: buffer too short, expected 8, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *Int8Serializer) GetType() ptypes.OID { return ptypes.PTypeInt8 }
