// types/date.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	"time"
)

var PgEpoch = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

type TypeDate struct {
	BufferPtr []byte
}

var _ BaseType[*time.Time] = (*TypeDate)(nil)
var _ OrderedType[*time.Time] = (*TypeDate)(nil)
var _ NullableType[*time.Time] = (*TypeDate)(nil)

func (t TypeDate) GetType() OID { return PTypeDate }

func (t TypeDate) Compare(other BaseType[*time.Time]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeDate) GetBuffer() []byte { return t.BufferPtr }

// daysFromTime вычисляет количество дней от PgEpoch до tm
func daysFromTime(tm time.Time) int32 {
	return int32(tm.UTC().Truncate(24*time.Hour).Sub(PgEpoch).Hours() / 24)
}

func (t TypeDate) IntoGo() *time.Time {
	if len(t.BufferPtr) < 4 {
		return nil
	}
	// читаем order-preserving значение и восстанавливаем знаковое int32
	raw := binary.BigEndian.Uint32(t.BufferPtr)
	days := int32(raw ^ 0x80000000)
	result := PgEpoch.AddDate(0, 0, int(days))
	return &result
}

// days читает количество дней из буфера
func (t TypeDate) days() int32 {
	if len(t.BufferPtr) < 4 {
		return 0
	}
	raw := binary.BigEndian.Uint32(t.BufferPtr)
	return int32(raw ^ 0x80000000)
}

// dateFromDays создаёт TypeDate из количества дней от PgEpoch
func dateFromDays(allocator pmem.Allocator, days int32) (TypeDate, error) {
	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return TypeDate{}, fmt.Errorf("date: alloc failed: %w", err)
	}
	// XOR с 0x80000000 для order-preserving: отрицательные < нулевых < положительных
	binary.BigEndian.PutUint32(buf, uint32(days)^0x80000000)
	return TypeDate{BufferPtr: buf}, nil
}

// NewTypeDate создаёт TypeDate из time.Time без аллокатора — для тестов и одноразовых значений
func NewTypeDate(allocator pmem.Allocator, tm time.Time) (TypeDate, error) {
	return dateFromDays(allocator, daysFromTime(tm))
}

// NullableType

func (t TypeDate) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeDate) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeDate) LessThan(other BaseType[*time.Time]) bool       { return t.Compare(other) < 0 }
func (t TypeDate) GreaterThan(other BaseType[*time.Time]) bool    { return t.Compare(other) > 0 }
func (t TypeDate) LessOrEqual(other BaseType[*time.Time]) bool    { return t.Compare(other) <= 0 }
func (t TypeDate) GreaterOrEqual(other BaseType[*time.Time]) bool { return t.Compare(other) >= 0 }
func (t TypeDate) Between(low, high BaseType[*time.Time]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// AddDays добавляет количество дней
func (t TypeDate) AddDays(allocator pmem.Allocator, days int32) (TypeDate, error) {
	return dateFromDays(allocator, t.days()+days)
}

// SubDays вычитает количество дней
func (t TypeDate) SubDays(allocator pmem.Allocator, days int32) (TypeDate, error) {
	return dateFromDays(allocator, t.days()-days)
}

// DiffDays возвращает разницу в днях — только читает
func (t TypeDate) DiffDays(other TypeDate) int32 {
	return t.days() - other.days()
}

// AddMonths добавляет месяцы
func (t TypeDate) AddMonths(allocator pmem.Allocator, months int) (TypeDate, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeDate{}, fmt.Errorf("date: nil buffer")
	}
	result := tm.AddDate(0, months, 0)
	return dateFromDays(allocator, daysFromTime(result))
}

// AddYears добавляет годы
func (t TypeDate) AddYears(allocator pmem.Allocator, years int) (TypeDate, error) {
	tm := t.IntoGo()
	if tm == nil {
		return TypeDate{}, fmt.Errorf("date: nil buffer")
	}
	result := tm.AddDate(years, 0, 0)
	return dateFromDays(allocator, daysFromTime(result))
}

// только читают

func (t TypeDate) Year() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Year()
}

func (t TypeDate) Month() time.Month {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Month()
}

func (t TypeDate) Day() int {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Day()
}

func (t TypeDate) Weekday() time.Weekday {
	tm := t.IntoGo()
	if tm == nil {
		return 0
	}
	return tm.Weekday()
}

func (t TypeDate) String() string {
	tm := t.IntoGo()
	if tm == nil {
		return "date(NULL)"
	}
	return "date(" + tm.Format("2006-01-02") + ")"
}
