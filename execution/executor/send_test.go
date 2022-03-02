package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

var (
	tSandbox     *sandbox.MockSandbox
	tVal1        *validator.Validator
	tAcc1        *account.Account
	tTotalCoin   int64
	tAcc1Balance int64
	tVal1Stake   int64
	tHash500000  hash.Hash
	tHash500001  hash.Hash
)

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	tSandbox = sandbox.MockingSandbox()

	tTotalCoin = int64(21000000 * 1e8)
	tAcc1Balance = int64(1000 * 1e8)
	tVal1Stake = int64(5000 * 1e8)
	treasuryAmt := tTotalCoin - (tAcc1Balance + tVal1Stake)
	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(treasuryAmt)
	tSandbox.UpdateAccount(acc0)

	addr := crypto.GenerateTestAddress()
	tAcc1 = account.NewAccount(addr, 0)
	tAcc1.AddToBalance(tAcc1Balance)
	tSandbox.UpdateAccount(tAcc1)
	assert.Equal(t, tSandbox.Account(tAcc1.Address()).Balance(), tAcc1Balance)

	pub, _ := bls.GenerateTestKeyPair()
	tVal1 = validator.NewValidator(pub, 0)
	tVal1.AddToStake(tVal1Stake)
	tVal1.UpdateLastBondingHeight(100001 - tSandbox.BondInterval()) // Check TestExecuteSortitionTx
	tSandbox.UpdateValidator(tVal1)
	assert.Equal(t, tSandbox.Validator(tVal1.Address()).Stake(), tVal1Stake)

	tHash500000 = hash.GenerateTestHash()
	tHash500001 = hash.GenerateTestHash()
	tSandbox.AppendNewBlock(500000, tHash500000)
	tSandbox.AppendNewBlock(500001, tHash500001)

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
	exe := NewSendExecutor(true)

	sender := bls.GenerateTestSigner()
	receiver := bls.GenerateTestSigner()

	t.Run("Should fail, Sender has no account", func(t *testing.T) {
		trx := tx.NewSendTx(tHash500000.Stamp(), 1, sender.Address(), sender.Address(), 3000, 1000, "non-existing account")
		sender.SignMsg(trx)

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("ok. Create sender account", func(t *testing.T) {
		trx := tx.NewSendTx(tHash500000.Stamp(), tSandbox.AccSeq(tAcc1.Address())+1, tAcc1.Address(), sender.Address(), 3000, 1000, "ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSendTx(tHash500000.Stamp(), 2, sender.Address(), receiver.Address(), 1000, 1000, "invalid sequence")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, insufficient balance", func(t *testing.T) {
		trx := tx.NewSendTx(tHash500000.Stamp(), 1, sender.Address(), receiver.Address(), 2001, 1000, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Ok", func(t *testing.T) {
		trx1 := tx.NewSendTx(tHash500000.Stamp(), 1, sender.Address(), receiver.Address(), 700, 1000, "ok")
		assert.NoError(t, exe.Execute(trx1, tSandbox))

		trx2 := tx.NewSendTx(tHash500000.Stamp(), 2, sender.Address(), receiver.Address(), 300, 1000, "ok")
		assert.NoError(t, exe.Execute(trx2, tSandbox))

		// Replay transactions
		assert.Error(t, exe.Execute(trx1, tSandbox))
		assert.Error(t, exe.Execute(trx2, tSandbox))
	})

	t.Run("Send to self", func(t *testing.T) {
		self := tAcc1.Address()
		bal := tSandbox.Account(self).Balance()
		trx := tx.NewSendTx(tHash500000.Stamp(), tSandbox.AccSeq(self)+1, self, self, 1000, 1000, "ok")
		assert.NoError(t, exe.Execute(trx, tSandbox))

		assert.Equal(t, tSandbox.Account(self).Balance(), bal-1000) /// Fee should be deducted
		assert.Equal(t, exe.Fee(), int64(1000))
	})

	assert.Equal(t, tSandbox.Account(sender.Address()).Balance(), int64(0))
	assert.Equal(t, tSandbox.Account(receiver.Address()).Balance(), int64(1000))

	checkTotalCoin(t, 4000)
}

func TestSendNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewSendExecutor(true)
	exe2 := NewSendExecutor(false)

	receiver1 := crypto.GenerateTestAddress()

	trx1 := tx.NewMintbaseTx(tHash500001.Stamp(), tSandbox.CurHeight, receiver1, 1, "")
	assert.Error(t, exe1.Execute(trx1, tSandbox)) // Invalid sequence
	assert.NoError(t, exe2.Execute(trx1, tSandbox))

	trx2 := tx.NewMintbaseTx(tHash500001.Stamp(), tSandbox.CurHeight+1, receiver1, 1, "")
	assert.Error(t, exe1.Execute(trx2, tSandbox)) // Invalid height
	assert.Error(t, exe2.Execute(trx2, tSandbox)) // Invalid height

}
