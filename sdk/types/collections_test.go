package ptypes_test

import (
	"testing"

	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

// ============================================================
// helpers
// ============================================================

func boolFactory(buf []byte) ptypes.TypeBool {
	return ptypes.TypeBool{BufferPtr: buf}
}

// makeArray создаёт TypeArray через сериализатор
func makeArray(t *testing.T, arena pmem.Allocator, vals ...bool) *ptypes.TypeArray[bool, ptypes.TypeBool] {
	t.Helper()
	elements := make([][]byte, len(vals))
	for i, v := range vals {
		buf, err := serializers.BoolSerializerInstance.Serialize(arena, v)
		if err != nil {
			t.Fatalf("serialize element %d: %v", i, err)
		}
		elements[i] = buf
	}
	buf, err := ptypes.SerializeArrayElements(arena, ptypes.PTypeBool, elements)
	if err != nil {
		t.Fatalf("serialize array: %v", err)
	}
	return &ptypes.TypeArray[bool, ptypes.TypeBool]{
		BufferPtr: buf,
		Factory:   boolFactory,
		// семантический компаратор для bool
		Comparator: func(a, b ptypes.TypeBool) int {
			return a.Compare(b)
		},
	}
}

// makeBoolElem создаёт BaseType[TypeBool] через сериализатор
func makeBoolElem(t *testing.T, arena pmem.Allocator, val bool) ptypes.BaseType[ptypes.TypeBool] {
	t.Helper()
	// Serialize(arena, bool) → []byte
	buf, err := serializers.BoolSerializerInstance.Serialize(arena, val)
	if err != nil {
		t.Fatalf("serialize bool elem: %v", err)
	}
	// Deserialize([]byte) → TypeBool
	typeBool, err := serializers.BoolSerializerInstance.Deserialize(buf)
	if err != nil {
		t.Fatalf("deserialize bool elem: %v", err)
	}
	return &boolElemWrapper{val: typeBool}
}

type boolElemWrapper struct {
	val ptypes.TypeBool
}

func (w *boolElemWrapper) GetType() ptypes.OID { return ptypes.PTypeBool }
func (w *boolElemWrapper) GetBuffer() []byte   { return w.val.GetBuffer() }
func (w *boolElemWrapper) Compare(other ptypes.BaseType[ptypes.TypeBool]) int {
	return w.val.Compare(other.IntoGo())
}
func (w *boolElemWrapper) IntoGo() ptypes.TypeBool { return w.val }

// ============================================================
// Len / GetType
// ============================================================

func TestTypeArray_Len_Empty(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena)

	if arr.Len() != 0 {
		t.Fatalf("expected Len=0, got %d", arr.Len())
	}
}

func TestTypeArray_Len(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true, false, true)

	if arr.Len() != 3 {
		t.Fatalf("expected Len=3, got %d", arr.Len())
	}
}

func TestTypeArray_GetType(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true)

	if arr.GetType() != ptypes.PTypeBool {
		t.Fatalf("expected PTypeBool, got %d", arr.GetType())
	}
}

// ============================================================
// NullableType
// ============================================================

func TestTypeArray_IsNull(t *testing.T) {
	arr := &ptypes.TypeArray[bool, ptypes.TypeBool]{BufferPtr: nil, Factory: boolFactory}

	if !arr.IsNull() {
		t.Fatal("expected IsNull=true for nil buffer")
	}
	if arr.IsNotNull() {
		t.Fatal("expected IsNotNull=false for nil buffer")
	}
}

func TestTypeArray_IsNotNull(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true)

	if arr.IsNull() {
		t.Fatal("expected IsNull=false")
	}
	if !arr.IsNotNull() {
		t.Fatal("expected IsNotNull=true")
	}
}

// ============================================================
// GetPos / GetPosBuffer
// ============================================================

func TestTypeArray_GetPosBuffer(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true, false, true)

	buf, err := arr.GetPosBuffer(0)
	if err != nil {
		t.Fatalf("GetPosBuffer(0): %v", err)
	}
	elem, _ := serializers.BoolSerializerInstance.Deserialize(buf)
	if !elem.IntoGo() {
		t.Fatal("expected element 0 = true")
	}

	buf, err = arr.GetPosBuffer(1)
	if err != nil {
		t.Fatalf("GetPosBuffer(1): %v", err)
	}
	elem, _ = serializers.BoolSerializerInstance.Deserialize(buf)
	if elem.IntoGo() {
		t.Fatal("expected element 1 = false")
	}
}

