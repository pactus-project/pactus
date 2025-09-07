package state

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

func (st *state) executeBlock(blk *block.Block, sbx sandbox.Sandbox, check bool) error {
	var subsidyTrx *tx.Tx
	for i, trx := range blk.Transactions() {
		// TODO: This check can be omitted during the old block check.
		// The first transaction should be subsidy transaction
		isSubsidyTx := (i == 0)
		if isSubsidyTx {
			err := st.checkSubsidy(blk.Header().Version(), trx)
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

//nolint:all // Remove me after enabling split reward forks
func (st *state) checkSubsidy(blockVersion protocol.Version, trx *tx.Tx) error {
	if !trx.IsSubsidyTx() {
		return ErrInvalidSubsidyTransaction
	}

	lockTime := trx.LockTime()
	if blockVersion >= protocol.ProtocolVersion2 {
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

		return nil
	} else {
		if !trx.IsTransferTx() {
			return ErrInvalidSubsidyTransaction
		}
	}

	return nil
}
