package block

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
)

func TestTxsMerkle(t *testing.T) {
	b, txs := GenerateTestBlock(nil)

	data := make([]crypto.Hash, len(txs))
	for i, tx := range txs {
		data[i] = tx.ID()
	}
	merkle := simpleMerkle.NewTreeFromHashes(data)
	assert.Equal(t, b.Header().TxIDsHash(), merkle.Root())
}
