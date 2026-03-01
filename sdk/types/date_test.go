package ptypes_test

import (
	"testing"
	"time"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeDate(t *testing.T, arena pmem.Allocator, year int, month time.Month, day int) ptypes.TypeDate {
	t.Helper()
	tm := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	buf, err := serializers.DateSerializerInstance.Serialize(arena, &tm)
	if err != nil {
		t.Fatalf("serialize date: %v", err)
	}
	result, err := serializers.DateSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize date: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeDate_SerializeDeserialize(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.January, 15)

	tm := d.IntoGo()
	if tm == nil {
		t.Fatal("expected non-nil time")
	}
	if tm.Year() != 2024 || tm.Month() != time.January || tm.Day() != 15 {
		t.Fatalf("expected 2024-01-15, got %v", tm)
	}
}

func TestTypeDate_SerializeDeserialize_PgEpoch(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2000, time.January, 1)

	tm := d.IntoGo()
	if !tm.Equal(ptypes.PgEpoch) {
		t.Fatalf("expected pg epoch, got %v", tm)
	}
}

func TestTypeDate_SerializeDeserialize_BeforeEpoch(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 1999, time.December, 31)

	tm := d.IntoGo()
	if tm.Year() != 1999 || tm.Month() != time.December || tm.Day() != 31 {
		t.Fatalf("expected 1999-12-31, got %v", tm)
	}
}

func TestTypeDate_Validate(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.March, 1)

	if err := serializers.DateSerializerInstance.Validate(d); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestTypeDate_Validate_ShortBuffer(t *testing.T) {
	d := ptypes.TypeDate{BufferPtr: []byte{0x01}}
	if err := serializers.DateSerializerInstance.Validate(d); err == nil {
		t.Fatal("expected validation error for short buffer")
	}
}

func TestTypeDate_GetType(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.January, 1)

	if d.GetType() != ptypes.PTypeDate {
		t.Fatalf("expected PTypeDate, got %d", d.GetType())
	}
}

func TestTypeDate_IntoGo_ShortBuffer(t *testing.T) {
	d := ptypes.TypeDate{BufferPtr: []byte{0x01}}
	if d.IntoGo() != nil {
		t.Fatal("expected nil for short buffer")
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeDate_IsNull(t *testing.T) {
	d := ptypes.TypeDate{BufferPtr: nil}
	if !d.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if d.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeDate_IsNotNull(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.January, 1)

	if d.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !d.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeDate_Compare(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2024, time.January, 1)
	d2 := makeDate(t, arena, 2024, time.June, 15)
	d3 := makeDate(t, arena, 2024, time.January, 1)

	if d1.Compare(d2) >= 0 {
		t.Fatal("expected d1 < d2")
	}
	if d2.Compare(d1) <= 0 {
		t.Fatal("expected d2 > d1")
	}
	if d1.Compare(d3) != 0 {
		t.Fatal("expected d1 == d3")
	}
}

func TestTypeDate_LessThan(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2020, time.January, 1)
	d2 := makeDate(t, arena, 2024, time.January, 1)

	if !d1.LessThan(d2) {
		t.Fatal("expected d1 < d2")
	}
	if d2.LessThan(d1) {
		t.Fatal("expected d2 not < d1")
	}
}

func TestTypeDate_GreaterThan(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2020, time.January, 1)
	d2 := makeDate(t, arena, 2024, time.January, 1)

	if !d2.GreaterThan(d1) {
		t.Fatal("expected d2 > d1")
	}
	if d1.GreaterThan(d2) {
		t.Fatal("expected d1 not > d2")
	}
}

func TestTypeDate_LessOrEqual(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2020, time.January, 1)
	d2 := makeDate(t, arena, 2024, time.January, 1)

	if !d1.LessOrEqual(d2) {
		t.Fatal("expected d1 <= d2")
	}
	if !d1.LessOrEqual(d1) {
		t.Fatal("expected d1 <= d1")
	}
	if d2.LessOrEqual(d1) {
		t.Fatal("expected d2 not <= d1")
	}
}

func TestTypeDate_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2020, time.January, 1)
	d2 := makeDate(t, arena, 2024, time.January, 1)

	if !d2.GreaterOrEqual(d1) {
		t.Fatal("expected d2 >= d1")
	}
	if !d2.GreaterOrEqual(d2) {
		t.Fatal("expected d2 >= d2")
	}
	if d1.GreaterOrEqual(d2) {
		t.Fatal("expected d1 not >= d2")
	}
}

