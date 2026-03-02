package pmem_test

import (
	"testing"

	"petacore/sdk/pmem"
)

// ============================================================
// ArenaAllocator tests
// ============================================================

func TestArenaAlloc_Basic(t *testing.T) {
	arena := pmem.NewArena(1024)

	buf, err := arena.Alloc(64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buf) != 64 {
		t.Fatalf("expected len 64, got %d", len(buf))
	}
}

func TestArenaAlloc_MultipleAllocs(t *testing.T) {
	arena := pmem.NewArena(1024)

	buf1, err := arena.Alloc(64)
	if err != nil {
		t.Fatalf("alloc1: %v", err)
	}
	buf2, err := arena.Alloc(64)
	if err != nil {
		t.Fatalf("alloc2: %v", err)
	}

	addr1 := pmem.BufAddr(buf1)
	addr2 := pmem.BufAddr(buf2)

	if addr1 == addr2 {
		t.Fatal("buffers start at same address")
	}
	if addr1+uintptr(len(buf1)) > addr2 {
		t.Fatalf("buf1 overlaps buf2: buf1=[%d, %d) buf2=[%d, %d)", addr1, addr1+uintptr(len(buf1)), addr2, addr2+uintptr(len(buf2)))
	}
}

func TestArenaAlloc_OOM(t *testing.T) {
	arena := pmem.NewArena(64)

	_, err := arena.Alloc(128)
	if err == nil {
		t.Fatal("expected OOM error, got nil")
	}
}

func TestArenaAlloc_ExactSize(t *testing.T) {
	arena := pmem.NewArena(64)

	buf, err := arena.Alloc(64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buf) != 64 {
		t.Fatalf("expected len 64, got %d", len(buf))
	}

	// следующая аллокация должна вернуть OOM
	_, err = arena.Alloc(1)
	if err == nil {
		t.Fatal("expected OOM after exact fill")
	}
}

