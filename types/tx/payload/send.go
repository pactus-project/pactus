package payload

import (
	"fmt"
	"io"

	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/errors"
)

type SendPayload struct {
	Sender   crypto.Address
	Receiver crypto.Address
	Amount   int64
}

func (p *SendPayload) Type() Type {
	return PayloadTypeSend
}

func (p *SendPayload) Signer() crypto.Address {
	return p.Sender
}

func (p *SendPayload) Value() int64 {
	return p.Amount
}

func (p *SendPayload) SanityCheck() error {
	if err := p.Sender.SanityCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}
	if err := p.Receiver.SanityCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}

	return nil
}

func (p *SendPayload) SerializeSize() int {
	return 42 + encoding.VarIntSerializeSize(uint64(p.Amount))
}

func (p *SendPayload) Encode(w io.Writer) error {
	err := encoding.WriteElements(w, &p.Sender, &p.Receiver)
	if err != nil {
		return err
	}
	return encoding.WriteVarInt(w, uint64(p.Amount))
}

func (p *SendPayload) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &p.Sender, &p.Receiver)
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

func (p *SendPayload) Fingerprint() string {
	return fmt.Sprintf("{Send 💸 %v->%v %v",
		p.Sender.Fingerprint(),
		p.Receiver.Fingerprint(),
		p.Amount)
}
