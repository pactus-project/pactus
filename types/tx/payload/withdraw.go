package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

type WithdrawPayload struct {
	From   crypto.Address // withdraw from validator address
	To     crypto.Address // deposit to account address
	Amount int64          // amount to deposit
}

func (p *WithdrawPayload) Type() Type {
	return PayloadTypeWithdraw
}

func (p *WithdrawPayload) Signer() crypto.Address {
	return p.From
}

func (p *WithdrawPayload) Value() int64 {
	return p.Amount
}

// TODO: write test for me
func (p *WithdrawPayload) SanityCheck() error {
	if err := p.From.SanityCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}
	if err := p.To.SanityCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}

	return nil
}

func (p *WithdrawPayload) SerializeSize() int {
	return 42 + encoding.VarIntSerializeSize(uint64(p.Amount))
}

func (p *WithdrawPayload) Encode(w io.Writer) error {
	err := encoding.WriteElements(w, &p.From, &p.To)
	if err != nil {
		return err
	}
	return encoding.WriteVarInt(w, uint64(p.Amount))
}

func (p *WithdrawPayload) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &p.From, &p.To)
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

func (p *WithdrawPayload) String() string {
	return fmt.Sprintf("{WithdrawPayload 🧾 %v->%v %v",
		p.From.ShortString(),
		p.To.ShortString(),
		p.Amount)
}
