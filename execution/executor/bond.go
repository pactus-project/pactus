package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
)

type BondExecutor struct {
	pld          *payload.BondPayload
	fee          amount.Amount
	sender       *account.Account
	receiver     *validator.Validator
	newValidator bool
}

func newBondExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*BondExecutor, error) {
	pld := trx.Payload().(*payload.BondPayload)

	sender := sbx.Account(pld.From)
	if sender == nil {
		return nil, AccountNotFoundError{Address: pld.From}
	}

	receiver := sbx.Validator(pld.To)
	newValidator := false
	if receiver == nil {
		if pld.PublicKey == nil {
			return nil, ErrPublicKeyNotSet
		}
		// sbx.MakeNewValidator is deferred to Execute, which only runs once
		// Check has accepted the transaction. Registering the validator here
		// would leak it into the sandbox even if Check later rejects the
		// transaction (e.g. insufficient stake), letting a later transaction
		// bond to this half-created validator without ever setting its
		// public key. This placeholder is for validation only and is never
		// written back.
		receiver = validator.NewValidator(pld.PublicKey, 0)
		newValidator = true
	} else if pld.PublicKey != nil {
		return nil, ErrPublicKeyAlreadySet
	}

	return &BondExecutor{
		pld:          pld,
		fee:          trx.Fee(),
		sender:       sender,
		receiver:     receiver,
		newValidator: newValidator,
	}, nil
}

func (e *BondExecutor) Check(sbx sandbox.SandboxReader, strict bool) error {
	if e.receiver.IsUnbonded() {
		return ErrValidatorUnbonded
	}

	if e.pld.IsDelegated() {
		if e.pld.Stake != sbx.Params().MaximumStake {
			return ErrInvalidDelegation
		}
		if e.pld.DelegateExpiry <= sbx.CurrentHeight() {
			return ErrDelegateExpiryInPast
		}
	}

	if e.sender.Balance() < e.pld.Value()+e.fee {
		return ErrInsufficientFunds
	}

	if e.pld.Stake < sbx.Params().MinimumStake {
		// This check prevents a potential attack where an attacker could send zero
		// or a small amount of stake to a full validator, effectively parking the
		// validator for the bonding period.
		if e.pld.Stake == 0 || e.pld.Stake+e.receiver.Stake() != sbx.Params().MaximumStake {
			return SmallStakeError{
				Minimum: sbx.Params().MinimumStake,
			}
		}
	}

	if e.receiver.Stake()+e.pld.Stake > sbx.Params().MaximumStake {
		return MaximumStakeError{
			Maximum: sbx.Params().MaximumStake,
		}
	}

	if strict {
		// In strict mode, bond transactions will be rejected if a validator is
		// already in the committee.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if sbx.Committee().Contains(e.pld.To) {
			return ErrValidatorInCommittee
		}

		// In strict mode, bond transactions will be rejected if a validator is
		// going to join the committee in the next height.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if sbx.IsJoinedCommittee(e.pld.To) {
			return ErrValidatorInCommittee
		}
	}

	return nil
}

func (e *BondExecutor) Execute(sbx sandbox.Sandbox) {
	if e.newValidator {
		e.receiver = sbx.MakeNewValidator(e.pld.PublicKey)
	}

	e.sender.SubtractFromBalance(e.pld.Stake + e.fee)
	e.receiver.AddToStake(e.pld.Stake)
	e.receiver.UpdateLastBondingHeight(sbx.CurrentHeight())

	if e.pld.IsDelegated() {
		e.receiver.SetDelegation(e.pld.DelegateOwner, e.pld.DelegateShare, e.pld.DelegateExpiry)
	}

	sbx.UpdatePowerDelta(int64(e.pld.Stake))
	sbx.UpdateAccount(e.pld.From, e.sender)
	sbx.UpdateValidator(e.receiver)
}
