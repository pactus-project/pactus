package payload

import (
	"fmt"
	"io"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/encoding"
	"github.com/zarbchain/zarb-go/errors"
)

type UnbondPayload struct {
	Validator crypto.Address
}

func (p *UnbondPayload) Type() Type {
	return PayloadTypeUnbond
}

func (p *UnbondPayload) Signer() crypto.Address {
	return p.Validator
}

func (p *UnbondPayload) Value() int64 {
	return 0
}

func (p *UnbondPayload) SanityCheck() error {
	if err := p.Validator.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidAddress, "invalid validator address")
	}

	return nil
}

func (p *UnbondPayload) Encode(w io.Writer) error {
	return encoding.WriteElements(w, &p.Validator)
}

func (p *UnbondPayload) Decode(r io.Reader) error {
	return encoding.ReadElements(r, &p.Validator)
}

func (p *UnbondPayload) Fingerprint() string {
	return fmt.Sprintf("{Unbond 🔓 %v",
		p.Validator.Fingerprint(),
	)
}
