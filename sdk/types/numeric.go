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

// NumericValue — распакованное значение
// Value хранит целое число = реальное_значение * 10^Scale
// Например: 123.456 при Scale=6 → Value = 123456000, Neg = false
type NumericValue struct {
	Value *big.Int
	Scale int32
	Neg   bool
}

func (n NumericValue) ToBigFloat() *big.Float {
	f := new(big.Float).SetPrec(256).SetInt(n.Value)
	if n.Scale > 0 {
		divisor := new(big.Float).SetPrec(256).SetInt(pow10(n.Scale))
		f.Quo(f, divisor)
	}
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

func (t TypeNumeric) GetType() OID      { return PTypeNumeric }
func (t TypeNumeric) GetBuffer() []byte { return t.BufferPtr }
func (t TypeNumeric) IntoGo() []byte    { return t.BufferPtr }

func (t TypeNumeric) Compare(other BaseType[[]byte]) int {
	return bytes.Compare(t.BufferPtr, other.GetBuffer())
}

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

// IsZero

func (t TypeNumeric) IsZero() bool {
	if len(t.BufferPtr) < 1 {
		return true
	}
	return t.BufferPtr[0] == 0x01
}

// ToNumericValue распаковывает буфер в NumericValue
// Value = abs(реальное_значение) * 10^Scale (целое)
func (t TypeNumeric) ToNumericValue() (*NumericValue, error) {
	if len(t.BufferPtr) < 1 {
		return nil, fmt.Errorf("numeric: empty buffer")
	}
	sign := t.BufferPtr[0]
	mag := make([]byte, len(t.BufferPtr)-1)
	copy(mag, t.BufferPtr[1:])

	// для отрицательных magnitude хранится инвертированным
	if sign == 0x00 {
		for i := range mag {
			mag[i] ^= 0xFF
		}
	}

	return &NumericValue{
		Value: new(big.Int).SetBytes(mag),
		Scale: t.Meta.Scale, // Scale берём из Meta — он не хранится в буфере
		Neg:   sign == 0x00,
	}, nil
}

// numericFromScaledInt создаёт TypeNumeric из целого Value = реальное * 10^Scale
func numericFromScaledInt(allocator pmem.Allocator, meta NumericMeta, value *big.Int, neg bool) (TypeNumeric, error) {
	// если значение нулевое — знак неважен
	if value.Sign() == 0 {
		buf, err := allocator.Alloc(1)
		if err != nil {
			return TypeNumeric{}, fmt.Errorf("numeric: alloc failed: %w", err)
		}
		buf[0] = 0x01 // zero sign byte
		return TypeNumeric{BufferPtr: buf, Meta: meta}, nil
	}

	mag := value.Bytes()

	var signByte byte
	if neg {
		signByte = 0x00
		// инвертируем magnitude для order-preserving
		for i := range mag {
			mag[i] ^= 0xFF
		}
	} else {
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

// numericFromBigFloat конвертирует big.Float → TypeNumeric
// f должен быть уже в реальных единицах (не масштабированным)
func numericFromBigFloat(allocator pmem.Allocator, meta NumericMeta, f *big.Float) (TypeNumeric, error) {
	neg := f.Sign() < 0
	if neg {
		f.Neg(f)
	}

	// умножаем на 10^Scale чтобы получить целое Value
	scale := new(big.Float).SetPrec(256).SetInt(pow10(meta.Scale))
	scaled := new(big.Float).SetPrec(256).Mul(f, scale)

	intVal, _ := scaled.Int(nil)
	return numericFromScaledInt(allocator, meta, intVal, neg)
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
	bOther, err := TypeNumeric{BufferPtr: other.GetBuffer(), Meta: t.Meta}.toBigFloat()
	if err != nil {
		return nil, err
	}
	result := new(big.Float).SetPrec(256).Add(a, bOther)
	return numericFromBigFloat(allocator, t.Meta, result)
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
	// Mod работает на целых Value (уже масштабированных)
	a, err := t.ToNumericValue()
	if err != nil {
		return nil, err
	}
	b, err := TypeNumeric{BufferPtr: other.GetBuffer(), Meta: t.Meta}.ToNumericValue()
	if err != nil {
		return nil, err
	}

	aVal := new(big.Int).Set(a.Value)
	bVal := new(big.Int).Set(b.Value)
	if a.Neg {
		aVal.Neg(aVal)
	}
	if b.Neg {
		bVal.Neg(bVal)
	}

	modVal := new(big.Int).Mod(aVal, bVal)
	neg := modVal.Sign() < 0
	if neg {
		modVal.Neg(modVal)
	}
	return numericFromScaledInt(allocator, t.Meta, modVal, neg)
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

func (t TypeNumeric) String() string {
	f, err := t.toBigFloat()
	if err != nil {
		return "numeric(invalid)"
	}
	return "numeric(" + f.Text('f', int(t.Meta.Scale)) + ")"
}
