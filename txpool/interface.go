package txpool

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
)

type Reader interface {
	PrepareBlockTransactions() block.Txs
	PendingTx(id tx.ID) *tx.Tx
	HasTx(id tx.ID) bool
	Size() int
	Fingerprint() string
}

type TxPool interface {
	Reader

	SetNewSandboxAndRecheck(sb sandbox.Sandbox)
	AppendTxAndBroadcast(trx *tx.Tx) error
	AppendTx(tx *tx.Tx) error
	RemoveTx(id tx.ID)
}
