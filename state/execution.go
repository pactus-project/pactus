package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/tx"
)

func (st *state) executeBlock(block block.Block, exe *execution.Execution) ([]*tx.Receipt, error) {
	hashes := block.TxHashes().Hashes()
	receipts := make([]*tx.Receipt, len(hashes))

	for i := 0; i < len(hashes); i++ {
		trx := st.txPool.PendingTx(hashes[i])
		if trx == nil {
			return nil, errors.Errorf(errors.ErrInvalidBlock, "Transaction not found")
		}
		if err := trx.SanityCheck(); err != nil {
			return nil, err
		}
		// Only first transaction is mintbase transaction
		err := exe.Execute(trx, (i == 0))
		if err != nil {
			return nil, err
		}
		receipt := trx.GenerateReceipt(tx.Ok, block.Hash())
		receipts[i] = receipt
	}

	// Now, check rewards + fee
	//tx, _ := st.txPool.PendingTx(hashes[0])

	return receipts, nil
}
