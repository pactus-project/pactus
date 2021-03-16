package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

func (st *state) executeBlock(block block.Block) ([]tx.CommittedTx, error) {
	st.execution.Reset()

	ids := block.TxIDs().IDs()
	twrs := make([]tx.CommittedTx, len(ids))
	var mintbaseTrx *tx.Tx
	for i := 0; i < len(ids); i++ {
		trx := st.txPool.QueryTx(ids[i])
		if trx == nil {
			return nil, errors.Errorf(errors.ErrInvalidBlock, "Transaction not found: %s", ids[i])
		}
		// Only first transaction should be subsidy transaction
		IsMintbaseTx := (i == 0)
		if IsMintbaseTx {
			if !trx.IsMintbaseTx() {
				return nil, errors.Errorf(errors.ErrInvalidTx, "First transaction should be a subsidy transaction")
			}
			mintbaseTrx = trx
		} else {
			if trx.IsMintbaseTx() {
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
	if mintbaseTrx.Payload().Value() != subsidyAmt {
		return nil, errors.Errorf(errors.ErrInvalidTx, "Invalid subsidy amount. Expected %v, got %v", subsidyAmt, mintbaseTrx.Payload().Value())
	}
	st.execution.ClaimAccumulatedFee()

	return twrs, nil
}
