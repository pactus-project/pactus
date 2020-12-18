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

func TestAppendAndPrepend(t *testing.T) {
	ids := NewTxIDs()
	h1 := crypto.GenerateTestHash()
	h2 := crypto.GenerateTestHash()
	h3 := crypto.GenerateTestHash()
	h4 := crypto.GenerateTestHash()
	ids.Append(h2)
	ids.Append(h3)
	ids.Prepend(h1)
	ids.Append(h4)

	assert.Equal(t, ids.data.IDs, []crypto.Hash{h1, h2, h3, h4})
}
