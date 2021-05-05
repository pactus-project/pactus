package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type UnbondExecutor struct {
	fee    int64
	strict bool
}

func NewUnbondExecutor(strict bool) *UnbondExecutor {
	return &UnbondExecutor{strict: strict}
}

func (e *UnbondExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.UnbondPayload)

	val := sb.Validator(pld.Validator)
	if val == nil {
		//if couldn't retrive the validator then cann't unbound it
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve validator assoiciated with this key")
	}
	if e.strict && sb.IsInCommittee(pld.Validator) {
		return errors.Errorf(errors.ErrInvalidTx, "Validator is in committee right now please wait")
	}
	if val.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
	}
	//should use the vlidator stake to pay fee??
	//can make the unbonding free

	val.IncSequence()
	// unbondingVal.SubtractFromBalance(trx.Fee())
	val.UpdateUnbondingHeight(sb.CurrentHeight())
	sb.UpdateValidator(val)

	e.fee = trx.Fee()

	return nil
}

//Fee will return unbound execution fee
func (e *UnbondExecutor) Fee() int64 {
	return 0
}
