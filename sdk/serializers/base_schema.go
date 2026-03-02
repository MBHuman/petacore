// serializers/schema.go
package serializers

import (
	"encoding/binary"
	"fmt"
	"petacore/sdk/pmem"
	ptypes "petacore/sdk/types"
)

// FieldDef описывает поле схемы
type FieldDef struct {
	Name       string
	OID        ptypes.OID
	TableAlias string // optional table name or alias (e.g. "t", "pg_type")
}

// BaseSchema описывает структуру Row — порядок и типы полей
type BaseSchema struct {
	Fields  []FieldDef
	nameIdx map[string]int // быстрый поиск по имени — O(1)
}

func NewBaseSchema(fields []FieldDef) *BaseSchema {
	s := &BaseSchema{Fields: fields}
	s.RebuildIndex()
	return s
}

// RebuildIndex rebuilds the name→index map after fields (particularly TableAlias) are modified.
// Indexes both the bare name ("oid") and the qualified name ("pg_type.oid").
func (s *BaseSchema) RebuildIndex() {
	idx := make(map[string]int, len(s.Fields)*2)
	for i, f := range s.Fields {
		idx[f.Name] = i
		if f.TableAlias != "" {
			idx[f.TableAlias+"."+f.Name] = i
		}
	}
	s.nameIdx = idx
}

func (s *BaseSchema) Equal(other *BaseSchema) bool {
	if len(s.Fields) != len(other.Fields) {
		return false
	}
	for i := range s.Fields {
		if s.Fields[i] != other.Fields[i] {
			return false
		}
	}
	return true
}

// FieldIndex возвращает индекс поля по имени — O(1)
func (s *BaseSchema) FieldIndex(name string) (int, bool) {
	idx, ok := s.nameIdx[name]
	return idx, ok
}

// Pack упаковывает значения в Row буфер через аллокатор
// values — уже сериализованные байты каждого поля в порядке схемы
func (s *BaseSchema) Pack(allocator pmem.Allocator, values [][]byte) (*ptypes.Row, error) { // добавь в начало
	if allocator == nil {
		return nil, fmt.Errorf("schema pack: allocator is nil")
	}
	if s == nil {
		return nil, fmt.Errorf("schema pack: schema is nil")
	}
	// ...
	if len(values) != len(s.Fields) {
		return nil, fmt.Errorf("schema pack: expected %d fields, got %d", len(s.Fields), len(values))
	}

	count := len(s.Fields)
	dataSize := 0
	for _, v := range values {
		dataSize += len(v)
	}

	// 4 (count) + count*4 (OIDs) + count*4 (offsets) + count*4 (lengths) + data
	headerSize := 4 + count*4 + count*4 + count*4
	totalSize := headerSize + dataSize

	buf, err := allocator.Alloc(totalSize)
	if err != nil {
		return nil, fmt.Errorf("schema pack: %w", err)
	}

	// field_count
	binary.BigEndian.PutUint32(buf[0:4], uint32(count))

	// OIDs
	for i, f := range s.Fields {
		oidPos := 4 + i*4
		binary.BigEndian.PutUint32(buf[oidPos:oidPos+4], uint32(f.OID))
	}

	// offsets + lengths + data
	dataOffset := headerSize
	for i, val := range values {
		offsetPos := 4 + count*4 + i*4
		lengthPos := 4 + count*4 + count*4 + i*4
		binary.BigEndian.PutUint32(buf[offsetPos:offsetPos+4], uint32(dataOffset))
		binary.BigEndian.PutUint32(buf[lengthPos:lengthPos+4], uint32(len(val)))
		copy(buf[dataOffset:], val)
		dataOffset += len(val)
	}

	return &ptypes.Row{BufferPtr: buf}, nil
}

// GetField возвращает байты поля по индексу — O(1)
func (s *BaseSchema) GetField(row *ptypes.Row, idx int) ([]byte, ptypes.OID, error) {
	if idx < 0 || idx >= len(s.Fields) {
		return nil, 0, fmt.Errorf("schema: field index %d out of bounds", idx)
	}
	buf, err := row.GetFieldBuffer(idx)
	if err != nil {
		return nil, 0, err
	}
	oid, err := row.GetFieldOID(idx)
	if err != nil {
		return nil, 0, err
	}
	return buf, oid, nil
}

// GetFieldByName возвращает байты поля по имени — O(1)
func (s *BaseSchema) GetFieldByName(row *ptypes.Row, name string) ([]byte, ptypes.OID, error) {
	idx, ok := s.FieldIndex(name)
	if !ok {
		return nil, 0, fmt.Errorf("schema: field %q not found", name)
	}
	return s.GetField(row, idx)
}

// SetField заменяет поле inplace — O(1) если размер совпадает
// Row должен быть mutable (не из read-only памяти)
func (s *BaseSchema) SetField(row *ptypes.Row, idx int, value []byte) error {
	count := row.FieldCount()
	if idx < 0 || idx >= count {
		return fmt.Errorf("schema setfield: index %d out of bounds", idx)
	}

	lengthPos := 4 + count*4 + count*4 + idx*4
	existingLen := int(binary.BigEndian.Uint32(row.BufferPtr[lengthPos : lengthPos+4]))

	if len(value) != existingLen {
		return fmt.Errorf(
			"schema setfield: size mismatch for field %d %q: existing %d bytes, new %d bytes — use PackNew for resize",
			idx, s.Fields[idx].Name, existingLen, len(value),
		)
	}

	offsetPos := 4 + count*4 + idx*4
	offset := int(binary.BigEndian.Uint32(row.BufferPtr[offsetPos : offsetPos+4]))

	// пишем прямо в буфер — O(1) zero allocation
	copy(row.BufferPtr[offset:offset+existingLen], value)
	return nil
}

// PackNew создаёт новый Row с изменённым полем — используй когда размер меняется
func (s *BaseSchema) PackNew(allocator pmem.Allocator, row *ptypes.Row, idx int, value []byte) (*ptypes.Row, error) {
	count := row.FieldCount()
	if idx < 0 || idx >= count {
		return nil, fmt.Errorf("schema packnew: index %d out of bounds", idx)
	}

	values := make([][]byte, count)
	for i := 0; i < count; i++ {
		buf, err := row.GetFieldBuffer(i)
		if err != nil {
			return nil, err
		}
		values[i] = buf
	}
	values[idx] = value

	return s.Pack(allocator, values)
}

// Validate проверяет что буфер соответствует схеме
func (s *BaseSchema) Validate(row *ptypes.Row) error {
	if row.FieldCount() != len(s.Fields) {
		return fmt.Errorf("schema validate: expected %d fields, got %d", len(s.Fields), row.FieldCount())
	}
	for i, field := range s.Fields {
		oid, err := row.GetFieldOID(i)
		if err != nil {
			return fmt.Errorf("schema validate: field %d: %w", i, err)
		}
		if oid != field.OID {
			return fmt.Errorf("schema validate: field %d %q OID mismatch: expected %d, got %d",
				i, field.Name, field.OID, oid)
		}
	}
	return nil
}
