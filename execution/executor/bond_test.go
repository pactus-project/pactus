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
	pub, _ := td.RandomBLSKeyPair()
	receiverAddr := pub.Address()
	amt, fee := td.randomAmountAndFee(td.sandbox.TestParams.MinimumStake, senderBalance)
	t.Run("Should fail, invalid sender", func(t *testing.T) {
		trx := tx.NewBondTx(td.randStamp, 1, td.RandomAddress(),
			receiverAddr, pub, amt, fee, "invalid sender")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, treasury address as receiver", func(t *testing.T) {
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
			crypto.TreasuryAddress, nil, amt, fee, "invalid ")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	t.Run("Should fail, invalid sequence", func(t *testing.T) {
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+2, senderAddr,
			receiverAddr, pub, amt, fee, "invalid sequence")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
			receiverAddr, pub, senderBalance+1, 0, "insufficient balance")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInsufficientFunds)
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		pub := td.sandbox.Committee().Proposer(0).PublicKey()
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
			pub.Address(), nil, amt, fee, "inside committee")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Should fail, unbonded before", func(t *testing.T) {
		pub, _ := td.RandomBLSKeyPair()
		val := td.sandbox.MakeNewValidator(pub)
		val.UpdateUnbondingHeight(td.sandbox.CurrentHeight())
		td.sandbox.UpdateValidator(val)
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
			pub.Address(), nil, amt, fee, "unbonded before")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	t.Run("Should fail, public key is not set", func(t *testing.T) {
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
			receiverAddr, nil, amt, fee, "no public key")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	t.Run("Should fail, amount less than MinimumStake", func(t *testing.T) {
		trx := tx.NewBondTx(td.RandomStamp(), senderAcc.Sequence()+1, senderAddr,
			receiverAddr, pub, 1000, fee, "less than MinimumStake")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.ErrInvalidTx, errors.Code(err))
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
			receiverAddr, pub, amt, fee, "ok")

		err := exe.Execute(trx, td.sandbox)
		assert.NoError(t, err, "Ok")
		err = exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence, "Execute again, should fail")
	})

	t.Run("Should fail, public key should not set for existing validators", func(t *testing.T) {
		trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+2, senderAddr,
			receiverAddr, pub, amt, fee, "with public key")

		err := exe.Execute(trx, td.sandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	assert.Equal(t, td.sandbox.Account(senderAddr).Balance(), senderBalance-(amt+fee))
	assert.Equal(t, td.sandbox.Validator(receiverAddr).Stake(), amt)
	assert.Equal(t, td.sandbox.Validator(receiverAddr).LastBondingHeight(), td.sandbox.CurrentHeight())
	assert.Equal(t, td.sandbox.PowerDelta(), amt)
	assert.Equal(t, exe.Fee(), fee)
	td.checkTotalCoin(t, fee)
}

// TestBondInsideCommittee checks if a validator inside the committee attempts to
// increase their stake.
// In non-strict mode it should be accepted.
func TestBondInsideCommittee(t *testing.T) {
	td := setup(t)

	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)
	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	amt, fee := td.randomAmountAndFee(td.sandbox.TestParams.MinimumStake, senderBalance)

	pub := td.sandbox.Committee().Proposer(0).PublicKey()
	trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
		pub.Address(), nil, amt, fee, "inside committee")

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
	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderBalance := senderAcc.Balance()
	pub, _ := td.RandomBLSKeyPair()
	amt, fee := td.randomAmountAndFee(td.sandbox.TestParams.MinimumStake, senderBalance)

	val := td.sandbox.MakeNewValidator(pub)
	val.UpdateLastSortitionHeight(td.sandbox.CurrentHeight())
	td.sandbox.UpdateValidator(val)
	td.sandbox.JoinedToCommittee(val.Address())

	trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
		pub.Address(), nil, amt, fee, "joining committee")

	assert.Error(t, exe1.Execute(trx, td.sandbox))
	assert.NoError(t, exe2.Execute(trx, td.sandbox))
}

// TestStakeExceeded checks if the validator's stake exceeded the MaximumStake parameter.
func TestStakeExceeded(t *testing.T) {
	td := setup(t)

	exe := NewBondExecutor(true)
	amt := td.sandbox.TestParams.MaximumStake + 1
	fee := int64(float64(amt) * td.sandbox.Params().FeeFraction)
	senderAddr, senderAcc := td.sandbox.TestStore.RandomTestAcc()
	senderAcc.AddToBalance(td.sandbox.TestParams.MaximumStake + 1)
	td.sandbox.UpdateAccount(senderAddr, senderAcc)
	pub, _ := td.RandomBLSKeyPair()

	trx := tx.NewBondTx(td.randStamp, senderAcc.Sequence()+1, senderAddr,
		pub.Address(), pub, amt, fee, "stake exceeded")

	err := exe.Execute(trx, td.sandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidAmount)
}
