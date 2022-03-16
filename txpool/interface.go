package txpool

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
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
