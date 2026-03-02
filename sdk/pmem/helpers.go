package pmem

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

// ============================================================
// Helpers
// ============================================================

const defaultAlign = 8

func alignUp(offset, align int) int {
	return (offset + align - 1) &^ (align - 1)
}

func mmapAlloc(size int) ([]byte, error) {
	data, err := unix.Mmap(
		-1, 0, size,
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_ANON|unix.MAP_PRIVATE,
	)
	if err != nil {
		return nil, fmt.Errorf("mmap failed: %w", err)
	}
	return data, nil
}

// UnsafeInt64 возвращает *int64 указывающий прямо в буфер — zero copy
func UnsafeInt64(buf []byte) *int64 {
	return (*int64)(unsafe.Pointer(&buf[0]))
}

// UnsafeUint64 возвращает *uint64 указывающий прямо в буфер — zero copy
func UnsafeUint64(buf []byte) *uint64 {
	return (*uint64)(unsafe.Pointer(&buf[0]))
}
