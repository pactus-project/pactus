package txpool

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
)

var tPool *txPool
var tSandbox *sandbox.MockSandbox
var tAcc1Addr crypto.Address
var tAcc1Signer crypto.Signer
var tCh chan payload.Payload

func setup(t *testing.T) {
	logger.InitLogger(logger.TestConfig())
	tCh = make(chan payload.Payload, 10)
	tSandbox = sandbox.MockingSandbox()
	p, err := NewTxPool(TestConfig(), tCh)
	assert.NoError(t, err)
	p.SetNewSandboxAndRecheck(tSandbox)
	tAcc1Signer = bls.GenerateTestSigner()
	tAcc1Addr = tAcc1Signer.Address()
	acc1 := account.NewAccount(tAcc1Addr, 0)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1)
	tPool = p.(*txPool)
}

func shouldPublishTransaction(t *testing.T, id tx.ID) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case pld := <-tCh:
			logger.Info("shouldPublishTransaction", "pld", pld)

			if pld.Type() == payload.PayloadTypeTransactions {
				pld := pld.(*payload.TransactionsPayload)
				assert.Equal(t, pld.Transactions[0].ID(), id)
				return
			}
		}
	}
}

func TestAppendAndRemove(t *testing.T) {
	setup(t)

	hash88 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(88, hash88)
	trx1 := tx.NewMintbaseTx(hash88.Stamp(), 89, tAcc1Addr, 25000000, "subsidy-tx")

	assert.NoError(t, tPool.AppendTx(trx1))
	assert.NoError(t, tPool.AppendTx(trx1))
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

	hash88 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(88, hash88)
	trx := tx.NewMintbaseTx(hash88.Stamp(), 89, tAcc1Addr, 25000000, "subsidy-tx")

	// Increat the waiting time for testing
	tPool.config.WaitingTimeout = 2 * time.Second

	go func(ch chan payload.Payload) {
		for {
			pld := <-ch
			fmt.Printf("Received a message payload: %v\n", pld.Fingerprint())
			p := pld.(*payload.QueryTransactionsPayload)
			if p.IDs[0].EqualsTo(trx.ID()) {
				assert.NoError(t, tPool.AppendTx(trx))
			}
		}
	}(tCh)

	assert.Nil(t, tPool.PendingTx(trx.ID()))
	assert.NotNil(t, tPool.QueryTx(trx.ID()))
	assert.True(t, tPool.pendings.Has(trx.ID()))

	invID := hash.GenerateTestHash()
	assert.Nil(t, tPool.PendingTx(invID))
}

func TestGetAllTransaction(t *testing.T) {
	setup(t)

	hash10000 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(10000, hash10000)
	trxs1 := make([]*tx.Tx, 10)

	t.Run("pool is empty", func(t *testing.T) {
		trxs0 := tPool.AllTransactions()
		assert.Empty(t, trxs0)
	})

	t.Run("Fill up the pool and get all transactions", func(t *testing.T) {
		for i := 0; i < len(trxs1); i++ {
			a := crypto.GenerateTestAddress()
			trx := tx.NewSendTx(hash10000.Stamp(), tSandbox.AccSeq(tAcc1Addr)+1, tAcc1Addr, a, 1000, 1000, "ok")
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
		a := crypto.GenerateTestAddress()
		trx := tx.NewSendTx(hash10000.Stamp(), tSandbox.AccSeq(tAcc1Addr)+1, tAcc1Addr, a, 1000, 1000, "ok")
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

	hash88 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(88, hash88)
	trx := tx.NewMintbaseTx(hash88.Stamp(), 89, tAcc1Addr, 25000000, "subsidy-tx")

	assert.NoError(t, tPool.AppendTxAndBroadcast(trx))
	shouldPublishTransaction(t, trx.ID())

	invTrx, _ := tx.GenerateTestBondTx()
	assert.Error(t, tPool.AppendTxAndBroadcast(invTrx))
}

func TestAddSubsidyTransactions(t *testing.T) {
	setup(t)

	hash88 := hash.GenerateTestHash()
	hash89 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(88, hash88)
	proposer1 := crypto.GenerateTestAddress()
	proposer2 := crypto.GenerateTestAddress()
	trx1 := tx.NewMintbaseTx(hash88.Stamp(), 88, proposer1, 25000000, "subsidy-tx-1")
	trx2 := tx.NewMintbaseTx(hash88.Stamp(), 89, proposer1, 25000000, "subsidy-tx-1")
	trx3 := tx.NewMintbaseTx(hash88.Stamp(), 89, proposer2, 25000000, "subsidy-tx-2")

	assert.Error(t, tPool.AppendTx(trx1), "Expired subsidy transaction")
	assert.NoError(t, tPool.AppendTx(trx2))
	assert.NoError(t, tPool.AppendTx(trx3))

	tSandbox.AppendNewBlock(89, hash89)

	tPool.SetNewSandboxAndRecheck(sandbox.MockingSandbox())
	assert.Zero(t, tPool.Size())
}
