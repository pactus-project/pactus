package util

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
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

func CompressSlice(s []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(s); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DecompressSlice(s []byte) ([]byte, error) {
	b := bytes.NewBuffer(s)
	var r io.Reader
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}

	var res bytes.Buffer
	if _, err = res.ReadFrom(r); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}
