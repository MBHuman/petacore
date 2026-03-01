package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

type BytesSerializer struct{}

var BytesSerializerInstance BaseSerializer[[]byte, ptypes.TypeBytea] = &BytesSerializer{}

func (s *BytesSerializer) Serialize(allocator pmem.Allocator, value []byte) ([]byte, error) {
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

func (s *BytesSerializer) Deserialize(data []byte) (ptypes.TypeBytea, error) {
	return ptypes.TypeBytea{BufferPtr: data}, nil
}

func (s *BytesSerializer) Validate(value ptypes.TypeBytea) error {
	return nil
}

func (s *BytesSerializer) GetType() ptypes.OID {
	return ptypes.PTypeBytea
}
