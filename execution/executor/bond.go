package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type BondExecutor struct {
	sandbox sandbox.Sandbox
	fee     int64
	strict  bool
}

func NewBondExecutor(sb sandbox.Sandbox, strict bool) *BondExecutor {
	return &BondExecutor{sandbox: sb, strict: strict}
}

func (e *BondExecutor) Execute(trx *tx.Tx) error {
	pld := trx.Payload().(*payload.BondPayload)

	bonderAcc := e.sandbox.Account(pld.Bonder)
	if bonderAcc == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve bonder account")
	}
	val := e.sandbox.Validator(pld.Validator.Address())
	if val == nil {
		val = e.sandbox.MakeNewValidator(pld.Validator)
	}
	if e.strict && e.sandbox.IsInCommittee(pld.Validator.Address()) {
		return errors.Errorf(errors.ErrInvalidTx, "Validator is in committee right now")
	}
	if bonderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", bonderAcc.Sequence()+1, trx.Sequence())
	}
	if bonderAcc.Balance() < pld.Stake+trx.Fee() {
		return errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}
	bonderAcc.IncSequence()
	bonderAcc.SubtractFromBalance(pld.Stake + trx.Fee())
	val.AddToStake(pld.Stake)

	e.sandbox.UpdateAccount(bonderAcc)
	e.sandbox.UpdateValidator(val)

	e.fee = trx.Fee()

	return nil
}

func (e *BondExecutor) Fee() int64 {
	return e.fee
}
