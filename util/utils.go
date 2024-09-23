package util

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/bits"

	"golang.org/x/exp/constraints"
)

const (
	MaxUint16 = ^uint16(0)
	MinUint16 = 0
	MaxInt16  = int16(MaxUint16 >> 1)
	MinInt16  = -MaxInt16 - 1
)

const (
	MaxUint32 = ^uint32(0)
	MinUint32 = 0
	MaxInt32  = int32(MaxUint32 >> 1)
	MinInt32  = -MaxInt32 - 1
)

const (
	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int64(MaxUint64 >> 1)
	MinInt64  = -MaxInt64 - 1
)

// Max returns the biggest of two integer numbers.
func Max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// Min returns the smallest of two integer numbers.
func Min[T constraints.Integer](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// RandInt16 returns a random int16 in between 0 and max.
// If max set to zero or negative, the max will set to MaxInt16.
func RandInt16(maxVal int16) int16 {
	return int16(RandUint64(uint64(maxVal)))
}

// RandUint16 returns a random uint16 in between 0 and max.
// If max set to zero or negative, the max will set to MaxUint16.
func RandUint16(maxVal uint32) uint16 {
	return uint16(RandUint64(uint64(maxVal)))
}

// RandInt32 returns a random int32 in between 0 and max.
// If max set to zero or negative, the max will set to MaxInt32.
func RandInt32(maxVal int32) int32 {
	return int32(RandUint64(uint64(maxVal)))
}

// RandUint32 returns a random uint32 in between 0 and max.
// If max set to zero or negative, the max will set to MaxUint32.
func RandUint32(maxVal uint32) uint32 {
	return uint32(RandUint64(uint64(maxVal)))
}

// RandInt64 returns a random int64 in between 0 and max.
// If max set to zero or negative, the max will set to MaxInt64.
func RandInt64(maxVal int64) int64 {
	return int64(RandUint64(uint64(maxVal)))
}

// RandUint64 returns a random uint64 in between 0 and max.
// If max set to zero or negative, the max will set to MaxUint64.
func RandUint64(maxVal uint64) uint64 {
	if maxVal <= 0 {
		maxVal = MaxUint64
	}

	bigMax := &big.Int{}
	bigMax.SetUint64(maxVal)
	bigRnd, _ := crand.Int(crand.Reader, bigMax)

	return bigRnd.Uint64()
}

// SetFlag applies mask to the flags.
func SetFlag[T constraints.Integer](flags, mask T) T {
	return flags | mask
}

// UnsetFlag removes mask from the flags.
func UnsetFlag[T constraints.Integer](flags, mask T) T {
	return flags & ^mask
}

// IsFlagSet checks if the mask is set for the given flags.
func IsFlagSet[T constraints.Integer](flags, mask T) bool {
	return flags&mask == mask
}

// OS2IP converts an octet string to a nonnegative integer.
// OS2IP: https://datatracker.ietf.org/doc/html/rfc8017#section-4.2
func OS2IP(x []byte) *big.Int {
	return new(big.Int).SetBytes(x)
}

// I2OSP converts a nonnegative integer to an octet string of a specified length.
// https://datatracker.ietf.org/doc/html/rfc8017#section-4.1
func I2OSP(x *big.Int, xLen int) []byte {
	if x.Sign() == -1 {
		return nil
	}
	buf := make([]byte, xLen)

	return x.FillBytes(buf)
}

// LogScale computes 2^⌈log₂(val)⌉, where ⌈x⌉ represents the ceiling of x.
// For more information, refer to: https://en.wikipedia.org/wiki/Logarithmic_scale
func LogScale(val int) int {
	bitlen := bits.Len(uint(val - 1))

	return 1 << bitlen
}

func FormatBytesToHumanReadable(bytes uint64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)
	unit := "Bytes"
	value := float64(bytes)

	switch {
	case bytes >= TB:
		unit = "TB"
		value /= TB
	case bytes >= GB:
		unit = "GB"
		value /= GB
	case bytes >= MB:
		unit = "MB"
		value /= MB
	case bytes >= KB:
		unit = "KB"
		value /= KB
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}
