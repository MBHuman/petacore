// types/time.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	"time"
)

type TypeTime struct {
	BufferPtr []byte
}

var _ BaseType[*time.Time] = (*TypeTime)(nil)
var _ OrderedType[*time.Time] = (*TypeTime)(nil)
var _ NullableType[*time.Time] = (*TypeTime)(nil)

func (t TypeTime) GetType() OID { return PTypeTime }

func (t TypeTime) Compare(other BaseType[*time.Time]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeTime) GetBuffer() []byte { return t.BufferPtr }

func (t TypeTime) IntoGo() *time.Time {
	if len(t.BufferPtr) < 8 {
		return nil
	}
	usec := int64(binary.BigEndian.Uint64(t.BufferPtr) ^ 0x8000000000000000)
	result := time.Date(0, 1, 1, 0, 0, 0, int(usec)*1000, time.UTC)
	return &result
}

// NullableType

func (t TypeTime) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeTime) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeTime) LessThan(other BaseType[*time.Time]) bool       { return t.Compare(other) < 0 }
func (t TypeTime) GreaterThan(other BaseType[*time.Time]) bool    { return t.Compare(other) > 0 }
func (t TypeTime) LessOrEqual(other BaseType[*time.Time]) bool    { return t.Compare(other) <= 0 }
func (t TypeTime) GreaterOrEqual(other BaseType[*time.Time]) bool { return t.Compare(other) >= 0 }
func (t TypeTime) Between(low, high BaseType[*time.Time]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func (t TypeTime) usec() int64 {
	if len(t.BufferPtr) < 8 {
		return 0
	}
	return int64(binary.BigEndian.Uint64(t.BufferPtr) ^ 0x8000000000000000)
}

func timeFromUsec(allocator pmem.Allocator, usec int64) (TypeTime, error) {
	// зажимаем в диапазон [0, 86399999999] мкс
	const maxUsec = int64(86_399_999_999)
	if usec < 0 {
		usec = 0
	}
	if usec > maxUsec {
		usec = maxUsec
	}
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return TypeTime{}, fmt.Errorf("time: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
	return TypeTime{BufferPtr: buf}, nil
}

// операции специфичные для времени

// AddMicroseconds добавляет микросекунды
func (t TypeTime) AddMicroseconds(allocator pmem.Allocator, usec int64) (TypeTime, error) {
	return timeFromUsec(allocator, t.usec()+usec)
}

// AddSeconds добавляет секунды
func (t TypeTime) AddSeconds(allocator pmem.Allocator, sec int64) (TypeTime, error) {
	return timeFromUsec(allocator, t.usec()+sec*1_000_000)
}

// AddMinutes добавляет минуты
func (t TypeTime) AddMinutes(allocator pmem.Allocator, min int64) (TypeTime, error) {
	return timeFromUsec(allocator, t.usec()+min*60_000_000)
}

// AddHours добавляет часы
func (t TypeTime) AddHours(allocator pmem.Allocator, hours int64) (TypeTime, error) {
	return timeFromUsec(allocator, t.usec()+hours*3_600_000_000)
}

// DiffMicroseconds возвращает разницу в микросекундах — только читает
func (t TypeTime) DiffMicroseconds(other TypeTime) int64 {
	return t.usec() - other.usec()
}

// Hour / Minute / Second / Microsecond — только читают

func (t TypeTime) Hour() int {
	return int(t.usec() / 3_600_000_000)
}

func (t TypeTime) Minute() int {
	return int((t.usec() % 3_600_000_000) / 60_000_000)
}

func (t TypeTime) Second() int {
	return int((t.usec() % 60_000_000) / 1_000_000)
}

func (t TypeTime) Microsecond() int {
	return int(t.usec() % 1_000_000)
}

func (t TypeTime) String() string {
	tm := t.IntoGo()
	if tm == nil {
		return "time(NULL)"
	}
	return "time(" + fmt.Sprintf("%02d:%02d:%02d.%06d", t.Hour(), t.Minute(), t.Second(), t.Microsecond()) + ")"
}

var _ CastableType[*time.Time] = (*TypeTime)(nil)

// CastableType

func (t TypeTime) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	switch targetType {
	case PTypeText:
		s := fmt.Sprintf("%02d:%02d:%02d.%06d", t.Hour(), t.Minute(), t.Second(), t.Microsecond())
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("time cast to text: %w", err)
		}
		copy(buf, s)
		return AnyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	case PTypeVarchar:
		s := fmt.Sprintf("%02d:%02d:%02d.%06d", t.Hour(), t.Minute(), t.Second(), t.Microsecond())
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("time cast to varchar: %w", err)
		}
		copy(buf, s)
		return AnyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	case PTypeTimestamp:
		// time → timestamp: комбинируем с PgEpoch (2000-01-01)
		usec := t.usec()
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("time cast to timestamp: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
		return AnyWrapper[*time.Time]{TypeTimestamp{BufferPtr: buf}}, nil

	case PTypeTimestampz:
		// time → timestampz: то же самое но интерпретируется как UTC
		usec := t.usec()
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("time cast to timestampz: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
		return AnyWrapper[*time.Time]{TypeTimestampz{BufferPtr: buf}}, nil
	default:
		return nil, fmt.Errorf("time: unsupported cast to OID %d", targetType)
	}
}

func TimeFactory(buf []byte) TypeTime {
	return TypeTime{BufferPtr: buf}
}

func TimeComparator(a, b TypeTime) int {
	return bytes.Compare(a.BufferPtr, b.BufferPtr)
}
