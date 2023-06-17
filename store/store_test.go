package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tStore *store

func setup(t *testing.T) {
	conf := &Config{
		Path: util.TempDirPath(),
	}
	s, err := NewStore(conf, 21)
	require.NoError(t, err)

	tStore = s.(*store)
	SaveTestBlocks(t, 10)
}

func SaveTestBlocks(t *testing.T, num int) {
	lastHeight, _ := tStore.LastCertificate()
	for i := 0; i < num; i++ {
		b := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())

		tStore.SaveBlock(lastHeight+uint32(i+1), b, c)
		assert.NoError(t, tStore.WriteBatch())
	}
}

func TestBlockHash(t *testing.T) {
	setup(t)

	sb, _ := tStore.Block(1)

	assert.Equal(t, tStore.BlockHash(0), hash.UndefHash)
	assert.Equal(t, tStore.BlockHash(util.MaxUint32), hash.UndefHash)
	assert.Equal(t, tStore.BlockHash(1), sb.BlockHash)
}

func TestBlockHeight(t *testing.T) {
	setup(t)

	sb, _ := tStore.Block(1)

	assert.Equal(t, tStore.BlockHeight(hash.UndefHash), uint32(0))
	assert.Equal(t, tStore.BlockHeight(hash.GenerateTestHash()), uint32(0))
	assert.Equal(t, tStore.BlockHeight(sb.BlockHash), uint32(1))
}

func TestUnknownTransactionID(t *testing.T) {
	setup(t)

	tx, err := tStore.Transaction(hash.GenerateTestHash())
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestWriteAndClosePeacefully(t *testing.T) {
	setup(t)

	// After closing db, we should not crash
	assert.NoError(t, tStore.Close())
	assert.Error(t, tStore.WriteBatch())
}
func TestRetrieveBlockAndTransactions(t *testing.T) {
	setup(t)

	height, _ := tStore.LastCertificate()
	storedBlock, err := tStore.Block(height)
	assert.NoError(t, err)
	assert.Equal(t, height, storedBlock.Height)
	block := storedBlock.ToBlock()
	for _, trx := range block.Transactions() {
		storedTx, err := tStore.Transaction(trx.ID())
		assert.NoError(t, err)
		assert.Equal(t, storedTx.TxID, trx.ID())
		assert.Equal(t, storedTx.BlockTime, block.Header().UnixTime())
		assert.Equal(t, storedTx.ToTx().ID(), trx.ID())
	}
}

func TestRecentBlockByStamp(t *testing.T) {
	setup(t)

	hash1 := tStore.BlockHash(1)

	h, b := tStore.RecentBlockByStamp(hash.UndefHash.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	h, b = tStore.RecentBlockByStamp(hash1.Stamp())
	assert.Equal(t, h, uint32(1))
	assert.Equal(t, b.Hash(), hash1)

	// Saving more blocks, blocks 11 to 22
	SaveTestBlocks(t, 12)
	hash2 := tStore.BlockHash(2)
	hash14 := tStore.BlockHash(14)
	hash22 := tStore.BlockHash(22)

	// First block should remove from the list
	h, b = tStore.RecentBlockByStamp(hash1.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	h, b = tStore.RecentBlockByStamp(hash2.Stamp())
	assert.Equal(t, h, uint32(2))
	assert.Equal(t, b.Hash(), hash2)

	h, b = tStore.RecentBlockByStamp(hash14.Stamp())
	assert.Equal(t, h, uint32(14))
	assert.Equal(t, b.Hash(), hash14)

	h, b = tStore.RecentBlockByStamp(hash22.Stamp())
	assert.Equal(t, h, uint32(22))
	assert.Equal(t, b.Hash(), hash22)

	// Reopen the store
	tStore.Close()
	s, _ := NewStore(tStore.config, 21)
	tStore = s.(*store)

	h, b = tStore.RecentBlockByStamp(hash2.Stamp())
	assert.Equal(t, h, uint32(2))
	assert.Equal(t, b.Hash(), hash2)

	// Saving one more blocks, block 23
	SaveTestBlocks(t, 1)

	// Second block should remove from the list
	h, b = tStore.RecentBlockByStamp(hash2.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	// Genesis block
	h, b = tStore.RecentBlockByStamp(hash.UndefHash.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)
}
