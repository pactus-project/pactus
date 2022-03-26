package sortition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/util"
)

func TestEvaluation(t *testing.T) {

	prv, err := bls.PrivateKeyFromString("0f09c13c87597d8e37a4070b9ee5f79cbda01404fdaadc2e8ea67b22531a568f")
	assert.NoError(t, err)
	signer := crypto.NewSigner(prv)
	seed, _ := VerifiableSeedFromString("90a3f39f4f15e89c312d9c88213acfeb8a5bfe196e39062b731a785d8ec651cf6a69a6c540e2bf03e71c55fb27c364fc")
	proof, _ := ProofFromString("b3020232051fa07763fdad6558619d3b54b0a38f1f69416ab6568bc88345ce08f372876b1d99df188e3fafa03ccfeb72")

	t.Run("Total stake is zero", func(t *testing.T) {
		seed := GenerateRandomSeed()
		threshold := util.RandInt64(1 * 10e14)
		ok, proof := EvaluateSortition(seed, signer, 0, threshold)
		require.True(t, ok)
		ok = VerifyProof(seed, proof, signer.PublicKey(), 0, threshold)
		require.True(t, ok)
	})

	t.Run("Total stake is not zero, but validator stake is zero", func(t *testing.T) {
		seed := GenerateRandomSeed()
		total := util.RandInt64(1 * 10e14)

		ok, _ := EvaluateSortition(seed, signer, total, 0)
		require.False(t, ok)
	})

	t.Run("Invalid proof (Infinity public key)", func(t *testing.T) {
		seed := GenerateRandomSeed()
		total := util.RandInt64(1 * 10e14)

		pub, _ := bls.PublicKeyFromString("c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		proof, _ := ProofFromString("c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		require.False(t, VerifyProof(seed, proof, pub, total, total))
	})

	t.Run("Invalid proof (Zero proof)", func(t *testing.T) {
		seed := GenerateRandomSeed()
		total := util.RandInt64(1 * 10e14)

		pub, _ := bls.GenerateTestKeyPair()
		proof, _ := ProofFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		require.False(t, VerifyProof(seed, proof, pub, total, total))
	})

	t.Run("Sortition ok", func(t *testing.T) {
		total := int64(1 * 1e9)

		ok, proof2 := EvaluateSortition(seed, signer, total, total/100)
		require.True(t, ok)
		require.Equal(t, proof, proof2)

		require.True(t, VerifyProof(seed, proof, signer.PublicKey(), total, total/10))
		require.False(t, VerifyProof(seed, proof, signer.PublicKey(), total, 0))
		require.False(t, VerifyProof(seed, GenerateRandomProof(), signer.PublicKey(), total, total/10))
		require.False(t, VerifyProof(seed, Proof{}, signer.PublicKey(), total, total/10))
		require.False(t, VerifyProof(GenerateRandomSeed(), proof, signer.PublicKey(), total, total/10))
	})

}

func TestSortitionMedian(t *testing.T) {
	total := int64(1 * 1e9)
	signer := bls.GenerateTestSigner()

	count := 1000
	median := 0
	for j := 0; j < count; j++ {
		seed := GenerateRandomSeed()
		ok, _ := EvaluateSortition(seed, signer, total, total/10)
		if ok {
			median++
		}
	}

	// Should be about 10%
	fmt.Printf("%v%% \n", median*100/count)
	assert.GreaterOrEqual(t, median*100/count, 5)
	assert.LessOrEqual(t, median*100/count, 15)
	assert.NotZero(t, median*100/count)
}
