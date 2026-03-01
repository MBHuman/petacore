// serializers/timestampz.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"time"
)

type TimestampzSerializer struct{}

var TimestampzSerializerInstance BaseSerializer[*time.Time, ptypes.TypeTimestampz] = &TimestampzSerializer{}

func (s *TimestampzSerializer) Serialize(allocator pmem.Allocator, value *time.Time) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("timestampz serialize: value is nil")
	}

	// всегда конвертируем в UTC — это главное отличие от timestamp
	utc := value.UTC()
	usec := utc.Sub(ptypes.PgEpoch).Microseconds()

	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return nil, fmt.Errorf("timestampz serialize: %w", err)
	}

	binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
	return buf, nil
}

func (s *TimestampzSerializer) Deserialize(data []byte) (ptypes.TypeTimestampz, error) {
	if len(data) < 8 {
		return ptypes.TypeTimestampz{}, fmt.Errorf("timestampz deserialize: expected 8 bytes, got %d", len(data))
	}
	return ptypes.TypeTimestampz{BufferPtr: data[:8]}, nil
}

func (s *TimestampzSerializer) Validate(value ptypes.TypeTimestampz) error {
	if len(value.BufferPtr) < 8 {
		return fmt.Errorf("timestampz validate: buffer too short, expected 8, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *TimestampzSerializer) GetType() ptypes.OID { return ptypes.PTypeTimestampz }
