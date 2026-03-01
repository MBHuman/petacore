// serializers/time.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"time"
)

type TimeSerializer struct{}

var TimeSerializerInstance BaseSerializer[*time.Time, ptypes.TypeTime] = &TimeSerializer{}

func (s *TimeSerializer) Serialize(allocator pmem.Allocator, value *time.Time) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("time serialize: value is nil")
	}

	// извлекаем только время суток — микросекунды с полуночи
	h, m, sec := value.Clock()
	nsec := value.Nanosecond()
	usec := int64(h)*3_600_000_000 +
		int64(m)*60_000_000 +
		int64(sec)*1_000_000 +
		int64(nsec)/1000

	if usec < 0 || usec > 86_399_999_999 {
		return nil, fmt.Errorf("time serialize: value out of range")
	}

	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return nil, fmt.Errorf("time serialize: %w", err)
	}

	binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
	return buf, nil
}

func (s *TimeSerializer) Deserialize(data []byte) (ptypes.TypeTime, error) {
	if len(data) < 8 {
		return ptypes.TypeTime{}, fmt.Errorf("time deserialize: expected 8 bytes, got %d", len(data))
	}
	return ptypes.TypeTime{BufferPtr: data[:8]}, nil
}

func (s *TimeSerializer) Validate(value ptypes.TypeTime) error {
	if len(value.BufferPtr) < 8 {
		return fmt.Errorf("time validate: buffer too short, expected 8, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *TimeSerializer) GetType() ptypes.OID { return ptypes.PTypeTime }
