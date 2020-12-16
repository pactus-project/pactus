package util

import (
	"encoding/binary"
)

func UIntToSlice(n uint) []byte {
	return UInt64ToSlice(uint64(n))
}
func IntToSlice(n int) []byte {
	return Int64ToSlice(int64(n))
}

func SliceToUInt(bs []byte) uint {
	return uint(SliceToUInt64(bs))
}

func SliceToInt(bs []byte) int {
	return int(SliceToInt64(bs))
}

func UInt64ToSlice(n uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(n))
	return bs
}
func Int64ToSlice(n int64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(n))
	return bs
}

func SliceToUInt64(bs []byte) uint64 {
	n := binary.LittleEndian.Uint64(bs)
	return n
}

func SliceToInt64(bs []byte) int64 {
	n := binary.LittleEndian.Uint64(bs)
	return int64(n)
}
