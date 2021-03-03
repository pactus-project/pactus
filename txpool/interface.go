package txpool

import (
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

type TxPoolReader interface {
	AllTransactions() []*tx.Tx
	PendingTx(id tx.ID) *tx.Tx
	BroadcastTxs(ids []tx.ID)
	HasTx(id tx.ID) bool
	Size() int

	Fingerprint() string
}

type TxPool interface {
	TxPoolReader

	SetSandbox(sandbox sandbox.Sandbox)
	AppendTx(tx *tx.Tx) error
	AppendTxAndBroadcast(trx *tx.Tx) error
	RemoveTx(id tx.ID)
	Recheck()
}
