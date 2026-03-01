// types/int4.go
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

type TypeInt4 struct {
	BufferPtr []byte
}

var _ BaseType[int32] = (*TypeInt4)(nil)
var _ NumericType[int32] = (*TypeInt4)(nil)
var _ OrderedType[int32] = (*TypeInt4)(nil)
var _ NullableType[int32] = (*TypeInt4)(nil)
var _ BitwiseType[int32] = (*TypeInt4)(nil)

func (t TypeInt4) GetType() OID { return PTypeInt4 }

func (t TypeInt4) Compare(other BaseType[int32]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeInt4) GetBuffer() []byte { return t.BufferPtr }

func (t TypeInt4) IntoGo() int32 {
	if len(t.BufferPtr) < 4 {
		return 0
	}
	return int32(binary.BigEndian.Uint32(t.BufferPtr) ^ 0x80000000)
}

// NullableType

func (t TypeInt4) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeInt4) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeInt4) LessThan(other BaseType[int32]) bool       { return t.Compare(other) < 0 }
func (t TypeInt4) GreaterThan(other BaseType[int32]) bool    { return t.Compare(other) > 0 }
func (t TypeInt4) LessOrEqual(other BaseType[int32]) bool    { return t.Compare(other) <= 0 }
func (t TypeInt4) GreaterOrEqual(other BaseType[int32]) bool { return t.Compare(other) >= 0 }
func (t TypeInt4) Between(low, high BaseType[int32]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func int4FromVal(allocator pmem.Allocator, v int32) (TypeInt4, error) {
	buf, err := allocator.AllocAligned(4, 4)
	if err != nil {
		return TypeInt4{}, fmt.Errorf("int4: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint32(buf, uint32(v)^0x80000000)
	return TypeInt4{BufferPtr: buf}, nil
}

// NumericType

func (t TypeInt4) Add(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	return int4FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeInt4) Sub(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	return int4FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeInt4) Mul(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	return int4FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeInt4) Div(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int4: division by zero")
	}
	return int4FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeInt4) Mod(allocator pmem.Allocator, other NumericType[int32]) (NumericType[int32], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int4: modulo by zero")
	}
	return int4FromVal(allocator, t.IntoGo()%other.IntoGo())
}

func (t TypeInt4) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeInt4) Neg(allocator pmem.Allocator) NumericType[int32] {
	v, _ := int4FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeInt4) Abs(allocator pmem.Allocator) NumericType[int32] {
	v := t.IntoGo()
	if v < 0 {
		v = -v
	}
	result, _ := int4FromVal(allocator, v)
	return result
}

// BitwiseType

func (t TypeInt4) And(allocator pmem.Allocator, other BitwiseType[int32]) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()&other.IntoGo())
	return v
}

func (t TypeInt4) Or(allocator pmem.Allocator, other BitwiseType[int32]) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()|other.IntoGo())
	return v
}

func (t TypeInt4) Xor(allocator pmem.Allocator, other BitwiseType[int32]) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()^other.IntoGo())
	return v
}

func (t TypeInt4) Not(allocator pmem.Allocator) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, ^t.IntoGo())
	return v
}

func (t TypeInt4) ShiftLeft(allocator pmem.Allocator, n uint) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()<<n)
	return v
}

func (t TypeInt4) ShiftRight(allocator pmem.Allocator, n uint) BitwiseType[int32] {
	v, _ := int4FromVal(allocator, t.IntoGo()>>n)
	return v
}

func (t TypeInt4) String() string {
	return "int4(" + fmt.Sprintf("%v", t.IntoGo()) + ")"
}

var _ CastableType[int32] = (*TypeInt4)(nil)

// CastableType

func (t TypeInt4) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	v := t.IntoGo()

	switch targetType {
	case PTypeInt2:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return nil, fmt.Errorf("int4 cast to int2: value %v out of range", v)
		}
		buf, err := allocator.AllocAligned(2, 2)
		if err != nil {
			return nil, fmt.Errorf("int4 cast to int2: %w", err)
		}
		binary.BigEndian.PutUint16(buf, uint16(int16(v))^0x8000)
		return anyWrapper[int16]{TypeInt2{BufferPtr: buf}}, nil

	case PTypeInt8:
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("int4 cast to int8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(int64(v))^0x8000000000000000)
		return anyWrapper[int64]{TypeInt8{BufferPtr: buf}}, nil

	case PTypeFloat4:
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("int4 cast to float4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, OrderableFloat32bits(float32(v)))
		return anyWrapper[float32]{TypeFloat4{BufferPtr: buf}}, nil

	case PTypeFloat8:
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("int4 cast to float8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, OrderableFloat64bits(float64(v)))
		return anyWrapper[float64]{TypeFloat8{BufferPtr: buf}}, nil

	case PTypeNumeric:
		meta := NumericMeta{Precision: 38, Scale: 0}
		f := new(big.Float).SetPrec(256).SetInt64(int64(v))
		result, err := numericFromBigFloat(allocator, meta, f)
		if err != nil {
			return nil, fmt.Errorf("int4 cast to numeric: %w", err)
		}
		return anyWrapper[[]byte]{result}, nil

	case PTypeBool:
		buf, err := allocator.Alloc(1)
		if err != nil {
			return nil, fmt.Errorf("int4 cast to bool: %w", err)
		}
		if v != 0 {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
		return anyWrapper[bool]{TypeBool{BufferPtr: buf}}, nil

	case PTypeText:
		s := strconv.FormatInt(int64(v), 10)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("int4 cast to text: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	case PTypeVarchar:
		s := strconv.FormatInt(int64(v), 10)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("int4 cast to varchar: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("int4: unsupported cast to OID %d", targetType)
	}
}
