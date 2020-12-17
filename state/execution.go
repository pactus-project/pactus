package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

func (st *state) executeBlock(block block.Block) ([]tx.CommittedTx, error) {
	ids := block.TxIDs().IDs()
	twrs := make([]tx.CommittedTx, len(ids))

	for i := 0; i < len(ids); i++ {
		trx := st.txPool.PendingTx(ids[i])
		if trx == nil {
			return nil, errors.Errorf(errors.ErrInvalidBlock, "Transaction not found")
		}
		if err := trx.SanityCheck(); err != nil {
			return nil, err
		}
		// Only first transaction should be mintbase transaction
		isMintbaseTx := (i == 0)
		if isMintbaseTx {
			if !trx.IsMintbaseTx() {
				return nil, errors.Errorf(errors.ErrInvalidTx, "Not a mintbase transaction")
			}
		} else {
			if trx.IsMintbaseTx() {
				return nil, errors.Errorf(errors.ErrInvalidTx, "Duplicated mintbase transaction")
			}
		}

		err := st.execution.Execute(trx)
		if err != nil {
			return nil, err
		}
		receipt := trx.GenerateReceipt(tx.Ok, block.Hash())
		twrs[i].Tx = trx
		twrs[i].Receipt = receipt
	}

	return twrs, nil
}
