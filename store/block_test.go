package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockStore(t *testing.T) {
	td := setup(t, nil)

	lastCert := td.store.LastCertificate()
	lastHeight := lastCert.Height()
	nextBlk, nextCert := td.GenerateTestBlock(lastHeight + 1)

	t.Run("Add block, don't batch write", func(t *testing.T) {
		td.store.SaveBlock(nextBlk, nextCert)
		b2, err := td.store.Block(lastHeight + 1)
		assert.Error(t, err)
		assert.Nil(t, b2)
	})

	t.Run("Add block, batch write", func(t *testing.T) {
		td.store.SaveBlock(nextBlk, nextCert)
		assert.NoError(t, td.store.WriteBatch())

		cBlk, err := td.store.Block(lastHeight + 1)
		assert.NoError(t, err)
		assert.Equal(t, lastHeight+1, cBlk.Height)

		d, _ := nextBlk.Bytes()
		assert.True(t, bytes.Equal(cBlk.Data, d))

		cert := td.store.LastCertificate()
		assert.NoError(t, err)
		assert.Equal(t, nextCert.Hash(), cert.Hash())
	})
}

func TestSortitionSeed(t *testing.T) {
	conf := testConfig()
	conf.SeedCacheWindow = 7

	td := setup(t, conf)
	lastHeight := td.store.LastCertificate().Height()

	t.Run("Test height zero", func(t *testing.T) {
		assert.Nil(t, td.store.SortitionSeed(0))
	})

	t.Run("Test non existing height", func(t *testing.T) {
		assert.Nil(t, td.store.SortitionSeed(lastHeight+1))
	})

	t.Run("Test not cached height", func(t *testing.T) {
		assert.Nil(t, td.store.SortitionSeed(3))
	})

	t.Run("OK", func(t *testing.T) {
		rndInt := td.RandUint32Max(conf.SeedCacheWindow)
		rndInt += lastHeight - conf.SeedCacheWindow + 1

		committedBlk, _ := td.store.Block(rndInt)
		blk, _ := committedBlk.ToBlock()
		expectedSortition := blk.Header().SortitionSeed()
		assert.Equal(t, &expectedSortition, td.store.SortitionSeed(rndInt))
	})
}
