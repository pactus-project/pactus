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

	if batchTrx.Recipients[0].Amount != st.params.FoundationReward {
		return ErrInvalidSubsidyTransaction
	}

	addressIndex := int(lockTime) % len(st.params.FoundationAddress)
	foundationAddress := st.params.FoundationAddress[addressIndex]
	if batchTrx.Recipients[0].To != foundationAddress {
		return ErrInvalidSubsidyTransaction
	}

	val, err := st.store.Validator(proposerAddr)
	if err != nil {
		return ErrInvalidSubsidyTransaction
	}

	if val.IsDelegated() {
		if val.DelegateShare() > 0 {
			if batchTrx.Recipients[1].To != val.DelegateOwner() || batchTrx.Recipients[1].Amount < val.DelegateShare() {
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
