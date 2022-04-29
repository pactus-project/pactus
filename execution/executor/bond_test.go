package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestExecuteBondTx(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(true)

	sender := tSandbox.TestStore.RandomTestAcc()
	senderBalance := sender.Balance()
	pub, _ := bls.GenerateTestKeyPair()
	fee, amt := randomAmountAndFee(senderBalance)

	t.Run("Should fail, invalid sender", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, 1, crypto.GenerateTestAddress(),
			pub.Address(), pub, amt, fee, "invalid sender")
		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
	})

	t.Run("Should fail, invalid sequence", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+2, sender.Address(),
			pub.Address(), pub, amt, fee, "invalid sequence")
		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(),
			pub.Address(), pub, senderBalance+1, 0, "insufficient balance")
		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInsufficientFunds)
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		pub := tSandbox.Committee().Proposer(0).PublicKey()
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(),
			pub.Address(), pub, amt, fee, "inside committee")
		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidTx)
	})

	t.Run("Unbonded before", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub)
		val.UpdateUnbondingHeight(tSandbox.CurrentHeight())
		tSandbox.UpdateValidator(val)

		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(),
			pub.Address(), pub, amt, fee, "unbonded before")
		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
	})

	t.Run("public key is not set", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(),
			pub.Address(), nil, amt, fee, "ok")

		err := exe.Execute(trx, tSandbox)
		assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(),
			pub.Address(), pub, amt, fee, "ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Account(sender.Address()).Balance(),
		senderBalance-(amt+fee))
	assert.Equal(t, tSandbox.Validator(pub.Address()).Stake(), amt)
	assert.Equal(t, tSandbox.Validator(pub.Address()).LastBondingHeight(),
		tSandbox.CurrentHeight())
	assert.Equal(t, exe.Fee(), fee)

	checkTotalCoin(t, fee)
}

func TestBondNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)
	sender := tSandbox.TestStore.RandomTestAcc()

	pub := tSandbox.Committee().Proposer(0).PublicKey()
	trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(),
		pub.Address(), pub, 1000, 1000, "")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}

func TestBondJoiningCommittee(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(true)
	sender := tSandbox.TestStore.RandomTestAcc()
	senderBalance := sender.Balance()
	pub, _ := bls.GenerateTestKeyPair()
	fee, amt := randomAmountAndFee(senderBalance)

	val := tSandbox.MakeNewValidator(pub)
	val.UpdateLastJoinedHeight(tSandbox.CurrentHeight())
	tSandbox.UpdateValidator(val)

	trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(),
		pub.Address(), pub, amt, fee, "joining committee")
	err := exe.Execute(trx, tSandbox)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidHeight)
}
