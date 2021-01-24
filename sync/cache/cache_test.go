package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

var tCache *Cache
var tStore *store.MockStore

func setup(t *testing.T) {
	var err error
	tStore = store.MockingStore()
	tCache, err = NewCache(10, tStore)
	assert.NoError(t, err)
}

func TestKeys(t *testing.T) {
	h, _ := crypto.HashFromString("75238478393bfea9e42a59c2cc52876da663ea9acf3873d0a096fd57d61797d4")
	assert.Equal(t, blockKey(1234), key{0x1, 0xd2, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	assert.Equal(t, commitKey(h), key{0x2, 0x75, 0x23, 0x84, 0x78, 0x39, 0x3b, 0xfe, 0xa9, 0xe4, 0x2a, 0x59, 0xc2, 0xcc, 0x52, 0x87, 0x6d, 0xa6, 0x63, 0xea, 0x9a, 0xcf, 0x38, 0x73, 0xd0, 0xa0, 0x96, 0xfd, 0x57, 0xd6, 0x17, 0x97})
	assert.Equal(t, txKey(h), key{0x3, 0x75, 0x23, 0x84, 0x78, 0x39, 0x3b, 0xfe, 0xa9, 0xe4, 0x2a, 0x59, 0xc2, 0xcc, 0x52, 0x87, 0x6d, 0xa6, 0x63, 0xea, 0x9a, 0xcf, 0x38, 0x73, 0xd0, 0xa0, 0x96, 0xfd, 0x57, 0xd6, 0x17, 0x97})
	assert.Equal(t, proposalKey(1234, 3), key{0x4, 0xd2, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
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

	tCache.AddCommit(b1.LastCommit())
	tCache.AddCommit(b2.LastCommit())

	assert.Equal(t, tCache.GetCommit(b1.Header().LastBlockHash()).Hash(), b1.LastCommit().Hash())
	assert.Equal(t, tCache.GetCommit(b2.Header().LastBlockHash()).Hash(), b2.LastCommit().Hash())
	assert.Nil(t, tCache.GetCommit(b3.Header().LastBlockHash()))
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

func TestCacheProposal(t *testing.T) {
	setup(t)

	p1, _ := vote.GenerateTestProposal(100, 0)
	p2, _ := vote.GenerateTestProposal(101, 1)

	tCache.AddProposal(p1)
	tCache.AddProposal(p2)

	assert.Equal(t, tCache.GetProposal(100, 0).Hash(), p1.Hash())
	assert.Equal(t, tCache.GetProposal(101, 1).Hash(), p2.Hash())
	assert.Nil(t, tCache.GetProposal(100, 1))
}

func TestClearCache(t *testing.T) {
	setup(t)

	b, trxs := block.GenerateTestBlock(nil, nil)

	tCache.AddBlock(2, b)
	tCache.AddTransactions(trxs)

	assert.Equal(t, tCache.Len(), 5)
	tCache.Clear()
	assert.Equal(t, tCache.Len(), 0)
	assert.Nil(t, tCache.GetBlock(2))
}

func TestCacheIsFull(t *testing.T) {
	setup(t)

	i := 0
	for ; i < 10; i++ {
		b, _ := block.GenerateTestBlock(nil, nil)
		tCache.AddBlock(i+1, b)
	}

	newBlock, _ := block.GenerateTestBlock(nil, nil)
	tCache.AddBlock(i+1, newBlock)

	assert.NotNil(t, tCache.GetBlock(i+1))
	assert.Nil(t, tCache.GetBlock(1))
}
