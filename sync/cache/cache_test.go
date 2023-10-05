package cache

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestCacheBlocks(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	b1 := ts.GenerateTestBlock()
	b2 := ts.GenerateTestBlock()
	testHeight := ts.RandHeight()

	cache.AddBlock(testHeight, b1)
	cache.AddBlock(testHeight+1, b2)

	assert.True(t, cache.HasBlockInCache(testHeight))
	assert.True(t, cache.HasBlockInCache(testHeight+1))
	assert.False(t, cache.HasBlockInCache(testHeight+3))

	assert.NotNil(t, cache.GetBlock(testHeight))
	assert.NotNil(t, cache.GetBlock(testHeight+1))
	assert.Nil(t, cache.GetBlock(testHeight+2))

	assert.Equal(t, cache.GetBlock(testHeight).Hash(), b1.Hash())
	assert.Equal(t, cache.GetBlock(testHeight+1).Hash(), b2.Hash())
	assert.Nil(t, cache.GetCertificate(0))
	assert.Equal(t, cache.GetCertificate(testHeight).Hash(), b2.PrevCertificate().Hash())
	assert.Nil(t, cache.GetCertificate(4))
}

func TestClearCache(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	b := ts.GenerateTestBlock()

	cache.AddBlock(2, b)

	assert.Equal(t, cache.Len(), 1)
	cache.Clear()
	assert.Equal(t, cache.Len(), 0)
	assert.Nil(t, cache.GetBlock(2))
}

func TestCacheIsFull(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	i := int32(0)
	for ; i < 10; i++ {
		b := ts.GenerateTestBlock()
		cache.AddBlock(uint32(i+1), b)
	}

	newBlock := ts.GenerateTestBlock()
	cache.AddBlock(uint32(i+1), newBlock)

	assert.NotNil(t, cache.GetBlock(uint32(i+1)))
	assert.Nil(t, cache.GetBlock(1))
}

func TestAddAgain(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	height := ts.RandHeight()
	blk1 := ts.GenerateTestBlock()
	blk2 := ts.GenerateTestBlock()

	cache.AddBlock(height, blk1)
	assert.Equal(t, blk1, cache.GetBlock(height))

	cache.AddBlock(height, blk2)
	assert.Equal(t, blk2, cache.GetBlock(height))
}