func TestTypeArray_GetPosBuffer_OutOfBounds(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true)

	_, err := arr.GetPosBuffer(5)
	if err == nil {
		t.Fatal("expected out of bounds error")
	}

	_, err = arr.GetPosBuffer(-1)
	if err == nil {
		t.Fatal("expected error for negative index")
	}
}

func TestTypeArray_GetPos(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, false, true)

	elem := arr.GetPos(1)
	if elem == nil {
		t.Fatal("expected non-nil element at index 1")
	}
	b, _ := serializers.BoolSerializerInstance.Deserialize(elem.GetBuffer())
	if !b.IntoGo() {
		t.Fatal("expected element 1 = true")
	}
}

func TestTypeArray_GetPos_OutOfBounds(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true)

	elem := arr.GetPos(10)
	if elem != nil {
		t.Fatal("expected nil for out of bounds GetPos")
	}
}

// ============================================================
// IntoGo
// ============================================================

func TestTypeArray_IntoGo(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true, false, true)

	vals := arr.IntoGo()
	if len(vals) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(vals))
	}
	if !vals[0].IntoGo() {
		t.Fatal("expected vals[0]=true")
	}
	if vals[1].IntoGo() {
		t.Fatal("expected vals[1]=false")
	}
	if !vals[2].IntoGo() {
		t.Fatal("expected vals[2]=true")
	}
}

// ============================================================
// SetPos
// ============================================================

func TestTypeArray_SetPos(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, true, true, true)

	newVal := makeBoolElem(t, arena, false)
	result := arr.SetPos(arena, 1, newVal)

	buf, err := result.(*ptypes.TypeArray[bool, ptypes.TypeBool]).GetPosBuffer(1)
	if err != nil {
		t.Fatalf("GetPosBuffer after SetPos: %v", err)
	}
	elem, _ := serializers.BoolSerializerInstance.Deserialize(buf)
	if elem.IntoGo() {
		t.Fatal("expected element 1 = false after SetPos")
	}
}

func TestTypeArray_SetPos_OutOfBounds(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true)

	newVal := makeBoolElem(t, arena, false)
	result := arr.SetPos(arena, 10, newVal)

	// должен вернуть оригинальный массив без изменений
	if result.Len() != arr.Len() {
		t.Fatal("expected unchanged array for out of bounds SetPos")
	}
}

// ============================================================
// DeletePos
// ============================================================

func TestTypeArray_DeletePos(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, true, false, true)

	result := arr.DeletePos(arena, 1)

	if result.Len() != 2 {
		t.Fatalf("expected Len=2 after delete, got %d", result.Len())
	}

	buf, _ := result.(*ptypes.TypeArray[bool, ptypes.TypeBool]).GetPosBuffer(0)
	elem, _ := serializers.BoolSerializerInstance.Deserialize(buf)
	if !elem.IntoGo() {
		t.Fatal("expected element 0 = true after delete")
	}

	buf, _ = result.(*ptypes.TypeArray[bool, ptypes.TypeBool]).GetPosBuffer(1)
	elem, _ = serializers.BoolSerializerInstance.Deserialize(buf)
	if !elem.IntoGo() {
		t.Fatal("expected element 1 = true after delete")
	}
}

func TestTypeArray_DeletePos_First(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, false, true, true)

	result := arr.DeletePos(arena, 0)

	if result.Len() != 2 {
		t.Fatalf("expected Len=2, got %d", result.Len())
	}
	buf, _ := result.(*ptypes.TypeArray[bool, ptypes.TypeBool]).GetPosBuffer(0)
	elem, _ := serializers.BoolSerializerInstance.Deserialize(buf)
	if !elem.IntoGo() {
		t.Fatal("expected element 0 = true after deleting first")
	}
}

func TestTypeArray_DeletePos_OutOfBounds(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true, false)

	result := arr.DeletePos(arena, 10)
	if result.Len() != 2 {
		t.Fatal("expected unchanged array for out of bounds DeletePos")
	}
}

// ============================================================
// Append
// ============================================================

func TestTypeArray_Append(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, true, false)

	newVal := makeBoolElem(t, arena, true)
	result, err := arr.Append(arena, newVal)
	if err != nil {
		t.Fatalf("Append: %v", err)
	}

	if result.Len() != 3 {
		t.Fatalf("expected Len=3 after append, got %d", result.Len())
	}

	buf, _ := result.(*ptypes.TypeArray[bool, ptypes.TypeBool]).GetPosBuffer(2)
	elem, _ := serializers.BoolSerializerInstance.Deserialize(buf)
	if !elem.IntoGo() {
		t.Fatal("expected appended element = true")
	}
}

