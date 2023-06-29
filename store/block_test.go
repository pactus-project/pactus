package store

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockStore(t *testing.T) {
	td := setup(t)

	lastHeight, _ := td.store.LastCertificate()
	b1 := td.GenerateTestBlock(nil, nil)
	c1 := td.GenerateTestCertificate(b1.Hash())

	t.Run("Missed block, Should panic ", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		td.store.SaveBlock(lastHeight+2, b1, c1)
	})

	t.Run("Add block, don't batch write", func(t *testing.T) {
		td.store.SaveBlock(lastHeight+1, b1, c1)
		b2, err := td.store.Block(lastHeight + 1)
		assert.Error(t, err)
		assert.Nil(t, b2)
	})

	t.Run("Add block, batch write", func(t *testing.T) {
		td.store.SaveBlock(lastHeight+1, b1, c1)
		assert.NoError(t, td.store.WriteBatch())
		sb, err := td.store.Block(lastHeight + 1)
		assert.NoError(t, err)
		d, _ := b1.Bytes()
		assert.Equal(t, sb.Height, lastHeight+1)
		assert.True(t, bytes.Equal(sb.Data, d))

		h, cert := td.store.LastCertificate()
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
		td.store.SaveBlock(lastHeight, b1, c1)
	})
}
