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

var (
	tSigner crypto.Signer
	tProof  Proof
	tSeed   VerifiableSeed
)

func setup(t *testing.T) {
	prv, err := bls.PrivateKeyFromString("39bc26dfcd0a5aec45cd2375122dffe46f713b6f93bc06c1fed759c251d4a13b")
	assert.NoError(t, err)
	tSigner = crypto.NewSigner(prv)
	tSeed, _ = VerifiableSeedFromString("94023a1ca49ab6583d2d382a2cee07ec65126af0311e67172f895fba4f1d1fa72e42b2d8918aa6c2d22bbaf0c80b26d8")
	tProof, _ = ProofFromString("98e7feee357ea2414f15665d42e223f0ca63e76516618c0586c4cfe83fa68b373cd301ba17f5e1aafe618c432c719c58")
}

func TestEvaluation(t *testing.T) {
	setup(t)

	t.Run("Total stake is zero", func(t *testing.T) {
		seed := GenerateRandomSeed()
		threshold := util.RandInt64(1 * 10e14)
		ok, proof := EvaluateSortition(seed, tSigner, 0, threshold)
		require.True(t, ok)
		ok = VerifyProof(seed, proof, tSigner.PublicKey(), 0, threshold)
		require.True(t, ok)
	})

	t.Run("Total stake is not zero, but validator stake is zero", func(t *testing.T) {
		seed := GenerateRandomSeed()
		total := util.RandInt64(1 * 10e14)

		ok, _ := EvaluateSortition(seed, tSigner, total, 0)
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
		proof, _ := ProofFromString("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		require.False(t, VerifyProof(seed, proof, pub, total, total))
	})

	t.Run("Sortition ok", func(t *testing.T) {
		total := int64(1 * 1e9)

		ok, _ := EvaluateSortition(tSeed, tSigner, total, total/10)
		require.True(t, ok)

		require.True(t, VerifyProof(tSeed, tProof, tSigner.PublicKey(), total, total/10))
		require.False(t, VerifyProof(tSeed, GenerateRandomProof(), tSigner.PublicKey(), total, total/10))
		require.False(t, VerifyProof(tSeed, Proof{}, tSigner.PublicKey(), total, total/10))
		require.False(t, VerifyProof(GenerateRandomSeed(), tProof, tSigner.PublicKey(), total, total/10))
	})
}

func TestSortitionMedian(t *testing.T) {
	setup(t)

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
