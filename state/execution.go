package state

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

func (st *state) executeBlock(b *block.Block, sb sandbox.Sandbox) error {
	exe := execution.NewExecutor()

	var mintbaseTrx *tx.Tx
	for i, trx := range b.Transactions() {
		// The first transaction should be subsidy transaction
		IsMintbaseTx := (i == 0)
		if IsMintbaseTx {
			if !trx.IsMintbaseTx() {
				return errors.Errorf(errors.ErrInvalidTx,
					"first transaction should be a subsidy transaction")
			}
			mintbaseTrx = trx
		} else {
			if trx.IsMintbaseTx() {
				return errors.Errorf(errors.ErrInvalidTx,
					"duplicated subsidy transaction")
			}
		}

		err := exe.Execute(trx, sb)
		if err != nil {
			return err
		}
	}

	accumulatedFee := exe.AccumulatedFee()
	subsidyAmt := st.params.BlockReward + exe.AccumulatedFee()
	if mintbaseTrx.Payload().Value() != subsidyAmt {
		return errors.Errorf(errors.ErrInvalidTx,
			"invalid subsidy amount, expected %v, got %v", subsidyAmt, mintbaseTrx.Payload().Value())
	}

	// Claim accumulated fees
	acc := sb.Account(crypto.TreasuryAddress)
	acc.AddToBalance(accumulatedFee)
	sb.UpdateAccount(acc)

	return nil
}
