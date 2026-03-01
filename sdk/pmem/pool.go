package pmem

import "fmt"

// ============================================================
// PoolAllocator
// Объединяет несколько SlabAllocator с разными размерами блоков.
// При Alloc выбирает минимально подходящий slab — general purpose аллокатор
// без GC для объектов переменного размера.
// Аналог tcmalloc / jemalloc size classes.
// ============================================================

type slabClass struct {
	slab      *SlabAllocator
	blockSize int
}

type PoolAllocator struct {
	classes []slabClass // отсортированы по blockSize
}

// NewPool создаёт pool с заданными size classes.
// Пример: NewPool([]int{64, 128, 256, 512, 1024, 4096}, 1024)
func NewPool(blockSizes []int, blocksPerClass int) *PoolAllocator {
	classes := make([]slabClass, len(blockSizes))
	for i, size := range blockSizes {
		classes[i] = slabClass{
			slab:      NewSlab(size, blocksPerClass),
			blockSize: size,
		}
	}
	return &PoolAllocator{classes: classes}
}

func (p *PoolAllocator) Alloc(size int) ([]byte, error) {
	return p.AllocAligned(size, defaultAlign)
}

func (p *PoolAllocator) AllocAligned(size, align int) ([]byte, error) {
	for _, class := range p.classes {
		if class.blockSize >= size {
			return class.slab.AllocAligned(size, align)
		}
	}
	return nil, fmt.Errorf("pool: no slab class for size %d", size)
}

func (p *PoolAllocator) Free(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	// пробуем найти slab которому принадлежит буфер
	for _, class := range p.classes {
		err := class.slab.Free(buf)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("pool: buffer does not belong to any slab class")
}

func (p *PoolAllocator) Reset() {
	for _, class := range p.classes {
		class.slab.Reset()
	}
}

func (p *PoolAllocator) Used() int {
	total := 0
	for _, class := range p.classes {
		total += class.slab.Used()
	}
	return total
}

func (p *PoolAllocator) Available() int {
	total := 0
	for _, class := range p.classes {
		total += class.slab.Available()
	}
	return total
}

func (p *PoolAllocator) Close() error {
	for _, class := range p.classes {
		if err := class.slab.Close(); err != nil {
			return err
		}
	}
	return nil
}
