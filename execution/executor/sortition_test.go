package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func TestExecuteSortitionTx(t *testing.T) {
	setup(t)
	exe := NewSortitionExecutor(true)

	seq := tSandbox.Validator(tVal1.Address()).Sequence() + 1
	proof := sortition.GenerateRandomProof()

	tVal1.UpdateLastBondingHeight(tSandbox.CurHeight - tSandbox.BondInterval() + 1)
	t.Run("Should fail, Bonding period", func(t *testing.T) {

		trx := tx.NewSortitionTx(tHash1000000.Stamp(), seq, tVal1.Address(), proof)
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = true

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	tVal1.UpdateLastBondingHeight(tSandbox.CurHeight - tSandbox.BondInterval())
	t.Run("Should fail, Invalid address", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		trx := tx.NewSortitionTx(tHash1000001.Stamp(), seq, addr, proof)
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = true

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(tHash1000001.Stamp(), seq+1, tVal1.Address(), proof)
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = true

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(tHash1000001.Stamp(), seq, tVal1.Address(), proof)

		// Check if can't join to committee
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = false
		assert.Error(t, exe.Execute(trx, tSandbox))

		// Check if proof is wrong
		tSandbox.AcceptSortition = false
		tSandbox.WelcomeToCommittee = true
		assert.Error(t, exe.Execute(trx, tSandbox))

		// Check if power is 0
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = true
		tSandbox.Validator(tVal1.Address()).UpdateUnbondingHeight(3)
		assert.Error(t, exe.Execute(trx, tSandbox))

		// Sounds good
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = true
		tSandbox.Validator(tVal1.Address()).UpdateUnbondingHeight(0)
		assert.NoError(t, exe.Execute(trx, tSandbox))

		// replay
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	val := tSandbox.Validator(tVal1.Address())
	assert.Equal(t, val.Sequence(), 1)
	assert.Equal(t, val.LastJoinedHeight(), tSandbox.CurHeight)
	assert.Zero(t, exe.Fee())

	checkTotalCoin(t, 0)
}

func TestSortitionNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewSortitionExecutor(true)
	exe2 := NewSortitionExecutor(false)

	proof1 := sortition.GenerateRandomProof()
	proof2 := sortition.GenerateRandomProof()

	tSandbox.AcceptSortition = true
	tSandbox.WelcomeToCommittee = false

	sortition1 := tx.NewSortitionTx(tHash1000000.Stamp(), util.RandInt(100000), tVal1.Address(), proof1)
	assert.Error(t, exe1.Execute(sortition1, tSandbox))
	assert.NoError(t, exe2.Execute(sortition1, tSandbox))

	sortition2 := tx.NewSortitionTx(tHash1000001.Stamp(), tVal1.Sequence()+1, tVal1.Address(), proof2)
	assert.Error(t, exe1.Execute(sortition2, tSandbox))
	assert.NoError(t, exe2.Execute(sortition2, tSandbox))
}
