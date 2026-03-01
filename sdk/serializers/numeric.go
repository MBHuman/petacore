// serializers/numeric.go
package serializers

import (
	"fmt"
	"math/big"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"strings"
)

type NumericSerializer struct {
	Meta ptypes.NumericMeta
}

func NewNumericSerializer(precision, scale int32) (*NumericSerializer, error) {
	meta := ptypes.NumericMeta{Precision: precision, Scale: scale}
	if err := meta.Validate(); err != nil {
		return nil, err
	}
	return &NumericSerializer{Meta: meta}, nil
}

var NumericSerializerInstance BaseSerializer[string, ptypes.TypeNumeric] = &NumericSerializer{
	Meta: ptypes.NumericMeta{Precision: 38, Scale: 10},
}

// Serialize принимает строку вида "123.456" или "-99.99"
func (s *NumericSerializer) Serialize(allocator pmem.Allocator, value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("numeric serialize: empty string")
	}

	// парсим через big.Float для точности
	f, _, err := new(big.Float).SetPrec(256).Parse(value, 10)
	if err != nil {
		return nil, fmt.Errorf("numeric serialize: invalid value %q: %w", value, err)
	}

	neg := f.Sign() < 0
	if neg {
		f.Neg(f)
	}

	// масштабируем: умножаем на 10^Scale чтобы убрать дробную часть
	scale := new(big.Float).SetPrec(256).SetInt(pow10Big(s.Meta.Scale))
	f.Mul(f, scale)

	// конвертируем в big.Int (усекаем дробную часть)
	intVal, _ := f.Int(nil)

	// проверяем precision
	digits := len(intVal.String())
	if neg && intVal.Sign() != 0 {
		digits = len(intVal.String())
	}
	if int32(digits) > s.Meta.Precision {
		return nil, fmt.Errorf("numeric serialize: value exceeds precision %d", s.Meta.Precision)
	}

	mag := intVal.Bytes() // big-endian magnitude

	// для отрицательных инвертируем байты — order-preserving
	if neg {
		for i := range mag {
			mag[i] ^= 0xFF
		}
	}

	// sign byte
	var signByte byte
	switch {
	case neg:
		signByte = 0x00
	case intVal.Sign() == 0:
		signByte = 0x01
	default:
		signByte = 0x02
	}

	totalSize := 1 + len(mag)
	buf, err := allocator.Alloc(totalSize)
	if err != nil {
		return nil, fmt.Errorf("numeric serialize: %w", err)
	}

	buf[0] = signByte
	copy(buf[1:], mag)
	return buf, nil
}

func (s *NumericSerializer) Deserialize(data []byte) (ptypes.TypeNumeric, error) {
	if len(data) < 1 {
		return ptypes.TypeNumeric{}, fmt.Errorf("numeric deserialize: empty data")
	}
	return ptypes.TypeNumeric{BufferPtr: data}, nil
}

func (s *NumericSerializer) Validate(value ptypes.TypeNumeric) error {
	if len(value.BufferPtr) < 1 {
		return fmt.Errorf("numeric validate: empty buffer")
	}
	sign := value.BufferPtr[0]
	if sign != 0x00 && sign != 0x01 && sign != 0x02 {
		return fmt.Errorf("numeric validate: invalid sign byte 0x%02X", sign)
	}
	return nil
}

func (s *NumericSerializer) GetType() ptypes.OID { return ptypes.PTypeNumeric }

func pow10Big(n int32) *big.Int {
	result := big.NewInt(1)
	ten := big.NewInt(10)
	for i := int32(0); i < n; i++ {
		result.Mul(result, ten)
	}
	return result
}
