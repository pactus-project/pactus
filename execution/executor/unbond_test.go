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

		td.sbx.FakeCommittee.EXPECT().Contains(val.Address()).Return(true).Times(1)

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Should fail, joining committee", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, val.Address())

		td.sbx.FakeCommittee.EXPECT().Contains(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(val.Address()).Return(true).Times(1)

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewUnbondTx(lockTime, val.Address())

		td.sbx.FakeCommittee.EXPECT().Contains(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(-1 * val.Stake().ToNanoPAC()).Return().Times(1)

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

	td.sbx.EXPECT().UpdatePowerDelta(-1 * val.Stake().ToNanoPAC()).Return().Times(1)

	td.execute(t, trx)
}

func TestExecuteDelegatedUnbondTx(t *testing.T) {
	td := setup(t)

	val := td.addTestValidator(t)
	val.AddToStake(td.sbx.Params().MaximumStake)
	owner := td.RandAccAddress()
	expiry := td.sbx.CurrentHeight() + td.RandHeight()
	val.SetDelegation(owner, amount.Amount(0.2e9), expiry)
	td.sbx.UpdateValidator(val)

	// Sanity: the on-chain validator is  delegated.
	assert.True(t, td.sbx.Validator(val.Address()).IsDelegated())

	t.Run("Should fail, invalid delegate owner", func(t *testing.T) {
		trx := tx.NewUnbondTx(td.RandHeight(), val.Address())
		pld := trx.Payload().(*payload.UnbondPayload)
		pld.DelegateOwner = td.RandAccAddress()

		td.check(t, trx, true, ErrInvalidDelegateOwner)
		td.check(t, trx, false, ErrInvalidDelegateOwner)
	})

	t.Run("Should fail, delegation is not expired", func(t *testing.T) {
		td.sbx.FakeHeight = expiry - 1
		trx := tx.NewUnbondTx(td.RandHeight(), val.Address())
		pld := trx.Payload().(*payload.UnbondPayload)
		pld.DelegateOwner = owner

		td.check(t, trx, true, ErrDelegationNotExpired)
		td.check(t, trx, false, ErrDelegationNotExpired)
	})

	t.Run("Ok", func(t *testing.T) {
		td.sbx.FakeHeight = expiry
		trx := tx.NewUnbondTx(td.RandHeight(), val.Address())
		pld := trx.Payload().(*payload.UnbondPayload)
		pld.DelegateOwner = owner

		td.sbx.FakeCommittee.EXPECT().Contains(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(val.Address()).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(-1 * val.Stake().ToNanoPAC()).Return().Times(1)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedVal := td.sbx.Validator(val.Address())
	assert.Equal(t, td.sbx.CurrentHeight(), updatedVal.UnbondingHeight())
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

	val := td.addTestValidator(t)
	val.AddToStake(td.sbx.Params().MaximumStake)
	td.sbx.UpdateValidator(val)
	lockTime := td.sbx.CurrentHeight()

	// Sanity: the on-chain validator is not delegated.
	assert.False(t, td.sbx.Validator(val.Address()).IsDelegated())

	trx := tx.NewUnbondTx(lockTime, val.Address())
	pld := trx.Payload().(*payload.UnbondPayload)
	// Attacker points DelegateOwner at their own account.
	pld.DelegateOwner = td.RandAccAddress()

	// The payload now reports itself as delegated, so its Signer() is the
	// attacker, but the on-chain validator is not delegated.
	assert.True(t, pld.IsDelegated())

	td.check(t, trx, true, ErrInvalidDelegateOwner)
	td.check(t, trx, false, ErrInvalidDelegateOwner)
}
