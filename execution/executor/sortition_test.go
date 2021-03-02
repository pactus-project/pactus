package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteSortitionTx(t *testing.T) {
	setup(t)
	exe := NewSortitionExecutor(tSandbox)

	stamp40 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(40, stamp40)
	proof1 := sortition.GenerateRandomProof()

	t.Run("Should fail, Bonding period", func(t *testing.T) {
		trx := tx.NewSortitionTx(stamp40, 1, tValSigner.Address(), proof1)

		assert.Error(t, exe.Execute(trx))
	})

	stamp41 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(41, stamp41)

	t.Run("Should fail, Invalid address", func(t *testing.T) {
		addr, _, _ := crypto.GenerateTestKeyPair()
		trx := tx.NewSortitionTx(stamp41, 1, addr, proof1)

		assert.Error(t, exe.Execute(trx))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(stamp41, 2, tValSigner.Address(), proof1)

		assert.Error(t, exe.Execute(trx))
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(stamp41, 1, tValSigner.Address(), proof1)

		// Check if can't join to validator set
		tSandbox.AcceptSortition = true
		tSandbox.ErrorAddToSet = true
		assert.Error(t, exe.Execute(trx))

		// Check if proof is wrong
		tSandbox.AcceptSortition = false
		tSandbox.ErrorAddToSet = true
		assert.Error(t, exe.Execute(trx))

		// Sounds good
		tSandbox.AcceptSortition = true
		tSandbox.ErrorAddToSet = false
		assert.NoError(t, exe.Execute(trx))

		// replay
		assert.Error(t, exe.Execute(trx))
	})

	val := tSandbox.Validator(tValSigner.Address())
	assert.Equal(t, val.Sequence(), 1)
	assert.Equal(t, val.LastJoinedHeight(), 42)
	assert.Zero(t, exe.Fee())

	checkTotalCoin(t, 0)
}
