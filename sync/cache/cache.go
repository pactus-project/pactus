package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
)

const (
	blockPrefix    = 0x01
	commitPrefix   = 0x02
	txPrefix       = 0x03
	proposalPrefix = 0x04
)

type key [32]byte

func blockKey(height int) key {
	var k key
	k[0] = blockPrefix
	copy(k[1:], util.IntToSlice(height))
	return k
}
func commitKey(hash crypto.Hash) key {
	var k key
	k[0] = commitPrefix
	copy(k[1:], hash.RawBytes())
	return k
}
func txKey(id tx.ID) key {
	var k key
	k[0] = txPrefix
	copy(k[1:], id.RawBytes())
	return k
}
func proposalKey(height, round int) key {
	var k key
	k[0] = proposalPrefix
	copy(k[1:], util.IntToSlice(height))
	copy(k[16:], util.IntToSlice(round))
	return k
}

type Cache struct {
	cache *lru.ARCCache // it's thread safe
	store store.StoreReader
}

func NewCache(size int, store store.StoreReader) (*Cache, error) {
	c, err := lru.NewARC(size)
	if err != nil {
		return nil, err
	}
	return &Cache{
		cache: c,
		store: store,
	}, nil
}

func (c *Cache) GetBlock(height int) *block.Block {
	i, ok := c.cache.Get(blockKey(height))
	if ok {
		return i.(*block.Block)
	}

	b, err := c.store.Block(height)
	if err == nil {
		c.cache.Add(blockKey(height), b)
		return b
	}

	return nil
}

func (c *Cache) AddBlock(height int, block *block.Block) {
	c.cache.Add(blockKey(height), block)
}

func (c *Cache) GetCommit(blockhash crypto.Hash) *block.Commit {
	i, ok := c.cache.Get(commitKey(blockhash))
	if ok {
		return i.(*block.Commit)
	}

	return nil
}

func (c *Cache) AddCommit(commit *block.Commit) {
	if commit != nil {
		c.cache.Add(commitKey(commit.BlockHash()), commit)
	}
}

func (c *Cache) GetTransaction(id tx.ID) *tx.Tx {
	i, ok := c.cache.Get(txKey(id))
	if ok {
		return i.(*tx.Tx)
	}

	ct, err := c.store.Transaction(id)
	if err == nil {
		c.cache.Add(txKey(id), ct.Tx)
		return ct.Tx
	}

	// Should we check txpool?
	// No, because transaction in txpool should be in cache.
	// TODO: write tests for me

	return nil
}
func (c *Cache) AddTransaction(trx *tx.Tx) {
	c.cache.Add(txKey(trx.ID()), trx)
}

func (c *Cache) AddTransactions(trxs []*tx.Tx) {
	for _, trx := range trxs {
		c.AddTransaction(trx)
	}
}

func (c *Cache) GetProposal(height, round int) *vote.Proposal {
	i, ok := c.cache.Get(proposalKey(height, round))
	if ok {
		return i.(*vote.Proposal)
	}

	return nil
}

func (c *Cache) AddProposal(p *vote.Proposal) {
	c.cache.Add(proposalKey(p.Height(), p.Round()), p)
}
func (c *Cache) Len() int {
	return c.cache.Len()
}

func (c *Cache) Clear() {
	c.cache.Purge()
}
