package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	assert.Equal(t, Min(1, 1), 1)
	assert.Equal(t, Min(1, 2), 1)
	assert.Equal(t, Min(2, 1), 1)
	assert.Equal(t, Max(2, 2), 2)
	assert.Equal(t, Max(1, 2), 2)
	assert.Equal(t, Max(2, 1), 2)
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
}

func TestRandomPeerID(t *testing.T) {
	id := RandomPeerID()
	assert.NoError(t, id.Validate())
}

func TestRandInt(t *testing.T) {
	rnd1 := RandInt(4)
	assert.GreaterOrEqual(t, rnd1, 0)
	assert.LessOrEqual(t, rnd1, 4)

	rnd2 := RandInt(0)
	assert.Equal(t, rnd2, 0)

	rnd3 := RandInt(-1)
	assert.Equal(t, rnd3, 0)

	rnd4 := RandInt(1000)
	rnd5 := RandInt(1000)
	assert.NotEqual(t, rnd4, rnd5)
}

func TestRandInt64(t *testing.T) {
	rnd1 := RandInt64(4)
	assert.GreaterOrEqual(t, rnd1, int64(0))
	assert.LessOrEqual(t, rnd1, int64(4))

	rnd2 := RandInt64(0)
	assert.Equal(t, rnd2, int64(0))

	rnd3 := RandInt64(-1)
	assert.Equal(t, rnd3, int64(0))

	rnd4 := RandInt64(1000)
	rnd5 := RandInt64(1000)
	assert.NotEqual(t, rnd4, rnd5)
}
