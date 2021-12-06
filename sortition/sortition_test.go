package sortition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestEvaluation(t *testing.T) {
	signer := bls.GenerateTestSigner()

	t.Run("Pool stake is zero", func(t *testing.T) {
		s := NewSortition()
		h := hash.GenerateTestHash()
		s.SetParams(h, GenerateRandomSeed(), 0)

		valStake := int64(1000000)
		ok, proof := s.EvaluateSortition(h, signer, valStake)
		require.True(t, ok)
		ok = s.VerifyProof(h, proof, signer.PublicKey(), valStake)
		require.True(t, ok)
	})

	t.Run("Pool stake is not zero, but validator stake is zero", func(t *testing.T) {
		s := NewSortition()
		h := hash.GenerateTestHash()
		s.SetParams(h, GenerateRandomSeed(), 1*1e9)

		ok, _ := s.EvaluateSortition(h, signer, 0)
		require.False(t, ok)
	})

	t.Run("Sortition ok", func(t *testing.T) {
		seed, _ := SeedFromString("8d019192c24224e2cafccae3a61fb586b14323a6bc8f9e7df1d929333ff993933bea6f5b3af6de0374366c4719e43a1b")
		prv, _ := bls.PrivateKeyFromString("39bc26dfcd0a5aec45cd2375122dffe46f713b6f93bc06c1fed759c251d4a13b")
		signer := crypto.NewSigner(prv)
		poolStake := int64(1 * 1e9)
		s := NewSortition()
		h := hash.GenerateTestHash()
		s.SetParams(h, seed, poolStake)

		ok, _ := s.EvaluateSortition(hash.GenerateTestHash(), signer, poolStake/10)
		require.False(t, ok)

		ok, proof := s.EvaluateSortition(h, signer, poolStake/10)
		require.True(t, ok)

		require.True(t, s.VerifyProof(h, proof, signer.PublicKey(), poolStake/10))
		require.False(t, s.VerifyProof(h, GenerateRandomProof(), signer.PublicKey(), poolStake/10))
		require.False(t, s.VerifyProof(h, Proof{}, signer.PublicKey(), poolStake/10))
		require.False(t, s.VerifyProof(hash.GenerateTestHash(), proof, signer.PublicKey(), poolStake/10))
	})
}

func TestVerifyProof(t *testing.T) {
	seed, _ := SeedFromString("8d019192c24224e2cafccae3a61fb586b14323a6bc8f9e7df1d929333ff993933bea6f5b3af6de0374366c4719e43a1b")
	pub, _ := bls.PublicKeyFromString("9a267cac764b1d860f1d587d0d5a61110c0c21bc6a57bdfdb8d4f2941e59fe709a017a32a599a35e81b91255d1b9d500f2427135a97d89a0a9431946d5db35d539bbe33f9f9b534c2cf88ef1a532f9d52a065a45221d18d6d4e6912680a5b58f")
	proof, _ := ProofFromString("2fbbe418b7b12068b2cfe43138e02453ea0146b1345381c72061274483af580f1c47a3e626c4927431c5447346860084")
	poolStake := int64(1 * 1e9)
	s := NewSortition()
	h := hash.GenerateTestHash()
	s.SetParams(h, seed, poolStake)

	assert.True(t, s.VerifyProof(h, proof, pub, poolStake/10))
	assert.False(t, s.VerifyProof(h, proof, pub, poolStake/30))
}

func TestSortitionMedian(t *testing.T) {
	poolStake := int64(1 * 1e9)
	valStake := poolStake / 10

	s := NewSortition()
	h := hash.GenerateTestHash()

	signer := bls.GenerateTestSigner()
	total := 1000
	median := 0
	for j := 0; j < total; j++ {
		seed := GenerateRandomSeed()
		s.SetParams(h, seed, poolStake)
		ok, _ := s.EvaluateSortition(h, signer, valStake)
		if ok {
			median++
		}
	}

	// Should be about 10%
	fmt.Printf("%v%% \n", median*100/total)
	assert.GreaterOrEqual(t, median*100/total, 5)
	assert.LessOrEqual(t, median*100/total, 15)
	assert.NotZero(t, median*100/total)
}

func TestExpiredProof(t *testing.T) {
	seed, _ := SeedFromString("65fd6c247d843cd80827a7a24cf01e1fbb697bd9e255fa259b745be24fe5bdce944c02b24d3a86b2c6460111f2876a88")
	pub, _ := bls.PublicKeyFromString("7002d6264285782be3ea70f231b123330ace6c6dc0b70a80fef4271e9379da2c60f63554e99bbf55877744c218e09a183368703ad432cc0a4b73509050f4a31695fc525468feee379339bd61fbc4b54d49ef997618be7c51c1ac3fd4ea185d97")
	proof, _ := ProofFromString("70e4951675331ce0bba3701f9c442889a6ff7b8364af1174cec27dedcbc90cfc9da1cf920ad6af64ffe70d9cfe826a0c")
	poolStake := int64(884 * 1e8)
	s := NewSortition()
	h := hash.GenerateTestHash()
	s.SetParams(h, seed, poolStake)

	for i := 0; i < 3; i++ {
		s.SetParams(hash.GenerateTestHash(), GenerateRandomSeed(), poolStake)
	}
	assert.True(t, s.VerifyProof(h, proof, pub, 21*1e8), "Sortition is valid")

	for i := 0; i < 4; i++ {
		s.SetParams(hash.GenerateTestHash(), GenerateRandomSeed(), poolStake)
	}
	assert.False(t, s.VerifyProof(h, proof, pub, 21*1e8), "Sortition expired")
}

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
