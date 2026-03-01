package ptypes

import (
	"bytes"
	"fmt"
	"petacore/sdk/pmem"
	"strings"
	"unicode/utf8"
)

type TypeText struct {
	BufferPtr []byte
}

var _ BaseType[string] = (*TypeText)(nil)
var _ OrderedType[string] = (*TypeText)(nil)
var _ NullableType[string] = (*TypeText)(nil)
var _ CastableType[string] = (*TypeText)(nil)

func NewTypeText(val string) TypeText {
	return TypeText{BufferPtr: []byte(val)}
}

func (t TypeText) GetType() OID      { return PTypeText }
func (t TypeText) GetBuffer() []byte { return t.BufferPtr }

func (t TypeText) IntoGo() string {
	return string(t.BufferPtr)
}

func (t TypeText) Compare(other BaseType[string]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

// OrderedType

func (t TypeText) LessThan(other BaseType[string]) bool       { return t.Compare(other) < 0 }
func (t TypeText) GreaterThan(other BaseType[string]) bool    { return t.Compare(other) > 0 }
func (t TypeText) LessOrEqual(other BaseType[string]) bool    { return t.Compare(other) <= 0 }
func (t TypeText) GreaterOrEqual(other BaseType[string]) bool { return t.Compare(other) >= 0 }
func (t TypeText) Between(low, high BaseType[string]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// NullableType

func (t TypeText) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeText) IsNotNull() bool { return t.BufferPtr != nil }

// ============================================================
// TextType операции — создают новый буфер
// ============================================================

// Concat объединяет две строки
func (t TypeText) Concat(allocator pmem.Allocator, other TypeText) (TypeText, error) {
	size := len(t.BufferPtr) + len(other.BufferPtr)
	buf, err := allocator.Alloc(size)
	if err != nil {
		return TypeText{}, fmt.Errorf("text concat: %w", err)
	}
	copy(buf, t.BufferPtr)
	copy(buf[len(t.BufferPtr):], other.BufferPtr)
	return TypeText{BufferPtr: buf}, nil
}

// ToUpper возвращает строку в верхнем регистре
func (t TypeText) ToUpper(allocator pmem.Allocator) (TypeText, error) {
	s := strings.ToUpper(string(t.BufferPtr))
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text toupper: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

// ToLower возвращает строку в нижнем регистре
func (t TypeText) ToLower(allocator pmem.Allocator) (TypeText, error) {
	s := strings.ToLower(string(t.BufferPtr))
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text tolower: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

// Substring возвращает подстроку по символьным позициям (не байтовым)
func (t TypeText) Substring(allocator pmem.Allocator, start, length int) (TypeText, error) {
	runes := []rune(string(t.BufferPtr))
	if start < 0 || start+length > len(runes) {
		return TypeText{}, fmt.Errorf("text substring: [%d:%d] out of bounds %d", start, start+length, len(runes))
	}
	s := string(runes[start : start+length])
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text substring: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

// Trim убирает пробелы с обоих концов
func (t TypeText) Trim(allocator pmem.Allocator) (TypeText, error) {
	s := strings.TrimSpace(string(t.BufferPtr))
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text trim: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

// Replace заменяет все вхождения old на new
func (t TypeText) Replace(allocator pmem.Allocator, old, new string) (TypeText, error) {
	s := strings.ReplaceAll(string(t.BufferPtr), old, new)
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text replace: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

// ============================================================
// TextType операции — только читают, аллокатор не нужен
// ============================================================

// Length возвращает длину в символах (Unicode code points)
func (t TypeText) Length() int {
	return utf8.RuneCount(t.BufferPtr)
}

// ByteLength возвращает длину в байтах
func (t TypeText) ByteLength() int {
	return len(t.BufferPtr)
}

// Contains проверяет наличие подстроки
func (t TypeText) Contains(sub string) bool {
	return strings.Contains(string(t.BufferPtr), sub)
}

// StartsWith проверяет начало строки
func (t TypeText) StartsWith(prefix string) bool {
	return strings.HasPrefix(string(t.BufferPtr), prefix)
}

// EndsWith проверяет конец строки
func (t TypeText) EndsWith(suffix string) bool {
	return strings.HasSuffix(string(t.BufferPtr), suffix)
}

// Like выполняет SQL LIKE сравнение: % — любая строка, _ — один символ
func (t TypeText) Like(pattern string) bool {
	return sqlLikeMatch(string(t.BufferPtr), pattern, false)
}

// ILike выполняет SQL ILIKE сравнение (case-insensitive)
func (t TypeText) ILike(pattern string) bool {
	return sqlLikeMatch(
		strings.ToLower(string(t.BufferPtr)),
		strings.ToLower(pattern),
		false,
	)
}

// Position возвращает позицию подстроки в символах (-1 если не найдено)
func (t TypeText) Position(sub string) int {
	s := string(t.BufferPtr)
	idx := strings.Index(s, sub)
	if idx < 0 {
		return -1
	}
	return utf8.RuneCountInString(s[:idx])
}

// ============================================================
// CastableType
// ============================================================

func (t TypeText) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	s := string(t.BufferPtr)

	switch targetType {
	case PTypeBool:
		var v bool
		switch strings.ToLower(strings.TrimSpace(s)) {
		case "true", "t", "yes", "y", "1", "on":
			v = true
		case "false", "f", "no", "n", "0", "off":
			v = false
		default:
			return nil, fmt.Errorf("text cast to bool: invalid value %q", s)
		}
		buf, err := allocator.Alloc(1)
		if err != nil {
			return nil, fmt.Errorf("text cast to bool: %w", err)
		}
		if v {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
		return anyWrapper[bool]{TypeBool{BufferPtr: buf}}, nil

	case PTypeBytea:
		buf, err := allocator.Alloc(len(t.BufferPtr))
		if err != nil {
			return nil, fmt.Errorf("text cast to bytea: %w", err)
		}
		copy(buf, t.BufferPtr)
		return anyWrapper[[]byte]{TypeBytea{BufferPtr: buf}}, nil

	case PTypeVarchar:
		buf, err := allocator.Alloc(len(t.BufferPtr))
		if err != nil {
			return nil, fmt.Errorf("text cast to varchar: %w", err)
		}
		copy(buf, t.BufferPtr)
		return anyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("text: unsupported cast to OID %d", targetType)
	}
}

func (t TypeText) String() string {
	return "text(" + string(t.BufferPtr) + ")"
}

// ============================================================
// SQL LIKE pattern matching
// ============================================================

// sqlLikeMatch реализует SQL LIKE: % — любая строка, _ — один символ
func sqlLikeMatch(s, pattern string, escape bool) bool {
	sr := []rune(s)
	pr := []rune(pattern)
	return likeMatchRunes(sr, pr)
}

func likeMatchRunes(s, p []rune) bool {
	for len(p) > 0 {
		switch p[0] {
		case '%':
			// пропускаем подряд идущие %
			for len(p) > 0 && p[0] == '%' {
				p = p[1:]
			}
			if len(p) == 0 {
				return true
			}
			// пробуем сопоставить остаток паттерна с каждой позиции строки
			for i := 0; i <= len(s); i++ {
				if likeMatchRunes(s[i:], p) {
					return true
				}
			}
			return false

		case '_':
			// один любой символ
			if len(s) == 0 {
				return false
			}
			s = s[1:]
			p = p[1:]

		default:
			if len(s) == 0 || s[0] != p[0] {
				return false
			}
			s = s[1:]
			p = p[1:]
		}
	}
	return len(s) == 0
}
