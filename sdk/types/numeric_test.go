package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

var defaultMeta = ptypes.NumericMeta{Precision: 20, Scale: 6}

func makeNumeric(t *testing.T, arena pmem.Allocator, val string, meta ptypes.NumericMeta) ptypes.TypeNumeric {
	t.Helper()
	ser, err := serializers.NewNumericSerializer(meta.Precision, meta.Scale)
	if err != nil {
		t.Fatalf("NewNumericSerializer: %v", err)
	}
	buf, err := ser.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize numeric(%s): %v", val, err)
	}
	// Deserialize на том же сериализаторе — он знает Meta
	result, err := ser.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize numeric(%s): %v", val, err)
	}
	return result
}

func numericToString(t *testing.T, n ptypes.TypeNumeric) string {
	t.Helper()
	v, err := n.ToNumericValue()
	if err != nil {
		t.Fatalf("ToNumericValue: %v", err)
	}
	f := v.ToBigFloat()
	return f.Text('f', int(n.Meta.Scale))
}

// ============================================================
// NumericMeta Validate
// ============================================================

func TestNumericMeta_Validate(t *testing.T) {
	cases := []struct {
		meta    ptypes.NumericMeta
		wantErr bool
	}{
		{ptypes.NumericMeta{Precision: 10, Scale: 2}, false},
		{ptypes.NumericMeta{Precision: 1, Scale: 0}, false},
		{ptypes.NumericMeta{Precision: 1000, Scale: 1000}, false},
		{ptypes.NumericMeta{Precision: 0, Scale: 0}, true},    // precision < 1
		{ptypes.NumericMeta{Precision: 1001, Scale: 0}, true}, // precision > 1000
		{ptypes.NumericMeta{Precision: 10, Scale: -1}, true},  // scale < 0
		{ptypes.NumericMeta{Precision: 5, Scale: 6}, true},    // scale > precision
	}

	for _, c := range cases {
		err := c.meta.Validate()
		if c.wantErr && err == nil {
			t.Fatalf("expected error for meta %+v", c.meta)
		}
		if !c.wantErr && err != nil {
			t.Fatalf("unexpected error for meta %+v: %v", c.meta, err)
		}
	}
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeNumeric_SerializeDeserialize_Positive(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "123.456789", defaultMeta)

	got := numericToString(t, n)
	if got != "123.456789" {
		t.Fatalf("expected 123.456789, got %s", got)
	}
}

func TestTypeNumeric_SerializeDeserialize_Negative(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "-123.456789", defaultMeta)

	got := numericToString(t, n)
	if got != "-123.456789" {
		t.Fatalf("expected -123.456789, got %s", got)
	}
}

func TestTypeNumeric_SerializeDeserialize_Zero(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "0", defaultMeta)

	if !n.IsZero() {
		t.Fatal("expected IsZero=true for 0")
	}
}

func TestTypeNumeric_SerializeDeserialize_Integer(t *testing.T) {
	arena := newArena()
	meta := ptypes.NumericMeta{Precision: 10, Scale: 0}
	n := makeNumeric(t, arena, "42", meta)

	got := numericToString(t, n)
	if got != "42" {
		t.Fatalf("expected 42, got %s", got)
	}
}

func TestTypeNumeric_GetType(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "1.0", defaultMeta)

	if n.GetType() != ptypes.PTypeNumeric {
		t.Fatalf("expected PTypeNumeric, got %d", n.GetType())
	}
}

func TestTypeNumeric_GetBuffer(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "1.5", defaultMeta)

	buf := n.GetBuffer()
	if len(buf) == 0 {
		t.Fatal("expected non-empty buffer")
	}
	// первый байт — sign: 0x00=neg, 0x01=zero, 0x02=pos
	if buf[0] != 0x02 {
		t.Fatalf("expected sign byte 0x02 for positive, got 0x%02x", buf[0])
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeNumeric_IsNull(t *testing.T) {
	n := ptypes.TypeNumeric{BufferPtr: nil, Meta: defaultMeta}
	if !n.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if n.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeNumeric_IsNotNull(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "1.0", defaultMeta)

	if n.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !n.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// IsZero
// ============================================================

func TestTypeNumeric_IsZero_True(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "0", defaultMeta)

	if !n.IsZero() {
		t.Fatal("expected IsZero=true")
	}
}

func TestTypeNumeric_IsZero_False(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "0.000001", defaultMeta)

	if n.IsZero() {
		t.Fatal("expected IsZero=false for 0.000001")
	}
}

