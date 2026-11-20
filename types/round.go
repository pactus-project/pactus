package types

import (
	"math"

	"github.com/pactus-project/pactus/util"
)

type Round int16

// RoundFromBytesLE creates a Round from a byte slice in little-endian format.
func RoundFromBytesLE(data []byte) Round {
	return Round(util.BytesToInt16LE(data))
}

// SafeIncrease returns a new Round that is the result of adding count to h,
// but it returns the maximum round if the result would overflow.
func (h Round) SafeIncrease(count uint32) Round {
	sum := int64(h) + int64(count)
	if sum > math.MaxInt16 {
		return Round(math.MaxInt16)
	}

	return Round(sum)
}

// SafeDecrease returns the result of subtracting count from h,
// but it returns 0 if the result would be negative.
func (h Round) SafeDecrease(count uint32) Round {
	if int64(h) < int64(count) {
		return 0
	}

	return Round(int64(h) - int64(count))
}

// SafeSub returns the result of subtracting other from h,
// but it returns 0 if the result would be negative.
func (h Round) SafeSub(other Round) uint32 {
	if h < other {
		return 0
	}

	return uint32(int32(h) - int32(other))
}

// BytesLE encodes the Round as a byte slice in little-endian format.
func (h Round) BytesLE() []byte {
	return util.Int16ToBytesLE(int16(h))
}
