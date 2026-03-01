package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeInt2(t *testing.T, arena pmem.Allocator, val int16) ptypes.TypeInt2 {
	t.Helper()
	buf, err := serializers.Int2SerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize int2: %v", err)
	}
	result, err := serializers.Int2SerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize int2: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeInt2_SerializeDeserialize(t *testing.T) {
	arena := newArena()

	cases := []int16{0, 1, -1, 100, -100, 32767, -32768}
	for _, v := range cases {
		result := makeInt2(t, arena, v)
		if result.IntoGo() != v {
			t.Fatalf("expected %d, got %d", v, result.IntoGo())
		}
	}
}

func TestTypeInt2_GetType(t *testing.T) {
	arena := newArena()
	v := makeInt2(t, arena, 1)

	if v.GetType() != ptypes.PTypeInt2 {
		t.Fatalf("expected PTypeInt2, got %d", v.GetType())
	}
}

func TestTypeInt2_IntoGo_ShortBuffer(t *testing.T) {
	v := ptypes.TypeInt2{BufferPtr: []byte{0x01}}
	if v.IntoGo() != 0 {
		t.Fatal("expected 0 for short buffer")
	}
}

func TestTypeInt2_GetBuffer(t *testing.T) {
	arena := newArena()
	v := makeInt2(t, arena, 42)

	buf := v.GetBuffer()
	if len(buf) != 2 {
		t.Fatalf("expected buffer len 2, got %d", len(buf))
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeInt2_IsNull(t *testing.T) {
	v := ptypes.TypeInt2{BufferPtr: nil}
	if !v.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if v.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeInt2_IsNotNull(t *testing.T) {
	arena := newArena()
	v := makeInt2(t, arena, 1)

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

func TestTypeInt2_Compare(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -1)
	b := makeInt2(t, arena, 1)
	c := makeInt2(t, arena, -1)

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

func TestTypeInt2_Compare_Boundaries(t *testing.T) {
	arena := newArena()
	min := makeInt2(t, arena, -32768)
	zero := makeInt2(t, arena, 0)
	max := makeInt2(t, arena, 32767)

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

func TestTypeInt2_LessThan(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -10)
	b := makeInt2(t, arena, 10)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeInt2_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -10)
	b := makeInt2(t, arena, 10)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeInt2_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 1)
	b := makeInt2(t, arena, 2)

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

func TestTypeInt2_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 1)
	b := makeInt2(t, arena, 2)

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

func TestTypeInt2_Between(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -10)
	b := makeInt2(t, arena, 0)
	c := makeInt2(t, arena, 10)

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

func TestTypeInt2_Add(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 10)
	b := makeInt2(t, arena, 20)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.IntoGo() != 30 {
		t.Fatalf("expected 30, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Add_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -10)
	b := makeInt2(t, arena, -20)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.IntoGo() != -30 {
		t.Fatalf("expected -30, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Sub(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 30)
	b := makeInt2(t, arena, 10)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	if result.IntoGo() != 20 {
		t.Fatalf("expected 20, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Sub_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -5)
	b := makeInt2(t, arena, 3)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	if result.IntoGo() != -8 {
		t.Fatalf("expected -8, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Mul(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 6)
	b := makeInt2(t, arena, 7)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	if result.IntoGo() != 42 {
		t.Fatalf("expected 42, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Mul_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -6)
	b := makeInt2(t, arena, 7)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	if result.IntoGo() != -42 {
		t.Fatalf("expected -42, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Div(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 10)
	b := makeInt2(t, arena, 3)

	result, err := a.Div(arena, b)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	if result.IntoGo() != 3 {
		t.Fatalf("expected 3, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Div_ByZero(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 10)
	b := makeInt2(t, arena, 0)

	_, err := a.Div(arena, b)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestTypeInt2_Mod(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 10)
	b := makeInt2(t, arena, 3)

	result, err := a.Mod(arena, b)
	if err != nil {
		t.Fatalf("Mod: %v", err)
	}
	if result.IntoGo() != 1 {
		t.Fatalf("expected 1, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Mod_ByZero(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 10)
	b := makeInt2(t, arena, 0)

	_, err := a.Mod(arena, b)
	if err == nil {
		t.Fatal("expected modulo by zero error")
	}
}

// ============================================================
// IsZero / Neg / Abs
// ============================================================

func TestTypeInt2_IsZero(t *testing.T) {
	arena := newArena()

	zero := makeInt2(t, arena, 0)
	if !zero.IsZero() {
		t.Fatal("expected IsZero=true")
	}

	nonZero := makeInt2(t, arena, 1)
	if nonZero.IsZero() {
		t.Fatal("expected IsZero=false")
	}
}

func TestTypeInt2_Neg(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 42)

	result := a.Neg(arena)
	if result.IntoGo() != -42 {
		t.Fatalf("expected -42, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Neg_AlreadyNegative(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -42)

	result := a.Neg(arena)
	if result.IntoGo() != 42 {
		t.Fatalf("expected 42, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Abs_Positive(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 42)

	result := a.Abs(arena)
	if result.IntoGo() != 42 {
		t.Fatalf("expected 42, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Abs_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, -42)

	result := a.Abs(arena)
	if result.IntoGo() != 42 {
		t.Fatalf("expected 42, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Abs_Zero(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 0)

	result := a.Abs(arena)
	if result.IntoGo() != 0 {
		t.Fatalf("expected 0, got %d", result.IntoGo())
	}
}

// ============================================================
// BitwiseType
// ============================================================

func TestTypeInt2_And(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 0b1100)
	b := makeInt2(t, arena, 0b1010)

	result := a.And(arena, b)
	if result.IntoGo() != 0b1000 {
		t.Fatalf("expected 0b1000=8, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Or(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 0b1100)
	b := makeInt2(t, arena, 0b1010)

	result := a.Or(arena, b)
	if result.IntoGo() != 0b1110 {
		t.Fatalf("expected 0b1110=14, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Xor(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 0b1100)
	b := makeInt2(t, arena, 0b1010)

	result := a.Xor(arena, b)
	if result.IntoGo() != 0b0110 {
		t.Fatalf("expected 0b0110=6, got %d", result.IntoGo())
	}
}

func TestTypeInt2_Not(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 0)

	result := a.Not(arena)
	if result.IntoGo() != -1 {
		t.Fatalf("expected -1, got %d", result.IntoGo())
	}
}

func TestTypeInt2_ShiftLeft(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 1)

	result := a.ShiftLeft(arena, 3)
	if result.IntoGo() != 8 {
		t.Fatalf("expected 8, got %d", result.IntoGo())
	}
}

func TestTypeInt2_ShiftRight(t *testing.T) {
	arena := newArena()
	a := makeInt2(t, arena, 16)

	result := a.ShiftRight(arena, 2)
	if result.IntoGo() != 4 {
		t.Fatalf("expected 4, got %d", result.IntoGo())
	}
}

// ============================================================
// OOM
// ============================================================

func TestTypeInt2_Add_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeInt2{BufferPtr: make([]byte, 2)}
	b := ptypes.TypeInt2{BufferPtr: make([]byte, 2)}

	_, err := a.Add(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeInt2_And_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeInt2{BufferPtr: make([]byte, 2)}
	b := ptypes.TypeInt2{BufferPtr: make([]byte, 2)}

	result := a.And(arena, b)
	if result == nil {
		t.Fatal("expected non-nil result even on OOM")
	}
}
