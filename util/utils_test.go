package util

import (
	"fmt"
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

func TestRandomPeerID(t *testing.T) {
	id := RandomPeerID()
	assert.NoError(t, id.Validate())
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

// TestRandomUint64 exercises the randomness of the random number generator on
// the system by ensuring the probability of the generated numbers.  If the RNG
// is evenly distributed as a proper cryptographic RNG should be, there really
// should only be 1 number < 2^56 in 2^8 tries for a 64-bit number.  However,
// use a higher number of 5 to really ensure the test doesn't fail unless the
// RNG is just horrendous.
func TestRandomUint64(t *testing.T) {
	tries := 1 << 8              // 2^8
	watermark := uint64(1 << 56) // 2^56
	maxHits := 5
	badRNG := "The random number generator on this system is clearly " +
		"terrible since we got %d values less than %d in %d runs " +
		"when only %d was expected"

	numHits := 0
	for i := 0; i < tries; i++ {
		nonce := RandUint64(MaxUint64)
		if nonce < watermark {
			numHits++
		}
		if numHits > maxHits {
			str := fmt.Sprintf(badRNG, numHits, watermark, tries, maxHits)
			t.Errorf("Random Uint64 iteration %d failed - %v %v", i,
				str, numHits)
			return
		}
	}
}
