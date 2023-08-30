package txpool

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	pool    *txPool
	sandbox *sandbox.MockSandbox
	ch      chan message.Message
	testTx  *tx.Tx
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	ch := make(chan message.Message, 10)
	sandbox := sandbox.MockingSandbox(ts)
	p := NewTxPool(DefaultConfig(), ch)
	p.SetNewSandboxAndRecheck(sandbox)
	pool := p.(*txPool)
	assert.NotNil(t, pool)

	block88 := sandbox.TestStore.AddTestBlock(88)
	testTx := tx.NewSubsidyTx(block88.Stamp(), 89, ts.RandAddress(), 25000000, "subsidy-tx")

	return &testData{
		TestSuite: ts,
		pool:      pool,
		sandbox:   sandbox,
		ch:        ch,
		testTx:    testTx,
	}
}

func (td *testData) shouldPublishTransaction(t *testing.T, id tx.ID) {
	t.Helper()

	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-td.ch:
			logger.Info("shouldPublishTransaction", "message", msg)

			if msg.Type() == message.TypeTransactions {
				m := msg.(*message.TransactionsMessage)
				assert.Equal(t, m.Transactions[0].ID(), id)
				return
			}
		}
	}
}

func TestAppendAndRemove(t *testing.T) {
	td := setup(t)

	assert.NoError(t, td.pool.AppendTx(td.testTx))
	// Appending the same transaction again, should not return any error
	assert.NoError(t, td.pool.AppendTx(td.testTx))
	td.pool.RemoveTx(td.testTx.ID())
	assert.False(t, td.pool.HasTx(td.testTx.ID()), "Transaction should be removed")
}

func TestAppendInvalidTransaction(t *testing.T) {
	td := setup(t)

	invalidTx, _ := td.GenerateTestTransferTx()
	assert.Error(t, td.pool.AppendTx(invalidTx))
}

// TestFullPool tests if the pool prunes the old transactions when it is full.
func TestFullPool(t *testing.T) {
	td := setup(t)

	block10000 := td.sandbox.TestStore.AddTestBlock(10000)
	trxs := make([]*tx.Tx, td.pool.config.sendPoolSize()+1)

	signer := td.RandSigner()
	acc := account.NewAccount(0)
	acc.AddToBalance(10000000000)
	td.sandbox.UpdateAccount(signer.Address(), acc)

	// Make sure the pool is empty
	assert.Equal(t, td.pool.Size(), 0)

	for i := 0; i < len(trxs); i++ {
		trx := tx.NewTransferTx(block10000.Stamp(), acc.Sequence()+int32(i+1), signer.Address(),
			td.RandAddress(), 1000, 1000, "ok")
		signer.SignMsg(trx)
		assert.NoError(t, td.pool.AppendTx(trx))
		trxs[i] = trx
	}

	assert.False(t, td.pool.HasTx(trxs[0].ID()))
	assert.True(t, td.pool.HasTx(trxs[1].ID()))
	assert.Equal(t, td.pool.Size(), td.pool.config.sendPoolSize())
}

func TestEmptyPool(t *testing.T) {
	td := setup(t)

	assert.Empty(t, td.pool.PrepareBlockTransactions(), "pool should be empty")
}

