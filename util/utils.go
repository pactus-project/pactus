package util

import (
	crand "crypto/rand"
	"math/big"
)

const MaxUint16 = ^uint16(0)
const MinUint16 = 0
const MaxInt16 = int16(MaxUint16 >> 1)
const MinInt16 = -MaxInt16 - 1

const MaxUint32 = ^uint32(0)
const MinUint32 = 0
const MaxInt32 = int32(MaxUint32 >> 1)
const MinInt32 = -MaxInt32 - 1

const MaxUint64 = ^uint64(0)
const MinUint64 = 0
const MaxInt64 = int64(MaxUint64 >> 1)
const MinInt64 = -MaxInt64 - 1

/// Max32 returns the biggest of two 32-bits numbers
func Max32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

/// Min32 returns the smallest of two 32-bits numbers
func Min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

/// Min64 returns the smallest of two 64-bits numbers
func Min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

/// Max64 returns the biggest of two 64-bits numbers
func Max64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

/// RandInt32 returns a random int16 in between 0 and max
/// If max set to zero or negative, the max will set to MaxInt16
func RandInt16(max int16) int16 {
	return int16(RandUint64(uint64(max)))
}

/// RandUint32 returns a random uint32 in between 0 and max
/// If max set to zero or negative, the max will set to MaxUint32
func RandUint16(max uint32) uint16 {
	return uint16(RandUint64(uint64(max)))
}

/// RandInt32 returns a random int32 in between 0 and max
/// If max set to zero or negative, the max will set to MaxInt32
func RandInt32(max int32) int32 {
	return int32(RandUint64(uint64(max)))
}

/// RandUint32 returns a random uint32 in between 0 and max
/// If max set to zero or negative, the max will set to MaxUint32
func RandUint32(max uint32) uint32 {
	return uint32(RandUint64(uint64(max)))
}

/// RandInt64 returns a random int64 in between 0 and max
/// If max set to zero or negative, the max will set to MaxInt64
func RandInt64(max int64) int64 {
	return int64(RandUint64(uint64(max)))
}

/// RandUint64 returns a random uint64 in between 0 and max
/// If max set to zero or negative, the max will set to MaxUint64
func RandUint64(max uint64) uint64 {
	if max <= 0 {
		max = MaxUint64
	}

	bigMax := &big.Int{}
	bigMax.SetUint64(max)
	bigRnd, _ := crand.Int(crand.Reader, bigMax)
	return bigRnd.Uint64()
}

/// SetFlag applies mask to the flags
func SetFlag(flags, mask int) int {
	flags = flags | mask
	return flags
}

/// UnsetFlag removes mask from the flags
func UnsetFlag(flags, mask int) int {
	flags = flags & ^mask
	return flags
}

/// IsFlagSet checks if flags is set or not
func IsFlagSet(flags, mask int) bool {
	return flags&mask == mask
}

// OS2IP converts an octet string to a nonnegative integer.
// OS2IP: https://datatracker.ietf.org/doc/html/rfc8017#section-4.2
func OS2IP(os []byte) *big.Int {
	return new(big.Int).SetBytes(os)
}

/// I2OSP converts a nonnegative integer to an octet string of a specified length.
/// https://datatracker.ietf.org/doc/html/rfc8017#section-4.1
func IS2OP(x *big.Int, xLen int) []byte {
	if x.Sign() == -1 {
		return nil
	}
	buf := make([]byte, xLen)
	return x.FillBytes(buf)
}
