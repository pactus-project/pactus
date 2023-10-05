package cache

import (
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestAddBlockOne(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	cache, _ := NewCache(10)

	blk1, _ := ts.GenerateTestBlock(1)
	cache.AddBlock(blk1)

	assert.True(t, cache.HasBlockInCache(1))
	assert.Equal(t, blk1, cache.GetBlock(1))
	assert.Nil(t, cache.GetCertificate(0))
}

func TestAddBlocks(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	cache, _ := NewCache(10)

	testHeight := ts.RandHeight()
	blk1, _ := ts.GenerateTestBlock(testHeight)
	cache.AddBlock(blk1)

	assert.True(t, cache.HasBlockInCache(testHeight))
	assert.Equal(t, blk1, cache.GetBlock(testHeight))
	assert.Equal(t, blk1.PrevCertificate(), cache.GetCertificate(testHeight-1))
}

func TestAddCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	cache, _ := NewCache(10)

	testHeight := ts.RandHeight()
	_, cert1 := ts.GenerateTestBlock(testHeight)
	cache.AddCertificate(cert1)

	assert.Equal(t, cert1, cache.GetCertificate(testHeight))
}

func TestClearCache(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	blk, _ := ts.GenerateTestBlock(ts.RandHeight())
	cache.AddBlock(blk)

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
		blk, _ := ts.GenerateTestBlock(uint32(i + 1))
		cache.AddBlock(blk)
	}

	newBlock, _ := ts.GenerateTestBlock(uint32(i + 1))
	cache.AddBlock(newBlock)

	assert.NotNil(t, cache.GetBlock(uint32(i+1)))
	assert.Nil(t, cache.GetBlock(1))
}

func TestAddAgain(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	height := ts.RandHeight()
	firstBlk, _ := ts.GenerateTestBlock(height)
	secondBlk, _ := ts.GenerateTestBlock(height)

	cache.AddBlock(firstBlk)
	assert.Equal(t, firstBlk, cache.GetBlock(height))

	cache.AddBlock(secondBlk)
	assert.Equal(t, secondBlk, cache.GetBlock(height))
}
