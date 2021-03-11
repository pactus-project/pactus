package txpool

import (
	"fmt"
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
var tAcc1Signer crypto.Signer
var tCh chan *message.Message

func setup(t *testing.T) {
	logger.InitLogger(logger.TestConfig())
	tCh = make(chan *message.Message, 10)
	p, _ := NewTxPool(TestConfig(), tCh)
	tSandbox = sandbox.MockingSandbox()
	tAcc1Signer = crypto.GenerateTestSigner()
	tAcc1Addr = tAcc1Signer.Address()
	acc1 := account.NewAccount(tAcc1Addr, 0)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)
	p.SetSandbox(tSandbox)
	tPool = p.(*txPool)
}

func shouldPublishTransaction(t *testing.T, id tx.ID) {
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
	trx1 := tx.NewMintbaseTx(stamp, 89, tAcc1Addr, 25000000, "subsidy-tx")

	assert.NoError(t, tPool.AppendTx(trx1))
	assert.Error(t, tPool.AppendTx(trx1))
	tPool.RemoveTx(trx1.ID())
	assert.False(t, tPool.HasTx(trx1.ID()))
}

func TestAppendInvalidTransaction(t *testing.T) {
	setup(t)

	invalidTx, _ := tx.GenerateTestSendTx()
	assert.Error(t, tPool.AppendTx(invalidTx))
}

func TestPending(t *testing.T) {
	setup(t)

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(88, stamp)
	trx := tx.NewMintbaseTx(stamp, 89, tAcc1Addr, 25000000, "subsidy-tx")

	go func() {
		for {
			msg := <-tCh
			pld := msg.Payload.(*payload.QueryTransactionsPayload)
			if pld.IDs[0].EqualsTo(trx.ID()) {
				assert.NoError(t, tPool.AppendTx(trx))
			}
		}
	}()

	assert.Nil(t, tPool.PendingTx(trx.ID()))
	assert.NotNil(t, tPool.QueryTx(trx.ID()))
	assert.True(t, tPool.pendings.Has(trx.ID()))

	invID := crypto.GenerateTestHash()
	assert.Nil(t, tPool.PendingTx(invID))
}

func TestGetAllTransaction(t *testing.T) {
	setup(t)

	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(10000, stamp)
	trxs1 := make([]*tx.Tx, 10)

	t.Run("pool is empty", func(t *testing.T) {
		trxs0 := tPool.AllTransactions()
		assert.Empty(t, trxs0)
	})

	t.Run("Fill up the pool and get all transactions", func(t *testing.T) {
		for i := 0; i < len(trxs1); i++ {
			a, _, _ := crypto.GenerateTestKeyPair()
			trx := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, tAcc1Addr, a, 1000, 1000, "ok")
			tAcc1Signer.SignMsg(trx)
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
		trx := tx.NewSendTx(stamp, tSandbox.AccSeq(tAcc1Addr)+1, tAcc1Addr, a, 1000, 1000, "ok")
		tAcc1Signer.SignMsg(trx)
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
	trx := tx.NewMintbaseTx(stamp, 89, tAcc1Addr, 25000000, "subsidy-tx")

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
	trx1 := tx.NewMintbaseTx(stamp1, 88, proposer1, 25000000, "subsidy-tx-1")
	trx2 := tx.NewMintbaseTx(stamp1, 89, proposer1, 25000000, "subsidy-tx-1")
	trx3 := tx.NewMintbaseTx(stamp1, 89, proposer2, 25000000, "subsidy-tx-2")

	// Recheck on empty pool
	tPool.Recheck()

	assert.Error(t, tPool.AppendTx(trx1))
	assert.NoError(t, tPool.AppendTx(trx2))
	assert.NoError(t, tPool.AppendTx(trx3))

	tSandbox.AppendStampAndUpdateHeight(89, stamp2)

	tPool.Recheck()
	assert.Zero(t, tPool.Size())
}
