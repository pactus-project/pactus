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
}

func testConfig() *Config {
	return &Config{
		MaxSize:   100,
		MinFeePAC: 0.000001,
	}
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	ch := make(chan message.Message, 10)
	sb := sandbox.MockingSandbox(ts)
	config := testConfig()
	p := NewTxPool(config, ch)
	p.SetNewSandboxAndRecheck(sb)
	pool := p.(*txPool)
	assert.NotNil(t, pool)

	return &testData{
		TestSuite: ts,
		pool:      pool,
		sandbox:   sb,
		ch:        ch,
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
			logger.Info("shouldPublishTransaction", "msg", msg)

			if msg.Type() == message.TypeTransaction {
				m := msg.(*message.TransactionsMessage)
				assert.Equal(t, m.Transactions[0].ID(), id)

				return
			}
		}
	}
}

func TestAppendAndRemove(t *testing.T) {
	td := setup(t)

	height := td.RandHeight()
	td.sandbox.TestStore.AddTestBlock(height)
	testTrx := tx.NewSubsidyTx(height+1, td.RandAccAddress(), 1e9, "subsidy-tx")

	assert.NoError(t, td.pool.AppendTx(testTrx))
	assert.True(t, td.pool.HasTx(testTrx.ID()))
	assert.Equal(t, testTrx, td.pool.PendingTx(testTrx.ID()))

	// Appending the same transaction again, should not return any error
	assert.NoError(t, td.pool.AppendTx(testTrx))

	td.pool.RemoveTx(testTrx.ID())
	assert.False(t, td.pool.HasTx(testTrx.ID()), "Transaction should be removed")
	assert.Nil(t, td.pool.PendingTx(testTrx.ID()))
}

func TestAppendInvalidTransaction(t *testing.T) {
	td := setup(t)

	invalidTx, _ := td.GenerateTestTransferTx()
	assert.Error(t, td.pool.AppendTx(invalidTx))
}

// TestFullPool tests if the pool prunes the old transactions when it is full.
func TestFullPool(t *testing.T) {
	td := setup(t)

	randHeight := td.RandHeight()
	_ = td.sandbox.TestStore.AddTestBlock(randHeight)
	trxs := make([]*tx.Tx, td.pool.config.transferPoolSize()+1)

	senderAddr := td.RandAccAddress()
	senderAcc := account.NewAccount(0)
	senderAcc.AddToBalance(1000e9)
	td.sandbox.UpdateAccount(senderAddr, senderAcc)

	// Make sure the pool is empty
	assert.Equal(t, td.pool.Size(), 0)

	for i := 0; i < len(trxs); i++ {
		trx := tx.NewTransferTx(randHeight+1, senderAddr,
			td.RandAccAddress(), 1e9, 1000_000_000, "ok")

		assert.NoError(t, td.pool.AppendTx(trx))
		trxs[i] = trx
	}

	assert.False(t, td.pool.HasTx(trxs[0].ID()))
	assert.True(t, td.pool.HasTx(trxs[1].ID()))
	assert.Equal(t, td.pool.Size(), td.pool.config.transferPoolSize())
}

func TestEmptyPool(t *testing.T) {
	td := setup(t)

	assert.Empty(t, td.pool.PrepareBlockTransactions(), "pool should be empty")
}

func TestPrepareBlockTransactions(t *testing.T) {
	td := setup(t)

	randHeight := td.RandHeight() + td.sandbox.TestParams.UnbondInterval
	_ = td.sandbox.TestStore.AddTestBlock(randHeight)

	acc1Addr := td.RandAccAddress()
	acc1 := account.NewAccount(0)
	acc1.AddToBalance(1000e9)
	td.sandbox.UpdateAccount(acc1Addr, acc1)

	val1PubKey, _ := td.RandBLSKeyPair()
	val1 := validator.NewValidator(val1PubKey, 0)
	val1.AddToStake(1000e9)
	td.sandbox.UpdateValidator(val1)

	val2PubKey, _ := td.RandBLSKeyPair()
	val2 := validator.NewValidator(val2PubKey, 0)
	val2.AddToStake(1000e9)
	val2.UpdateUnbondingHeight(1)
	td.sandbox.UpdateValidator(val2)

	val3PubKey, _ := td.RandBLSKeyPair()
	val3 := validator.NewValidator(val3PubKey, 0)
	val3.AddToStake(1000e9)
	td.sandbox.UpdateValidator(val3)

	transferTx := tx.NewTransferTx(randHeight+1, acc1Addr,
		td.RandAccAddress(), 1e9, 100_000_000, "transfer-tx")

	pub, _ := td.RandBLSKeyPair()
	bondTx := tx.NewBondTx(randHeight+2, acc1Addr,
		pub.ValidatorAddress(), pub, 1e9, 100_000_000, "bond-tx")

	unbondTx := tx.NewUnbondTx(randHeight+3, val1.Address(), "unbond-tx")

	withdrawTx := tx.NewWithdrawTx(randHeight+4, val2.Address(),
		td.RandAccAddress(), 1e9, 100_000_000, "withdraw-tx")

	td.sandbox.TestAcceptSortition = true
	sortitionTx := tx.NewSortitionTx(randHeight, val3.Address(),
		td.RandProof())

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

	height := td.RandHeight()
	td.sandbox.TestStore.AddTestBlock(height)
	testTrx := tx.NewSubsidyTx(height+1, td.RandAccAddress(), 1e9, "subsidy-tx")

	assert.NoError(t, td.pool.AppendTxAndBroadcast(testTrx))
	td.shouldPublishTransaction(t, testTrx.ID())

	invTrx, _ := td.GenerateTestBondTx()
	assert.Error(t, td.pool.AppendTxAndBroadcast(invTrx))
}

func TestAddSubsidyTransactions(t *testing.T) {
	td := setup(t)

	randHeight := td.RandHeight()
	td.sandbox.TestStore.AddTestBlock(randHeight)
	proposer1 := td.RandAccAddress()
	proposer2 := td.RandAccAddress()
	trx1 := tx.NewSubsidyTx(randHeight, proposer1, 1e9, "subsidy-tx-1")
	trx2 := tx.NewSubsidyTx(randHeight+1, proposer1, 1e9, "subsidy-tx-1")
	trx3 := tx.NewSubsidyTx(randHeight+1, proposer2, 1e9, "subsidy-tx-2")

	err := td.pool.AppendTx(trx1)
	assert.ErrorIs(t, err, AppendError{
		Err: execution.PastLockTimeError{
			LockTime: randHeight,
		},
	})

	err = td.pool.AppendTx(trx2)
	assert.NoError(t, err)

	err = td.pool.AppendTx(trx3)
	assert.NoError(t, err)

	td.sandbox.TestStore.AddTestBlock(randHeight + 1)

	td.pool.SetNewSandboxAndRecheck(td.sandbox)
	assert.Zero(t, td.pool.Size())
}
