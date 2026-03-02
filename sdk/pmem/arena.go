package pmem

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// ============================================================
// ArenaAllocator
// Линейный аллокатор. Выделяет последовательно, освобождает всё сразу.
// Идеален для lifetime одного запроса / транзакции.
// ============================================================

type ArenaAllocator struct {
	data   []byte
	offset int
	mmaped bool
}

func NewArena(size int) *ArenaAllocator {
	return &ArenaAllocator{
		data: make([]byte, size),
	}
}

func NewMmapArena(size int) (*ArenaAllocator, error) {
	data, err := mmapAlloc(size)
	if err != nil {
		return nil, err
	}
	return &ArenaAllocator{data: data, mmaped: true}, nil
}

func (a *ArenaAllocator) Alloc(size int) ([]byte, error) {
	return a.AllocAligned(size, defaultAlign)
}

func (a *ArenaAllocator) AllocAligned(size, align int) ([]byte, error) {
	if align == 0 || align&(align-1) != 0 {
		return nil, fmt.Errorf("align must be power of 2, got %d", align)
	}
	aligned := alignUp(a.offset, align)
	if aligned+size > len(a.data) {
		return nil, fmt.Errorf("arena oom: need %d, available %d", size, len(a.data)-aligned)
	}
	buf := a.data[aligned : aligned+size]
	a.offset = aligned + size
	return buf, nil
}

// Free — arena не поддерживает индивидуальное освобождение
func (a *ArenaAllocator) Free(_ []byte) error { return nil }

func (a *ArenaAllocator) Reset() {
	clear(a.data[:a.offset])
	a.offset = 0
}

func (a *ArenaAllocator) Used() int      { return a.offset }
func (a *ArenaAllocator) Available() int { return len(a.data) - a.offset }

func (a *ArenaAllocator) Close() error {
	if a.mmaped {
		return unix.Munmap(a.data)
	}
	return nil
}
