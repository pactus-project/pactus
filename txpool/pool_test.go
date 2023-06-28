package txpool

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tPool *txPool
var tSandbox *sandbox.MockSandbox
var tCh chan message.Message
var tTestTx *tx.Tx

func setup(t *testing.T) {
	tCh = make(chan message.Message, 10)
	tSandbox = sandbox.MockingSandbox()
	p := NewTxPool(DefaultConfig(), tCh)
	p.SetNewSandboxAndRecheck(tSandbox)
	tPool = p.(*txPool)
	assert.NotNil(t, tPool)

	block88 := tSandbox.TestStore.AddTestBlock(88)
	tTestTx = tx.NewSubsidyTx(block88.Stamp(), 89, crypto.GenerateTestAddress(), 25000000, "subsidy-tx")
}

func shouldPublishTransaction(t *testing.T, id tx.ID) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-tCh:
			logger.Info("shouldPublishTransaction", "message", msg)

			if msg.Type() == message.MessageTypeTransactions {
				m := msg.(*message.TransactionsMessage)
				assert.Equal(t, m.Transactions[0].ID(), id)
				return
			}
		}
	}
}

func TestAppendAndRemove(t *testing.T) {
	setup(t)

	assert.NoError(t, tPool.AppendTx(tTestTx))
	// Appending the same transaction again, should not return any error
	assert.NoError(t, tPool.AppendTx(tTestTx))
	tPool.RemoveTx(tTestTx.ID())
	assert.False(t, tPool.HasTx(tTestTx.ID()), "Transaction should be removed")
}

func TestAppendInvalidTransaction(t *testing.T) {
	setup(t)

	invalidTx, _ := tx.GenerateTestSendTx()
	assert.Error(t, tPool.AppendTx(invalidTx))
}

// TestFullPool tests if the pool prunes the old transactions when it is full.
func TestFullPool(t *testing.T) {
	setup(t)

	block10000 := tSandbox.TestStore.AddTestBlock(10000)
	trxs := make([]*tx.Tx, tPool.config.sendPoolSize()+1)

	signer := bls.GenerateTestSigner()
	acc := account.NewAccount(0)
	acc.AddToBalance(10000000000)
	tSandbox.UpdateAccount(signer.Address(), acc)

	// Make sure the pool is empty
	assert.Equal(t, tPool.Size(), 0)

	for i := 0; i < len(trxs); i++ {
		trx := tx.NewTransferTx(block10000.Stamp(), acc.Sequence()+int32(i+1), signer.Address(),
			crypto.GenerateTestAddress(), 1000, 1000, "ok")
		signer.SignMsg(trx)
		assert.NoError(t, tPool.AppendTx(trx))
		trxs[i] = trx
	}

	assert.False(t, tPool.HasTx(trxs[0].ID()))
	assert.True(t, tPool.HasTx(trxs[1].ID()))
	assert.Equal(t, tPool.Size(), tPool.config.sendPoolSize())
}

func TestEmptyPool(t *testing.T) {
	setup(t)

	assert.Empty(t, tPool.PrepareBlockTransactions(), "pool should be empty")
}

