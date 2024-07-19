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
	sb       sandbox.Sandbox
	pld      *payload.BondPayload
	fee      amount.Amount
	sender   *account.Account
	receiver *validator.Validator
}

func newBondExecutor(trx *tx.Tx, sb sandbox.Sandbox) (*BondExecutor, error) {
	pld := trx.Payload().(*payload.BondPayload)

	sender := sb.Account(pld.From)
	if sender == nil {
		return nil, AccountNotFoundError{Address: pld.From}
	}

	receiver := sb.Validator(pld.To)
	if receiver == nil {
		if pld.PublicKey == nil {
			return nil, ErrPublicKeyNotSet
		}
		receiver = sb.MakeNewValidator(pld.PublicKey)
	} else if pld.PublicKey != nil {
		return nil, ErrPublicKeyAlreadySet
	}

	return &BondExecutor{
		sb:       sb,
		pld:      pld,
		fee:      trx.Fee(),
		sender:   sender,
		receiver: receiver,
	}, nil
}

func (e *BondExecutor) Check(strict bool) error {
	if e.receiver.UnbondingHeight() > 0 {
		return ErrValidatorUnbonded
	}

	if e.sender.Balance() < e.pld.Stake+e.fee {
		return ErrInsufficientFunds
	}

	if e.pld.Stake < e.sb.Params().MinimumStake {
		if e.pld.Stake == 0 || e.pld.Stake+e.receiver.Stake() != e.sb.Params().MaximumStake {
			return SmallStakeError{
				Minimum: e.sb.Params().MinimumStake,
			}
		}
	}

	if e.receiver.Stake()+e.pld.Stake > e.sb.Params().MaximumStake {
		return MaximumStakeError{
			Maximum: e.sb.Params().MaximumStake,
		}
	}

	if strict {
		// In strict mode, bond transactions will be rejected if a validator is
		// already in the committee.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if e.sb.Committee().Contains(e.pld.To) {
			return ErrValidatorInCommittee
		}

		// In strict mode, bond transactions will be rejected if a validator is
		// going to join the committee in the next height.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if e.sb.IsJoinedCommittee(e.pld.To) {
			return ErrValidatorInCommittee
		}
	}

	return nil
}

func (e *BondExecutor) Execute() {
	e.sender.SubtractFromBalance(e.pld.Stake + e.fee)
	e.receiver.AddToStake(e.pld.Stake)
	e.receiver.UpdateLastBondingHeight(e.sb.CurrentHeight())

	e.sb.UpdatePowerDelta(int64(e.pld.Stake))
	e.sb.UpdateAccount(e.pld.From, e.sender)
	e.sb.UpdateValidator(e.receiver)
}
