package txpool

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

type TxPoolReader interface {
	AllTransactions() []*tx.Tx
	PendingTx(id crypto.Hash) *tx.Tx
	HasTx(id crypto.Hash) bool
	Size() int

	Fingerprint() string
}

type TxPool interface {
	TxPoolReader

	SetSandbox(sandbox sandbox.Sandbox)
	AppendTx(tx *tx.Tx) error
	AppendTxAndBroadcast(trx *tx.Tx) error
	RemoveTx(id crypto.Hash)
	Recheck()
}
