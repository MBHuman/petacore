package pmem

import (
	"sync"
)

// ============================================================
// ArenaPool
// Пул аллокаторов Arena для переиспользования между запросами.
// Вместо создания нового аллокатора на каждый запрос,
// берем из пула, используем, сбрасываем и возвращаем в пул.
// ============================================================

type ArenaPool struct {
	pool      sync.Pool
	arenaSize int
}

func NewArenaPool(arenaSize int) *ArenaPool {
	return &ArenaPool{
		arenaSize: arenaSize,
		pool: sync.Pool{
			New: func() any {
				a, err := NewMmapArena(arenaSize)
				if err != nil {
					panic("mmap arena alloc: " + err.Error())
				}
				return a
			},
		},
	}
}

// Get возвращает аллокатор из пула (или создает новый)
func (ap *ArenaPool) Get() *ArenaAllocator {
	arena := ap.pool.Get().(*ArenaAllocator)
	arena.Reset() // Сбрасываем состояние перед использованием
	return arena
}

// Put возвращает аллокатор в пул после использования
func (ap *ArenaPool) Put(arena *ArenaAllocator) {
	if arena == nil {
		return
	}
	// Проверяем, что arena имеет правильный размер
	// Если arena была создана вне пула с другим размером, не кладем её обратно
	if len(arena.data) == ap.arenaSize {
		arena.Reset() // Очищаем перед возвратом в пул
		ap.pool.Put(arena)
	}
}