func TestPrepareBlockTransactions(t *testing.T) {
	td := setup(t)

	block1000000 := td.sandbox.TestStore.AddTestBlock(1000000)

	acc1Signer := td.RandSigner()
	acc1 := account.NewAccount(0)
	acc1.AddToBalance(10000000000)
	td.sandbox.UpdateAccount(acc1Signer.Address(), acc1)

	val1Signer := td.RandSigner()
	val1Pub := val1Signer.PublicKey().(*bls.PublicKey)
	val1 := validator.NewValidator(val1Pub, 0)
	val1.AddToStake(10000000000)
	td.sandbox.UpdateValidator(val1)

	val2Signer := td.RandSigner()
	val2Pub := val2Signer.PublicKey().(*bls.PublicKey)
	val2 := validator.NewValidator(val2Pub, 0)
	val2.AddToStake(10000000000)
	val2.UpdateUnbondingHeight(1)
	td.sandbox.UpdateValidator(val2)

	val3Signer := td.RandSigner()
	val3Pub := val3Signer.PublicKey().(*bls.PublicKey)
	val3 := validator.NewValidator(val3Pub, 0)
	val3.AddToStake(10000000000)
	td.sandbox.UpdateValidator(val3)

	transferTx := tx.NewTransferTx(block1000000.Stamp(), acc1.Sequence()+1, acc1Signer.Address(),
		td.RandAddress(), 1000, 1000, "send-tx")
	acc1Signer.SignMsg(transferTx)

	pub, _ := td.RandBLSKeyPair()
	bondTx := tx.NewBondTx(block1000000.Stamp(), acc1.Sequence()+2, acc1Signer.Address(),
		pub.Address(), pub, 1000000000, 100000, "bond-tx")
	acc1Signer.SignMsg(bondTx)

	unbondTx := tx.NewUnbondTx(block1000000.Stamp(), val1.Sequence()+1, val1.Address(), "unbond-tx")
	val1Signer.SignMsg(unbondTx)

	withdrawTx := tx.NewWithdrawTx(block1000000.Stamp(), val2.Sequence()+1, val2.Address(),
		td.RandAddress(), 1000, 1000, "withdraw-tx")
	val2Signer.SignMsg(withdrawTx)

	td.sandbox.TestAcceptSortition = true
	sortitionTx := tx.NewSortitionTx(block1000000.Stamp(), val3.Sequence()+1, val3.Address(),
		td.RandProof())
	val3Signer.SignMsg(sortitionTx)

	assert.NoError(t, td.pool.AppendTx(transferTx))
	assert.NoError(t, td.pool.AppendTx(unbondTx))
	assert.NoError(t, td.pool.AppendTx(withdrawTx))
	assert.NoError(t, td.pool.AppendTx(bondTx))
	assert.NoError(t, td.pool.AppendTx(sortitionTx))

	trxs := td.pool.PrepareBlockTransactions()
	assert.Len(t, trxs, 5)
	assert.Equal(t, trxs[0].ID(), sortitionTx.ID())
	assert.Equal(t, trxs[1].ID(), bondTx.ID())
	assert.Equal(t, trxs[2].ID(), unbondTx.ID())
	assert.Equal(t, trxs[3].ID(), withdrawTx.ID())
	assert.Equal(t, trxs[4].ID(), transferTx.ID())
}

func TestAppendAndBroadcast(t *testing.T) {
	td := setup(t)

	assert.NoError(t, td.pool.AppendTxAndBroadcast(td.testTx))
	td.shouldPublishTransaction(t, td.testTx.ID())

	invTrx, _ := td.GenerateTestBondTx()
	assert.Error(t, td.pool.AppendTxAndBroadcast(invTrx))
}

func TestAddSubsidyTransactions(t *testing.T) {
	td := setup(t)

	block88 := td.sandbox.TestStore.AddTestBlock(88)
	proposer1 := td.RandAddress()
	proposer2 := td.RandAddress()
	trx1 := tx.NewSubsidyTx(block88.Stamp(), 88, proposer1, 25000000, "subsidy-tx-1")
	trx2 := tx.NewSubsidyTx(block88.Stamp(), 89, proposer1, 25000000, "subsidy-tx-1")
	trx3 := tx.NewSubsidyTx(block88.Stamp(), 89, proposer2, 25000000, "subsidy-tx-2")

	assert.Error(t, td.pool.AppendTx(trx1), "Expired subsidy transaction")
	assert.NoError(t, td.pool.AppendTx(trx2))
	assert.NoError(t, td.pool.AppendTx(trx3))

	td.sandbox.TestStore.AddTestBlock(89)

	td.pool.SetNewSandboxAndRecheck(sandbox.MockingSandbox(td.TestSuite))
	assert.Zero(t, td.pool.Size())
}
