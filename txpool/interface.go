package txpool

import (
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

type Reader interface {
	PrepareBlockTransactions() []*tx.Tx
	PendingTx(id tx.ID) *tx.Tx
	QueryTx(id tx.ID) *tx.Tx
	HasTx(id tx.ID) bool
	Size() int
	Fingerprint() string
}

type TxPool interface {
	Reader

	SetNewSandboxAndRecheck(sb sandbox.Sandbox)
	AppendTx(tx *tx.Tx) error
	AppendTxAndBroadcast(trx *tx.Tx) error
	RemoveTx(id tx.ID)
}
