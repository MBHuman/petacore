// types/row.go
package ptypes

import (
	"encoding/binary"
	"fmt"
)

// Layout:
// [4 bytes: field_count]
// [field_count * 4 bytes: OIDs]
// [field_count * 4 bytes: offsets от начала буфера]
// [field_count * 4 bytes: lengths]
// [data...]

type Row struct {
	BufferPtr []byte
}

func RowFactory(buf []byte) *Row {
	return &Row{BufferPtr: buf}
}

func (r *Row) FieldCount() int {
	if len(r.BufferPtr) < 4 {
		return 0
	}
	return int(binary.BigEndian.Uint32(r.BufferPtr[:4]))
}

// GetFieldBuffer возвращает байты поля по индексу — O(1) zero copy
func (r *Row) GetFieldBuffer(idx int) ([]byte, error) {
	count := r.FieldCount()
	if idx < 0 || idx >= count {
		return nil, fmt.Errorf("row: index %d out of bounds [0, %d)", idx, count)
	}

	offsetPos := 4 + count*4 + idx*4
	lengthPos := 4 + count*4 + count*4 + idx*4

	if len(r.BufferPtr) < lengthPos+4 {
		return nil, fmt.Errorf("row: buffer too short")
	}

	offset := int(binary.BigEndian.Uint32(r.BufferPtr[offsetPos : offsetPos+4]))
	length := int(binary.BigEndian.Uint32(r.BufferPtr[lengthPos : lengthPos+4]))

	if offset+length > len(r.BufferPtr) {
		return nil, fmt.Errorf("row: field %d out of buffer bounds", idx)
	}

	return r.BufferPtr[offset : offset+length], nil
}

// GetFieldOID возвращает OID поля по индексу — O(1)
func (r *Row) GetFieldOID(idx int) (OID, error) {
	count := r.FieldCount()
	if idx < 0 || idx >= count {
		return 0, fmt.Errorf("row: index %d out of bounds [0, %d)", idx, count)
	}
	oidPos := 4 + idx*4
	return OID(binary.BigEndian.Uint32(r.BufferPtr[oidPos : oidPos+4])), nil
}

func (t Row) String() string {
	count := t.FieldCount()
	result := "row{"
	for i := 0; i < count; i++ {
		oid, err := t.GetFieldOID(i)
		if err != nil {
			result += fmt.Sprintf("field%d: error(%v)", i, err)
			continue
		}
		buf, err := t.GetFieldBuffer(i)
		if err != nil {
			result += fmt.Sprintf("field%d(oid %d): error(%v)", i, oid, err)
			continue
		}
		result += fmt.Sprintf("field%d(oid %d): %d bytes", i, oid, len(buf))
		if i < count-1 {
			result += ", "
		}
	}
	result += "}"
	return result
}
