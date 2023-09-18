package util

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"

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
func RandInt16(max int16) int16 {
	return int16(RandUint64(uint64(max)))
}

// RandUint16 returns a random uint16 in between 0 and max.
// If max set to zero or negative, the max will set to MaxUint16.
func RandUint16(max uint32) uint16 {
	return uint16(RandUint64(uint64(max)))
}

// RandInt32 returns a random int32 in between 0 and max.
// If max set to zero or negative, the max will set to MaxInt32.
func RandInt32(max int32) int32 {
	return int32(RandUint64(uint64(max)))
}

// RandUint32 returns a random uint32 in between 0 and max.
// If max set to zero or negative, the max will set to MaxUint32.
func RandUint32(max uint32) uint32 {
	return uint32(RandUint64(uint64(max)))
}

// RandInt64 returns a random int64 in between 0 and max.
// If max set to zero or negative, the max will set to MaxInt64.
func RandInt64(max int64) int64 {
	return int64(RandUint64(uint64(max)))
}

// RandUint64 returns a random uint64 in between 0 and max.
// If max set to zero or negative, the max will set to MaxUint64.
func RandUint64(max uint64) uint64 {
	if max <= 0 {
		max = MaxUint64
	}

	bigMax := &big.Int{}
	bigMax.SetUint64(max)
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
func OS2IP(os []byte) *big.Int {
	return new(big.Int).SetBytes(os)
}

// I2OSP converts a nonnegative integer to an octet string of a specified length.
// https://datatracker.ietf.org/doc/html/rfc8017#section-4.1
func IS2OP(x *big.Int, xLen int) []byte {
	if x.Sign() == -1 {
		return nil
	}
	buf := make([]byte, xLen)
	return x.FillBytes(buf)
}

// ExitOnErr exit the software immediately if an error happens.
// Panics are not safe because panics print a stack trace,
// which may not be relevant to the error at all.
func ExitOnErr(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

const changeFactor = float64(1000000000)

// CoinToChange converts a coin amount to its corresponding change value.
// Example: CoinToChange(2.75) returns 2750000000.
func CoinToChange(coin float64) int64 {
	return int64(coin * changeFactor)
}

// ChangeToCoin converts a change value to its corresponding coin amount.
// Example: ChangeToCoin(2750000000) returns 2.75.
func ChangeToCoin(change int64) float64 {
	return float64(change) / changeFactor
}

// StringToChange converts a string representation of a coin amount to its corresponding change value.
// It returns the change value as an int64 and an error if the string cannot be parsed.
// Example: StringToChange("2.75") returns 2750000000, nil.
func StringToChange(amount string) (int64, error) {
	coin, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}
	return CoinToChange(coin), nil
}

// ChangeToStringWithTrailingZeros converts a change value to its string representation with trailing zeros.
// The returned string will have up to 9 decimal places.
// Example: ChangeToStringWithTrailingZeros(2750000000) returns "2.750000000".
func ChangeToStringWithTrailingZeros(change int64) string {
	coin := ChangeToCoin(change)
	return strconv.FormatFloat(coin, 'f', 9, 64)
}

// ChangeToString converts a change value to its string representation.
// Example: ChangeToString(2750000000) returns "2.75".
func ChangeToString(change int64) string {
	coin := ChangeToCoin(change)
	return strconv.FormatFloat(coin, 'f', -1, 64)
}
