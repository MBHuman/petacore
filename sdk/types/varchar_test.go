package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeVarchar(t *testing.T, arena pmem.Allocator, val string, maxLen int32) ptypes.TypeVarchar {
	t.Helper()
	ser := serializers.NewVarcharSerializer(maxLen)
	buf, err := ser.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize varchar: %v", err)
	}
	result, err := ser.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize varchar: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeVarchar_SerializeDeserialize(t *testing.T) {
	arena := newArena()

	cases := []string{"hello", "world", "hello world", ""}
	for _, v := range cases {
		result := makeVarchar(t, arena, v, 100)
		if result.IntoGo() != v {
			t.Fatalf("expected %q, got %q", v, result.IntoGo())
		}
	}
}

func TestTypeVarchar_SerializeDeserialize_Unicode(t *testing.T) {
	arena := newArena()

	cases := []string{"привет мир", "日本語", "🎉🚀💡", "café"}
	for _, v := range cases {
		result := makeVarchar(t, arena, v, 100)
		if result.IntoGo() != v {
			t.Fatalf("expected %q, got %q", v, result.IntoGo())
		}
	}
}

func TestTypeVarchar_Serialize_ExceedsMaxLength(t *testing.T) {
	arena := newArena()
	ser := serializers.NewVarcharSerializer(5)

	_, err := ser.Serialize(arena, "hello world")
	if err == nil {
		t.Fatal("expected error for value exceeding max length")
	}
}

func TestTypeVarchar_Serialize_UnlimitedMaxLength(t *testing.T) {
	arena := newArena()
	// maxLength=0 — без ограничения
	result := makeVarchar(t, arena, "very long string that would exceed any reasonable limit", 0)
	if result.IntoGo() != "very long string that would exceed any reasonable limit" {
		t.Fatal("expected unlimited varchar to work")
	}
}

func TestTypeVarchar_GetType(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 100)

	if v.GetType() != ptypes.PTypeVarchar {
		t.Fatalf("expected PTypeVarchar, got %d", v.GetType())
	}
}

func TestTypeVarchar_GetBuffer(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 100)

	buf := v.GetBuffer()
	if len(buf) != 5 {
		t.Fatalf("expected buffer len 5, got %d", len(buf))
	}
}

func TestTypeVarchar_MaxLength(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 50)

	if v.MaxLength() != 50 {
		t.Fatalf("expected MaxLength=50, got %d", v.MaxLength())
	}
}

func TestTypeVarchar_NewTypeVarchar(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v, err := ptypes.NewTypeVarchar(arena, "hello", 10)
	if err != nil {
		t.Fatalf("NewTypeVarchar: %v", err)
	}
	if v.IntoGo() != "hello" {
		t.Fatalf("expected 'hello', got %q", v.IntoGo())
	}
}

