package txpool

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type Reader interface {
	PrepareBlockTransactions() block.Txs
	PendingTx(id tx.ID) *tx.Tx
	HasTx(id tx.ID) bool
	Size() int
	EstimatedFee(amt amount.Amount, payloadType payload.Type) amount.Amount
	AllPendingTxs() []*tx.Tx
}

type TxPool interface {
	Reader

	SetNewSandboxAndRecheck(sb sandbox.Sandbox)
	AppendTxAndBroadcast(trx *tx.Tx) error
	AppendTx(trx *tx.Tx) error
	HandleCommittedBlock(block *block.Block) error
}
