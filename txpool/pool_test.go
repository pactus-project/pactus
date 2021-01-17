package txpool

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
)

var tPool *txPool
var tSandbox *sandbox.MockSandbox
var tAcc1Addr crypto.Address
var tAcc1Pub crypto.PublicKey
var tAcc1Priv crypto.PrivateKey
var tCh chan *message.Message

func setup(t *testing.T) {
	logger.InitLogger(logger.DefaultConfig())
	tCh = make(chan *message.Message, 10)
	p, _ := NewTxPool(TestConfig(), tCh)
	tSandbox = sandbox.MockingSandbox()
	tAcc1Addr, tAcc1Pub, tAcc1Priv = crypto.GenerateTestKeyPair()
	acc1 := account.NewAccount(tAcc1Addr, 0)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)
	p.SetSandbox(tSandbox)
	tPool = p.(*txPool)
}

func shouldPublishTransaction(t *testing.T, id crypto.Hash) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-tCh:
			logger.Info("shouldPublishTransaction", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeTransactions {
				pld := msg.Payload.(*payload.TransactionsPayload)
				assert.Equal(t, pld.Transactions[0].ID(), id)
				return
			}
		}
	}
}

func TestAppendAndRemove(t *testing.T) {
	setup(t)

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(88, stamp)
	trx1 := tx.NewSubsidyTx(stamp, 89, tAcc1Addr, 25000000, "subsidy-tx")

	assert.NoError(t, tPool.appendTx(trx1))
	assert.Error(t, tPool.appendTx(trx1))
	tPool.RemoveTx(trx1.ID())
	assert.False(t, tPool.HasTx(trx1.ID()))
}

func TestSendTxValidity(t *testing.T) {
	setup(t)

	stamp := crypto.GenerateTestHash()
	senderAddr, senderPub, senderPriv := tAcc1Pub.Address(), tAcc1Pub, tAcc1Priv
	receiverAddr, _, _ := crypto.GenerateTestKeyPair()
	bigMemo := strings.Repeat("a", 1025)

	tSandbox.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1000, 1000, bigMemo, &senderPub, nil)
	trx1.SetSignature(senderPriv.Sign(trx1.SignBytes()))
	assert.Error(t, tPool.appendTx(trx1))

	trx2 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1000, 999, "invalid-fee", &senderPub, nil)
	trx2.SetSignature(senderPriv.Sign(trx2.SignBytes()))
	assert.Error(t, tPool.appendTx(trx2))

	trx3 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 10000000, 1000, "invalid-fee", &senderPub, nil)
	trx3.SetSignature(senderPriv.Sign(trx3.SignBytes()))
	assert.Error(t, tPool.appendTx(trx3))

	trx4 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 123456789, 123456, "ok", &senderPub, nil)
	trx4.SetSignature(senderPriv.Sign(trx4.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx4))

	trx5 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 123456000, 123456, "ok", &senderPub, nil)
	trx5.SetSignature(senderPriv.Sign(trx5.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx5))

	trx6 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1, 1000, "ok", &senderPub, nil)
	trx6.SetSignature(senderPriv.Sign(trx6.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx6))

	trx7 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1000, 1000, "ok", &senderPub, nil)
	trx7.SetSignature(senderPriv.Sign(trx7.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx7))

	trx8 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1000000, 1000, "ok", &senderPub, nil)
	trx8.SetSignature(senderPriv.Sign(trx8.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx8))

	trx9 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 10000000, 10000, "ok", &senderPub, nil)
	trx9.SetSignature(senderPriv.Sign(trx9.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx9))

	trx10 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+2, senderAddr, receiverAddr, 10000000, 10000, "invalid-sequence", &senderPub, nil)
	trx10.SetSignature(senderPriv.Sign(trx10.SignBytes()))
	assert.Error(t, tPool.appendTx(trx10))

	trx11 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, -10000000, 10000, "negative-amount", &senderPub, nil)
	trx11.SetSignature(senderPriv.Sign(trx11.SignBytes()))
	assert.Error(t, tPool.appendTx(trx11))

	trx12 := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 10000000, -10000, "negative-fee", &senderPub, nil)
	trx12.SetSignature(senderPriv.Sign(trx12.SignBytes()))
	assert.Error(t, tPool.appendTx(trx12))
}

func TestStampValidity(t *testing.T) {
	setup(t)

	stamp1 := crypto.GenerateTestHash()
	stamp2 := crypto.GenerateTestHash()
	stamp3 := crypto.GenerateTestHash()
	stamp4 := crypto.GenerateTestHash()
	stamp5 := crypto.GenerateTestHash()
	senderAddr, senderPub, senderPriv := tAcc1Pub.Address(), tAcc1Pub, tAcc1Priv
	receiverAddr, _, _ := crypto.GenerateTestKeyPair()

	tSandbox.AppendStampAndUpdateHeight(100, stamp1)
	tSandbox.AppendStampAndUpdateHeight(101, stamp2)
	tSandbox.AppendStampAndUpdateHeight(102, stamp3)
	tSandbox.AppendStampAndUpdateHeight(103, stamp4)

	trx1 := tx.NewSendTx(stamp1, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1000, 1000, "stamp1-ok", &senderPub, nil)
	trx1.SetSignature(senderPriv.Sign(trx1.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx1))

	tSandbox.AppendStampAndUpdateHeight(104, stamp5)

	trx2 := tx.NewSendTx(stamp1, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1000, 1000, "stamp1-invalid", &senderPub, nil)
	trx2.SetSignature(senderPriv.Sign(trx2.SignBytes()))
	assert.Error(t, tPool.appendTx(trx2))

	trx3 := tx.NewSendTx(stamp2, tSandbox.AccSeq(tAcc1Addr)+1, senderAddr, receiverAddr, 1000, 1000, "stamp2-ok", &senderPub, nil)
	trx3.SetSignature(senderPriv.Sign(trx3.SignBytes()))
	assert.NoError(t, tPool.appendTx(trx3))
}

