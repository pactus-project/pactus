package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

var tSandbox *sandbox.MockSandbox
var tVal1 *validator.Validator
var tAcc1 *account.Account
var tValSigner crypto.Signer
var tTotalCoin int64

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	tSandbox = sandbox.MockingSandbox()

	tTotalCoin = 21 * 1e14
	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(tTotalCoin - 10000000000 - 5000000000)
	tSandbox.UpdateAccount(acc0)

	signer1 := crypto.GenerateTestSigner()
	tValSigner = crypto.GenerateTestSigner()

	tAcc1 = account.NewAccount(signer1.Address(), 0)
	tAcc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(tAcc1)

	tVal1 = validator.NewValidator(tValSigner.PublicKey(), 0, 0)
	tVal1.AddToStake(5000000000)
	tSandbox.UpdateValidator(tVal1)
}

func checkTotalCoin(t *testing.T, fee int64) {
	total := int64(0)
	for _, acc := range tSandbox.Accounts {
		total += acc.Balance()
	}
	for _, val := range tSandbox.Validators {
		total += val.Stake()
	}
	assert.Equal(t, total+fee, tTotalCoin)
}

func TestExecuteSendTx(t *testing.T) {
	setup(t)
	exe := NewSendExecutor(tSandbox, true)

	sender := crypto.GenerateTestSigner()
	receiver := crypto.GenerateTestSigner()
	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewSendTx(stamp, 1, sender.Address(), sender.Address(), 3000, 1000, "non-existing account")
		sender.SignMsg(trx)

		assert.Error(t, exe.Execute(trx))
	})

	t.Run("ok. Create sender account", func(t *testing.T) {
		trx := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1.Address())+1, tAcc1.Address(), sender.Address(), 3000, 1000, "ok")

		assert.NoError(t, exe.Execute(trx))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSendTx(stamp, 2, sender.Address(), receiver.Address(), 1000, 1000, "invalid sequence")

		assert.Error(t, exe.Execute(trx))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewSendTx(stamp, 1, sender.Address(), receiver.Address(), 2001, 1000, "insufficient balance")

		assert.Error(t, exe.Execute(trx))
	})

	t.Run("Ok", func(t *testing.T) {
		trx1 := tx.NewSendTx(stamp, 1, sender.Address(), receiver.Address(), 700, 1000, "ok")
		assert.NoError(t, exe.Execute(trx1))

		trx2 := tx.NewSendTx(stamp, 2, sender.Address(), receiver.Address(), 300, 1000, "ok")
		assert.NoError(t, exe.Execute(trx2))

		// Replay transactions
		assert.Error(t, exe.Execute(trx1))
		assert.Error(t, exe.Execute(trx2))
	})

	t.Run("Send to self", func(t *testing.T) {
		self := tAcc1.Address()
		bal := tSandbox.Account(self).Balance()
		trx := tx.NewSendTx(stamp, tSandbox.AccSeq(self)+1, self, self, 1000, 1000, "ok")
		assert.NoError(t, exe.Execute(trx))

		assert.Equal(t, tSandbox.Account(self).Balance(), bal-1000) /// Fee should be deducted
		assert.Equal(t, exe.Fee(), int64(1000))
	})

	assert.Equal(t, tSandbox.Account(sender.Address()).Balance(), int64(0))
	assert.Equal(t, tSandbox.Account(receiver.Address()).Balance(), int64(1000))

	checkTotalCoin(t, 4000)
}

func TestSendNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewSendExecutor(tSandbox, false)
	exe2 := NewSendExecutor(tSandbox, true)

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)
	receiver1, _, _ := crypto.GenerateTestKeyPair()
	receiver2, _, _ := crypto.GenerateTestKeyPair()

	mintbase1 := tx.NewMintbaseTx(stamp, 101, receiver1, 5, "")
	mintbase2 := tx.NewMintbaseTx(stamp, 101, receiver2, 5, "")
	mintbase3 := tx.NewMintbaseTx(stamp, 102, receiver2, 5, "")

	assert.NoError(t, exe1.Execute(mintbase1))
	assert.NoError(t, exe1.Execute(mintbase2))
	assert.Error(t, exe1.Execute(mintbase3)) // Invalid sequence

	assert.Error(t, exe2.Execute(mintbase1))
	assert.Error(t, exe2.Execute(mintbase2))
}
