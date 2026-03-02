// serializers/varchar.go
package serializers

import (
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
	"unicode/utf8"
)

type VarcharSerializer struct {
	Meta ptypes.VarcharMeta
}

// NewVarcharSerializer создаёт сериализатор с ограничением длины
// maxLength=0 означает неограниченный varchar
func NewVarcharSerializer(maxLength int32) *VarcharSerializer {
	return &VarcharSerializer{Meta: ptypes.VarcharMeta{MaxLength: maxLength}}
}

var VarcharSerializerInstance = &VarcharSerializer{
	Meta: ptypes.VarcharMeta{MaxLength: 0}, // без ограничения по умолчанию
}

func (s *VarcharSerializer) Serialize(allocator pmem.Allocator, value string) ([]byte, error) {
	if s.Meta.MaxLength > 0 {
		runeLen := utf8.RuneCountInString(value)
		if int32(runeLen) > s.Meta.MaxLength {
			return nil, fmt.Errorf("varchar serialize: value length %d exceeds max %d", runeLen, s.Meta.MaxLength)
		}
	}

	b := []byte(value)
	buf, err := allocator.Alloc(len(b))
	if err != nil {
		return nil, fmt.Errorf("varchar serialize: %w", err)
	}
	copy(buf, b)
	return buf, nil
}

func (s *VarcharSerializer) Deserialize(data []byte) (ptypes.TypeVarchar, error) {
	return ptypes.TypeVarchar{
		BufferPtr: data,
		Meta:      s.Meta,
	}, nil
}

func (s *VarcharSerializer) Validate(value ptypes.TypeVarchar) error {
	if value.BufferPtr == nil {
		return fmt.Errorf("varchar validate: nil buffer")
	}
	if s.Meta.MaxLength > 0 {
		runeLen := utf8.RuneCount(value.BufferPtr)
		if int32(runeLen) > s.Meta.MaxLength {
			return fmt.Errorf("varchar validate: value length %d exceeds max %d", runeLen, s.Meta.MaxLength)
		}
	}
	return nil
}

func (s *VarcharSerializer) GetType() ptypes.OID { return ptypes.PTypeVarchar }
