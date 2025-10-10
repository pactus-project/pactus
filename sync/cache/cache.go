package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util"
)

type Cache struct {
	blocks *lru.Cache[uint32, *block.Block] // it's thread safe
	certs  *lru.Cache[uint32, *certificate.Certificate]
}

func NewCache(size int) (*Cache, error) {
	blockCache, err := lru.New[uint32, *block.Block](size)
	if err != nil {
		return nil, err
	}

	certCache, err := lru.New[uint32, *certificate.Certificate](size)
	if err != nil {
		return nil, err
	}

	return &Cache{
		blocks: blockCache,
		certs:  certCache,
	}, nil
}

func (c *Cache) HasBlockInCache(height uint32) bool {
	return c.blocks.Contains(height)
}

func (c *Cache) GetBlock(height uint32) *block.Block {
	blk, ok := c.blocks.Get(height)
	if ok {
		return blk
	}

	return nil
}

func (c *Cache) AddBlock(blk *block.Block) {
	prvCert := blk.PrevCertificate()
	if prvCert == nil {
		c.blocks.Add(1, blk)
	} else {
		c.blocks.Add(prvCert.Height()+1, blk)
		c.certs.Add(prvCert.Height(), prvCert)
	}
}

func (c *Cache) GetCertificate(height uint32) *certificate.Certificate {
	cert, ok := c.certs.Get(height)
	if ok {
		return cert
	}

	return nil
}

func (c *Cache) AddCertificate(cert *certificate.Certificate) {
	if cert != nil {
		c.certs.Add(cert.Height(), cert)
	}
}

// RemoveBlock removes the block and certificates at the specified height from the cache.
func (c *Cache) RemoveBlock(height uint32) {
	c.blocks.Remove(height)
	c.certs.Remove(height)
}

// Len returns the maximum number of items in the blocks and certificates cache.
func (c *Cache) Len() int {
	return util.Max(c.blocks.Len(), c.certs.Len())
}

func (c *Cache) Clear() {
	c.blocks.Purge()
	c.certs.Purge()
}
