package tx_test

import (
	"bytes"
	"sync"
	"testing"

	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

// A transaction is shared between the services of the node.
// Run with the race detector: the first calls of ID and BasicCheck store their result.
func TestConcurrentIDAndBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx := ts.GenerateTestTransferTx()
	want := trx.ID()

	decoded := ts.GenerateTestTransferTx()
	data, err := trx.Bytes()
	assert.NoError(t, err)
	assert.NoError(t, decoded.Decode(bytes.NewReader(data)))

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			assert.Equal(t, want, decoded.ID())
			assert.NoError(t, decoded.BasicCheck())
		})
	}
	wg.Wait()
}