func TestArenaAlloc_WriteAndRead(t *testing.T) {
	arena := pmem.NewArena(1024)

	buf, err := arena.Alloc(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf[0] = 0xDE
	buf[1] = 0xAD
	buf[2] = 0xBE
	buf[3] = 0xEF

	if buf[0] != 0xDE || buf[1] != 0xAD || buf[2] != 0xBE || buf[3] != 0xEF {
		t.Fatal("data mismatch after write")
	}
}

func TestArenaAllocAligned_Alignment(t *testing.T) {
	arena := pmem.NewArena(1024)

	// выделяем 1 байт чтобы сдвинуть offset
	_, err := arena.Alloc(1)
	if err != nil {
		t.Fatalf("alloc 1 byte: %v", err)
	}

	// следующая аллокация с выравниванием 8
	buf, err := arena.AllocAligned(8, 8)
	if err != nil {
		t.Fatalf("aligned alloc: %v", err)
	}

	addr := uintptr(len(buf)) // косвенная проверка через размер
	_ = addr

	// проверяем что адрес выровнен
	if uintptr(pmem.BufAddr(buf))%8 != 0 {
		t.Fatalf("buffer not aligned to 8, addr=%d", pmem.BufAddr(buf))
	}
}

func TestArenaAllocAligned_InvalidAlign(t *testing.T) {
	arena := pmem.NewArena(1024)

	_, err := arena.AllocAligned(8, 3) // 3 не степень двойки
	if err == nil {
		t.Fatal("expected error for non-power-of-2 align")
	}
}

func TestArenaAllocAligned_ZeroAlign(t *testing.T) {
	arena := pmem.NewArena(1024)

	_, err := arena.AllocAligned(8, 0)
	if err == nil {
		t.Fatal("expected error for zero align")
	}
}

func TestArenaReset(t *testing.T) {
	arena := pmem.NewArena(1024)

	buf1, _ := arena.Alloc(256)
	buf1[0] = 0xFF

	if arena.Used() != 256 {
		t.Fatalf("expected used=256, got %d", arena.Used())
	}

	arena.Reset()

	if arena.Used() != 0 {
		t.Fatalf("expected used=0 after reset, got %d", arena.Used())
	}
	if arena.Available() != 1024 {
		t.Fatalf("expected available=1024 after reset, got %d", arena.Available())
	}

	// после reset память должна быть обнулена
	if buf1[0] != 0x00 {
		t.Fatal("memory not zeroed after reset")
	}
}

func TestArenaReset_ReuseAfterReset(t *testing.T) {
	arena := pmem.NewArena(128)

	for i := 0; i < 10; i++ {
		buf, err := arena.Alloc(64)
		if err != nil {
			t.Fatalf("iteration %d alloc failed: %v", i, err)
		}
		buf[0] = byte(i)
		arena.Reset()
	}
}

func TestArenaUsedAndAvailable(t *testing.T) {
	arena := pmem.NewArena(1024)

	if arena.Used() != 0 {
		t.Fatalf("expected used=0, got %d", arena.Used())
	}
	if arena.Available() != 1024 {
		t.Fatalf("expected available=1024, got %d", arena.Available())
	}

	_, _ = arena.Alloc(100)

	// used >= 100 (может быть больше из-за выравнивания)
	if arena.Used() < 100 {
		t.Fatalf("expected used>=100, got %d", arena.Used())
	}
	if arena.Used()+arena.Available() != 1024 {
		t.Fatalf("used+available != total: %d+%d != 1024", arena.Used(), arena.Available())
	}
}

func TestArenaFree_NoOp(t *testing.T) {
	arena := pmem.NewArena(1024)

	buf, _ := arena.Alloc(64)
	usedBefore := arena.Used()

	// Free для arena ничего не делает
	err := arena.Free(buf)
	if err != nil {
		t.Fatalf("Free returned error: %v", err)
	}
	if arena.Used() != usedBefore {
		t.Fatal("Free changed used counter in arena")
	}
}

// ============================================================
// SlabAllocator tests
// ============================================================

func TestSlabAlloc_Basic(t *testing.T) {
	slab := pmem.NewSlab(64, 16)

	buf, err := slab.Alloc(64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buf) != 64 {
		t.Fatalf("expected len 64, got %d", len(buf))
	}
}

func TestSlabAlloc_OOM(t *testing.T) {
	slab := pmem.NewSlab(64, 2)

	_, _ = slab.Alloc(64)
	_, _ = slab.Alloc(64)

	_, err := slab.Alloc(64)
	if err == nil {
		t.Fatal("expected OOM error")
	}
}

func TestSlabAlloc_SizeExceedsBlock(t *testing.T) {
	slab := pmem.NewSlab(64, 16)

	_, err := slab.Alloc(128) // больше blockSize
	if err == nil {
		t.Fatal("expected error for size > blockSize")
	}
}

func TestSlabFree_ReturnToPool(t *testing.T) {
	slab := pmem.NewSlab(64, 2)

	buf1, _ := slab.Alloc(64)
	_, _ = slab.Alloc(64)

	// пул исчерпан
	_, err := slab.Alloc(64)
	if err == nil {
		t.Fatal("expected OOM")
	}

	// возвращаем buf1
	if err := slab.Free(buf1); err != nil {
		t.Fatalf("free failed: %v", err)
	}

	// теперь можно снова выделить
	buf3, err := slab.Alloc(64)
	if err != nil {
		t.Fatalf("alloc after free failed: %v", err)
	}
	if len(buf3) != 64 {
		t.Fatalf("expected len 64, got %d", len(buf3))
	}
}

func TestSlabFree_InvalidBuffer(t *testing.T) {
	slab := pmem.NewSlab(64, 16)

	external := make([]byte, 64)
	err := slab.Free(external)
	if err == nil {
		t.Fatal("expected error for external buffer")
	}
}

func TestSlabFree_ZeroesMemory(t *testing.T) {
	slab := pmem.NewSlab(64, 4)

	buf, _ := slab.Alloc(64)
	buf[0] = 0xFF
	buf[63] = 0xAB

	_ = slab.Free(buf)

	// после free память должна быть обнулена
	if buf[0] != 0x00 || buf[63] != 0x00 {
		t.Fatal("memory not zeroed after free")
	}
}

func TestSlabReset(t *testing.T) {
	slab := pmem.NewSlab(64, 4)

	_, _ = slab.Alloc(64)
	_, _ = slab.Alloc(64)
	_, _ = slab.Alloc(64)
	_, _ = slab.Alloc(64)

	if slab.Available() != 0 {
		t.Fatal("expected available=0 after exhausting slab")
	}

	slab.Reset()

	if slab.Available() != 64*4 {
		t.Fatalf("expected available=%d after reset, got %d", 64*4, slab.Available())
	}
}

func TestSlabUsedAndAvailable(t *testing.T) {
	slab := pmem.NewSlab(64, 4)

	if slab.Used() != 0 {
		t.Fatalf("expected used=0, got %d", slab.Used())
	}

	_, _ = slab.Alloc(64)
	if slab.Used() != 64 {
		t.Fatalf("expected used=64, got %d", slab.Used())
	}

	if slab.Used()+slab.Available() != 64*4 {
		t.Fatalf("used+available != total")
	}
}

// ============================================================
// PoolAllocator tests
// ============================================================

func TestPoolAlloc_SelectsSmallestFit(t *testing.T) {
	pool := pmem.NewPool([]int{64, 128, 256, 512}, 8)

	buf, err := pool.Alloc(100) // должен выбрать slab 128
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buf) != 128 {
		t.Fatalf("expected block size 128, got %d", len(buf))
	}
}

