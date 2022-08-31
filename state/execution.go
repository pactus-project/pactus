package state

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/types/block"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/util/errors"
)

func (st *state) executeBlock(b *block.Block, sb sandbox.Sandbox) error {
	exe := execution.NewExecutor()

	var subsidyTrx *tx.Tx
	for i, trx := range b.Transactions() {
		// The first transaction should be subsidy transaction
		isSubsidyTx := (i == 0)
		if isSubsidyTx {
			if !trx.IsSubsidyTx() {
				return errors.Errorf(errors.ErrInvalidTx,
					"first transaction should be a subsidy transaction")
			}
			subsidyTrx = trx
		} else {
			if trx.IsSubsidyTx() {
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
	if subsidyTrx.Payload().Value() != subsidyAmt {
		return errors.Errorf(errors.ErrInvalidTx,
			"invalid subsidy amount, expected %v, got %v", subsidyAmt, subsidyTrx.Payload().Value())
	}

	// Claim accumulated fees
	acc := sb.Account(crypto.TreasuryAddress)
	acc.AddToBalance(accumulatedFee)
	sb.UpdateAccount(acc)

	return nil
}
