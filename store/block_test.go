package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/util"
)

func TestLastBlockHeight(t *testing.T) {
	store, _ := newBlockStore(util.TempDirPath())

	assert.False(t, store.hasAnyBlock())

	b1, _ := block.GenerateTestBlock(nil, nil)
	assert.NoError(t, store.saveBlock(b1, 1))

	assert.True(t, store.hasAnyBlock())
}
