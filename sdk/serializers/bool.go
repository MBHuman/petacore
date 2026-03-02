package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

type BoolSerializer struct{}

var BoolSerializerInstance BaseSerializer[bool, ptypes.TypeBool] = &BoolSerializer{}

func (s *BoolSerializer) Serialize(allocator pmem.Allocator, value bool) ([]byte, error) {
	buf, err := allocator.Alloc(1)
	if err != nil {
		return nil, fmt.Errorf("bool serialize: %w", err)
	}

	if value {
		buf[0] = 1
	} else {
		buf[0] = 0
	}

	return buf, nil
}

func (s *BoolSerializer) Deserialize(data []byte) (ptypes.TypeBool, error) {
	if len(data) == 0 {
		return ptypes.TypeBool{BufferPtr: []byte{0}}, nil
	}
	if data[0] != 0 && data[0] != 1 {
		return ptypes.TypeBool{}, fmt.Errorf("bool deserialize: invalid byte %d", data[0])
	}
	return ptypes.TypeBool{BufferPtr: data}, nil
}

func (s *BoolSerializer) Validate(value ptypes.TypeBool) error {
	if len(value.BufferPtr) == 0 {
		return fmt.Errorf("bool validate: empty buffer")
	}
	if value.BufferPtr[0] != 0 && value.BufferPtr[0] != 1 {
		return fmt.Errorf("bool validate: invalid byte %d, expected 0 or 1", value.BufferPtr[0])
	}
	return nil
}

func (s *BoolSerializer) GetType() ptypes.OID {
	return ptypes.PTypeBool
}
