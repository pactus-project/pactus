package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
)

type WithdrawExecutor struct {
	sb       sandbox.Sandbox
	pld      *payload.WithdrawPayload
	fee      amount.Amount
	sender   *validator.Validator
	receiver *account.Account
}

func newWithdrawExecutor(trx *tx.Tx, sb sandbox.Sandbox) (*WithdrawExecutor, error) {
	pld := trx.Payload().(*payload.WithdrawPayload)

	sender := sb.Validator(pld.From)
	if sender == nil {
		return nil, ValidatorNotFoundError{Address: pld.From}
	}

	receiver := sb.Account(pld.To)
	if receiver == nil {
		receiver = sb.MakeNewAccount(pld.To)
	}

	return &WithdrawExecutor{
		sb:       sb,
		pld:      pld,
		fee:      trx.Fee(),
		sender:   sender,
		receiver: receiver,
	}, nil
}

func (e *WithdrawExecutor) Check(strict bool) error {
	if e.sender.Stake() < e.pld.Amount+e.fee {
		return ErrInsufficientFunds
	}

	if e.sender.UnbondingHeight() == 0 {
		return ErrValidatorBonded
	}

	if e.sb.CurrentHeight() < e.sender.UnbondingHeight()+e.sb.Params().UnbondInterval {
		return ErrUnbondingPeriod
	}

	return nil
}

func (e *WithdrawExecutor) Execute() {
	e.sender.SubtractFromStake(e.pld.Amount + e.fee)
	e.receiver.AddToBalance(e.pld.Amount)

	e.sb.UpdateValidator(e.sender)
	e.sb.UpdateAccount(e.pld.To, e.receiver)
}
