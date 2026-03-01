package ptypes_test

import (
	"testing"
	"time"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeTimestampz(t *testing.T, arena pmem.Allocator, year int, month time.Month, day, hour, minute, second int) ptypes.TypeTimestampz {
	t.Helper()
	tm := time.Date(year, month, day, hour, minute, second, 0, time.UTC)
	buf, err := serializers.TimestampzSerializerInstance.Serialize(arena, &tm)
	if err != nil {
		t.Fatalf("serialize timestampz: %v", err)
	}
	result, err := serializers.TimestampzSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize timestampz: %v", err)
	}
	return result
}

func makeTimestampzInLoc(t *testing.T, arena pmem.Allocator, year int, month time.Month, day, hour, minute, second int, loc *time.Location) ptypes.TypeTimestampz {
	t.Helper()
	tm := time.Date(year, month, day, hour, minute, second, 0, loc)
	buf, err := serializers.TimestampzSerializerInstance.Serialize(arena, &tm)
	if err != nil {
		t.Fatalf("serialize timestampz in loc: %v", err)
	}
	result, err := serializers.TimestampzSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize timestampz in loc: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeTimestampz_SerializeDeserialize(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.June, 15, 13, 30, 45)

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

func TestTypeTimestampz_SerializeDeserialize_UTC(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

	tm := ts.IntoGo()
	if tm.Location() != time.UTC {
		t.Fatalf("expected UTC, got %v", tm.Location())
	}
}

func TestTypeTimestampz_SerializeDeserialize_PgEpoch(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2000, time.January, 1, 0, 0, 0)

	tm := ts.IntoGo()
	if !tm.Equal(ptypes.PgEpoch) {
		t.Fatalf("expected PgEpoch, got %v", tm)
	}
}

func TestTypeTimestampz_SerializeDeserialize_BeforeEpoch(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 1999, time.December, 31, 23, 59, 59)

	if ts.Year() != 1999 || ts.Month() != time.December || ts.Day() != 31 {
		t.Fatalf("expected 1999-12-31, got %d-%v-%d", ts.Year(), ts.Month(), ts.Day())
	}
}

func TestTypeTimestampz_NormalizesToUTC(t *testing.T) {
	arena := newArena()
	// UTC+3: 15:00 local = 12:00 UTC
	loc := time.FixedZone("UTC+3", 3*60*60)
	ts := makeTimestampzInLoc(t, arena, 2024, time.June, 15, 15, 0, 0, loc)

	if ts.Hour() != 12 {
		t.Fatalf("expected UTC Hour=12, got %d", ts.Hour())
	}
}

func TestTypeTimestampz_GetType(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

	if ts.GetType() != ptypes.PTypeTimestampz {
		t.Fatalf("expected PTypeTimestampz, got %d", ts.GetType())
	}
}

func TestTypeTimestampz_IntoGo_ShortBuffer(t *testing.T) {
	ts := ptypes.TypeTimestampz{BufferPtr: []byte{0x01}}
	if ts.IntoGo() != nil {
		t.Fatal("expected nil for short buffer")
	}
}

