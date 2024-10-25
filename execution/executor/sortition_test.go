package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/stretchr/testify/assert"
)

func updateCommittee(td *testData) {
	joiningCommittee := make([]*validator.Validator, 0)
	td.sbx.IterateValidators(func(val *validator.Validator, _ bool, joined bool) {
		if joined {
			joiningCommittee = append(joiningCommittee, val)
		}
	})

	td.sbx.TestCommittee.Update(0, joiningCommittee)
	td.sbx.TestJoinedValidators = make(map[crypto.Address]bool)
}

func TestExecuteSortitionTx(t *testing.T) {
	td := setup(t)

	bonderAddr, bonderAcc := td.sbx.TestStore.RandomTestAcc()
	bonderBalance := bonderAcc.Balance()
	stake := td.RandAmountRange(
		td.sbx.TestParams.MinimumStake,
		bonderBalance)
	bonderAcc.SubtractFromBalance(stake)
	td.sbx.UpdateAccount(bonderAddr, bonderAcc)

	valPub, _ := td.RandBLSKeyPair()
	valAddr := valPub.ValidatorAddress()
	val := td.sbx.MakeNewValidator(valPub)
	val.AddToStake(stake)
	td.sbx.UpdateValidator(val)

	curHeight := td.sbx.CurrentHeight()
	lockTime := td.sbx.CurrentHeight()
	proof := td.RandProof()

	val.UpdateLastBondingHeight(curHeight - td.sbx.Params().BondInterval)
	td.sbx.UpdateValidator(val)

	assert.Zero(t, val.LastSortitionHeight())
	assert.False(t, td.sbx.IsJoinedCommittee(val.Address()))

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandAccAddress()
		trx := tx.NewSortitionTx(lockTime, randomAddr, proof)

		td.check(t, trx, true, ValidatorNotFoundError{Address: randomAddr})
		td.check(t, trx, false, ValidatorNotFoundError{Address: randomAddr})
	})

	val.UpdateLastBondingHeight(curHeight - td.sbx.Params().BondInterval + 1)
	td.sbx.UpdateValidator(val)

	t.Run("Should fail, Bonding period", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)

		td.check(t, trx, true, ErrBondingPeriod)
		td.check(t, trx, false, ErrBondingPeriod)
	})

	// Let's add one more block
	td.sbx.TestStore.AddTestBlock(curHeight + 1)

	t.Run("Should fail, invalid proof", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)
		td.sbx.TestAcceptSortition = false
		td.check(t, trx, true, ErrInvalidSortitionProof)
		td.check(t, trx, false, ErrInvalidSortitionProof)
	})

	t.Run("Should fail, committee has free seats and validator is in the committee", func(t *testing.T) {
		val0 := td.sbx.Committee().Proposer(0)
		trx := tx.NewSortitionTx(lockTime, val0.Address(), proof)
		td.sbx.TestAcceptSortition = true

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)
		td.sbx.TestAcceptSortition = true

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	t.Run("Should fail, expired sortition", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime-1, val.Address(), proof)
		td.sbx.TestAcceptSortition = true
		td.check(t, trx, true, ErrExpiredSortition)
		td.check(t, trx, false, ErrExpiredSortition)
	})

	t.Run("Should fail, duplicated sortition", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)

		td.check(t, trx, true, ErrExpiredSortition)
		td.check(t, trx, false, ErrExpiredSortition)
	})

	updatedVal := td.sbx.Validator(valAddr)

	assert.Equal(t, lockTime, updatedVal.LastSortitionHeight())
	assert.True(t, td.sbx.IsJoinedCommittee(val.Address()))

	td.checkTotalCoin(t, 0)
}

func TestChangePower1(t *testing.T) {
	td := setup(t)

	// This moves proposer to next validator
	updateCommittee(td)

	// Let's create validators first
	pub1, _ := td.RandBLSKeyPair()
	amt1 := td.sbx.Committee().TotalPower() / 3
	val1 := td.sbx.MakeNewValidator(pub1)
	val1.AddToStake(amount.Amount(amt1 - 1))
	val1.UpdateLastBondingHeight(td.sbx.CurrentHeight() - td.sbx.Params().BondInterval)
	td.sbx.UpdateValidator(val1)
	proof1 := td.RandProof()

	pub2, _ := td.RandBLSKeyPair()
	val2 := td.sbx.MakeNewValidator(pub2)
	val2.AddToStake(2)
	val2.UpdateLastBondingHeight(td.sbx.CurrentHeight() - td.sbx.Params().BondInterval)
	td.sbx.UpdateValidator(val2)
	lockTime := td.sbx.CurrentHeight()
	proof2 := td.RandProof()

	val3 := td.sbx.Committee().Proposer(0)
	proof3 := td.RandProof()

	td.sbx.TestParams.CommitteeSize = 4
	td.sbx.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(lockTime, val1.Address(), proof1)
	td.check(t, trx1, true, nil)
	td.check(t, trx1, false, nil)
	td.execute(t, trx1)

	trx2 := tx.NewSortitionTx(lockTime, val2.Address(), proof2)
	td.check(t, trx2, true, ErrCommitteeJoinLimitExceeded)
	td.check(t, trx2, false, nil)

	// Val3 is a Committee member
	trx3 := tx.NewSortitionTx(lockTime, val3.Address(), proof3)
	td.check(t, trx3, true, nil)
	td.check(t, trx3, false, nil)
}

