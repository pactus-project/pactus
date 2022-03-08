package sortition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

func TestVRF(t *testing.T) {
	pk, pv := bls.GenerateTestKeyPair()
	signer := crypto.NewSigner(pv)
	for i := 0; i < 100; i++ {
		seed := GenerateRandomSeed()
		fmt.Printf("seed is: %x \n", seed)

		//max := int64(i * 1000)
		max := int64(1 * 1e6)
		index, proof := evaluate(seed, signer, max)
		// fmt.Printf("index is : %v \n", index)

		assert.LessOrEqual(t, index, max)

		index2, result := verify(seed, pk, proof, max)

		assert.Equal(t, result, true)
		assert.Equal(t, index, index2)
	}
}

/// TestRandomness is a naive test to check the randomness of the getIndex function
/// If we call getIndex 100 times, how many times we have unique number?
func TestRandomness(t *testing.T) {
	_, pv := bls.GenerateTestKeyPair()
	signer := crypto.NewSigner(pv)

	max := int64(100)

	entropy := make([]bool, max)
	for i := int64(0); i < max; i++ {
		seed := GenerateRandomSeed()

		index, _ := evaluate(seed, signer, max)
		assert.LessOrEqual(t, index, max)

		entropy[index] = true
	}

	hits := int64(0)
	for _, b := range entropy {
		if b == true {
			hits++
		}
	}

	fmt.Printf("Randomness is : %v%% \n", hits*100/max)
	assert.Greater(t, hits, int64(50))
}

func TestGetIndex(t *testing.T) {
	//  Total: 1000000

	// proof: 1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03
	// proofH: a7b8166584387f4ea76f9caa0969bd6b0bb8df4c3bb8e87f8b6e4dad62bf3359
	//
	// 0xa7b8166584387f4e & 0x7fffffffffffffff = 0x27b8166584387f4e
	// 0x27b8166584387f4e * 1000000 / 0x7fffffffffffffff = 310305.40425166817391776726
	proof1, _ := ProofFromString("1719b896ec1cc66a0f44c4bf90890d988e341cb2c1a808907780af844c854291536c12fdaef9a526bb7ef80da17c0b03")
	assert.Equal(t, getIndex(proof1, 1*1e6), int64(310305))

	// proof: 45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396
	// proofH: 80212979d1de1ca4ce1258fc0be66a4453b3804e64a5ca8d95f7def2c291c7fe
	//
	// 0x80212979d1de1ca4 & 0x7fffffffffffffff = 0x00212979d1de1ca4
	// 0x00212979d1de1ca4 * 1000000 / 0x7fffffffffffffff = 1012.02438575933141931421
	proof2, _ := ProofFromString("45180defab2daae377977bf09dcdd7d76ff4fc96d1b50cc8ac5a1601c0522fb11641c3ed0fefd4b1e1808c498d699396")
	assert.Equal(t, getIndex(proof2, 1*1e6), int64(1012))
}
