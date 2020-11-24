package txpool

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)


type TxPoolReader interface {
	PendingTx(hash crypto.Hash) *tx.Tx
	HasTx(hash crypto.Hash) bool
	Size() int

	Fingerprint() string
}

type TxPool interface {
	TxPoolReader

	AppendTxs(txs []tx.Tx)
	AppendTx(tx tx.Tx)
	AppendTxAndBroadcast(trx tx.Tx)
	RemoveTx(hash crypto.Hash) *tx.Tx
}
