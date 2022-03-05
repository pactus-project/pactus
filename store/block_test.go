package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
)

func TestBlockStore(t *testing.T) {
	setup(t)

	b1, _ := block.GenerateTestBlock(nil, nil)
	c1 := block.GenerateTestCertificate(b1.Hash())

	t.Run("Missed block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		tStore.SaveBlock(tLastHeight+2, b1, c1)
	})

	t.Run("Add block, but don't apply batch.", func(t *testing.T) {
		tStore.SaveBlock(tLastHeight+1, b1, c1)
		b2, err := tStore.Block(tLastHeight + 1)
		assert.Error(t, err)
		assert.Nil(t, b2)
	})

	t.Run("Add block and apply batch", func(t *testing.T) {
		tStore.SaveBlock(tLastHeight+1, b1, c1)
		assert.NoError(t, tStore.WriteBatch())
		b2, err := tStore.Block(tLastHeight + 1)
		assert.NoError(t, err)
		assert.Equal(t, b1.Hash(), b2.Hash())

		h, cert, err := tStore.LastCertificate()
		assert.NoError(t, err)
		assert.Equal(t, h, tLastHeight+1)
		assert.Equal(t, cert.Hash(), c1.Hash())
	})

}
