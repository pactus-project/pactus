package state

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

func (st *state) executeBlock(blk *block.Block, sbx sandbox.Sandbox, check bool) error {
	if blk.Transactions().Len() == 0 {
		// There is no subsidy transaction. The blocks of the network are checked
		// before, but a block certified by a committee is not trusted blindly.
		return ErrInvalidSubsidyTransaction
	}

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

	if check && sbx.BlockVersion().HasHardenedRules() {
		err := st.checkExactSubsidy(subsidyTrx, proposerAddr, accumulatedFee)
		if err != nil {
			return err
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

	if len(batchTrx.Recipients) < 2 {
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

// checkExactSubsidy checks that the subsidy transaction pays exactly what `createSubsidyTx` makes:
//   - the foundation gets its part,
//   - the owner of a delegated validator gets the share of the delegation,
//     and the accumulated fee if the share is the maximum share (PIP-49),
//   - the rest of the reward goes to the reward address of the proposer. It has no recipient if it is empty.
//
// The total value of the subsidy is checked before, but without this check the proposer
// could change the distribution, for example by keeping the fee that belongs to the stake owner.
// The accumulated fee is known after the execution of all the transactions of the block.
func (st *state) checkExactSubsidy(trx *tx.Tx, proposerAddr crypto.Address, accumulatedFee amount.Amount) error {
	batchTrx, ok := trx.Payload().(*payload.BatchTransferPayload)
	if !ok {
		return ErrInvalidSubsidyTransaction
	}

	val, err := st.store.Validator(proposerAddr)
	if err != nil {
		return ErrInvalidSubsidyTransaction
	}

	lockTime := trx.LockTime()
	blockReward := st.params.BlockReward(lockTime)
	foundationReward := st.params.FoundationReward(lockTime)

	expected := []payload.BatchRecipient{
		{To: st.params.FoundationAddress(lockTime), Amount: foundationReward},
	}
	proposerReward := blockReward + accumulatedFee - foundationReward

	if val.IsDelegated() {
		rewardCoeff := st.params.RewardCoefficient(lockTime)
		dlgOwnerShare := val.DelegateShare().MulF64(rewardCoeff)
		maxOwnerShare := param.MaxDelegateOwnerRewardShare.MulF64(rewardCoeff)

		if dlgOwnerShare > 0 {
			ownerAmount := dlgOwnerShare
			if dlgOwnerShare == maxOwnerShare {
				ownerAmount += accumulatedFee
			}

			expected = append(expected, payload.BatchRecipient{To: val.DelegateOwner(), Amount: ownerAmount})
			proposerReward -= ownerAmount
		}
	}

	if len(batchTrx.Recipients) != len(expected) && len(batchTrx.Recipients) != len(expected)+1 {
		return ErrInvalidSubsidyTransaction
	}

	for i, recipient := range expected {
		if batchTrx.Recipients[i] != recipient {
			return ErrInvalidSubsidyTransaction
		}
	}

	// The reward address of the proposer is chosen by the proposer, but the amount is not.
	hasProposerReward := len(batchTrx.Recipients) == len(expected)+1
	if hasProposerReward != (proposerReward > 0) {
		return ErrInvalidSubsidyTransaction
	}
	if hasProposerReward && batchTrx.Recipients[len(expected)].Amount != proposerReward {
		return ErrInvalidSubsidyTransaction
	}

	return nil
}
