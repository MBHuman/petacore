package ptypes

import (
	"bytes"
	"fmt"
	"petacore/sdk/pmem"
	"strings"
	"unicode/utf8"
)

type VarcharMeta struct {
	MaxLength int32
}

type TypeVarchar struct {
	BufferPtr []byte
	Meta      VarcharMeta
}

var _ BaseType[string] = (*TypeVarchar)(nil)
var _ OrderedType[string] = (*TypeVarchar)(nil)
var _ NullableType[string] = (*TypeVarchar)(nil)
var _ CastableType[string] = (*TypeVarchar)(nil)

func (t TypeVarchar) GetType() OID      { return PTypeVarchar }
func (t TypeVarchar) GetBuffer() []byte { return t.BufferPtr }

func (t TypeVarchar) IntoGo() string {
	return string(t.BufferPtr)
}

func (t TypeVarchar) Compare(other BaseType[string]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

// OrderedType

func (t TypeVarchar) LessThan(other BaseType[string]) bool       { return t.Compare(other) < 0 }
func (t TypeVarchar) GreaterThan(other BaseType[string]) bool    { return t.Compare(other) > 0 }
func (t TypeVarchar) LessOrEqual(other BaseType[string]) bool    { return t.Compare(other) <= 0 }
func (t TypeVarchar) GreaterOrEqual(other BaseType[string]) bool { return t.Compare(other) >= 0 }
func (t TypeVarchar) Between(low, high BaseType[string]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// NullableType

func (t TypeVarchar) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeVarchar) IsNotNull() bool { return t.BufferPtr != nil }

// ============================================================
// helpers
// ============================================================

// varcharFromString создаёт TypeVarchar с проверкой MaxLength
func varcharFromString(allocator pmem.Allocator, meta VarcharMeta, s string) (TypeVarchar, error) {
	runeLen := utf8.RuneCountInString(s)
	if meta.MaxLength > 0 && int32(runeLen) > meta.MaxLength {
		return TypeVarchar{}, fmt.Errorf("varchar: value length %d exceeds max %d", runeLen, meta.MaxLength)
	}
	b := []byte(s)
	buf, err := allocator.Alloc(len(b))
	if err != nil {
		return TypeVarchar{}, fmt.Errorf("varchar: alloc failed: %w", err)
	}
	copy(buf, b)
	return TypeVarchar{BufferPtr: buf, Meta: meta}, nil
}

// NewTypeVarchar создаёт TypeVarchar из строки
func NewTypeVarchar(allocator pmem.Allocator, val string, maxLen int32) (TypeVarchar, error) {
	return varcharFromString(allocator, VarcharMeta{MaxLength: maxLen}, val)
}

// ============================================================
// операции — создают новый буфер
// ============================================================

// Concat объединяет две строки с проверкой MaxLength результата
func (t TypeVarchar) Concat(allocator pmem.Allocator, other TypeVarchar) (TypeVarchar, error) {
	s := string(t.BufferPtr) + string(other.BufferPtr)
	return varcharFromString(allocator, t.Meta, s)
}

// ToUpper возвращает строку в верхнем регистре
func (t TypeVarchar) ToUpper(allocator pmem.Allocator) (TypeVarchar, error) {
	return varcharFromString(allocator, t.Meta, strings.ToUpper(string(t.BufferPtr)))
}

// ToLower возвращает строку в нижнем регистре
func (t TypeVarchar) ToLower(allocator pmem.Allocator) (TypeVarchar, error) {
	return varcharFromString(allocator, t.Meta, strings.ToLower(string(t.BufferPtr)))
}

// Substring возвращает подстроку по символьным позициям
func (t TypeVarchar) Substring(allocator pmem.Allocator, start, length int) (TypeVarchar, error) {
	runes := []rune(string(t.BufferPtr))
	if start < 0 || start+length > len(runes) {
		return TypeVarchar{}, fmt.Errorf("varchar substring: [%d:%d] out of bounds %d", start, start+length, len(runes))
	}
	return varcharFromString(allocator, t.Meta, string(runes[start:start+length]))
}

// Trim убирает пробелы с обоих концов
func (t TypeVarchar) Trim(allocator pmem.Allocator) (TypeVarchar, error) {
	return varcharFromString(allocator, t.Meta, strings.TrimSpace(string(t.BufferPtr)))
}

// Replace заменяет все вхождения old на new
func (t TypeVarchar) Replace(allocator pmem.Allocator, old, new string) (TypeVarchar, error) {
	s := strings.ReplaceAll(string(t.BufferPtr), old, new)
	return varcharFromString(allocator, t.Meta, s)
}

// ============================================================
// операции — только читают
// ============================================================

// Length возвращает длину в символах
func (t TypeVarchar) Length() int {
	return utf8.RuneCount(t.BufferPtr)
}

// ByteLength возвращает длину в байтах
func (t TypeVarchar) ByteLength() int {
	return len(t.BufferPtr)
}

// MaxLength возвращает ограничение длины из Meta
func (t TypeVarchar) MaxLength() int32 {
	return t.Meta.MaxLength
}

// Contains проверяет наличие подстроки
func (t TypeVarchar) Contains(sub string) bool {
	return strings.Contains(string(t.BufferPtr), sub)
}

// StartsWith проверяет начало строки
func (t TypeVarchar) StartsWith(prefix string) bool {
	return strings.HasPrefix(string(t.BufferPtr), prefix)
}

// EndsWith проверяет конец строки
func (t TypeVarchar) EndsWith(suffix string) bool {
	return strings.HasSuffix(string(t.BufferPtr), suffix)
}

// Like выполняет SQL LIKE: % — любая строка, _ — один символ
func (t TypeVarchar) Like(pattern string) bool {
	return sqlLikeMatch(string(t.BufferPtr), pattern, false)
}

// ILike выполняет SQL ILIKE (case-insensitive)
func (t TypeVarchar) ILike(pattern string) bool {
	return sqlLikeMatch(
		strings.ToLower(string(t.BufferPtr)),
		strings.ToLower(pattern),
		false,
	)
}

// Position возвращает позицию подстроки в символах (-1 если не найдено)
func (t TypeVarchar) Position(sub string) int {
	s := string(t.BufferPtr)
	idx := strings.Index(s, sub)
	if idx < 0 {
		return -1
	}
	return utf8.RuneCountInString(s[:idx])
}

// ToText конвертирует в TypeText без аллокации — zero copy
func (t TypeVarchar) ToText() TypeText {
	return TypeText{BufferPtr: t.BufferPtr}
}

// ============================================================
// CastableType
// ============================================================

func (t TypeVarchar) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	// делегируем в TypeText — логика кастов одинакова
	result, err := TypeText{BufferPtr: t.BufferPtr}.CastTo(allocator, targetType)
	if err != nil {
		return nil, fmt.Errorf("varchar cast: %w", err)
	}
	return result, nil
}

func (t TypeVarchar) String() string {
	return fmt.Sprintf("varchar(%d)(%s)", t.Meta.MaxLength, string(t.BufferPtr))
}
