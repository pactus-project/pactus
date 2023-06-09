package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/errors"
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
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve validator")
	}
	if val.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidSequence,
			"expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
	}
	if val.UnbondingHeight() > 0 {
		return errors.Errorf(errors.ErrInvalidHeight,
			"validator has unbonded at height %v", val.UnbondingHeight())
	}
	if e.strict {
		// In strict mode, the unbond transaction will be rejected if the
		// validator is in the committee.
		// In non-strict mode, we accept it and keep it inside the transaction pool to
		// process it when the validator leaves the committee.
		if sb.Committee().Contains(pld.Validator) {
			return errors.Errorf(errors.ErrInvalidTx,
				"validator %v is in committee", pld.Validator)
		}

		// In strict mode, unbond transactions will be rejected if a validator is
		// going to be in the committee for the next height.
		// In non-strict mode, we accept it and keep it inside the transaction pool to
		// process it when the validator leaves the committee.
		if val.LastJoinedHeight() == sb.CurrentHeight() {
			return errors.Errorf(errors.ErrInvalidHeight,
				"validator %v joins committee in the next height", pld.Validator)
		}
	}

	val.IncSequence()
	val.UpdateUnbondingHeight(sb.CurrentHeight())
	sb.UpdateValidator(val)

	return nil
}

// Fee will return unbond execution fee.
func (e *UnbondExecutor) Fee() int64 {
	return 0
}
