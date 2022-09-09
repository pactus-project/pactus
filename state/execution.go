package state

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/errors"
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
