package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestExecuteBondTx(t *testing.T) {
	td := setup(t)
	exe := NewBondExecutor(true)

	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	pub, _ := td.RandBLSKeyPair()
	receiverAddr := pub.ValidatorAddress()
	amt := td.RandAmountRange(
		td.sandbox.TestParams.MinimumStake,
		td.sandbox.TestParams.MaximumStake)
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	t.Run("Should fail, invalid sender", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, td.RandAccAddress(),
			receiverAddr, pub, amt, fee, "invalid sender")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidAddress, errors.Code(err))
	})

	t.Run("Should fail, treasury address as receiver", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			crypto.TreasuryAddress, nil, amt, fee, "invalid ")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidPublicKey, errors.Code(err))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, pub, senderBalance+1, 0, "insufficient balance")

		err := exe.Execute(trx, td.sandbox)
		assert.ErrorIs(t, err, ErrInsufficientFunds)
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		pub0 := td.sandbox.Committee().Proposer(0).PublicKey()
		trx := tx.NewBondTx(lockTime, senderAddr,
			pub0.ValidatorAddress(), nil, amt, fee, "inside committee")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidTx, errors.Code(err))
	})

	t.Run("Should fail, unbonded before", func(t *testing.T) {
		unbondedPub, _ := td.RandBLSKeyPair()
		val := td.sandbox.MakeNewValidator(unbondedPub)
		val.UpdateLastBondingHeight(1)
		val.UpdateUnbondingHeight(td.sandbox.CurrentHeight())
		td.sandbox.UpdateValidator(val)
		trx := tx.NewBondTx(lockTime, senderAddr,
			unbondedPub.ValidatorAddress(), nil, amt, fee, "unbonded before")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidHeight, errors.Code(err))
	})

	t.Run("Should fail, public key is not set", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, amt, fee, "no public key")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidPublicKey, errors.Code(err))
	})

	t.Run("Should fail, amount less than MinimumStake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, pub, 1000, fee, "less than MinimumStake")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidTx, errors.Code(err))
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, pub, amt, fee, "ok")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err, "Ok")
	})

	t.Run("Should fail, public key should not set for existing validators", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, pub, amt, fee, "with public key")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidPublicKey, errors.Code(err))
	})

	assert.Equal(t, senderBalance-(amt+fee), td.sandbox.Account(senderAddr).Balance())
	assert.Equal(t, amt, td.sandbox.Validator(receiverAddr).Stake())
	assert.Equal(t, td.sandbox.CurrentHeight(), td.sandbox.Validator(receiverAddr).LastBondingHeight())
	assert.Equal(t, int64(amt), td.sandbox.PowerDelta())
	td.checkTotalCoin(t, fee)
}

// TestBondInsideCommittee checks if a validator inside the committee attempts to
// increase their stake.
// In non-strict mode it should be accepted.
func TestBondInsideCommittee(t *testing.T) {
	td := setup(t)

	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)
	senderAddr, _ := td.sandbox.TestStore.RandomTestAcc()
	amt := td.RandAmountRange(
		td.sandbox.TestParams.MinimumStake,
		td.sandbox.TestParams.MaximumStake-10e9) // it has 10e9 stake
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	pub := td.sandbox.Committee().Proposer(0).PublicKey()
	trx := tx.NewBondTx(lockTime, senderAddr,
		pub.ValidatorAddress(), nil, amt, fee, "inside committee")

	assert.Error(t, exe1.Execute(trx, td.sandbox))
	assert.NoError(t, exe2.Execute(trx, td.sandbox))
}

// TestBondJoiningCommittee checks if a validator attempts to increase their
// stake after evaluating sortition.
// In non-strict mode, it should be accepted.
func TestBondJoiningCommittee(t *testing.T) {
	td := setup(t)

	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)
	senderAddr, _ := td.sandbox.TestStore.RandomTestAcc()
	pub, _ := td.RandBLSKeyPair()
	amt := td.RandAmountRange(
		td.sandbox.TestParams.MinimumStake,
		td.sandbox.TestParams.MaximumStake-10e9) // it has 10e9 stake
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()

	val := td.sandbox.MakeNewValidator(pub)
	val.UpdateLastBondingHeight(1)
	val.UpdateLastSortitionHeight(td.sandbox.CurrentHeight())
	td.sandbox.UpdateValidator(val)
	td.sandbox.JoinedToCommittee(val.Address())

	trx := tx.NewBondTx(lockTime, senderAddr,
		pub.ValidatorAddress(), nil, amt, fee, "joining committee")

	assert.Error(t, exe1.Execute(trx, td.sandbox))
	assert.NoError(t, exe2.Execute(trx, td.sandbox))
}

