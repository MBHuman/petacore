package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"petacore/sdk/pmem"
	"regexp"
	"strconv"
	"strings"
	"time"
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

var _ TextType[string] = (*TypeText)(nil)

func (t TypeText) Concat(allocator pmem.Allocator, other TextType[string]) (TextType[string], error) {
	size := len(t.BufferPtr) + len(other.GetBuffer())
	buf, err := allocator.Alloc(size)
	if err != nil {
		return TypeText{}, fmt.Errorf("text concat: %w", err)
	}
	copy(buf, t.BufferPtr)
	copy(buf[len(t.BufferPtr):], other.GetBuffer())
	return TypeText{BufferPtr: buf}, nil
}

func (t TypeText) ToUpper(allocator pmem.Allocator) (TextType[string], error) {
	s := strings.ToUpper(string(t.BufferPtr))
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text toupper: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

func (t TypeText) ToLower(allocator pmem.Allocator) (TextType[string], error) {
	s := strings.ToLower(string(t.BufferPtr))
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text tolower: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

func (t TypeText) Substring(allocator pmem.Allocator, start, length int) (TextType[string], error) {
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

func (t TypeText) Trim(allocator pmem.Allocator) (TextType[string], error) {
	s := strings.TrimSpace(string(t.BufferPtr))
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text trim: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

func (t TypeText) Replace(allocator pmem.Allocator, old, new string) (TextType[string], error) {
	s := strings.ReplaceAll(string(t.BufferPtr), old, new)
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text replace: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

func (t TypeText) RegexpMatch(pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("text regexp: invalid pattern %q: %w", pattern, err)
	}
	return re.Match(t.BufferPtr), nil
}

func (t TypeText) RegexpMatchCompiled(re *regexp.Regexp) bool {
	return re.Match(t.BufferPtr)
}

func (t TypeText) RegexpReplace(allocator pmem.Allocator, pattern, replacement string) (TextType[string], error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return TypeText{}, fmt.Errorf("text regexp replace: invalid pattern %q: %w", pattern, err)
	}
	return t.RegexpReplaceCompiled(allocator, re, replacement)
}

func (t TypeText) RegexpReplaceCompiled(allocator pmem.Allocator, re *regexp.Regexp, replacement string) (TextType[string], error) {
	s := re.ReplaceAllString(string(t.BufferPtr), replacement)
	buf, err := allocator.Alloc(len(s))
	if err != nil {
		return TypeText{}, fmt.Errorf("text regexp replace: %w", err)
	}
	copy(buf, s)
	return TypeText{BufferPtr: buf}, nil
}

func (t TypeText) AsStr() string {
	return string(t.BufferPtr)
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

var _ CastableType[string] = (*TypeText)(nil)

func (t TypeText) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	s := strings.TrimSpace(string(t.BufferPtr))

	switch targetType {
	case PTypeBool:
		var v bool
		switch strings.ToLower(s) {
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
		return AnyWrapper[bool]{TypeBool{BufferPtr: buf}}, nil

	case PTypeBytea:
		buf, err := allocator.Alloc(len(t.BufferPtr))
		if err != nil {
			return nil, fmt.Errorf("text cast to bytea: %w", err)
		}
		copy(buf, t.BufferPtr)
		return AnyWrapper[[]byte]{TypeBytea{BufferPtr: buf}}, nil

	case PTypeVarchar:
		buf, err := allocator.Alloc(len(t.BufferPtr))
		if err != nil {
			return nil, fmt.Errorf("text cast to varchar: %w", err)
		}
		copy(buf, t.BufferPtr)
		return AnyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	case PTypeInt2:
		v, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("text cast to int2: invalid value %q", s)
		}
		buf, err := allocator.AllocAligned(2, 2)
		if err != nil {
			return nil, fmt.Errorf("text cast to int2: %w", err)
		}
		binary.BigEndian.PutUint16(buf, uint16(int16(v))^0x8000)
		return AnyWrapper[int16]{TypeInt2{BufferPtr: buf}}, nil

	case PTypeInt4:
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("text cast to int4: invalid value %q", s)
		}
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("text cast to int4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, uint32(int32(v))^0x80000000)
		return AnyWrapper[int32]{TypeInt4{BufferPtr: buf}}, nil

	case PTypeInt8:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("text cast to int8: invalid value %q", s)
		}
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("text cast to int8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(v)^0x8000000000000000)
		return AnyWrapper[int64]{TypeInt8{BufferPtr: buf}}, nil

	case PTypeFloat4:
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, fmt.Errorf("text cast to float4: invalid value %q", s)
		}
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("text cast to float4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, OrderableFloat32bits(float32(v)))
		return AnyWrapper[float32]{TypeFloat4{BufferPtr: buf}}, nil

	case PTypeFloat8:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("text cast to float8: invalid value %q", s)
		}
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("text cast to float8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, OrderableFloat64bits(v))
		return AnyWrapper[float64]{TypeFloat8{BufferPtr: buf}}, nil

	case PTypeNumeric:
		meta := NumericMeta{Precision: 38, Scale: 10}
		f, _, err := new(big.Float).SetPrec(256).Parse(s, 10)
		if err != nil {
			return nil, fmt.Errorf("text cast to numeric: invalid value %q", s)
		}
		result, err := numericFromBigFloat(allocator, meta, f)
		if err != nil {
			return nil, fmt.Errorf("text cast to numeric: %w", err)
		}
		return AnyWrapper[[]byte]{result}, nil

	case PTypeDate:
		tm, err := time.Parse("2006-01-02", s)
		if err != nil {
			return nil, fmt.Errorf("text cast to date: invalid value %q, expected YYYY-MM-DD", s)
		}
		result, err := NewTypeDate(allocator, tm)
		if err != nil {
			return nil, fmt.Errorf("text cast to date: %w", err)
		}
		return AnyWrapper[*time.Time]{result}, nil

	case PTypeTimestamp:
		tm, err := parseTimestampString(s)
		if err != nil {
			return nil, fmt.Errorf("text cast to timestamp: %w", err)
		}
		usec := tm.Sub(PgEpoch).Microseconds()
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("text cast to timestamp: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
		return AnyWrapper[*time.Time]{TypeTimestamp{BufferPtr: buf}}, nil

	case PTypeTimestampz:
		tm, err := parseTimestampString(s)
		if err != nil {
			return nil, fmt.Errorf("text cast to timestampz: %w", err)
		}
		usec := tm.UTC().Sub(PgEpoch).Microseconds()
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("text cast to timestampz: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(usec)^0x8000000000000000)
		return AnyWrapper[*time.Time]{TypeTimestampz{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("text: unsupported cast to OID %d", targetType)
	}
}

// parseTimestampString парсит строку в несколько форматов ISO 8601
var timestampFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.999999",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.999999",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.999999Z07:00",
}

func parseTimestampString(s string) (time.Time, error) {
	for _, layout := range timestampFormats {
		if tm, err := time.Parse(layout, s); err == nil {
			return tm, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid timestamp %q, expected ISO 8601", s)
}

func TextFactory(buf []byte) TypeText {
	return TypeText{BufferPtr: buf}
}

func TextComparator(a, b TypeText) int {
	return bytes.Compare(a.BufferPtr, b.BufferPtr)
}