func TestPending(t *testing.T) {
	setup(t)

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(88, stamp)
	trx := tx.NewSubsidyTx(stamp, 89, tAcc1Addr, 25000000, "subsidy-tx")

	go func() {
		for {
			msg := <-tCh
			pld := msg.Payload.(*payload.QueryTransactionsPayload)
			if pld.IDs[0].EqualsTo(trx.ID()) {
				assert.NoError(t, tPool.AppendTx(trx))
			}
		}
	}()

	assert.NotNil(t, tPool.PendingTx(trx.ID()))
	assert.True(t, tPool.pendings.Has(trx.ID()))

	// For second time it should response immediately
	assert.NotNil(t, tPool.PendingTx(trx.ID()))

	invID := crypto.GenerateTestHash()
	assert.Nil(t, tPool.PendingTx(invID))
}

func TestGetAllTransaction(t *testing.T) {
	setup(t)

	go func() {
		for {
			<-tPool.appendTxCh
		}
	}()

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)
	trxs1 := make([]*tx.Tx, 10)

	t.Run("pool is empty", func(t *testing.T) {
		trxs0 := tPool.AllTransactions()
		assert.Empty(t, trxs0)
	})

	t.Run("Fill up the pool and get all transactions", func(t *testing.T) {
		for i := 0; i < len(trxs1); i++ {
			a, _, _ := crypto.GenerateTestKeyPair()
			trx := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, tAcc1Pub.Address(), a, 1000, 1000, "stamp1-ok", &tAcc1Pub, nil)
			trx.SetSignature(tAcc1Priv.Sign(trx.SignBytes()))
			assert.NoError(t, tPool.AppendTx(trx))
			trxs1[i] = trx
		}

		trxs2 := tPool.AllTransactions()
		for i := 0; i < 10; i++ {
			// Should be in same order
			assert.Equal(t, trxs1[i].ID(), trxs2[i].ID())
		}
		assert.Equal(t, tPool.Size(), 10)
	})

	t.Run("Add one more transaction, when pool is full", func(t *testing.T) {
		a, _, _ := crypto.GenerateTestKeyPair()
		trx := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, tAcc1Pub.Address(), a, 1000, 1000, "stamp1-ok", &tAcc1Pub, nil)
		trx.SetSignature(tAcc1Priv.Sign(trx.SignBytes()))
		assert.NoError(t, tPool.AppendTx(trx))

		trxs3 := tPool.AllTransactions()
		for i := 0; i < 9; i++ {
			assert.Equal(t, trxs1[i+1].ID(), trxs3[i].ID())
		}
		assert.Equal(t, trx.ID(), trxs3[9].ID())
		assert.Equal(t, tPool.Size(), 10)
	})
}

func TestAppendAndBroadcast(t *testing.T) {
	setup(t)

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(88, stamp)
	trx := tx.NewSubsidyTx(stamp, 89, tAcc1Addr, 25000000, "subsidy-tx")

	assert.NoError(t, tPool.AppendTxAndBroadcast(trx))
	shouldPublishTransaction(t, trx.ID())

	invTrx, _ := tx.GenerateTestBondTx()
	assert.Error(t, tPool.AppendTxAndBroadcast(invTrx))
}

func TestAddSubsidyTransactions(t *testing.T) {
	setup(t)

	stamp1 := crypto.GenerateTestHash()
	stamp2 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(88, stamp1)
	proposer1, _, _ := crypto.GenerateTestKeyPair()
	proposer2, _, _ := crypto.GenerateTestKeyPair()
	trx1 := tx.NewSubsidyTx(stamp1, 88, proposer1, 25000000, "subsidy-tx-1")
	trx2 := tx.NewSubsidyTx(stamp1, 89, proposer1, 25000000, "subsidy-tx-1")
	trx3 := tx.NewSubsidyTx(stamp1, 89, proposer2, 25000000, "subsidy-tx-2")

	// Recheck on empty pool
	tPool.Recheck()

	assert.Error(t, tPool.AppendTx(trx1))
	assert.NoError(t, tPool.AppendTx(trx2))
	assert.NoError(t, tPool.AppendTx(trx3))

	tSandbox.AppendStampAndUpdateHeight(89, stamp2)

	tPool.Recheck()
	assert.Zero(t, tPool.Size())
}
