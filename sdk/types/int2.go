// types/int2.go
package ptypes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"petacore/sdk/pmem"
	"strconv"
)

type TypeInt2 struct {
	BufferPtr []byte
}

var _ BaseType[int16] = (*TypeInt2)(nil)
var _ NumericType[int16] = (*TypeInt2)(nil)
var _ OrderedType[int16] = (*TypeInt2)(nil)
var _ NullableType[int16] = (*TypeInt2)(nil)
var _ BitwiseType[int16] = (*TypeInt2)(nil)

func (t TypeInt2) GetType() OID { return PTypeInt2 }

func (t TypeInt2) Compare(other BaseType[int16]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeInt2) GetBuffer() []byte { return t.BufferPtr }

func (t TypeInt2) IntoGo() int16 {
	if len(t.BufferPtr) < 2 {
		return 0
	}
	return int16(binary.BigEndian.Uint16(t.BufferPtr) ^ 0x8000)
}

// NullableType

func (t TypeInt2) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeInt2) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeInt2) LessThan(other BaseType[int16]) bool       { return t.Compare(other) < 0 }
func (t TypeInt2) GreaterThan(other BaseType[int16]) bool    { return t.Compare(other) > 0 }
func (t TypeInt2) LessOrEqual(other BaseType[int16]) bool    { return t.Compare(other) <= 0 }
func (t TypeInt2) GreaterOrEqual(other BaseType[int16]) bool { return t.Compare(other) >= 0 }
func (t TypeInt2) Between(low, high BaseType[int16]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// helpers

func int2FromVal(allocator pmem.Allocator, v int16) (TypeInt2, error) {
	buf, err := allocator.AllocAligned(2, 2)
	if err != nil {
		return TypeInt2{}, fmt.Errorf("int2: alloc failed: %w", err)
	}
	binary.BigEndian.PutUint16(buf, uint16(v)^0x8000)
	return TypeInt2{BufferPtr: buf}, nil
}

// NumericType

func (t TypeInt2) Add(allocator pmem.Allocator, other NumericType[int16]) (NumericType[int16], error) {
	return int2FromVal(allocator, t.IntoGo()+other.IntoGo())
}

func (t TypeInt2) Sub(allocator pmem.Allocator, other NumericType[int16]) (NumericType[int16], error) {
	return int2FromVal(allocator, t.IntoGo()-other.IntoGo())
}

func (t TypeInt2) Mul(allocator pmem.Allocator, other NumericType[int16]) (NumericType[int16], error) {
	return int2FromVal(allocator, t.IntoGo()*other.IntoGo())
}

func (t TypeInt2) Div(allocator pmem.Allocator, other NumericType[int16]) (NumericType[int16], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int2: division by zero")
	}
	return int2FromVal(allocator, t.IntoGo()/other.IntoGo())
}

func (t TypeInt2) Mod(allocator pmem.Allocator, other NumericType[int16]) (NumericType[int16], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("int2: modulo by zero")
	}
	return int2FromVal(allocator, t.IntoGo()%other.IntoGo())
}

func (t TypeInt2) IsZero() bool { return t.IntoGo() == 0 }

func (t TypeInt2) Neg(allocator pmem.Allocator) NumericType[int16] {
	v, _ := int2FromVal(allocator, -t.IntoGo())
	return v
}

func (t TypeInt2) Abs(allocator pmem.Allocator) NumericType[int16] {
	v := t.IntoGo()
	if v < 0 {
		v = -v
	}
	result, _ := int2FromVal(allocator, v)
	return result
}

// BitwiseType

func (t TypeInt2) And(allocator pmem.Allocator, other BitwiseType[int16]) BitwiseType[int16] {
	v, _ := int2FromVal(allocator, t.IntoGo()&other.IntoGo())
	return v
}

func (t TypeInt2) Or(allocator pmem.Allocator, other BitwiseType[int16]) BitwiseType[int16] {
	v, _ := int2FromVal(allocator, t.IntoGo()|other.IntoGo())
	return v
}

func (t TypeInt2) Xor(allocator pmem.Allocator, other BitwiseType[int16]) BitwiseType[int16] {
	v, _ := int2FromVal(allocator, t.IntoGo()^other.IntoGo())
	return v
}

func (t TypeInt2) Not(allocator pmem.Allocator) BitwiseType[int16] {
	v, _ := int2FromVal(allocator, ^t.IntoGo())
	return v
}

func (t TypeInt2) ShiftLeft(allocator pmem.Allocator, n uint) BitwiseType[int16] {
	v, _ := int2FromVal(allocator, t.IntoGo()<<n)
	return v
}

func (t TypeInt2) ShiftRight(allocator pmem.Allocator, n uint) BitwiseType[int16] {
	v, _ := int2FromVal(allocator, t.IntoGo()>>n)
	return v
}

func (t TypeInt2) String() string {
	return "int2(" + fmt.Sprintf("%v", t.IntoGo()) + ")"
}

var _ CastableType[int16] = (*TypeInt2)(nil)

// CastableType

func (t TypeInt2) CastTo(allocator pmem.Allocator, targetType OID) (BaseType[any], error) {
	v := t.IntoGo()

	switch targetType {
	case PTypeInt4:
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("int2 cast to int4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, uint32(int32(v))^0x80000000)
		return anyWrapper[int32]{TypeInt4{BufferPtr: buf}}, nil

	case PTypeInt8:
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("int2 cast to int8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, uint64(int64(v))^0x8000000000000000)
		return anyWrapper[int64]{TypeInt8{BufferPtr: buf}}, nil

	case PTypeFloat4:
		buf, err := allocator.AllocAligned(4, 4)
		if err != nil {
			return nil, fmt.Errorf("int2 cast to float4: %w", err)
		}
		binary.BigEndian.PutUint32(buf, OrderableFloat32bits(float32(v)))
		return anyWrapper[float32]{TypeFloat4{BufferPtr: buf}}, nil

	case PTypeFloat8:
		buf, err := allocator.AllocAligned(8, 8)
		if err != nil {
			return nil, fmt.Errorf("int2 cast to float8: %w", err)
		}
		binary.BigEndian.PutUint64(buf, OrderableFloat64bits(float64(v)))
		return anyWrapper[float64]{TypeFloat8{BufferPtr: buf}}, nil

	case PTypeNumeric:
		meta := NumericMeta{Precision: 38, Scale: 0}
		f := new(big.Float).SetPrec(256).SetInt64(int64(v))
		result, err := numericFromBigFloat(allocator, meta, f)
		if err != nil {
			return nil, fmt.Errorf("int2 cast to numeric: %w", err)
		}
		return anyWrapper[[]byte]{result}, nil

	case PTypeBool:
		buf, err := allocator.Alloc(1)
		if err != nil {
			return nil, fmt.Errorf("int2 cast to bool: %w", err)
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
			return nil, fmt.Errorf("int2 cast to text: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeText{BufferPtr: buf}}, nil

	case PTypeVarchar:
		s := strconv.FormatInt(int64(v), 10)
		buf, err := allocator.Alloc(len(s))
		if err != nil {
			return nil, fmt.Errorf("int2 cast to varchar: %w", err)
		}
		copy(buf, s)
		return anyWrapper[string]{TypeVarchar{BufferPtr: buf}}, nil

	default:
		return nil, fmt.Errorf("int2: unsupported cast to OID %d", targetType)
	}
}
