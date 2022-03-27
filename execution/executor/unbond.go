package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type UnbondExecutor struct {
	strict bool
}

func NewUnbondExecutor(strict bool) *UnbondExecutor {
	return &UnbondExecutor{strict: strict}
}

func (e *UnbondExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.UnbondPayload)

	val := sb.Validator(pld.Signer())
	if val == nil {
		return errors.Errorf(errors.ErrInvalidAddress, "unable to retrieve validator")
	}
	if val.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidSequence, "expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
	}
	if val.UnbondingHeight() > 0 {
		return errors.Errorf(errors.ErrInvalidHeight, "validator has unbonded at height %v", val.UnbondingHeight())
	}
	if e.strict {
		// In strict mode, unbond transaction will be rejected if a validator is in committee.
		// In non-restrict mode, we accept it and keep it inside tx pool to process it later
		if sb.Committee().Contains(pld.Validator) {
			return errors.Errorf(errors.ErrInvalidTx, "validator %v is in committee", pld.Validator)
		}

		// In strict mode, a validator can not evaluate sortition after unbonding.
		// In non-restrict mode, we accept it and keep it inside tx pool to process it later
		if val.LastJoinedHeight() == sb.CurrentHeight() {
			return errors.Errorf(errors.ErrInvalidHeight, "validator %v will join committee", pld.Validator)
		}
	}

	val.IncSequence()
	val.UpdateUnbondingHeight(sb.CurrentHeight())
	sb.UpdateValidator(val)

	return nil
}

//Fee will return unbond execution fee
func (e *UnbondExecutor) Fee() int64 {
	return 0
}
