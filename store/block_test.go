package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func TestBlockStore(t *testing.T) {
	setup(t)

	lastHeight, _ := tStore.LastCertificate()
	b1 := block.GenerateTestBlock(nil, nil)
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
		b2, err := tStore.Block(b1.Hash())
		assert.Error(t, err)
		assert.Nil(t, b2)
	})

	t.Run("Add block, batch write", func(t *testing.T) {
		tStore.SaveBlock(lastHeight+1, b1, c1)
		assert.NoError(t, tStore.WriteBatch())
		sb, err := tStore.Block(b1.Hash())
		assert.NoError(t, err)
		d, _ := b1.Bytes()
		assert.Equal(t, sb.height, lastHeight+1)
		assert.True(t, bytes.Equal(sb.data, d))

		h, cert := tStore.LastCertificate()
		assert.NoError(t, err)
		assert.Equal(t, h, lastHeight+1)
		assert.Equal(t, cert.Hash(), c1.Hash())
	})

	t.Run("Duplicated block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		tStore.SaveBlock(lastHeight, b1, c1)
	})
}

func TestBlockHeightByStamp(t *testing.T) {
	setup(t)

	assert.Zero(t, tStore.BlockHeightByStamp(hash.UndefHash.Stamp()))
	assert.Equal(t, tStore.BlockHeightByStamp(hash.GenerateTestStamp()), int32(-1))

	SaveTestBlocks(t, 11)

	assert.Zero(t, tStore.BlockHeightByStamp(hash.UndefHash.Stamp()))
	assert.Equal(t, tStore.BlockHeightByStamp(tStore.BlockHash(15).Stamp()), int32(15))
}
