package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

type TransferPayload struct {
	From   crypto.Address
	To     crypto.Address
	Amount amount.Amount
}

func (*TransferPayload) Type() Type {
	return TypeTransfer
}

func (p *TransferPayload) Signer() crypto.Address {
	return p.From
}

func (p *TransferPayload) Value() amount.Amount {
	return p.Amount
}

// BasicCheck performs basic checks on the Transfer payload.
func (p *TransferPayload) BasicCheck() error {
	if !p.From.IsAccountAddress() {
		return BasicCheckError{
			Reason: "sender is not an account address: " + p.From.String(),
		}
	}
	if !p.To.IsAccountAddress() {
		return BasicCheckError{
			Reason: "receiver is not an account address: " + p.To.String(),
		}
	}

	return nil
}

func (p *TransferPayload) SerializeSize() int {
	return p.From.SerializeSize() +
		p.To.SerializeSize() +
		encoding.VarIntSerializeSize(uint64(p.Amount))
}

func (p *TransferPayload) Encode(w io.Writer) error {
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

func (p *TransferPayload) Decode(r io.Reader) error {
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

// LogString returns a concise string representation intended for use in logs.
func (p *TransferPayload) LogString() string {
	return fmt.Sprintf("{Transfer 💸 %s->%s %s",
		p.From.LogString(),
		p.To.LogString(),
		p.Amount)
}
