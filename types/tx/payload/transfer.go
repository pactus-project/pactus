package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/encoding"
)

type TransferPayload struct {
	Sender   crypto.Address
	Receiver crypto.Address
	Amount   int64
}

func (p *TransferPayload) Type() Type {
	return PayloadTypeTransfer
}

func (p *TransferPayload) Signer() crypto.Address {
	return p.Sender
}

func (p *TransferPayload) Value() int64 {
	return p.Amount
}

func (p *TransferPayload) BasicCheck() error {
	if err := p.Sender.BasicCheck(); err != nil {
		return err
	}
	return p.Receiver.BasicCheck()
}

func (p *TransferPayload) SerializeSize() int {
	if p.Sender.EqualsTo(crypto.TreasuryAddress) {
		return 22 + encoding.VarIntSerializeSize(uint64(p.Amount))
	}
	return 42 + encoding.VarIntSerializeSize(uint64(p.Amount))
}

func (p *TransferPayload) Encode(w io.Writer) error {
	// If the transaction is a subsidy transaction (sender is treasury address)
	// compress the address to one byte.
	// This helps to reduce the size of each block by 20 bytes.
	if p.Sender.EqualsTo(crypto.TreasuryAddress) {
		err := encoding.WriteElement(w, uint8(0))
		if err != nil {
			return err
		}
	} else {
		err := encoding.WriteElement(w, &p.Sender)
		if err != nil {
			return err
		}
	}

	err := encoding.WriteElement(w, &p.Receiver)
	if err != nil {
		return err
	}
	return encoding.WriteVarInt(w, uint64(p.Amount))
}

func (p *TransferPayload) Decode(r io.Reader) error {
	var sigType uint8
	err := encoding.ReadElement(r, &sigType)
	if err != nil {
		return err
	}

	if sigType == crypto.SignatureTypeTreasury {
		p.Sender = crypto.TreasuryAddress
	} else {
		p.Sender[0] = sigType

		err := encoding.ReadElement(r, p.Sender[1:])
		if err != nil {
			return err
		}
	}
	err = encoding.ReadElement(r, &p.Receiver)
	if err != nil {
		return err
	}
	amount, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	p.Amount = int64(amount)
	return nil
}

func (p *TransferPayload) String() string {
	return fmt.Sprintf("{Send 💸 %v->%v %v",
		p.Sender.ShortString(),
		p.Receiver.ShortString(),
		p.Amount)
}
