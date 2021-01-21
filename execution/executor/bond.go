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
}

func NewBondExecutor(sb sandbox.Sandbox) *BondExecutor {
	return &BondExecutor{sandbox: sb}
}

func (e *BondExecutor) Execute(trx *tx.Tx) error {
	pld := trx.Payload().(*payload.BondPayload)

	bonderAcc := e.sandbox.Account(pld.Bonder)
	if bonderAcc == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve bonder account")
	}
	bondVal := e.sandbox.Validator(pld.Validator.Address())
	if bondVal == nil {
		bondVal = e.sandbox.MakeNewValidator(pld.Validator)
	}
	if bonderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", bonderAcc.Sequence()+1, trx.Sequence())
	}
	if bonderAcc.Balance() < pld.Stake+trx.Fee() {
		return errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}
	bonderAcc.IncSequence()
	bonderAcc.SubtractFromBalance(pld.Stake + trx.Fee())
	bondVal.AddToStake(pld.Stake)

	e.sandbox.UpdateAccount(bonderAcc)
	e.sandbox.UpdateValidator(bondVal)

	e.fee = trx.Fee()

	return nil
}

func (e *BondExecutor) Fee() int64 {
	return e.fee
}