func TestTypeTimestampz_GetBuffer(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

	buf := ts.GetBuffer()
	if len(buf) != 8 {
		t.Fatalf("expected buffer len 8, got %d", len(buf))
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeTimestampz_IsNull(t *testing.T) {
	ts := ptypes.TypeTimestampz{BufferPtr: nil}
	if !ts.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if ts.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeTimestampz_IsNotNull(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

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

func TestTypeTimestampz_Compare(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)
	c := makeTimestampz(t, arena, 2020, time.January, 1, 0, 0, 0)

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

func TestTypeTimestampz_Compare_DifferentTimezones_SameInstant(t *testing.T) {
	arena := newArena()
	// одинаковый момент времени в разных зонах — после нормализации буферы одинаковы
	locPlus3 := time.FixedZone("UTC+3", 3*60*60)
	locMinus5 := time.FixedZone("UTC-5", -5*60*60)

	a := makeTimestampzInLoc(t, arena, 2024, time.June, 15, 15, 0, 0, locPlus3) // 12:00 UTC
	b := makeTimestampzInLoc(t, arena, 2024, time.June, 15, 7, 0, 0, locMinus5) // 12:00 UTC

	if a.Compare(b) != 0 {
		t.Fatal("expected same UTC instant to compare equal")
	}
}

func TestTypeTimestampz_LessThan(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeTimestampz_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeTimestampz_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

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

func TestTypeTimestampz_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

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

func TestTypeTimestampz_Between(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2020, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2022, time.June, 15, 12, 0, 0)
	c := makeTimestampz(t, arena, 2024, time.December, 31, 23, 59, 59)

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
// InLocation
// ============================================================

func TestTypeTimestampz_InLocation(t *testing.T) {
	arena := newArena()
	// 12:00 UTC
	ts := makeTimestampz(t, arena, 2024, time.June, 15, 12, 0, 0)

	locPlus5 := time.FixedZone("UTC+5", 5*60*60)
	result := ts.InLocation(locPlus5)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// UTC+5: должно быть 17:00
	if result.Hour() != 17 {
		t.Fatalf("expected Hour=17 in UTC+5, got %d", result.Hour())
	}
}

func TestTypeTimestampz_InLocation_Negative(t *testing.T) {
	arena := newArena()
	// 12:00 UTC
	ts := makeTimestampz(t, arena, 2024, time.June, 15, 12, 0, 0)

	locMinus8 := time.FixedZone("UTC-8", -8*60*60)
	result := ts.InLocation(locMinus8)
	// UTC-8: должно быть 04:00
	if result.Hour() != 4 {
		t.Fatalf("expected Hour=4 in UTC-8, got %d", result.Hour())
	}
}

func TestTypeTimestampz_InLocation_NilBuffer(t *testing.T) {
	ts := ptypes.TypeTimestampz{BufferPtr: nil}
	result := ts.InLocation(time.UTC)
	if result != nil {
		t.Fatal("expected nil for nil buffer")
	}
}

func TestTypeTimestampz_InLocation_UTC(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.June, 15, 12, 0, 0)

	result := ts.InLocation(time.UTC)
	if result.Hour() != 12 {
		t.Fatalf("expected Hour=12 in UTC, got %d", result.Hour())
	}
}

// ============================================================
// Year / Month / Day / Hour / Minute / Second / Weekday (UTC)
// ============================================================

func TestTypeTimestampz_DateComponents(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.March, 15, 0, 0, 0)

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

func TestTypeTimestampz_TimeComponents(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 13, 30, 45)

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

func TestTypeTimestampz_Weekday(t *testing.T) {
	arena := newArena()
	// 2024-01-01 — понедельник
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)

	if ts.Weekday() != time.Monday {
		t.Fatalf("expected Monday, got %v", ts.Weekday())
	}
}

func TestTypeTimestampz_Components_AlwaysUTC(t *testing.T) {
	arena := newArena()
	// UTC+6: 18:00 local = 12:00 UTC
	locPlus6 := time.FixedZone("UTC+6", 6*60*60)
	ts := makeTimestampzInLoc(t, arena, 2024, time.June, 15, 18, 0, 0, locPlus6)

	// все компоненты должны быть в UTC
	if ts.Hour() != 12 {
		t.Fatalf("expected UTC Hour=12, got %d", ts.Hour())
	}
}

func TestTypeTimestampz_Components_NilBuffer(t *testing.T) {
	ts := ptypes.TypeTimestampz{BufferPtr: nil}

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

func TestTypeTimestampz_AddMicroseconds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 12, 0, 0)

	result, err := ts.AddMicroseconds(arena, 500_000)
	if err != nil {
		t.Fatalf("AddMicroseconds: %v", err)
	}
	tm := result.IntoGo()
	if tm.Nanosecond() != 500_000*1000 {
		t.Fatalf("expected 500000 usec, got nanoseconds=%d", tm.Nanosecond())
	}
}

func TestTypeTimestampz_AddSeconds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 12, 0, 0)

	result, err := ts.AddSeconds(arena, 90)
	if err != nil {
		t.Fatalf("AddSeconds: %v", err)
	}
	if result.Minute() != 1 || result.Second() != 30 {
		t.Fatalf("expected 12:01:30, got %d:%d:%d", result.Hour(), result.Minute(), result.Second())
	}
}

func TestTypeTimestampz_AddMinutes(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 12, 0, 0)

	result, err := ts.AddMinutes(arena, 75)
	if err != nil {
		t.Fatalf("AddMinutes: %v", err)
	}
	if result.Hour() != 13 || result.Minute() != 15 {
		t.Fatalf("expected 13:15, got %d:%d", result.Hour(), result.Minute())
	}
}

func TestTypeTimestampz_AddHours(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 1, 20, 0, 0)

	result, err := ts.AddHours(arena, 5)
	if err != nil {
		t.Fatalf("AddHours: %v", err)
	}
	if result.Day() != 2 || result.Hour() != 1 {
		t.Fatalf("expected 2024-01-02 01:00 UTC, got day=%d hour=%d", result.Day(), result.Hour())
	}
}

func TestTypeTimestampz_AddDays(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 15, 12, 0, 0)

	result, err := ts.AddDays(arena, 20)
	if err != nil {
		t.Fatalf("AddDays: %v", err)
	}
	if result.Month() != time.February || result.Day() != 4 {
		t.Fatalf("expected 2024-02-04, got %v-%d", result.Month(), result.Day())
	}
}

