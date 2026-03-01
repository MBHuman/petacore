package ptypes_test

import (
	"testing"
	"time"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeTimestamp(t *testing.T, arena pmem.Allocator, year int, month time.Month, day, hour, minute, second int) ptypes.TypeTimestamp {
	t.Helper()
	tm := time.Date(year, month, day, hour, minute, second, 0, time.UTC)
	buf, err := serializers.TimestampSerializerInstance.Serialize(arena, &tm)
	if err != nil {
		t.Fatalf("serialize timestamp: %v", err)
	}
	result, err := serializers.TimestampSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize timestamp: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeTimestamp_SerializeDeserialize(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2024, time.June, 15, 13, 30, 45)

	tm := ts.IntoGo()
	if tm == nil {
		t.Fatal("expected non-nil time")
	}
	if tm.Year() != 2024 || tm.Month() != time.June || tm.Day() != 15 {
		t.Fatalf("expected 2024-06-15, got %v", tm)
	}
	if tm.Hour() != 13 || tm.Minute() != 30 || tm.Second() != 45 {
		t.Fatalf("expected 13:30:45, got %v", tm)
	}
}

func TestTypeTimestamp_SerializeDeserialize_PgEpoch(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2000, time.January, 1, 0, 0, 0)

	tm := ts.IntoGo()
	if !tm.Equal(ptypes.PgEpoch) {
		t.Fatalf("expected PgEpoch, got %v", tm)
	}
}

func TestTypeTimestamp_SerializeDeserialize_BeforeEpoch(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 1999, time.December, 31, 23, 59, 59)

	tm := ts.IntoGo()
	if tm.Year() != 1999 || tm.Month() != time.December || tm.Day() != 31 {
		t.Fatalf("expected 1999-12-31, got %v", tm)
	}
	if tm.Hour() != 23 || tm.Minute() != 59 || tm.Second() != 59 {
		t.Fatalf("expected 23:59:59, got %v", tm)
	}
}

func TestTypeTimestamp_GetType(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	if ts.GetType() != ptypes.PTypeTimestamp {
		t.Fatalf("expected PTypeTimestamp, got %d", ts.GetType())
	}
}

func TestTypeTimestamp_IntoGo_ShortBuffer(t *testing.T) {
	ts := ptypes.TypeTimestamp{BufferPtr: []byte{0x01}}
	if ts.IntoGo() != nil {
		t.Fatal("expected nil for short buffer")
	}
}

func TestTypeTimestamp_GetBuffer(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	buf := ts.GetBuffer()
	if len(buf) != 8 {
		t.Fatalf("expected buffer len 8, got %d", len(buf))
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeTimestamp_IsNull(t *testing.T) {
	ts := ptypes.TypeTimestamp{BufferPtr: nil}
	if !ts.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if ts.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeTimestamp_IsNotNull(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	if ts.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !ts.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeTimestamp_Compare(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)
	c := makeTimestamp(t, arena, 2020, time.January, 1, 0, 0, 0)

	if a.Compare(b) >= 0 {
		t.Fatal("expected a < b")
	}
	if b.Compare(a) <= 0 {
		t.Fatal("expected b > a")
	}
	if a.Compare(c) != 0 {
		t.Fatal("expected a == c")
	}
}

func TestTypeTimestamp_Compare_SameDay_DifferentTime(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2024, time.June, 15, 8, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.June, 15, 20, 0, 0)

	if a.Compare(b) >= 0 {
		t.Fatal("expected morning < evening")
	}
}

func TestTypeTimestamp_LessThan(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeTimestamp_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeTimestamp_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	if !a.LessOrEqual(b) {
		t.Fatal("expected a <= b")
	}
	if !a.LessOrEqual(a) {
		t.Fatal("expected a <= a")
	}
	if b.LessOrEqual(a) {
		t.Fatal("expected b not <= a")
	}
}

func TestTypeTimestamp_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	if !b.GreaterOrEqual(a) {
		t.Fatal("expected b >= a")
	}
	if !b.GreaterOrEqual(b) {
		t.Fatal("expected b >= b")
	}
	if a.GreaterOrEqual(b) {
		t.Fatal("expected a not >= b")
	}
}

func TestTypeTimestamp_Between(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2022, time.June, 15, 12, 0, 0)
	c := makeTimestamp(t, arena, 2024, time.December, 31, 23, 59, 59)

	if !b.Between(a, c) {
		t.Fatal("expected b between a and c")
	}
	if !a.Between(a, c) {
		t.Fatal("expected a between a and c (inclusive)")
	}
	if !c.Between(a, c) {
		t.Fatal("expected c between a and c (inclusive)")
	}
	if a.Between(b, c) {
		t.Fatal("expected a not between b and c")
	}
}

// ============================================================
// Year / Month / Day / Hour / Minute / Second / Weekday
// ============================================================

