package sortition

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluation(t *testing.T) {
	prv, _ := bls.PrivateKeyFromString(
		"SECRET1P838V87AW42JS8YWLYYK0AYFJQ9445VR72H23D6LR7GEJ8KW9UQ0QVE8WHE")
	seed, _ := VerifiableSeedFromString(
		"b63179137423ab2da8279d7aa3726d7ad05ae7d3ab3f744db0a9a719d12a720e72dc1d1e9222360243007f2f4adf7009")
	proof, _ := ProofFromString(
		"8cb689ec126465ddadd32493b71dc7ee3bfa2ef5a0a0f4b9b8aa777fb915a5f88def3305a3579e97b96ac862a6d67316")
	signer := crypto.NewSigner(prv)

	t.Run("Total stake is zero", func(t *testing.T) {
		threshold := util.RandInt64(1 * 1e14)
		ok, proof := EvaluateSortition(seed, signer, 0, threshold)
		require.True(t, ok)
		ok = VerifyProof(seed, proof, signer.PublicKey(), 0, threshold)
		require.True(t, ok)
	})

	t.Run("Total stake is not zero, but validator stake is zero", func(t *testing.T) {
		total := util.RandInt64(1 * 1e14)

		ok, _ := EvaluateSortition(seed, signer, total, 0)
		require.False(t, ok)
	})

	t.Run("OK!", func(t *testing.T) {
		total := int64(1 * 1e14)

		ok, proof2 := EvaluateSortition(seed, signer, total, total/100)
		require.True(t, ok)
		require.Equal(t, proof, proof2)

		require.True(t, VerifyProof(seed, proof, signer.PublicKey(), total, total/100))
		require.False(t, VerifyProof(seed, proof, signer.PublicKey(), total, 0))
		require.False(t, VerifyProof(seed, GenerateRandomProof(), signer.PublicKey(), total, total/10))
		require.False(t, VerifyProof(seed, Proof{}, signer.PublicKey(), total, total/10))
		require.False(t, VerifyProof(GenerateRandomSeed(), proof, signer.PublicKey(), total, total/10))
	})
}

func TestInvalidProof(t *testing.T) {
	t.Run("Invalid proof (Zero proof)", func(t *testing.T) {
		total := util.RandInt64(1 * 1e14)
		seed := GenerateRandomSeed()
		pub, _ := bls.GenerateTestKeyPair()
		proof, _ := ProofFromString(
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		require.False(t, VerifyProof(seed, proof, pub, total, total))
	})

	t.Run("Invalid proof (Infinity proof)", func(t *testing.T) {
		total := util.RandInt64(1 * 1e14)
		seed := GenerateRandomSeed()
		pub, _ := bls.GenerateTestKeyPair()
		proof, _ := ProofFromString(
			"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		require.False(t, VerifyProof(seed, proof, pub, total, total))
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
