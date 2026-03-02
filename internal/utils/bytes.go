package utils

import "unsafe"

// BytesToString converts []byte to string without copying.
// WARNING: The returned string shares memory with the input slice.
// The caller must ensure the slice is not modified while the string is in use.
// This is safe for read-only operations and when the slice lifetime exceeds string usage.
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToBytes converts string to []byte without copying.
// WARNING: The returned slice shares memory with the input string.
// The caller must NOT modify the returned slice.
// This is safe only for read-only operations.
func StringToBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
