package store

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/sortition"
	"github.com/stretchr/testify/assert"
)

func TestBlockStore(t *testing.T) {
	td := setup(t)

	lastCert := td.store.LastCertificate()
	lastHeight := lastCert.Height()
	nextBlk, nextCert := td.GenerateTestBlock(lastHeight + 1)
	nextNextBlk, nextNextCert := td.GenerateTestBlock(lastHeight + 2)

	t.Run("Missed block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		td.store.SaveBlock(nextNextBlk, nextNextCert)
	})

	t.Run("Add block, don't batch write", func(t *testing.T) {
		td.store.SaveBlock(nextBlk, nextCert)
		b2, err := td.store.Block(lastHeight + 1)
		assert.Error(t, err)
		assert.Nil(t, b2)
	})

	t.Run("Add block, batch write", func(t *testing.T) {
		td.store.SaveBlock(nextBlk, nextCert)
		assert.NoError(t, td.store.WriteBatch())

		committedBlock, err := td.store.Block(lastHeight + 1)
		assert.NoError(t, err)
		assert.Equal(t, committedBlock.Height, lastHeight+1)

		d, _ := nextBlk.Bytes()
		assert.True(t, bytes.Equal(committedBlock.Data, d))

		cert := td.store.LastCertificate()
		assert.NoError(t, err)
		assert.Equal(t, cert.Hash(), nextCert.Hash())
	})

	t.Run("Duplicated block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		td.store.SaveBlock(nextBlk, nextCert)
	})

	t.Run("Should not be old sortition seed in cache plus check last sortition seed be in cache",
		func(t *testing.T) {
			var oldSortitionSeed sortition.VerifiableSeed
			var lastSortitionSeed sortition.VerifiableSeed
			lastHeight = td.store.LastCertificate().Height()
			generateBlkCount := 100

			for i := 0; i <= generateBlkCount; i++ {
				lastHeight++
				nNextBlk, nNextCert := td.GenerateTestBlock(lastHeight)

				td.store.SaveBlock(nNextBlk, nNextCert)
				assert.NoError(t, td.store.WriteBatch())

				if i == 0 {
					oldSortitionSeed = nNextBlk.Header().SortitionSeed()
				}

				if i == generateBlkCount {
					lastSortitionSeed = nNextBlk.Header().SortitionSeed()
				}
			}

			// check old sortition seed doesn't exist in cache
			for _, seed := range td.store.blockStore.sortitionSeedCache.All() {
				assert.NotEqual(t, &oldSortitionSeed, seed)
			}

			// last sortition seed should exist at last index of cache
			assert.Equal(t, &lastSortitionSeed, td.store.blockStore.sortitionSeedCache.Last().SecondElement)
		})
}
