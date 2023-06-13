package cache

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

var tCache *Cache

func setup(t *testing.T) {
	var err error
	tCache, err = NewCache(10)
	assert.NoError(t, err)
}

func TestKeys(t *testing.T) {
	assert.Equal(t, blockKey(1234),
		key{0x1, 0xd2, 0x4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, certificateKey(1234),
		key{0x2, 0xd2, 0x4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func TestCacheBlocks(t *testing.T) {
	setup(t)

	b1 := block.GenerateTestBlock(nil, nil)
	h1 := b1.Hash()
	b2 := block.GenerateTestBlock(nil, &h1)
	testHeight := util.RandUint32(0)

	tCache.AddBlock(testHeight, b1)
	tCache.AddBlock(testHeight+1, b2)

	assert.True(t, tCache.HasBlockInCache(testHeight))
	assert.True(t, tCache.HasBlockInCache(testHeight+1))
	assert.False(t, tCache.HasBlockInCache(testHeight+3))

	assert.NotNil(t, tCache.GetBlock(testHeight))
	assert.NotNil(t, tCache.GetBlock(testHeight+1))
	assert.Nil(t, tCache.GetBlock(testHeight+2))

	assert.Equal(t, tCache.GetBlock(testHeight).Hash(), b1.Hash())
	assert.Equal(t, tCache.GetBlock(testHeight+1).Hash(), b2.Hash())
	assert.Nil(t, tCache.GetCertificate(0))
	assert.Equal(t, tCache.GetCertificate(testHeight).Hash(), b2.PrevCertificate().Hash())
	assert.Nil(t, tCache.GetCertificate(4))
}

func TestClearCache(t *testing.T) {
	setup(t)

	b := block.GenerateTestBlock(nil, nil)

	tCache.AddBlock(2, b)

	assert.Equal(t, tCache.Len(), 2) // block + certificate
	tCache.Clear()
	assert.Equal(t, tCache.Len(), 0)
	assert.Nil(t, tCache.GetBlock(2))
}

func TestCacheIsFull(t *testing.T) {
	setup(t)

	i := int32(0)
	for ; i < 10; i++ {
		b := block.GenerateTestBlock(nil, nil)
		tCache.AddBlock(uint32(i+1), b)
	}

	newBlock := block.GenerateTestBlock(nil, nil)
	tCache.AddBlock(uint32(i+1), newBlock)

	assert.NotNil(t, tCache.GetBlock(uint32(i+1)))
	assert.Nil(t, tCache.GetBlock(1))
}
