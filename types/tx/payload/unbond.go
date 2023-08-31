package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

type UnbondPayload struct {
	Validator crypto.Address
}

func (p *UnbondPayload) Type() Type {
	return TypeUnbound
}

func (p *UnbondPayload) Signer() crypto.Address {
	return p.Validator
}

func (p *UnbondPayload) Value() int64 {
	return 0
}

// TODO: write test for me.
func (p *UnbondPayload) BasicCheck() error {
	if err := p.Validator.BasicCheck(); err != nil {
		return errors.Error(errors.ErrInvalidAddress)
	}

	return nil
}

func (p *UnbondPayload) SerializeSize() int {
	return 21
}

func (p *UnbondPayload) Encode(w io.Writer) error {
	return encoding.WriteElements(w, &p.Validator)
}

func (p *UnbondPayload) Decode(r io.Reader) error {
	return encoding.ReadElements(r, &p.Validator)
}

func (p *UnbondPayload) String() string {
	return fmt.Sprintf("{Unbond 🔓 %v",
		p.Validator.ShortString(),
	)
}
