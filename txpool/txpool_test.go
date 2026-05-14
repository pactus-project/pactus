package txpool

import (
	"testing"
	"time"

	"github.com/ezex-io/gopkg/pipeline"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	pool          *txPool
	sbx           *sandbox.MockSandbox
	store         *store.MockStore
	broadcastPipe pipeline.Pipeline[message.Message]
	eventPipe     pipeline.Pipeline[any]
}

func testDefaultConfig() *Config {
	return DefaultConfig()
}

func testConsumptionalConfig() *Config {
	return &Config{
		MaxSize:           10,
		ConsumptionWindow: 3,
		Fee: &FeeConfig{
			FixedFee:   0,
			DailyLimit: 360,
			UnitPrice:  0.000005,
		},
	}
}

func setup(t *testing.T, cfg *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	broadcastPipe := pipeline.New[message.Message](t.Context())
	eventPipe := pipeline.New[any](t.Context())
	sbx := sandbox.NewMockSandbox(ts.Ctrl)
	store := store.MockingStore(ts)
	config := testDefaultConfig()
	if cfg != nil {
		config = cfg
	}
	poolInt := NewTxPool(config, store, broadcastPipe, eventPipe)
	poolInt.SetNewSandboxAndRecheck(sbx)
	pool := poolInt.(*txPool)
	assert.NotNil(t, pool)

	// Mock sandbox methods
	params := &param.Params{
		UnbondInterval: 10,
		BondInterval:   10,
	}
	currentHeight := types.Height(100)
	sbx.EXPECT().Params().Return(params).AnyTimes()
	sbx.EXPECT().CurrentHeight().Return(currentHeight).AnyTimes()
	sbx.EXPECT().IsBanned(crypto.Address{}).Return(false).AnyTimes()
	sbx.EXPECT().RecentTransaction(tx.ID{}).Return(false).AnyTimes()
	sbx.EXPECT().MakeNewAccount(gomock.Any()).Return(account.NewAccount(0)).AnyTimes()
	sbx.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Return().AnyTimes()
	sbx.EXPECT().MakeNewValidator(gomock.Any()).Return(validator.NewValidator(bls.PublicKey{}, 0)).AnyTimes()
	sbx.EXPECT().UpdateValidator(gomock.Any()).Return().AnyTimes()

	randHeight := ts.RandHeight(
		testsuite.HeightWithMin(params.UnbondInterval))
	_ = store.AddTestBlock(randHeight)

	return &testData{
		TestSuite:     ts,
		pool:          pool,
		sbx:           sbx,
		store:         store,
		broadcastPipe: broadcastPipe,
		eventPipe:     eventPipe,
	}
}

func (td *testData) shouldPublishTransaction(t *testing.T, txID tx.ID) {
	t.Helper()

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timer.C:
			require.Fail(t, "Publish transaction timeout")

			return

		case msg := <-td.broadcastPipe.UnsafeGetChannel():
			logger.Info("shouldPublishTransaction", "msg", msg)

			if msg.Type() == message.TypeTransaction {
				m := msg.(*message.TransactionsMessage)
				assert.Equal(t, txID, m.Transactions[0].ID())

				return
			}
		}
	}
}

