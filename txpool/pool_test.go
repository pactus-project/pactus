package txpool

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/execution"
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
	testTx := tx.NewSubsidyTx(block88.Stamp(), 89, ts.RandAccAddress(), 25000000, "subsidy-tx")

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

	valKey := td.RandValKey()
	acc := account.NewAccount(0)
	acc.AddToBalance(10000000000)
	td.sandbox.UpdateAccount(valKey.Address(), acc)

	// Make sure the pool is empty
	assert.Equal(t, td.pool.Size(), 0)

	for i := 0; i < len(trxs); i++ {
		trx := tx.NewTransferTx(block10000.Stamp(), td.sandbox.CurrentHeight()+1, valKey.Address(),
			td.RandAccAddress(), 1000, 1000, "ok")
		valKey.Sign(trx.SignBytes())
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

	randHeight := td.RandHeight() + td.sandbox.TestParams.UnbondInterval
	randBlock := td.sandbox.TestStore.AddTestBlock(randHeight)

	acc1ValKey := td.RandValKey()
	acc1 := account.NewAccount(0)
	acc1.AddToBalance(10000000000)
	td.sandbox.UpdateAccount(acc1ValKey.Address(), acc1)

	valKey1 := td.RandValKey()
	val1Pub := valKey1.PublicKey()
	val1 := validator.NewValidator(val1Pub, 0)
	val1.AddToStake(10000000000)
	td.sandbox.UpdateValidator(val1)

	valKey2 := td.RandValKey()
	val2Pub := valKey2.PublicKey()
	val2 := validator.NewValidator(val2Pub, 0)
	val2.AddToStake(10000000000)
	val2.UpdateUnbondingHeight(1)
	td.sandbox.UpdateValidator(val2)

	valKey3 := td.RandValKey()
	val3Pub := valKey3.PublicKey()
	val3 := validator.NewValidator(val3Pub, 0)
	val3.AddToStake(10000000000)
	td.sandbox.UpdateValidator(val3)


	transferTx := tx.NewTransferTx(randBlock.Stamp(), randHeight+1, acc1ValKey.Address(),
		td.RandAccAddress(), 1000, 1000, "send-tx")
	td.HelperSignTransaction(acc1ValKey.PrivateKey(), transferTx)

	pub, _ := td.RandBLSKeyPair()
	bondTx := tx.NewBondTx(randBlock.Stamp(), randHeight+2, acc1ValKey.Address(),
		pub.ValidatorAddress(), pub, 1000000000, 100000, "bond-tx")
	td.HelperSignTransaction(acc1ValKey.PrivateKey(), bondTx)

	unbondTx := tx.NewUnbondTx(randBlock.Stamp(), randHeight+3, val1.Address(), "unbond-tx")
	td.HelperSignTransaction(valKey1.PrivateKey(), unbondTx)

	withdrawTx := tx.NewWithdrawTx(randBlock.Stamp(), randHeight+4, val2.Address(),
		td.RandAccAddress(), 1000, 1000, "withdraw-tx")
	td.HelperSignTransaction(valKey2.PrivateKey(), withdrawTx)

	td.sandbox.TestAcceptSortition = true
	sortitionTx := tx.NewSortitionTx(randBlock.Stamp(), randHeight, val3.Address(),
		td.RandProof())
	td.HelperSignTransaction(valKey3.PrivateKey(), sortitionTx)

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

	randHeight := td.RandHeight()
	randBlock := td.sandbox.TestStore.AddTestBlock(randHeight)
	proposer1 := td.RandAccAddress()
	proposer2 := td.RandAccAddress()
	trx1 := tx.NewSubsidyTx(randBlock.Stamp(), randHeight, proposer1, 25000000, "subsidy-tx-1")
	trx2 := tx.NewSubsidyTx(randBlock.Stamp(), randHeight+1, proposer1, 25000000, "subsidy-tx-1")
	trx3 := tx.NewSubsidyTx(randBlock.Stamp(), randHeight+1, proposer2, 25000000, "subsidy-tx-2")

	err := td.pool.AppendTx(trx1)
	assert.ErrorIs(t, err, execution.PastLockTimeError{LockTime: randHeight})

	err = td.pool.AppendTx(trx2)
	assert.NoError(t, err)

	err = td.pool.AppendTx(trx3)
	assert.NoError(t, err)

	td.sandbox.TestStore.AddTestBlock(randHeight + 1)

	td.pool.SetNewSandboxAndRecheck(td.sandbox)
	assert.Zero(t, td.pool.Size())
}
