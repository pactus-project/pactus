package txpool

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

func TestValidity(t *testing.T) {
	sb := sandbox.NewMockSandbox()
	pool.SetSandbox(sb)

	stamp := crypto.GenerateTestHash()
	stampInv := crypto.GenerateTestHash()
	senderAddr, senderPub, senderPriv := crypto.GenerateTestKeyPair()
	receiverAddr, _, _ := crypto.GenerateTestKeyPair()
	bigMemo := strings.Repeat("a", 1025)

	sb.TTLInterval = 2
	sb.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000, 1000, bigMemo, &senderPub, nil)
	trx1.SetSignature(senderPriv.Sign(trx1.SignBytes()))
	assert.Error(t, pool.validateTx(trx1))

	trx2 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000, 999, "invalid-fee", &senderPub, nil)
	trx2.SetSignature(senderPriv.Sign(trx2.SignBytes()))
	assert.Error(t, pool.validateTx(trx2))

	trx3 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 10000000, 1000, "invalid-fee", &senderPub, nil)
	trx3.SetSignature(senderPriv.Sign(trx3.SignBytes()))
	assert.Error(t, pool.validateTx(trx3))

	trx4 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 123456789, 123456, "ok", &senderPub, nil)
	trx4.SetSignature(senderPriv.Sign(trx4.SignBytes()))
	assert.NoError(t, pool.validateTx(trx4))

	trx5 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 123456000, 123456, "ok", &senderPub, nil)
	trx5.SetSignature(senderPriv.Sign(trx5.SignBytes()))
	assert.NoError(t, pool.validateTx(trx5))

	trx6 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1, 1000, "ok", &senderPub, nil)
	trx6.SetSignature(senderPriv.Sign(trx6.SignBytes()))
	assert.NoError(t, pool.validateTx(trx6))

	trx7 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000, 1000, "ok", &senderPub, nil)
	trx7.SetSignature(senderPriv.Sign(trx7.SignBytes()))
	assert.NoError(t, pool.validateTx(trx7))

	trx8 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 1000000, 1000, "ok", &senderPub, nil)
	trx8.SetSignature(senderPriv.Sign(trx8.SignBytes()))
	assert.NoError(t, pool.validateTx(trx8))

	trx9 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 10000000, 10000, "ok", &senderPub, nil)
	trx9.SetSignature(senderPriv.Sign(trx9.SignBytes()))
	assert.NoError(t, pool.validateTx(trx9))

	trx10 := tx.NewSendTx(stampInv, 1, senderAddr, receiverAddr, 10000000, 10000, "invalid-stamp", &senderPub, nil)
	trx10.SetSignature(senderPriv.Sign(trx10.SignBytes()))
	assert.Error(t, pool.validateTx(trx10))

	trx11 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, -10000000, 10000, "negative-amount", &senderPub, nil)
	trx11.SetSignature(senderPriv.Sign(trx11.SignBytes()))
	assert.Error(t, pool.validateTx(trx11))

	trx12 := tx.NewSendTx(stamp, 1, senderAddr, receiverAddr, 10000000, -10000, "negative-fee", &senderPub, nil)
	trx12.SetSignature(senderPriv.Sign(trx12.SignBytes()))
	assert.Error(t, pool.validateTx(trx12))

	// Test stamp validity (ttl)
	sb.AppendStampAndUpdateHeight(101, crypto.GenerateTestHash())
	assert.NoError(t, pool.validateTx(trx9))
	sb.AppendStampAndUpdateHeight(102, crypto.GenerateTestHash())
	assert.Error(t, pool.validateTx(trx9))
}
