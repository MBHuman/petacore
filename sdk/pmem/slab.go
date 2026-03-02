package pmem

import (
	"fmt"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ============================================================
// SlabAllocator
// Выделяет блоки фиксированного размера из заранее выделенного пула.
// O(1) Alloc и Free. Идеален когда размер объектов известен заранее
// (строки фиксированной длины, page descriptors, row headers и т.д.)
// ============================================================

type slabBlock struct {
	offset int // смещение блока внутри data
}

type SlabAllocator struct {
	data      []byte
	blockSize int
	free      []slabBlock // стек свободных блоков
	mmaped    bool
	mu        sync.Mutex
}

func NewSlab(blockSize, blockCount int) *SlabAllocator {
	size := blockSize * blockCount
	data := make([]byte, size)
	return newSlabFromData(data, blockSize, blockCount, false)
}

func NewMmapSlab(blockSize, blockCount int) (*SlabAllocator, error) {
	size := blockSize * blockCount
	data, err := mmapAlloc(size)
	if err != nil {
		return nil, err
	}
	return newSlabFromData(data, blockSize, blockCount, true), nil
}

func newSlabFromData(data []byte, blockSize, blockCount int, mmaped bool) *SlabAllocator {
	free := make([]slabBlock, blockCount)
	for i := 0; i < blockCount; i++ {
		free[i] = slabBlock{offset: i * blockSize}
	}
	return &SlabAllocator{
		data:      data,
		blockSize: blockSize,
		free:      free,
		mmaped:    mmaped,
	}
}

func (s *SlabAllocator) Alloc(size int) ([]byte, error) {
	if size > s.blockSize {
		return nil, fmt.Errorf("slab: requested %d > block size %d", size, s.blockSize)
	}
	return s.AllocAligned(size, defaultAlign)
}

func (s *SlabAllocator) AllocAligned(size, _ int) ([]byte, error) {
	if size > s.blockSize {
		return nil, fmt.Errorf("slab: requested %d > block size %d", size, s.blockSize)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.free) == 0 {
		return nil, fmt.Errorf("slab oom: no free blocks")
	}

	block := s.free[len(s.free)-1]
	s.free = s.free[:len(s.free)-1]

	return s.data[block.offset : block.offset+s.blockSize], nil
}

func (s *SlabAllocator) Free(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}

	// вычисляем offset через unsafe pointer арифметику
	bufPtr := uintptr(unsafe.Pointer(&buf[0]))
	dataPtr := uintptr(unsafe.Pointer(&s.data[0]))

	if bufPtr < dataPtr || bufPtr >= dataPtr+uintptr(len(s.data)) {
		return fmt.Errorf("slab: buffer does not belong to this allocator")
	}

	offset := int(bufPtr - dataPtr)
	if offset%s.blockSize != 0 {
		return fmt.Errorf("slab: invalid buffer offset %d", offset)
	}

	// обнуляем блок перед возвратом
	clear(s.data[offset : offset+s.blockSize])

	s.mu.Lock()
	s.free = append(s.free, slabBlock{offset: offset})
	s.mu.Unlock()

	return nil
}

func (s *SlabAllocator) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.data)
	blockCount := len(s.data) / s.blockSize
	s.free = s.free[:0]
	for i := 0; i < blockCount; i++ {
		s.free = append(s.free, slabBlock{offset: i * s.blockSize})
	}
}

func (s *SlabAllocator) Used() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return (len(s.data)/s.blockSize - len(s.free)) * s.blockSize
}

func (s *SlabAllocator) Available() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.free) * s.blockSize
}

func (s *SlabAllocator) Close() error {
	if s.mmaped {
		return unix.Munmap(s.data)
	}
	return nil
}
