package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
)

type UnbondExecutor struct {
	sb        sandbox.Sandbox
	pld       *payload.UnbondPayload
	validator *validator.Validator
}

func NewUnbondExecutor(trx *tx.Tx, sb sandbox.Sandbox) (*UnbondExecutor, error) {
	pld := trx.Payload().(*payload.UnbondPayload)

	val := sb.Validator(pld.Signer())
	if val == nil {
		return nil, ValidatorNotFoundError{Address: pld.Validator}
	}

	return &UnbondExecutor{
		sb:        sb,
		pld:       pld,
		validator: val,
	}, nil
}

func (e *UnbondExecutor) Check(strict bool) error {
	if e.validator.UnbondingHeight() > 0 {
		return ErrValidatorUnbonded
	}

	if strict {
		// In strict mode, the unbond transaction will be rejected if the
		// validator is in the committee.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if e.sb.Committee().Contains(e.pld.Validator) {
			return ErrValidatorInCommittee
		}

		// In strict mode, unbond transactions will be rejected if a validator is
		// going to be in the committee for the next height.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if e.sb.IsJoinedCommittee(e.pld.Validator) {
			return ErrValidatorInCommittee
		}
	}

	return nil
}

func (e *UnbondExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) {
	unbondedPower := e.validator.Power()
	e.validator.UpdateUnbondingHeight(sb.CurrentHeight())

	// The validator's power is reduced to zero,
	// so we update the power delta with the negative value of the validator's power.
	sb.UpdatePowerDelta(-1 * unbondedPower)
	sb.UpdateValidator(e.validator)
}
