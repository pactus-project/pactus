package txpool

import (
	"errors"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/pipeline"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	pool *txPool
	sbx  *sandbox.MockSandbox
	pipe *pipeline.MockPipeline[message.Message]
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

	pipe := pipeline.MockingPipeline[message.Message]()
	sbx := sandbox.MockingSandbox(ts)
	config := testDefaultConfig()
	if cfg != nil {
		config = cfg
	}
	poolInt := NewTxPool(config, sbx.TestStore, pipe)
	poolInt.SetNewSandboxAndRecheck(sbx)
	pool := poolInt.(*txPool)
	assert.NotNil(t, pool)

	sbx.TestAcceptSortition = true

	return &testData{
		TestSuite: ts,
		pool:      pool,
		sbx:       sbx,
		pipe:      pipe,
	}
}

func (td *testData) shouldPublishTransaction(t *testing.T, txID tx.ID) {
	t.Helper()

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timer.C:
			require.NoError(t, errors.New("Timeout"))

			return

		case msg := <-td.pipe.UnsafeGetChannel():
			logger.Info("shouldPublishTransaction", "msg", msg)

			if msg.Type() == message.TypeTransaction {
				m := msg.(*message.TransactionsMessage)
				assert.Equal(t, txID, m.Transactions[0].ID())

				return
			}
		}
	}
}

