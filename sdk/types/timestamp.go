// types/timestamp.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	"time"
)

type TypeTimestamp struct {
	BufferPtr []byte
}

var _ BaseType[*time.Time] = (*TypeTimestamp)(nil)
var _ OrderedType[*time.Time] = (*TypeTimestamp)(nil)
var _ NullableType[*time.Time] = (*TypeTimestamp)(nil)

func (t TypeTimestamp) GetType() OID { return PTypeTimestamp }

func (t TypeTimestamp) Compare(other BaseType[*time.Time]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeTimestamp) GetBuffer() []byte { return t.BufferPtr }

func (t TypeTimestamp) IntoGo() *time.Time {
	if len(t.BufferPtr) < 8 {
		return nil
	}
	usec := int64(binary.BigEndian.Uint64(t.BufferPtr) ^ 0x8000000000000000)
	result := PgEpoch.Add(time.Duration(usec) * time.Microsecond)
	return &result
}

// NullableType

func (t TypeTimestamp) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeTimestamp) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeTimestamp) LessThan(other BaseType[*time.Time]) bool       { return t.Compare(other) < 0 }
func (t TypeTimestamp) GreaterThan(other BaseType[*time.Time]) bool    { return t.Compare(other) > 0 }
func (t TypeTimestamp) LessOrEqual(other BaseType[*time.Time]) bool    { return t.Compare(other) <= 0 }
func (t TypeTimestamp) GreaterOrEqual(other BaseType[*time.Time]) bool { return t.Compare(other) >= 0 }
func (t TypeTimestamp) Between(low, high BaseType[*time.Time]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func (t TypeTimestamp) usec() int64 {
	if len(t.BufferPtr) < 8 {
		return 0
	}
	return int64(binary.BigEndian.Uint64(t.BufferPtr) ^ 0x8000000000000000)
}

func timestampFromUsec(allocator pmem.Allocator, usec int64) (TypeTimestamp, error) {
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return TypeTimestamp{}, fmt.Errorf("timestamp: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
	return TypeTimestamp{BufferPtr: buf}, nil
}

// операции специфичные для timestamp

// AddMicroseconds добавляет микросекунды
func (t TypeTimestamp) AddMicroseconds(allocator pmem.Allocator, usec int64) (TypeTimestamp, error) {
	return timestampFromUsec(allocator, t.usec()+usec)
}

// AddSeconds добавляет секунды
func (t TypeTimestamp) AddSeconds(allocator pmem.Allocator, sec int64) (TypeTimestamp, error) {
	return timestampFromUsec(allocator, t.usec()+sec*1_000_000)
}

// AddMinutes добавляет минуты
func (t TypeTimestamp) AddMinutes(allocator pmem.Allocator, min int64) (TypeTimestamp, error) {
	return timestampFromUsec(allocator, t.usec()+min*60_000_000)
}

// AddHours добавляет часы
func (t TypeTimestamp) AddHours(allocator pmem.Allocator, hours int64) (TypeTimestamp, error) {
	return timestampFromUsec(allocator, t.usec()+hours*3_600_000_000)
}

// AddDays добавляет дни
func (t TypeTimestamp) AddDays(allocator pmem.Allocator, days int64) (TypeTimestamp, error) {
	return timestampFromUsec(allocator, t.usec()+days*86_400_000_000)
}

// AddMonths добавляет месяцы — календарная арифметика
func (t TypeTimestamp) AddMonths(allocator pmem.Allocator, months int) (TypeTimestamp, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeTimestamp{}, fmt.Errorf("timestamp: nil buffer")
	}
	result := tm.AddDate(0, months, 0)
	usec := result.Sub(PgEpoch).Microseconds()
	return timestampFromUsec(allocator, usec)
}

// AddYears добавляет годы — календарная арифметика
func (t TypeTimestamp) AddYears(allocator pmem.Allocator, years int) (TypeTimestamp, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeTimestamp{}, fmt.Errorf("timestamp: nil buffer")
	}
	result := tm.AddDate(years, 0, 0)
	usec := result.Sub(PgEpoch).Microseconds()
	return timestampFromUsec(allocator, usec)
}

// DiffMicroseconds возвращает разницу в микросекундах — только читает
func (t TypeTimestamp) DiffMicroseconds(other TypeTimestamp) int64 {
	return t.usec() - other.usec()
}

// Truncate усекает до заданной точности — создаёт новый timestamp
func (t TypeTimestamp) Truncate(allocator pmem.Allocator, d time.Duration) (TypeTimestamp, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeTimestamp{}, fmt.Errorf("timestamp: nil buffer")
	}
	result := tm.Truncate(d)
	usec := result.Sub(PgEpoch).Microseconds()
	return timestampFromUsec(allocator, usec)
}

// только читают — аллокатор не нужен

func (t TypeTimestamp) Year() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Year()
}

func (t TypeTimestamp) Month() time.Month {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Month()
}

func (t TypeTimestamp) Day() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Day()
}

func (t TypeTimestamp) Hour() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Hour()
}

func (t TypeTimestamp) Minute() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Minute()
}

func (t TypeTimestamp) Second() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Second()
}

func (t TypeTimestamp) Weekday() time.Weekday {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Weekday()
}

func (t TypeTimestamp) String() string {
	tm := t.IntoGo()
	if tm == nil {
		return "timestamp(NULL)"
	}
	return "timestamp(" + fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.usec()%1_000_000) + ")"
}

var _ CastableType[*time.Time] = (*TypeTimestamp)(nil)

// CastableType

func (t TypeTimestamp) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	switch targetType {
	case PTypeDate:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestamp cast to date: nil buffer")
		}
		result, err := NewTypeDate(allocator, *tm)
		if err != nil {
			return nil, fmt.Errorf("timestamp cast to date: %w", err)
		}
		return AnyWrapper[*time.Time]{result}, nil

	case PTypeTime:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestamp cast to time: nil buffer")
		}
		// извлекаем только время суток в микросекундах
		usec := int64(tm.Hour())*3_600_000_000 +
			int64(tm.Minute())*60_000_000 +
			int64(tm.Second())*1_000_000 +
			int64(tm.Nanosecond()/1000)
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("timestamp cast to time: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
		return AnyWrapper[*time.Time]{TypeTime{BufferPtr: buf}}, nil

	case PTypeTimestampz:
		// timestamp → timestampz: интерпретируем как UTC, буфер идентичен
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("timestamp cast to timestampz: %w", err)
		}
		copy(buf, t.BufferPtr)
		return AnyWrapper[*time.Time]{TypeTimestampz{BufferPtr: buf}}, nil

	case PTypeText:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestamp cast to text: nil buffer")
		}
		s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d",
			tm.Year(), tm.Month(), tm.Day(),
			tm.Hour(), tm.Minute(), tm.Second(),
			tm.Nanosecond()/1000,
		)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("timestamp cast to text: %w", err)
		}
		copy(buf, s)
		return AnyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	case PTypeVarchar:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestamp cast to varchar: nil buffer")
		}
		s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d",
			tm.Year(), tm.Month(), tm.Day(),
			tm.Hour(), tm.Minute(), tm.Second(),
			tm.Nanosecond()/1000,
		)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("timestamp cast to varchar: %w", err)
		}
		copy(buf, s)
		return AnyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("timestamp: unsupported cast to OID %d", targetType)
	}
}

func TimestampFactory(buf []byte) TypeTimestamp {
	return TypeTimestamp{BufferPtr: buf}
}

func TimestampComparator(a, b TypeTimestamp) int {
	return bytes.Compare(a.BufferPtr, b.BufferPtr)
}
