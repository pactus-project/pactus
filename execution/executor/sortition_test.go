package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func updateCommittee(td *testData) {
	joiningCommittee := make([]*validator.Validator, 0)
	td.sandbox.IterateValidators(func(val *validator.Validator, _ bool, joined bool) {
		if joined {
			joiningCommittee = append(joiningCommittee, val)
		}
	})

	td.sandbox.TestCommittee.Update(0, joiningCommittee)
	td.sandbox.TestJoinedValidators = make(map[crypto.Address]bool)
}

func TestExecuteSortitionTx(t *testing.T) {
	td := setup(t)
	exe := NewSortitionExecutor(true)

	existingVal := td.sandbox.TestStore.RandomTestVal()
	pub, _ := td.RandBLSKeyPair()
	newVal := td.sandbox.MakeNewValidator(pub)
	accAddr, acc := td.sandbox.TestStore.RandomTestAcc()
	amt, fee := td.randomAmountAndFee(0, acc.Balance())
	newVal.AddToStake(amt + fee)
	acc.SubtractFromBalance(amt + fee)
	td.sandbox.UpdateAccount(accAddr, acc)
	td.sandbox.UpdateValidator(newVal)
	lockTime := td.sandbox.CurrentHeight()
	proof := td.RandProof()

	newVal.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(newVal)
	assert.Zero(t, td.sandbox.Validator(newVal.Address()).LastSortitionHeight())
	assert.False(t, td.sandbox.IsJoinedCommittee(newVal.Address()))

	t.Run("Should fail, Invalid address", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, td.RandAccAddress(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	newVal.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval + 1)
	td.sandbox.UpdateValidator(newVal)

	t.Run("Should fail, Bonding period", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	// Let's add one more block
	td.sandbox.TestStore.AddTestBlock(td.randHeight + 1)

	t.Run("Should fail, Invalid proof", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = false
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidProof)
	})

	t.Run("Should fail, Committee has free seats and validator is in the committee", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, existingVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)
	})

	t.Run("Should fail, duplicated sortition", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	assert.Equal(t, td.sandbox.CurrentHeight(), td.randHeight+2)
	assert.Equal(t, td.sandbox.Validator(newVal.Address()).LastSortitionHeight(), lockTime)
	assert.True(t, td.sandbox.IsJoinedCommittee(newVal.Address()))

	td.checkTotalCoin(t, 0)
}

func TestSortitionNonStrictMode(t *testing.T) {
	td := setup(t)
	exe1 := NewSortitionExecutor(true)
	exe2 := NewSortitionExecutor(false)

	val := td.sandbox.TestStore.RandomTestVal()
	lockTime := td.sandbox.CurrentHeight()
	proof := td.RandProof()

	td.sandbox.TestAcceptSortition = true
	trx := tx.NewSortitionTx(lockTime, val.Address(), proof)
	err := exe1.Execute(trx, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	err = exe2.Execute(trx, td.sandbox)
	assert.NoError(t, err)
}

func TestChangePower1(t *testing.T) {
	td := setup(t)

	exe := NewSortitionExecutor(true)

	// This moves proposer to next validator
	updateCommittee(td)

	// Let's create validators first
	pub1, _ := td.RandBLSKeyPair()
	amt1 := td.sandbox.Committee().TotalPower() / 3
	val1 := td.sandbox.MakeNewValidator(pub1)
	val1.AddToStake(amount.Amount(amt1 - 1))
	val1.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val1)
	proof1 := td.RandProof()

	pub2, _ := td.RandBLSKeyPair()
	val2 := td.sandbox.MakeNewValidator(pub2)
	val2.AddToStake(2)
	val2.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val2)
	lockTime := td.sandbox.CurrentHeight()
	proof2 := td.RandProof()
	val3 := td.sandbox.Committee().Validators()[0]
	proof3 := td.RandProof()

	td.sandbox.TestParams.CommitteeSize = 4
	td.sandbox.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(lockTime, val1.Address(), proof1)
	err := exe.Execute(trx1, td.sandbox)
	assert.NoError(t, err)

	trx2 := tx.NewSortitionTx(lockTime, val2.Address(), proof2)
	err = exe.Execute(trx2, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidTx, "More than 1/3 of power is joining at the same height")

	// Val3 is a Committee member
	trx3 := tx.NewSortitionTx(lockTime, val3.Address(), proof3)
	err = exe.Execute(trx3, td.sandbox)
	assert.NoError(t, err)
}

