package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util"
)

const (
	blockPrefix       = 0x01
	certificatePrefix = 0x02
)

type key [32]byte

func blockKey(height uint32) key {
	var k key
	k[0] = blockPrefix
	copy(k[1:], util.Uint32ToSlice(height))
	return k
}

func certificateKey(height uint32) key {
	var k key
	k[0] = certificatePrefix
	copy(k[1:], util.Uint32ToSlice(height))
	return k
}

type Cache struct {
	cache *lru.Cache // it's thread safe
}

func NewCache(size int) (*Cache, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &Cache{
		cache: c,
	}, nil
}

func (c *Cache) HasBlockInCache(height uint32) bool {
	_, ok := c.cache.Get(blockKey(height))
	return ok
}

func (c *Cache) GetBlock(height uint32) *block.Block {
	i, ok := c.cache.Get(blockKey(height))
	if ok {
		return i.(*block.Block)
	}

	return nil
}

func (c *Cache) AddBlock(height uint32, block *block.Block) {
	c.cache.Add(blockKey(height), block)
	c.AddCertificate(height-1, block.PrevCertificate())
}

func (c *Cache) GetCertificate(height uint32) *block.Certificate {
	i, ok := c.cache.Get(certificateKey(height))
	if ok {
		return i.(*block.Certificate)
	}

	return nil
}

func (c *Cache) AddCertificate(height uint32, cert *block.Certificate) {
	if cert != nil {
		c.cache.Add(certificateKey(height), cert)
	}
}

func (c *Cache) Len() int {
	return c.cache.Len()
}

func (c *Cache) Clear() {
	c.cache.Purge()
}
