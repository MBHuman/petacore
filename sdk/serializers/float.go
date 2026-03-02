// serializers/float.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

type Float4Serializer struct{}

var Float4SerializerInstance BaseSerializer[float32, ptypes.TypeFloat4] = &Float4Serializer{}

func (s *Float4Serializer) Serialize(allocator pmem.Allocator, value float32) ([]byte, error) {
	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return nil, fmt.Errorf("float4 serialize: %w", err)
	}
	binary.BigEndian.PutUint32(buf, ptypes.OrderableFloat32bits(value))
	return buf, nil
}

func (s *Float4Serializer) Deserialize(data []byte) (ptypes.TypeFloat4, error) {
	if len(data) < 4 {
		return ptypes.TypeFloat4{}, fmt.Errorf("float4 deserialize: expected 4 bytes, got %d", len(data))
	}
	return ptypes.TypeFloat4{BufferPtr: data[:4]}, nil
}

func (s *Float4Serializer) Validate(value ptypes.TypeFloat4) error {
	if len(value.BufferPtr) < 4 {
		return fmt.Errorf("float4 validate: buffer too short, expected 4, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *Float4Serializer) GetType() ptypes.OID {
	return ptypes.PTypeFloat4
}

// ----

type Float8Serializer struct{}

var Float8SerializerInstance BaseSerializer[float64, ptypes.TypeFloat8] = &Float8Serializer{}

func (s *Float8Serializer) Serialize(allocator pmem.Allocator, value float64) ([]byte, error) {
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return nil, fmt.Errorf("float8 serialize: %w", err)
	}
	binary.BigEndian.PutUint64(buf, ptypes.OrderableFloat64bits(value))
	return buf, nil
}

func (s *Float8Serializer) Deserialize(data []byte) (ptypes.TypeFloat8, error) {
	if len(data) < 8 {
		return ptypes.TypeFloat8{}, fmt.Errorf("float8 deserialize: expected 8 bytes, got %d", len(data))
	}
	return ptypes.TypeFloat8{BufferPtr: data[:8]}, nil
}

func (s *Float8Serializer) Validate(value ptypes.TypeFloat8) error {
	if len(value.BufferPtr) < 8 {
		return fmt.Errorf("float8 validate: buffer too short, expected 8, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *Float8Serializer) GetType() ptypes.OID {
	return ptypes.PTypeFloat8
}
