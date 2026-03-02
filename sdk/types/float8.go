// types/float8.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"petacore/sdk/pmem"
	"strconv"
)

type TypeFloat8 struct {
	BufferPtr []byte
}

var _ BaseType[float64] = (*TypeFloat8)(nil)
var _ NumericType[float64] = (*TypeFloat8)(nil)
var _ OrderedType[float64] = (*TypeFloat8)(nil)
var _ NullableType[float64] = (*TypeFloat8)(nil)

func (t TypeFloat8) GetType() OID { return PTypeFloat8 }

func (t TypeFloat8) Compare(other BaseType[float64]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeFloat8) GetBuffer() []byte { return t.BufferPtr }

func (t TypeFloat8) IntoGo() float64 {
	if len(t.BufferPtr) < 8 {
		return 0
	}
	return Float64fromOrderableBits(binary.BigEndian.Uint64(t.BufferPtr))
}

// NullableType

func (t TypeFloat8) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeFloat8) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeFloat8) LessThan(other BaseType[float64]) bool       { return t.Compare(other) < 0 }
func (t TypeFloat8) GreaterThan(other BaseType[float64]) bool    { return t.Compare(other) > 0 }
func (t TypeFloat8) LessOrEqual(other BaseType[float64]) bool    { return t.Compare(other) <= 0 }
func (t TypeFloat8) GreaterOrEqual(other BaseType[float64]) bool { return t.Compare(other) >= 0 }
func (t TypeFloat8) Between(low, high BaseType[float64]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func float8FromVal(allocator pmem.Allocator, v float64) (TypeFloat8, error) {
	buf, err := allocator.AllocAligned(8, 8)
	if err != nil {
		return TypeFloat8{}, fmt.Errorf("float8: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint64(buf, OrderableFloat64bits(v))
	return TypeFloat8{BufferPtr: buf}, nil
}

// NumericType

func (t TypeFloat8) Add(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	return float8FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeFloat8) Sub(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	return float8FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeFloat8) Mul(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	return float8FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeFloat8) Div(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float8: division by zero")
	}
	return float8FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeFloat8) Mod(allocator pmem.Allocator, other NumericType[float64]) (NumericType[float64], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float8: modulo by zero")
	}
	return float8FromVal(allocator, math.Mod(t.IntoGo(), other.IntoGo()))
}

func (t TypeFloat8) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeFloat8) Neg(allocator pmem.Allocator) NumericType[float64] {
	v, _ := float8FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeFloat8) Abs(allocator pmem.Allocator) NumericType[float64] {
	v, _ := float8FromVal(allocator, math.Abs(t.IntoGo()))
	return v
}

func (t TypeFloat8) String() string {
	return "float8(" + fmt.Sprintf("%v", t.IntoGo()) + ")"
}

var _ CastableType[float64] = (*TypeFloat8)(nil)

// CastableType

func (t TypeFloat8) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	v := t.IntoGo()

	switch targetType {
	case PTypeFloat4:
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("float8 cast to float4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, OrderableFloat32bits(float32(v)))
		return AnyWrapper[float32]{TypeFloat4{BufferPtr: buf}}, nil

	case PTypeInt2:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return nil, fmt.Errorf("float8 cast to int2: value %v out of range", v)
		}
		buf, err := allocator.AllocAligned(2, 2)
		if err != nil {
			return nil, fmt.Errorf("float8 cast to int2: %w", err)
		}
		binary.BigEndian.PutUint16(buf, uint16(int16(v))^0x8000)
		return AnyWrapper[int16]{TypeInt2{BufferPtr: buf}}, nil

	case PTypeInt4:
		if v < math.MinInt32 || v > math.MaxInt32 {
			return nil, fmt.Errorf("float8 cast to int4: value %v out of range", v)
		}
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("float8 cast to int4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, uint32(int32(v))^0x80000000)
		return AnyWrapper[int32]{TypeInt4{BufferPtr: buf}}, nil

	case PTypeInt8:
		if v < math.MinInt64 || v > math.MaxInt64 {
			return nil, fmt.Errorf("float8 cast to int8: value %v out of range", v)
		}
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("float8 cast to int8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(int64(v))^0x8000000000000000)
		return AnyWrapper[int64]{TypeInt8{BufferPtr: buf}}, nil

	case PTypeNumeric:
		meta := NumericMeta{Precision: 38, Scale: 10}
		s := strconv.FormatFloat(v, 'f', 10, 64)
		f, _, err := new(big.Float).SetPrec(256).Parse(s, 10)
		if err != nil {
			return nil, fmt.Errorf("float8 cast to numeric: %w", err)
		}
		result, err := numericFromBigFloat(allocator, meta, f)
		if err != nil {
			return nil, fmt.Errorf("float8 cast to numeric: %w", err)
		}
		return AnyWrapper[[]byte]{result}, nil

	case PTypeText:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("float8 cast to text: %w", err)
		}
		copy(buf, s)
		return AnyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	case PTypeVarchar:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("float8 cast to varchar: %w", err)
		}
		copy(buf, s)
		return AnyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("float8: unsupported cast to OID %d", targetType)
	}
}

func Float8Factory(buf []byte) TypeFloat8 {
	return TypeFloat8{BufferPtr: buf}
}

func Float8Comparator(a, b TypeFloat8) int {
	return bytes.Compare(a.BufferPtr, b.BufferPtr)
}
