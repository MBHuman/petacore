package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

type ByteaSerializer struct{}

var ByteaSerializerInstance BaseSerializer[[]byte, ptypes.TypeBytea] = &ByteaSerializer{}

func (s *ByteaSerializer) Serialize(allocator pmem.Allocator, value []byte) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	buf, err := allocator.Alloc(len(value))
	if err != nil {
		return nil, fmt.Errorf("bytea serialize: %w", err)
	}

	copy(buf, value)
	return buf, nil
}

func (s *ByteaSerializer) Deserialize(data []byte) (ptypes.TypeBytea, error) {
	return ptypes.TypeBytea{BufferPtr: data}, nil
}

func (s *ByteaSerializer) Validate(value ptypes.TypeBytea) error {
	return nil
}

func (s *ByteaSerializer) GetType() ptypes.OID {
	return ptypes.PTypeBytea
}
