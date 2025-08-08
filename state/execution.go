package state

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
)

func (st *state) executeBlock(blk *block.Block, sbx sandbox.Sandbox, check bool) error {
	var subsidyTrx *tx.Tx
	for i, trx := range blk.Transactions() {
		// The first transaction should be subsidy transaction
		isSubsidyTx := (i == 0)
		if isSubsidyTx {
			err := st.checkSubsidy(blk, trx)
			if err != nil {
				return err
			}

			subsidyTrx = trx
		} else if trx.IsSubsidyTx() {
			return ErrDuplicatedSubsidyTransaction
		}

		if check {
			err := execution.CheckAndExecute(trx, sbx, true)
			if err != nil {
				return err
			}
		} else {
			err := execution.Execute(trx, sbx)
			if err != nil {
				return err
			}
		}
	}

	accumulatedFee := sbx.AccumulatedFee()
	subsidyAmt := st.params.BlockReward + sbx.AccumulatedFee()
	if subsidyTrx.Payload().Value() != subsidyAmt {
		return InvalidSubsidyAmountError{
			Expected: subsidyAmt,
			Got:      subsidyTrx.Payload().Value(),
		}
	}

	// Claim accumulated fees
	acc := sbx.Account(crypto.TreasuryAddress)
	acc.AddToBalance(accumulatedFee)
	sbx.UpdateAccount(crypto.TreasuryAddress, acc)

	return nil
}

func (st *state) checkSubsidy(blk *block.Block, trx *tx.Tx) error {
	if !trx.IsSubsidyTx() {
		return ErrInvalidSubsidyTransaction
	}

	if st.params.SplitRewardForkHeight > 0 && blk.Height() > st.params.SplitRewardForkHeight {
		if !trx.IsBatchTransferTx() {
			return ErrInvalidSubsidyTransaction
		}
	} else {
		if !trx.IsTransferTx() {
			return ErrInvalidSubsidyTransaction
		}
	}

	return nil
}
