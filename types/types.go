package types

import "github.com/pactus-project/pactus/util"

type Height uint32

// HeightFromSlice creates a Height from a byte slice in little-endian format.
func HeightFromSlice(data []byte) Height {
	return Height(util.SliceToUint32(data))
}

// SafeIncrease returns a new Height that is the result of adding count to h.
func (h Height) SafeIncrease(count uint32) Height {
	return Height(uint32(h) + count)
}

// SafeAdd returns the result of subtracting other from h,
// but it returns 0 if the result would be negative.
func (h Height) SafeDecrease(count uint32) Height {
	if uint32(h) < count {
		return 0
	}

	return Height(uint32(h) - count)
}

// SafeSub returns the result of subtracting other from h,
// but it returns 0 if the result would be negative.
func (h Height) SafeSub(other Height) uint32 {
	if h < other {
		return 0
	}

	return uint32(h - other)
}

// EncodeAsSlice encodes the height as a byte slice in little-endian format.
func (h Height) EncodeAsSlice() []byte {
	return util.Uint32ToSlice(uint32(h))
}
