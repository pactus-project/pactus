package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func updateCommittee(td *testData) {
	joiningCommittee := make([]*validator.Validator, 0)
	td.sandbox.IterateValidators(func(val *validator.Validator, updated bool, joined bool) {
		if joined {
			joiningCommittee = append(joiningCommittee, val)
		}
	})

	td.sandbox.TestCommittee.Update(0, joiningCommittee)
	td.sandbox.TestJoinedValidators = make(map[crypto.Address]bool)
	// fmt.Println(td.sandbox.TestCommittee.String())
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

	proof := td.RandProof()

	newVal.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(newVal)
	assert.Zero(t, td.sandbox.Validator(newVal.Address()).LastSortitionHeight())
	assert.False(t, td.sandbox.IsJoinedCommittee(newVal.Address()))

	t.Run("Should fail, Invalid address", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.randStamp, 1, td.RandAddress(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	newVal.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval + 1)
	td.sandbox.UpdateValidator(newVal)

	t.Run("Should fail, Bonding period", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.randStamp, newVal.Sequence()+1, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	// Let's add one more block
	td.sandbox.TestStore.AddTestBlock(td.randHeight + 1)

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.randStamp, newVal.Sequence(), newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidLockTime)
	})

	t.Run("Should fail, Invalid proof", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.randStamp, newVal.Sequence()+1, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = false
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidProof)
	})

	t.Run("Should fail, Committee has free seats and validator is in the committee", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.randStamp, existingVal.Sequence()+1, existingVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.randStamp, newVal.Sequence()+1, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err)

		// Execute again, should fail
		err = exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidLockTime)
	})

	t.Run("Should fail, duplicated sortition", func(t *testing.T) {
		trx := tx.NewSortitionTx(td.randStamp, newVal.Sequence()+2, newVal.Address(), proof)
		td.sandbox.TestAcceptSortition = true
		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	assert.Equal(t, td.sandbox.CurrentHeight(), td.randHeight+2)
	assert.Equal(t, td.sandbox.Validator(newVal.Address()).LastSortitionHeight(), td.randHeight)
	assert.True(t, td.sandbox.IsJoinedCommittee(newVal.Address()))
	assert.Zero(t, exe.Fee())

	td.checkTotalCoin(t, 0)
}

func TestSortitionNonStrictMode(t *testing.T) {
	td := setup(t)
	exe1 := NewSortitionExecutor(true)
	exe2 := NewSortitionExecutor(false)

	val := td.sandbox.TestStore.RandomTestVal()
	proof := td.RandProof()

	td.sandbox.TestAcceptSortition = true
	trx := tx.NewSortitionTx(td.randStamp, val.Sequence(), val.Address(), proof)
	err := exe1.Execute(trx, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidLockTime)
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
	val1.AddToStake(amt1 - 1)
	val1.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val1)
	proof1 := td.RandProof()

	pub2, _ := td.RandBLSKeyPair()
	val2 := td.sandbox.MakeNewValidator(pub2)
	val2.AddToStake(2)
	val2.UpdateLastBondingHeight(td.sandbox.CurrentHeight() - td.sandbox.Params().BondInterval)
	td.sandbox.UpdateValidator(val2)
	proof2 := td.RandProof()

	val3 := td.sandbox.Committee().Validators()[0]
	proof3 := td.RandProof()

	td.sandbox.TestParams.CommitteeSize = 4
	td.sandbox.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(td.randStamp, val1.Sequence()+1, val1.Address(), proof1)
	err := exe.Execute(trx1, td.sandbox)
	assert.NoError(t, err)

	trx2 := tx.NewSortitionTx(td.randStamp, val2.Sequence()+1, val2.Address(), proof2)
	err = exe.Execute(trx2, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidTx, "More than 1/3 of power is joining at the same height")

	// Committee member
	trx3 := tx.NewSortitionTx(td.randStamp, val3.Sequence()+1, val3.Address(), proof3)
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
	proof3 := td.RandProof()

	val4 := td.sandbox.Committee().Validators()[0]
	proof4 := td.RandProof()

	td.sandbox.TestParams.CommitteeSize = 7
	td.sandbox.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(td.randStamp, val1.Sequence()+1, val1.Address(), proof1)
	err := exe.Execute(trx1, td.sandbox)
	assert.NoError(t, err)

	trx2 := tx.NewSortitionTx(td.randStamp, val2.Sequence()+1, val2.Address(), proof2)
	err = exe.Execute(trx2, td.sandbox)
	assert.NoError(t, err)

	trx3 := tx.NewSortitionTx(td.randStamp, val3.Sequence()+1, val3.Address(), proof3)
	err = exe.Execute(trx3, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidTx, "More than 1/3 of power is leaving at the same height")

	// Committee member
	trx4 := tx.NewSortitionTx(td.randStamp, val4.Sequence()+1, val4.Address(), proof4)
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
		b := td.sandbox.TestStore.AddTestBlock(height)
		stamp := b.Stamp()

		trx1 := tx.NewSortitionTx(stamp,
			vals[i].Sequence()+1,
			vals[i].Address(), td.RandProof())
		err := exe.Execute(trx1, td.sandbox)
		assert.NoError(t, err)

		updateCommittee(td)
	}

	height++
	b := td.sandbox.TestStore.AddTestBlock(height)
	stamp := b.Stamp()

	trx1 := tx.NewSortitionTx(stamp, vals[7].Sequence()+1, vals[7].Address(), td.RandProof())
	err := exe.Execute(trx1, td.sandbox)
	assert.NoError(t, err)

	trx2 := tx.NewSortitionTx(stamp, vals[8].Sequence()+1, vals[8].Address(), td.RandProof())
	err = exe.Execute(trx2, td.sandbox)
	assert.NoError(t, err)
	updateCommittee(td)

	height++
	b = td.sandbox.TestStore.AddTestBlock(height)
	stamp = b.Stamp()

	// Entering validator 16
	trx3 := tx.NewSortitionTx(stamp, vals[8].Sequence()+2, vals[8].Address(), td.RandProof())
	err = exe.Execute(trx3, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
}
