package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util"
)

// trust tells where a cached block or certificate comes from.
// An entry is replaced by an entry that is trusted more, never by one that is trusted less.
type trust int

const (
	// trustAnnounced is for an announcement that can't be verified yet. Any peer can send it.
	trustAnnounced trust = iota
	// trustPrevious is for the certificate of the previous block that a trusted block carries.
	trustPrevious
	// trustVerified is for a block or a certificate that is verified,
	// or that comes from the peer we asked for it.
	trustVerified
)

type entry[T any] struct {
	value T
	trust trust
}

type Cache struct {
	blocks *lru.Cache[types.Height, entry[*block.Block]] // it's thread safe
	certs  *lru.Cache[types.Height, entry[*certificate.Certificate]]
}

func NewCache(size int) (*Cache, error) {
	blockCache, err := lru.New[types.Height, entry[*block.Block]](size)
	if err != nil {
		return nil, err
	}

	certCache, err := lru.New[types.Height, entry[*certificate.Certificate]](size)
	if err != nil {
		return nil, err
	}

	return &Cache{
		blocks: blockCache,
		certs:  certCache,
	}, nil
}

func (c *Cache) HasBlockInCache(height types.Height) bool {
	return c.blocks.Contains(height)
}

// HasAnnouncedBlock reports whether the cached block at the given height comes from
// an announcement that was not verified.
func (c *Cache) HasAnnouncedBlock(height types.Height) bool {
	ent, ok := c.blocks.Peek(height)

	return ok && ent.trust == trustAnnounced
}

func (c *Cache) GetBlock(height types.Height) *block.Block {
	ent, ok := c.blocks.Get(height)
	if ok {
		return ent.value
	}

	return nil
}

func blockHeight(blk *block.Block) types.Height {
	prvCert := blk.PrevCertificate()
	if prvCert == nil {
		return 1
	}

	return prvCert.Height() + 1
}

// addCertificate caches the certificate, unless a certificate that is trusted at least
// as much is cached for the same height.
func (c *Cache) addCertificate(cert *certificate.Certificate, level trust) {
	if cert == nil {
		return
	}

	if ent, ok := c.certs.Peek(cert.Height()); ok && ent.trust >= level && level != trustVerified {
		return
	}

	c.certs.Add(cert.Height(), entry[*certificate.Certificate]{value: cert, trust: level})
}

// AddBlock caches a block that is verified, or that comes from the peer we asked for it.
// It replaces the block that is cached for the same height.
func (c *Cache) AddBlock(blk *block.Block) {
	c.blocks.Add(blockHeight(blk), entry[*block.Block]{value: blk, trust: trustVerified})

	// The certificate of the previous block is not verified yet.
	// It should not replace a certificate that is cached already,
	// which can be verified, or made by the peers that we trust more.
	// `AddCertificate` replaces it. It replaces the certificate of an announcement.
	c.addCertificate(blk.PrevCertificate(), trustPrevious)
}

// AddAnnouncedBlock caches the block and the certificate of an announcement that can't be
// verified yet. They never replace what is cached, and what comes later from a trusted
// source replaces them: a peer can't keep a block or a certificate out of the cache.
func (c *Cache) AddAnnouncedBlock(blk *block.Block, cert *certificate.Certificate) {
	height := blockHeight(blk)
	if !c.blocks.Contains(height) {
		c.blocks.Add(height, entry[*block.Block]{value: blk, trust: trustAnnounced})
	}

	c.addCertificate(cert, trustAnnounced)
	c.addCertificate(blk.PrevCertificate(), trustAnnounced)
}

func (c *Cache) GetCertificate(height types.Height) *certificate.Certificate {
	ent, ok := c.certs.Get(height)
	if ok {
		return ent.value
	}

	return nil
}

// AddCertificate caches a certificate that is verified, or that comes from the peer
// we asked for it. It replaces the certificate that is cached for the same height.
func (c *Cache) AddCertificate(cert *certificate.Certificate) {
	c.addCertificate(cert, trustVerified)
}

// RemoveBlock removes the block and certificates at the specified height from the cache.
func (c *Cache) RemoveBlock(height types.Height) {
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
