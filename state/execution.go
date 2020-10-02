package state

import (
	"gitlab.com/zarb-chain/zarb-go/block"
	"gitlab.com/zarb-chain/zarb-go/errors"
	"gitlab.com/zarb-chain/zarb-go/execution"
	"gitlab.com/zarb-chain/zarb-go/tx"
)

func (st *State) executeBlock(block *block.Block, exe *execution.Executor) ([]*tx.Receipt, error) {
	if block == nil {
		return nil, errors.Errorf(errors.ErrInvalidBlock, "Block is empty")
	}

	hashes := block.Txs().TxHashes()
	receipts := make([]*tx.Receipt, len(hashes))

	for i := 0; i < len(hashes); i++ {
		trx, found := st.txPool.PendingTx(hashes[i])
		if !found {
			return nil, errors.Errorf(errors.ErrInvalidBlock, "We don't have transaction to validate the block")
		}
		if err := trx.SanityCheck(); err != nil {
			return nil, errors.Errorf(errors.ErrInvalidBlock, "Invalid transaction %v", err)
		}
		// Only first transaction is mintbase transaction
		receipt, err := exe.Execute(trx, (i == 0))
		if err != nil {
			return nil, errors.Errorf(errors.ErrInvalidBlock, "Invalid transaction inside the block: %v", err)
		}

		receipt.SetBlockHash(block.Hash())
		receipts[i] = receipt
	}

	// Now, check rewards + fee
	//tx, _ := st.txPool.PendingTx(hashes[0])


	return receipts, nil
}
