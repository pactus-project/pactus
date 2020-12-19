package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

const (
	blockPrefix  = 0x01
	commitPrefix = 0x02
	txPrefix     = 0x03
)

func blockKey(height int) [9]byte {
	var k [9]byte
	k[0] = blockPrefix
	copy(k[1:], util.IntToSlice(height))
	return k
}
func commitKey(hash crypto.Hash) [33]byte {
	var k [33]byte
	k[0] = commitPrefix
	copy(k[1:], hash.RawBytes())
	return k
}
func txKey(id crypto.Hash) [33]byte {
	var k [33]byte
	k[0] = txPrefix
	copy(k[1:], id.RawBytes())
	return k
}

type Cache struct {
	cache *lru.ARCCache
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

	// TODO: get block commit from store. Good idea?

	return nil
}

func (c *Cache) AddCommit(blockhash crypto.Hash, commit *block.Commit) {
	c.cache.Add(commitKey(blockhash), commit)
}

func (c *Cache) GetTransaction(id crypto.Hash) *tx.Tx {
	i, ok := c.cache.Get(txKey(id))
	if ok {
		return i.(*tx.Tx)
	}

	ct, err := c.store.Transaction(id)
	if err == nil {
		c.cache.Add(txKey(id), ct.Tx)
		return ct.Tx
	}

	return nil
}
func (c *Cache) AddTransaction(trx *tx.Tx) {
	c.cache.Add(txKey(trx.ID()), trx)
}

func (c *Cache) Len() int {
	return c.cache.Len()
}
