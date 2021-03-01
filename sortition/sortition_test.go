package sortition

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestEvaluation(t *testing.T) {
	signer := crypto.GenerateTestSigner()
	seed := [48]byte{}
	rand.Read(seed[:])

	t.Run("Total stake is zero", func(t *testing.T) {
		sortition := NewSortition()

		require.NotNil(t, sortition.EvaluateSortition(seed, signer, 0))
	})

	t.Run("Total stake is not zero, but validator stake is zero", func(t *testing.T) {
		sortition := NewSortition()
		sortition.SetTotalStake(1000)

		require.Nil(t, sortition.EvaluateSortition(seed, signer, 0))
	})

	t.Run("Sortition ok", func(t *testing.T) {
		seed := [48]byte{}
		d, _ := hex.DecodeString("8d019192c24224e2cafccae3a61fb586b14323a6bc8f9e7df1d929333ff993933bea6f5b3af6de0374366c4719e43a1b")
		copy(seed[:], d)
		priv, _ := crypto.PrivateKeyFromString("39bc26dfcd0a5aec45cd2375122dffe46f713b6f93bc06c1fed759c251d4a13b")
		signer := crypto.NewSigner(priv)
		totalStake := int64(1000000000)
		sortition := NewSortition()
		sortition.AddToTotalStake(totalStake)

		proof := sortition.EvaluateSortition(seed, signer, totalStake/10)
		fmt.Printf("proof: %x\n", proof)

		require.NotNil(t, proof)
	})
}

func TestVerifyProof(t *testing.T) {
	seed := [48]byte{}
	d, _ := hex.DecodeString("8d019192c24224e2cafccae3a61fb586b14323a6bc8f9e7df1d929333ff993933bea6f5b3af6de0374366c4719e43a1b")
	copy(seed[:], d)
	pub, _ := crypto.PublicKeyFromString("9a267cac764b1d860f1d587d0d5a61110c0c21bc6a57bdfdb8d4f2941e59fe709a017a32a599a35e81b91255d1b9d500f2427135a97d89a0a9431946d5db35d539bbe33f9f9b534c2cf88ef1a532f9d52a065a45221d18d6d4e6912680a5b58f")
	proof, _ := hex.DecodeString("2fbbe418b7b12068b2cfe43138e02453ea0146b1345381c72061274483af580f1c47a3e626c4927431c5447346860084")
	totalStake := int64(1000000000)
	sortition := NewSortition()
	sortition.AddToTotalStake(totalStake)

	assert.True(t, sortition.VerifyProof(seed, proof, pub, totalStake/10))
	assert.False(t, sortition.VerifyProof(seed, proof, pub, totalStake/30))
}

func TestEvaluationSortition(t *testing.T) {
	stake := int64(100000000) // 1/10 of total stake
	totalStake := 10 * stake

	s := NewSortition()
	s.SetTotalStake(totalStake)

	signer := crypto.GenerateTestSigner()
	seed := [48]byte{}
	total := 500
	median := 0
	for j := 0; j < total; j++ {
		rand.Read(seed[:])
		proof := s.EvaluateSortition(seed, signer, stake)
		if proof != nil {
			median++
		}
	}

	// Should be about 10%
	fmt.Printf("%v%% ", median*100/total)
	assert.NotZero(t, median*100/total)
}
