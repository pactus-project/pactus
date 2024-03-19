package util

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert.Equal(t, Min(int32(1), 1), int32(1))
	assert.Equal(t, Min(int32(1), 2), int32(1))
	assert.Equal(t, Min(2, int32(1)), int32(1))
	assert.Equal(t, Max(int32(2), 2), int32(2))
	assert.Equal(t, Max(1, int32(2)), int32(2))
	assert.Equal(t, Max(int32(2), 1), int32(2))

	assert.Equal(t, Min(uint32(1), 1), uint32(1))
	assert.Equal(t, Min(uint32(1), 2), uint32(1))
	assert.Equal(t, Min(2, uint32(1)), uint32(1))
	assert.Equal(t, Max(uint32(2), 2), uint32(2))
	assert.Equal(t, Max(1, uint32(2)), uint32(2))
	assert.Equal(t, Max(uint32(2), 1), uint32(2))

	assert.Equal(t, MaxUint32, uint32(0xffffffff))
	assert.Equal(t, MaxUint64, uint64(0xffffffffffffffff))
	assert.Equal(t, MaxInt32, int32(0x7fffffff))
	assert.Equal(t, MaxInt64, int64(0x7fffffffffffffff))
	assert.Equal(t, Max(MaxInt64, 1), MaxInt64)
	assert.Equal(t, Max(MinInt64, MaxInt64), MaxInt64)
	assert.Equal(t, Min(MaxInt64, 1), int64(1))
	assert.Equal(t, Min(MinInt64, MaxInt64), MinInt64)
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
	assert.Nil(t, I2OSP(big.NewInt(int64(-1)), 2))

	assert.Equal(t, I2OSP(big.NewInt(int64(0)), 2), []byte{0, 0})
	assert.Equal(t, I2OSP(big.NewInt(int64(1)), 2), []byte{0, 1})
	assert.Equal(t, I2OSP(big.NewInt(int64(255)), 2), []byte{0, 255})
	assert.Equal(t, I2OSP(big.NewInt(int64(256)), 2), []byte{1, 0})
	assert.Equal(t, I2OSP(big.NewInt(int64(65535)), 2), []byte{255, 255})
}

func TestIS2OP(t *testing.T) {
	assert.Equal(t, OS2IP([]byte{0, 0}).Int64(), int64(0))
	assert.Equal(t, OS2IP([]byte{0, 1}).Int64(), int64(1))
	assert.Equal(t, OS2IP([]byte{0, 255}).Int64(), int64(255))
	assert.Equal(t, OS2IP([]byte{1, 0}).Int64(), int64(256))
	assert.Equal(t, OS2IP([]byte{255, 255}).Int64(), int64(65535))
}

func TestLogScale(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{1, 1},
		{2, 2},
		{3, 4},
		{7, 8},
		{8, 8},
	}

	for _, testCase := range testCases {
		result := LogScale(testCase.input)
		assert.Equal(t, testCase.expected, result, "LogScale(%d) failed", testCase.input)
	}
}
