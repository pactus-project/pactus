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
	PendingTx(txID tx.ID) *tx.Tx
	HasTx(txID tx.ID) bool
	Size() int
	EstimatedFee(amt amount.Amount, payloadType payload.Type) amount.Amount
	AllPendingTxs() []*tx.Tx
}

type TxPool interface {
	Reader

	SetNewSandboxAndRecheck(sbx sandbox.Sandbox)
	AppendTxAndBroadcast(trx *tx.Tx) error
	AppendTx(trx *tx.Tx) error
	HandleCommittedBlock(blk *block.Block)
}
