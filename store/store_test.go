package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	store *store
}

func setup(t *testing.T) *testData {
	ts := testsuite.NewTestSuite(t)

	conf := &Config{
		Path: util.TempDirPath(),
	}
	s, err := NewStore(conf, 21)
	require.NoError(t, err)

	td := &testData{
		TestSuite: ts,
		store:     s.(*store),
	}

	td.saveTestBlocks(t, 10)

	return td
}

func (td *testData) saveTestBlocks(t *testing.T, num int) {
	lastHeight, _ := td.store.LastCertificate()
	for i := 0; i < num; i++ {
		b := td.GenerateTestBlock(nil, nil)
		c := td.GenerateTestCertificate(b.Hash())

		td.store.SaveBlock(lastHeight+uint32(i+1), b, c)
		assert.NoError(t, td.store.WriteBatch())
	}
}

func TestBlockHash(t *testing.T) {
	td := setup(t)

	sb, _ := td.store.Block(1)

	assert.Equal(t, td.store.BlockHash(0), hash.UndefHash)
	assert.Equal(t, td.store.BlockHash(util.MaxUint32), hash.UndefHash)
	assert.Equal(t, td.store.BlockHash(1), sb.BlockHash)
}

func TestBlockHeight(t *testing.T) {
	td := setup(t)

	sb, _ := td.store.Block(1)

	assert.Equal(t, td.store.BlockHeight(hash.UndefHash), uint32(0))
	assert.Equal(t, td.store.BlockHeight(td.RandomHash()), uint32(0))
	assert.Equal(t, td.store.BlockHeight(sb.BlockHash), uint32(1))
}

func TestUnknownTransactionID(t *testing.T) {
	td := setup(t)

	tx, err := td.store.Transaction(td.RandomHash())
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestWriteAndClosePeacefully(t *testing.T) {
	td := setup(t)

	// After closing db, we should not crash
	assert.NoError(t, td.store.Close())
	assert.Error(t, td.store.WriteBatch())
}

func TestRetrieveBlockAndTransactions(t *testing.T) {
	td := setup(t)

	height, _ := td.store.LastCertificate()
	storedBlock, err := td.store.Block(height)
	assert.NoError(t, err)
	assert.Equal(t, height, storedBlock.Height)
	block := storedBlock.ToBlock()
	for _, trx := range block.Transactions() {
		storedTx, err := td.store.Transaction(trx.ID())
		assert.NoError(t, err)
		assert.Equal(t, storedTx.TxID, trx.ID())
		assert.Equal(t, storedTx.BlockTime, block.Header().UnixTime())
		assert.Equal(t, storedTx.ToTx().ID(), trx.ID())
	}
}

func TestRecentBlockByStamp(t *testing.T) {
	td := setup(t)

	hash1 := td.store.BlockHash(1)

	h, b := td.store.RecentBlockByStamp(hash.UndefHash.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	h, b = td.store.RecentBlockByStamp(hash1.Stamp())
	assert.Equal(t, h, uint32(1))
	assert.Equal(t, b.Hash(), hash1)

	// Saving more blocks, blocks 11 to 22
	td.saveTestBlocks(t, 12)
	hash2 := td.store.BlockHash(2)
	hash14 := td.store.BlockHash(14)
	hash22 := td.store.BlockHash(22)

	// First block should remove from the list
	h, b = td.store.RecentBlockByStamp(hash1.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	h, b = td.store.RecentBlockByStamp(hash2.Stamp())
	assert.Equal(t, h, uint32(2))
	assert.Equal(t, b.Hash(), hash2)

	h, b = td.store.RecentBlockByStamp(hash14.Stamp())
	assert.Equal(t, h, uint32(14))
	assert.Equal(t, b.Hash(), hash14)

	h, b = td.store.RecentBlockByStamp(hash22.Stamp())
	assert.Equal(t, h, uint32(22))
	assert.Equal(t, b.Hash(), hash22)

	// Reopen the store
	td.store.Close()
	s, _ := NewStore(td.store.config, 21)
	td.store = s.(*store)

	h, b = td.store.RecentBlockByStamp(hash2.Stamp())
	assert.Equal(t, h, uint32(2))
	assert.Equal(t, b.Hash(), hash2)

	// Saving one more blocks, block 23
	td.saveTestBlocks(t, 1)

	// Second block should remove from the list
	h, b = td.store.RecentBlockByStamp(hash2.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	// Genesis block
	h, b = td.store.RecentBlockByStamp(hash.UndefHash.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)
}
