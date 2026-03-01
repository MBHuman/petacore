package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeBytea(t *testing.T, arena pmem.Allocator, val []byte) ptypes.TypeBytea {
	t.Helper()
	buf, err := serializers.BytesSerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize bytea: %v", err)
	}
	result, err := serializers.BytesSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize bytea: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeBytea_SerializeDeserialize(t *testing.T) {
	arena := newArena()
	data := []byte{0x01, 0x02, 0x03, 0x04}
	b := makeBytea(t, arena, data)

	if string(b.IntoGo()) != string(data) {
		t.Fatalf("expected %v, got %v", data, b.IntoGo())
	}
}

func TestTypeBytea_SerializeDeserialize_Empty(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{})

	if len(b.IntoGo()) != 0 {
		t.Fatalf("expected empty, got %v", b.IntoGo())
	}
}

func TestTypeBytea_SerializeDeserialize_Nil(t *testing.T) {
	arena := newArena()

	buf, err := serializers.BytesSerializerInstance.Serialize(arena, []byte(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf != nil {
		t.Fatalf("expected nil buf for nil value, got %v", buf)
	}
}

func TestTypeBytea_GetType(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{1, 2, 3})

	if b.GetType() != ptypes.PTypeBytea {
		t.Fatalf("expected PTypeBytea, got %d", b.GetType())
	}
}

func TestTypeBytea_GetBuffer(t *testing.T) {
	arena := newArena()
	data := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	b := makeBytea(t, arena, data)

	buf := b.GetBuffer()
	if string(buf) != string(data) {
		t.Fatalf("GetBuffer mismatch: expected %v, got %v", data, buf)
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeBytea_IsNull(t *testing.T) {
	b := ptypes.TypeBytea{BufferPtr: nil}
	if !b.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if b.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeBytea_IsNotNull(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{1})

	if b.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !b.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeBytea_Compare(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01})
	b := makeBytea(t, arena, []byte{0x02})
	c := makeBytea(t, arena, []byte{0x01})

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

func TestTypeBytea_LessThan(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01})
	b := makeBytea(t, arena, []byte{0x02})

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeBytea_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01})
	b := makeBytea(t, arena, []byte{0x02})

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeBytea_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01})
	b := makeBytea(t, arena, []byte{0x02})

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

func TestTypeBytea_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01})
	b := makeBytea(t, arena, []byte{0x02})

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

func TestTypeBytea_Between(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01})
	b := makeBytea(t, arena, []byte{0x02})
	c := makeBytea(t, arena, []byte{0x03})

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
// Concat
// ============================================================

func TestTypeBytea_Concat(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01, 0x02})
	b := makeBytea(t, arena, []byte{0x03, 0x04})

	result, err := a.Concat(arena, b)
	if err != nil {
		t.Fatalf("concat: %v", err)
	}

	expected := []byte{0x01, 0x02, 0x03, 0x04}
	if string(result.IntoGo()) != string(expected) {
		t.Fatalf("expected %v, got %v", expected, result.IntoGo())
	}
}

func TestTypeBytea_Concat_Empty(t *testing.T) {
	arena := newArena()
	a := makeBytea(t, arena, []byte{0x01, 0x02})
	b := makeBytea(t, arena, []byte{})

	result, err := a.Concat(arena, b)
	if err != nil {
		t.Fatalf("concat with empty: %v", err)
	}
	if string(result.IntoGo()) != string(a.IntoGo()) {
		t.Fatalf("expected unchanged, got %v", result.IntoGo())
	}
}

func TestTypeBytea_Concat_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	a := ptypes.TypeBytea{BufferPtr: []byte{0x01}}
	b := ptypes.TypeBytea{BufferPtr: []byte{0x02}}

	_, err := a.Concat(arena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// Slice
// ============================================================

func TestTypeBytea_Slice(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03, 0x04})

	result, err := b.Slice(1, 2)
	if err != nil {
		t.Fatalf("slice: %v", err)
	}

	expected := []byte{0x02, 0x03}
	if string(result.IntoGo()) != string(expected) {
		t.Fatalf("expected %v, got %v", expected, result.IntoGo())
	}
}

