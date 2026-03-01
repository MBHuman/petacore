// serializers/date.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"time"
)

type DateSerializer struct{}

var DateSerializerInstance BaseSerializer[*time.Time, ptypes.TypeDate] = &DateSerializer{}

func (d *DateSerializer) Serialize(allocator pmem.Allocator, value *time.Time) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("date serialize: value is nil")
	}

	v := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
	days := int32(v.Sub(ptypes.PgEpoch).Hours() / 24)

	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return nil, fmt.Errorf("date serialize: %w", err)
	}

	binary.BigEndian.PutUint32(buf, uint32(days))
	return buf, nil
}

func (d *DateSerializer) Deserialize(data []byte) (ptypes.TypeDate, error) {
	if len(data) < 4 {
		return ptypes.TypeDate{}, fmt.Errorf("date deserialize: expected 4 bytes, got %d", len(data))
	}
	return ptypes.TypeDate{BufferPtr: data[:4]}, nil
}

func (d *DateSerializer) Validate(value ptypes.TypeDate) error {
	if len(value.BufferPtr) < 4 {
		return fmt.Errorf("date validate: buffer too short, expected 4, got %d", len(value.BufferPtr))
	}
	return nil
}

func (d *DateSerializer) GetType() ptypes.OID {
	return ptypes.PTypeDate
}
