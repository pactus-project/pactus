package sortition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

func TestVRF(t *testing.T) {
	pk, pv := bls.GenerateTestKeyPair()
	signer := crypto.NewSigner(pv)
	for i := 0; i < 100; i++ {
		seed := GenerateRandomSeed()
		fmt.Printf("seed is: %x \n", seed)

		max := uint64(1 * 1e6)
		index, proof := evaluate(seed, signer, max)

		assert.LessOrEqual(t, index, max)

		index2, result := verify(seed, pk, proof, max)

		assert.Equal(t, result, true)
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
	tries := 1 << 8              // 2^8
	watermark := uint64(1 << 56) // 2^56
	maxHits := 5
	badRNG := "The random number generator on this system is clearly " +
		"terrible since we got %d values less than %d in %d runs " +
		"when only %d was expected"
	_, pv := bls.GenerateTestKeyPair()

	signer := crypto.NewSigner(pv)

	numHits := 0
	for i := 0; i < tries; i++ {
		h := hash.GenerateTestHash()
		sig := signer.SignData(h.Bytes())
		seed, _ := VerifiableSeedFromBytes(sig.Bytes())

		nonce, _ := evaluate(seed, signer, util.MaxUint64)
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

func TestGetIndex(t *testing.T) {
	//  Total: 1000000

	// proof: 1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03
	// proofH: a7b8166584387f4ea76f9caa0969bd6b0bb8df4c3bb8e87f8b6e4dad62bf3359
	//
	// 0x4e7f38846516b8a7 * 1000000 / 0xffffffffffffffff = 306628.73489662957134845753
	proof1, _ := ProofFromString("1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03")
	assert.Equal(t, getIndex(proof1, 1*1e6), uint64(306628))

	// proof: 45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396
	// proofH: 80212979d1de1ca4ce1258fc0be66a4453b3804e64a5ca8d95f7def2c291c7fe
	//
	// 0xa41cded179292180 * 1000000 / 0xffffffffffffffff = 641065.52709668964549350137
	proof2, _ := ProofFromString("45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396")
	assert.Equal(t, getIndex(proof2, 1*1e6), uint64(641065))
}
