package util

import (
	"encoding/binary"
)

func UIntToSlice(n uint) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(n))
	return bs
}

func SliceToUInt(bs []byte) uint {
	n := binary.LittleEndian.Uint32(bs)
	return uint(n)
}

func IntToSlice(n int) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(n))
	return bs
}

func SliceToInt(bs []byte) int {
	n := binary.LittleEndian.Uint32(bs)
	return int(n)
}

func SliceToUInt64(bs []byte) uint64 {
	n := binary.LittleEndian.Uint64(bs)
	return n
}

func SliceToInt64(bs []byte) int64 {
	n := binary.LittleEndian.Uint64(bs)
	return int64(n)
}
