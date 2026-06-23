package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
)

type UnbondExecutor struct {
	pld       *payload.UnbondPayload
	validator *validator.Validator
}

func newUnbondExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*UnbondExecutor, error) {
	pld := trx.Payload().(*payload.UnbondPayload)

	val := sbx.Validator(pld.Validator)
	if val == nil {
		return nil, ValidatorNotFoundError{Address: pld.Validator}
	}

	return &UnbondExecutor{
		pld:       pld,
		validator: val,
	}, nil
}

func (e *UnbondExecutor) Check(sbx sandbox.SandboxReader, strict bool) error {
	if e.validator.IsUnbonded() {
		return ErrValidatorUnbonded
	}

	if e.validator.IsDelegated() {
		// A delegated validator can only be unbonded by its real owner.
		if e.pld.DelegateOwner != e.validator.DelegateOwner() {
			return ErrInvalidDelegateOwner
		}

		if !e.validator.DelegateExpired(sbx.CurrentHeight()) {
			return ErrDelegationNotExpired
		}
	} else if e.pld.IsDelegated() {
		// A non-delegated validator can only be unbonded by its own key.
		// The payload must not be delegated, forcing Signer() == Validator.

		return ErrInvalidDelegateOwner
	}

	if strict {
		// In strict mode, the unbond transaction will be rejected if the
		// validator is in the committee.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if sbx.Committee().Contains(e.pld.Validator) {
			return ErrValidatorInCommittee
		}

		// In strict mode, unbond transactions will be rejected if a validator is
		// going to be in the committee for the next height.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if sbx.IsJoinedCommittee(e.pld.Validator) {
			return ErrValidatorInCommittee
		}
	}

	return nil
}

func (e *UnbondExecutor) Execute(sbx sandbox.Sandbox) {
	unbondedPower := e.validator.Power()
	e.validator.UpdateUnbondingHeight(sbx.CurrentHeight())

	// The validator's power is reduced to zero,
	// so we update the power delta with the negative value of the validator's power.
	sbx.UpdatePowerDelta(-1 * unbondedPower)
	sbx.UpdateValidator(e.validator)
}
