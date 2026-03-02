package ptypes_test

import (
	"testing"
	"time"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeTime(t *testing.T, arena pmem.Allocator, hour, minute, second, microsecond int) ptypes.TypeTime {
	t.Helper()
	tm := time.Date(0, 1, 1, hour, minute, second, microsecond*1000, time.UTC)
	buf, err := serializers.TimeSerializerInstance.Serialize(arena, &tm)
	if err != nil {
		t.Fatalf("serialize time: %v", err)
	}
	result, err := serializers.TimeSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize time: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeTime_SerializeDeserialize(t *testing.T) {
	arena := newArena()
	tm := makeTime(t, arena, 13, 30, 45, 123456)

	if tm.Hour() != 13 {
		t.Fatalf("expected Hour=13, got %d", tm.Hour())
	}
	if tm.Minute() != 30 {
		t.Fatalf("expected Minute=30, got %d", tm.Minute())
	}
	if tm.Second() != 45 {
		t.Fatalf("expected Second=45, got %d", tm.Second())
	}
	if tm.Microsecond() != 123456 {
		t.Fatalf("expected Microsecond=123456, got %d", tm.Microsecond())
	}
}

func TestTypeTime_SerializeDeserialize_Midnight(t *testing.T) {
	arena := newArena()
	tm := makeTime(t, arena, 0, 0, 0, 0)

	if tm.Hour() != 0 || tm.Minute() != 0 || tm.Second() != 0 || tm.Microsecond() != 0 {
		t.Fatalf("expected midnight, got %d:%d:%d.%d", tm.Hour(), tm.Minute(), tm.Second(), tm.Microsecond())
	}
}

func TestTypeTime_SerializeDeserialize_EndOfDay(t *testing.T) {
	arena := newArena()
	tm := makeTime(t, arena, 23, 59, 59, 999999)

	if tm.Hour() != 23 {
		t.Fatalf("expected Hour=23, got %d", tm.Hour())
	}
	if tm.Minute() != 59 {
		t.Fatalf("expected Minute=59, got %d", tm.Minute())
	}
	if tm.Second() != 59 {
		t.Fatalf("expected Second=59, got %d", tm.Second())
	}
	if tm.Microsecond() != 999999 {
		t.Fatalf("expected Microsecond=999999, got %d", tm.Microsecond())
	}
}

func TestTypeTime_GetType(t *testing.T) {
	arena := newArena()
	tm := makeTime(t, arena, 12, 0, 0, 0)

	if tm.GetType() != ptypes.PTypeTime {
		t.Fatalf("expected PTypeTime, got %d", tm.GetType())
	}
}

func TestTypeTime_IntoGo_ShortBuffer(t *testing.T) {
	tm := ptypes.TypeTime{BufferPtr: []byte{0x01}}
	if tm.IntoGo() != nil {
		t.Fatal("expected nil for short buffer")
	}
}

func TestTypeTime_GetBuffer(t *testing.T) {
	arena := newArena()
	tm := makeTime(t, arena, 12, 0, 0, 0)

	buf := tm.GetBuffer()
	if len(buf) != 8 {
		t.Fatalf("expected buffer len 8, got %d", len(buf))
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeTime_IsNull(t *testing.T) {
	tm := ptypes.TypeTime{BufferPtr: nil}
	if !tm.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if tm.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeTime_IsNotNull(t *testing.T) {
	arena := newArena()
	tm := makeTime(t, arena, 12, 0, 0, 0)

	if tm.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !tm.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeTime_Compare(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 10, 0, 0, 0)
	b := makeTime(t, arena, 20, 0, 0, 0)
	c := makeTime(t, arena, 10, 0, 0, 0)

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

func TestTypeTime_Compare_Microseconds(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 12, 0, 0, 0)
	b := makeTime(t, arena, 12, 0, 0, 1)

	if a.Compare(b) >= 0 {
		t.Fatal("expected a < b (1 microsecond difference)")
	}
}

func TestTypeTime_LessThan(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 8, 0, 0, 0)
	b := makeTime(t, arena, 18, 0, 0, 0)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeTime_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 8, 0, 0, 0)
	b := makeTime(t, arena, 18, 0, 0, 0)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeTime_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 8, 0, 0, 0)
	b := makeTime(t, arena, 18, 0, 0, 0)

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

func TestTypeTime_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 8, 0, 0, 0)
	b := makeTime(t, arena, 18, 0, 0, 0)

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

func TestTypeTime_Between(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 8, 0, 0, 0)
	b := makeTime(t, arena, 12, 0, 0, 0)
	c := makeTime(t, arena, 18, 0, 0, 0)

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
// Hour / Minute / Second / Microsecond
// ============================================================

func TestTypeTime_Hour(t *testing.T) {
	arena := newArena()

	for _, h := range []int{0, 1, 12, 23} {
		tm := makeTime(t, arena, h, 0, 0, 0)
		if tm.Hour() != h {
			t.Fatalf("expected Hour=%d, got %d", h, tm.Hour())
		}
	}
}

func TestTypeTime_Minute(t *testing.T) {
	arena := newArena()

	for _, m := range []int{0, 1, 30, 59} {
		tm := makeTime(t, arena, 0, m, 0, 0)
		if tm.Minute() != m {
			t.Fatalf("expected Minute=%d, got %d", m, tm.Minute())
		}
	}
}

func TestTypeTime_Second(t *testing.T) {
	arena := newArena()

	for _, s := range []int{0, 1, 30, 59} {
		tm := makeTime(t, arena, 0, 0, s, 0)
		if tm.Second() != s {
			t.Fatalf("expected Second=%d, got %d", s, tm.Second())
		}
	}
}

func TestTypeTime_Microsecond(t *testing.T) {
	arena := newArena()

	for _, us := range []int{0, 1, 500000, 999999} {
		tm := makeTime(t, arena, 0, 0, 0, us)
		if tm.Microsecond() != us {
			t.Fatalf("expected Microsecond=%d, got %d", us, tm.Microsecond())
		}
	}
}

// ============================================================
// AddMicroseconds / AddSeconds / AddMinutes / AddHours
// ============================================================

func TestTypeTime_AddMicroseconds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	tm := makeTime(t, arena, 12, 0, 0, 0)

	result, err := tm.AddMicroseconds(arena, 500000)
	if err != nil {
		t.Fatalf("AddMicroseconds: %v", err)
	}
	if result.Microsecond() != 500000 {
		t.Fatalf("expected Microsecond=500000, got %d", result.Microsecond())
	}
}

func TestTypeTime_AddSeconds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	tm := makeTime(t, arena, 12, 0, 0, 0)

	result, err := tm.AddSeconds(arena, 30)
	if err != nil {
		t.Fatalf("AddSeconds: %v", err)
	}
	if result.Second() != 30 {
		t.Fatalf("expected Second=30, got %d", result.Second())
	}
	if result.Hour() != 12 {
		t.Fatalf("expected Hour=12, got %d", result.Hour())
	}
}