// makeValidTransferTx makes a valid Transfer transaction for testing purpose.
func (td *testData) makeValidTransferTx(options ...func(tm *testsuite.TransactionMaker)) *tx.Tx {
	options = append(options, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	trx := td.GenerateTestTransferTx(options...)
	signer := trx.Payload().Signer()

	acc := td.sbx.MakeNewAccount(signer)
	acc.AddToBalance(trx.Payload().Value() + trx.Fee())
	td.sbx.UpdateAccount(signer, acc)

	return trx
}

// makeValidBatchTransferTx make a valid Batch transfer transaction for test purpose.
func (td *testData) makeValidBatchTransferTx(options ...func(tm *testsuite.TransactionMaker)) *tx.Tx {
	options = append(options, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	trx := td.GenerateTestBatchTransferTx(options...)
	signer := trx.Payload().Signer()

	acc := td.sbx.MakeNewAccount(signer)
	acc.AddToBalance(trx.Payload().Value() + trx.Fee())
	td.sbx.UpdateAccount(signer, acc)

	return trx
}

// makeValidSubsidyTx make a valid Batch transfer transaction for test purpose.
func (td *testData) makeValidSubsidyTx(options ...func(tm *testsuite.TransactionMaker)) *tx.Tx {
	return td.GenerateTestSubsidyTx(options...)
}

// makeValidBondTx makes a valid Bond transaction for testing purpose.
func (td *testData) makeValidBondTx(options ...func(tm *testsuite.TransactionMaker)) *tx.Tx {
	options = append(options, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	trx := td.GenerateTestBondTx(options...)
	signer := trx.Payload().Signer()

	acc := td.sbx.MakeNewAccount(signer)
	acc.AddToBalance(trx.Payload().Value() + trx.Fee())
	td.sbx.UpdateAccount(signer, acc)

	return trx
}

// makeValidUnbondTx makes a valid Unbond transaction for testing purpose.
// Ensure that the signer key is set through the options.
func (td *testData) makeValidUnbondTx(options ...func(tm *testsuite.TransactionMaker)) *tx.Tx {
	options = append(options, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	trx := td.GenerateTestUnbondTx(options...)

	validatorPublicKey := trx.PublicKey().(*bls.PublicKey)
	val := td.sbx.MakeNewValidator(validatorPublicKey)
	td.sbx.UpdateValidator(val)

	return trx
}

// makeValidWithdrawTx makes a valid Withdraw transaction for testing purpose.
// Ensure that the signer key is set through the options.
func (td *testData) makeValidWithdrawTx(options ...func(tm *testsuite.TransactionMaker)) *tx.Tx {
	options = append(options, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	trx := td.GenerateTestWithdrawTx(options...)

	validatorPublicKey := trx.PublicKey().(*bls.PublicKey)
	val := td.sbx.MakeNewValidator(validatorPublicKey)
	val.AddToStake(trx.Payload().Value() + trx.Fee())
	val.UpdateUnbondingHeight(td.sbx.CurrentHeight() - (td.sbx.TestParams.UnbondInterval + 1))
	td.sbx.UpdateValidator(val)

	return trx
}

// makeValidSortitionTx makes a valid Sortition transaction for testing purpose.
// Ensure that the signer key is set through the options.
func (td *testData) makeValidSortitionTx(options ...func(tm *testsuite.TransactionMaker)) *tx.Tx {
	options = append(options, testsuite.TransactionWithLockTime(td.sbx.CurrentHeight()))
	trx := td.GenerateTestSortitionTx(options...)

	validatorPublicKey := trx.PublicKey().(*bls.PublicKey)
	val := td.sbx.MakeNewValidator(validatorPublicKey)
	val.UpdateLastBondingHeight(td.sbx.CurrentHeight() - (td.sbx.TestParams.BondInterval + 1))
	td.sbx.UpdateValidator(val)

	return trx
}

func TestAppendAndRemove(t *testing.T) {
	td := setup(t, nil)

	trx := td.makeValidTransferTx()

	assert.NoError(t, td.pool.AppendTx(trx))
	assert.True(t, td.pool.HasTx(trx.ID()))
	assert.Equal(t, trx, td.pool.PendingTx(trx.ID()))

	td.pool.removeTx(trx.ID())
	assert.False(t, td.pool.HasTx(trx.ID()), "Transaction should be removed")
	assert.Nil(t, td.pool.PendingTx(trx.ID()))
}

func TestAppendSameTransaction(t *testing.T) {
	td := setup(t, nil)

	trx := td.makeValidTransferTx()

	err := td.pool.AppendTx(trx)
	assert.NoError(t, err)

	err = td.pool.AppendTx(trx)
	assert.ErrorIs(t, err, execution.TransactionCommittedError{ID: trx.ID()})
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
	trx10 := td.makeValidSubsidyTx()
	trx11 := td.makeValidTransferTx(testsuite.TransactionWithSigner(prv1))
	trx12 := td.makeValidBondTx(testsuite.TransactionWithSigner(prv2))
	trx13 := td.makeValidUnbondTx(testsuite.TransactionWithSigner(prv4))
	trx20 := td.makeValidSubsidyTx()
	trx21 := td.makeValidTransferTx(testsuite.TransactionWithSigner(prv1))
	trx30 := td.makeValidSubsidyTx()
	trx31 := td.makeValidBondTx(testsuite.TransactionWithSigner(prv3))
	trx32 := td.makeValidSortitionTx(testsuite.TransactionWithSigner(prv4))
	trx40 := td.makeValidSubsidyTx()
	trx41 := td.makeValidUnbondTx(testsuite.TransactionWithSigner(prv3))
	trx42 := td.makeValidTransferTx(testsuite.TransactionWithSigner(prv2))
	trx50 := td.makeValidSubsidyTx()
	trx51 := td.makeValidWithdrawTx(testsuite.TransactionWithSigner(prv3))
	trx52 := td.makeValidTransferTx(testsuite.TransactionWithSigner(prv2))
	trx53 := td.makeValidBatchTransferTx(testsuite.TransactionWithSigner(prv2))

	// Commit the first block
	blk1, cert1 := td.GenerateTestBlock(1,
		testsuite.BlockWithTransactions([]*tx.Tx{trx10, trx11, trx12, trx13}))
	td.sbx.TestStore.SaveBlock(blk1, cert1)

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
		height uint32
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
		td.sbx.TestStore.SaveBlock(blk, cert)

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
		trx := td.makeValidTransferTx(testsuite.TransactionWithSigner(accPrv))
		blk, cert := td.GenerateTestBlock(td.RandHeight(), testsuite.BlockWithTransactions([]*tx.Tx{trx}))
		td.sbx.TestStore.SaveBlock(blk, cert)

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
			testTrx := td.makeValidTransferTx(
				testsuite.TransactionWithSigner(accPrv),
				testsuite.TransactionWithFee(tt.fee))

			err := td.pool.AppendTx(testTrx)
			if tt.withErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
	})

	t.Run("Test non-indexed signer", func(t *testing.T) {
		trx := td.makeValidTransferTx(testsuite.TransactionWithFee(0))

		err := td.pool.AppendTx(trx)
		assert.Error(t, err)
	})
}

func TestAppendInvalidTransaction(t *testing.T) {
	td := setup(t, nil)

	invTrx := td.GenerateTestTransferTx()

	err := td.pool.AppendTx(invTrx)
	assert.Error(t, err)
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
		trx := td.makeValidTransferTx()

		assert.NoError(t, td.pool.AppendTx(trx))
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

	transferTx := td.makeValidTransferTx()
	bondTx := td.makeValidBondTx(testsuite.TransactionWithValidatorPublicKey(pub1))
	unbondTx := td.makeValidUnbondTx(testsuite.TransactionWithSigner(prv2))
	withdrawTx := td.makeValidWithdrawTx(testsuite.TransactionWithSigner(prv3))
	sortitionTx := td.makeValidSortitionTx(testsuite.TransactionWithSigner(prv4))
	batchTransferTx := td.makeValidBatchTransferTx(testsuite.TransactionWithSigner(prv5))

	assert.NoError(t, td.pool.AppendTx(transferTx))
	assert.NoError(t, td.pool.AppendTx(unbondTx))
	assert.NoError(t, td.pool.AppendTx(withdrawTx))
	assert.NoError(t, td.pool.AppendTx(bondTx))
	assert.NoError(t, td.pool.AppendTx(sortitionTx))
	assert.NoError(t, td.pool.AppendTx(batchTransferTx))

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
		td.sbx.TestStore.AddTestBlock(randHeight)
		trx := td.makeValidSubsidyTx(testsuite.TransactionWithLockTime(randHeight))

		err := td.pool.AppendTx(trx)
		assert.ErrorIs(t, err, execution.LockTimeExpiredError{
			LockTime: randHeight,
		})
	})

	t.Run("valid transaction: Should add it to the pool", func(t *testing.T) {
		td := setup(t, nil)

		randHeight := td.RandHeight()
		td.sbx.TestStore.AddTestBlock(randHeight)
		trx := td.makeValidSubsidyTx(testsuite.TransactionWithLockTime(randHeight + 1))

		err := td.pool.AppendTx(trx)
		assert.NoError(t, err)
	})
}

func TestRecheckTransactions(t *testing.T) {
	td := setup(t, nil)

	trx := td.makeValidSubsidyTx()

	err := td.pool.AppendTx(trx)
	assert.NoError(t, err)
	assert.Equal(t, 1, td.pool.Size())

	td.pool.SetNewSandboxAndRecheck(td.sbx)
	assert.Equal(t, 0, td.pool.Size())
}

func TestAppendAndBroadcast(t *testing.T) {
	t.Run("Invalid transaction: Should return error", func(t *testing.T) {
		td := setup(t, nil)

		invTrx := td.GenerateTestTransferTx()
		assert.Error(t, td.pool.AppendTxAndBroadcast(invTrx))
	})

	t.Run("Valid transaction with valid fee: Should add to pool and broadcast", func(t *testing.T) {
		td := setup(t, nil)

		trx := td.makeValidTransferTx()

		err := td.pool.AppendTxAndBroadcast(trx)
		assert.NoError(t, err)

		assert.Equal(t, 1, td.pool.Size())
		td.shouldPublishTransaction(t, trx.ID())
	})

	t.Run("Valid transaction with zero fee: Should broadcast but not add to the pool", func(t *testing.T) {
		td := setup(t, nil)

		trx := td.makeValidTransferTx(testsuite.TransactionWithFee(0))

		err := td.pool.AppendTxAndBroadcast(trx)
		assert.NoError(t, err)

		assert.Zero(t, td.pool.Size())
		td.shouldPublishTransaction(t, trx.ID())
	})
}

func TestAllPendingTxs(t *testing.T) {
	td := setup(t, nil)

	trxs := td.pool.AllPendingTxs()
	assert.Empty(t, trxs, 0)

	pub1, _ := td.RandBLSKeyPair()
	_, prv2 := td.RandBLSKeyPair()

	transferTx := td.makeValidTransferTx()
	bondTx := td.makeValidBondTx(testsuite.TransactionWithValidatorPublicKey(pub1))
	unbondTx := td.makeValidUnbondTx()
	withdrawTx := td.makeValidWithdrawTx()
	sortitionTx := td.makeValidSortitionTx(testsuite.TransactionWithSigner(prv2))
	batchTransferTx := td.makeValidBatchTransferTx()

	assert.NoError(t, td.pool.AppendTx(transferTx))
	assert.NoError(t, td.pool.AppendTx(bondTx))
	assert.NoError(t, td.pool.AppendTx(unbondTx))
	assert.NoError(t, td.pool.AppendTx(withdrawTx))
	assert.NoError(t, td.pool.AppendTx(sortitionTx))
	assert.NoError(t, td.pool.AppendTx(batchTransferTx))

	trxs = td.pool.AllPendingTxs()
	assert.Len(t, trxs, 6)
}

func TestEstimatedFee(t *testing.T) {
	td := setup(t, nil)

	estimatedFee := td.pool.EstimatedFee(td.RandAmount(), payload.TypeTransfer)
	assert.Equal(t, td.pool.config.fixedFee(), estimatedFee)
}
