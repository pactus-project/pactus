package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	store *store
}

func testConfig() *Config {
	return &Config{
		Path:               util.TempDirPath(),
		TxCacheWindow:      1024,
		SeedCacheWindow:    1024,
		AccountCacheSize:   1024,
		PublicKeyCacheSize: 1024,
		BannedAddrs:        make(map[crypto.Address]bool),
	}
}

func setup(t *testing.T, config *Config) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	if config == nil {
		config = testConfig()
	}

	storeInt, err := NewStore(config)
	require.NoError(t, err)
	assert.False(t, storeInt.IsPruned(), "empty store should not be in prune mode")
	assert.Zero(t, storeInt.PruningHeight(), "pruning height should be zero for an empty store")

	td := &testData{
		TestSuite: ts,
		store:     storeInt.(*store),
	}

	// Save 10 blocks
	for height := uint32(0); height < 10; height++ {
		blk, cert := td.GenerateTestBlock(height + 1)
		td.store.SaveBlock(blk, cert)
		assert.NoError(t, td.store.WriteBatch())
	}

	return td
}

func TestReopenStore(t *testing.T) {
	td := setup(t, nil)
	td.store.Close()
	store, _ := NewStore(td.store.config)

	assert.False(t, store.IsPruned())
	assert.Zero(t, store.PruningHeight())
	assert.Equal(t, uint32(10), store.LastCertificate().Height())
}

func TestBlockHash(t *testing.T) {
	td := setup(t, nil)

	sb, _ := td.store.Block(1)

	assert.Equal(t, hash.UndefHash, td.store.BlockHash(0))
	assert.Equal(t, hash.UndefHash, td.store.BlockHash(util.MaxUint32))
	assert.Equal(t, sb.BlockHash, td.store.BlockHash(1))
}

func TestBlockHeight(t *testing.T) {
	td := setup(t, nil)

	sb, _ := td.store.Block(1)

	assert.Equal(t, uint32(0), td.store.BlockHeight(hash.UndefHash))
	assert.Equal(t, uint32(0), td.store.BlockHeight(td.RandHash()))
	assert.Equal(t, uint32(1), td.store.BlockHeight(sb.BlockHash))
}

func TestUnknownTransactionID(t *testing.T) {
	td := setup(t, nil)

	trx, err := td.store.Transaction(td.RandHash())
	assert.Error(t, err)
	assert.Nil(t, trx)
}

func TestWriteAndClosePeacefully(t *testing.T) {
	td := setup(t, nil)

	// After closing the database, writing will result in an error
	td.store.Close()
	assert.Error(t, td.store.WriteBatch())
}

func TestRetrieveBlockAndTransactions(t *testing.T) {
	td := setup(t, nil)

	lastCert := td.store.LastCertificate()
	lastHeight := lastCert.Height()
	cBlk, err := td.store.Block(lastHeight)
	assert.NoError(t, err)
	assert.Equal(t, lastHeight, cBlk.Height)
	blk, _ := cBlk.ToBlock()
	assert.Equal(t, lastHeight-1, blk.PrevCertificate().Height())

	for _, trx := range blk.Transactions() {
		committedTx, err := td.store.Transaction(trx.ID())
		assert.NoError(t, err)
		assert.Equal(t, committedTx.BlockTime, blk.Header().UnixTime())
		assert.Equal(t, committedTx.TxID, trx.ID())
		assert.Equal(t, committedTx.Height, lastHeight)
		trx2, _ := committedTx.ToTx()
		assert.Equal(t, trx.ID(), trx2.ID())
	}
}

func TestIndexingPublicKeys(t *testing.T) {
	td := setup(t, nil)

	t.Run("Query existing public key", func(t *testing.T) {
		cBlk, _ := td.store.Block(1)
		blk, _ := cBlk.ToBlock()
		for _, trx := range blk.Transactions() {
			addr := trx.Payload().Signer()
			pub, err := td.store.PublicKey(addr)
			assert.NoError(t, err)

			assert.True(t, trx.PublicKey().EqualsTo(pub))
		}
	})

	t.Run("Query non existing public key", func(t *testing.T) {
		randValAddress := td.RandValAddress()
		pubKey, err := td.store.PublicKey(randValAddress)
		assert.Error(t, err)
		assert.Nil(t, pubKey)
	})
}

