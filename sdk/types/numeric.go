// types/numeric.go
package ptypes

import (
	"bytes"
	"fmt"
	"math/big"
	"petacore/sdk/pmem"
)

type NumericMeta struct {
	Precision int32
	Scale     int32
}

func (m NumericMeta) Validate() error {
	if m.Precision < 1 || m.Precision > 1000 {
		return fmt.Errorf("numeric: precision must be between 1 and 1000, got %d", m.Precision)
	}
	if m.Scale < 0 || m.Scale > m.Precision {
		return fmt.Errorf("numeric: scale must be between 0 and precision (%d), got %d", m.Precision, m.Scale)
	}
	return nil
}

type NumericValue struct {
	Value *big.Int
	Scale int32
	Neg   bool
}

func (n NumericValue) ToBigFloat() *big.Float {
	f := new(big.Float).SetInt(n.Value)
	divisor := new(big.Float).SetInt(pow10(n.Scale))
	f.Quo(f, divisor)
	if n.Neg {
		f.Neg(f)
	}
	return f
}

func pow10(n int32) *big.Int {
	result := big.NewInt(1)
	ten := big.NewInt(10)
	for i := int32(0); i < n; i++ {
		result.Mul(result, ten)
	}
	return result
}

type TypeNumeric struct {
	BufferPtr []byte
	Meta      NumericMeta
}

var _ BaseType[[]byte] = (*TypeNumeric)(nil)
var _ NumericType[[]byte] = (*TypeNumeric)(nil)
var _ OrderedType[[]byte] = (*TypeNumeric)(nil)
var _ NullableType[[]byte] = (*TypeNumeric)(nil)

func (t TypeNumeric) GetType() OID { return PTypeNumeric }

