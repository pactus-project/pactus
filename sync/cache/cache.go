package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pactus-project/pactus/types/block"
)

type Cache struct {
	blocks *lru.Cache[uint32, *block.Block] // it's thread safe
	certs  *lru.Cache[uint32, *block.Certificate]
}

func NewCache(size int) (*Cache, error) {
	b, err := lru.New[uint32, *block.Block](size)
	if err != nil {
		return nil, err
	}

	c, err := lru.New[uint32, *block.Certificate](size)
	if err != nil {
		return nil, err
	}

	return &Cache{
		blocks: b,
		certs:  c,
	}, nil
}

func (c *Cache) HasBlockInCache(height uint32) bool {
	return c.blocks.Contains(height)
}

func (c *Cache) GetBlock(height uint32) *block.Block {
	block, ok := c.blocks.Get(height)
	if ok {
		return block
	}

	return nil
}

func (c *Cache) AddBlock(height uint32, block *block.Block) {
	c.blocks.Add(height, block)
	c.AddCertificate(height-1, block.PrevCertificate())
}

func (c *Cache) GetCertificate(height uint32) *block.Certificate {
	certificate, ok := c.certs.Get(height)
	if ok {
		return certificate
	}

	return nil
}

func (c *Cache) AddCertificate(height uint32, cert *block.Certificate) {
	if cert != nil {
		c.certs.Add(height, cert)
	}
}

// Len returns the maximum number of items in the blocks and certificates cache.
func (c *Cache) Len() int {
	if c.blocks.Len() > c.certs.Len() {
		return c.blocks.Len()
	}
	return c.certs.Len()
}

func (c *Cache) Clear() {
	c.blocks.Purge()
	c.certs.Purge()
}
