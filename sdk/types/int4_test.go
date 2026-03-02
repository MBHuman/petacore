package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeInt4(t *testing.T, arena pmem.Allocator, val int32) ptypes.TypeInt4 {
	t.Helper()
	buf, err := serializers.Int4SerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize int4: %v", err)
	}
	result, err := serializers.Int4SerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize int4: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeInt4_SerializeDeserialize(t *testing.T) {
	arena := newArena()

	cases := []int32{0, 1, -1, 1000, -1000, 2147483647, -2147483648}
	for _, v := range cases {
		result := makeInt4(t, arena, v)
		if result.IntoGo() != v {
			t.Fatalf("expected %d, got %d", v, result.IntoGo())
		}
	}
}

func TestTypeInt4_GetType(t *testing.T) {
	arena := newArena()
	v := makeInt4(t, arena, 1)

	if v.GetType() != ptypes.PTypeInt4 {
		t.Fatalf("expected PTypeInt4, got %d", v.GetType())
	}
}

func TestTypeInt4_IntoGo_ShortBuffer(t *testing.T) {
	v := ptypes.TypeInt4{BufferPtr: []byte{0x01, 0x02}}
	if v.IntoGo() != 0 {
		t.Fatal("expected 0 for short buffer")
	}
}

func TestTypeInt4_GetBuffer(t *testing.T) {
	arena := newArena()
	v := makeInt4(t, arena, 42)

	buf := v.GetBuffer()
	if len(buf) != 4 {
		t.Fatalf("expected buffer len 4, got %d", len(buf))
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeInt4_IsNull(t *testing.T) {
	v := ptypes.TypeInt4{BufferPtr: nil}
	if !v.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if v.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeInt4_IsNotNull(t *testing.T) {
	arena := newArena()
	v := makeInt4(t, arena, 1)

	if v.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !v.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeInt4_Compare(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -1)
	b := makeInt4(t, arena, 1)
	c := makeInt4(t, arena, -1)

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

func TestTypeInt4_Compare_Boundaries(t *testing.T) {
	arena := newArena()
	min := makeInt4(t, arena, -2147483648)
	zero := makeInt4(t, arena, 0)
	max := makeInt4(t, arena, 2147483647)

	if min.Compare(zero) >= 0 {
		t.Fatal("expected min < zero")
	}
	if zero.Compare(max) >= 0 {
		t.Fatal("expected zero < max")
	}
	if min.Compare(max) >= 0 {
		t.Fatal("expected min < max")
	}
}

func TestTypeInt4_LessThan(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -10)
	b := makeInt4(t, arena, 10)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeInt4_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -10)
	b := makeInt4(t, arena, 10)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeInt4_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 1)
	b := makeInt4(t, arena, 2)

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

func TestTypeInt4_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 1)
	b := makeInt4(t, arena, 2)

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

func TestTypeInt4_Between(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -100)
	b := makeInt4(t, arena, 0)
	c := makeInt4(t, arena, 100)

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
// NumericType — Add / Sub / Mul / Div / Mod
// ============================================================

func TestTypeInt4_Add(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 1000)
	b := makeInt4(t, arena, 2000)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.IntoGo() != 3000 {
		t.Fatalf("expected 3000, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Add_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -1000)
	b := makeInt4(t, arena, -2000)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.IntoGo() != -3000 {
		t.Fatalf("expected -3000, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Sub(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 5000)
	b := makeInt4(t, arena, 3000)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	if result.IntoGo() != 2000 {
		t.Fatalf("expected 2000, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Sub_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -5)
	b := makeInt4(t, arena, 3)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	if result.IntoGo() != -8 {
		t.Fatalf("expected -8, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Mul(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 1000)
	b := makeInt4(t, arena, 1000)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	if result.IntoGo() != 1000000 {
		t.Fatalf("expected 1000000, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Mul_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -1000)
	b := makeInt4(t, arena, 1000)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	if result.IntoGo() != -1000000 {
		t.Fatalf("expected -1000000, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Div(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 100)
	b := makeInt4(t, arena, 7)

	result, err := a.Div(arena, b)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	if result.IntoGo() != 14 {
		t.Fatalf("expected 14, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Div_ByZero(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 10)
	b := makeInt4(t, arena, 0)

	_, err := a.Div(arena, b)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestTypeInt4_Mod(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 100)
	b := makeInt4(t, arena, 7)

	result, err := a.Mod(arena, b)
	if err != nil {
		t.Fatalf("Mod: %v", err)
	}
	if result.IntoGo() != 2 {
		t.Fatalf("expected 2, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Mod_ByZero(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 10)
	b := makeInt4(t, arena, 0)

	_, err := a.Mod(arena, b)
	if err == nil {
		t.Fatal("expected modulo by zero error")
	}
}

// ============================================================
// IsZero / Neg / Abs
// ============================================================

func TestTypeInt4_IsZero(t *testing.T) {
	arena := newArena()

	zero := makeInt4(t, arena, 0)
	if !zero.IsZero() {
		t.Fatal("expected IsZero=true")
	}

	nonZero := makeInt4(t, arena, 1)
	if nonZero.IsZero() {
		t.Fatal("expected IsZero=false")
	}
}

func TestTypeInt4_Neg(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 12345)

	result := a.Neg(arena)
	if result.IntoGo() != -12345 {
		t.Fatalf("expected -12345, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Neg_AlreadyNegative(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -12345)

	result := a.Neg(arena)
	if result.IntoGo() != 12345 {
		t.Fatalf("expected 12345, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Abs_Positive(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 12345)

	result := a.Abs(arena)
	if result.IntoGo() != 12345 {
		t.Fatalf("expected 12345, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Abs_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, -12345)

	result := a.Abs(arena)
	if result.IntoGo() != 12345 {
		t.Fatalf("expected 12345, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Abs_Zero(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 0)

	result := a.Abs(arena)
	if result.IntoGo() != 0 {
		t.Fatalf("expected 0, got %d", result.IntoGo())
	}
}

// ============================================================
// BitwiseType
// ============================================================

func TestTypeInt4_And(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 0b11001100)
	b := makeInt4(t, arena, 0b10101010)

	result := a.And(arena, b)
	if result.IntoGo() != 0b10001000 {
		t.Fatalf("expected 0b10001000=136, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Or(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 0b11001100)
	b := makeInt4(t, arena, 0b10101010)

	result := a.Or(arena, b)
	if result.IntoGo() != 0b11101110 {
		t.Fatalf("expected 0b11101110=238, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Xor(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 0b11001100)
	b := makeInt4(t, arena, 0b10101010)

	result := a.Xor(arena, b)
	if result.IntoGo() != 0b01100110 {
		t.Fatalf("expected 0b01100110=102, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Not(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 0)

	result := a.Not(arena)
	if result.IntoGo() != -1 {
		t.Fatalf("expected -1, got %d", result.IntoGo())
	}
}

func TestTypeInt4_Not_NonZero(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 1)

	result := a.Not(arena)
	if result.IntoGo() != -2 {
		t.Fatalf("expected -2, got %d", result.IntoGo())
	}
}

func TestTypeInt4_ShiftLeft(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 1)

	result := a.ShiftLeft(arena, 10)
	if result.IntoGo() != 1024 {
		t.Fatalf("expected 1024, got %d", result.IntoGo())
	}
}

func TestTypeInt4_ShiftRight(t *testing.T) {
	arena := newArena()
	a := makeInt4(t, arena, 1024)

	result := a.ShiftRight(arena, 10)
	if result.IntoGo() != 1 {
		t.Fatalf("expected 1, got %d", result.IntoGo())
	}
}

// ============================================================
// OOM
// ============================================================

func TestTypeInt4_Add_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeInt4{BufferPtr: make([]byte, 4)}
	b := ptypes.TypeInt4{BufferPtr: make([]byte, 4)}

	_, err := a.Add(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeInt4_And_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeInt4{BufferPtr: make([]byte, 4)}
	b := ptypes.TypeInt4{BufferPtr: make([]byte, 4)}

	result := a.And(arena, b)
	if result == nil {
		t.Fatal("expected non-nil result even on OOM")
	}
}
