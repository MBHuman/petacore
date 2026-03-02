package pmem

// ============================================================
// Allocator interface
// ============================================================

type Allocator interface {
	// Alloc выделяет size байт
	Alloc(size int) ([]byte, error)
	// AllocAligned выделяет size байт с заданным выравниванием
	AllocAligned(size, align int) ([]byte, error)
	// Free освобождает память (не все аллокаторы поддерживают)
	Free(buf []byte) error
	// Reset сбрасывает все выделения
	Reset()
	// Used возвращает количество использованных байт
	Used() int
	// Available возвращает количество свободных байт
	Available() int
	// Close освобождает ресурсы аллокатора
	Close() error
}
