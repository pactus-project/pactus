package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteBondTx(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(true)

	sender := tSandbox.RandomTestAcc()
	senderBalance := sender.Balance()
	pub, _ := bls.GenerateTestKeyPair()
	fee, amt := randomAmountandFee(senderBalance)

	t.Run("Should fail, invalid sender", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, 1, pub.Address(), pub, amt, fee, "invalid sender")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, invalid sequence", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+2, sender.Address(), pub, amt, fee, "invalid sequence")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(), pub, senderBalance+1, 0, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, inside committee", func(t *testing.T) {
		pub := tSandbox.Committee().Proposer(0).PublicKey()
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(), pub, amt, fee, "inside committee")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(), pub, amt, fee, "ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Unbonded before", func(t *testing.T) {
		val := tSandbox.Validator(pub.Address())
		val.UpdateUnbondingHeight(tSandbox.CurHeight)
		tSandbox.UpdateValidator(val)

		trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(), val.PublicKey(), amt, fee, "unbonded before")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, sender.Balance(), senderBalance-(amt+fee))
	assert.Equal(t, tSandbox.Validator(pub.Address()).Stake(), amt)
	assert.Equal(t, tSandbox.Validator(pub.Address()).LastBondingHeight(), tSandbox.CurHeight)
	assert.Equal(t, exe.Fee(), fee)

	checkTotalCoin(t, fee)
}

func TestBondNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)
	sender := tSandbox.RandomTestAcc()

	pub := tSandbox.Committee().Proposer(0).PublicKey()
	trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(), pub, 1000, 1000, "")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}

func TestBondJoiningCommittee(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(true)
	sender := tSandbox.RandomTestAcc()
	senderBalance := sender.Balance()
	pub, _ := bls.GenerateTestKeyPair()
	fee, amt := randomAmountandFee(senderBalance)

	val := tSandbox.MakeNewValidator(pub)
	val.UpdateLastJoinedHeight(tSandbox.CurHeight)
	tSandbox.UpdateValidator(val)

	trx := tx.NewBondTx(tStamp500000, sender.Sequence()+1, sender.Address(), pub, amt, fee, "joining committee")

	assert.Error(t, exe.Execute(trx, tSandbox))
}
