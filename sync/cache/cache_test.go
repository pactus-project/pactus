package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
)

var tCache *Cache
var tState *state.MockState

func setup(t *testing.T) {
	var err error
	tState = state.MockingState()
	tCache, err = NewCache(10, tState)
	assert.NoError(t, err)
}

func TestKeys(t *testing.T) {
	assert.Equal(t, blockKey(1234), key{0x1, 0xd2, 0x4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, certificateKey(1234), key{0x2, 0xd2, 0x4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func TestCacheBlocks(t *testing.T) {
	setup(t)

	b1 := block.GenerateTestBlock(nil, &hash.UndefHash)
	h1 := b1.Hash()
	b2 := block.GenerateTestBlock(nil, &h1)
	h2 := b2.Hash()
	b3 := block.GenerateTestBlock(nil, &h2)

	tState.TestStore.SaveBlock(1, b1, block.GenerateTestCertificate(b1.Hash()))
	tCache.AddBlocks(2, []*block.Block{b2, b3})

	assert.False(t, tCache.HasBlockInCache(1), "Block 1 is not cached")
	assert.True(t, tCache.HasBlockInCache(2))
	assert.True(t, tCache.HasBlockInCache(3))
	assert.False(t, tCache.HasBlockInCache(4))

	assert.NotNil(t, tCache.GetBlock(1))
	assert.NotNil(t, tCache.GetBlock(2))
	assert.NotNil(t, tCache.GetBlock(3))
	assert.Nil(t, tCache.GetBlock(4))

	assert.Equal(t, tCache.GetBlock(1).Hash(), b1.Hash())
	assert.Equal(t, tCache.GetBlock(2).Hash(), b2.Hash())
	assert.Nil(t, tCache.GetCertificate(0))
	assert.Equal(t, tCache.GetCertificate(1).Hash(), b2.PrevCertificate().Hash())
	assert.Equal(t, tCache.GetCertificate(2).Hash(), b3.PrevCertificate().Hash())
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
