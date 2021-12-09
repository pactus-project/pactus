package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

const (
	blockPrefix       = 0x01
	certificatePrefix = 0x02
	txPrefix          = 0x03
	proposalPrefix    = 0x04
)

type key [32]byte

func blockKey(height int) key {
	var k key
	k[0] = blockPrefix
	copy(k[1:], util.IntToSlice(height))
	return k
}
func certificateKey(hash hash.Hash) key {
	var k key
	k[0] = certificatePrefix
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
	state state.Facade
}

func NewCache(size int, state state.Facade) (*Cache, error) {
	c, err := lru.NewARC(size)
	if err != nil {
		return nil, err
	}
	return &Cache{
		cache: c,
		state: state,
	}, nil
}

func (c *Cache) GetBlock(height int) *block.Block {
	i, ok := c.cache.Get(blockKey(height))
	if ok {
		return i.(*block.Block)
	}

	b := c.state.Block(height)
	if b != nil {
		c.AddBlock(height, b)
		return b
	}

	return nil
}

func (c *Cache) AddBlock(height int, block *block.Block) {
	c.cache.Add(blockKey(height), block)
	c.AddCertificate(block.PrevCertificate())
}

func (c *Cache) AddBlocks(height int, blocks []*block.Block) {
	for _, block := range blocks {
		c.AddBlock(height, block)
		height++
	}
}

func (c *Cache) GetCertificate(blockhash hash.Hash) *block.Certificate {
	i, ok := c.cache.Get(certificateKey(blockhash))
	if ok {
		return i.(*block.Certificate)
	}

	return nil
}

func (c *Cache) AddCertificate(cert *block.Certificate) {
	if cert != nil {
		c.cache.Add(certificateKey(cert.BlockHash()), cert)
	}
}

func (c *Cache) GetTransaction(id tx.ID) *tx.Tx {
	i, ok := c.cache.Get(txKey(id))
	if ok {
		return i.(*tx.Tx)
	}

	pendingTrx := c.state.PendingTx(id)
	if pendingTrx != nil {
		c.cache.Add(txKey(id), pendingTrx)
		return pendingTrx
	}

	cacheTrx := c.state.Transaction(id)
	if cacheTrx != nil {
		c.cache.Add(txKey(id), cacheTrx)
		return cacheTrx
	}

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

func (c *Cache) GetProposal(height, round int) *proposal.Proposal {
	i, ok := c.cache.Get(proposalKey(height, round))
	if ok {
		return i.(*proposal.Proposal)
	}

	return nil
}

func (c *Cache) AddProposal(p *proposal.Proposal) {
	c.cache.Add(proposalKey(p.Height(), p.Round()), p)
}
func (c *Cache) Len() int {
	return c.cache.Len()
}

func (c *Cache) Clear() {
	c.cache.Purge()
}
