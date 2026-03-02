// serializers/text.go
package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

type TextSerializer struct{}

var TextSerializerInstance = &TextSerializer{}

func (s *TextSerializer) Serialize(allocator pmem.Allocator, value string) ([]byte, error) {
	b := []byte(value)
	buf, err := allocator.Alloc(len(b))
	if err != nil {
		return nil, fmt.Errorf("text serialize: %w", err)
	}
	copy(buf, b)
	return buf, nil
}

func (s *TextSerializer) Deserialize(data []byte) (ptypes.TypeText, error) {
	return ptypes.TypeText{BufferPtr: data}, nil
}

func (s *TextSerializer) Validate(value ptypes.TypeText) error {
	if value.BufferPtr == nil {
		return fmt.Errorf("text validate: nil buffer")
	}
	return nil
}
