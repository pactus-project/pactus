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
var tAcc1Signer, tVal1Signer crypto.Signer
var tAcc1Pub, tVal1Pub crypto.PublicKey
var tTotalCoin int64

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	tSandbox = sandbox.MockingSandbox()

	tTotalCoin = 2100000000000000
	acc0 := account.NewAccount(crypto.TreasuryAddress, 0)
	acc0.AddToBalance(tTotalCoin - 10000000000 - 5000000000)
	tSandbox.UpdateAccount(acc0)

	addr1, pub1, priv1 := crypto.GenerateTestKeyPair()
	_, pub2, priv2 := crypto.GenerateTestKeyPair()

	tAcc1 = account.NewAccount(addr1, 0)
	tAcc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(tAcc1)

	tVal1 = validator.NewValidator(pub2, 0, 0)
	tVal1.AddToStake(5000000000)
	tSandbox.UpdateValidator(tVal1)

	tAcc1Signer = crypto.NewSigner(priv1)
	tVal1Signer = crypto.NewSigner(priv2)
	tAcc1Pub = pub1
	tVal1Pub = pub2
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
	exe := NewSendExecutor(tSandbox)

	addr1, pub1, priv1 := crypto.GenerateTestKeyPair()
	addr2, _, _ := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv1)
	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, 1, addr1, addr2, 3000, 1000, "invalid-sender", &tAcc1Pub, nil)
	signer.SignMsg(trx1)
	assert.Error(t, exe.Execute(trx1))

	trx2 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Signer.Address())+1, tAcc1Signer.Address(), addr1, 3000, 1000, "ok", &tAcc1Pub, nil)
	signer.SignMsg(trx2)
	assert.NoError(t, exe.Execute(trx2))

	trx3 := tx.NewSendTx(stamp, tSandbox.AccSeq(addr1)-2, addr1, addr2, 1000, 1000, "invalid sequence", &pub1, nil)
	signer.SignMsg(trx3)
	assert.Error(t, exe.Execute(trx3))

	trx4 := tx.NewSendTx(stamp, tSandbox.AccSeq(addr1)+1, addr1, addr2, 2001, 1000, "insufficient balance", &pub1, nil)
	signer.SignMsg(trx4)
	assert.Error(t, exe.Execute(trx4))

	trx5 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Signer.Address())+1, tAcc1Signer.Address(), addr1, 3000, 1000, "ok", &tAcc1Pub, nil)
	signer.SignMsg(trx5)
	assert.NoError(t, exe.Execute(trx5))

	signer.SignMsg(trx4)
	assert.NoError(t, exe.Execute(trx4)) // Now has sufficient balance
	assert.Equal(t, exe.Fee(), int64(1000))

	// Duplicated. Invalid sequence
	assert.Error(t, exe.Execute(trx5))

	assert.Equal(t, tSandbox.Account(addr1).Balance(), int64(2999))
	assert.Equal(t, tSandbox.Account(addr2).Balance(), int64(2001))

	checkTotalCoin(t, 3000)
}

func TestSendToSelf(t *testing.T) {
	setup(t)
	exe := NewSendExecutor(tSandbox)

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)
	self := tAcc1Signer.Address()

	bal := tSandbox.Account(self).Balance()
	trx := tx.NewSendTx(stamp, tSandbox.AccSeq(self)+1, self, self, 1000, 1000, "ok", &tAcc1Pub, nil)
	tAcc1Signer.SignMsg(trx)
	assert.NoError(t, exe.Execute(trx))

	assert.Equal(t, tSandbox.Account(self).Balance(), bal-1000) /// He just pay the fee
	assert.Equal(t, exe.Fee(), int64(1000))
}
