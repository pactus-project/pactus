package execution

import (
	"testing"

	"github.com/zarbchain/zarb-go/sortition"

	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/validator"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sandbox"
)

var tExec *Execution
var tAcc1 *account.Account
var tVal1 *validator.Validator
var tPriv1 crypto.PrivateKey
var tPub1 crypto.PublicKey
var tSandbox *sandbox.MockSandbox

func setup(t *testing.T) {
	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)

	tSandbox = sandbox.NewMockSandbox()

	tAcc1, tPriv1 = account.GenerateTestAccount(1)
	tPub1 = tPriv1.PublicKey()
	tAcc1.SubtractFromBalance(tAcc1.Balance()) // make balance zero
	tAcc1.AddToBalance(3000)
	tSandbox.UpdateAccount(tAcc1)
	tVal1 = validator.NewValidator(tPub1, 0, 0)
	tSandbox.UpdateValidator(tVal1)

	tExec = NewExecution(tSandbox)
}

func TestExecuteSendTx(t *testing.T) {
	setup(t)

	rcvAddr, recPub, rcvPriv := crypto.GenerateTestKeyPair()
	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, 1, rcvAddr, rcvAddr, 100, 1000, "invalid sender", &recPub, nil)
	trx1.SetSignature(rcvPriv.Sign(trx1.SignBytes()))
	assert.Error(t, tExec.Execute(trx1))

	trx2 := tx.NewSendTx(stamp, tAcc1.Sequence()+2, tAcc1.Address(), rcvAddr, 1000, 1000, "invalid sequence", &tPub1, nil)
	trx2.SetSignature(tPriv1.Sign(trx2.SignBytes()))
	assert.Error(t, tExec.Execute(trx2))

	trx3 := tx.NewSendTx(stamp, tAcc1.Sequence()+1, tAcc1.Address(), rcvAddr, 2001, 1000, "insufficient balance", &tPub1, nil)
	trx3.SetSignature(tPriv1.Sign(trx3.SignBytes()))
	assert.Error(t, tExec.Execute(trx3))

	trx4 := tx.NewSendTx(stamp, tAcc1.Sequence()+1, tAcc1.Address(), rcvAddr, 1000, 999, "invalid fee", &tPub1, nil)
	trx4.SetSignature(tPriv1.Sign(trx4.SignBytes()))
	assert.Error(t, tExec.Execute(trx4))

	trx5 := tx.NewSendTx(stamp, tAcc1.Sequence()+1, tAcc1.Address(), rcvAddr, 1000, 1000, "ok", &tPub1, nil)
	trx5.SetSignature(tPriv1.Sign(trx5.SignBytes()))
	assert.NoError(t, tExec.Execute(trx5))

	// Duplicated. Invalid sequence
	assert.Error(t, tExec.Execute(trx5))

	trx6 := tx.NewSendTx(stamp, tAcc1.Sequence()+1, tAcc1.Address(), rcvAddr, 1, 1000, "insufficient balance", &tPub1, nil)
	trx6.SetSignature(tPriv1.Sign(trx6.SignBytes()))
	assert.Error(t, tExec.Execute(trx6))
	assert.Equal(t, tSandbox.Account(tAcc1.Address()).Balance(), int64(1000))
	assert.Equal(t, tSandbox.Account(rcvAddr).Balance(), int64(1000))
}

func TestExecuteBondTx(t *testing.T) {
	setup(t)

	valAddr, valPub, valPriv := crypto.GenerateTestKeyPair()
	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewBondTx(stamp, 1, valAddr, valPub, 1000, "invalid boner", &valPub, nil)
	trx1.SetSignature(valPriv.Sign(trx1.SignBytes()))
	assert.Error(t, tExec.Execute(trx1))

	trx2 := tx.NewBondTx(stamp, tAcc1.Sequence()+2, tAcc1.Address(), valPub, 1000, "invalid sequence", &tPub1, nil)
	trx2.SetSignature(tPriv1.Sign(trx2.SignBytes()))
	assert.Error(t, tExec.Execute(trx2))

	trx3 := tx.NewBondTx(stamp, tAcc1.Sequence()+1, tAcc1.Address(), valPub, 3001, "insufficient balance", &tPub1, nil)
	trx3.SetSignature(tPriv1.Sign(trx3.SignBytes()))
	assert.Error(t, tExec.Execute(trx3))

	trx4 := tx.NewBondTx(stamp, tAcc1.Sequence()+1, tAcc1.Address(), valPub, 1000, "ok", &tPub1, nil)
	trx4.SetSignature(tPriv1.Sign(trx4.SignBytes()))
	assert.NoError(t, tExec.Execute(trx4))

	// Duplicated. Invalid sequence
	assert.Error(t, tExec.Execute(trx4))

	assert.Equal(t, tSandbox.Account(tAcc1.Address()).Balance(), int64(2000))
	assert.Equal(t, tSandbox.Validator(valAddr).Stake(), int64(1000))
}

func TestExecuteSortitionTx(t *testing.T) {
	setup(t)

	valAddr, valPub, valPriv := crypto.GenerateTestKeyPair()
	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)
	proof := [48]byte{}

	trx1 := tx.NewSortitionTx(stamp, 1, valAddr, proof[:], "invalid validator", &valPub, nil)
	trx1.SetSignature(valPriv.Sign(trx1.SignBytes()))
	assert.Error(t, tExec.Execute(trx1))

	val := validator.NewValidator(valPub, 0, 0)
	tSandbox.UpdateValidator(val)

	trx2 := tx.NewSortitionTx(stamp, 1, valAddr, proof[:], "invalid proof", &valPub, nil)
	trx2.SetSignature(valPriv.Sign(trx2.SignBytes()))
	assert.Error(t, tExec.Execute(trx2))

	sortition := sortition.NewSortition(crypto.NewSigner(valPriv))
	trx3 := sortition.EvaluateTransaction(stamp, val)
	assert.NotNil(t, trx3)
	assert.NoError(t, tExec.Execute(trx3))

}
