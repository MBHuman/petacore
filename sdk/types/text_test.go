package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

func makeText(t *testing.T, arena pmem.Allocator, val string) ptypes.TypeText {
	t.Helper()
	buf, err := serializers.TextSerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize text: %v", err)
	}
	result, err := serializers.TextSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize text: %v", err)
	}
	return result
}

// ============================================================
// Serialize / Deserialize
// ============================================================

func TestTypeText_SerializeDeserialize(t *testing.T) {
	arena := newArena()

	cases := []string{"hello", "world", "hello world", ""}
	for _, v := range cases {
		result := makeText(t, arena, v)
		if result.IntoGo() != v {
			t.Fatalf("expected %q, got %q", v, result.IntoGo())
		}
	}
}

func TestTypeText_SerializeDeserialize_Unicode(t *testing.T) {
	arena := newArena()

	cases := []string{"привет мир", "日本語", "🎉🚀💡", "café"}
	for _, v := range cases {
		result := makeText(t, arena, v)
		if result.IntoGo() != v {
			t.Fatalf("expected %q, got %q", v, result.IntoGo())
		}
	}
}

func TestTypeText_NewTypeText(t *testing.T) {
	text := ptypes.NewTypeText("hello")
	if text.IntoGo() != "hello" {
		t.Fatalf("expected hello, got %q", text.IntoGo())
	}
}

func TestTypeText_GetType(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello")

	if text.GetType() != ptypes.PTypeText {
		t.Fatalf("expected PTypeText, got %d", text.GetType())
	}
}

