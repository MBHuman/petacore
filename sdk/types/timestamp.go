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
