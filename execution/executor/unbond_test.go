package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestExecuteUnbondTx(t *testing.T) {
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
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandValAddress()
		trx := tx.NewUnbondTx(lockTime, randomAddr)

		td.check(t, trx, true, ValidatorNotFoundError{Address: randomAddr})
		td.check(t, trx, false, ValidatorNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		val0 := td.sbx.Committee().Proposer(0)
		trx := tx.NewUnbondTx(lockTime, val0.Address())

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Should fail, joining committee", func(t *testing.T) {
		randPub, _ := td.RandBLSKeyPair()
		randVal := td.sbx.MakeNewValidator(randPub)
		td.sbx.UpdateValidator(randVal)
		td.sbx.JoinedToCommittee(randVal.Address())
		trx := tx.NewUnbondTx(lockTime, randPub.ValidatorAddress())

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, valAddr)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	t.Run("Should fail, Cannot unbond if already unbonded", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, valAddr)

		td.check(t, trx, true, ErrValidatorUnbonded)
		td.check(t, trx, false, ErrValidatorUnbonded)
	})

	updatedVal := td.sbx.Validator(valAddr)

	assert.Equal(t, stake, updatedVal.Stake())
	assert.Zero(t, updatedVal.Power())
	assert.Equal(t, lockTime, updatedVal.UnbondingHeight())
	assert.Equal(t, int64(-stake), td.sbx.PowerDelta())

	td.checkTotalCoin(t, 0)
}

func TestPowerDeltaUnbond(t *testing.T) {
	td := setup(t)

	pub, _ := td.RandBLSKeyPair()
	valAddr := pub.ValidatorAddress()
	val := td.sbx.MakeNewValidator(pub)
	amt := td.RandAmount()
	val.AddToStake(amt)
	td.sbx.UpdateValidator(val)
	lockTime := td.sbx.CurrentHeight()
	trx := tx.NewUnbondTx(lockTime, valAddr)

	td.execute(t, trx)

	assert.Equal(t, int64(-amt), td.sbx.PowerDelta())
}
