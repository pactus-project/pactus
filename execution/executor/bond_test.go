package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestExecuteBondTx(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(true)

	senderAddr, senderAcc := tSandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	pub, _ := bls.GenerateTestKeyPair()
	fee, amt := randomAmountAndFee(senderBalance / 2)

	t.Run("Should fail, invalid sender", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, 1, crypto.GenerateTestAddress(),
			pub.Address(), pub, amt, fee, "invalid sender")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, treasury address as receiver", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
			crypto.TreasuryAddress, nil, amt, fee, "invalid ")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	t.Run("Should fail, invalid sequence", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+2, senderAddr,
			pub.Address(), pub, amt, fee, "invalid sequence")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
			pub.Address(), pub, senderBalance+1, 0, "insufficient balance")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInsufficientFunds)
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		pub := tSandbox.Committee().Proposer(0).PublicKey()
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
			pub.Address(), nil, amt, fee, "inside committee")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Should fail, unbonded before", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub)
		val.UpdateUnbondingHeight(tSandbox.CurrentHeight())
		tSandbox.UpdateValidator(val)
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
			pub.Address(), nil, amt, fee, "unbonded before")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	t.Run("Should fail, public key is not set", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
			pub.Address(), nil, amt, fee, "no public key")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
			pub.Address(), pub, amt, fee, "ok")

		assert.NoError(t, exe.Execute(trx, tSandbox), "Ok")
		assert.Error(t, exe.Execute(trx, tSandbox), "Execute again, should fail")
	})

	t.Run("Should fail, public key should not set for existing validators", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+2, senderAddr,
			pub.Address(), pub, amt, fee, "with public key")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	assert.Equal(t, tSandbox.Account(senderAddr).Balance(),
		senderBalance-(amt+fee))
	assert.Equal(t, tSandbox.Validator(pub.Address()).Stake(), amt)
	assert.Equal(t, tSandbox.Validator(pub.Address()).LastBondingHeight(),
		tSandbox.CurrentHeight())
	assert.Equal(t, exe.Fee(), fee)

	checkTotalCoin(t, fee)
}

// TestBondInsideCommittee checks if a validator inside the committee attempts to
// increase their stake.
// In non-strict mode it should be accepted.
func TestBondInsideCommittee(t *testing.T) {
	setup(t)

	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)
	senderAddr, senderAcc := tSandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	fee, amt := randomAmountAndFee(senderBalance)

	pub := tSandbox.Committee().Proposer(0).PublicKey()
	trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
		pub.Address(), nil, amt, fee, "inside committee")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}

// TestBondJoiningCommittee checks if a validator attempts to increase their
// stake after evaluating sortition.
// In non-strict mode, it should be accepted.
func TestBondJoiningCommittee(t *testing.T) {
	setup(t)

	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)
	senderAddr, senderAcc := tSandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	pub, _ := bls.GenerateTestKeyPair()
	fee, amt := randomAmountAndFee(senderBalance)

	val := tSandbox.MakeNewValidator(pub)
	val.UpdateLastSortitionHeight(tHeight500000)
	tSandbox.UpdateValidator(val)
	tSandbox.JoinedToCommittee(val.Address())

	trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
		pub.Address(), nil, amt, fee, "joining committee")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}

// TestStakeExceeded checks if the validator's stake exceeded the MaximumStake parameter.
func TestStakeExceeded(t *testing.T) {
	setup(t)

	exe := NewBondExecutor(true)
	amt := tSandbox.TestParams.MaximumStake + 1
	fee := int64(float64(amt) * tSandbox.Params().FeeFraction)
	senderAddr, senderAcc := tSandbox.TestStore.RandomTestAcc()
	senderAcc.AddToBalance(tSandbox.TestParams.MaximumStake + 1)
	tSandbox.UpdateAccount(senderAddr, senderAcc)
	pub, _ := bls.GenerateTestKeyPair()

	trx := tx.NewBondTx(tStamp500000, senderAcc.Sequence()+1, senderAddr,
		pub.Address(), pub, amt, fee, "stake exceeded")

	assert.Error(t, exe.Execute(trx, tSandbox))
}