func TestTypeText_GetBuffer(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello")

	buf := text.GetBuffer()
	if len(buf) != 5 {
		t.Fatalf("expected buffer len 5, got %d", len(buf))
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeText_IsNull(t *testing.T) {
	text := ptypes.TypeText{BufferPtr: nil}
	if !text.IsNull() {
		t.Fatal("expected IsNull=true")
	}
	if text.IsNotNull() {
		t.Fatal("expected IsNotNull=false")
	}
}

func TestTypeText_IsNotNull(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello")

	if text.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !text.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

func TestTypeText_IsNotNull_Empty(t *testing.T) {
	// пустая строка — не null
	text := ptypes.TypeText{BufferPtr: []byte{}}
	if text.IsNull() {
		t.Fatal("expected IsNull=false for empty string")
	}
	if !text.IsNotNull() {
		t.Fatal("expected IsNotNull=true for empty string")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeText_Compare(t *testing.T) {
	arena := newArena()
	a := makeText(t, arena, "apple")
	b := makeText(t, arena, "banana")
	c := makeText(t, arena, "apple")

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

func TestTypeText_Compare_EmptyString(t *testing.T) {
	arena := newArena()
	empty := makeText(t, arena, "")
	nonempty := makeText(t, arena, "a")

	if empty.Compare(nonempty) >= 0 {
		t.Fatal("expected empty < nonempty")
	}
}

func TestTypeText_LessThan(t *testing.T) {
	arena := newArena()
	a := makeText(t, arena, "apple")
	b := makeText(t, arena, "banana")

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeText_GreaterThan(t *testing.T) {
	arena := newArena()
	a := makeText(t, arena, "apple")
	b := makeText(t, arena, "banana")

	if !b.GreaterThan(a) {
		t.Fatal("expected b > a")
	}
	if a.GreaterThan(b) {
		t.Fatal("expected a not > b")
	}
}

func TestTypeText_LessOrEqual(t *testing.T) {
	arena := newArena()
	a := makeText(t, arena, "apple")
	b := makeText(t, arena, "banana")

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

func TestTypeText_GreaterOrEqual(t *testing.T) {
	arena := newArena()
	a := makeText(t, arena, "apple")
	b := makeText(t, arena, "banana")

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

func TestTypeText_Between(t *testing.T) {
	arena := newArena()
	a := makeText(t, arena, "apple")
	b := makeText(t, arena, "mango")
	c := makeText(t, arena, "zebra")

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
// Length / ByteLength
// ============================================================

func TestTypeText_Length_ASCII(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello")

	if text.Length() != 5 {
		t.Fatalf("expected Length=5, got %d", text.Length())
	}
}

func TestTypeText_Length_Unicode(t *testing.T) {
	arena := newArena()
	// "привет" — 6 символов но 12 байт в UTF-8
	text := makeText(t, arena, "привет")

	if text.Length() != 6 {
		t.Fatalf("expected Length=6, got %d", text.Length())
	}
	if text.ByteLength() != 12 {
		t.Fatalf("expected ByteLength=12, got %d", text.ByteLength())
	}
}

func TestTypeText_Length_Emoji(t *testing.T) {
	arena := newArena()
	// каждый эмодзи — 1 руна но 4 байта
	text := makeText(t, arena, "🎉🚀")

	if text.Length() != 2 {
		t.Fatalf("expected Length=2, got %d", text.Length())
	}
	if text.ByteLength() != 8 {
		t.Fatalf("expected ByteLength=8, got %d", text.ByteLength())
	}
}

func TestTypeText_Length_Empty(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "")

	if text.Length() != 0 {
		t.Fatalf("expected Length=0, got %d", text.Length())
	}
	if text.ByteLength() != 0 {
		t.Fatalf("expected ByteLength=0, got %d", text.ByteLength())
	}
}

// ============================================================
// Contains / StartsWith / EndsWith
// ============================================================

func TestTypeText_Contains(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello world")

	if !text.Contains("world") {
		t.Fatal("expected Contains=true for 'world'")
	}
	if !text.Contains("hello") {
		t.Fatal("expected Contains=true for 'hello'")
	}
	if !text.Contains("") {
		t.Fatal("expected Contains=true for empty string")
	}
	if text.Contains("xyz") {
		t.Fatal("expected Contains=false for 'xyz'")
	}
}

func TestTypeText_StartsWith(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello world")

	if !text.StartsWith("hello") {
		t.Fatal("expected StartsWith=true for 'hello'")
	}
	if !text.StartsWith("") {
		t.Fatal("expected StartsWith=true for empty prefix")
	}
	if text.StartsWith("world") {
		t.Fatal("expected StartsWith=false for 'world'")
	}
}

func TestTypeText_EndsWith(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello world")

	if !text.EndsWith("world") {
		t.Fatal("expected EndsWith=true for 'world'")
	}
	if !text.EndsWith("") {
		t.Fatal("expected EndsWith=true for empty suffix")
	}
	if text.EndsWith("hello") {
		t.Fatal("expected EndsWith=false for 'hello'")
	}
}

// ============================================================
// Position
// ============================================================

func TestTypeText_Position(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello world")

	if text.Position("world") != 6 {
		t.Fatalf("expected Position=6, got %d", text.Position("world"))
	}
	if text.Position("hello") != 0 {
		t.Fatalf("expected Position=0, got %d", text.Position("hello"))
	}
	if text.Position("xyz") != -1 {
		t.Fatalf("expected Position=-1, got %d", text.Position("xyz"))
	}
}

func TestTypeText_Position_Unicode(t *testing.T) {
	arena := newArena()
	// "привет мир" — "мир" начинается на позиции 7 (в символах)
	text := makeText(t, arena, "привет мир")

	pos := text.Position("мир")
	if pos != 7 {
		t.Fatalf("expected Position=7, got %d", pos)
	}
}

// ============================================================
// Concat
// ============================================================

func TestTypeText_Concat(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeText(t, arena, "hello")
	b := makeText(t, arena, " world")

	result, err := a.Concat(arena, b)
	if err != nil {
		t.Fatalf("Concat: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", result.IntoGo())
	}
}

func TestTypeText_Concat_Empty(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeText(t, arena, "hello")
	b := makeText(t, arena, "")

	result, err := a.Concat(arena, b)
	if err != nil {
		t.Fatalf("Concat empty: %v", err)
	}
	if result.IntoGo() != "hello" {
		t.Fatalf("expected 'hello', got %q", result.IntoGo())
	}
}

func TestTypeText_Concat_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	a := makeText(t, bigArena, "hello")
	b := makeText(t, bigArena, " world")

	oomArena := pmem.NewArena(0)
	_, err := a.Concat(oomArena, b)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// ToUpper / ToLower
// ============================================================

func TestTypeText_ToUpper(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "Hello World")

	result, err := text.ToUpper(arena)
	if err != nil {
		t.Fatalf("ToUpper: %v", err)
	}
	if result.IntoGo() != "HELLO WORLD" {
		t.Fatalf("expected 'HELLO WORLD', got %q", result.IntoGo())
	}
}

func TestTypeText_ToLower(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "Hello World")

	result, err := text.ToLower(arena)
	if err != nil {
		t.Fatalf("ToLower: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", result.IntoGo())
	}
}

func TestTypeText_ToUpper_Unicode(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "привет")

	result, err := text.ToUpper(arena)
	if err != nil {
		t.Fatalf("ToUpper unicode: %v", err)
	}
	if result.IntoGo() != "ПРИВЕТ" {
		t.Fatalf("expected 'ПРИВЕТ', got %q", result.IntoGo())
	}
}

func TestTypeText_ToUpper_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	text := makeText(t, bigArena, "hello")

	oomArena := pmem.NewArena(0)
	_, err := text.ToUpper(oomArena)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// Substring
// ============================================================

func TestTypeText_Substring(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello world")

	result, err := text.Substring(arena, 6, 5)
	if err != nil {
		t.Fatalf("Substring: %v", err)
	}
	if result.IntoGo() != "world" {
		t.Fatalf("expected 'world', got %q", result.IntoGo())
	}
}

func TestTypeText_Substring_Unicode(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	// "привет мир" — берём "мир" (позиции 7..9)
	text := makeText(t, arena, "привет мир")

	result, err := text.Substring(arena, 7, 3)
	if err != nil {
		t.Fatalf("Substring unicode: %v", err)
	}
	if result.IntoGo() != "мир" {
		t.Fatalf("expected 'мир', got %q", result.IntoGo())
	}
}

func TestTypeText_Substring_OutOfBounds(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello")

	_, err := text.Substring(arena, 3, 5)
	if err == nil {
		t.Fatal("expected out of bounds error")
	}
}

func TestTypeText_Substring_NegativeStart(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello")

	_, err := text.Substring(arena, -1, 3)
	if err == nil {
		t.Fatal("expected negative start error")
	}
}

func TestTypeText_Substring_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	text := makeText(t, bigArena, "hello world")

	oomArena := pmem.NewArena(0)
	_, err := text.Substring(oomArena, 0, 5)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// Trim
// ============================================================

func TestTypeText_Trim(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "  hello world  ")

	result, err := text.Trim(arena)
	if err != nil {
		t.Fatalf("Trim: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", result.IntoGo())
	}
}

func TestTypeText_Trim_AlreadyTrimmed(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello")

	result, err := text.Trim(arena)
	if err != nil {
		t.Fatalf("Trim already trimmed: %v", err)
	}
	if result.IntoGo() != "hello" {
		t.Fatalf("expected 'hello', got %q", result.IntoGo())
	}
}

func TestTypeText_Trim_OnlySpaces(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "   ")

	result, err := text.Trim(arena)
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

func TestTypeText_Replace(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello world hello")

	result, err := text.Replace(arena, "hello", "bye")
	if err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if result.IntoGo() != "bye world bye" {
		t.Fatalf("expected 'bye world bye', got %q", result.IntoGo())
	}
}

func TestTypeText_Replace_NotFound(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello world")

	result, err := text.Replace(arena, "xyz", "abc")
	if err != nil {
		t.Fatalf("Replace not found: %v", err)
	}
	if result.IntoGo() != "hello world" {
		t.Fatalf("expected unchanged string, got %q", result.IntoGo())
	}
}

func TestTypeText_Replace_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	text := makeText(t, bigArena, "hello world")

	oomArena := pmem.NewArena(0)
	_, err := text.Replace(oomArena, "hello", "bye")
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// Like / ILike
// ============================================================

func TestTypeText_Like(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "hello world")

	cases := []struct {
		pattern string
		want    bool
	}{
		{"hello world", true},
		{"hello%", true},
		{"%world", true},
		{"%o w%", true},
		{"hello_world", true},
		{"hello%xyz", false},
		{"%xyz", false},
		{"hello", false},
		{"%", true},
		{"", false},
	}

	for _, c := range cases {
		got := text.Like(c.pattern)
		if got != c.want {
			t.Fatalf("Like(%q): expected %v, got %v", c.pattern, c.want, got)
		}
	}
}

func TestTypeText_Like_Unicode(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "привет мир")

	if !text.Like("привет%") {
		t.Fatal("expected Like=true for 'привет%'")
	}
	if !text.Like("%мир") {
		t.Fatal("expected Like=true for '%мир'")
	}
	if !text.Like("привет_мир") {
		t.Fatal("expected Like=true for 'привет_мир'")
	}
}

func TestTypeText_ILike(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "Hello World")

	if !text.ILike("hello%") {
		t.Fatal("expected ILike=true for 'hello%'")
	}
	if !text.ILike("%WORLD") {
		t.Fatal("expected ILike=true for '%WORLD'")
	}
	if !text.ILike("HELLO_WORLD") {
		t.Fatal("expected ILike=true for 'HELLO_WORLD'")
	}
	if text.ILike("xyz%") {
		t.Fatal("expected ILike=false for 'xyz%'")
	}
}

func TestTypeText_Like_MultipleWildcards(t *testing.T) {
	arena := newArena()
	text := makeText(t, arena, "abcdef")

	pattern1 := "%b%d%"
	if !text.Like(pattern1) {
		t.Fatalf("expected Like=true for %q", pattern1)
	}

	pattern2 := "a%%f"
	if !text.Like(pattern2) {
		t.Fatalf("expected Like=true for %q", pattern2)
	}

	pattern3 := "%b%z%"
	if text.Like(pattern3) {
		t.Fatalf("expected Like=false for %q", pattern3)
	}
}

// ============================================================
// CastTo
// ============================================================

func TestTypeText_CastTo_Bool_True(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)

	trueCases := []string{"true", "t", "yes", "y", "1", "on", "TRUE", "YES"}
	for _, v := range trueCases {
		text := makeText(t, arena, v)
		result, err := text.CastTo(arena, ptypes.PTypeBool)
		if err != nil {
			t.Fatalf("CastTo bool(%q): %v", v, err)
		}
		b, ok := result.IntoGo().(bool)
		if !ok || !b {
			t.Fatalf("CastTo bool(%q): expected true", v)
		}
	}
}

func TestTypeText_CastTo_Bool_False(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)

	falseCases := []string{"false", "f", "no", "n", "0", "off", "FALSE", "NO"}
	for _, v := range falseCases {
		text := makeText(t, arena, v)
		result, err := text.CastTo(arena, ptypes.PTypeBool)
		if err != nil {
			t.Fatalf("CastTo bool(%q): %v", v, err)
		}
		b, ok := result.IntoGo().(bool)
		if !ok || b {
			t.Fatalf("CastTo bool(%q): expected false", v)
		}
	}
}

func TestTypeText_CastTo_Bool_Invalid(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "maybe")

	_, err := text.CastTo(arena, ptypes.PTypeBool)
	if err == nil {
		t.Fatal("expected error for invalid bool value")
	}
}

func TestTypeText_CastTo_Bytea(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello")

	result, err := text.CastTo(arena, ptypes.PTypeBytea)
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

func TestTypeText_CastTo_Varchar(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello")

	result, err := text.CastTo(arena, ptypes.PTypeVarchar)
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

func TestTypeText_CastTo_Unsupported(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	text := makeText(t, arena, "hello")

	_, err := text.CastTo(arena, ptypes.PTypeDate)
	if err == nil {
		t.Fatal("expected error for unsupported cast")
	}
}

func TestTypeText_CastTo_OOM(t *testing.T) {
	bigArena := pmem.NewArena(64 * 1024)
	text := makeText(t, bigArena, "true")

	oomArena := pmem.NewArena(0)
	_, err := text.CastTo(oomArena, ptypes.PTypeBool)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}
