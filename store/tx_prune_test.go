package store

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPruneCacheOfAnEmptyCache(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	txStore := newTxStore(nil, 10)

	assert.NotPanics(t, func() { txStore.pruneCache(ts.RandHeight()) })

	id := ts.RandHash()
	txStore.addToCache(id, 1)
	txStore.pruneCache(types.Height(100))
	assert.False(t, txStore.recentTransaction(id), "the transaction is out of the window")
	assert.NotPanics(t, func() { txStore.pruneCache(types.Height(200)) })
}
