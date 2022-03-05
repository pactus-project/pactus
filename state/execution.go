package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

func (st *state) executeBlock(block *block.Block, sb sandbox.Sandbox) ([]*tx.Tx, error) {
	exe := execution.NewExecution()

	ids := block.TxIDs().IDs()
	trxs := make([]*tx.Tx, len(ids))
	var mintbaseTrx *tx.Tx
	for i := 0; i < len(ids); i++ {
		trx := st.txPool.QueryTx(ids[i])
		if trx == nil {
			return nil, errors.Errorf(errors.ErrInvalidBlock,
				"transaction not found: %s", ids[i])
		}
		// The first transaction should be subsidy transaction
		IsMintbaseTx := (i == 0)
		if IsMintbaseTx {
			if !trx.IsMintbaseTx() {
				return nil, errors.Errorf(errors.ErrInvalidTx,
					"first transaction should be a subsidy transaction")
			}
			mintbaseTrx = trx
		} else {
			if trx.IsMintbaseTx() {
				return nil, errors.Errorf(errors.ErrInvalidTx,
					"duplicated subsidy transaction")
			}
		}

		err := exe.Execute(trx, sb)
		if err != nil {
			return nil, err
		}
		trxs[i] = trx
	}

	accumulatedFee := exe.AccumulatedFee()
	subsidyAmt := st.params.BlockReward + exe.AccumulatedFee()
	if mintbaseTrx.Payload().Value() != subsidyAmt {
		return nil, errors.Errorf(errors.ErrInvalidTx,
			"invalid subsidy amount. Expected %v, got %v", subsidyAmt, mintbaseTrx.Payload().Value())
	}

	// Claim accumulated fees
	acc := sb.Account(crypto.TreasuryAddress)
	acc.AddToBalance(accumulatedFee)
	sb.UpdateAccount(acc)

	return trxs, nil
}
