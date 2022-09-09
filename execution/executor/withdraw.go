package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/errors"
)

type WithdrawExecutor struct {
	fee    int64
	strict bool
}

func NewWithdrawExecutor(strict bool) *WithdrawExecutor {
	return &WithdrawExecutor{strict: strict}
}

func (e *WithdrawExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.WithdrawPayload)

	val := sb.Validator(pld.From)
	if val == nil {
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve validator account")
	}

	if val.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidSequence,
			"expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
	}
	if val.Stake() < pld.Amount+trx.Fee() {
		return errors.Error(errors.ErrInsufficientFunds)
	}
	if val.UnbondingHeight() == 0 {
		return errors.Errorf(errors.ErrInvalidHeight,
			"need to unbond first")
	}
	if sb.CurrentHeight() < val.UnbondingHeight()+sb.Params().UnbondInterval {
		return errors.Errorf(errors.ErrInvalidHeight,
			"hasn't passed unbonding period, expected: %v, got: %v",
			val.UnbondingHeight()+sb.Params().UnbondInterval, sb.CurrentHeight())
	}

	acc := sb.Account(pld.To)
	if acc == nil {
		acc = sb.MakeNewAccount(pld.To)
	}

	val.IncSequence()
	val.SubtractFromStake(pld.Amount + trx.Fee())
	acc.AddToBalance(pld.Amount)

	sb.UpdateValidator(val)
	sb.UpdateAccount(acc)

	e.fee = trx.Fee()

	return nil
}

func (e *WithdrawExecutor) Fee() int64 {
	return e.fee
}
