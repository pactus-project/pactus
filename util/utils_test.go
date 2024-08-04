package util

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert.Equal(t, int32(1), Min(int32(1), 1))
	assert.Equal(t, int32(1), Min(int32(1), 2))
	assert.Equal(t, int32(1), Min(2, int32(1)))
	assert.Equal(t, int32(2), Max(int32(2), 2))
	assert.Equal(t, int32(2), Max(1, int32(2)))
	assert.Equal(t, int32(2), Max(int32(2), 1))

	assert.Equal(t, uint32(1), Min(uint32(1), 1))
	assert.Equal(t, uint32(1), Min(uint32(1), 2))
	assert.Equal(t, uint32(1), Min(2, uint32(1)))
	assert.Equal(t, uint32(2), Max(uint32(2), 2))
	assert.Equal(t, uint32(2), Max(1, uint32(2)))
	assert.Equal(t, uint32(2), Max(uint32(2), 1))

	assert.Equal(t, MaxUint32, uint32(0xffffffff))
	assert.Equal(t, MaxUint64, uint64(0xffffffffffffffff))
	assert.Equal(t, MaxInt32, int32(0x7fffffff))
	assert.Equal(t, MaxInt64, int64(0x7fffffffffffffff))
	assert.Equal(t, MaxInt64, Max(MaxInt64, 1))
	assert.Equal(t, MaxInt64, Max(MinInt64, MaxInt64))
	assert.Equal(t, int64(1), Min(MaxInt64, 1))
	assert.Equal(t, MinInt64, Min(MinInt64, MaxInt64))
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

	assert.Equal(t, []byte{0, 0}, I2OSP(big.NewInt(int64(0)), 2))
	assert.Equal(t, []byte{0, 1}, I2OSP(big.NewInt(int64(1)), 2))
	assert.Equal(t, []byte{0, 255}, I2OSP(big.NewInt(int64(255)), 2))
	assert.Equal(t, []byte{1, 0}, I2OSP(big.NewInt(int64(256)), 2))
	assert.Equal(t, []byte{255, 255}, I2OSP(big.NewInt(int64(65535)), 2))
}

func TestIS2OP(t *testing.T) {
	assert.Equal(t, int64(0), OS2IP([]byte{0, 0}).Int64())
	assert.Equal(t, int64(1), OS2IP([]byte{0, 1}).Int64())
	assert.Equal(t, int64(255), OS2IP([]byte{0, 255}).Int64())
	assert.Equal(t, int64(256), OS2IP([]byte{1, 0}).Int64())
	assert.Equal(t, int64(65535), OS2IP([]byte{255, 255}).Int64())
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

func TestFormatBytesToHumanReadable(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected string
	}{
		{1048576, "1.00 MB"},
		{3145728, "3.00 MB"},
		{1024, "1.00 KB"},
		{512, "512.00 Bytes"},
		{1_073_741_824, "1.00 GB"},
		{1_099_511_627_776, "1.00 TB"},
	}

	for _, test := range tests {
		result := FormatBytesToHumanReadable(test.bytes)
		if result != test.expected {
			t.Errorf("FormatBytesToHumanReadable(%d) returned %s, expected %s", test.bytes, result, test.expected)
		}
	}
}