func TestTypeNumeric_IsZero_EmptyBuffer(t *testing.T) {
	n := ptypes.TypeNumeric{BufferPtr: []byte{}, Meta: defaultMeta}
	if !n.IsZero() {
		t.Fatal("expected IsZero=true for empty buffer")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeNumeric_Compare(t *testing.T) {
	arena := newArena()
	a := makeNumeric(t, arena, "-1.5", defaultMeta)
	b := makeNumeric(t, arena, "1.5", defaultMeta)
	c := makeNumeric(t, arena, "-1.5", defaultMeta)

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

func TestTypeNumeric_Compare_NegativeZeroPositive(t *testing.T) {
	arena := newArena()
	neg := makeNumeric(t, arena, "-999.999999", defaultMeta)
	zero := makeNumeric(t, arena, "0", defaultMeta)
	pos := makeNumeric(t, arena, "999.999999", defaultMeta)

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

func TestTypeNumeric_LessThan(t *testing.T) {
	arena := newArena()
	a := makeNumeric(t, arena, "-10.5", defaultMeta)
	b := makeNumeric(t, arena, "10.5", defaultMeta)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeNumeric_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeNumeric(t, arena, "-10.5", defaultMeta)
	b := makeNumeric(t, arena, "10.5", defaultMeta)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeNumeric_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeNumeric(t, arena, "1.5", defaultMeta)
	b := makeNumeric(t, arena, "2.5", defaultMeta)

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

func TestTypeNumeric_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeNumeric(t, arena, "1.5", defaultMeta)
	b := makeNumeric(t, arena, "2.5", defaultMeta)

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

func TestTypeNumeric_Between(t *testing.T) {
	arena := newArena()
	a := makeNumeric(t, arena, "-10.0", defaultMeta)
	b := makeNumeric(t, arena, "0.0", defaultMeta)
	c := makeNumeric(t, arena, "10.0", defaultMeta)

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

func TestTypeNumeric_Add(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "1.5", defaultMeta)
	b := makeNumeric(t, arena, "2.5", defaultMeta)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "4.000000" {
		t.Fatalf("expected 4.000000, got %s", got)
	}
}

func TestTypeNumeric_Add_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "-1.5", defaultMeta)
	b := makeNumeric(t, arena, "-2.5", defaultMeta)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "-4.000000" {
		t.Fatalf("expected -4.000000, got %s", got)
	}
}

func TestTypeNumeric_Add_MixedSign(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "10.0", defaultMeta)
	b := makeNumeric(t, arena, "-3.5", defaultMeta)

	result, err := a.Add(arena, b)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "6.500000" {
		t.Fatalf("expected 6.500000, got %s", got)
	}
}

func TestTypeNumeric_Sub(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "10.5", defaultMeta)
	b := makeNumeric(t, arena, "3.5", defaultMeta)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "7.000000" {
		t.Fatalf("expected 7.000000, got %s", got)
	}
}

func TestTypeNumeric_Sub_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "-5.0", defaultMeta)
	b := makeNumeric(t, arena, "3.0", defaultMeta)

	result, err := a.Sub(arena, b)
	if err != nil {
		t.Fatalf("Sub: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "-8.000000" {
		t.Fatalf("expected -8.000000, got %s", got)
	}
}

func TestTypeNumeric_Mul(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "3.0", defaultMeta)
	b := makeNumeric(t, arena, "4.0", defaultMeta)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "12.000000" {
		t.Fatalf("expected 12.000000, got %s", got)
	}
}

func TestTypeNumeric_Mul_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "-3.0", defaultMeta)
	b := makeNumeric(t, arena, "4.0", defaultMeta)

	result, err := a.Mul(arena, b)
	if err != nil {
		t.Fatalf("Mul: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "-12.000000" {
		t.Fatalf("expected -12.000000, got %s", got)
	}
}

