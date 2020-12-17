package sortition

import (
	"fmt"
	"testing"

	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/validator"

	"github.com/zarbchain/zarb-go/crypto"
)

func TestEvaluationTotalStakeZero(t *testing.T) {
	_, pub, priv := crypto.GenerateTestKeyPair()
	_, invPub, invPriv := crypto.GenerateTestKeyPair()
	val := validator.NewValidator(pub, 0, 0)
	invVal := validator.NewValidator(invPub, 0, 0)

	s := NewSortition(crypto.NewSigner(priv))
	h := crypto.GenerateTestHash()
	invHash := crypto.GenerateTestHash()

	invSortition := NewSortition(crypto.NewSigner(invPriv))
	trx := invSortition.EvaluateTransaction(h, val)
	require.Nil(t, trx)

	trx = s.EvaluateTransaction(h, val)
	require.NotNil(t, trx)
	proof := trx.Payload().(*payload.SortitionPayload).Proof
	assert.True(t, s.VerifyProof(h, proof, val))
	assert.False(t, s.VerifyProof(invHash, proof, val))
	assert.False(t, s.VerifyProof(h, proof, invVal))
	proof[0] = proof[0] + 1
	assert.False(t, s.VerifyProof(h, proof, val))

}

func TestEvaluationTotalStakeNotZero(t *testing.T) {
	_, pub, priv := crypto.GenerateTestKeyPair()
	val := validator.NewValidator(pub, 0, 0)
	val.AddToStake(100000000) // 1/10 of total stake

	s := NewSortition(crypto.NewSigner(priv))
	s.AddToTotalStake(1000000000)

	total := 100
	median := 0
	for j := 0; j < total; j++ {

		var h crypto.Hash
		var trx *tx.Tx
		i := 0
		for ; ; i++ {
			h = crypto.GenerateTestHash()
			trx = s.EvaluateTransaction(h, val)
			if trx != nil {
				break
			}
		}
		median += i
		require.NotNil(t, trx)

		proof := trx.Payload().(*payload.SortitionPayload).Proof
		assert.True(t, s.VerifyProof(h, proof, val))
	}

	median /= total
	fmt.Printf("%v ", median)

}