func TestTypeTimestamp_DateComponents(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2024, time.March, 15, 0, 0, 0)

	if ts.Year() != 2024 {
		t.Fatalf("expected Year=2024, got %d", ts.Year())
	}
	if ts.Month() != time.March {
		t.Fatalf("expected Month=March, got %v", ts.Month())
	}
	if ts.Day() != 15 {
		t.Fatalf("expected Day=15, got %d", ts.Day())
	}
}

func TestTypeTimestamp_TimeComponents(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 13, 30, 45)

	if ts.Hour() != 13 {
		t.Fatalf("expected Hour=13, got %d", ts.Hour())
	}
	if ts.Minute() != 30 {
		t.Fatalf("expected Minute=30, got %d", ts.Minute())
	}
	if ts.Second() != 45 {
		t.Fatalf("expected Second=45, got %d", ts.Second())
	}
}

func TestTypeTimestamp_Weekday(t *testing.T) {
	arena := newArena()
	// 2024-01-01 — понедельник
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)

	if ts.Weekday() != time.Monday {
		t.Fatalf("expected Monday, got %v", ts.Weekday())
	}
}

func TestTypeTimestamp_Components_NilBuffer(t *testing.T) {
	ts := ptypes.TypeTimestamp{BufferPtr: nil}

	if ts.Year() != 0 {
		t.Fatalf("expected Year=0, got %d", ts.Year())
	}
	if ts.Month() != 0 {
		t.Fatalf("expected Month=0, got %v", ts.Month())
	}
	if ts.Day() != 0 {
		t.Fatalf("expected Day=0, got %d", ts.Day())
	}
	if ts.Hour() != 0 {
		t.Fatalf("expected Hour=0, got %d", ts.Hour())
	}
	if ts.Minute() != 0 {
		t.Fatalf("expected Minute=0, got %d", ts.Minute())
	}
	if ts.Second() != 0 {
		t.Fatalf("expected Second=0, got %d", ts.Second())
	}
}

// ============================================================
// AddMicroseconds / AddSeconds / AddMinutes / AddHours / AddDays
// ============================================================

func TestTypeTimestamp_AddMicroseconds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 12, 0, 0)

	result, err := ts.AddMicroseconds(arena, 500_000)
	if err != nil {
		t.Fatalf("AddMicroseconds: %v", err)
	}
	tm := result.IntoGo()
	if tm.Nanosecond() != 500_000*1000 {
		t.Fatalf("expected 500000 usec, got nanoseconds=%d", tm.Nanosecond())
	}
}

func TestTypeTimestamp_AddSeconds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 12, 0, 0)

	result, err := ts.AddSeconds(arena, 90)
	if err != nil {
		t.Fatalf("AddSeconds: %v", err)
	}
	if result.Minute() != 1 || result.Second() != 30 {
		t.Fatalf("expected 12:01:30, got %d:%d:%d", result.Hour(), result.Minute(), result.Second())
	}
}

func TestTypeTimestamp_AddMinutes(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 12, 0, 0)

	result, err := ts.AddMinutes(arena, 75)
	if err != nil {
		t.Fatalf("AddMinutes: %v", err)
	}
	if result.Hour() != 13 || result.Minute() != 15 {
		t.Fatalf("expected 13:15, got %d:%d", result.Hour(), result.Minute())
	}
}

func TestTypeTimestamp_AddHours(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 1, 20, 0, 0)

	result, err := ts.AddHours(arena, 5)
	if err != nil {
		t.Fatalf("AddHours: %v", err)
	}
	// переход через полночь
	if result.Day() != 2 || result.Hour() != 1 {
		t.Fatalf("expected 2024-01-02 01:00, got day=%d hour=%d", result.Day(), result.Hour())
	}
}

func TestTypeTimestamp_AddDays(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 15, 12, 0, 0)

	result, err := ts.AddDays(arena, 20)
	if err != nil {
		t.Fatalf("AddDays: %v", err)
	}
	if result.Month() != time.February || result.Day() != 4 {
		t.Fatalf("expected 2024-02-04, got %v-%d", result.Month(), result.Day())
	}
}

func TestTypeTimestamp_AddDays_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.February, 1, 0, 0, 0)

	result, err := ts.AddDays(arena, -1)
	if err != nil {
		t.Fatalf("AddDays negative: %v", err)
	}
	if result.Month() != time.January || result.Day() != 31 {
		t.Fatalf("expected 2024-01-31, got %v-%d", result.Month(), result.Day())
	}
}

// ============================================================
// AddMonths / AddYears
// ============================================================

func TestTypeTimestamp_AddMonths(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 15, 12, 0, 0)

	result, err := ts.AddMonths(arena, 3)
	if err != nil {
		t.Fatalf("AddMonths: %v", err)
	}
	if result.Month() != time.April || result.Day() != 15 {
		t.Fatalf("expected 2024-04-15, got %v-%d", result.Month(), result.Day())
	}
}