func TestTypeTimestampz_AddDays_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.February, 1, 0, 0, 0)

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

func TestTypeTimestampz_AddMonths(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 15, 12, 0, 0)

	result, err := ts.AddMonths(arena, 3)
	if err != nil {
		t.Fatalf("AddMonths: %v", err)
	}
	if result.Month() != time.April || result.Day() != 15 {
		t.Fatalf("expected 2024-04-15, got %v-%d", result.Month(), result.Day())
	}
}

func TestTypeTimestampz_AddMonths_YearRollover(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.November, 1, 0, 0, 0)

	result, err := ts.AddMonths(arena, 3)
	if err != nil {
		t.Fatalf("AddMonths year rollover: %v", err)
	}
	if result.Year() != 2025 || result.Month() != time.February {
		t.Fatalf("expected 2025-02, got %d-%v", result.Year(), result.Month())
	}
}

func TestTypeTimestampz_AddYears(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2020, time.March, 15, 10, 30, 0)

	result, err := ts.AddYears(arena, 4)
	if err != nil {
		t.Fatalf("AddYears: %v", err)
	}
	if result.Year() != 2024 || result.Month() != time.March || result.Day() != 15 {
		t.Fatalf("expected 2024-03-15, got %d-%v-%d", result.Year(), result.Month(), result.Day())
	}
	if result.Hour() != 10 || result.Minute() != 30 {
		t.Fatalf("expected time 10:30 UTC, got %d:%d", result.Hour(), result.Minute())
	}
}

func TestTypeTimestampz_AddYears_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.June, 1, 0, 0, 0)

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

func TestTypeTimestampz_DiffMicroseconds(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2024, time.January, 1, 1, 0, 0)

	diff := b.DiffMicroseconds(a)
	if diff != 3_600_000_000 {
		t.Fatalf("expected 3600000000, got %d", diff)
	}
}

func TestTypeTimestampz_DiffMicroseconds_Negative(t *testing.T) {
	arena := newArena()
	a := makeTimestampz(t, arena, 2024, time.January, 1, 0, 0, 0)
	b := makeTimestampz(t, arena, 2024, time.January, 1, 1, 0, 0)

	diff := a.DiffMicroseconds(b)
	if diff != -3_600_000_000 {
		t.Fatalf("expected -3600000000, got %d", diff)
	}
}

func TestTypeTimestampz_DiffMicroseconds_Same(t *testing.T) {
	arena := newArena()
	ts := makeTimestampz(t, arena, 2024, time.June, 15, 12, 0, 0)

	diff := ts.DiffMicroseconds(ts)
	if diff != 0 {
		t.Fatalf("expected 0, got %d", diff)
	}
}

func TestTypeTimestampz_DiffMicroseconds_CrossTimezone(t *testing.T) {
	arena := newArena()
	// разные таймзоны, но одинаковый момент UTC — diff должен быть 0
	locPlus3 := time.FixedZone("UTC+3", 3*60*60)
	locMinus3 := time.FixedZone("UTC-3", -3*60*60)

	a := makeTimestampzInLoc(t, arena, 2024, time.June, 15, 15, 0, 0, locPlus3) // 12:00 UTC
	b := makeTimestampzInLoc(t, arena, 2024, time.June, 15, 9, 0, 0, locMinus3) // 12:00 UTC

	diff := a.DiffMicroseconds(b)
	if diff != 0 {
		t.Fatalf("expected 0 for same UTC instant, got %d", diff)
	}
}

// ============================================================
// Truncate
// ============================================================

func TestTypeTimestampz_Truncate_Hour(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 15, 13, 45, 30)

	result, err := ts.Truncate(arena, time.Hour)
	if err != nil {
		t.Fatalf("Truncate hour: %v", err)
	}
	if result.Hour() != 13 || result.Minute() != 0 || result.Second() != 0 {
		t.Fatalf("expected 13:00:00 UTC, got %d:%d:%d", result.Hour(), result.Minute(), result.Second())
	}
}

func TestTypeTimestampz_Truncate_Minute(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, arena, 2024, time.January, 15, 13, 45, 30)

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

func TestTypeTimestampz_AddSeconds_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	ts := ptypes.TypeTimestampz{BufferPtr: make([]byte, 8)}

	_, err := ts.AddSeconds(arena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeTimestampz_AddMonths_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, bigArena, 2024, time.January, 1, 0, 0, 0)

	oomArena := pmem.NewArena(0)
	_, err := ts.AddMonths(oomArena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeTimestampz_Truncate_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	ts := makeTimestampz(t, bigArena, 2024, time.January, 1, 12, 30, 0)

	oomArena := pmem.NewArena(0)
	_, err := ts.Truncate(oomArena, time.Hour)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}
