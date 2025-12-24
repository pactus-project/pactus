package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

const maxBatchRecipients = 8

type BatchRecipient struct {
	To     crypto.Address
	Amount amount.Amount
}

func (tr *BatchRecipient) BasicCheck() error {
	if !tr.To.IsAccountAddress() {
		return BasicCheckError{
			Reason: "receiver is not an account address: " + tr.To.String(),
		}
	}

	if tr.Amount <= 0 {
		return BasicCheckError{
			Reason: fmt.Sprintf("amount must be greater than zero: %d", tr.Amount),
		}
	}

	return nil
}

func (tr *BatchRecipient) SerializeSize() int {
	return tr.To.SerializeSize() +
		encoding.VarIntSerializeSize(uint64(tr.Amount))
}

func (tr *BatchRecipient) Encode(w io.Writer) error {
	err := tr.To.Encode(w)
	if err != nil {
		return err
	}

	return encoding.WriteVarInt(w, uint64(tr.Amount))
}

func (tr *BatchRecipient) Decode(r io.Reader) error {
	err := tr.To.Decode(r)
	if err != nil {
		return err
	}

	amt, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	tr.Amount = amount.Amount(amt)

	return nil
}

type BatchTransferPayload struct {
	From       crypto.Address
	Recipients []BatchRecipient
}

func (*BatchTransferPayload) Type() Type {
	return TypeBatchTransfer
}

func (p *BatchTransferPayload) Signer() crypto.Address {
	return p.From
}

func (p *BatchTransferPayload) Value() amount.Amount {
	value := amount.Amount(0)
	for _, r := range p.Recipients {
		value += r.Amount
	}

	return value
}

// BasicCheck performs basic checks on the Batch Transfer payload.
func (p *BatchTransferPayload) BasicCheck() error {
	if !p.From.IsAccountAddress() {
		return BasicCheckError{
			Reason: "sender is not an account address: " + p.From.String(),
		}
	}

	if len(p.Recipients) < 2 || len(p.Recipients) > maxBatchRecipients {
		return BasicCheckError{
			Reason: fmt.Sprintf("recipients must be between 2 and %d", maxBatchRecipients),
		}
	}

	for _, r := range p.Recipients {
		if err := r.BasicCheck(); err != nil {
			return err
		}
	}

	return nil
}

func (p *BatchTransferPayload) SerializeSize() int {
	size := p.From.SerializeSize()
	size += encoding.VarIntSerializeSize(uint64(len(p.Recipients)))
	for _, r := range p.Recipients {
		size += r.SerializeSize()
	}

	return size
}

func (p *BatchTransferPayload) Encode(w io.Writer) error {
	err := p.From.Encode(w)
	if err != nil {
		return err
	}

	err = encoding.WriteVarInt(w, uint64(len(p.Recipients)))
	if err != nil {
		return err
	}

	for _, r := range p.Recipients {
		err = r.Encode(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *BatchTransferPayload) Decode(r io.Reader) error {
	err := p.From.Decode(r)
	if err != nil {
		return err
	}

	numberOfRecipients, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}

	if numberOfRecipients > maxBatchRecipients {
		return ErrTooManyRecipients
	}

	p.Recipients = make([]BatchRecipient, numberOfRecipients)
	for i := uint64(0); i < numberOfRecipients; i++ {
		err := p.Recipients[i].Decode(r)
		if err != nil {
			return err
		}
	}

	return nil
}

// LogString returns a concise string representation intended for use in logs.
func (p *BatchTransferPayload) LogString() string {
	return fmt.Sprintf("{BatchTransfer 💸 %s->[%d] %s",
		p.From.LogString(),
		len(p.Recipients),
		p.Value())
}
