package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pactus-project/pactus/types/block"
)

type Cache struct {
	block       *lru.Cache[uint32, *block.Block] // it's thread safe
	certificate *lru.Cache[uint32, *block.Certificate]
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
		block:       b,
		certificate: c,
	}, nil
}

func (c *Cache) HasBlockInCache(height uint32) bool {
	return c.block.Contains(height)
}

func (c *Cache) GetBlock(height uint32) *block.Block {
	block, ok := c.block.Get(height)
	if ok {
		return block
	}

	return nil
}

func (c *Cache) AddBlock(height uint32, block *block.Block) {
	c.block.Add(height, block)
	c.AddCertificate(height-1, block.PrevCertificate())
}

func (c *Cache) GetCertificate(height uint32) *block.Certificate {
	certificate, ok := c.certificate.Get(height)
	if ok {
		return certificate
	}

	return nil
}

func (c *Cache) AddCertificate(height uint32, cert *block.Certificate) {
	if cert != nil {
		c.certificate.Add(height, cert)
	}
}

func (c *Cache) Len() int {
	if c.block.Len() > c.certificate.Len() {
		return c.block.Len()
	}
	return c.certificate.Len()
}

func (c *Cache) Clear() {
	c.block.Purge()
	c.certificate.Purge()
}