func TestTypeNumeric_Div(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "10.0", defaultMeta)
	b := makeNumeric(t, arena, "4.0", defaultMeta)

	result, err := a.Div(arena, b)
	if err != nil {
		t.Fatalf("Div: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "2.500000" {
		t.Fatalf("expected 2.500000, got %s", got)
	}
}

func TestTypeNumeric_Div_ByZero(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "10.0", defaultMeta)
	b := makeNumeric(t, arena, "0", defaultMeta)

	_, err := a.Div(arena, b)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
}

func TestTypeNumeric_Mod(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	meta := ptypes.NumericMeta{Precision: 10, Scale: 0}
	a := makeNumeric(t, arena, "10", meta)
	b := makeNumeric(t, arena, "3", meta)

	result, err := a.Mod(arena, b)
	if err != nil {
		t.Fatalf("Mod: %v", err)
	}
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "1" {
		t.Fatalf("expected 1, got %s", got)
	}
}

func TestTypeNumeric_Mod_ByZero(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "10.0", defaultMeta)
	b := makeNumeric(t, arena, "0", defaultMeta)

	_, err := a.Mod(arena, b)
	if err == nil {
		t.Fatal("expected modulo by zero error")
	}
}

// ============================================================
// Neg / Abs
// ============================================================

func TestTypeNumeric_Neg(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "3.14", defaultMeta)

	result := a.Neg(arena)
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "-3.140000" {
		t.Fatalf("expected -3.140000, got %s", got)
	}
}

func TestTypeNumeric_Neg_AlreadyNegative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "-3.14", defaultMeta)

	result := a.Neg(arena)
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "3.140000" {
		t.Fatalf("expected 3.140000, got %s", got)
	}
}

func TestTypeNumeric_Neg_Zero(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "0", defaultMeta)

	result := a.Neg(arena)
	if !result.IsZero() {
		t.Fatal("expected neg(0) = 0")
	}
}

func TestTypeNumeric_Abs_Positive(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "3.14", defaultMeta)

	result := a.Abs(arena)
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "3.140000" {
		t.Fatalf("expected 3.140000, got %s", got)
	}
}

func TestTypeNumeric_Abs_Negative(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, arena, "-3.14", defaultMeta)

	result := a.Abs(arena)
	got := numericToString(t, result.(ptypes.TypeNumeric))
	if got != "3.140000" {
		t.Fatalf("expected 3.140000, got %s", got)
	}
}

// ============================================================
// ToNumericValue
// ============================================================

func TestTypeNumeric_ToNumericValue_Positive(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "123.456", defaultMeta)

	v, err := n.ToNumericValue()
	if err != nil {
		t.Fatalf("ToNumericValue: %v", err)
	}
	if v.Neg {
		t.Fatal("expected Neg=false for positive")
	}
	if v.Scale != defaultMeta.Scale {
		t.Fatalf("expected Scale=%d, got %d", defaultMeta.Scale, v.Scale)
	}
}

func TestTypeNumeric_ToNumericValue_Negative(t *testing.T) {
	arena := newArena()
	n := makeNumeric(t, arena, "-123.456", defaultMeta)

	v, err := n.ToNumericValue()
	if err != nil {
		t.Fatalf("ToNumericValue: %v", err)
	}
	if !v.Neg {
		t.Fatal("expected Neg=true for negative")
	}
}

func TestTypeNumeric_ToNumericValue_EmptyBuffer(t *testing.T) {
	n := ptypes.TypeNumeric{BufferPtr: []byte{}, Meta: defaultMeta}
	_, err := n.ToNumericValue()
	if err == nil {
		t.Fatal("expected error for empty buffer")
	}
}

// ============================================================
// OOM
// ============================================================

func TestTypeNumeric_Add_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	a := makeNumeric(t, bigArena, "1.5", defaultMeta)
	b := makeNumeric(t, bigArena, "2.5", defaultMeta)

	oomArena := pmem.NewArena(0)
	_, err := a.Add(oomArena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}