func TestTypeArray_Append_ToEmpty(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena)

	newVal := makeBoolElem(t, arena, false)
	result, err := arr.Append(arena, newVal)
	if err != nil {
		t.Fatalf("Append to empty: %v", err)
	}

	if result.Len() != 1 {
		t.Fatalf("expected Len=1, got %d", result.Len())
	}
}

func TestTypeArray_Append_OOM(t *testing.T) {
	arena := pmem.NewArena(0)
	arr := &ptypes.TypeArray[bool, ptypes.TypeBool]{
		BufferPtr: mustSerializeArray(t, pmem.NewArena(1024), ptypes.PTypeBool),
		Factory:   boolFactory,
	}
	newVal := makeBoolElem(t, pmem.NewArena(64), true)

	_, err := arr.Append(arena, newVal)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

// ============================================================
// Slice
// ============================================================

func TestTypeArray_Slice(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, true, false, true, false)

	result, err := arr.Slice(arena, 1, 2)
	if err != nil {
		t.Fatalf("Slice: %v", err)
	}

	if result.Len() != 2 {
		t.Fatalf("expected Len=2, got %d", result.Len())
	}

	buf, _ := result.(*ptypes.TypeArray[bool, ptypes.TypeBool]).GetPosBuffer(0)
	elem, _ := serializers.BoolSerializerInstance.Deserialize(buf)
	if elem.IntoGo() {
		t.Fatal("expected slice[0]=false")
	}

	buf, _ = result.(*ptypes.TypeArray[bool, ptypes.TypeBool]).GetPosBuffer(1)
	elem, _ = serializers.BoolSerializerInstance.Deserialize(buf)
	if !elem.IntoGo() {
		t.Fatal("expected slice[1]=true")
	}
}

func TestTypeArray_Slice_OutOfBounds(t *testing.T) {
	arena := newArena()
	arr := makeArray(t, arena, true, false)

	_, err := arr.Slice(arena, 0, 10)
	if err == nil {
		t.Fatal("expected out of bounds error")
	}

	_, err = arr.Slice(arena, -1, 1)
	if err == nil {
		t.Fatal("expected error for negative start")
	}
}

func TestTypeArray_Slice_Full(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, true, false, true)

	result, err := arr.Slice(arena, 0, 3)
	if err != nil {
		t.Fatalf("Slice full: %v", err)
	}
	if result.Len() != 3 {
		t.Fatalf("expected Len=3, got %d", result.Len())
	}
}

// ============================================================
// Contains
// ============================================================

func TestTypeArray_Contains_True(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, true, false, true)

	falseElem := makeBoolElem(t, arena, false)
	if !arr.Contains(falseElem) {
		t.Fatal("expected Contains=true for false element")
	}
}

func TestTypeArray_Contains_False(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena, true, true, true)

	falseElem := makeBoolElem(t, arena, false)
	if arr.Contains(falseElem) {
		t.Fatal("expected Contains=false")
	}
}

func TestTypeArray_Contains_Empty(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	arr := makeArray(t, arena)

	trueElem := makeBoolElem(t, arena, true)
	if arr.Contains(trueElem) {
		t.Fatal("expected Contains=false for empty array")
	}
}

// ============================================================
// OrderedType
// ============================================================

func TestTypeArray_Compare(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeArray(t, arena, false)
	b := makeArray(t, arena, true)

	if a.Compare(b) >= 0 {
		t.Fatal("expected a < b")
	}
	if b.Compare(a) <= 0 {
		t.Fatal("expected b > a")
	}
	if a.Compare(a) != 0 {
		t.Fatal("expected a == a")
	}
}

func TestTypeArray_LessThan(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeArray(t, arena, false)
	b := makeArray(t, arena, true)

	if !a.LessThan(b) {
		t.Fatal("expected a < b")
	}
	if b.LessThan(a) {
		t.Fatal("expected b not < a")
	}
}

func TestTypeArray_Between(t *testing.T) {
	arena := pmem.NewArena(64 * 1024)
	a := makeArray(t, arena, false)
	b := makeArray(t, arena, false, true)
	c := makeArray(t, arena, true)

	if !b.Between(a, c) {
		t.Fatal("expected b between a and c")
	}
	if !a.Between(a, c) {
		t.Fatal("expected a between a and c (inclusive)")
	}
}

// ============================================================
// helpers для тестов
// ============================================================

func mustSerializeArray(t *testing.T, arena pmem.Allocator, innerType ptypes.OID, elements ...[]byte) []byte {
	t.Helper()
	buf, err := ptypes.SerializeArrayElements(arena, innerType, elements)
	if err != nil {
		t.Fatalf("mustSerializeArray: %v", err)
	}
	return buf
}
