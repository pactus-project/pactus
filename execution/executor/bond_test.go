package executor

import (
	"testing"

	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestExecuteBondTx(t *testing.T) {
	td := setup(t)

	senderAcc, senderAddr := td.addTestAccount(t,
		testsuite.AccountWithBalance(10_000e9))
	senderBalance := senderAcc.Balance()
	valPub, _ := td.RandBLSKeyPair()
	receiverAddr := valPub.ValidatorAddress()

	stake := td.randStake()
	fee := td.RandFee()
	lockTime := td.sbx.CurrentHeight()

	t.Run("Should fail, unknown address", func(t *testing.T) {
		randomAddr := td.RandAccAddress()
		trx := tx.NewBondTx(lockTime, randomAddr, receiverAddr, valPub, stake, fee)

		td.check(t, trx, true, AccountNotFoundError{Address: randomAddr})
		td.check(t, trx, false, AccountNotFoundError{Address: randomAddr})
	})

	t.Run("Should fail, public key is not set", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, nil, stake, fee)

		td.check(t, trx, true, ErrPublicKeyNotSet)
		td.check(t, trx, false, ErrPublicKeyNotSet)
	})

	t.Run("Should fail, public key should not set for existing validators", func(t *testing.T) {
		val := td.addTestValidator(t)

		trx := tx.NewBondTx(lockTime, senderAddr, val.Address(), val.PublicKey(), stake, fee)

		td.check(t, trx, true, ErrPublicKeyAlreadySet)
		td.check(t, trx, false, ErrPublicKeyAlreadySet)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, senderBalance+1, 0)

		td.check(t, trx, true, ErrInsufficientFunds)
		td.check(t, trx, false, ErrInsufficientFunds)
	})

	t.Run("Should fail, unbonded before", func(t *testing.T) {
		unbondedVal := td.addTestValidator(t)
		unbondedVal.UpdateUnbondingHeight(td.RandHeight())

		trx := tx.NewBondTx(lockTime, senderAddr, unbondedVal.Address(), nil, stake, fee)

		td.check(t, trx, true, ErrValidatorUnbonded)
		td.check(t, trx, false, ErrValidatorUnbonded)
	})

	t.Run("Should fail, amount less than MinimumStake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, td.sbx.Params().MinimumStake-1, fee)

		td.check(t, trx, true, SmallStakeError{td.sbx.Params().MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sbx.Params().MinimumStake})
	})

	t.Run("Should fail, validator's stake exceeds the MaximumStake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, td.sbx.Params().MaximumStake+1, fee)

		td.check(t, trx, true, MaximumStakeError{td.sbx.Params().MaximumStake})
		td.check(t, trx, false, MaximumStakeError{td.sbx.Params().MaximumStake})
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, stake, fee)

		td.sbx.SbxCommittee.EXPECT().Contains(receiverAddr).Return(true).Times(1)

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Should fail, joining committee", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, stake, fee)

		td.sbx.SbxCommittee.EXPECT().Contains(receiverAddr).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(receiverAddr).Return(true).Times(1)

		td.check(t, trx, true, ErrValidatorInCommittee)
		td.check(t, trx, false, nil)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, stake, fee)

		td.sbx.SbxCommittee.EXPECT().Contains(receiverAddr).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(receiverAddr).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(stake.ToNanoPAC()).Times(1)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedSenderAcc := td.sbx.Account(senderAddr)
	updatedReceiverVal := td.sbx.Validator(receiverAddr)
	assert.Equal(t, senderBalance-(stake+fee), updatedSenderAcc.Balance())
	assert.Equal(t, stake, updatedReceiverVal.Stake())
	assert.Equal(t, lockTime, updatedReceiverVal.LastBondingHeight())

	td.checkTotalCoin(t, fee)
}

func TestPowerDeltaBond(t *testing.T) {
	td := setup(t)

	_, senderAddr := td.addTestAccount(t,
		testsuite.AccountWithBalance(10_000e9))
	pub, _ := td.RandBLSKeyPair()
	receiverAddr := pub.ValidatorAddress()
	stake := td.randStake()
	fee := td.RandFee()
	lockTime := td.sbx.CurrentHeight()
	trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, pub, stake, fee)

	td.sbx.EXPECT().UpdatePowerDelta(stake.ToNanoPAC()).Times(1)

	td.execute(t, trx)
}

