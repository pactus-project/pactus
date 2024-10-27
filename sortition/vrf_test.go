package sortition_test

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestVRF(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandBLSKeyPair()
	valKey := bls.NewValidatorKey(prv)
	for i := 0; i < 100; i++ {
		seed := ts.RandSeed()
		t.Logf("seed is: %x \n", seed)

		maxSize := uint64(1 * 1e6)
		index, proof := sortition.Evaluate(seed, valKey.PrivateKey(), maxSize)

		assert.LessOrEqual(t, index, maxSize)

		index2, result := sortition.Verify(seed, pub, proof, maxSize)

		assert.True(t, result)
		assert.Equal(t, index, index2)
	}
}

// TestRandomUint64 exercises the randomness of the random number generator on
// the system by ensuring the probability of the generated numbers.  If the RNG
// is evenly distributed as a proper cryptographic RNG should be, there really
// should only be 1 number < 2^56 in 2^8 tries for a 64-bit number.  However,
// use a higher number of 5 to really ensure the test doesn't fail unless the
// RNG is just horrendous.
func TestRandomUint64(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tries := 1 << 8              // 2^8
	watermark := uint64(1 << 56) // 2^56
	maxHits := 5
	badRNG := "The random number generator on this system is clearly " +
		"terrible since we got %d values less than %d in %d runs " +
		"when only %d was expected"
	_, pv := ts.RandBLSKeyPair()

	valKey := bls.NewValidatorKey(pv)

	numHits := 0
	for i := 0; i < tries; i++ {
		seed := ts.RandSeed()

		nonce, _ := sortition.Evaluate(seed, valKey.PrivateKey(), util.MaxUint64)
		if nonce < watermark {
			numHits++
		}
		if numHits > maxHits {
			str := fmt.Sprintf(badRNG, numHits, watermark, tries, maxHits)
			t.Errorf("Random Uint64 iteration %d failed - %v", i, str)

			return
		}
	}
}

func TestGetIndex(t *testing.T) {
	// The expected values
	//
	// Total: 1,000,000
	//
	// proof: 0x1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03
	// proofH: 0xa7b8166584387f4ea76f9caa0969bd6b0bb8df4c3bb8e87f8b6e4dad62bf3359
	//
	// proofH * 1000000 / denominator = 655152.7021258341
	proof1, _ := sortition.ProofFromString(
		"1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03")
	assert.Equal(t, uint64(655152), sortition.GetIndex(proof1, 1*1e6))

	// proof: 45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396
	// proofH: 80212979d1de1ca4ce1258fc0be66a4453b3804e64a5ca8d95f7def2c291c7fe
	//
	// proofH * 1000000 / denominator = 500506.0121928797
	proof2, _ := sortition.ProofFromString(
		"45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396")
	assert.Equal(t, uint64(500506), sortition.GetIndex(proof2, 1*1e6))
}
