package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	simplemerkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/tx"
)

func TestTxsMerkle(t *testing.T) {
	b, txs := GenerateTestBlock(nil, nil)

	data := make([]tx.ID, len(txs))
	for i, tx := range txs {
		data[i] = tx.ID()
	}
	merkle := simplemerkle.NewTreeFromHashes(data)
	assert.Equal(t, b.Header().TxIDsHash(), merkle.Root())
}

func TestAppendAndPrepend(t *testing.T) {
	ids := NewTxIDs()
	h1 := hash.GenerateTestHash()
	h2 := hash.GenerateTestHash()
	h3 := hash.GenerateTestHash()
	h4 := hash.GenerateTestHash()
	ids.Append(h2)
	ids.Append(h3)
	ids.Prepend(h1)
	ids.Append(h4)

	assert.Equal(t, ids.data.IDs, []hash.Hash{h1, h2, h3, h4})
}