func TestPoolAlloc_ExactSize(t *testing.T) {
	pool := pmem.NewPool([]int{64, 128, 256}, 8)

	buf, err := pool.Alloc(64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buf) != 64 {
		t.Fatalf("expected block size 64, got %d", len(buf))
	}
}

func TestPoolAlloc_TooLarge(t *testing.T) {
	pool := pmem.NewPool([]int{64, 128, 256}, 8)

	_, err := pool.Alloc(512) // больше максимального класса
	if err == nil {
		t.Fatal("expected error for oversized alloc")
	}
}

func TestPoolFree_ReturnsToCorrectSlab(t *testing.T) {
	pool := pmem.NewPool([]int{64, 128}, 2)

	// исчерпываем slab 64
	b1, _ := pool.Alloc(64)
	b2, _ := pool.Alloc(64)
	_, err := pool.Alloc(64)
	if err == nil {
		t.Fatal("expected OOM for slab 64")
	}

	_ = pool.Free(b1)
	_ = pool.Free(b2)

	// теперь снова доступны
	_, err = pool.Alloc(64)
	if err != nil {
		t.Fatalf("alloc after free failed: %v", err)
	}
}

func TestPoolReset(t *testing.T) {
	pool := pmem.NewPool([]int{64, 128, 256}, 2)

	_, _ = pool.Alloc(64)
	_, _ = pool.Alloc(64)
	_, _ = pool.Alloc(128)
	_, _ = pool.Alloc(128)

	pool.Reset()

	// после reset всё доступно
	_, err := pool.Alloc(64)
	if err != nil {
		t.Fatalf("alloc after reset failed: %v", err)
	}
}

func TestPoolUsedAndAvailable(t *testing.T) {
	pool := pmem.NewPool([]int{64, 128}, 4)

	totalBefore := pool.Available()

	_, _ = pool.Alloc(64)
	_, _ = pool.Alloc(128)

	if pool.Used()+pool.Available() != totalBefore {
		t.Fatalf("used+available changed: %d+%d != %d", pool.Used(), pool.Available(), totalBefore)
	}
}
