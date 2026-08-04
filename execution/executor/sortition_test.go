package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestExecuteSortitionTx(t *testing.T) {
	td := setup(t)

	lockTime := td.sbx.CurrentHeight()
	proof := td.RandProof()
	val := td.addTestValidator(t)

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandValAddress()
		trx := tx.NewSortitionTx(lockTime, randomAddr, proof)

		td.check(t, trx, true, ValidatorNotFoundError{Address: randomAddr})
		td.check(t, trx, false, ValidatorNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, Bonding period", func(t *testing.T) {
		val.UpdateLastBondingHeight(td.sbx.CurrentHeight().SafeDecrease(td.sbx.Params().BondInterval - 1))
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)

		td.check(t, trx, true, ErrBondingPeriod)
		td.check(t, trx, false, ErrBondingPeriod)
	})

	val.UpdateLastBondingHeight(td.sbx.CurrentHeight().SafeDecrease(td.sbx.Params().BondInterval))

	t.Run("Should fail, invalid proof", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)

		td.sbx.EXPECT().VerifyProof(lockTime, proof, val).Return(false).Times(2)

		td.check(t, trx, true, ErrInvalidSortitionProof)
		td.check(t, trx, false, ErrInvalidSortitionProof)
	})

	t.Run("Should fail, committee has free seats and validator is in the committee", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)

		td.sbx.EXPECT().VerifyProof(lockTime, proof, val).Return(true).Times(2)
		td.sbx.FakeCommittee.EXPECT().Contains(val.Address()).Return(true).Times(1)
		td.sbx.FakeCommittee.EXPECT().Size().Return(td.sbx.Params().CommitteeSize - 1).Times(1)

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)

		td.sbx.EXPECT().VerifyProof(lockTime, proof, val).Return(true).Times(2)
		td.sbx.FakeCommittee.EXPECT().Contains(val.Address()).Return(false).Times(1)
		td.sbx.FakeCommittee.EXPECT().Size().Return(td.sbx.Params().CommitteeSize - 1).Times(1)
		td.sbx.EXPECT().JoinToCommittee(val.Address()).Return().Times(1)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	t.Run("Should fail, expired sortition", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime-1, val.Address(), proof)

		td.sbx.EXPECT().VerifyProof(lockTime-1, proof, val).Return(true).Times(2)

		td.check(t, trx, true, ErrExpiredSortition)
		td.check(t, trx, false, ErrExpiredSortition)
	})

	t.Run("Should fail, duplicated sortition", func(t *testing.T) {
		trx := tx.NewSortitionTx(lockTime, val.Address(), proof)

		td.sbx.EXPECT().VerifyProof(lockTime, proof, val).Return(true).Times(2)

		td.check(t, trx, true, ErrExpiredSortition)
		td.check(t, trx, false, ErrExpiredSortition)
	})

	updatedVal := td.sbx.Validator(val.Address())

	assert.Equal(t, lockTime, updatedVal.LastSortitionHeight())

	td.checkTotalCoin(t, 0)
}

