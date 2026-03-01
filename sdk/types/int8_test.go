package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeInt8(t *testing.T, arena pmem.Allocator, val int64) ptypes.TypeInt8 {
	t.Helper()
	buf, err := serializers.Int8SerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize int8: %v", err)
	}
	result, err := serializers.Int8SerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize int8: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeInt8_SerializeDeserialize(t *testing.T) {
	arena := newArena()

	cases := []int64{0, 1, -1, 1000000, -1000000, 9223372036854775807, -9223372036854775808}
	for _, v := range cases {
		result := makeInt8(t, arena, v)
		if result.IntoGo() != v {
			t.Fatalf("expected %d, got %d", v, result.IntoGo())
		}
	}
}

func TestTypeInt8_GetType(t *testing.T) {
	arena := newArena()
	v := makeInt8(t, arena, 1)

	if v.GetType() != ptypes.PTypeInt8 {
		t.Fatalf("expected PTypeInt8, got %d", v.GetType())
	}
}

func TestTypeInt8_IntoGo_ShortBuffer(t *testing.T) {
	v := ptypes.TypeInt8{BufferPtr: []byte{0x01, 0x02, 0x03}}
	if v.IntoGo() != 0 {
		t.Fatal("expected 0 for short buffer")
	}
}

func TestTypeInt8_GetBuffer(t *testing.T) {
	arena := newArena()
	v := makeInt8(t, arena, 42)

	buf := v.GetBuffer()
	if len(buf) != 8 {
		t.Fatalf("expected buffer len 8, got %d", len(buf))
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeInt8_IsNull(t *testing.T) {
	v := ptypes.TypeInt8{BufferPtr: nil}
	if !v.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if v.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeInt8_IsNotNull(t *testing.T) {
	arena := newArena()
	v := makeInt8(t, arena, 1)

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

func TestTypeInt8_Compare(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -1)
	b := makeInt8(t, arena, 1)
	c := makeInt8(t, arena, -1)

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

func TestTypeInt8_Compare_Boundaries(t *testing.T) {
	arena := newArena()
	min := makeInt8(t, arena, -9223372036854775808)
	zero := makeInt8(t, arena, 0)
	max := makeInt8(t, arena, 9223372036854775807)

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

func TestTypeInt8_LessThan(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -1000000)
	b := makeInt8(t, arena, 1000000)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeInt8_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -1000000)
	b := makeInt8(t, arena, 1000000)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeInt8_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1)
	b := makeInt8(t, arena, 2)

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

func TestTypeInt8_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1)
	b := makeInt8(t, arena, 2)

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

func TestTypeInt8_Between(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -1000000)
	b := makeInt8(t, arena, 0)
	c := makeInt8(t, arena, 1000000)

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

func TestTypeInt8_Add(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1000000000)
	b := makeInt8(t, arena, 2000000000)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.IntoGo() != 3000000000 {
		t.Fatalf("expected 3000000000, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Add_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -1000000000)
	b := makeInt8(t, arena, -2000000000)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if result.IntoGo() != -3000000000 {
		t.Fatalf("expected -3000000000, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Sub(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 5000000000)
	b := makeInt8(t, arena, 3000000000)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	if result.IntoGo() != 2000000000 {
		t.Fatalf("expected 2000000000, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Sub_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -5)
	b := makeInt8(t, arena, 3)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	if result.IntoGo() != -8 {
		t.Fatalf("expected -8, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Mul(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1000000)
	b := makeInt8(t, arena, 1000000)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	if result.IntoGo() != 1000000000000 {
		t.Fatalf("expected 1000000000000, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Mul_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -1000000)
	b := makeInt8(t, arena, 1000000)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	if result.IntoGo() != -1000000000000 {
		t.Fatalf("expected -1000000000000, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Div(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1000000000000)
	b := makeInt8(t, arena, 1000000)

	result, err := a.Div(arena, b)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	if result.IntoGo() != 1000000 {
		t.Fatalf("expected 1000000, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Div_ByZero(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 10)
	b := makeInt8(t, arena, 0)

	_, err := a.Div(arena, b)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestTypeInt8_Mod(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1000000000007)
	b := makeInt8(t, arena, 1000000000)

	result, err := a.Mod(arena, b)
	if err != nil {
		t.Fatalf("Mod: %v", err)
	}
	if result.IntoGo() != 7 {
		t.Fatalf("expected 7, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Mod_ByZero(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 10)
	b := makeInt8(t, arena, 0)

	_, err := a.Mod(arena, b)
	if err == nil {
		t.Fatal("expected modulo by zero error")
	}
}

// ============================================================
// IsZero / Neg / Abs
// ============================================================

func TestTypeInt8_IsZero(t *testing.T) {
	arena := newArena()

	zero := makeInt8(t, arena, 0)
	if !zero.IsZero() {
		t.Fatal("expected IsZero=true")
	}

	nonZero := makeInt8(t, arena, 1)
	if nonZero.IsZero() {
		t.Fatal("expected IsZero=false")
	}
}

func TestTypeInt8_Neg(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 9223372036854775806)

	result := a.Neg(arena)
	if result.IntoGo() != -9223372036854775806 {
		t.Fatalf("expected -9223372036854775806, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Neg_AlreadyNegative(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -9223372036854775806)

	result := a.Neg(arena)
	if result.IntoGo() != 9223372036854775806 {
		t.Fatalf("expected 9223372036854775806, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Abs_Positive(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 9223372036854775806)

	result := a.Abs(arena)
	if result.IntoGo() != 9223372036854775806 {
		t.Fatalf("expected 9223372036854775806, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Abs_Negative(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, -9223372036854775806)

	result := a.Abs(arena)
	if result.IntoGo() != 9223372036854775806 {
		t.Fatalf("expected 9223372036854775806, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Abs_Zero(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 0)

	result := a.Abs(arena)
	if result.IntoGo() != 0 {
		t.Fatalf("expected 0, got %d", result.IntoGo())
	}
}

// ============================================================
// BitwiseType
// ============================================================

func TestTypeInt8_And(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 0x00FF00FF00FF00FF)
	b := makeInt8(t, arena, 0x0F0F0F0F0F0F0F0F)

	result := a.And(arena, b)
	if result.IntoGo() != 0x000F000F000F000F {
		t.Fatalf("expected 0x000F000F000F000F, got %x", result.IntoGo())
	}
}

func TestTypeInt8_Or(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 0x00FF00FF00FF00FF)
	b := makeInt8(t, arena, 0x0F0F0F0F0F0F0F0F)

	result := a.Or(arena, b)
	if result.IntoGo() != 0x0FFF0FFF0FFF0FFF {
		t.Fatalf("expected 0x0FFF0FFF0FFF0FFF, got %x", result.IntoGo())
	}
}

func TestTypeInt8_Xor(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 0b11001100)
	b := makeInt8(t, arena, 0b10101010)

	result := a.Xor(arena, b)
	if result.IntoGo() != 0b01100110 {
		t.Fatalf("expected 0b01100110=102, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Not(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 0)

	result := a.Not(arena)
	if result.IntoGo() != -1 {
		t.Fatalf("expected -1, got %d", result.IntoGo())
	}
}

func TestTypeInt8_Not_NonZero(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1)

	result := a.Not(arena)
	if result.IntoGo() != -2 {
		t.Fatalf("expected -2, got %d", result.IntoGo())
	}
}

func TestTypeInt8_ShiftLeft(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 1)

	result := a.ShiftLeft(arena, 32)
	if result.IntoGo() != 4294967296 {
		t.Fatalf("expected 4294967296, got %d", result.IntoGo())
	}
}

func TestTypeInt8_ShiftRight(t *testing.T) {
	arena := newArena()
	a := makeInt8(t, arena, 4294967296)

	result := a.ShiftRight(arena, 32)
	if result.IntoGo() != 1 {
		t.Fatalf("expected 1, got %d", result.IntoGo())
	}
}

// ============================================================
// OOM
// ============================================================

func TestTypeInt8_Add_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeInt8{BufferPtr: make([]byte, 8)}
	b := ptypes.TypeInt8{BufferPtr: make([]byte, 8)}

	_, err := a.Add(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeInt8_And_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeInt8{BufferPtr: make([]byte, 8)}
	b := ptypes.TypeInt8{BufferPtr: make([]byte, 8)}

	result := a.And(arena, b)
	if result == nil {
		t.Fatal("expected non-nil result even on OOM")
	}
}