func TestStrippedPublicKey(t *testing.T) {
	td := setup(t, nil)

	lastHeight := td.store.LastCertificate().Height()
	_, blsPrv := td.RandBLSKeyPair()
	committedTrx1 := td.GenerateTestTransferTx(
		testsuite.TransactionWithBLSSigner(blsPrv),
	)
	_, ed25519Prv := td.RandEd25519KeyPair()
	committedTrx2 := td.GenerateTestTransferTx(
		testsuite.TransactionWithEd25519Signer(ed25519Prv),
	)
	blk0, cert0 := td.GenerateTestBlock(lastHeight+1,
		testsuite.BlockWithTransactions([]*tx.Tx{committedTrx1, committedTrx2}))
	td.store.SaveBlock(blk0, cert0)
	err := td.store.writeBatch()
	require.NoError(t, err)

	// We have some known and index public key, run tests...
	trx1 := td.GenerateTestTransferTx(
		testsuite.TransactionWithBLSSigner(blsPrv),
	)
	trx2 := td.GenerateTestTransferTx(
		testsuite.TransactionWithEd25519Signer(ed25519Prv),
	)
	trx3 := td.GenerateTestTransferTx(
		testsuite.TransactionWithBLSSigner(blsPrv),
	)
	trx4 := td.GenerateTestTransferTx(
		testsuite.TransactionWithEd25519Signer(ed25519Prv),
	)
	trx5 := td.GenerateTestTransferTx()

	trx3.StripPublicKey()
	trx4.StripPublicKey()
	trx5.StripPublicKey()

	tests := []struct {
		trx    *tx.Tx
		failed bool
	}{
		{trx1, false}, // indexed public key and not stripped
		{trx2, false}, // indexed public key and not stripped
		{trx3, false}, // indexed public key and stripped
		{trx4, false}, // indexed public key and stripped
		{trx5, true},  // unknown public key and stripped
	}

	for no, tt := range tests {
		trxs := block.Txs{tt.trx}
		blockHeight := td.store.LastCertificate().Height()
		blk, cert := td.GenerateTestBlock(blockHeight+1, testsuite.BlockWithTransactions(trxs))
		td.store.SaveBlock(blk, cert)
		err := td.store.writeBatch()
		require.NoError(t, err)

		cBlk, err := td.store.Block(blockHeight + 1)
		require.NoError(t, err)

		cTrx, err := td.store.Transaction(tt.trx.ID())
		require.NoError(t, err)

		//
		if tt.failed {
			_, err := cBlk.ToBlock()
			assert.ErrorIs(t, err, PublicKeyNotFoundError{
				Address: tt.trx.Payload().Signer(),
			}, "test %d failed, expected error", no+1)

			_, err = cTrx.ToTx()
			assert.ErrorIs(t, err, PublicKeyNotFoundError{
				Address: tt.trx.Payload().Signer(),
			}, "test %d failed, expected error", no+1)
		} else {
			_, err := cBlk.ToBlock()
			assert.NoError(t, err, "test %d failed, not expected error", no+1)

			_, err = cTrx.ToTx()
			assert.NoError(t, err, "test %d failed, not expected error", no+1)
		}
	}
}

func TestIsBanned(t *testing.T) {
	conf := testConfig()
	td := setup(t, conf)

	bannedAddr := td.RandValAddress()
	conf.BannedAddrs[bannedAddr] = true

	assert.False(t, td.store.IsBanned(td.RandAccAddress()))
	assert.True(t, td.store.IsBanned(bannedAddr))
}

func TestPruneBlock(t *testing.T) {
	conf := testConfig()
	td := setup(t, conf)

	t.Run("Prune existing block", func(t *testing.T) {
		height := uint32(1)
		cBlkOne, _ := td.store.Block(height)
		blkOne, _ := cBlkOne.ToBlock()
		pruned, err := td.store.pruneBlock(height)
		assert.True(t, pruned)
		assert.NoError(t, err)

		err = td.store.WriteBatch()
		assert.NoError(t, err)

		cBlk, _ := td.store.Block(height)
		assert.Nil(t, cBlk)

		h := td.store.BlockHash(height)
		assert.Equal(t, hash.UndefHash, h)

		for _, trx := range blkOne.Transactions() {
			cTrx, _ := td.store.Transaction(trx.ID())
			assert.Nil(t, cTrx)
		}
	})

	t.Run("Prune non existing block", func(t *testing.T) {
		height := uint32(11)
		pruned, err := td.store.pruneBlock(height)
		assert.False(t, pruned)
		assert.NoError(t, err)

		err = td.store.WriteBatch()
		assert.NoError(t, err)
	})
}