func TestChangePower(t *testing.T) {
	td := setup(t)

	lockTime := types.Height(1011)
	proof := td.RandProof()
	td.sbx.EXPECT().VerifyProof(lockTime, proof, gomock.Any()).Return(true).AnyTimes()

	// Existing validator in the committee.
	vals := make([]*validator.Validator, 7)
	vals[0] = td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(1001))
	vals[1] = td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(999))
	vals[2] = td.addTestValidator(t, testsuite.ValidatorWithStake(2), testsuite.ValidatorWithSortitionHeight(1002))
	vals[3] = td.addTestValidator(t, testsuite.ValidatorWithStake(3), testsuite.ValidatorWithSortitionHeight(997))
	vals[4] = td.addTestValidator(t, testsuite.ValidatorWithStake(3), testsuite.ValidatorWithSortitionHeight(1004))
	vals[5] = td.addTestValidator(t, testsuite.ValidatorWithStake(3), testsuite.ValidatorWithSortitionHeight(1007))
	vals[6] = td.addTestValidator(t, testsuite.ValidatorWithStake(2), testsuite.ValidatorWithSortitionHeight(1009))

	td.sbx.FakeParams.CommitteeSize = len(vals)
	td.sbx.FakeCommittee.EXPECT().Size().Return(td.sbx.Params().CommitteeSize).AnyTimes()
	td.sbx.FakeCommittee.EXPECT().Validators().DoAndReturn(
		func() []*validator.Validator {
			vals2 := make([]*validator.Validator, 7)
			for i, v := range vals {
				vals2[i] = v.Clone()
			}

			return vals2
		},
	).AnyTimes()
	td.sbx.FakeCommittee.EXPECT().Contains(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) bool {
			for _, v := range vals {
				if v.Address() == addr {
					return true
				}
			}

			return false
		},
	).AnyTimes()

	td.sbx.FakeCommittee.EXPECT().Power().DoAndReturn(
		func() int64 {
			power := int64(0)
			for _, v := range vals {
				power += v.Power()
			}

			return power
		},
	).AnyTimes()

	t.Run("join power exceeds 1/3 threshold", func(t *testing.T) {
		// New validator attempting to join the committee.
		newVal1 := td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(lockTime-3))
		newVal2 := td.addTestValidator(t, testsuite.ValidatorWithStake(3), testsuite.ValidatorWithSortitionHeight(lockTime-2))
		newVal3 := td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(lockTime-1))

		td.sbx.EXPECT().IterateValidators(gomock.Any()).DoAndReturn(func(consume func(*validator.Validator, bool, bool)) {
			consume(vals[0], false, false)
			consume(vals[1], false, true)
			consume(newVal1, false, true)
			consume(newVal2, false, true)
		}).Times(1)

		trx1 := tx.NewSortitionTx(lockTime, newVal3.Address(), proof)

		td.check(t, trx1, true, ErrCommitteeJoinLimitExceeded)
		td.check(t, trx1, false, nil)
	})

	t.Run("leaving power exceeds 1/3 threshold", func(t *testing.T) {
		// New validator attempting to join the committee.
		newVal1 := td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(lockTime-3))
		newVal2 := td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(lockTime-2))
		newVal3 := td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(lockTime-1))

		td.sbx.EXPECT().IterateValidators(gomock.Any()).DoAndReturn(func(consume func(*validator.Validator, bool, bool)) {
			consume(vals[0], false, false)
			consume(vals[1], false, true)
			consume(newVal1, false, true)
			consume(newVal2, false, true)
		}).Times(1)

		trx1 := tx.NewSortitionTx(lockTime, newVal3.Address(), proof)

		td.check(t, trx1, true, ErrCommitteeLeaveLimitExceeded)
		td.check(t, trx1, false, nil)
	})

	t.Run("oldest validator has not proposed", func(t *testing.T) {
		// New validator attempting to join the committee.
		newVal := td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(lockTime-1))

		td.sbx.EXPECT().IterateValidators(gomock.Any()).DoAndReturn(func(func(*validator.Validator, bool, bool)) {
		}).Times(1)

		td.sbx.FakeCommittee.EXPECT().Proposer(types.Round(0)).Return(vals[3]).Times(1)
		trx1 := tx.NewSortitionTx(lockTime, newVal.Address(), proof)

		td.check(t, trx1, true, ErrOldestValidatorNotProposed)
		td.check(t, trx1, false, nil)
	})

	t.Run("no error when validator can join committee safely", func(t *testing.T) {
		// New validator attempting to join the committee.
		newVal := td.addTestValidator(t, testsuite.ValidatorWithStake(1), testsuite.ValidatorWithSortitionHeight(lockTime-1))

		td.sbx.EXPECT().IterateValidators(gomock.Any()).DoAndReturn(func(func(*validator.Validator, bool, bool)) {
		}).Times(1)
		td.sbx.FakeCommittee.EXPECT().Proposer(types.Round(0)).Return(vals[0]).Times(1)

		trx1 := tx.NewSortitionTx(lockTime, newVal.Address(), proof)

		td.check(t, trx1, true, nil)
		td.check(t, trx1, false, nil)

		td.sbx.EXPECT().JoinToCommittee(newVal.Address()).Return().Times(1)
		td.execute(t, trx1)
	})
}