func (t TypeNumeric) Compare(other BaseType[[]byte]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

func (t TypeNumeric) GetBuffer() []byte { return t.BufferPtr }
func (t TypeNumeric) IntoGo() []byte    { return t.BufferPtr }

// NullableType

func (t TypeNumeric) IsNull() bool    { return t.BufferPtr == nil }
func (t TypeNumeric) IsNotNull() bool { return t.BufferPtr != nil }

// OrderedType

func (t TypeNumeric) LessThan(other BaseType[[]byte]) bool       { return t.Compare(other) < 0 }
func (t TypeNumeric) GreaterThan(other BaseType[[]byte]) bool    { return t.Compare(other) > 0 }
func (t TypeNumeric) LessOrEqual(other BaseType[[]byte]) bool    { return t.Compare(other) <= 0 }
func (t TypeNumeric) GreaterOrEqual(other BaseType[[]byte]) bool { return t.Compare(other) >= 0 }
func (t TypeNumeric) Between(low, high BaseType[[]byte]) bool {
	return t.GreaterOrEqual(low) && t.LessOrEqual(high)
}

// ToNumericValue распаковывает буфер в NumericValue
func (t TypeNumeric) ToNumericValue() (*NumericValue, error) {
	if len(t.BufferPtr) < 1 {
		return nil, fmt.Errorf("numeric: empty buffer")
	}
	sign := t.BufferPtr[0]
	mag := make([]byte, len(t.BufferPtr)-1)
	copy(mag, t.BufferPtr[1:])
	if sign == 0x00 {
		for i := range mag {
			mag[i] ^= 0xFF
		}
	}
	return &NumericValue{
		Value: new(big.Int).SetBytes(mag),
		Scale: t.Meta.Scale,
		Neg:   sign == 0x00,
	}, nil
}

// IsZero — только читает

func (t TypeNumeric) IsZero() bool {
	if len(t.BufferPtr) < 1 {
		return true
	}
	return t.BufferPtr[0] == 0x01 // sign byte для нуля
}

// helpers

func numericFromBigFloat(allocator pmem.Allocator, meta NumericMeta, f *big.Float) (TypeNumeric, error) {
	neg := f.Sign() < 0
	if neg {
		f.Neg(f)
	}

	scale := new(big.Float).SetPrec(256).SetInt(pow10(meta.Scale))
	f.Mul(f, scale)

	intVal, _ := f.Int(nil)
	mag := intVal.Bytes()

	if neg {
		for i := range mag {
			mag[i] ^= 0xFF
		}
	}

	var signByte byte
	switch {
	case neg:
		signByte = 0x00
	case intVal.Sign() == 0:
		signByte = 0x01
	default:
		signByte = 0x02
	}

	buf, err := allocator.Alloc(1 + len(mag))
	if err != nil {
		return TypeNumeric{}, fmt.Errorf("numeric: alloc failed: %w", err)
	}
	buf[0] = signByte
	copy(buf[1:], mag)

	return TypeNumeric{BufferPtr: buf, Meta: meta}, nil
}

func (t TypeNumeric) toBigFloat() (*big.Float, error) {
	v, err := t.ToNumericValue()
	if err != nil {
		return nil, err
	}
	return v.ToBigFloat(), nil
}

// NumericType

func (t TypeNumeric) Add(allocator pmem.Allocator, other NumericType[[]byte]) (NumericType[[]byte], error) {
	a, err := t.toBigFloat()
	if err != nil {
		return nil, err
	}
	b := new(big.Float).SetPrec(256)
	bOther, err := TypeNumeric{BufferPtr: other.GetBuffer(), Meta: t.Meta}.toBigFloat()
	if err != nil {
		return nil, err
	}
	b.Add(a, bOther)
	return numericFromBigFloat(allocator, t.Meta, b)
}

func (t TypeNumeric) Sub(allocator pmem.Allocator, other NumericType[[]byte]) (NumericType[[]byte], error) {
	a, err := t.toBigFloat()
	if err != nil {
		return nil, err
	}
	bOther, err := TypeNumeric{BufferPtr: other.GetBuffer(), Meta: t.Meta}.toBigFloat()
	if err != nil {
		return nil, err
	}
	result := new(big.Float).SetPrec(256).Sub(a, bOther)
	return numericFromBigFloat(allocator, t.Meta, result)
}

func (t TypeNumeric) Mul(allocator pmem.Allocator, other NumericType[[]byte]) (NumericType[[]byte], error) {
	a, err := t.toBigFloat()
	if err != nil {
		return nil, err
	}
	bOther, err := TypeNumeric{BufferPtr: other.GetBuffer(), Meta: t.Meta}.toBigFloat()
	if err != nil {
		return nil, err
	}
	result := new(big.Float).SetPrec(256).Mul(a, bOther)
	return numericFromBigFloat(allocator, t.Meta, result)
}

func (t TypeNumeric) Div(allocator pmem.Allocator, other NumericType[[]byte]) (NumericType[[]byte], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("numeric: division by zero")
	}
	a, err := t.toBigFloat()
	if err != nil {
		return nil, err
	}
	bOther, err := TypeNumeric{BufferPtr: other.GetBuffer(), Meta: t.Meta}.toBigFloat()
	if err != nil {
		return nil, err
	}
	result := new(big.Float).SetPrec(256).Quo(a, bOther)
	return numericFromBigFloat(allocator, t.Meta, result)
}

func (t TypeNumeric) Mod(allocator pmem.Allocator, other NumericType[[]byte]) (NumericType[[]byte], error) {
	if other.IsZero() {
		return nil, fmt.Errorf("numeric: modulo by zero")
	}
	a, err := t.ToNumericValue()
	if err != nil {
		return nil, err
	}
	b, err := TypeNumeric{BufferPtr: other.GetBuffer(), Meta: t.Meta}.ToNumericValue()
	if err != nil {
		return nil, err
	}
	// Mod через big.Int — точная операция
	result := new(big.Int).Mod(a.Value, b.Value)
	f := new(big.Float).SetPrec(256).SetInt(result)
	divisor := new(big.Float).SetPrec(256).SetInt(pow10(t.Meta.Scale))
	f.Quo(f, divisor)
	return numericFromBigFloat(allocator, t.Meta, f)
}

func (t TypeNumeric) Neg(allocator pmem.Allocator) NumericType[[]byte] {
	f, err := t.toBigFloat()
	if err != nil {
		return t
	}
	result, _ := numericFromBigFloat(allocator, t.Meta, f.Neg(f))
	return result
}

func (t TypeNumeric) Abs(allocator pmem.Allocator) NumericType[[]byte] {
	f, err := t.toBigFloat()
	if err != nil {
		return t
	}
	result, _ := numericFromBigFloat(allocator, t.Meta, f.Abs(f))
	return result
}
