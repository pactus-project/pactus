package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockStore(t *testing.T) {
	td := setup(t)

	lastCert := td.store.LastCertificate()
	lastHeight := lastCert.Height()
	nextBlk, nextCrert := td.GenerateTestBlock(lastHeight + 1)
	nextNextBlk, nextNextCrert := td.GenerateTestBlock(lastHeight + 2)

	t.Run("Missed block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		td.store.SaveBlock(nextNextBlk, nextNextCrert)
	})

	t.Run("Add block, don't batch write", func(t *testing.T) {
		td.store.SaveBlock(nextBlk, nextCrert)
		b2, err := td.store.Block(lastHeight + 1)
		assert.Error(t, err)
		assert.Nil(t, b2)
	})

	t.Run("Add block, batch write", func(t *testing.T) {
		td.store.SaveBlock(nextBlk, nextCrert)
		assert.NoError(t, td.store.WriteBatch())

		committedBlock, err := td.store.Block(lastHeight + 1)
		assert.NoError(t, err)
		assert.Equal(t, committedBlock.Height, lastHeight+1)

		d, _ := nextBlk.Bytes()
		assert.True(t, bytes.Equal(committedBlock.Data, d))

		cert := td.store.LastCertificate()
		assert.NoError(t, err)
		assert.Equal(t, cert.Hash(), nextCrert.Hash())
	})

	t.Run("Duplicated block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		td.store.SaveBlock(nextBlk, nextCrert)
	})
}
