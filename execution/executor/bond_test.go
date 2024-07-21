package executor

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/assert"
)

func TestExecuteBondTx(t *testing.T) {
	td := setup(t)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	valPub, _ := td.RandBLSKeyPair()
	receiverAddr := valPub.ValidatorAddress()

	amt := td.RandAmountRange(
		td.sandbox.TestParams.MinimumStake,
		td.sandbox.TestParams.MaximumStake)
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandAccAddress()
		trx := tx.NewBondTx(lockTime, randomAddr,
			receiverAddr, valPub, amt, fee, "unknown address")

		td.check(t, trx, true, AccountNotFoundError{Address: randomAddr})
		td.check(t, trx, false, AccountNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, public key is not set", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, amt, fee, "no public key")

		td.check(t, trx, true, ErrPublicKeyNotSet)
		td.check(t, trx, false, ErrPublicKeyNotSet)
	})

	t.Run("Should fail, public key should not set for existing validators", func(t *testing.T) {
		randPub, _ := td.RandBLSKeyPair()
		val := td.sandbox.MakeNewValidator(randPub)
		td.sandbox.UpdateValidator(val)

		trx := tx.NewBondTx(lockTime, senderAddr,
			randPub.ValidatorAddress(), randPub, amt, fee, "with public key")

		td.check(t, trx, true, ErrPublicKeyAlreadySet)
		td.check(t, trx, false, ErrPublicKeyAlreadySet)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, valPub, senderBalance+1, 0, "insufficient balance")

		td.check(t, trx, true, ErrInsufficientFunds)
		td.check(t, trx, false, ErrInsufficientFunds)
	})

	t.Run("Should fail, unbonded before", func(t *testing.T) {
		randPub, _ := td.RandBLSKeyPair()
		val := td.sandbox.MakeNewValidator(randPub)
		val.UpdateUnbondingHeight(td.RandHeight())
		td.sandbox.UpdateValidator(val)
		trx := tx.NewBondTx(lockTime, senderAddr,
			randPub.ValidatorAddress(), nil, amt, fee, "unbonded before")

		td.check(t, trx, true, ErrValidatorUnbonded)
		td.check(t, trx, false, ErrValidatorUnbonded)
	})

	t.Run("Should fail, amount less than MinimumStake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, valPub, td.sandbox.TestParams.MinimumStake-1, fee, "less than MinimumStake")

		td.check(t, trx, true, SmallStakeError{td.sandbox.TestParams.MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sandbox.TestParams.MinimumStake})
	})

	t.Run("Should fail, validator's stake exceeds the MaximumStake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, valPub, td.sandbox.TestParams.MaximumStake+1, fee, "more than MaximumStake")

		td.check(t, trx, true, MaximumStakeError{td.sandbox.TestParams.MaximumStake})
		td.check(t, trx, false, MaximumStakeError{td.sandbox.TestParams.MaximumStake})
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		pub0 := td.sandbox.Committee().Proposer(0).PublicKey()
		trx := tx.NewBondTx(lockTime, senderAddr,
			pub0.ValidatorAddress(), nil, amt, fee, "inside committee")

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Should fail, joining committee", func(t *testing.T) {
		randPub, _ := td.RandBLSKeyPair()
		val := td.sandbox.MakeNewValidator(randPub)
		td.sandbox.UpdateValidator(val)
		td.sandbox.JoinedToCommittee(val.Address())
		trx := tx.NewBondTx(lockTime, senderAddr,
			randPub.ValidatorAddress(), nil, amt, fee, "inside committee")

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, amt, fee, "ok")

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedSenderAcc := td.sandbox.Account(senderAddr)
	updatedReceiverVal := td.sandbox.Validator(receiverAddr)
	assert.Equal(t, senderBalance-(amt+fee), updatedSenderAcc.Balance())
	assert.Equal(t, amt, updatedReceiverVal.Stake())
	assert.Equal(t, lockTime, updatedReceiverVal.LastBondingHeight())

	td.checkTotalCoin(t, fee)
}

func TestPowerDeltaBond(t *testing.T) {
	td := setup(t)

	senderAddr, _ := td.sandbox.TestStore.RandomTestAcc()
	pub, _ := td.RandBLSKeyPair()
	receiverAddr := pub.ValidatorAddress()
	amt := td.RandAmountRange(
		td.sandbox.TestParams.MinimumStake,
		td.sandbox.TestParams.MaximumStake)
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()
	trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, pub, amt, fee, "ok")

	td.execute(t, trx)

	assert.Equal(t, int64(amt), td.sandbox.PowerDelta())
}

// TestSmallBond tests scenarios involving small and zero stake amounts in bond transactions.
// This test suite is designed to address the issue reported on GitHub:
// https://github.com/pactus-project/pactus/issues/1223
func TestSmallBond(t *testing.T) {
	td := setup(t)

	senderAddr, _ := td.sandbox.TestStore.RandomTestAcc()
	receiverPub, _ := td.RandBLSKeyPair()
	receiverAddr := receiverPub.ValidatorAddress()
	receiverVal := td.sandbox.MakeNewValidator(receiverPub)
	receiverVal.AddToStake(td.sandbox.TestParams.MaximumStake - 2)
	td.sandbox.UpdateValidator(receiverVal)
	lockTime := td.sandbox.CurrentHeight()
	fee := td.RandFee()

	t.Run("Rejects bond transaction with zero amount", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 0, fee, "attacking validator")

		td.check(t, trx, true, SmallStakeError{td.sandbox.TestParams.MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sandbox.TestParams.MinimumStake})
	})

	t.Run("Rejects bond transaction below full validator stake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 1, fee, "attacking validator")

		td.check(t, trx, true, SmallStakeError{td.sandbox.TestParams.MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sandbox.TestParams.MinimumStake})
	})

	t.Run("Accepts bond transaction reaching full validator stake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 2, fee, "fulfilling validator stake")

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	t.Run("Rejects bond transaction with zero amount on full validator", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 0, fee, "attacking validator")

		td.check(t, trx, true, SmallStakeError{td.sandbox.TestParams.MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sandbox.TestParams.MinimumStake})
	})

	receiverValAfterExecution, _ := td.sandbox.TestStore.Validator(receiverVal.Address())
	assert.Equal(t, td.sandbox.Params().MaximumStake, receiverValAfterExecution.Stake())
}
