package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestBlockStore(t *testing.T) {
	setup(t)

	lastHeight, _ := tStore.LastCertificate()
	b1, _ := block.GenerateTestBlock(nil, nil)
	c1 := block.GenerateTestCertificate(b1.Hash())

	t.Run("Missed block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		tStore.SaveBlock(lastHeight+2, b1, c1)
	})

	t.Run("Add block, don't batch write", func(t *testing.T) {
		tStore.SaveBlock(lastHeight+1, b1, c1)
		b2, err := tStore.Block(lastHeight + 1)
		assert.Error(t, err)
		assert.Nil(t, b2)
	})

	t.Run("Add block, batch write", func(t *testing.T) {
		tStore.SaveBlock(lastHeight+1, b1, c1)
		assert.NoError(t, tStore.WriteBatch())
		b2, err := tStore.Block(lastHeight + 1)
		assert.NoError(t, err)
		assert.Equal(t, b1.Hash(), b2.Hash())

		h, cert := tStore.LastCertificate()
		assert.NoError(t, err)
		assert.Equal(t, h, lastHeight+1)
		assert.Equal(t, cert.Hash(), c1.Hash())
	})
}

func TestBlockHeightByStamp(t *testing.T) {
	setup(t)

	assert.Equal(t, tStore.BlockHeightByStamp(hash.UndefHash.Stamp()), 0)
	assert.Equal(t, tStore.BlockHeightByStamp(hash.GenerateTestStamp()), -1)
}
