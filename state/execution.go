package state

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
)

func (st *state) executeBlock(b *block.Block, sb sandbox.Sandbox, check bool) error {
	var subsidyTrx *tx.Tx
	for i, trx := range b.Transactions() {
		// The first transaction should be subsidy transaction
		isSubsidyTx := (i == 0)
		if isSubsidyTx {
			if !trx.IsSubsidyTx() {
				return ErrInvalidSubsidyTransaction
			}
			subsidyTrx = trx
		} else if trx.IsSubsidyTx() {
			return ErrDuplicatedSubsidyTransaction
		}

		if check {
			if err := st.checkEd25519Fork(trx); err != nil {
				return err
			}

			err := execution.CheckAndExecute(trx, sb, true)
			if err != nil {
				return err
			}
		} else {
			err := execution.Execute(trx, sb)
			if err != nil {
				return err
			}
		}
	}

	accumulatedFee := sb.AccumulatedFee()
	subsidyAmt := st.params.BlockReward + sb.AccumulatedFee()
	if subsidyTrx.Payload().Value() != subsidyAmt {
		return InvalidSubsidyAmountError{
			Expected: subsidyAmt,
			Got:      subsidyTrx.Payload().Value(),
		}
	}

	// Claim accumulated fees
	acc := sb.Account(crypto.TreasuryAddress)
	acc.AddToBalance(accumulatedFee)
	sb.UpdateAccount(crypto.TreasuryAddress, acc)

	return nil
}

func (st *state) checkEd25519Fork(trx *tx.Tx) error {
	// TODO: remove me after enabling Ed255519
	if trx.Payload().Signer().Type() == crypto.AddressTypeEd25519Account {
		if st.genDoc.ChainType().IsMainnet() {
			return fmt.Errorf("ed255519 not supported yet")
		}

		if st.genDoc.ChainType().IsTestnet() {
			if st.lastInfo.BlockHeight() < 1_320_000 {
				return fmt.Errorf("ed255519 not supported yet")
			}
		}
	}

	return nil
}
