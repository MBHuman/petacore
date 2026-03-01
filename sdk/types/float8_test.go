package ptypes_test

import (
	"math"
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeFloat8(t *testing.T, arena pmem.Allocator, val float64) ptypes.TypeFloat8 {
	t.Helper()
	buf, err := serializers.Float8SerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize float8: %v", err)
	}
	result, err := serializers.Float8SerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize float8: %v", err)
	}
	return result
}

const float8Eps = 1e-12

func assertFloat8Near(t *testing.T, got, expected float64) {
	t.Helper()
	if math.Abs(got-expected) > float8Eps {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeFloat8_SerializeDeserialize(t *testing.T) {
	arena := newArena()

	cases := []float64{0, 1.0, -1.0, 3.14159265358979, -3.14159265358979, 1e100, -1e100}
	for _, v := range cases {
		result := makeFloat8(t, arena, v)
		assertFloat8Near(t, result.IntoGo(), v)
	}
}

func TestTypeFloat8_SerializeDeserialize_Extremes(t *testing.T) {
	arena := newArena()

	inf := makeFloat8(t, arena, math.Inf(1))
	if !math.IsInf(inf.IntoGo(), 1) {
		t.Fatal("expected +Inf")
	}

	negInf := makeFloat8(t, arena, math.Inf(-1))
	if !math.IsInf(negInf.IntoGo(), -1) {
		t.Fatal("expected -Inf")
	}
}

func TestTypeFloat8_GetType(t *testing.T) {
	arena := newArena()
	f := makeFloat8(t, arena, 1.0)

	if f.GetType() != ptypes.PTypeFloat8 {
		t.Fatalf("expected PTypeFloat8, got %d", f.GetType())
	}
}

func TestTypeFloat8_IntoGo_ShortBuffer(t *testing.T) {
	f := ptypes.TypeFloat8{BufferPtr: []byte{0x01, 0x02}}
	if f.IntoGo() != 0 {
		t.Fatal("expected 0 for short buffer")
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeFloat8_IsNull(t *testing.T) {
	f := ptypes.TypeFloat8{BufferPtr: nil}
	if !f.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if f.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeFloat8_IsNotNull(t *testing.T) {
	arena := newArena()
	f := makeFloat8(t, arena, 1.0)

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

func TestTypeFloat8_Compare(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -1.0)
	b := makeFloat8(t, arena, 1.0)
	c := makeFloat8(t, arena, -1.0)

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

func TestTypeFloat8_Compare_NegativePositive(t *testing.T) {
	arena := newArena()
	neg := makeFloat8(t, arena, -1e10)
	zero := makeFloat8(t, arena, 0.0)
	pos := makeFloat8(t, arena, 1e10)

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

func TestTypeFloat8_Compare_Infinity(t *testing.T) {
	arena := newArena()
	negInf := makeFloat8(t, arena, math.Inf(-1))
	posInf := makeFloat8(t, arena, math.Inf(1))
	large := makeFloat8(t, arena, 1e308)

	if negInf.Compare(large) >= 0 {
		t.Fatal("expected -Inf < large")
	}
	if large.Compare(posInf) >= 0 {
		t.Fatal("expected large < +Inf")
	}
}

func TestTypeFloat8_LessThan(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -1.0)
	b := makeFloat8(t, arena, 1.0)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeFloat8_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -1.0)
	b := makeFloat8(t, arena, 1.0)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeFloat8_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 1.0)
	b := makeFloat8(t, arena, 2.0)

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

func TestTypeFloat8_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 1.0)
	b := makeFloat8(t, arena, 2.0)

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

func TestTypeFloat8_Between(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -1.0)
	b := makeFloat8(t, arena, 0.0)
	c := makeFloat8(t, arena, 1.0)

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

func TestTypeFloat8_Add(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 1.5)
	b := makeFloat8(t, arena, 2.5)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), 4.0)
}

func TestTypeFloat8_Add_Negative(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -1.5)
	b := makeFloat8(t, arena, -2.5)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), -4.0)
}

func TestTypeFloat8_Sub(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 5.0)
	b := makeFloat8(t, arena, 3.0)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), 2.0)
}

func TestTypeFloat8_Sub_Negative(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -5.0)
	b := makeFloat8(t, arena, 3.0)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), -8.0)
}

func TestTypeFloat8_Mul(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 3.0)
	b := makeFloat8(t, arena, 4.0)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), 12.0)
}

func TestTypeFloat8_Mul_NegativePositive(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -3.0)
	b := makeFloat8(t, arena, 4.0)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), -12.0)
}

func TestTypeFloat8_Div(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 10.0)
	b := makeFloat8(t, arena, 4.0)

	result, err := a.Div(arena, b)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), 2.5)
}

func TestTypeFloat8_Div_ByZero(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 10.0)
	b := makeFloat8(t, arena, 0.0)

	_, err := a.Div(arena, b)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestTypeFloat8_Mod(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 10.0)
	b := makeFloat8(t, arena, 3.0)

	result, err := a.Mod(arena, b)
	if err != nil {
		t.Fatalf("Mod: %v", err)
	}
	assertFloat8Near(t, result.IntoGo(), 1.0)
}

func TestTypeFloat8_Mod_ByZero(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 10.0)
	b := makeFloat8(t, arena, 0.0)

	_, err := a.Mod(arena, b)
	if err == nil {
		t.Fatal("expected modulo by zero error")
	}
}

// ============================================================
// IsZero / Neg / Abs
// ============================================================

func TestTypeFloat8_IsZero(t *testing.T) {
	arena := newArena()

	zero := makeFloat8(t, arena, 0.0)
	if !zero.IsZero() {
		t.Fatal("expected IsZero=true")
	}

	nonZero := makeFloat8(t, arena, 1.0)
	if nonZero.IsZero() {
		t.Fatal("expected IsZero=false")
	}
}

func TestTypeFloat8_Neg(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 3.14159265358979)

	result := a.Neg(arena)
	assertFloat8Near(t, result.IntoGo(), -3.14159265358979)
}

func TestTypeFloat8_Neg_AlreadyNegative(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -3.14159265358979)

	result := a.Neg(arena)
	assertFloat8Near(t, result.IntoGo(), 3.14159265358979)
}

func TestTypeFloat8_Abs_Positive(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 3.14159265358979)

	result := a.Abs(arena)
	assertFloat8Near(t, result.IntoGo(), 3.14159265358979)
}

func TestTypeFloat8_Abs_Negative(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, -3.14159265358979)

	result := a.Abs(arena)
	assertFloat8Near(t, result.IntoGo(), 3.14159265358979)
}

func TestTypeFloat8_Abs_Zero(t *testing.T) {
	arena := newArena()
	a := makeFloat8(t, arena, 0.0)

	result := a.Abs(arena)
	assertFloat8Near(t, result.IntoGo(), 0.0)
}

// ============================================================
// OOM
// ============================================================

func TestTypeFloat8_Add_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeFloat8{BufferPtr: make([]byte, 8)}
	b := ptypes.TypeFloat8{BufferPtr: make([]byte, 8)}

	_, err := a.Add(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeFloat8_Neg_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeFloat8{BufferPtr: make([]byte, 8)}

	result := a.Neg(arena)
	if result == nil {
		t.Fatal("expected non-nil result even on OOM")
	}
}
