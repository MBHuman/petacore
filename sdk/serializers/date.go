// serializers/date.go
package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"time"
)

type DateSerializer struct{}

var DateSerializerInstance BaseSerializer[*time.Time, ptypes.TypeDate] = &DateSerializer{}

func (s DateSerializer) Serialize(allocator pmem.Allocator, val *time.Time) ([]byte, error) {
	if val == nil {
		return nil, nil
	}
	d, err := ptypes.NewTypeDate(allocator, *val)
	if err != nil {
		return nil, fmt.Errorf("date serialize: %w", err)
	}
	return d.GetBuffer(), nil
}

func (s DateSerializer) Deserialize(buf []byte) (ptypes.TypeDate, error) {
	if buf == nil {
		return ptypes.TypeDate{}, nil
	}
	if len(buf) < 4 {
		return ptypes.TypeDate{}, fmt.Errorf("date deserialize: buffer too short %d", len(buf))
	}
	return ptypes.TypeDate{BufferPtr: buf}, nil
}

func (s DateSerializer) Validate(val ptypes.TypeDate) error {
	if len(val.GetBuffer()) != 4 {
		return fmt.Errorf("date validate: expected 4 bytes, got %d", len(val.GetBuffer()))
	}
	return nil
}
