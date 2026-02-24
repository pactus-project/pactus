package sortition_test

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	prv, _ := bls.PrivateKeyFromString(
		"SECRET1P838V87AW42JS8YWLYYK0AYFJQ9445VR72H23D6LR7GEJ8KW9UQ0QVE8WHE")
	seed, _ := sortition.VerifiableSeedFromString(
		"b63179137423ab2da8279d7aa3726d7ad05ae7d3ab3f744db0a9a719d12a720e72dc1d1e9222360243007f2f4adf7009")
	valKey := bls.NewValidatorKey(prv)

	t.Run("Total stake is zero", func(t *testing.T) {
		threshold := ts.RandInt64Max(int64(1e14))
		ok, proof := sortition.EvaluateSortition(seed, valKey.PrivateKey(), 0, threshold)
		require.True(t, ok)
		ok = sortition.VerifyProof(seed, proof, valKey.PublicKey(), 0, threshold)
		require.True(t, ok)
	})

	t.Run("Total stake is not zero, but validator stake is zero", func(t *testing.T) {
		total := ts.RandInt64Max(int64(1e14))

		ok, _ := sortition.EvaluateSortition(seed, valKey.PrivateKey(), total, 0)
		require.False(t, ok)
	})

	t.Run("OK!", func(t *testing.T) {
		proof1, _ := sortition.ProofFromString(
			"8cb689ec126465ddadd32493b71dc7ee3bfa2ef5a0a0f4b9b8aa777fb915a5f88def3305a3579e97b96ac862a6d67316")
		total := int64(1 * 1e14)

		ok, proof2 := sortition.EvaluateSortition(seed, valKey.PrivateKey(), total, total/100)
		require.True(t, ok)
		require.Equal(t, proof1, proof2)

		require.True(t, sortition.VerifyProof(seed, proof1, valKey.PublicKey(), total, total/100))
		require.False(t, sortition.VerifyProof(seed, proof1, valKey.PublicKey(), total, 0))
		require.False(t, sortition.VerifyProof(seed, ts.RandProof(), valKey.PublicKey(), total, total/10))
		require.False(t, sortition.VerifyProof(seed, sortition.Proof{}, valKey.PublicKey(), total, total/10))
		require.False(t, sortition.VerifyProof(ts.RandSeed(), proof1, valKey.PublicKey(), total, total/10))
	})
}

func TestInvalidProof(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("Invalid proof (Zero proof)", func(t *testing.T) {
		total := ts.RandInt64Max(int64(1e14))
		seed := ts.RandSeed()
		pub, _ := ts.RandBLSKeyPair()
		proof, _ := sortition.ProofFromString(
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		require.False(t, sortition.VerifyProof(seed, proof, pub, total, total))
	})

	t.Run("Invalid proof (Infinity proof)", func(t *testing.T) {
		total := ts.RandInt64Max(int64(1e14))
		seed := ts.RandSeed()
		pub, _ := ts.RandBLSKeyPair()
		proof, _ := sortition.ProofFromString(
			"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

		require.False(t, sortition.VerifyProof(seed, proof, pub, total, total))
	})
}

func TestSortitionMedian(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	total := int64(1 * 1e9)
	valKey := ts.RandValKey()

	count := 1000
	median := 0
	for j := 0; j < count; j++ {
		seed := ts.RandSeed()
		ok, _ := sortition.EvaluateSortition(seed, valKey.PrivateKey(), total, total/10)
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
