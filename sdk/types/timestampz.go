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

func (t TypeTimestampz) String() string {
	tm := t.IntoGo()
	if tm == nil {
		return "timestampz(NULL)"
	}
	return "timestampz(" + fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d %s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.usec()%1_000_000, tm.Location()) + ")"
}

var _ CastableType[*time.Time] = (*TypeTimestampz)(nil)

// CastableType

func (t TypeTimestampz) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	switch targetType {
	case PTypeDate:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestampz cast to date: nil buffer")
		}
		// дата всегда в UTC
		result, err := NewTypeDate(allocator, tm.UTC())
		if err != nil {
			return nil, fmt.Errorf("timestampz cast to date: %w", err)
		}
		return anyWrapper[*time.Time]{result}, nil

	case PTypeTime:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestampz cast to time: nil buffer")
		}
		// время суток в UTC
		utc := tm.UTC()
		usec := int64(utc.Hour())*3_600_000_000 +
			int64(utc.Minute())*60_000_000 +
			int64(utc.Second())*1_000_000 +
			int64(utc.Nanosecond()/1000)
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("timestampz cast to time: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
		return anyWrapper[*time.Time]{TypeTime{BufferPtr: buf}}, nil

	case PTypeTimestamp:
		// timestampz → timestamp: снимаем timezone, оставляем UTC момент
		// буфер идентичен — оба хранят usec от PgEpoch в UTC
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("timestampz cast to timestamp: %w", err)
		}
		copy(buf, t.BufferPtr)
		return anyWrapper[*time.Time]{TypeTimestamp{BufferPtr: buf}}, nil

	case PTypeText:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestampz cast to text: nil buffer")
		}
		utc := tm.UTC()
		s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d+00",
			utc.Year(), utc.Month(), utc.Day(),
			utc.Hour(), utc.Minute(), utc.Second(),
			utc.Nanosecond()/1000,
		)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("timestampz cast to text: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	case PTypeVarchar:
		tm := t.IntoGo()
		if tm == nil {
			return nil, fmt.Errorf("timestampz cast to varchar: nil buffer")
		}
		utc := tm.UTC()
		s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d+00",
			utc.Year(), utc.Month(), utc.Day(),
			utc.Hour(), utc.Minute(), utc.Second(),
			utc.Nanosecond()/1000,
		)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("timestampz cast to varchar: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("timestampz: unsupported cast to OID %d", targetType)
	}
}
