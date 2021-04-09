package store

import (
	"testing"

	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/block"
	"github.com/stretchr/testify/assert"
)

func TestLastBlockHeight(t *testing.T) {
	store, _ := newBlockStore(util.TempDirPath())

	assert.False(t, store.hasAnyBlock())

	b1, _ := block.GenerateTestBlock(nil, nil)
	assert.NoError(t, store.saveBlock(b1, 1))

	assert.True(t, store.hasAnyBlock())
}