func TestTypeBytea_Slice_ZeroCopy(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03, 0x04})

	result, _ := b.Slice(0, 4)

	// zero copy — должен указывать на тот же буфер
	if pmem.BufAddr(result.GetBuffer()) != pmem.BufAddr(b.GetBuffer()) {
		t.Fatal("expected zero-copy slice to point to same buffer")
	}
}

func TestTypeBytea_Slice_OutOfBounds(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02})

	_, err := b.Slice(0, 10)
	if err == nil {
		t.Fatal("expected out of bounds error")
	}
}

func TestTypeBytea_Slice_NegativeStart(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02})

	_, err := b.Slice(-1, 1)
	if err == nil {
		t.Fatal("expected error for negative start")
	}
}

// ============================================================
// Length
// ============================================================

func TestTypeBytea_Length(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03})

	if b.Length() != 3 {
		t.Fatalf("expected length 3, got %d", b.Length())
	}
}

func TestTypeBytea_Length_Empty(t *testing.T) {
	b := ptypes.TypeBytea{BufferPtr: []byte{}}
	if b.Length() != 0 {
		t.Fatalf("expected length 0, got %d", b.Length())
	}
}

// ============================================================
// Contains / StartsWith / EndsWith
// ============================================================

func TestTypeBytea_Contains(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03, 0x04})

	if !b.Contains([]byte{0x02, 0x03}) {
		t.Fatal("expected Contains=true")
	}
	if b.Contains([]byte{0x05}) {
		t.Fatal("expected Contains=false")
	}
}

func TestTypeBytea_StartsWith(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03})

	if !b.StartsWith([]byte{0x01, 0x02}) {
		t.Fatal("expected StartsWith=true")
	}
	if b.StartsWith([]byte{0x02}) {
		t.Fatal("expected StartsWith=false")
	}
}

func TestTypeBytea_EndsWith(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03})

	if !b.EndsWith([]byte{0x02, 0x03}) {
		t.Fatal("expected EndsWith=true")
	}
	if b.EndsWith([]byte{0x01}) {
		t.Fatal("expected EndsWith=false")
	}
}

// ============================================================
// Overlay
// ============================================================

func TestTypeBytea_Overlay(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03, 0x04})

	result, err := b.Overlay(arena, []byte{0xAA, 0xBB}, 1, 2)
	if err != nil {
		t.Fatalf("overlay: %v", err)
	}

	expected := []byte{0x01, 0xAA, 0xBB, 0x04}
	if string(result.IntoGo()) != string(expected) {
		t.Fatalf("expected %v, got %v", expected, result.IntoGo())
	}
}

func TestTypeBytea_Overlay_ShrinkReplacement(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02, 0x03, 0x04})

	// заменяем 2 байта на 1
	result, err := b.Overlay(arena, []byte{0xFF}, 1, 2)
	if err != nil {
		t.Fatalf("overlay: %v", err)
	}

	expected := []byte{0x01, 0xFF, 0x04}
	if string(result.IntoGo()) != string(expected) {
		t.Fatalf("expected %v, got %v", expected, result.IntoGo())
	}
}

func TestTypeBytea_Overlay_OutOfBounds(t *testing.T) {
	arena := newArena()
	b := makeBytea(t, arena, []byte{0x01, 0x02})

	_, err := b.Overlay(arena, []byte{0xFF}, 1, 5)
	if err == nil {
		t.Fatal("expected out of bounds error")
	}
}

func TestTypeBytea_Overlay_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	b := ptypes.TypeBytea{BufferPtr: []byte{0x01, 0x02, 0x03}}

	_, err := b.Overlay(arena, []byte{0xFF}, 0, 1)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}
