package cache

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	height := types.Height(0)
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

// The certificate that comes with a block is not verified.
// It must not replace the certificate that is cached already,
// otherwise a block of the next height can poison the certificate of the current one.
func TestAddBlockDoesNotReplaceCertificate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	height := ts.RandHeight()
	cert := ts.GenerateTestCertificate(height - 1)
	cache.AddCertificate(cert)

	blk, _ := ts.GenerateTestBlock(height)
	require.NotEqual(t, cert, blk.PrevCertificate())

	cache.AddBlock(blk)

	assert.Equal(t, blk, cache.GetBlock(height))
	assert.Equal(t, cert, cache.GetCertificate(height-1), "the cached certificate should be kept")
}

// A certificate can be replaced explicitly, for example by a verified one.
func TestAddCertificateReplaces(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, _ := NewCache(10)

	height := ts.RandHeight()
	blk, _ := ts.GenerateTestBlock(height)
	cache.AddBlock(blk)
	require.Equal(t, blk.PrevCertificate(), cache.GetCertificate(height-1))

	verified := ts.GenerateTestCertificate(height - 1)
	cache.AddCertificate(verified)

	assert.Equal(t, verified, cache.GetCertificate(height-1))
}

func TestAddNilCertificate(t *testing.T) {
	cache, _ := NewCache(10)

	cache.AddCertificate(nil)

	assert.Equal(t, 0, cache.Len())
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

func TestNewCacheWithInvalidSize(t *testing.T) {
	for _, size := range []int{0, -1} {
		c, err := NewCache(size)

		require.Error(t, err, "size %d", size)
		assert.Nil(t, c)
	}
}

func TestAnnouncedBlockDoesNotMaskTrustedData(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cache, err := NewCache(10)
	require.NoError(t, err)

	height := ts.RandHeight()
	fakeBlk, fakeCert := ts.GenerateTestBlock(height)
	realBlk, _ := ts.GenerateTestBlock(height)
	nextBlk, _ := ts.GenerateTestBlock(height + 1)
	realCert := nextBlk.PrevCertificate()
	require.Equal(t, height, realCert.Height())

	// An announcement that can't be verified is cached.
	cache.AddAnnouncedBlock(fakeBlk, fakeCert)
	assert.True(t, cache.HasAnnouncedBlock(height))
	assert.Equal(t, fakeBlk, cache.GetBlock(height))
	assert.Equal(t, fakeCert, cache.GetCertificate(height))

	// Another announcement doesn't replace it.
	otherBlk, otherCert := ts.GenerateTestBlock(height)
	cache.AddAnnouncedBlock(otherBlk, otherCert)
	assert.Equal(t, fakeBlk, cache.GetBlock(height))
	assert.Equal(t, fakeCert, cache.GetCertificate(height))

	// The blocks that we asked for replace the block and its certificate.
	cache.AddBlock(realBlk)
	cache.AddBlock(nextBlk)
	assert.False(t, cache.HasAnnouncedBlock(height))
	assert.Equal(t, realBlk, cache.GetBlock(height))
	assert.Equal(t, realCert, cache.GetCertificate(height))

	// And an announcement doesn't replace them back.
	cache.AddAnnouncedBlock(fakeBlk, fakeCert)
	assert.Equal(t, realBlk, cache.GetBlock(height))
	assert.Equal(t, realCert, cache.GetCertificate(height))

	// The certificate of the previous block of a trusted block doesn't replace
	// a certificate that was added on its own.
	verifiedCert := ts.GenerateTestCertificate(height)
	cache.AddCertificate(verifiedCert)
	cache.AddBlock(nextBlk)
	assert.Equal(t, verifiedCert, cache.GetCertificate(height))
}

func TestCertificateTrustLevels(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("a certificate added on its own replaces the one added before", func(t *testing.T) {
		cache, err := NewCache(10)
		require.NoError(t, err)

		height := ts.RandHeight()
		oldCert := ts.GenerateTestCertificate(height)
		newCert := ts.GenerateTestCertificate(height)

		cache.AddCertificate(oldCert)
		cache.AddCertificate(newCert)
		assert.Equal(t, newCert, cache.GetCertificate(height))
	})

	t.Run("the previous certificate of an announcement is replaced by the one of a trusted block", func(t *testing.T) {
		cache, err := NewCache(10)
		require.NoError(t, err)

		height := ts.RandHeight()
		announced, _ := ts.GenerateTestBlock(height + 1)
		trusted, _ := ts.GenerateTestBlock(height + 1)
		require.NotEqual(t, announced.PrevCertificate(), trusted.PrevCertificate())

		// The certificate of the announced block is for another height.
		cache.AddAnnouncedBlock(announced, ts.GenerateTestCertificate(height+1))
		assert.Equal(t, announced.PrevCertificate(), cache.GetCertificate(height))

		cache.AddBlock(trusted)
		assert.Equal(t, trusted.PrevCertificate(), cache.GetCertificate(height))
	})
}
