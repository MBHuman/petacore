// serializers/timestamp.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"time"
)

type TimestampSerializer struct{}

var TimestampSerializerInstance BaseSerializer[*time.Time, ptypes.TypeTimestamp] = &TimestampSerializer{}

func (s *TimestampSerializer) Serialize(allocator pmem.Allocator, value *time.Time) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("timestamp serialize: value is nil")
	}

	// микросекунды от PostgreSQL epoch
	usec := value.Sub(ptypes.PgEpoch).Microseconds()

	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return nil, fmt.Errorf("timestamp serialize: %w", err)
	}

	binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
	return buf, nil
}

func (s *TimestampSerializer) Deserialize(data []byte) (ptypes.TypeTimestamp, error) {
	if len(data) < 8 {
		return ptypes.TypeTimestamp{}, fmt.Errorf("timestamp deserialize: expected 8 bytes, got %d", len(data))
	}
	return ptypes.TypeTimestamp{BufferPtr: data[:8]}, nil
}

func (s *TimestampSerializer) Validate(value ptypes.TypeTimestamp) error {
	if len(value.BufferPtr) < 8 {
		return fmt.Errorf("timestamp validate: buffer too short, expected 8, got %d", len(value.BufferPtr))
	}
	return nil
}

func (s *TimestampSerializer) GetType() ptypes.OID { return ptypes.PTypeTimestamp }