// TestStakeExceeded checks if the validator's stake exceeded the MaximumStake parameter.
func TestStakeExceeded(t *testing.T) {
	td := setup(t)

	exe := NewBondExecutor(true)
	amt := td.sandbox.TestParams.MaximumStake + 1
	fee := td.RandFee()
	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderAcc.AddToBalance(td.sandbox.TestParams.MaximumStake + 1)
	td.sandbox.UpdateAccount(senderAddr, senderAcc)
	pub, _ := td.RandBLSKeyPair()
	lockTime := td.sandbox.CurrentHeight()

	trx := tx.NewBondTx(lockTime, senderAddr,
		pub.ValidatorAddress(), pub, amt, fee, "stake exceeded")

	err := exe.Execute(trx, td.sandbox)
	assert.Equal(t, errors.ErrInvalidAmount, errors.Code(err))
}

func TestPowerDeltaBond(t *testing.T) {
	td := setup(t)
	exe := NewBondExecutor(true)

	senderAddr, _ := td.sandbox.TestStore.RandomTestAcc()
	pub, _ := td.RandBLSKeyPair()
	receiverAddr := pub.ValidatorAddress()
	amt := td.RandAmountRange(
		td.sandbox.TestParams.MinimumStake,
		td.sandbox.TestParams.MaximumStake)
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()
	trx := tx.NewBondTx(lockTime, senderAddr,
		receiverAddr, pub, amt, fee, "ok")

	err := exe.Execute(trx, td.sandbox)
	assert.NoError(t, err, "Ok")

	assert.Equal(t, int64(amt), td.sandbox.PowerDelta())
}

// TestSmallBond tests scenarios involving small and zero stake amounts in bond transactions.
// This test suite is designed to address the issue reported on GitHub:
// https://github.com/pactus-project/pactus/issues/1223
func TestSmallBond(t *testing.T) {
	td := setup(t)
	exe := NewBondExecutor(false)

	td.sandbox.TestStore.AddTestBlock(752000 + 1) // TODO: remove me in future
	senderAddr, _ := td.sandbox.TestStore.RandomTestAcc()
	receiverVal := td.sandbox.TestStore.RandomTestVal()
	receiverAddr := receiverVal.Address()
	fee := td.RandFee()
	lockTime := td.sandbox.CurrentHeight()
	trxBond := tx.NewBondTx(lockTime, senderAddr,
		receiverAddr, nil, 1000e9-receiverVal.Stake()-2, fee, "ok")

	err := exe.Execute(trxBond, td.sandbox)
	assert.NoError(t, err, "Ok")

	t.Run("Rejects bond transaction with zero amount", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 0, fee, "attacking validator")

		err := exe.Execute(trx, td.sandbox)
		assert.Error(t, err, "Zero bond amount should be rejected")
	})

	t.Run("Rejects bond transaction below full validator stake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 1, fee, "attacking validator")

		err := exe.Execute(trx, td.sandbox)
		assert.Error(t, err, "Bond amount below full stake should be rejected")
	})

	t.Run("Accepts bond transaction reaching full validator stake", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 2, fee, "fulfilling validator stake")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err, "Bond reaching full stake should be accepted")
	})

	t.Run("Accepts bond transaction with zero amount on full validator", func(t *testing.T) {
		trx := tx.NewBondTx(lockTime, senderAddr,
			receiverAddr, nil, 0, fee, "attacking validator")

		err := exe.Execute(trx, td.sandbox)
		assert.Error(t, err, "Zero bond amount on full stake should be rejected")
	})

	val, _ := td.sandbox.TestStore.Validator(receiverVal.Address())
	assert.Equal(t, val.Stake(), td.sandbox.Params().MaximumStake)
}
