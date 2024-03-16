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
}

func (p *WithdrawPayload) Type() Type {
	return TypeWithdraw
}

func (p *WithdrawPayload) Signer() crypto.Address {
	return p.From
}

func (p *WithdrawPayload) Value() amount.Amount {
	return p.Amount
}

// TODO: write test for me.
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

	return nil
}

func (p *WithdrawPayload) SerializeSize() int {
	return 42 + encoding.VarIntSerializeSize(uint64(p.Amount))
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

	return encoding.WriteVarInt(w, uint64(p.Amount))
}

func (p *WithdrawPayload) Decode(r io.Reader) error {
	err := p.From.Decode(r)
	if err != nil {
		return err
	}

	err = p.To.Decode(r)
	if err != nil {
		return err
	}

	amt, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	p.Amount = amount.Amount(amt)

	return nil
}

func (p *WithdrawPayload) String() string {
	return fmt.Sprintf("{WithdrawPayload 🧾 %s->%s %s",
		p.From.ShortString(),
		p.To.ShortString(),
		p.Amount)
}

func (p *WithdrawPayload) Receiver() *crypto.Address {
	return &p.To
}
