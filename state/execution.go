package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

func (st *state) executeBlock(block block.Block) ([]tx.CommittedTx, error) {
	st.txPoolSandbox.Clear()
	st.executionSandbox.Clear()
	st.execution.ResetFee()

	ids := block.TxIDs().IDs()
	twrs := make([]tx.CommittedTx, len(ids))
	var subsidyTrx *tx.Tx
	for i := 0; i < len(ids); i++ {
		trx := st.txPool.PendingTx(ids[i])
		if trx == nil {
			return nil, errors.Errorf(errors.ErrInvalidBlock, "Transaction not found")
		}
		// Only first transaction should be subsidy transaction
		isSubsidyTx := (i == 0)
		if isSubsidyTx {
			if !trx.IsSubsidyTx() {
				return nil, errors.Errorf(errors.ErrInvalidTx, "First transaction should be a subsidy transaction")
			}
			subsidyTrx = trx
		} else {
			if trx.IsSubsidyTx() {
				return nil, errors.Errorf(errors.ErrInvalidTx, "Duplicated subsidy transaction")
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

	subsidyAmt := calcBlockSubsidy(st.lastBlockHeight+1, st.params.SubsidyReductionInterval) + st.execution.AccumulatedFee()
	if subsidyTrx.Payload().Value() != subsidyAmt {
		return nil, errors.Errorf(errors.ErrInvalidTx, "Invalid subsidy amount. Expected %v, got %v", subsidyAmt, subsidyTrx.Payload().Value())
	}

	return twrs, nil
}
