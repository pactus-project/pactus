package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
)

type BondExecutor struct {
	sbx      sandbox.Sandbox
	pld      *payload.BondPayload
	fee      amount.Amount
	sender   *account.Account
	receiver *validator.Validator
}

func newBondExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*BondExecutor, error) {
	pld := trx.Payload().(*payload.BondPayload)

	sender := sbx.Account(pld.From)
	if sender == nil {
		return nil, AccountNotFoundError{Address: pld.From}
	}

	receiver := sbx.Validator(pld.To)
	if receiver == nil {
		if pld.PublicKey == nil {
			return nil, ErrPublicKeyNotSet
		}
		receiver = sbx.MakeNewValidator(pld.PublicKey)
	} else if pld.PublicKey != nil {
		return nil, ErrPublicKeyAlreadySet
	}

	return &BondExecutor{
		sbx:      sbx,
		pld:      pld,
		fee:      trx.Fee(),
		sender:   sender,
		receiver: receiver,
	}, nil
}

func (e *BondExecutor) Check(strict bool) error {
	if e.receiver.IsUnbonded() {
		return ErrValidatorUnbonded
	}

	if e.sender.Balance() < e.pld.Value()+e.fee {
		return ErrInsufficientFunds
	}

	if e.pld.Stake < e.sbx.Params().MinimumStake {
		// This check prevents a potential attack where an attacker could send zero
		// or a small amount of stake to a full validator, effectively parking the
		// validator for the bonding period.
		if e.pld.Stake == 0 || e.pld.Stake+e.receiver.Stake() != e.sbx.Params().MaximumStake {
			return SmallStakeError{
				Minimum: e.sbx.Params().MinimumStake,
			}
		}
	}

	if e.receiver.Stake()+e.pld.Stake > e.sbx.Params().MaximumStake {
		return MaximumStakeError{
			Maximum: e.sbx.Params().MaximumStake,
		}
	}

	if strict {
		// In strict mode, bond transactions will be rejected if a validator is
		// already in the committee.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if e.sbx.Committee().Contains(e.pld.To) {
			return ErrValidatorInCommittee
		}

		// In strict mode, bond transactions will be rejected if a validator is
		// going to join the committee in the next height.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if e.sbx.IsJoinedCommittee(e.pld.To) {
			return ErrValidatorInCommittee
		}
	}

	return nil
}

func (e *BondExecutor) Execute() {
	e.sender.SubtractFromBalance(e.pld.Stake + e.fee)
	e.receiver.AddToStake(e.pld.Stake)
	e.receiver.UpdateLastBondingHeight(e.sbx.CurrentHeight())

	e.sbx.UpdatePowerDelta(int64(e.pld.Stake))
	e.sbx.UpdateAccount(e.pld.From, e.sender)
	e.sbx.UpdateValidator(e.receiver)
}