func TestChangePower2(t *testing.T) {
	td := setup(t)

	exe := NewSortitionExecutor(true)

	// This moves proposer to next validator
	updateCommittee(td)

	// Let's create validators first
	pub1, _ := td.RandBLSKeyPair()
	val1 := td.sandbox.MakeNewValidator(pub1)
	val1.AddToStake(1)
	val1.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val1)
	proof1 := td.RandProof()

	pub2, _ := td.RandBLSKeyPair()
	val2 := td.sandbox.MakeNewValidator(pub2)
	val2.AddToStake(1)
	val2.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val2)
	proof2 := td.RandProof()

	pub3, _ := td.RandBLSKeyPair()
	val3 := td.sandbox.MakeNewValidator(pub3)
	val3.AddToStake(1)
	val3.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val3)
	lockTime := td.sandbox.CurrentHeight()
	proof3 := td.RandProof()
	val4 := td.sandbox.Committee().Validators()[0]
	proof4 := td.RandProof()

	td.sandbox.TestParams.CommitteeSize = 7
	td.sandbox.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(lockTime, val1.Address(), proof1)
	err := exe.Execute(trx1, td.sandbox)
	assert.NoError(t, err)

	trx2 := tx.NewSortitionTx(lockTime, val2.Address(), proof2)
	err = exe.Execute(trx2, td.sandbox)
	assert.NoError(t, err)

	trx3 := tx.NewSortitionTx(lockTime, val3.Address(), proof3)
	err = exe.Execute(trx3, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidTx, "More than 1/3 of power is leaving at the same height")

	// Committee member
	trx4 := tx.NewSortitionTx(lockTime, val4.Address(), proof4)
	err = exe.Execute(trx4, td.sandbox)
	assert.NoError(t, err)
}

// TestOldestDidNotPropose tests if the oldest validator in the committee had
// chance to propose a block or not.
func TestOldestDidNotPropose(t *testing.T) {
	td := setup(t)

	exe := NewSortitionExecutor(true)

	// Let's create validators first
	vals := make([]*validator.Validator, 9)
	for i := 0; i < 9; i++ {
		pub, _ := td.RandBLSKeyPair()
		val := td.sandbox.MakeNewValidator(pub)
		val.AddToStake(10 * 1e9)
		val.UpdateLastBondingHeight(
			td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
		td.sandbox.UpdateValidator(val)
		vals[i] = val
	}

	td.sandbox.TestParams.CommitteeSize = 7
	td.sandbox.TestAcceptSortition = true

	// This moves proposer to the next validator
	updateCommittee(td)

	// Let's update committee
	height := td.randHeight
	for i := uint32(0); i < 7; i++ {
		height++
		_ = td.sandbox.TestStore.AddTestBlock(height)

		lockTime := td.sandbox.CurrentHeight()
		trx1 := tx.NewSortitionTx(lockTime,
			vals[i].Address(), td.RandProof())
		err := exe.Execute(trx1, td.sandbox)
		assert.NoError(t, err)

		updateCommittee(td)
	}

	height++
	_ = td.sandbox.TestStore.AddTestBlock(height)
	lockTime := td.sandbox.CurrentHeight()

	trx1 := tx.NewSortitionTx(lockTime, vals[7].Address(), td.RandProof())
	err := exe.Execute(trx1, td.sandbox)
	assert.NoError(t, err)

	trx2 := tx.NewSortitionTx(lockTime, vals[8].Address(), td.RandProof())
	err = exe.Execute(trx2, td.sandbox)
	assert.NoError(t, err)
	updateCommittee(td)

	height++
	_ = td.sandbox.TestStore.AddTestBlock(height)
	// Entering validator 16
	trx3 := tx.NewSortitionTx(lockTime, vals[8].Address(), td.RandProof())
	err = exe.Execute(trx3, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
}