func TestChangePower2(t *testing.T) {
	td := setup(t)

	// This moves proposer to next validator
	updateCommittee(td)

	// Let's create validators first
	pub1, _ := td.RandBLSKeyPair()
	val1 := td.sbx.MakeNewValidator(pub1)
	val1.AddToStake(1)
	val1.UpdateLastBondingHeight(td.sbx.CurrentHeight() - td.sbx.Params().BondInterval)
	td.sbx.UpdateValidator(val1)
	proof1 := td.RandProof()

	pub2, _ := td.RandBLSKeyPair()
	val2 := td.sbx.MakeNewValidator(pub2)
	val2.AddToStake(1)
	val2.UpdateLastBondingHeight(td.sbx.CurrentHeight() - td.sbx.Params().BondInterval)
	td.sbx.UpdateValidator(val2)
	proof2 := td.RandProof()

	pub3, _ := td.RandBLSKeyPair()
	val3 := td.sbx.MakeNewValidator(pub3)
	val3.AddToStake(1)
	val3.UpdateLastBondingHeight(td.sbx.CurrentHeight() - td.sbx.Params().BondInterval)
	td.sbx.UpdateValidator(val3)
	lockTime := td.sbx.CurrentHeight()
	proof3 := td.RandProof()

	val4 := td.sbx.Committee().Proposer(0)
	proof4 := td.RandProof()

	td.sbx.TestParams.CommitteeSize = 7
	td.sbx.TestAcceptSortition = true
	trx1 := tx.NewSortitionTx(lockTime, val1.Address(), proof1)
	td.check(t, trx1, true, nil)
	td.check(t, trx1, false, nil)
	td.execute(t, trx1)

	trx2 := tx.NewSortitionTx(lockTime, val2.Address(), proof2)
	td.check(t, trx2, true, nil)
	td.check(t, trx2, false, nil)
	td.execute(t, trx2)

	trx3 := tx.NewSortitionTx(lockTime, val3.Address(), proof3)
	td.check(t, trx3, true, ErrCommitteeLeaveLimitExceeded)
	td.check(t, trx3, false, nil)

	// Committee member
	trx4 := tx.NewSortitionTx(lockTime, val4.Address(), proof4)
	td.check(t, trx4, true, nil)
	td.check(t, trx4, false, nil)
	td.execute(t, trx4)
}

// TestOldestDidNotPropose tests if the oldest validator in the committee had
// chance to propose a block or not.
func TestOldestDidNotPropose(t *testing.T) {
	td := setup(t)

	// Let's create validators first
	vals := make([]*validator.Validator, 9)
	for index := 0; index < 9; index++ {
		pub, _ := td.RandBLSKeyPair()
		val := td.sbx.MakeNewValidator(pub)
		val.AddToStake(10 * 1e9)
		val.UpdateLastBondingHeight(
			td.sbx.CurrentHeight() - td.sbx.Params().BondInterval)
		td.sbx.UpdateValidator(val)
		vals[index] = val
	}

	td.sbx.TestParams.CommitteeSize = 7
	td.sbx.TestAcceptSortition = true

	// This moves proposer to the next validator
	updateCommittee(td)

	// Let's update committee
	height := td.sbx.CurrentHeight()
	for i := uint32(0); i < 7; i++ {
		height++
		_ = td.sbx.TestStore.AddTestBlock(height)

		lockTime := height
		trx := tx.NewSortitionTx(lockTime, vals[i].Address(), td.RandProof())

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)

		updateCommittee(td)
	}

	height++
	_ = td.sbx.TestStore.AddTestBlock(height)
	lockTime := td.sbx.CurrentHeight()

	trx1 := tx.NewSortitionTx(lockTime, vals[7].Address(), td.RandProof())
	td.check(t, trx1, true, nil)
	td.check(t, trx1, false, nil)
	td.execute(t, trx1)

	trx2 := tx.NewSortitionTx(lockTime, vals[8].Address(), td.RandProof())
	td.check(t, trx2, true, nil)
	td.check(t, trx2, false, nil)
	td.execute(t, trx2)

	updateCommittee(td)

	height++
	_ = td.sbx.TestStore.AddTestBlock(height)
	// Entering validator 16
	trx3 := tx.NewSortitionTx(lockTime+1, vals[8].Address(), td.RandProof())
	td.check(t, trx3, true, ErrOldestValidatorNotProposed)
	td.check(t, trx3, false, nil)
	td.execute(t, trx3)
}