func TestTypeTimestamp_AddMonths_YearRollover(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.November, 1, 0, 0, 0)

	result, err := ts.AddMonths(arena, 3)
	if err != nil {
		t.Fatalf("AddMonths year rollover: %v", err)
	}
	if result.Year() != 2025 || result.Month() != time.February {
		t.Fatalf("expected 2025-02, got %d-%v", result.Year(), result.Month())
	}
}

func TestTypeTimestamp_AddYears(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2020, time.March, 15, 10, 30, 0)

	result, err := ts.AddYears(arena, 4)
	if err != nil {
		t.Fatalf("AddYears: %v", err)
	}
	if result.Year() != 2024 || result.Month() != time.March || result.Day() != 15 {
		t.Fatalf("expected 2024-03-15, got %d-%v-%d", result.Year(), result.Month(), result.Day())
	}
	// время должно сохраниться
	if result.Hour() != 10 || result.Minute() != 30 {
		t.Fatalf("expected time 10:30, got %d:%d", result.Hour(), result.Minute())
	}
}

func TestTypeTimestamp_AddYears_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.June, 1, 0, 0, 0)

	result, err := ts.AddYears(arena, -10)
	if err != nil {
		t.Fatalf("AddYears negative: %v", err)
	}
	if result.Year() != 2014 {
		t.Fatalf("expected 2014, got %d", result.Year())
	}
}

// ============================================================
// DiffMicroseconds
// ============================================================

func TestTypeTimestamp_DiffMicroseconds(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 1, 1, 0, 0)

	diff := b.DiffMicroseconds(a)
	if diff != 3_600_000_000 {
		t.Fatalf("expected 3600000000, got %d", diff)
	}
}

func TestTypeTimestamp_DiffMicroseconds_Negative(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 1, 1, 0, 0)

	diff := a.DiffMicroseconds(b)
	if diff != -3_600_000_000 {
		t.Fatalf("expected -3600000000, got %d", diff)
	}
}

func TestTypeTimestamp_DiffMicroseconds_Days(t *testing.T) {
	arena := newArena()
	a := makeTimestamp(t, arena, 2024, time.January, 1, 0, 0, 0)
	b := makeTimestamp(t, arena, 2024, time.January, 2, 0, 0, 0)

	diff := b.DiffMicroseconds(a)
	if diff != 86_400_000_000 {
		t.Fatalf("expected 86400000000, got %d", diff)
	}
}

func TestTypeTimestamp_DiffMicroseconds_Same(t *testing.T) {
	arena := newArena()
	ts := makeTimestamp(t, arena, 2024, time.June, 15, 12, 0, 0)

	diff := ts.DiffMicroseconds(ts)
	if diff != 0 {
		t.Fatalf("expected 0, got %d", diff)
	}
}

// ============================================================
// Truncate
// ============================================================

func TestTypeTimestamp_Truncate_Hour(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 15, 13, 45, 30)

	result, err := ts.Truncate(arena, time.Hour)
	if err != nil {
		t.Fatalf("Truncate hour: %v", err)
	}
	if result.Hour() != 13 || result.Minute() != 0 || result.Second() != 0 {
		t.Fatalf("expected 13:00:00, got %d:%d:%d", result.Hour(), result.Minute(), result.Second())
	}
}

func TestTypeTimestamp_Truncate_Day(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 15, 13, 45, 30)

	result, err := ts.Truncate(arena, 24*time.Hour)
	if err != nil {
		t.Fatalf("Truncate day: %v", err)
	}
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Fatalf("expected 00:00:00 after day truncate, got %d:%d:%d",
			result.Hour(), result.Minute(), result.Second())
	}
}

func TestTypeTimestamp_Truncate_Minute(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, arena, 2024, time.January, 15, 13, 45, 30)

	result, err := ts.Truncate(arena, time.Minute)
	if err != nil {
		t.Fatalf("Truncate minute: %v", err)
	}
	if result.Minute() != 45 || result.Second() != 0 {
		t.Fatalf("expected :45:00, got :%d:%d", result.Minute(), result.Second())
	}
}

// ============================================================
// OOM
// ============================================================

func TestTypeTimestamp_AddSeconds_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	ts := ptypes.TypeTimestamp{BufferPtr: make([]byte, 8)}

	_, err := ts.AddSeconds(arena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeTimestamp_AddMonths_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, bigArena, 2024, time.January, 1, 0, 0, 0)

	oomArena := pmem.NewArena(0)
	_, err := ts.AddMonths(oomArena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeTimestamp_Truncate_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	ts := makeTimestamp(t, bigArena, 2024, time.January, 1, 12, 30, 0)

	oomArena := pmem.NewArena(0)
	_, err := ts.Truncate(oomArena, time.Hour)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}
