package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
)

func TestLastBlockHeight(t *testing.T) {
	store, _ := NewStore(TestConfig())

	assert.False(t, store.HasAnyBlock())

	b1, _ := block.GenerateTestBlock(nil, nil)
	assert.NoError(t, store.SaveBlock(1, b1))

	assert.True(t, store.HasAnyBlock())
}
