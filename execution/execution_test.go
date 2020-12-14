package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/validator"
)

var exe *Execution
var sb *sandbox.MockSandbox
var acc1 *account.Account
var val1 *validator.Validator
var priv1 crypto.PrivateKey
var pub1 crypto.PublicKey

func setup() {
	sb = sandbox.NewMockSandbox()
	exe = NewExecution(sb)

	acc1, priv1 = account.GenerateTestAccount(1)
	pub1 = priv1.PublicKey()
	acc1.SetBalance(3000)
	sb.UpdateAccount(acc1)
}

func TestExecuteSendTx(t *testing.T) {
	setup()

	rcvAddr, recPub, rcvPriv := crypto.GenerateTestKeyPair()
	stamp := crypto.GenerateTestHash()
	sb.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, 1, rcvAddr, rcvAddr, 100, 1000, "invalid sender", &recPub, nil)
	trx1.SetSignature(rcvPriv.Sign(trx1.SignBytes()))
	assert.Error(t, exe.Execute(trx1))

	trx2 := tx.NewSendTx(stamp, acc1.Sequence()+2, acc1.Address(), rcvAddr, 1000, 1000, "invalid sequence", &pub1, nil)
	trx2.SetSignature(priv1.Sign(trx2.SignBytes()))
	assert.Error(t, exe.Execute(trx2))

	trx3 := tx.NewSendTx(stamp, acc1.Sequence()+1, acc1.Address(), rcvAddr, 2001, 1000, "insuficent balance", &pub1, nil)
	trx3.SetSignature(priv1.Sign(trx3.SignBytes()))
	assert.Error(t, exe.Execute(trx3))

	trx4 := tx.NewSendTx(stamp, acc1.Sequence()+1, acc1.Address(), rcvAddr, 1000, 999, "invalid fee", &pub1, nil)
	trx4.SetSignature(priv1.Sign(trx4.SignBytes()))
	assert.Error(t, exe.Execute(trx4))

	trx5 := tx.NewSendTx(stamp, acc1.Sequence()+1, acc1.Address(), rcvAddr, 1000, 1000, "ok", &pub1, nil)
	trx5.SetSignature(priv1.Sign(trx5.SignBytes()))
	assert.NoError(t, exe.Execute(trx5))

	// Duplicated. Invalid sequence
	assert.Error(t, exe.Execute(trx5))

	trx6 := tx.NewSendTx(stamp, acc1.Sequence()+1, acc1.Address(), rcvAddr, 1, 1000, "insuficent balance", &pub1, nil)
	trx6.SetSignature(priv1.Sign(trx6.SignBytes()))
	assert.Error(t, exe.Execute(trx6))
	assert.Equal(t, sb.Account(acc1.Address()).Balance(), int64(1000))
	assert.Equal(t, sb.Account(rcvAddr).Balance(), int64(1000))
}

func TestExecuteBondTx(t *testing.T) {
	setup()

	valAddr, valPub, valPriv := crypto.GenerateTestKeyPair()
	stamp := crypto.GenerateTestHash()
	sb.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewBondTx(stamp, 1, valAddr, valPub, 1000, "invalid boner", &valPub, nil)
	trx1.SetSignature(valPriv.Sign(trx1.SignBytes()))
	assert.Error(t, exe.Execute(trx1))

	trx2 := tx.NewBondTx(stamp, acc1.Sequence()+2, acc1.Address(), valPub, 1000, "invalid sequence", &pub1, nil)
	trx2.SetSignature(priv1.Sign(trx2.SignBytes()))
	assert.Error(t, exe.Execute(trx2))

	trx3 := tx.NewBondTx(stamp, acc1.Sequence()+1, acc1.Address(), valPub, 3001, "insuficent balance", &pub1, nil)
	trx3.SetSignature(priv1.Sign(trx3.SignBytes()))
	assert.Error(t, exe.Execute(trx3))

	trx4 := tx.NewBondTx(stamp, acc1.Sequence()+1, acc1.Address(), valPub, 1000, "ok", &pub1, nil)
	trx4.SetSignature(priv1.Sign(trx4.SignBytes()))
	assert.NoError(t, exe.Execute(trx4))

	// Duplicated. Invalid sequence
	assert.Error(t, exe.Execute(trx4))

	assert.Equal(t, sb.Account(acc1.Address()).Balance(), int64(2000))
	assert.Equal(t, sb.Validator(valAddr).Stake(), int64(1000))
}
