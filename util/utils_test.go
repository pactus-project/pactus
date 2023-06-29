package util

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert.Equal(t, Min32(1, 1), int32(1))
	assert.Equal(t, Min32(1, 2), int32(1))
	assert.Equal(t, Min32(2, 1), int32(1))
	assert.Equal(t, Max32(2, 2), int32(2))
	assert.Equal(t, Max32(1, 2), int32(2))
	assert.Equal(t, Max32(2, 1), int32(2))

	assert.Equal(t, MinU32(1, 1), uint32(1))
	assert.Equal(t, MinU32(1, 2), uint32(1))
	assert.Equal(t, MinU32(2, 1), uint32(1))
	assert.Equal(t, MaxU32(2, 2), uint32(2))
	assert.Equal(t, MaxU32(1, 2), uint32(2))
	assert.Equal(t, MaxU32(2, 1), uint32(2))

	assert.Equal(t, MaxUint32, uint32(0xffffffff))
	assert.Equal(t, MaxUint64, uint64(0xffffffffffffffff))
	assert.Equal(t, MaxInt32, int32(0x7fffffff))
	assert.Equal(t, MaxInt64, int64(0x7fffffffffffffff))
	assert.Equal(t, Max64(MaxInt64, 1), MaxInt64)
	assert.Equal(t, Max64(MinInt64, MaxInt64), MaxInt64)
	assert.Equal(t, Min64(MaxInt64, 1), int64(1))
	assert.Equal(t, Min64(MinInt64, MaxInt64), MinInt64)
}

func TestSetFlags(t *testing.T) {
	flags := 0
	flags = SetFlag(flags, 0x2)
	flags = SetFlag(flags, 0x8)
	assert.Equal(t, flags, 0xa)
	assert.True(t, IsFlagSet(flags, 0x2))
	assert.False(t, IsFlagSet(flags, 0x4))
	flags = UnsetFlag(flags, 0x2)
	assert.False(t, IsFlagSet(flags, 0x2))
	assert.Equal(t, flags, 0x8)
}

func TestRandUint16(t *testing.T) {
	rnd := RandUint16(4)
	assert.GreaterOrEqual(t, rnd, uint16(0))
	assert.LessOrEqual(t, rnd, uint16(4))
}

func TestRandInt16(t *testing.T) {
	rnd := RandInt16(4)
	assert.GreaterOrEqual(t, rnd, int16(0))
	assert.LessOrEqual(t, rnd, int16(4))
}

func TestRandUint32(t *testing.T) {
	rnd := RandUint32(4)
	assert.GreaterOrEqual(t, rnd, uint32(0))
	assert.LessOrEqual(t, rnd, uint32(4))
}

func TestRandInt32(t *testing.T) {
	rnd := RandInt32(4)
	assert.GreaterOrEqual(t, rnd, int32(0))
	assert.LessOrEqual(t, rnd, int32(4))
}

func TestRandInt64(t *testing.T) {
	rnd := RandInt64(4)
	assert.GreaterOrEqual(t, rnd, int64(0))
	assert.LessOrEqual(t, rnd, int64(4))
}

func TestRandUint64(t *testing.T) {
	rnd1 := RandUint64(4)
	assert.GreaterOrEqual(t, rnd1, uint64(0))
	assert.LessOrEqual(t, rnd1, uint64(4))

	rnd2 := RandUint64(0)
	assert.NotZero(t, rnd2)
}

func TestI2OSP(t *testing.T) {
	assert.Nil(t, IS2OP(big.NewInt(int64(-1)), 2))

	assert.Equal(t, IS2OP(big.NewInt(int64(0)), 2), []byte{0, 0})
	assert.Equal(t, IS2OP(big.NewInt(int64(1)), 2), []byte{0, 1})
	assert.Equal(t, IS2OP(big.NewInt(int64(255)), 2), []byte{0, 255})
	assert.Equal(t, IS2OP(big.NewInt(int64(256)), 2), []byte{1, 0})
	assert.Equal(t, IS2OP(big.NewInt(int64(65535)), 2), []byte{255, 255})
}

func TestIS2OP(t *testing.T) {
	assert.Equal(t, OS2IP([]byte{0, 0}).Int64(), int64(0))
	assert.Equal(t, OS2IP([]byte{0, 1}).Int64(), int64(1))
	assert.Equal(t, OS2IP([]byte{0, 255}).Int64(), int64(255))
	assert.Equal(t, OS2IP([]byte{1, 0}).Int64(), int64(256))
	assert.Equal(t, OS2IP([]byte{255, 255}).Int64(), int64(65535))
}

func TestCoinToChangeConversion(t *testing.T) {
	tests := []struct {
		amount  string
		coin    float64
		change  int64
		str1    string
		str2    string
		parsErr error
	}{
		{"0", 0, 0, "0.000000000", "0", nil},
		{"1", 1, 1000000000, "1.000000000", "1", nil},
		{"123.123", 123.123, 123123000000, "123.123000000", "123.123", nil},
		{"123.0123", 123.0123, 123012300000, "123.012300000", "123.0123", nil},
		{"123.01230", 123.0123, 123012300000, "123.012300000", "123.0123", nil},
		{"123.000123", 123.000123, 123000123000, "123.000123000", "123.000123", nil},
		{"123.000000123", 123.000000123, 123000000123, "123.000000123", "123.000000123", nil},
		{"-123.000000123", -123.000000123, -123000000123, "-123.000000123", "-123.000000123", nil},
		{"0123.000000123", 123.000000123, 123000000123, "123.000000123", "123.000000123", nil},
		{"+123.000000123", 123.000000123, 123000000123, "123.000000123", "123.000000123", nil},
		{"123.0000001234", 123.000000123, 123000000123, "123.000000123", "123.000000123", nil},
		{"1coin", 0, 0, "0.000000000", "0", strconv.ErrSyntax},
	}
	for _, test := range tests {
		change, err := StringToChange(test.amount)
		if test.parsErr == nil {
			assert.NoError(t, err)
			assert.Equal(t, change, test.change)
			assert.Equal(t, ChangeToCoin(change), test.coin)
			assert.Equal(t, ChangeToStringWithTrailingZeros(change), test.str1)
			assert.Equal(t, ChangeToString(change), test.str2)
		} else {
			assert.ErrorIs(t, err, test.parsErr)
		}
	}
}
