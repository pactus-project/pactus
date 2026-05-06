package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/assert"
)

func TestExecuteUnbondTx(t *testing.T) {
	td := setup(t)

	val := td.addTestValidator(t)
	stake := val.Stake()
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandValAddress()
		trx := tx.NewUnbondTx(lockTime, randomAddr)

		td.check(t, trx, true, ValidatorNotFoundError{Address: randomAddr})
		td.check(t, trx, false, ValidatorNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, val.Address())

		td.committee.EXPECT().Contains(val.Address()).Return(true).Times(1)

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Should fail, joining committee", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, val.Address())

		td.committee.EXPECT().Contains(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(val.Address()).Return(true).Times(1)

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, val.Address())

		td.committee.EXPECT().Contains(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(int64(-stake)).Return().Times(1)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	t.Run("Should fail, Cannot unbond if already unbonded", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, val.Address())

		td.check(t, trx, true, ErrValidatorUnbonded)
		td.check(t, trx, false, ErrValidatorUnbonded)
	})

	updatedVal := td.sbx.Validator(val.Address())

	assert.Equal(t, stake, updatedVal.Stake())
	assert.Zero(t, updatedVal.Power())
	assert.Equal(t, td.sbx.CurrentHeight(), updatedVal.UnbondingHeight())

	td.checkTotalCoin(t, 0)
}

func TestPowerDeltaUnbond(t *testing.T) {
	td := setup(t)

	val := td.addTestValidator(t)
	trx := tx.NewUnbondTx(td.RandHeight(), val.Address())

	td.sbx.EXPECT().UpdatePowerDelta(int64(-val.Stake())).Return().Times(1)

	td.execute(t, trx)
}

func TestExecuteDelegatedUnbondTx(t *testing.T) {
	td := setup(t)

	valPub, _ := td.RandBLSKeyPair()
	valAddr := valPub.ValidatorAddress()
	val := td.sbx.MakeNewValidator(valPub)
	val.AddToStake(td.params.MaximumStake)
	owner := td.RandAccAddress()
	val.SetDelegation(owner, amount.Amount(0.2e9), td.sbx.CurrentHeight()+10)
	td.sbx.UpdateValidator(val)
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, invalid delegate owner", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, valAddr)
		pld := trx.Payload().(*payload.UnbondPayload)
		pld.DelegateOwner = td.RandAccAddress()

		td.check(t, trx, true, ErrInvalidDelegateOwner)
		td.check(t, trx, false, ErrInvalidDelegateOwner)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, valAddr)
		pld := trx.Payload().(*payload.UnbondPayload)
		pld.DelegateOwner = owner

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedVal := td.sbx.Validator(valAddr)
	assert.Equal(t, lockTime, updatedVal.UnbondingHeight())
}
