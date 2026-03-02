package pmem

import "unsafe"

// BufAddr возвращает адрес первого байта буфера — только для тестов
func BufAddr(buf []byte) uintptr {
	return uintptr(unsafe.Pointer(&buf[0]))
}
