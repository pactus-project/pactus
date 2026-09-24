package state

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

func (st *state) executeBlock(blk *block.Block, sbx sandbox.Sandbox, check bool) error {
	proposerAddr := blk.Header().ProposerAddress()
	for i, trx := range blk.Transactions() {
		if check {
			// The first transaction should be subsidy transaction
			shouldBeSubsidyTx := (i == 0)
			err := st.checkSubsidy(trx, proposerAddr, shouldBeSubsidyTx)
			if err != nil {
				return err
			}

			err = execution.CheckAndExecute(trx, sbx, true)
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

	subsidyTrx := blk.Transactions().Subsidy()
	accumulatedFee := sbx.AccumulatedFee()
	blockReward := st.params.BlockReward(blk.Height())
	subsidyAmt := blockReward + sbx.AccumulatedFee()
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

func (st *state) checkSubsidy(trx *tx.Tx, proposerAddr crypto.Address, shouldBeSubsidyTx bool) error {
	if !shouldBeSubsidyTx {
		if trx.IsSubsidyTx() {
			return ErrDuplicatedSubsidyTransaction
		}

		return nil
	}

	if !trx.IsSubsidyTx() {
		return ErrInvalidSubsidyTransaction
	}

	lockTime := trx.LockTime()
	batchTrx, ok := trx.Payload().(*payload.BatchTransferPayload)
	if !ok {
		return ErrInvalidSubsidyTransaction
	}

	// LockTime for the subsidy transaction is same as block height.
	foundationReward := st.params.FoundationReward(lockTime)
	if batchTrx.Recipients[0].Amount != foundationReward {
		return ErrInvalidSubsidyTransaction
	}

	foundationAddress := st.params.FoundationAddress(lockTime)
	if batchTrx.Recipients[0].To != foundationAddress {
		return ErrInvalidSubsidyTransaction
	}

	val, err := st.store.Validator(proposerAddr)
	if err != nil {
		return ErrInvalidSubsidyTransaction
	}

	if val.IsDelegated() {
		if val.DelegateShare() > 0 {
			rewardCoeff := st.params.RewardCoefficient(lockTime)
			dlgOwner := val.DelegateOwner()
			dlgShare := val.DelegateShare().MulF64(rewardCoeff)
			if batchTrx.Recipients[1].To != dlgOwner ||
				batchTrx.Recipients[1].Amount < dlgShare {
				return ErrInvalidSubsidyTransaction
			}
		}
	} else {
		// 2 recipients: foundation + validator
		if len(batchTrx.Recipients) != 2 {
			return ErrInvalidSubsidyTransaction
		}
	}

	return nil
}
