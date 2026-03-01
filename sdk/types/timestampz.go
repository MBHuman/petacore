// types/timestampz.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	"time"
)

type TypeTimestampz struct {
	BufferPtr []byte
}

var _ BaseType[*time.Time] = (*TypeTimestampz)(nil)
var _ OrderedType[*time.Time] = (*TypeTimestampz)(nil)
var _ NullableType[*time.Time] = (*TypeTimestampz)(nil)

func (t TypeTimestampz) GetType() OID { return PTypeTimestampz }

func (t TypeTimestampz) Compare(other BaseType[*time.Time]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeTimestampz) GetBuffer() []byte { return t.BufferPtr }

func (t TypeTimestampz) IntoGo() *time.Time {
	if len(t.BufferPtr) < 8 {
		return nil
	}
	usec := int64(binary.BigEndian.Uint64(t.BufferPtr) ^ 0x8000000000000000)
	result := PgEpoch.Add(time.Duration(usec) * time.Microsecond)
	return &result
}

// NullableType

func (t TypeTimestampz) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeTimestampz) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeTimestampz) LessThan(other BaseType[*time.Time]) bool       { return t.Compare(other) < 0 }
func (t TypeTimestampz) GreaterThan(other BaseType[*time.Time]) bool    { return t.Compare(other) > 0 }
func (t TypeTimestampz) LessOrEqual(other BaseType[*time.Time]) bool    { return t.Compare(other) <= 0 }
func (t TypeTimestampz) GreaterOrEqual(other BaseType[*time.Time]) bool { return t.Compare(other) >= 0 }
func (t TypeTimestampz) Between(low, high BaseType[*time.Time]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func (t TypeTimestampz) usec() int64 {
	if len(t.BufferPtr) < 8 {
		return 0
	}
	return int64(binary.BigEndian.Uint64(t.BufferPtr) ^ 0x8000000000000000)
}

func timestampzFromUsec(allocator pmem.Allocator, usec int64) (TypeTimestampz, error) {
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return TypeTimestampz{}, fmt.Errorf("timestampz: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
	return TypeTimestampz{BufferPtr: buf}, nil
}

// InLocation конвертирует UTC timestamp в заданный timezone — только читает
func (t TypeTimestampz) InLocation(loc *time.Location) *time.Time {
	tm := t.IntoGo()
	if tm == nil {
		return nil
	}
	result := tm.In(loc)
	return &result
}

// операции — аналогичны TypeTimestamp но всегда работают в UTC

func (t TypeTimestampz) AddMicroseconds(allocator pmem.Allocator, usec int64) (TypeTimestampz, error) {
	return timestampzFromUsec(allocator, t.usec()+usec)
}

func (t TypeTimestampz) AddSeconds(allocator pmem.Allocator, sec int64) (TypeTimestampz, error) {
	return timestampzFromUsec(allocator, t.usec()+sec*1_000_000)
}

func (t TypeTimestampz) AddMinutes(allocator pmem.Allocator, min int64) (TypeTimestampz, error) {
	return timestampzFromUsec(allocator, t.usec()+min*60_000_000)
}

func (t TypeTimestampz) AddHours(allocator pmem.Allocator, hours int64) (TypeTimestampz, error) {
	return timestampzFromUsec(allocator, t.usec()+hours*3_600_000_000)
}

func (t TypeTimestampz) AddDays(allocator pmem.Allocator, days int64) (TypeTimestampz, error) {
	return timestampzFromUsec(allocator, t.usec()+days*86_400_000_000)
}

func (t TypeTimestampz) AddMonths(allocator pmem.Allocator, months int) (TypeTimestampz, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeTimestampz{}, fmt.Errorf("timestampz: nil buffer")
	}
	result := tm.UTC().AddDate(0, months, 0)
	usec := result.Sub(PgEpoch).Microseconds()
	return timestampzFromUsec(allocator, usec)
}

func (t TypeTimestampz) AddYears(allocator pmem.Allocator, years int) (TypeTimestampz, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeTimestampz{}, fmt.Errorf("timestampz: nil buffer")
	}
	result := tm.UTC().AddDate(years, 0, 0)
	usec := result.Sub(PgEpoch).Microseconds()
	return timestampzFromUsec(allocator, usec)
}

func (t TypeTimestampz) DiffMicroseconds(other TypeTimestampz) int64 {
	return t.usec() - other.usec()
}

func (t TypeTimestampz) Truncate(allocator pmem.Allocator, d time.Duration) (TypeTimestampz, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeTimestampz{}, fmt.Errorf("timestampz: nil buffer")
	}
	result := tm.UTC().Truncate(d)
	usec := result.Sub(PgEpoch).Microseconds()
	return timestampzFromUsec(allocator, usec)
}

// только читают

func (t TypeTimestampz) Year() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.UTC().Year()
}

func (t TypeTimestampz) Month() time.Month {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.UTC().Month()
}

func (t TypeTimestampz) Day() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.UTC().Day()
}

func (t TypeTimestampz) Hour() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.UTC().Hour()
}

func (t TypeTimestampz) Minute() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.UTC().Minute()
}

func (t TypeTimestampz) Second() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.UTC().Second()
}

func (t TypeTimestampz) Weekday() time.Weekday {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.UTC().Weekday()
}