func TestTypeDate_Between(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2020, time.January, 1)
	d2 := makeDate(t, arena, 2022, time.June, 15)
	d3 := makeDate(t, arena, 2024, time.December, 31)

	if !d2.Between(d1, d3) {
		t.Fatal("expected d2 between d1 and d3")
	}
	if !d1.Between(d1, d3) {
		t.Fatal("expected d1 between d1 and d3 (inclusive)")
	}
	if !d3.Between(d1, d3) {
		t.Fatal("expected d3 between d1 and d3 (inclusive)")
	}
	if d1.Between(d2, d3) {
		t.Fatal("expected d1 not between d2 and d3")
	}
}

// ============================================================
// AddDays / SubDays / DiffDays
// ============================================================

func TestTypeDate_AddDays(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.January, 1)

	result, err := d.AddDays(arena, 30)
	if err != nil {
		t.Fatalf("AddDays: %v", err)
	}

	tm := result.IntoGo()
	if tm.Month() != time.January || tm.Day() != 31 {
		t.Fatalf("expected 2024-01-31, got %v", tm)
	}
}

func TestTypeDate_AddDays_Negative(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.February, 1)

	result, err := d.AddDays(arena, -1)
	if err != nil {
		t.Fatalf("AddDays negative: %v", err)
	}

	tm := result.IntoGo()
	if tm.Month() != time.January || tm.Day() != 31 {
		t.Fatalf("expected 2024-01-31, got %v", tm)
	}
}

func TestTypeDate_SubDays(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.February, 1)

	result, err := d.SubDays(arena, 1)
	if err != nil {
		t.Fatalf("SubDays: %v", err)
	}

	tm := result.IntoGo()
	if tm.Month() != time.January || tm.Day() != 31 {
		t.Fatalf("expected 2024-01-31, got %v", tm)
	}
}

func TestTypeDate_DiffDays(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2024, time.January, 1)
	d2 := makeDate(t, arena, 2024, time.February, 1)

	diff := d2.DiffDays(d1)
	if diff != 31 {
		t.Fatalf("expected diff=31, got %d", diff)
	}
}

func TestTypeDate_DiffDays_Negative(t *testing.T) {
	arena := newArena()
	d1 := makeDate(t, arena, 2024, time.January, 1)
	d2 := makeDate(t, arena, 2024, time.February, 1)

	diff := d1.DiffDays(d2)
	if diff != -31 {
		t.Fatalf("expected diff=-31, got %d", diff)
	}
}

func TestTypeDate_AddDays_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	d := ptypes.TypeDate{BufferPtr: []byte{0x80, 0x00, 0x00, 0x00}}

	_, err := d.AddDays(arena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// AddMonths / AddYears
// ============================================================

func TestTypeDate_AddMonths(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.January, 31)

	result, err := d.AddMonths(arena, 1)
	if err != nil {
		t.Fatalf("AddMonths: %v", err)
	}

	tm := result.IntoGo()
	// январь + 1 месяц = февраль, но 31 февраля нет — Go усекает до последнего дня
	if tm.Month() != time.February && tm.Month() != time.March {
		t.Fatalf("unexpected month: %v", tm)
	}
}

func TestTypeDate_AddYears(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2020, time.March, 15)

	result, err := d.AddYears(arena, 4)
	if err != nil {
		t.Fatalf("AddYears: %v", err)
	}

	tm := result.IntoGo()
	if tm.Year() != 2024 || tm.Month() != time.March || tm.Day() != 15 {
		t.Fatalf("expected 2024-03-15, got %v", tm)
	}
}

func TestTypeDate_AddMonths_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	d := ptypes.TypeDate{BufferPtr: []byte{0x80, 0x00, 0x00, 0x00}}

	_, err := d.AddMonths(arena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// Year / Month / Day / Weekday
// ============================================================

func TestTypeDate_YearMonthDay(t *testing.T) {
	arena := newArena()
	d := makeDate(t, arena, 2024, time.June, 15)

	if d.Year() != 2024 {
		t.Fatalf("expected Year=2024, got %d", d.Year())
	}
	if d.Month() != time.June {
		t.Fatalf("expected Month=June, got %v", d.Month())
	}
	if d.Day() != 15 {
		t.Fatalf("expected Day=15, got %d", d.Day())
	}
}

func TestTypeDate_Weekday(t *testing.T) {
	arena := newArena()
	// 2024-01-01 — понедельник
	d := makeDate(t, arena, 2024, time.January, 1)

	if d.Weekday() != time.Monday {
		t.Fatalf("expected Monday, got %v", d.Weekday())
	}
}

func TestTypeDate_YearMonthDay_NilBuffer(t *testing.T) {
	d := ptypes.TypeDate{BufferPtr: nil}

	if d.Year() != 0 {
		t.Fatalf("expected Year=0 for nil buffer, got %d", d.Year())
	}
	if d.Month() != 0 {
		t.Fatalf("expected Month=0 for nil buffer, got %v", d.Month())
	}
	if d.Day() != 0 {
		t.Fatalf("expected Day=0 for nil buffer, got %d", d.Day())
	}
}
