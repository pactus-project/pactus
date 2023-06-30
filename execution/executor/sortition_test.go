package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestExecuteSortitionTx(t *testing.T) {
	td := setup(t)
	exe := NewSortitionExecutor(true)

	existingVal := td.sandbox.TestStore.RandomTestVal()
	pub, _ := td.RandomBLSKeyPair()
	newVal := td.sandbox.MakeNewValidator(pub)
	accAddr, acc := td.sandbox.TestStore.RandomTestAcc()
	amt, fee := td.randomAmountAndFee(acc.Balance())
	newVal.AddToStake(amt + fee)
	acc.SubtractFromBalance(amt + fee)
	td.sandbox.UpdateAccount(accAddr, acc)
	td.sandbox.UpdateValidator(newVal)

	proof := td.RandomProof()

	newVal.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(newVal)

	t.Run("Should fail, Invalid address", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.stamp500000, 1, td.RandomAddress(), proof)
		td.sandbox.TestAcceptSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidAddress)
	})

	newVal.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval + 1)
	td.sandbox.UpdateValidator(newVal)
	t.Run("Should fail, Bonding period", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.stamp500000, newVal.Sequence()+1, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidHeight)
	})

	td.sandbox.TestStore.AddTestBlock(500001)

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.stamp500000, newVal.Sequence()+2, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, Invalid proof", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.stamp500000, newVal.Sequence()+1, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = false
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidProof)
	})

	t.Run("Should fail, Committee has free seats and validator is in the committee", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.stamp500000, existingVal.Sequence()+1, existingVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, td.sandbox)), errors.ErrInvalidTx)
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.stamp500000, newVal.Sequence()+1, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		assert.NoError(t, exe.Execute(trx, td.sandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, td.sandbox))
	})

	assert.Equal(t, td.sandbox.Validator(newVal.Address()).LastJoinedHeight(), td.sandbox.CurrentHeight())
	assert.Zero(t, exe.Fee())

	td.checkTotalCoin(t, 0)
}

func TestSortitionNonStrictMode(t *testing.T) {
	td := setup(t)
	exe1 := NewSortitionExecutor(true)
	exe2 := NewSortitionExecutor(false)

	val := td.sandbox.TestStore.RandomTestVal()
	proof := td.RandomProof()

	td.sandbox.TestAcceptSortition = true
	trx := tx.NewSortitionTx(td.stamp500000, val.Sequence(), val.Address(), proof)
	assert.Error(t, exe1.Execute(trx, td.sandbox))
	assert.NoError(t, exe2.Execute(trx, td.sandbox))
}

func TestChangePower1(t *testing.T) {
	td := setup(t)

	exe := NewSortitionExecutor(true)

	// Let's create validators first
	pub1, _ := td.RandomBLSKeyPair()
	amt1 := td.sandbox.Committee().TotalPower() / 3
	val1 := td.sandbox.MakeNewValidator(pub1)
	val1.AddToStake(amt1 - 1)
	val1.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val1)
	proof1 := td.RandomProof()

	pub2, _ := td.RandomBLSKeyPair()
	val2 := td.sandbox.MakeNewValidator(pub2)
	val2.AddToStake(2)
	val2.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val2)
	proof2 := td.RandomProof()

	val3 := td.sandbox.Committee().Validators()[0]
	proof3 := td.RandomProof()

	td.sandbox.TestParams.CommitteeSize = 4
	td.sandbox.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(td.stamp500000, val1.Sequence()+1, val1.Address(), proof1)
	assert.NoError(t, exe.Execute(trx1, td.sandbox))

	trx2 := tx.NewSortitionTx(td.stamp500000, val2.Sequence()+1, val2.Address(), proof2)
	assert.Error(t, exe.Execute(trx2, td.sandbox), "More than 1/3 of power is joining at the same height")

	// Committee member
	trx3 := tx.NewSortitionTx(td.stamp500000, val3.Sequence()+1, val3.Address(), proof3)
	assert.NoError(t, exe.Execute(trx3, td.sandbox))
}

func TestChangePower2(t *testing.T) {
	td := setup(t)

	exe := NewSortitionExecutor(true)

	// Let's create validators first
	pub1, _ := td.RandomBLSKeyPair()
	val1 := td.sandbox.MakeNewValidator(pub1)
	val1.AddToStake(1)
	val1.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val1)
	proof1 := td.RandomProof()

	pub2, _ := td.RandomBLSKeyPair()
	val2 := td.sandbox.MakeNewValidator(pub2)
	val2.AddToStake(1)
	val2.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val2)
	proof2 := td.RandomProof()

	pub3, _ := td.RandomBLSKeyPair()
	val3 := td.sandbox.MakeNewValidator(pub3)
	val3.AddToStake(1)
	val3.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val3)
	proof3 := td.RandomProof()

	val4 := td.sandbox.Committee().Validators()[0]
	proof4 := td.RandomProof()

	td.sandbox.TestParams.CommitteeSize = 7
	td.sandbox.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(td.stamp500000, val1.Sequence()+1, val1.Address(), proof1)
	assert.NoError(t, exe.Execute(trx1, td.sandbox))

	trx2 := tx.NewSortitionTx(td.stamp500000, val2.Sequence()+1, val2.Address(), proof2)
	assert.NoError(t, exe.Execute(trx2, td.sandbox))

	trx3 := tx.NewSortitionTx(td.stamp500000, val3.Sequence()+1, val3.Address(), proof3)
	assert.Error(t, exe.Execute(trx3, td.sandbox), "More than 1/3 of power is leaving at the same height")

	// Committee member
	trx4 := tx.NewSortitionTx(td.stamp500000, val4.Sequence()+1, val4.Address(), proof4)
	assert.NoError(t, exe.Execute(trx4, td.sandbox))
}

// TestOldestDidNotPropose tests if the oldest validator in the committee had
// chance to propose a block or not.
func TestOldestDidNotPropose(t *testing.T) {
	td := setup(t)

	exe := NewSortitionExecutor(true)

	// Let's create validators first
	vals := make([]*validator.Validator, 9)
	for i := 0; i < 9; i++ {
		pub, _ := td.RandomBLSKeyPair()
		val := td.sandbox.MakeNewValidator(pub)
		val.AddToStake(10 * 1e9)
		val.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
		td.sandbox.UpdateValidator(val)
		vals[i] = val
	}

	td.sandbox.TestParams.CommitteeSize = 7
	td.sandbox.TestAcceptSortition = true

	stamp := td.stamp500000
	for i := 0; i < 8; i = i + 2 {
		b := td.sandbox.TestStore.AddTestBlock(uint32(500001 + (i / 2)))
		stamp = b.Stamp()

		trx1 := tx.NewSortitionTx(stamp, vals[i].Sequence()+1, vals[i].Address(), td.RandomProof())
		assert.NoError(t, exe.Execute(trx1, td.sandbox))

		trx2 := tx.NewSortitionTx(stamp, vals[i+1].Sequence()+1, vals[i+1].Address(), td.RandomProof())
		assert.NoError(t, exe.Execute(trx2, td.sandbox))

		joined := make([]*validator.Validator, 0)
		td.sandbox.IterateValidators(func(val *validator.Validator, updated bool) {
			if val.LastJoinedHeight() == td.sandbox.CurrentHeight() {
				joined = append(joined, val)
			}
		})
		td.sandbox.TestCommittee.Update(0, joined)
	}

	trx := tx.NewSortitionTx(stamp, vals[8].Sequence()+1, vals[8].Address(), td.RandomProof())
	assert.Error(t, exe.Execute(trx, td.sandbox))
}
