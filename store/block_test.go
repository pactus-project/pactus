package store

import (
	"testing"

	"github.com/zarbchain/zarb-go/util"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/tx"

	"github.com/stretchr/testify/assert"
)

func TestRetreiveBlockAndTransactions(t *testing.T) {
	conf := TestConfig()
	store, err := NewStore(conf)
	assert.NoError(t, err)

	b, txs := block.GenerateTestBlock(nil)
	h := util.RandInt(10000)
	err = store.SaveBlock(b, h)
	assert.NoError(t, err)

	for _, trx := range txs {
		r := trx.GenerateReceipt(tx.Ok, b.Hash())
		ctrx := tx.CommittedTx{Tx: trx, Receipt: r}
		err = store.SaveTransaction(ctrx)
		assert.NoError(t, err)
	}

	b2, h2, err := store.BlockByHash(b.Hash())
	assert.NoError(t, err)
	assert.Equal(t, b.Hash(), b2.Hash())
	bz1, _ := b.Encode()
	bz2, _ := b2.Encode()
	assert.Equal(t, bz1, bz2)
	assert.Equal(t, h, h2)

	for _, trx := range txs {
		r := trx.GenerateReceipt(tx.Ok, b.Hash())
		ctrx2, err := store.Transaction(trx.Hash())
		assert.NoError(t, err)

		assert.Equal(t, trx.Hash(), ctrx2.Tx.Hash())
		assert.Equal(t, r, ctrx2.Receipt)
	}
}
