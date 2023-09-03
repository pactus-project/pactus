package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

type WithdrawPayload struct {
	From crypto.Address // withdraw from validator address
	To   crypto.Address // deposit to account address
}

func (p *WithdrawPayload) Type() Type {
	return TypeWithdraw
}

func (p *WithdrawPayload) Signer() crypto.Address {
	return p.From
}

func (p *WithdrawPayload) Value() int64 {
	return 0
}

// TODO: write test for me.
func (p *WithdrawPayload) BasicCheck() error {
	if err := p.From.BasicCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}
	if err := p.To.BasicCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}

	return nil
}

func (p *WithdrawPayload) SerializeSize() int {
	return 42 + encoding.VarIntSerializeSize(uint64(0))
}

func (p *WithdrawPayload) Encode(w io.Writer) error {
	err := encoding.WriteElements(w, &p.From, &p.To)
	if err != nil {
		return err
	}
	return encoding.WriteVarInt(w, uint64(0))
}

func (p *WithdrawPayload) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &p.From, &p.To)
	if err != nil {
		return err
	}
	return nil
}

func (p *WithdrawPayload) String() string {
	return fmt.Sprintf("{WithdrawPayload 🧾 %v->%v",
		p.From.ShortString(),
		p.To.ShortString())
}
