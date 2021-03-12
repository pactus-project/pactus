package txpool

import (
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/tx"
)

type TxPoolReader interface {
	AllTransactions() []*tx.Tx
	PendingTx(id tx.ID) *tx.Tx
	HasTx(id tx.ID) bool
	Size() int
	Fingerprint() string
}

type TxPool interface {
	TxPoolReader

	SetChecker(checker *execution.Execution)
	AppendTx(tx *tx.Tx) error
	AppendTxAndBroadcast(trx *tx.Tx) error
	QueryTx(id tx.ID) *tx.Tx
	RemoveTx(id tx.ID)
	Recheck()
}
