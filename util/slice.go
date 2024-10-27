package util

import (
	"bytes"
	"compress/gzip"
	"crypto/subtle"
	"encoding/binary"
	"math/rand"
)

func Uint16ToSlice(n uint16) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, n)

	return bs
}

func Int16ToSlice(n int16) []byte {
	return Uint16ToSlice(uint16(n))
}

func SliceToUint16(bs []byte) uint16 {
	return binary.LittleEndian.Uint16(bs)
}

func SliceToInt16(bs []byte) int16 {
	return int16(SliceToUint16(bs))
}

func Uint32ToSlice(n uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, n)

	return bs
}

func Int32ToSlice(n int32) []byte {
	return Uint32ToSlice(uint32(n))
}

func SliceToUint32(bs []byte) uint32 {
	return binary.LittleEndian.Uint32(bs)
}

func SliceToInt32(bs []byte) int32 {
	return int32(SliceToUint32(bs))
}

func Uint64ToSlice(n uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, n)

	return bs
}

func Int64ToSlice(n int64) []byte {
	return Uint64ToSlice(uint64(n))
}

func SliceToUint64(bs []byte) uint64 {
	n := binary.LittleEndian.Uint64(bs)

	return n
}

func SliceToInt64(bs []byte) int64 {
	return int64(SliceToUint64(bs))
}

// StringToBytes converts a string to a slice of bytes.
func StringToBytes(s string) []byte {
	return []byte(s)
}

func CompressBuffer(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecompressBuffer(s []byte) ([]byte, error) {
	buf := bytes.NewBuffer(s)
	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	var res bytes.Buffer
	if _, err = res.ReadFrom(reader); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

// Subtracts subtracts slice2 from slice1 in order.
// Examples:
//
//	[1,2,3,4] - [1,2] = [3,4]
//	[1,2,3,4] - [2,4] = [1,3]
//	[1,2,3,4] - [4,2] = [1,3]
//	[1,2,3,4] - [4,5] = [1,2,3]
//
// .
func Subtracts(slice1, slice2 []int32) []int32 {
	sub := []int32{}
	if slice2 == nil {
		return slice1
	}

	for _, num1 := range slice1 {
		found := false
		for _, num2 := range slice2 {
			if num1 == num2 {
				found = true

				break
			}
		}
		if !found {
			sub = append(sub, num1)
		}
	}

	return sub
}

// Contains checks whether the given slice has a specific item.
func Contains[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}

	return false
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal[T comparable](a, b []T) bool {
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

// SafeCmp compares two slices with constant time.
// Note that we are using the subtle.ConstantTimeCompare() function for this
// to help prevent timing attacks.
func SafeCmp(left, right []byte) bool {
	return subtle.ConstantTimeCompare(left, right) == 1
}

// Merge accepts multiple slices and returns a single merged slice.
func Merge[T any](slices ...[]T) []T {
	var totalLength int

	// Calculate the total length of the merged slice.
	for _, slice := range slices {
		totalLength += len(slice)
	}

	// Create a merged slice with the appropriate capacity.
	merged := make([]T, 0, totalLength)

	// Append each input slice to the merged slice.
	for _, slice := range slices {
		merged = append(merged, slice...)
	}

	return merged
}

// Reverse replace the contents of a slice with the same elements but in
// reverse order.
func Reverse[S ~[]E, E any](slice S) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Extend extends the slice to the given length by appending zero-valued elements.
func Extend[T any](slice []T, length int) []T {
	if len(slice) < length {
		pad := make([]T, length-len(slice), length+len(slice))
		slice = append(pad, slice...)
	}

	return slice
}

// IsSubset checks if subSet is a subset of parentSet.
// It returns true if all elements of subSet are in parentSet.
func IsSubset[T comparable](parentSet, subSet []T) bool {
	lastIndex := 0
	for i := 0; i < len(subSet); i++ {
		matchFound := false
		for j := lastIndex; j < len(parentSet); j++ {
			if subSet[i] == parentSet[j] {
				matchFound = true
				lastIndex = j

				break
			}
		}
		if !matchFound {
			return false
		}
	}

	return true
}

// RemoveFirstOccurrenceOf removes the first occurrence of element from slice.
// It returns the modified slice and a boolean indicating whether an element was removed.
func RemoveFirstOccurrenceOf[T comparable](slice []T, element T) ([]T, bool) {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...), true
		}
	}

	return slice, false
}

func Trim[T any](slice []T, newLength int) []T {
	if newLength <= len(slice) {
		return slice[:newLength]
	}

	return slice
}

// Shuffle shuffles a slice of any type.
func Shuffle[T any](slice []T) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
