package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
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
		TxCacheSize:        1024,
		SortitionCacheSize: 1024,
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

	s, err := NewStore(config)
	require.NoError(t, err)
	assert.False(t, s.IsPruned(), "empty store should not be in prune mode")

	td := &testData{
		TestSuite: ts,
		store:     s.(*store),
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
	assert.Equal(t, uint32(10), store.LastCertificate().Height())
}

func TestBlockHash(t *testing.T) {
	td := setup(t, nil)

	sb, _ := td.store.Block(1)

	assert.Equal(t, td.store.BlockHash(0), hash.UndefHash)
	assert.Equal(t, td.store.BlockHash(util.MaxUint32), hash.UndefHash)
	assert.Equal(t, td.store.BlockHash(1), sb.BlockHash)
}

func TestBlockHeight(t *testing.T) {
	td := setup(t, nil)

	sb, _ := td.store.Block(1)

	assert.Equal(t, td.store.BlockHeight(hash.UndefHash), uint32(0))
	assert.Equal(t, td.store.BlockHeight(td.RandHash()), uint32(0))
	assert.Equal(t, td.store.BlockHeight(sb.BlockHash), uint32(1))
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
	assert.Equal(t, blk.PrevCertificate().Height(), lastHeight-1)

	for _, trx := range blk.Transactions() {
		committedTx, err := td.store.Transaction(trx.ID())
		assert.NoError(t, err)
		assert.Equal(t, blk.Header().UnixTime(), committedTx.BlockTime)
		assert.Equal(t, trx.ID(), committedTx.TxID)
		assert.Equal(t, lastHeight, committedTx.Height)
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
			pubKey, err := td.store.PublicKey(addr)
			assert.NoError(t, err)

			if addr.IsAccountAddress() {
				assert.Equal(t, pubKey.AccountAddress(), addr)
			} else if addr.IsValidatorAddress() {
				assert.Equal(t, pubKey.ValidatorAddress(), addr)
			}
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

	// Find a public key that we have already indexed in the database.
	cBlkOne, _ := td.store.Block(1)
	blkOne, _ := cBlkOne.ToBlock()
	trx0PubKey := blkOne.Transactions()[0].PublicKey()
	assert.NotNil(t, trx0PubKey)
	knownPubKey := trx0PubKey.(*bls.PublicKey)

	lastCert := td.store.LastCertificate()
	lastHeight := lastCert.Height()
	randPubKey, _ := td.RandBLSKeyPair()

	trx0 := tx.NewTransferTx(lastHeight, knownPubKey.AccountAddress(), td.RandAccAddress(), 1, 1, "")
	trx1 := tx.NewTransferTx(lastHeight, randPubKey.AccountAddress(), td.RandAccAddress(), 1, 1, "")
	trx2 := tx.NewTransferTx(lastHeight, randPubKey.AccountAddress(), td.RandAccAddress(), 1, 1, "")

	trx0.StripPublicKey()
	trx1.SetPublicKey(randPubKey)
	trx2.StripPublicKey()

	tests := []struct {
		trx    *tx.Tx
		failed bool
	}{
		{trx0, false}, // indexed public key and stripped
		{trx1, false}, // not stripped
		{trx2, true},  // unknown public key and stripped
	}

	for _, test := range tests {
		trxs := block.Txs{test.trx}
		blk, _ := td.GenerateTestBlock(td.RandHeight(), testsuite.BlockWithTransactions(trxs))

		trxData, _ := test.trx.Bytes()
		blkData, _ := blk.Bytes()

		cTrx := CommittedTx{
			store:  td.store,
			TxID:   test.trx.ID(),
			Height: lastHeight + 1,
			Data:   trxData,
		}
		cBlk := CommittedBlock{
			store:     td.store,
			BlockHash: blk.Hash(),
			Height:    lastHeight + 1,
			Data:      blkData,
		}

		//
		if test.failed {
			_, err := cBlk.ToBlock()
			assert.ErrorIs(t, err, PublicKeyNotFoundError{
				Address: test.trx.Payload().Signer(),
			})

			_, err = cTrx.ToTx()
			assert.ErrorIs(t, err, PublicKeyNotFoundError{
				Address: test.trx.Payload().Signer(),
			})
		} else {
			_, err := cBlk.ToBlock()
			assert.NoError(t, err)

			_, err = cTrx.ToTx()
			assert.NoError(t, err)
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
	cb := func(pruned bool, pruningHeight uint32) bool {
		if pruned {
			totalPruned++
		}
		lastPruningHeight = pruningHeight

		return false
	}

	t.Run("Not enough block to prune", func(t *testing.T) {
		totalPruned = uint32(0)
		lastPruningHeight = uint32(0)

		err := td.store.Prune(cb)
		assert.NoError(t, err)
		assert.False(t, td.store.isPruned)

		assert.Zero(t, totalPruned)
		assert.Zero(t, lastPruningHeight)
	})

	t.Run("Prune database", func(t *testing.T) {
		totalPruned = uint32(0)
		lastPruningHeight = uint32(0)

		blk, cert := td.GenerateTestBlock(blockPerDay + 9)
		td.store.SaveBlock(blk, cert)
		err := td.store.WriteBatch()
		require.NoError(t, err)

		err = td.store.Prune(cb)
		assert.NoError(t, err)
		assert.True(t, td.store.isPruned)

		assert.Equal(t, uint32(9), totalPruned)
		assert.Equal(t, uint32(1), lastPruningHeight)
	})

	t.Run("Commit new block", func(t *testing.T) {
		blk, cert := td.GenerateTestBlock(blockPerDay + 10)
		td.store.SaveBlock(blk, cert)
		err := td.store.WriteBatch()
		require.NoError(t, err)

		cBlk, err := td.store.Block(10)
		assert.Error(t, err)
		assert.Nil(t, cBlk)
	})
}
