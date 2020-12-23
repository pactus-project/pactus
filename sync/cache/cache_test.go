package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
)

var tCache *Cache
var tStore *store.MockStore

func setup(t *testing.T) {
	var err error
	tStore = store.NewMockStore()
	tCache, err = NewCache(10, tStore)
	assert.NoError(t, err)
}

func TestCacheBlock(t *testing.T) {
	setup(t)

	b1, _ := block.GenerateTestBlock(nil, nil)
	b2, _ := block.GenerateTestBlock(nil, nil)

	tStore.Blocks[1] = b1
	tCache.AddBlock(2, b2)

	assert.Equal(t, tCache.GetBlock(1).Hash(), b1.Hash())
	assert.Equal(t, tCache.GetBlock(2).Hash(), b2.Hash())
	assert.Nil(t, tCache.GetBlock(3))
}

func TestCacheCommit(t *testing.T) {
	setup(t)

	b1, _ := block.GenerateTestBlock(nil, nil)
	b2, _ := block.GenerateTestBlock(nil, nil)
	b3, _ := block.GenerateTestBlock(nil, nil)

	tStore.Blocks[1] = b1
	tStore.Blocks[2] = b2
	tCache.AddCommit(b1.Hash(), b2.LastCommit())
	tCache.AddCommit(b2.Hash(), b3.LastCommit())

	assert.Equal(t, tCache.GetCommit(b1.Hash()).Hash(), b2.LastCommit().Hash())
	assert.Equal(t, tCache.GetCommit(b2.Hash()).Hash(), b3.LastCommit().Hash())
	assert.Nil(t, tCache.GetCommit(b3.Hash()))
}

func TestCacheTx(t *testing.T) {
	setup(t)

	trx1, _ := tx.GenerateTestSendTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, _ := tx.GenerateTestSendTx()

	tStore.Transactions[trx1.ID()] = &tx.CommittedTx{Tx: trx1}
	tCache.AddTransaction(trx2)

	assert.Equal(t, tCache.GetTransaction(trx1.ID()).ID(), trx1.ID())
	assert.Equal(t, tCache.GetTransaction(trx2.ID()).ID(), trx2.ID())
	assert.Nil(t, tCache.GetTransaction(trx3.ID()))
}
