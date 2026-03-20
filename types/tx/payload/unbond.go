package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
)

type UnbondPayload struct {
	Validator crypto.Address

	// Delegation fields
	DelegateOwner crypto.Address
}

func (*UnbondPayload) Type() Type {
	return TypeUnbond
}

func (p *UnbondPayload) Signer() crypto.Address {
	if p.IsDelegated() {
		return p.DelegateOwner
	}

	return p.Validator
}

func (*UnbondPayload) Value() amount.Amount {
	return 0
}

// BasicCheck performs basic checks on the Unbond payload.
func (p *UnbondPayload) BasicCheck() error {
	if !p.Validator.IsValidatorAddress() {
		return BasicCheckError{
			Reason: "address is not a validator address: " + p.Validator.String(),
		}
	}

	if p.IsDelegated() {
		if !p.DelegateOwner.IsAccountAddress() {
			return BasicCheckError{
				Reason: "delegation owner is not an account address: " + p.DelegateOwner.String(),
			}
		}
	}

	return nil
}

func (p *UnbondPayload) SerializeSize() int {
	size := 21 // Validator address size
	if p.IsDelegated() {
		size += 21 // Delegate owner address size
	}

	return size
}

func (p *UnbondPayload) Encode(w io.Writer) error {
	if err := p.Validator.Encode(w); err != nil {
		return err
	}
	if p.IsDelegated() {
		if err := p.DelegateOwner.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (p *UnbondPayload) Decode(ctx DecodeContext, r io.Reader) error {
	if err := p.Validator.Decode(r); err != nil {
		return err
	}
	if ctx.WithDelegation {
		if err := p.DelegateOwner.Decode(r); err != nil {
			return err
		}
	}

	return nil
}

func (p *UnbondPayload) IsDelegated() bool {
	return p.DelegateOwner != crypto.TreasuryAddress
}

// LogString returns a concise string representation intended for use in logs.
func (p *UnbondPayload) LogString() string {
	return fmt.Sprintf("{Unbond 🔓 %s",
		p.Validator.LogString(),
	)
}
