package ptypes_test

import (
	"math"
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeFloat4(t *testing.T, arena pmem.Allocator, val float32) ptypes.TypeFloat4 {
	t.Helper()
	buf, err := serializers.Float4SerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize float4: %v", err)
	}
	result, err := serializers.Float4SerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize float4: %v", err)
	}
	return result
}

const float4Eps = float32(1e-6)

func assertFloat4Near(t *testing.T, got, expected float32) {
	t.Helper()
	if math.Abs(float64(got-expected)) > float64(float4Eps) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeFloat4_SerializeDeserialize(t *testing.T) {
	arena := newArena()

	cases := []float32{0, 1.0, -1.0, 3.14, -3.14, 1e10, -1e10}
	for _, v := range cases {
		result := makeFloat4(t, arena, v)
		assertFloat4Near(t, result.IntoGo(), v)
	}
}

func TestTypeFloat4_SerializeDeserialize_Extremes(t *testing.T) {
	arena := newArena()

	inf := makeFloat4(t, arena, float32(math.Inf(1)))
	if !math.IsInf(float64(inf.IntoGo()), 1) {
		t.Fatal("expected +Inf")
	}

	negInf := makeFloat4(t, arena, float32(math.Inf(-1)))
	if !math.IsInf(float64(negInf.IntoGo()), -1) {
		t.Fatal("expected -Inf")
	}
}

func TestTypeFloat4_GetType(t *testing.T) {
	arena := newArena()
	f := makeFloat4(t, arena, 1.0)

	if f.GetType() != ptypes.PTypeFloat4 {
		t.Fatalf("expected PTypeFloat4, got %d", f.GetType())
	}
}

func TestTypeFloat4_IntoGo_ShortBuffer(t *testing.T) {
	f := ptypes.TypeFloat4{BufferPtr: []byte{0x01}}
	if f.IntoGo() != 0 {
		t.Fatal("expected 0 for short buffer")
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeFloat4_IsNull(t *testing.T) {
	f := ptypes.TypeFloat4{BufferPtr: nil}
	if !f.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if f.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeFloat4_IsNotNull(t *testing.T) {
	arena := newArena()
	f := makeFloat4(t, arena, 1.0)

	if f.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !f.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeFloat4_Compare(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -1.0)
	b := makeFloat4(t, arena, 1.0)
	c := makeFloat4(t, arena, -1.0)

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

func TestTypeFloat4_Compare_NegativePositive(t *testing.T) {
	arena := newArena()
	neg := makeFloat4(t, arena, -100.0)
	pos := makeFloat4(t, arena, 100.0)
	zero := makeFloat4(t, arena, 0.0)

	if neg.Compare(zero) >= 0 {
		t.Fatal("expected neg < zero")
	}
	if zero.Compare(pos) >= 0 {
		t.Fatal("expected zero < pos")
	}
	if neg.Compare(pos) >= 0 {
		t.Fatal("expected neg < pos")
	}
}

func TestTypeFloat4_LessThan(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -1.0)
	b := makeFloat4(t, arena, 1.0)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeFloat4_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -1.0)
	b := makeFloat4(t, arena, 1.0)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeFloat4_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 1.0)
	b := makeFloat4(t, arena, 2.0)

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

func TestTypeFloat4_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 1.0)
	b := makeFloat4(t, arena, 2.0)

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

func TestTypeFloat4_Between(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -1.0)
	b := makeFloat4(t, arena, 0.0)
	c := makeFloat4(t, arena, 1.0)

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

func TestTypeFloat4_Add(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 1.5)
	b := makeFloat4(t, arena, 2.5)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), 4.0)
}

func TestTypeFloat4_Add_Negative(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -1.5)
	b := makeFloat4(t, arena, -2.5)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), -4.0)
}

func TestTypeFloat4_Sub(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 5.0)
	b := makeFloat4(t, arena, 3.0)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), 2.0)
}

func TestTypeFloat4_Sub_Negative(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -5.0)
	b := makeFloat4(t, arena, 3.0)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), -8.0)
}

func TestTypeFloat4_Mul(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 3.0)
	b := makeFloat4(t, arena, 4.0)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), 12.0)
}

func TestTypeFloat4_Mul_NegativePositive(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -3.0)
	b := makeFloat4(t, arena, 4.0)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), -12.0)
}

func TestTypeFloat4_Div(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 10.0)
	b := makeFloat4(t, arena, 4.0)

	result, err := a.Div(arena, b)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), 2.5)
}

func TestTypeFloat4_Div_ByZero(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 10.0)
	b := makeFloat4(t, arena, 0.0)

	_, err := a.Div(arena, b)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestTypeFloat4_Mod(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 10.0)
	b := makeFloat4(t, arena, 3.0)

	result, err := a.Mod(arena, b)
	if err != nil {
		t.Fatalf("Mod: %v", err)
	}
	assertFloat4Near(t, result.IntoGo(), 1.0)
}

func TestTypeFloat4_Mod_ByZero(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 10.0)
	b := makeFloat4(t, arena, 0.0)

	_, err := a.Mod(arena, b)
	if err == nil {
		t.Fatal("expected modulo by zero error")
	}
}

// ============================================================
// IsZero / Neg / Abs
// ============================================================

func TestTypeFloat4_IsZero(t *testing.T) {
	arena := newArena()

	zero := makeFloat4(t, arena, 0.0)
	if !zero.IsZero() {
		t.Fatal("expected IsZero=true")
	}

	nonZero := makeFloat4(t, arena, 1.0)
	if nonZero.IsZero() {
		t.Fatal("expected IsZero=false")
	}
}

func TestTypeFloat4_Neg(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 3.14)

	result := a.Neg(arena)
	assertFloat4Near(t, result.IntoGo(), -3.14)
}

func TestTypeFloat4_Neg_AlreadyNegative(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -3.14)

	result := a.Neg(arena)
	assertFloat4Near(t, result.IntoGo(), 3.14)
}

func TestTypeFloat4_Abs_Positive(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 3.14)

	result := a.Abs(arena)
	assertFloat4Near(t, result.IntoGo(), 3.14)
}

func TestTypeFloat4_Abs_Negative(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, -3.14)

	result := a.Abs(arena)
	assertFloat4Near(t, result.IntoGo(), 3.14)
}

func TestTypeFloat4_Abs_Zero(t *testing.T) {
	arena := newArena()
	a := makeFloat4(t, arena, 0.0)

	result := a.Abs(arena)
	assertFloat4Near(t, result.IntoGo(), 0.0)
}

// ============================================================
// OOM
// ============================================================

func TestTypeFloat4_Add_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeFloat4{BufferPtr: make([]byte, 4)}
	b := ptypes.TypeFloat4{BufferPtr: make([]byte, 4)}

	_, err := a.Add(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeFloat4_Neg_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeFloat4{BufferPtr: make([]byte, 4)}

	result := a.Neg(arena)
	// Neg игнорирует ошибку — возвращает zero value
	if result == nil {
		t.Fatal("expected non-nil result even on OOM")
	}
}