func TestTypeTime_AddMinutes(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	tm := makeTime(t, arena, 12, 0, 0, 0)

	result, err := tm.AddMinutes(arena, 45)
	if err != nil {
		t.Fatalf("AddMinutes: %v", err)
	}
	if result.Hour() != 12 || result.Minute() != 45 {
		t.Fatalf("expected 12:45, got %d:%d", result.Hour(), result.Minute())
	}
}

func TestTypeTime_AddHours(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	tm := makeTime(t, arena, 10, 0, 0, 0)

	result, err := tm.AddHours(arena, 3)
	if err != nil {
		t.Fatalf("AddHours: %v", err)
	}
	if result.Hour() != 13 {
		t.Fatalf("expected Hour=13, got %d", result.Hour())
	}
}

func TestTypeTime_AddHours_ClampAtEndOfDay(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	tm := makeTime(t, arena, 23, 0, 0, 0)

	// добавляем 2 часа — должно зажать до конца дня
	result, err := tm.AddHours(arena, 2)
	if err != nil {
		t.Fatalf("AddHours clamp: %v", err)
	}
	// максимум 23:59:59.999999
	if result.Hour() > 23 {
		t.Fatalf("expected Hour <= 23, got %d", result.Hour())
	}
}

func TestTypeTime_AddMicroseconds_ClampAtZero(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	tm := makeTime(t, arena, 0, 0, 0, 0)

	// вычитаем больше чем начало дня — должно зажать до 0
	result, err := tm.AddMicroseconds(arena, -1_000_000)
	if err != nil {
		t.Fatalf("AddMicroseconds clamp: %v", err)
	}
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 || result.Microsecond() != 0 {
		t.Fatalf("expected midnight after clamp, got %d:%d:%d.%d",
			result.Hour(), result.Minute(), result.Second(), result.Microsecond())
	}
}

// ============================================================
// DiffMicroseconds
// ============================================================

func TestTypeTime_DiffMicroseconds(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 10, 0, 0, 0)
	b := makeTime(t, arena, 11, 0, 0, 0)

	diff := b.DiffMicroseconds(a)
	if diff != 3_600_000_000 {
		t.Fatalf("expected diff=3600000000, got %d", diff)
	}
}

func TestTypeTime_DiffMicroseconds_Negative(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 10, 0, 0, 0)
	b := makeTime(t, arena, 11, 0, 0, 0)

	diff := a.DiffMicroseconds(b)
	if diff != -3_600_000_000 {
		t.Fatalf("expected diff=-3600000000, got %d", diff)
	}
}

func TestTypeTime_DiffMicroseconds_Same(t *testing.T) {
	arena := newArena()
	a := makeTime(t, arena, 12, 30, 0, 0)

	diff := a.DiffMicroseconds(a)
	if diff != 0 {
		t.Fatalf("expected diff=0, got %d", diff)
	}
}

// ============================================================
// OOM
// ============================================================

func TestTypeTime_AddSeconds_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	tm := ptypes.TypeTime{BufferPtr: make([]byte, 8)}

	_, err := tm.AddSeconds(arena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeTime_AddHours_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	tm := ptypes.TypeTime{BufferPtr: make([]byte, 8)}

	_, err := tm.AddHours(arena, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}
