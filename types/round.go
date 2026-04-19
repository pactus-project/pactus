package types

import "github.com/pactus-project/pactus/util"

type Round int16

// RoundFromSlice creates a Round from a byte slice in little-endian format.
func RoundFromSlice(data []byte) Round {
	return Round(util.SliceToInt16(data))
}

// SafeIncrease returns a new Round that is the result of adding count to h.
func (h Round) SafeIncrease(count uint32) Round {
	return Round(uint32(h) + count)
}

// SafeDecrease returns the result of subtracting other from h,
// but it returns 0 if the result would be negative.
func (h Round) SafeDecrease(count uint32) Round {
	if uint32(h) < count {
		return 0
	}

	return Round(uint32(h) - count)
}

// SafeSub returns the result of subtracting other from h,
// but it returns 0 if the result would be negative.
func (h Round) SafeSub(other Round) uint32 {
	if h < other {
		return 0
	}

	return uint32(h - other)
}

// EncodeAsSlice encodes the Round as a byte slice in little-endian format.
func (h Round) EncodeAsSlice() []byte {
	return util.Int16ToSlice(int16(h))
}
