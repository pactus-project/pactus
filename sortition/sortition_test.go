package sortition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

func TestEvaluation(t *testing.T) {
	t.Run("Signer should be same as validator", func(t *testing.T) {
		val, _ := validator.GenerateTestValidator(5)
		invSortition := NewSortition(crypto.GenerateTestSigner())
		h := crypto.GenerateTestHash()

		require.Nil(t, invSortition.EvaluateTransaction(h, val))
	})

	t.Run("Sortition on zero total stake", func(t *testing.T) {
		signer := crypto.GenerateTestSigner()
		val := validator.NewValidator(signer.PublicKey(), 1, 0)
		sortition := NewSortition(signer)
		h := crypto.GenerateTestHash()

		require.NotNil(t, sortition.EvaluateTransaction(h, val))
	})

	t.Run("Sortition on non-zero total stake", func(t *testing.T) {
		signer := crypto.GenerateTestSigner()
		val := validator.NewValidator(signer.PublicKey(), 1, 0)
		sortition := NewSortition(signer)
		sortition.AddToTotalStake(100000000)
		h := crypto.GenerateTestHash()

		require.Nil(t, sortition.EvaluateTransaction(h, val))
	})

	t.Run("Sortition ok", func(t *testing.T) {
		h, _ := crypto.HashFromString("a8fee3ae118a7c0c12a2a7a894af719655f2e9cbdfd2bbcd346c8cb99a0e71ba")
		priv, _ := crypto.PrivateKeyFromString("2b973a9589bd251341288cd8b19e62397ae5dd0367d45dbcbdfbe9b28253dd68")
		signer := crypto.NewSigner(priv)
		val := validator.NewValidator(signer.PublicKey(), 1, 0)
		val.AddToStake(1000)
		sortition := NewSortition(signer)
		sortition.AddToTotalStake(100000000)

		require.NotNil(t, sortition.EvaluateTransaction(h, val))
	})
}

func TestVerifyProof(t *testing.T) {
	h, _ := crypto.HashFromString("a8fee3ae118a7c0c12a2a7a894af719655f2e9cbdfd2bbcd346c8cb99a0e71ba")
	priv, _ := crypto.PrivateKeyFromString("2b973a9589bd251341288cd8b19e62397ae5dd0367d45dbcbdfbe9b28253dd68")
	signer := crypto.NewSigner(priv)
	val := validator.NewValidator(signer.PublicKey(), 1, 0)
	val.AddToStake(1000)
	sortition := NewSortition(signer)
	sortition.AddToTotalStake(100000000)

	index, proof := sortition.vrf.Evaluate(h)
	assert.Less(t, index, val.Stake())
	assert.True(t, sortition.VerifyProof(h, proof, val))

	t.Run("Invalid validator", func(t *testing.T) {
		anotherValidator, _ := validator.GenerateTestValidator(2)
		assert.False(t, sortition.VerifyProof(h, proof, anotherValidator))
	})

	t.Run("Less stake ", func(t *testing.T) {
		val2 := validator.NewValidator(signer.PublicKey(), 1, 0)
		val2.AddToStake(1000 / 2)
		assert.False(t, sortition.VerifyProof(h, proof, val2))
	})

}

func TestEvaluationTotalStakeNotZero(t *testing.T) {
	stake := int64(100000000)
	_, pub, priv := crypto.GenerateTestKeyPair()
	val := validator.NewValidator(pub, 0, 0)
	val.AddToStake(stake) // 1/10 of total stake

	s := NewSortition(crypto.NewSigner(priv))
	s.SetTotalStake(stake)
	s.AddToTotalStake(stake)
	assert.Equal(t, s.TotalStake(), 2*stake)

	total := 100
	median := 0
	for j := 0; j < total; j++ {
		var h crypto.Hash
		var trx *tx.Tx

		h = crypto.GenerateTestHash()
		trx = s.EvaluateTransaction(h, val)
		if trx != nil {
			median++
		}
	}

	// Should be in range of 40 to 60
	fmt.Printf("%v ", median*100/total)
}
