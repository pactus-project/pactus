package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
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
		return errors.Errorf(errors.ErrInvalidAddress, "Unable to retrieve validator account")
	}

	if val.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidSequence, "Invalid sequence, Expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
	}
	if val.Stake() < pld.Amount+trx.Fee() {
		return errors.Errorf(errors.ErrInsufficientFunds, "Insufficient balance")
	}
	if val.UnbondingHeight() == 0 {
		return errors.Errorf(errors.ErrInvalidTx, "Need to unbond first")
	}
	if sb.CurrentHeight() < val.UnbondingHeight()+sb.UnbondInterval() {
		return errors.Errorf(errors.ErrInvalidTx, "Hasn't passed unbonding period, Expected: %v, got: %v", val.UnbondingHeight()+sb.UnbondInterval(), sb.CurrentHeight())
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
