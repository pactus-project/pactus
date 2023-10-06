package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
)

type UnbondPayload struct {
	Validator crypto.Address
}

func (p *UnbondPayload) Type() Type {
	return TypeUnbond
}

func (p *UnbondPayload) Signer() crypto.Address {
	return p.Validator
}

func (p *UnbondPayload) Value() int64 {
	return 0
}

// TODO: write test for me.
func (p *UnbondPayload) BasicCheck() error {
	if !p.Validator.IsValidatorAddress() {
		return BasicCheckError{
			Reason: "address is not a validator address",
		}
	}

	return nil
}

func (p *UnbondPayload) SerializeSize() int {
	return 21
}

func (p *UnbondPayload) Encode(w io.Writer) error {
	return p.Validator.Encode(w)
}

func (p *UnbondPayload) Decode(r io.Reader) error {
	return p.Validator.Decode(r)
}

func (p *UnbondPayload) String() string {
	return fmt.Sprintf("{Unbond 🔓 %v",
		p.Validator,
	)
}

func (p *UnbondPayload) Receiver() *crypto.Address {
	return nil
}
