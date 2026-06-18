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

	bonderAddr, bonderAcc := td.sbx.TestStore.RandomTestAcc()
	bonderBalance := bonderAcc.Balance()
	stake := td.RandAmountRange(
		td.sbx.TestParams.MinimumStake,
		bonderBalance,
	)
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

func TestExecuteDelegatedUnbondTx(t *testing.T) {
	td := setup(t)

	valPub, _ := td.RandBLSKeyPair()
	valAddr := valPub.ValidatorAddress()
	val := td.sbx.MakeNewValidator(valPub)
	val.AddToStake(td.sbx.TestParams.MaximumStake)
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

// TestUnbondNonDelegatedWithDelegateOwner is a regression test for an
// authorization-bypass vulnerability. The on-chain validator is NOT delegated,
// but the attacker crafts an unbond payload with a non-Treasury DelegateOwner
// set to their own account. Because the payload reports itself as delegated,
// its Signer() becomes the attacker's account, allowing the attacker to sign
// the transaction with their own key. Before the fix, Check() skipped the
// owner-match validation for non-delegated validators, letting the attacker
// force-unbond a validator they do not own. Check() must reject this mismatch.
func TestUnbondNonDelegatedWithDelegateOwner(t *testing.T) {
	td := setup(t)

	valPub, _ := td.RandBLSKeyPair()
	valAddr := valPub.ValidatorAddress()
	val := td.sbx.MakeNewValidator(valPub)
	val.AddToStake(td.sbx.TestParams.MaximumStake)
	td.sbx.UpdateValidator(val)
	lockTime := td.sbx.CurrentHeight()

	// Sanity: the on-chain validator is not delegated.
	assert.False(t, td.sbx.Validator(valAddr).IsDelegated())

	t.Run("Should fail, attacker sets DelegateOwner on a non-delegated validator", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, valAddr)
		pld := trx.Payload().(*payload.UnbondPayload)
		// Attacker points DelegateOwner at their own account.
		pld.DelegateOwner = td.RandAccAddress()

		// The payload now reports itself as delegated, so its Signer() is the
		// attacker, but the on-chain validator is not delegated.
		assert.True(t, pld.IsDelegated())

		td.check(t, trx, true, ErrInvalidDelegateOwner)
		td.check(t, trx, false, ErrInvalidDelegateOwner)
	})

	t.Run("Ok, self-unbond with no DelegateOwner still works", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, valAddr)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
	})
}