func TestTypeVarchar_NewTypeVarchar_ExceedsMax(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	_, err := ptypes.NewTypeVarchar(arena, "hello world", 5)
	if err == nil {
		t.Fatal("expected error for value exceeding max length")
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeVarchar_IsNull(t *testing.T) {
	v := ptypes.TypeVarchar{BufferPtr: nil}
	if !v.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if v.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeVarchar_IsNotNull(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 100)

	if v.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !v.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

func TestTypeVarchar_IsNotNull_Empty(t *testing.T) {
	v := ptypes.TypeVarchar{BufferPtr: []byte{}}
	if v.IsNull() {
		t.Fatal("expected IsNull=false for empty string")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeVarchar_Compare(t *testing.T) {
	arena := newArena()
	a := makeVarchar(t, arena, "apple", 100)
	b := makeVarchar(t, arena, "banana", 100)
	c := makeVarchar(t, arena, "apple", 100)

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

func TestTypeVarchar_LessThan(t *testing.T) {
	arena := newArena()
	a := makeVarchar(t, arena, "apple", 100)
	b := makeVarchar(t, arena, "banana", 100)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeVarchar_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeVarchar(t, arena, "apple", 100)
	b := makeVarchar(t, arena, "banana", 100)

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeVarchar_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeVarchar(t, arena, "apple", 100)
	b := makeVarchar(t, arena, "banana", 100)

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

func TestTypeVarchar_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeVarchar(t, arena, "apple", 100)
	b := makeVarchar(t, arena, "banana", 100)

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

func TestTypeVarchar_Between(t *testing.T) {
	arena := newArena()
	a := makeVarchar(t, arena, "apple", 100)
	b := makeVarchar(t, arena, "mango", 100)
	c := makeVarchar(t, arena, "zebra", 100)

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
// Length / ByteLength / MaxLength
// ============================================================

func TestTypeVarchar_Length_ASCII(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 100)

	if v.Length() != 5 {
		t.Fatalf("expected Length=5, got %d", v.Length())
	}
}

func TestTypeVarchar_Length_Unicode(t *testing.T) {
	arena := newArena()
	// "привет" — 6 символов но 12 байт в UTF-8
	v := makeVarchar(t, arena, "привет", 100)

	if v.Length() != 6 {
		t.Fatalf("expected Length=6, got %d", v.Length())
	}
	if v.ByteLength() != 12 {
		t.Fatalf("expected ByteLength=12, got %d", v.ByteLength())
	}
}

func TestTypeVarchar_Length_Empty(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "", 100)

	if v.Length() != 0 {
		t.Fatalf("expected Length=0, got %d", v.Length())
	}
}

// ============================================================
// Contains / StartsWith / EndsWith
// ============================================================

func TestTypeVarchar_Contains(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello world", 100)

	if !v.Contains("world") {
		t.Fatal("expected Contains=true for 'world'")
	}
	if v.Contains("xyz") {
		t.Fatal("expected Contains=false for 'xyz'")
	}
}

func TestTypeVarchar_StartsWith(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello world", 100)

	if !v.StartsWith("hello") {
		t.Fatal("expected StartsWith=true for 'hello'")
	}
	if v.StartsWith("world") {
		t.Fatal("expected StartsWith=false for 'world'")
	}
}

func TestTypeVarchar_EndsWith(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello world", 100)

	if !v.EndsWith("world") {
		t.Fatal("expected EndsWith=true for 'world'")
	}
	if v.EndsWith("hello") {
		t.Fatal("expected EndsWith=false for 'hello'")
	}
}

// ============================================================
// Position
// ============================================================

func TestTypeVarchar_Position(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello world", 100)

	if v.Position("world") != 6 {
		t.Fatalf("expected Position=6, got %d", v.Position("world"))
	}
	if v.Position("hello") != 0 {
		t.Fatalf("expected Position=0, got %d", v.Position("hello"))
	}
	if v.Position("xyz") != -1 {
		t.Fatalf("expected Position=-1, got %d", v.Position("xyz"))
	}
}

func TestTypeVarchar_Position_Unicode(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "привет мир", 100)

	pos := v.Position("мир")
	if pos != 7 {
		t.Fatalf("expected Position=7, got %d", pos)
	}
}

// ============================================================
// Concat
// ============================================================

func TestTypeVarchar_Concat(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeVarchar(t, arena, "hello", 100)
	b := makeVarchar(t, arena, " world", 100)

	result, err := a.Concat(arena, b)
	if err != nil {
		t.Fatalf("Concat: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_Concat_ExceedsMaxLength(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeVarchar(t, arena, "hello", 8)
	b := makeVarchar(t, arena, " world", 8)

	// "hello world" = 11 символов > maxLen=8
	_, err := a.Concat(arena, b)
	if err == nil {
		t.Fatal("expected error: concat exceeds max length")
	}
}

func TestTypeVarchar_Concat_Empty(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeVarchar(t, arena, "hello", 100)
	b := makeVarchar(t, arena, "", 100)

	result, err := a.Concat(arena, b)
	if err != nil {
		t.Fatalf("Concat empty: %v", err)
	}
	if result.IntoGo() != "hello" {
		t.Fatalf("expected 'hello', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_Concat_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	a := makeVarchar(t, bigArena, "hello", 100)
	b := makeVarchar(t, bigArena, " world", 100)

	oomArena := pmem.NewArena(0)
	_, err := a.Concat(oomArena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// ToUpper / ToLower
// ============================================================

func TestTypeVarchar_ToUpper(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "Hello World", 100)

	result, err := v.ToUpper(arena)
	if err != nil {
		t.Fatalf("ToUpper: %v", err)
	}
	if result.IntoGo() != "HELLO WORLD" {
		t.Fatalf("expected 'HELLO WORLD', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_ToLower(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "Hello World", 100)

	result, err := v.ToLower(arena)
	if err != nil {
		t.Fatalf("ToLower: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_ToUpper_PreservesMeta(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello", 50)

	result, err := v.ToUpper(arena)
	if err != nil {
		t.Fatalf("ToUpper: %v", err)
	}
	if result.MaxLength() != 50 {
		t.Fatalf("expected MaxLength=50 preserved, got %d", result.MaxLength())
	}
}

// ============================================================
// Substring
// ============================================================

func TestTypeVarchar_Substring(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello world", 100)

	result, err := v.Substring(arena, 6, 5)
	if err != nil {
		t.Fatalf("Substring: %v", err)
	}
	if result.IntoGo() != "world" {
		t.Fatalf("expected 'world', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_Substring_Unicode(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "привет мир", 100)

	result, err := v.Substring(arena, 7, 3)
	if err != nil {
		t.Fatalf("Substring unicode: %v", err)
	}
	if result.IntoGo() != "мир" {
		t.Fatalf("expected 'мир', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_Substring_OutOfBounds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello", 100)

	_, err := v.Substring(arena, 3, 5)
	if err == nil {
		t.Fatal("expected out of bounds error")
	}
}

func TestTypeVarchar_Substring_NegativeStart(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello", 100)

	_, err := v.Substring(arena, -1, 3)
	if err == nil {
		t.Fatal("expected negative start error")
	}
}

// ============================================================
// Trim
// ============================================================

func TestTypeVarchar_Trim(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "  hello world  ", 100)

	result, err := v.Trim(arena)
	if err != nil {
		t.Fatalf("Trim: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_Trim_OnlySpaces(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "   ", 100)

	result, err := v.Trim(arena)
	if err != nil {
		t.Fatalf("Trim only spaces: %v", err)
	}
	if result.IntoGo() != "" {
		t.Fatalf("expected empty string, got %q", result.IntoGo())
	}
}

// ============================================================
// Replace
// ============================================================

func TestTypeVarchar_Replace(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello world hello", 100)

	result, err := v.Replace(arena, "hello", "bye")
	if err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if result.IntoGo() != "bye world bye" {
		t.Fatalf("expected 'bye world bye', got %q", result.IntoGo())
	}
}

func TestTypeVarchar_Replace_ExceedsMaxLength(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hi", 10)

	// "hi" → "hello world" = 11 символов > maxLen=10
	_, err := v.Replace(arena, "hi", "hello world")
	if err == nil {
		t.Fatal("expected error: replace result exceeds max length")
	}
}

func TestTypeVarchar_Replace_NotFound(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello world", 100)

	result, err := v.Replace(arena, "xyz", "abc")
	if err != nil {
		t.Fatalf("Replace not found: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected unchanged string, got %q", result.IntoGo())
	}
}

// ============================================================
// Like / ILike
// ============================================================

func TestTypeVarchar_Like(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello world", 100)

	cases := []struct {
		pattern string
		want    bool
	}{
		{"hello world", true},
		{"hello%", true},
		{"%world", true},
		{"hello_world", true},
		{"%", true},
		{"hello", false},
		{"%xyz", false},
	}

	for _, c := range cases {
		got := v.Like(c.pattern)
		if got != c.want {
			t.Fatalf("Like(%q): expected %v, got %v", c.pattern, c.want, got)
		}
	}
}

func TestTypeVarchar_ILike(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "Hello World", 100)

	if !v.ILike("hello%") {
		t.Fatal("expected ILike=true for 'hello%'")
	}
	if !v.ILike("%WORLD") {
		t.Fatal("expected ILike=true for '%WORLD'")
	}
	if v.ILike("xyz%") {
		t.Fatal("expected ILike=false for 'xyz%'")
	}
}

// ============================================================
// ToText
// ============================================================

func TestTypeVarchar_ToText(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 100)

	text := v.ToText()
	if text.IntoGo() != "hello" {
		t.Fatalf("expected 'hello', got %q", text.IntoGo())
	}
	// zero copy — тот же буфер
	if &text.GetBuffer()[0] != &v.GetBuffer()[0] {
		t.Fatal("expected ToText to share buffer (zero copy)")
	}
}

func TestTypeVarchar_ToText_GetType(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 100)

	text := v.ToText()
	if text.GetType() != ptypes.PTypeText {
		t.Fatalf("expected PTypeText after ToText, got %d", text.GetType())
	}
}

// ============================================================
// CastTo
// ============================================================

func TestTypeVarchar_CastTo_Bool_True(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "true", 100)

	result, err := v.CastTo(arena, ptypes.PTypeBool)
	if err != nil {
		t.Fatalf("CastTo bool: %v", err)
	}
	b, ok := result.IntoGo().(bool)
	if !ok || !b {
		t.Fatal("expected true")
	}
}

func TestTypeVarchar_CastTo_Bool_False(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "false", 100)

	result, err := v.CastTo(arena, ptypes.PTypeBool)
	if err != nil {
		t.Fatalf("CastTo bool: %v", err)
	}
	b, ok := result.IntoGo().(bool)
	if !ok || b {
		t.Fatal("expected false")
	}
}

func TestTypeVarchar_CastTo_Bool_Invalid(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "maybe", 100)

	_, err := v.CastTo(arena, ptypes.PTypeBool)
	if err == nil {
		t.Fatal("expected error for invalid bool value")
	}
}

func TestTypeVarchar_CastTo_Bytea(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello", 100)

	result, err := v.CastTo(arena, ptypes.PTypeBytea)
	if err != nil {
		t.Fatalf("CastTo bytea: %v", err)
	}
	b, ok := result.IntoGo().([]byte)
	if !ok {
		t.Fatal("expected []byte")
	}
	if string(b) != "hello" {
		t.Fatalf("expected 'hello', got %q", string(b))
	}
}

func TestTypeVarchar_CastTo_Text(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello", 100)

	result, err := v.CastTo(arena, ptypes.PTypeVarchar)
	if err != nil {
		t.Fatalf("CastTo varchar: %v", err)
	}
	s, ok := result.IntoGo().(string)
	if !ok {
		t.Fatal("expected string")
	}
	if s != "hello" {
		t.Fatalf("expected 'hello', got %q", s)
	}
}

func TestTypeVarchar_CastTo_Unsupported(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, arena, "hello", 100)

	_, err := v.CastTo(arena, ptypes.PTypeDate)
	if err == nil {
		t.Fatal("expected error for unsupported cast")
	}
}

func TestTypeVarchar_CastTo_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	v := makeVarchar(t, bigArena, "true", 100)

	oomArena := pmem.NewArena(0)
	_, err := v.CastTo(oomArena, ptypes.PTypeBool)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// String
// ============================================================

func TestTypeVarchar_String(t *testing.T) {
	arena := newArena()
	v := makeVarchar(t, arena, "hello", 50)

	s := v.String()
	if s != "varchar(50)(hello)" {
		t.Fatalf("expected 'varchar(50)(hello)', got %q", s)
	}
}
