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

func CompressBuffer(s []byte) ([]byte, error) {
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

func DecompressBuffer(s []byte) ([]byte, error) {
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

// Subtracts subtracts slice2 from slice 1 in order
// Examples:
//  [1,2,3,4] - [1,2] = [3,4]
//  [1,2,3,4] - [2,4] = [1,3]
//  [1,2,3,4] - [4,2] = [1,3]
//  [1,2,3,4] - [4,5] = [1,2,3]
func Subtracts(slice1 []int, slice2 []int) []int {
	sub := []int{}
	if slice2 == nil {
		return slice1
	}

	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			sub = append(sub, s1)
		}
	}

	return sub
}

// HasItem checks whether the given slice has a specific item.
func HasItem(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
