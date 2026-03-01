// types/float.go
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

type TypeFloat4 struct {
	BufferPtr []byte
}

var _ BaseType[float32] = (*TypeFloat4)(nil)
var _ NumericType[float32] = (*TypeFloat4)(nil)
var _ OrderedType[float32] = (*TypeFloat4)(nil)
var _ NullableType[float32] = (*TypeFloat4)(nil)

func (t TypeFloat4) GetType() OID { return PTypeFloat4 }

func (t TypeFloat4) Compare(other BaseType[float32]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeFloat4) GetBuffer() []byte { return t.BufferPtr }

func (t TypeFloat4) IntoGo() float32 {
	if len(t.BufferPtr) < 4 {
		return 0
	}
	return Float32fromOrderableBits(binary.BigEndian.Uint32(t.BufferPtr))
}

// NullableType

func (t TypeFloat4) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeFloat4) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeFloat4) LessThan(other BaseType[float32]) bool       { return t.Compare(other) < 0 }
func (t TypeFloat4) GreaterThan(other BaseType[float32]) bool    { return t.Compare(other) > 0 }
func (t TypeFloat4) LessOrEqual(other BaseType[float32]) bool    { return t.Compare(other) <= 0 }
func (t TypeFloat4) GreaterOrEqual(other BaseType[float32]) bool { return t.Compare(other) >= 0 }
func (t TypeFloat4) Between(low, high BaseType[float32]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func float4FromVal(allocator pmem.Allocator, v float32) (TypeFloat4, error) {
	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return TypeFloat4{}, fmt.Errorf("float4: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint32(buf, OrderableFloat32bits(v))
	return TypeFloat4{BufferPtr: buf}, nil
}

// NumericType

func (t TypeFloat4) Add(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	return float4FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeFloat4) Sub(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	return float4FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeFloat4) Mul(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	return float4FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeFloat4) Div(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float4: division by zero")
	}
	return float4FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeFloat4) Mod(allocator pmem.Allocator, other NumericType[float32]) (NumericType[float32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("float4: modulo by zero")
	}
	return float4FromVal(allocator, float32(math.Mod(float64(t.IntoGo()), float64(other.IntoGo()))))
}

func (t TypeFloat4) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeFloat4) Neg(allocator pmem.Allocator) NumericType[float32] {
	v, _ := float4FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeFloat4) Abs(allocator pmem.Allocator) NumericType[float32] {
	v, _ := float4FromVal(allocator, float32(math.Abs(float64(t.IntoGo()))))
	return v
}

func (t TypeFloat4) String() string {
	return "float4(" + fmt.Sprintf("%v", t.IntoGo()) + ")"
}

var _ CastableType[float32] = (*TypeFloat4)(nil)

// CastableType

func (t TypeFloat4) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	v := t.IntoGo()

	switch targetType {
	case PTypeFloat8:
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("float4 cast to float8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, OrderableFloat64bits(float64(v)))
		return anyWrapper[float64]{TypeFloat8{BufferPtr: buf}}, nil

	case PTypeInt2:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return nil, fmt.Errorf("float4 cast to int2: value %v out of range", v)
		}
		buf, err := allocator.AllocAligned(2, 2)
		if err != nil {
			return nil, fmt.Errorf("float4 cast to int2: %w", err)
		}
		binary.BigEndian.PutUint16(buf, uint16(int16(v))^0x8000)
		return anyWrapper[int16]{TypeInt2{BufferPtr: buf}}, nil

	case PTypeInt4:
		if v < math.MinInt32 || v > math.MaxInt32 {
			return nil, fmt.Errorf("float4 cast to int4: value %v out of range", v)
		}
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("float4 cast to int4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, uint32(int32(v))^0x80000000)
		return anyWrapper[int32]{TypeInt4{BufferPtr: buf}}, nil

	case PTypeInt8:
		if v < math.MinInt64 || v > math.MaxInt64 {
			return nil, fmt.Errorf("float4 cast to int8: value %v out of range", v)
		}
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("float4 cast to int8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(int64(v))^0x8000000000000000)
		return anyWrapper[int64]{TypeInt8{BufferPtr: buf}}, nil

	case PTypeNumeric:
		meta := NumericMeta{Precision: 38, Scale: 10}
		s := strconv.FormatFloat(float64(v), 'f', 10, 32)
		f, _, err := new(big.Float).SetPrec(256).Parse(s, 10)
		if err != nil {
			return nil, fmt.Errorf("float4 cast to numeric: %w", err)
		}
		result, err := numericFromBigFloat(allocator, meta, f)
		if err != nil {
			return nil, fmt.Errorf("float4 cast to numeric: %w", err)
		}
		return anyWrapper[[]byte]{result}, nil

	case PTypeText:
		s := strconv.FormatFloat(float64(v), 'f', -1, 32)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("float4 cast to text: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	case PTypeVarchar:
		s := strconv.FormatFloat(float64(v), 'f', -1, 32)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("float4 cast to varchar: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("float4: unsupported cast to OID %d", targetType)
	}
}