func TestPrune(t *testing.T) {
	conf := testConfig()
	conf.RetentionDays = 1
	td := setup(t, conf)

	totalPruned := uint32(0)
	lastPruningHeight := uint32(0)
	callback := func(pruned bool, pruningHeight uint32) bool {
		if pruned {
			totalPruned++
		}
		lastPruningHeight = pruningHeight

		return false
	}

	t.Run("Not enough block to prune", func(t *testing.T) {
		totalPruned = uint32(0)
		lastPruningHeight = uint32(0)

		// Store doesn't have blocks for one day
		err := td.store.Prune(callback)
		assert.NoError(t, err)

		assert.Zero(t, totalPruned)
		assert.Zero(t, lastPruningHeight)
	})

	t.Run("Prune database", func(t *testing.T) {
		totalPruned = uint32(0)
		lastPruningHeight = uint32(0)

		blk, cert := td.GenerateTestBlock(blockPerDay + 7)
		td.store.SaveBlock(blk, cert)
		err := td.store.WriteBatch()
		require.NoError(t, err)

		blk, cert = td.GenerateTestBlock(blockPerDay + 8)
		td.store.SaveBlock(blk, cert)
		err = td.store.WriteBatch()
		require.NoError(t, err)

		// It should remove blocks [1..8]
		err = td.store.Prune(callback)
		assert.NoError(t, err)

		assert.Equal(t, uint32(8), totalPruned)
		assert.Equal(t, uint32(1), lastPruningHeight)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		td.store.Close()
		td.store.config.TxCacheWindow = 1
		s, err := NewStore(td.store.config)
		require.NoError(t, err)
		td.store = s.(*store)

		assert.True(t, td.store.IsPruned(), "store should be in prune mode")
		assert.Equal(t, uint32(8), td.store.PruningHeight())
	})

	t.Run("Commit new block", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(blockPerDay + 9)
		td.store.SaveBlock(blk, cert)
		err := td.store.WriteBatch()
		require.NoError(t, err)

		cBlk, err := td.store.Block(9)
		assert.Error(t, err)
		assert.Nil(t, cBlk)

		assert.Equal(t, uint32(9), td.store.PruningHeight())
	})
}

func TestCancelPrune(t *testing.T) {
	conf := testConfig()
	conf.RetentionDays = 1
	td := setup(t, conf)

	hits := uint32(0)
	callback := func(_ bool, _ uint32) bool {
		hits++

		return true // Cancel pruning
	}

	t.Run("Cancel Pruning database", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(blockPerDay + 7)
		td.store.SaveBlock(blk, cert)
		err := td.store.WriteBatch()
		require.NoError(t, err)

		err = td.store.Prune(callback)
		assert.NoError(t, err)

		assert.Equal(t, uint32(1), hits)
	})
}

func TestRecentTransaction(t *testing.T) {
	td := setup(t, nil)

	lastHeight := td.store.LastCertificate().Height()
	oldTrx := td.GenerateTestTransferTx()
	blkOld, certOld := td.GenerateTestBlock(lastHeight+1,
		testsuite.BlockWithTransactions([]*tx.Tx{oldTrx}))
	td.store.SaveBlock(blkOld, certOld)
	err := td.store.writeBatch()
	require.NoError(t, err)
	assert.True(t, td.store.RecentTransaction(oldTrx.ID()))

	blk, cert := td.GenerateTestBlock(lastHeight + td.store.txStore.txCacheWindow + 2)
	td.store.SaveBlock(blk, cert)
	err = td.store.writeBatch()
	require.NoError(t, err)

	assert.False(t, td.store.RecentTransaction(oldTrx.ID()))
	assert.False(t, td.store.RecentTransaction(td.RandHash()))
}
