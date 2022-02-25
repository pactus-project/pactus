package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteSortitionTx(t *testing.T) {
	setup(t)
	exe := NewSortitionExecutor(true)

	hash40 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(40, hash40)
	proof1 := sortition.GenerateRandomProof()

	t.Run("Should fail, Bonding period", func(t *testing.T) {
		trx := tx.NewSortitionTx(hash40.Stamp(), 1, tValSigner.Address(), proof1)
		tSandbox.AcceptSortition = true

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	hash41 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(41, hash41)

	t.Run("Should fail, Invalid address", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		trx := tx.NewSortitionTx(hash41.Stamp(), 1, addr, proof1)
		tSandbox.AcceptSortition = true

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(hash41.Stamp(), 2, tValSigner.Address(), proof1)
		tSandbox.AcceptSortition = true

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(hash41.Stamp(), 1, tValSigner.Address(), proof1)

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
		tSandbox.Validator(tValSigner.Address()).UpdateLastBondingHeight(3)
		assert.Error(t, exe.Execute(trx, tSandbox))

		// Sounds good
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = true
		tSandbox.Validator(tValSigner.Address()).UpdateLastBondingHeight(0)
		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Check unbond state
		tSandbox.AcceptSortition = true
		tSandbox.WelcomeToCommittee = true
		tSandbox.Validator(tValSigner.Address()).UpdateUnbondingHeight(1)
		assert.Error(t, exe.Execute(trx, tSandbox))

		// replay
		assert.Error(t, exe.Execute(trx, tSandbox))

	})

	val := tSandbox.Validator(tValSigner.Address())
	assert.Equal(t, val.Sequence(), 1)
	assert.Equal(t, val.LastJoinedHeight(), 42)
	assert.Zero(t, exe.Fee())

	checkTotalCoin(t, 0)
}

func TestSortitionNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewSortitionExecutor(false)

	hash100 := hash.GenerateTestHash()
	hash101 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(100, hash100)
	tSandbox.AppendNewBlock(101, hash101)
	proof1 := sortition.GenerateRandomProof()
	proof2 := sortition.GenerateRandomProof()

	tSandbox.AcceptSortition = true
	tSandbox.WelcomeToCommittee = false

	sortition1 := tx.NewSortitionTx(hash100.Stamp(), 102, tValSigner.Address(), proof1)
	sortition2 := tx.NewSortitionTx(hash101.Stamp(), 102, tValSigner.Address(), proof2)

	assert.NoError(t, exe1.Execute(sortition1, tSandbox))
	assert.NoError(t, exe1.Execute(sortition2, tSandbox))
}
