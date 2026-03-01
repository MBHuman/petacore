package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func newArena() *pmem.ArenaAllocator {
	return pmem.NewArena(1024)
}

func makeBool(t *testing.T, arena pmem.Allocator, val bool) ptypes.TypeBool {
	t.Helper()
	buf, err := serializers.BoolSerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize bool: %v", err)
	}
	result, err := serializers.BoolSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize bool: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeBool_SerializeDeserialize_True(t *testing.T) {
	arena := newArena()
	b := makeBool(t, arena, true)

	if !b.IntoGo() {
		t.Fatal("expected true, got false")
	}
}

func TestTypeBool_SerializeDeserialize_False(t *testing.T) {
	arena := newArena()
	b := makeBool(t, arena, false)

	if b.IntoGo() {
		t.Fatal("expected false, got true")
	}
}

func TestTypeBool_Deserialize_EmptyBuffer(t *testing.T) {
	b, err := serializers.BoolSerializerInstance.Deserialize([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.IntoGo() {
		t.Fatal("empty buffer should deserialize as false")
	}
}

func TestTypeBool_Validate_InvalidByte(t *testing.T) {
	b := ptypes.TypeBool{BufferPtr: []byte{0x42}}
	err := serializers.BoolSerializerInstance.Validate(b)
	if err == nil {
		t.Fatal("expected validation error for invalid byte 0x42")
	}
}

func TestTypeBool_Validate_Valid(t *testing.T) {
	arena := newArena()

	for _, val := range []bool{true, false} {
		b := makeBool(t, arena, val)
		if err := serializers.BoolSerializerInstance.Validate(b); err != nil {
			t.Fatalf("validate(%v): %v", val, err)
		}
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeBool_IsNull_NilBuffer(t *testing.T) {
	b := ptypes.TypeBool{BufferPtr: nil}
	if !b.IsNull() {
		t.Fatal("expected IsNull=true for nil buffer")
	}
	if b.IsNotNull() {
		t.Fatal("expected IsNotNull=false for nil buffer")
	}
}

func TestTypeBool_IsNotNull_ValidBuffer(t *testing.T) {
	arena := newArena()
	b := makeBool(t, arena, true)

	if b.IsNull() {
		t.Fatal("expected IsNull=false for valid buffer")
	}
	if !b.IsNotNull() {
		t.Fatal("expected IsNotNull=true for valid buffer")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeBool_Compare(t *testing.T) {
	arena := newArena()
	f := makeBool(t, arena, false)
	tr := makeBool(t, arena, true)

	if f.Compare(tr) >= 0 {
		t.Fatal("expected false < true")
	}
	if tr.Compare(f) <= 0 {
		t.Fatal("expected true > false")
	}
	if f.Compare(f) != 0 {
		t.Fatal("expected false == false")
	}
	if tr.Compare(tr) != 0 {
		t.Fatal("expected true == true")
	}
}

func TestTypeBool_LessThan(t *testing.T) {
	arena := newArena()
	f := makeBool(t, arena, false)
	tr := makeBool(t, arena, true)

	if !f.LessThan(tr) {
		t.Fatal("expected false < true")
	}
	if tr.LessThan(f) {
		t.Fatal("expected true not < false")
	}
}

func TestTypeBool_GreaterThan(t *testing.T) {
	arena := newArena()
	f := makeBool(t, arena, false)
	tr := makeBool(t, arena, true)

	if !tr.GreaterThan(f) {
		t.Fatal("expected true > false")
	}
	if f.GreaterThan(tr) {
		t.Fatal("expected false not > true")
	}
}

func TestTypeBool_LessOrEqual(t *testing.T) {
	arena := newArena()
	f := makeBool(t, arena, false)
	tr := makeBool(t, arena, true)

	if !f.LessOrEqual(tr) {
		t.Fatal("expected false <= true")
	}
	if !f.LessOrEqual(f) {
		t.Fatal("expected false <= false")
	}
	if tr.LessOrEqual(f) {
		t.Fatal("expected true not <= false")
	}
}

func TestTypeBool_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	f := makeBool(t, arena, false)
	tr := makeBool(t, arena, true)

	if !tr.GreaterOrEqual(f) {
		t.Fatal("expected true >= false")
	}
	if !tr.GreaterOrEqual(tr) {
		t.Fatal("expected true >= true")
	}
	if f.GreaterOrEqual(tr) {
		t.Fatal("expected false not >= true")
	}
}

func TestTypeBool_Between(t *testing.T) {
	arena := newArena()
	f := makeBool(t, arena, false)
	tr := makeBool(t, arena, true)

	if !f.Between(f, tr) {
		t.Fatal("expected false between false and true")
	}
	if !tr.Between(f, tr) {
		t.Fatal("expected true between false and true")
	}
	if !f.Between(f, f) {
		t.Fatal("expected false between false and false")
	}
}

// ============================================================
// Logical operations
// ============================================================

func TestTypeBool_And(t *testing.T) {
	arena := newArena()

	cases := []struct {
		a, b     bool
		expected bool
	}{
		{true, true, true},
		{true, false, false},
		{false, true, false},
		{false, false, false},
	}

	for _, c := range cases {
		a := makeBool(t, arena, c.a)
		b := makeBool(t, arena, c.b)

		result, err := a.And(arena, b)
		if err != nil {
			t.Fatalf("And(%v, %v): %v", c.a, c.b, err)
		}
		if result.IntoGo() != c.expected {
			t.Fatalf("And(%v, %v): expected %v, got %v", c.a, c.b, c.expected, result.IntoGo())
		}
	}
}

func TestTypeBool_Or(t *testing.T) {
	arena := newArena()

	cases := []struct {
		a, b     bool
		expected bool
	}{
		{true, true, true},
		{true, false, true},
		{false, true, true},
		{false, false, false},
	}

	for _, c := range cases {
		a := makeBool(t, arena, c.a)
		b := makeBool(t, arena, c.b)

		result, err := a.Or(arena, b)
		if err != nil {
			t.Fatalf("Or(%v, %v): %v", c.a, c.b, err)
		}
		if result.IntoGo() != c.expected {
			t.Fatalf("Or(%v, %v): expected %v, got %v", c.a, c.b, c.expected, result.IntoGo())
		}
	}
}

func TestTypeBool_Not(t *testing.T) {
	arena := newArena()

	cases := []struct {
		input    bool
		expected bool
	}{
		{true, false},
		{false, true},
	}

	for _, c := range cases {
		b := makeBool(t, arena, c.input)
		result, err := b.Not(arena)
		if err != nil {
			t.Fatalf("Not(%v): %v", c.input, err)
		}
		if result.IntoGo() != c.expected {
			t.Fatalf("Not(%v): expected %v, got %v", c.input, c.expected, result.IntoGo())
		}
	}
}

func TestTypeBool_And_OOM(t *testing.T) {
	arena := pmem.NewArena(0) // пустая арена
	a := ptypes.TypeBool{BufferPtr: []byte{1}}
	b := ptypes.TypeBool{BufferPtr: []byte{0}}

	_, err := a.And(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeBool_Or_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeBool{BufferPtr: []byte{1}}
	b := ptypes.TypeBool{BufferPtr: []byte{0}}

	_, err := a.Or(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestTypeBool_Not_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	b := ptypes.TypeBool{BufferPtr: []byte{1}}

	_, err := b.Not(arena)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// GetBuffer / GetType
// ============================================================

func TestTypeBool_GetBuffer(t *testing.T) {
	arena := newArena()
	b := makeBool(t, arena, true)

	buf := b.GetBuffer()
	if len(buf) != 1 {
		t.Fatalf("expected buffer len 1, got %d", len(buf))
	}
	if buf[0] != 1 {
		t.Fatalf("expected buf[0]=1, got %d", buf[0])
	}
}

func TestTypeBool_GetType(t *testing.T) {
	arena := newArena()
	b := makeBool(t, arena, true)

	if b.GetType() != ptypes.PTypeBool {
		t.Fatalf("expected OID PTypeBool, got %d", b.GetType())
	}
}
