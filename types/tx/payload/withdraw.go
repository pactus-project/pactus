package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

type WithdrawPayload struct {
	From   crypto.Address // withdraw from validator address
	To     crypto.Address // deposit to account address
	Amount amount.Amount  // amount to deposit

	// Delegation fields
	DelegateOwner crypto.Address
}

func (*WithdrawPayload) Type() Type {
	return TypeWithdraw
}

// Signer returns the address that must sign this transaction.
// For a delegated validator the stake owner is the signer, so the operator
// cannot withdraw the owner's stake (for example, by taking it as a fee).
func (p *WithdrawPayload) Signer() crypto.Address {
	if p.IsDelegated() {
		return p.DelegateOwner
	}

	return p.From
}

func (p *WithdrawPayload) Value() amount.Amount {
	return p.Amount
}

// BasicCheck performs basic checks on the Withdraw payload.
func (p *WithdrawPayload) BasicCheck() error {
	if !p.From.IsValidatorAddress() {
		return BasicCheckError{
			Reason: "sender is not a validator address: " + p.From.String(),
		}
	}
	if !p.To.IsAccountAddress() {
		return BasicCheckError{
			Reason: "receiver is not an account address: " + p.To.String(),
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

func (p *WithdrawPayload) SerializeSize() int {
	size := p.From.SerializeSize() +
		p.To.SerializeSize() +
		encoding.VarIntSerializeSize(uint64(p.Amount))
	if p.IsDelegated() {
		size += p.DelegateOwner.SerializeSize()
	}

	return size
}

func (p *WithdrawPayload) Encode(w io.Writer) error {
	err := p.From.Encode(w)
	if err != nil {
		return err
	}

	err = p.To.Encode(w)
	if err != nil {
		return err
	}

	err = p.Amount.Encode(w)
	if err != nil {
		return err
	}

	if p.IsDelegated() {
		if err := p.DelegateOwner.Encode(w); err != nil {
			return err
		}
	}

	return nil
}

func (p *WithdrawPayload) Decode(ctx DecodeContext, r io.Reader) error {
	err := p.From.Decode(r)
	if err != nil {
		return err
	}

	err = p.To.Decode(r)
	if err != nil {
		return err
	}

	err = p.Amount.Decode(r)
	if err != nil {
		return err
	}

	if ctx.WithDelegation {
		if err := p.DelegateOwner.Decode(r); err != nil {
			return err
		}
	}

	return nil
}

func (p *WithdrawPayload) IsDelegated() bool {
	return p.DelegateOwner != crypto.TreasuryAddress
}

// LogString returns a concise string representation intended for use in logs.
func (p *WithdrawPayload) LogString() string {
	return fmt.Sprintf("{Withdraw 🧾 %s->%s %s",
		p.From.LogString(),
		p.To.LogString(),
		p.Amount)
}
