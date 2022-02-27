package sortition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestEvaluation(t *testing.T) {
	seed, _ := VerifiableSeedFromString("8d019192c24224e2cafccae3a61fb586b14323a6bc8f9e7df1d929333ff993933bea6f5b3af6de0374366c4719e43a1b")
	prv, _ := bls.PrivateKeyFromString("39bc26dfcd0a5aec45cd2375122dffe46f713b6f93bc06c1fed759c251d4a13b")
	signer := crypto.NewSigner(prv)
	totalStake := int64(100 * changeCoefficient)
	srt := NewSortition()
	blockHash := hash.GenerateTestHash()
	srt.SetParams(blockHash, seed, totalStake)

	t.Run("Validator stake is zero", func(t *testing.T) {
		ok, proof := srt.EvaluateSortition(blockHash, signer, 0)
		require.False(t, ok)
		require.Empty(t, proof)
	})

	t.Run("Sortition ok", func(t *testing.T) {
		ok, proof := srt.EvaluateSortition(blockHash, signer, totalStake/21)
		require.True(t, ok)

		require.True(t, srt.VerifyProof(blockHash, proof, signer.PublicKey(), totalStake/20))
		require.False(t, srt.VerifyProof(blockHash, GenerateRandomProof(), signer.PublicKey(), totalStake/20))
		require.False(t, srt.VerifyProof(blockHash, Proof{}, signer.PublicKey(), totalStake/20))
		require.False(t, srt.VerifyProof(hash.GenerateTestHash(), proof, signer.PublicKey(), totalStake/20))
	})
}

// func TestVerifyProof(t *testing.T) {
// 	seed, _ := VerifiableSeedFromString("836928e6a797858b17206018fba69d6d2640e84091d11aa345a260cda344d8c1b953b1dda2a415ad037b6cbe73e031da")
// 	pub, _ := bls.PublicKeyFromString("a2f1c33977381af6ec8e0ca68a5acfe61feacac89bf6035117b25727c82fb735e851151b100b3e0a056915fb9819dca906451c5e63149ae5dce8648ccbf372b19b17a344aaca1474bad9a0a65061bd43d6fb573c0203d31a27978bb0600dc8fa")
// 	proof, _ := ProofFromString("aab16ac3d8167e44ce9c842bf03e1d3bc5435ee9b7c72d6c6da8b862d54217895dc864f87930738cbcc1f60c2e4efa12")
// 	poolStake := int64(1 * 1e9)
// 	s := NewSortition()
// 	h := hash.GenerateTestHash()
// 	s.SetParams(h, seed, poolStake)

// 	assert.True(t, s.VerifyProof(h, proof, pub, poolStake/10))
// 	assert.False(t, s.VerifyProof(h, proof, pub, poolStake/30))
// }

// func TestSortitionMedian(t *testing.T) {
// 	poolStake := int64(1 * 1e9)
// 	valStake := poolStake / 10

// 	s := NewSortition()
// 	h := hash.GenerateTestHash()

// 	signer := bls.GenerateTestSigner()
// 	total := 1000
// 	median := 0
// 	for j := 0; j < total; j++ {
// 		seed := GenerateRandomSeed()
// 		s.SetParams(h, seed, poolStake)
// 		ok, _ := s.EvaluateSortition(h, signer, valStake)
// 		if ok {
// 			median++
// 		}
// 	}

// 	// Should be about 10%
// 	fmt.Printf("%v%% \n", median*100/total)
// 	assert.GreaterOrEqual(t, median*100/total, 5)
// 	assert.LessOrEqual(t, median*100/total, 15)
// 	assert.NotZero(t, median*100/total)
// }

// func TestExpiredProof(t *testing.T) {
// 	seed, _ := VerifiableSeedFromString("af38341e12ab8db809076ab49dd0879e26860ca906a085b1c9c476f2f203ccbf79c68dcc8ed0d0fb38c9317e84b92cd3")
// 	pub, _ := bls.PublicKeyFromString("b131e25335446039f76f1bdde662b07ae42d7504aa0bcef4d7703b07c0489cba603b138ef7a3146e63a9397c8d0681030652a5859cfca77b1cd29f60ec2093f326d030fec458dc941c40f836d8ca1d73873fe832e9ebb8cd7f4692eaaa09be55")
// 	proof, _ := ProofFromString("aec528b4647cad0083375a3395642400fbd2748f1d7f49d1643b67511ddeb4ae4e480bb8f9342446052d735381690d38")
// 	poolStake := int64(884 * 1e8)
// 	s := NewSortition()
// 	h := hash.GenerateTestHash()
// 	s.SetParams(h, seed, poolStake)

// 	for i := 0; i < 3; i++ {
// 		s.SetParams(hash.GenerateTestHash(), GenerateRandomSeed(), poolStake)
// 	}
// 	assert.True(t, s.VerifyProof(h, proof, pub, 21*1e8), "Sortition is valid")

// 	for i := 0; i < 4; i++ {
// 		s.SetParams(hash.GenerateTestHash(), GenerateRandomSeed(), poolStake)
// 	}
// 	assert.False(t, s.VerifyProof(h, proof, pub, 21*1e8), "Sortition expired")
// }

func TestGetParam(t *testing.T) {
	h1 := hash.GenerateTestHash()
	h2 := hash.GenerateTestHash()
	s1 := GenerateRandomSeed()
	s2 := GenerateRandomSeed()

	s := NewSortition()
	s.SetParams(h1, s1, 1000)
	s.SetParams(h2, s2, 2000)

	seed, stake := s.GetParams(h1)
	assert.Equal(t, seed, s1)
	assert.Equal(t, stake, int64(1000))

	seed, stake = s.GetParams(h2)
	assert.Equal(t, seed, s2)
	assert.Equal(t, stake, int64(2000))
}