func TestPrepareBlockTransactions(t *testing.T) {
	setup(t)

	block1000000 := tSandbox.TestStore.AddTestBlock(1000000)

	acc1Signer := bls.GenerateTestSigner()
	acc1 := account.NewAccount(0)
	acc1.AddToBalance(10000000000)
	tSandbox.UpdateAccount(acc1Signer.Address(), acc1)

	val1Signer := bls.GenerateTestSigner()
	val1Pub := val1Signer.PublicKey().(*bls.PublicKey)
	val1 := validator.NewValidator(val1Pub, 0)
	val1.AddToStake(10000000000)
	tSandbox.UpdateValidator(val1)

	val2Signer := bls.GenerateTestSigner()
	val2Pub := val2Signer.PublicKey().(*bls.PublicKey)
	val2 := validator.NewValidator(val2Pub, 0)
	val2.AddToStake(10000000000)
	val2.UpdateUnbondingHeight(1)
	tSandbox.UpdateValidator(val2)

	val3Signer := bls.GenerateTestSigner()
	val3Pub := val3Signer.PublicKey().(*bls.PublicKey)
	val3 := validator.NewValidator(val3Pub, 0)
	val3.AddToStake(10000000000)
	tSandbox.UpdateValidator(val3)

	sendTx := tx.NewTransferTx(block1000000.Stamp(), acc1.Sequence()+1, acc1Signer.Address(),
		crypto.GenerateTestAddress(), 1000, 1000, "send-tx")
	acc1Signer.SignMsg(sendTx)

	pub, _ := bls.GenerateTestKeyPair()
	bondTx := tx.NewBondTx(block1000000.Stamp(), acc1.Sequence()+2, acc1Signer.Address(),
		pub.Address(), pub, 1000, 1000, "bond-tx")
	acc1Signer.SignMsg(bondTx)

	unbondTx := tx.NewUnbondTx(block1000000.Stamp(), val1.Sequence()+1, val1.Address(), "unbond-tx")
	val1Signer.SignMsg(unbondTx)

	withdrawTx := tx.NewWithdrawTx(block1000000.Stamp(), val2.Sequence()+1, val2.Address(),
		crypto.GenerateTestAddress(), 1000, 1000, "withdraw-tx")
	val2Signer.SignMsg(withdrawTx)

	tSandbox.TestAcceptSortition = true
	sortitionTx := tx.NewSortitionTx(block1000000.Stamp(), val3.Sequence()+1, val3.Address(),
		sortition.GenerateRandomProof())
	val3Signer.SignMsg(sortitionTx)

	assert.NoError(t, tPool.AppendTx(sendTx))
	assert.NoError(t, tPool.AppendTx(unbondTx))
	assert.NoError(t, tPool.AppendTx(withdrawTx))
	assert.NoError(t, tPool.AppendTx(bondTx))
	assert.NoError(t, tPool.AppendTx(sortitionTx))

	trxs := tPool.PrepareBlockTransactions()
	assert.Len(t, trxs, 5)
	assert.Equal(t, trxs[0].ID(), sortitionTx.ID())
	assert.Equal(t, trxs[1].ID(), bondTx.ID())
	assert.Equal(t, trxs[2].ID(), unbondTx.ID())
	assert.Equal(t, trxs[3].ID(), withdrawTx.ID())
	assert.Equal(t, trxs[4].ID(), sendTx.ID())
}

func TestAppendAndBroadcast(t *testing.T) {
	setup(t)

	assert.NoError(t, tPool.AppendTxAndBroadcast(tTestTx))
	shouldPublishTransaction(t, tTestTx.ID())

	invTrx, _ := tx.GenerateTestBondTx()
	assert.Error(t, tPool.AppendTxAndBroadcast(invTrx))
}

func TestAddSubsidyTransactions(t *testing.T) {
	setup(t)

	block88 := tSandbox.TestStore.AddTestBlock(88)
	proposer1 := crypto.GenerateTestAddress()
	proposer2 := crypto.GenerateTestAddress()
	trx1 := tx.NewSubsidyTx(block88.Stamp(), 88, proposer1, 25000000, "subsidy-tx-1")
	trx2 := tx.NewSubsidyTx(block88.Stamp(), 89, proposer1, 25000000, "subsidy-tx-1")
	trx3 := tx.NewSubsidyTx(block88.Stamp(), 89, proposer2, 25000000, "subsidy-tx-2")

	assert.Error(t, tPool.AppendTx(trx1), "Expired subsidy transaction")
	assert.NoError(t, tPool.AppendTx(trx2))
	assert.NoError(t, tPool.AppendTx(trx3))

	tSandbox.TestStore.AddTestBlock(89)

	tPool.SetNewSandboxAndRecheck(sandbox.MockingSandbox())
	assert.Zero(t, tPool.Size())
}
