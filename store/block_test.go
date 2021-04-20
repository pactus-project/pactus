package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
)

func TestBlockStore(t *testing.T) {
	store, _ := NewStore(TestConfig())

	assert.False(t, store.HasAnyBlock())

	t.Run("Add block, but not write batch.", func(t *testing.T) {
		b1, _ := block.GenerateTestBlock(nil, nil)
		store.SaveBlock(1, b1)
		assert.False(t, store.HasAnyBlock())
	})

	t.Run("Add block and write batch", func(t *testing.T) {
		b1, _ := block.GenerateTestBlock(nil, nil)
		store.SaveBlock(1, b1)
		assert.NoError(t, store.WriteBatch())
		assert.True(t, store.HasAnyBlock())
	})

}
