package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/zarbchain/zarb-go/types/proposal"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

const (
	blockPrefix       = 0x01
	certificatePrefix = 0x02
	proposalPrefix    = 0x03
)

type key [32]byte

// TODO, create keys without copy?
func blockKey(height int32) key {
	var k key
	k[0] = blockPrefix
	copy(k[1:], util.Int32ToSlice(height))
	return k
}
func certificateKey(height int32) key {
	var k key
	k[0] = certificatePrefix
	copy(k[1:], util.Int32ToSlice(height))
	return k
}
func proposalKey(height int32, round int16) key {
	var k key
	k[0] = proposalPrefix
	copy(k[1:], util.Int32ToSlice(height))
	copy(k[5:], util.Int16ToSlice(round))
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

func (c *Cache) HasBlockInCache(height int32) bool {
	_, ok := c.cache.Get(blockKey(height))
	return ok
}

func (c *Cache) GetBlock(height int32) *block.Block {
	i, ok := c.cache.Get(blockKey(height))
	if ok {
		return i.(*block.Block)
	}

	h := c.state.BlockHash(height)
	if h != hash.UndefHash {
		b := c.state.Block(h)
		if b != nil {
			c.AddBlock(height, b)
			return b
		}
	}

	return nil
}

func (c *Cache) AddBlock(height int32, block *block.Block) {
	c.cache.Add(blockKey(height), block)
	c.AddCertificate(height-1, block.PrevCertificate())
}

func (c *Cache) AddBlocks(height int32, blocks []*block.Block) {
	for _, block := range blocks {
		c.AddBlock(height, block)
		height++
	}
}

func (c *Cache) GetCertificate(height int32) *block.Certificate {
	i, ok := c.cache.Get(certificateKey(height))
	if ok {
		return i.(*block.Certificate)
	}

	return nil
}

func (c *Cache) AddCertificate(height int32, cert *block.Certificate) {
	if cert != nil {
		c.cache.Add(certificateKey(height), cert)
	}
}

func (c *Cache) GetProposal(height int32, round int16) *proposal.Proposal {
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