// TestSmallBond tests scenarios involving small and zero stake amounts in bond transactions.
// This test suite is designed to address the issue reported on GitHub:
// https://github.com/pactus-project/pactus/issues/1223
func TestSmallBond(t *testing.T) {
	td := setup(t)

	_, senderAddr := td.addTestAccount(t)
	val := td.addTestValidator(t,
		testsuite.ValidatorWithStake(td.sbx.Params().MaximumStake-2))
	receiverAddr := val.Address()
	lockTime := td.sbx.CurrentHeight()
	fee := td.RandFee()

	td.sbx.SbxCommittee.EXPECT().Contains(receiverAddr).Return(false).AnyTimes()
	td.sbx.EXPECT().IsJoinedCommittee(receiverAddr).Return(false).AnyTimes()

	t.Run("Rejects bond transaction with zero amount", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, nil, 0, fee)

		td.check(t, trx, true, SmallStakeError{td.sbx.Params().MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sbx.Params().MinimumStake})
	})

	t.Run("Rejects bond transaction below full validator stake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, nil, 1, fee)

		td.check(t, trx, true, SmallStakeError{td.sbx.Params().MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sbx.Params().MinimumStake})
	})

	t.Run("Accepts bond transaction reaching full validator stake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, nil, 2, fee)

		td.sbx.EXPECT().UpdatePowerDelta(int64(2)).Times(1)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	t.Run("Rejects bond transaction with zero amount on full validator", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, nil, 0, fee)

		td.check(t, trx, true, SmallStakeError{td.sbx.Params().MinimumStake})
		td.check(t, trx, false, SmallStakeError{td.sbx.Params().MinimumStake})
	})

	receiverValAfterExecution := td.sbx.Validator(receiverAddr)
	assert.Equal(t, td.sbx.Params().MaximumStake, receiverValAfterExecution.Stake())
}

func TestExecuteDelegatedBondTx(t *testing.T) {
	td := setup(t)

	senderAcc, senderAddr := td.addTestAccount(t,
		testsuite.AccountWithBalance(10_000e9))
	senderBalance := senderAcc.Balance()
	valPub, _ := td.RandBLSKeyPair()
	receiverAddr := valPub.ValidatorAddress()
	lockTime := td.sbx.CurrentHeight()
	fee := td.RandFee()
	owner := td.RandAccAddress()
	delegateShare := td.RandAmountRange(0, param.MaxDelegateOwnerRewardShare)
	delegateExpiry := td.sbx.CurrentHeight() + 1

	makeDelegatedBond := func(stake amount.Amount) *tx.Tx {
		trx := tx.NewBondTx(lockTime, senderAddr, receiverAddr, valPub, stake, fee)
		pld := trx.Payload().(*payload.BondPayload)
		pld.DelegateOwner = owner
		pld.DelegateShare = delegateShare
		pld.DelegateExpiry = delegateExpiry

		return trx
	}

	t.Run("Should fail, delegation stake must equal maximum", func(t *testing.T) {
		trx := makeDelegatedBond(td.sbx.Params().MaximumStake - 1)

		td.check(t, trx, true, ErrInvalidDelegation)
		td.check(t, trx, false, ErrInvalidDelegation)
	})

	t.Run("Should fail, delegate expiry is in past/current height", func(t *testing.T) {
		trx := makeDelegatedBond(td.sbx.Params().MaximumStake)
		pld := trx.Payload().(*payload.BondPayload)
		pld.DelegateExpiry = td.sbx.CurrentHeight()

		td.check(t, trx, true, ErrDelegateExpiryInPast)
		td.check(t, trx, false, ErrDelegateExpiryInPast)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := makeDelegatedBond(td.sbx.Params().MaximumStake)

		td.sbx.SbxCommittee.EXPECT().Contains(receiverAddr).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(receiverAddr).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(td.sbx.Params().MaximumStake.ToNanoPAC()).Times(1)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
	})

	updatedSenderAcc := td.sbx.Account(senderAddr)
	updatedReceiverVal := td.sbx.Validator(receiverAddr)
	assert.Equal(t, senderBalance-(td.sbx.Params().MaximumStake+fee), updatedSenderAcc.Balance())
	assert.Equal(t, td.sbx.Params().MaximumStake, updatedReceiverVal.Stake())
	assert.Equal(t, owner, updatedReceiverVal.DelegateOwner())
	assert.Equal(t, delegateShare, updatedReceiverVal.DelegateShare())
	assert.Equal(t, delegateExpiry, updatedReceiverVal.DelegateExpiry())
	assert.True(t, updatedReceiverVal.IsDelegated())
}
