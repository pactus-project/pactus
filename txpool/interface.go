package txpool

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sandbox"
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

	SetSandbox(sandbox sandbox.Sandbox)
	AppendTxs(txs []*tx.Tx)
	AppendTx(tx *tx.Tx) error
	AppendTxAndBroadcast(trx *tx.Tx) error
	RemoveTx(id crypto.Hash)
	AllTransactions() []*tx.Tx
}
