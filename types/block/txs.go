package block

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/simplemerkle"
)

type Txs []*tx.Tx

func NewTxs() Txs {
	return make([]*tx.Tx, 0)
}

func (txs *Txs) Append(trx *tx.Tx) {
	*txs = append(*txs, trx)
}

func (txs *Txs) Prepend(trx *tx.Tx) {
	*txs = append(*txs, nil)
	copy((*txs)[1:], (*txs)[0:])
	(*txs)[0] = trx
}

func (txs *Txs) Remove(i int) {
	// https://github.com/golang/go/wiki/SliceTricks#delete
	copy((*txs)[i:], (*txs)[i+1:])
	(*txs)[txs.Len()-1] = nil
	*txs = (*txs)[:txs.Len()-1]
}

func (txs Txs) Root() hash.Hash {
	hashes := make([]hash.Hash, txs.Len())
	for i, trx := range txs {
		hashes[i] = trx.ID()
	}
	merkle := simplemerkle.NewTreeFromHashes(hashes)

	return merkle.Root()
}

func (txs Txs) IsEmpty() bool {
	return txs.Len() == 0
}

func (txs Txs) Len() int {
	return len(txs)
}

func (txs Txs) Get(i int) *tx.Tx {
	return txs[i]
}

func (txs Txs) Subsidy() *tx.Tx {
	return txs[0]
}