// makeTransferTx makes a Transfer transaction for testing purpose.
func (td *testData) makeTransferTx(opts ...testsuite.TransactionMakerOption) *tx.Tx {
	opts = util.Prepend(opts, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	return td.GenerateTestTransferTx(opts...)
}

// makeBatchTransferTx make a Batch transfer transaction for test purpose.
func (td *testData) makeBatchTransferTx(opts ...testsuite.TransactionMakerOption) *tx.Tx {
	opts = util.Prepend(opts, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	return td.GenerateTestBatchTransferTx(opts...)
}

// makeSubsidyTx make a valid Batch transfer transaction for test purpose.
func (td *testData) makeSubsidyTx(opts ...testsuite.TransactionMakerOption) *tx.Tx {
	opts = util.Prepend(opts, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))

	return td.GenerateTestSubsidyTx(opts...)
}

// makeBondTx makes a Bond transaction for testing purpose.
func (td *testData) makeBondTx(opts ...testsuite.TransactionMakerOption) *tx.Tx {
	opts = util.Prepend(opts, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	return td.GenerateTestBondTx(opts...)
}

// makeUnbondTx makes a Unbond transaction for testing purpose.
// Ensure that the signer key is set through the opts.
func (td *testData) makeUnbondTx(opts ...testsuite.TransactionMakerOption) *tx.Tx {
	opts = util.Prepend(opts, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	return td.GenerateTestUnbondTx(opts...)
}

// makeWithdrawTx makes a Withdraw transaction for testing purpose.
// Ensure that the signer key is set through the opts.
func (td *testData) makeWithdrawTx(opts ...testsuite.TransactionMakerOption) *tx.Tx {
	opts = util.Prepend(opts, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	return td.GenerateTestWithdrawTx(opts...)
}

// makeSortitionTx makes a Sortition transaction for testing purpose.
// Ensure that the signer key is set through the opts.
func (td *testData) makeSortitionTx(opts ...testsuite.TransactionMakerOption) *tx.Tx {
	opts = util.Prepend(opts, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	return td.GenerateTestSortitionTx(opts...)
}

func TestAppendAndRemove(t *testing.T) {
	td := setup(t, nil)

	trx := td.makeTransferTx()

	require.NoError(t, td.pool.AppendTx(trx))
	assert.True(t, td.pool.HasTx(trx.ID()))
	assert.Equal(t, trx, td.pool.PendingTx(trx.ID()))

	td.pool.removeTx(trx.ID())
	assert.False(t, td.pool.HasTx(trx.ID()), "Transaction should be removed")
	assert.Nil(t, td.pool.PendingTx(trx.ID()))
}

func TestAppendSameTransaction(t *testing.T) {
	td := setup(t, nil)

	trx := td.makeTransferTx()

	err := td.pool.AppendTx(trx)
	require.NoError(t, err)

	err = td.pool.AppendTx(trx)
	require.ErrorIs(t, err, execution.TransactionCommittedError{ID: trx.ID()})
}

func TestDisableConsumption(t *testing.T) {
	td := setup(t, testDefaultConfig())

	trx := td.GenerateTestTransferTx()
	assert.Zero(t, td.pool.consumptionalFee(trx))
}

func TestCalculatingConsumption(t *testing.T) {
	td := setup(t, testConsumptionalConfig())

	// Generate keys for different transaction signers
	_, prv1 := td.RandEd25519KeyPair()
	pub2, prv2 := td.RandEd25519KeyPair()
	pub3, prv3 := td.RandBLSKeyPair()
	pub4, prv4 := td.RandBLSKeyPair()

	// Generate different types of transactions
	trx10 := td.makeSubsidyTx()
	trx11 := td.makeTransferTx(testsuite.TransactionWithSigner(prv1))
	trx12 := td.makeBondTx(testsuite.TransactionWithSigner(prv2))
	trx13 := td.makeUnbondTx(testsuite.TransactionWithSigner(prv4))
	trx20 := td.makeSubsidyTx()
	trx21 := td.makeTransferTx(testsuite.TransactionWithSigner(prv1))
	trx30 := td.makeSubsidyTx()
	trx31 := td.makeBondTx(testsuite.TransactionWithSigner(prv3))
	trx32 := td.makeSortitionTx(testsuite.TransactionWithSigner(prv4))
	trx40 := td.makeSubsidyTx()
	trx41 := td.makeUnbondTx(testsuite.TransactionWithSigner(prv3))
	trx42 := td.makeTransferTx(testsuite.TransactionWithSigner(prv2))
	trx50 := td.makeSubsidyTx()
	trx51 := td.makeWithdrawTx(testsuite.TransactionWithSigner(prv3))
	trx52 := td.makeTransferTx(testsuite.TransactionWithSigner(prv2))
	trx53 := td.makeBatchTransferTx(testsuite.TransactionWithSigner(prv2))

	// Commit the first block
	blk1, cert1 := td.GenerateTestBlock(1,
		testsuite.BlockWithTransactions([]*tx.Tx{trx10, trx11, trx12, trx13}))
	td.store.SaveBlock(blk1, cert1)

	// Expected consumption map after transactions
	diff2 := 0
	if trx42.SerializeSize() > trx12.SerializeSize() {
		diff2 = trx42.SerializeSize() - trx12.SerializeSize()
	}

	expected := map[crypto.Address]int{
		pub2.AccountAddress():   trx52.SerializeSize() + trx53.SerializeSize() + diff2,
		pub3.AccountAddress():   trx31.SerializeSize(),
		pub3.ValidatorAddress(): trx41.SerializeSize() + trx51.SerializeSize(),
		pub4.ValidatorAddress(): trx32.SerializeSize() - trx13.SerializeSize(),
	}

	tests := []struct {
		height types.Height
		txs    []*tx.Tx
	}{
		{2, []*tx.Tx{trx20, trx21}},
		{3, []*tx.Tx{trx30, trx31, trx32}},
		{4, []*tx.Tx{trx40, trx41, trx42}},
		{5, []*tx.Tx{trx50, trx51, trx52, trx53}},
	}

	for _, tt := range tests {
		// Generate a block with the transactions for the given height
		blk, cert := td.GenerateTestBlock(tt.height, testsuite.BlockWithTransactions(tt.txs))
		td.store.SaveBlock(blk, cert)

		// Handle the block in the transaction pool
		td.pool.HandleCommittedBlock(blk)

		t.Logf("consumption Map: %v\n", td.pool.consumptionMap)
	}

	require.Equal(t, expected, td.pool.consumptionMap)
}

func TestEstimatedConsumptionalFee(t *testing.T) {
	td := setup(t, testConsumptionalConfig())

	t.Run("Test indexed signer", func(t *testing.T) {
		_, accPrv := td.RandEd25519KeyPair()
		trx := td.makeTransferTx(testsuite.TransactionWithSigner(accPrv))
		blk, cert := td.GenerateTestBlock(td.RandHeight(), testsuite.BlockWithTransactions([]*tx.Tx{trx}))
		td.store.SaveBlock(blk, cert)

		tests := []struct {
			fee     amount.Amount
			withErr bool
		}{
			{0, false},
			{0, false},
			{2_310_000, false},
			{3_090_000, false},
			{7_740_000, false},
			{9_300_000, false},
			{0, true},
		}

		for _, tt := range tests {
			testTrx := td.makeTransferTx(
				testsuite.TransactionWithSigner(accPrv),
				testsuite.TransactionWithFee(tt.fee))

			err := td.pool.AppendTx(testTrx)
			if tt.withErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	})

	t.Run("Test non-indexed signer", func(t *testing.T) {
		trx := td.makeTransferTx(testsuite.TransactionWithFee(0))

		err := td.pool.AppendTx(trx)
		require.Error(t, err)
	})
}

func TestAppendInvalidTransaction(t *testing.T) {
	td := setup(t, nil)

	t.Run("basic check error", func(t *testing.T) {
		trx := td.makeTransferTx()
		trx.SetSignature(nil)

		err := td.pool.AppendTx(trx)
		require.ErrorIs(t, err, tx.BasicCheckError{
			Reason: "no signature",
		})
	})

	t.Run("execution error", func(t *testing.T) {
		invTrx := td.GenerateTestTransferTx()

		err := td.pool.AppendTx(invTrx)
		require.ErrorIs(t, err, executor.AccountNotFoundError{
			Address: invTrx.Payload().Signer(),
		})
	})
}

// TestFullPool tests if the pool prunes the old transactions when it is full.
func TestFullPool(t *testing.T) {
	conf := testDefaultConfig()
	conf.MaxSize = 10
	td := setup(t, conf)

	trxs := make([]*tx.Tx, td.pool.config.transferPoolSize()+1)

	// Make sure the pool is empty
	assert.Equal(t, 0, td.pool.Size())

	for i := 0; i < len(trxs); i++ {
		trx := td.makeTransferTx()

		require.NoError(t, td.pool.AppendTx(trx))
		trxs[i] = trx
	}

	assert.False(t, td.pool.HasTx(trxs[0].ID()))
	assert.True(t, td.pool.HasTx(trxs[1].ID()))
	assert.Equal(t, td.pool.config.transferPoolSize(), td.pool.Size())
}

func TestEmptyPool(t *testing.T) {
	td := setup(t, nil)

	trxs := td.pool.AllPendingTxs()
	assert.Empty(t, trxs, "pool should be empty")
}

func TestPrepareBlockTransactions(t *testing.T) {
	td := setup(t, nil)

	pub1, _ := td.RandBLSKeyPair()
	_, prv2 := td.RandBLSKeyPair()
	_, prv3 := td.RandBLSKeyPair()
	_, prv4 := td.RandBLSKeyPair()
	_, prv5 := td.RandBLSKeyPair()

	transferTx := td.makeTransferTx()
	bondTx := td.makeBondTx(testsuite.TransactionWithValidatorPublicKey(pub1))
	unbondTx := td.makeUnbondTx(testsuite.TransactionWithSigner(prv2))
	withdrawTx := td.makeWithdrawTx(testsuite.TransactionWithSigner(prv3))
	sortitionTx := td.makeSortitionTx(testsuite.TransactionWithSigner(prv4))
	batchTransferTx := td.makeBatchTransferTx(testsuite.TransactionWithSigner(prv5))

	require.NoError(t, td.pool.AppendTx(transferTx))
	require.NoError(t, td.pool.AppendTx(unbondTx))
	require.NoError(t, td.pool.AppendTx(withdrawTx))
	require.NoError(t, td.pool.AppendTx(bondTx))
	require.NoError(t, td.pool.AppendTx(sortitionTx))
	require.NoError(t, td.pool.AppendTx(batchTransferTx))

	trxs := td.pool.PrepareBlockTransactions()
	assert.Len(t, trxs, 6)
	assert.Equal(t, sortitionTx.ID(), trxs[0].ID())
	assert.Equal(t, bondTx.ID(), trxs[1].ID())
	assert.Equal(t, unbondTx.ID(), trxs[2].ID())
	assert.Equal(t, withdrawTx.ID(), trxs[3].ID())
	assert.Equal(t, transferTx.ID(), trxs[4].ID())
	assert.Equal(t, batchTransferTx.ID(), trxs[5].ID())
}

func TestAddSubsidyTransactions(t *testing.T) {
	t.Run("invalid transaction: Should return error", func(t *testing.T) {
		td := setup(t, nil)

		randHeight := td.RandHeight()
		td.store.AddTestBlock(randHeight)
		trx := td.makeSubsidyTx(testsuite.TransactionWithLockTime(randHeight))

		err := td.pool.AppendTx(trx)
		require.ErrorIs(t, err, execution.LockTimeExpiredError{
			LockTime: randHeight,
		})
	})

	t.Run("valid transaction: Should add it to the pool", func(t *testing.T) {
		td := setup(t, nil)

		randHeight := td.RandHeight()
		td.store.AddTestBlock(randHeight)
		trx := td.makeSubsidyTx(testsuite.TransactionWithLockTime(randHeight + 1))

		err := td.pool.AppendTx(trx)
		require.NoError(t, err)
	})
}

func TestRecheckTransactions(t *testing.T) {
	td := setup(t, nil)

	trx := td.makeValidSubsidyTx()

	err := td.pool.AppendTx(trx)
	require.NoError(t, err)
	assert.Equal(t, 1, td.pool.Size())

	td.pool.SetNewSandboxAndRecheck(td.sbx)
	assert.Equal(t, 0, td.pool.Size())
}

func TestAppendAndBroadcast(t *testing.T) {
	t.Run("Invalid transaction: Should return error", func(t *testing.T) {
		td := setup(t, nil)

		invTrx := td.GenerateTestTransferTx()
		require.Error(t, td.pool.AppendTxAndBroadcast(invTrx))
	})

	t.Run("Valid transaction with valid fee: Should add to pool and broadcast", func(t *testing.T) {
		td := setup(t, nil)

		trx := td.makeTransferTx()

		err := td.pool.AppendTxAndBroadcast(trx)
		require.NoError(t, err)

		assert.Equal(t, 1, td.pool.Size())
		td.shouldPublishTransaction(t, trx.ID())
	})

	t.Run("Valid transaction with zero fee: Should broadcast but not add to the pool", func(t *testing.T) {
		td := setup(t, nil)

		trx := td.makeTransferTx(testsuite.TransactionWithFee(0))

		err := td.pool.AppendTxAndBroadcast(trx)
		require.NoError(t, err)

		assert.Zero(t, td.pool.Size())
		td.shouldPublishTransaction(t, trx.ID())
	})
}

func TestAllPendingTxs(t *testing.T) {
	td := setup(t, nil)

	trxs := td.pool.AllPendingTxs()
	assert.Empty(t, trxs, "%+v", 0)

	pub1, _ := td.RandBLSKeyPair()
	_, prv2 := td.RandBLSKeyPair()

	transferTx := td.makeTransferTx()
	bondTx := td.makeBondTx(testsuite.TransactionWithValidatorPublicKey(pub1))
	unbondTx := td.makeUnbondTx()
	withdrawTx := td.makeWithdrawTx()
	sortitionTx := td.makeSortitionTx(testsuite.TransactionWithSigner(prv2))
	batchTransferTx := td.makeBatchTransferTx()

	require.NoError(t, td.pool.AppendTx(transferTx))
	require.NoError(t, td.pool.AppendTx(bondTx))
	require.NoError(t, td.pool.AppendTx(unbondTx))
	require.NoError(t, td.pool.AppendTx(withdrawTx))
	require.NoError(t, td.pool.AppendTx(sortitionTx))
	require.NoError(t, td.pool.AppendTx(batchTransferTx))

	trxs = td.pool.AllPendingTxs()
	assert.Len(t, trxs, 6)
}

func TestEstimatedFee(t *testing.T) {
	td := setup(t, nil)

	estimatedFee := td.pool.EstimatedFee(td.RandAmount(), payload.TypeTransfer)
	assert.Equal(t, td.pool.config.fixedFee(), estimatedFee)
}
