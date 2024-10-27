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

	assert.Equal(t, 1, cache.Len())
	cache.Clear()
	assert.Equal(t, 0, cache.Len())
	assert.Nil(t, cache.GetBlock(2))
}

func TestCacheIsFull(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	height := uint32(0)
	for ; height < 10; height++ {
		blk, _ := ts.GenerateTestBlock(height + 1)
		cache.AddBlock(blk)
	}

	newBlock, _ := ts.GenerateTestBlock(height + 1)
	cache.AddBlock(newBlock)

	assert.NotNil(t, cache.GetBlock(height+1))
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

func TestRemoveBlock(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	height := ts.RandHeight()
	blk1, _ := ts.GenerateTestBlock(height)
	blk2, _ := ts.GenerateTestBlock(height + 1)
	cache.AddBlock(blk1)
	cache.AddBlock(blk2)

	cache.RemoveBlock(height)
	assert.Nil(t, cache.GetBlock(height))
	assert.Nil(t, cache.GetCertificate(height))
}
