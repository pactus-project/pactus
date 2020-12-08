package txpool

import (
	"strings"
	"testing"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/sandbox"

	"github.com/zarbchain/zarb-go/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/tx"

	"github.com/zarbchain/zarb-go/logger"
)

var pool *txPool
var sb *sandbox.MockSandbox
var acc1 *account.Account
var acc1Pub crypto.PublicKey
var acc1Priv crypto.PrivateKey

func init() {
	logger.InitLogger(logger.DefaultConfig())
	conf := TestConfig()
	p, _ := NewTxPool(conf, nil)
	sb = sandbox.NewMockSandbox()
	addr, pub, priv := crypto.GenerateTestKeyPair()
	acc1 = account.NewAccount(addr)
	acc1.AddToBalance(10000000000)
	sb.UpdateAccount(acc1)
	acc1Priv = priv
	acc1Pub = pub
	p.SetSandbox(sb)
	pool = p.(*txPool)
}

func TestAppendAndRemove(t *testing.T) {
	stamp := crypto.GenerateTestHash()
	sb.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, acc1.Sequence()+1, acc1.Address(), acc1.Address(), 1000, 1000, "acc1->acc1: ok", &acc1Pub, nil)
	trx1.SetSignature(acc1Priv.Sign(trx1.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx1))
	assert.Error(t, pool.appendTx(*trx1))
	pool.RemoveTx(trx1.Hash())
	assert.False(t, pool.HasTx(trx1.Hash()))
}

func TestSendTxValidity(t *testing.T) {
	stamp := crypto.GenerateTestHash()
	senderAddr, senderPub, senderPriv := acc1Pub.Address(), acc1Pub, acc1Priv
	receiverAddr, _, _ := crypto.GenerateTestKeyPair()
	bigMemo := strings.Repeat("a", 1025)

	sb.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 1000, 1000, bigMemo, &senderPub, nil)
	trx1.SetSignature(senderPriv.Sign(trx1.SignBytes()))
	assert.Error(t, pool.appendTx(*trx1))

	trx2 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 1000, 999, "invalid-fee", &senderPub, nil)
	trx2.SetSignature(senderPriv.Sign(trx2.SignBytes()))
	assert.Error(t, pool.appendTx(*trx2))

	trx3 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 10000000, 1000, "invalid-fee", &senderPub, nil)
	trx3.SetSignature(senderPriv.Sign(trx3.SignBytes()))
	assert.Error(t, pool.appendTx(*trx3))

	trx4 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 123456789, 123456, "ok", &senderPub, nil)
	trx4.SetSignature(senderPriv.Sign(trx4.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx4))

	trx5 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 123456000, 123456, "ok", &senderPub, nil)
	trx5.SetSignature(senderPriv.Sign(trx5.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx5))

	trx6 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 1, 1000, "ok", &senderPub, nil)
	trx6.SetSignature(senderPriv.Sign(trx6.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx6))

	trx7 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 1000, 1000, "ok", &senderPub, nil)
	trx7.SetSignature(senderPriv.Sign(trx7.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx7))

	trx8 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 1000000, 1000, "ok", &senderPub, nil)
	trx8.SetSignature(senderPriv.Sign(trx8.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx8))

	trx9 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 10000000, 10000, "ok", &senderPub, nil)
	trx9.SetSignature(senderPriv.Sign(trx9.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx9))

	trx10 := tx.NewSendTx(stamp, acc1.Sequence()+2, senderAddr, receiverAddr, 10000000, 10000, "invalid-sequence", &senderPub, nil)
	trx10.SetSignature(senderPriv.Sign(trx10.SignBytes()))
	assert.Error(t, pool.appendTx(*trx10))

	trx11 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, -10000000, 10000, "negative-amount", &senderPub, nil)
	trx11.SetSignature(senderPriv.Sign(trx11.SignBytes()))
	assert.Error(t, pool.appendTx(*trx11))

	trx12 := tx.NewSendTx(stamp, acc1.Sequence()+1, senderAddr, receiverAddr, 10000000, -10000, "negative-fee", &senderPub, nil)
	trx12.SetSignature(senderPriv.Sign(trx12.SignBytes()))
	assert.Error(t, pool.appendTx(*trx12))
}

func TestStampValidity(t *testing.T) {

	stamp1 := crypto.GenerateTestHash()
	stamp2 := crypto.GenerateTestHash()
	stamp3 := crypto.GenerateTestHash()
	stamp4 := crypto.GenerateTestHash()
	stamp5 := crypto.GenerateTestHash()
	senderAddr, senderPub, senderPriv := acc1Pub.Address(), acc1Pub, acc1Priv
	receiverAddr, _, _ := crypto.GenerateTestKeyPair()

	sb.AppendStampAndUpdateHeight(100, stamp1)
	sb.AppendStampAndUpdateHeight(101, stamp2)
	sb.AppendStampAndUpdateHeight(102, stamp3)
	sb.AppendStampAndUpdateHeight(103, stamp4)

	trx1 := tx.NewSendTx(stamp1, acc1.Sequence()+1, senderAddr, receiverAddr, 1000, 1000, "stamp1-ok", &senderPub, nil)
	trx1.SetSignature(senderPriv.Sign(trx1.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx1))

	sb.AppendStampAndUpdateHeight(104, stamp5)

	trx2 := tx.NewSendTx(stamp1, acc1.Sequence()+1, senderAddr, receiverAddr, 1000, 1000, "stamp1-invalid", &senderPub, nil)
	trx2.SetSignature(senderPriv.Sign(trx2.SignBytes()))
	assert.Error(t, pool.appendTx(*trx2))

	trx3 := tx.NewSendTx(stamp2, acc1.Sequence()+1, senderAddr, receiverAddr, 1000, 1000, "stamp2-ok", &senderPub, nil)
	trx3.SetSignature(senderPriv.Sign(trx3.SignBytes()))
	assert.NoError(t, pool.appendTx(*trx3))
}
