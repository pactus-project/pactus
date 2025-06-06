package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
)

type UnbondPayload struct {
	Validator crypto.Address
}

func (*UnbondPayload) Type() Type {
	return TypeUnbond
}

func (p *UnbondPayload) Signer() crypto.Address {
	return p.Validator
}

func (*UnbondPayload) Value() amount.Amount {
	return 0
}

// BasicCheck performs basic checks on the Unbond payload.
func (p *UnbondPayload) BasicCheck() error {
	if !p.Validator.IsValidatorAddress() {
		return BasicCheckError{
			Reason: "address is not a validator address",
		}
	}

	return nil
}

func (*UnbondPayload) SerializeSize() int {
	return 21
}

func (p *UnbondPayload) Encode(w io.Writer) error {
	return p.Validator.Encode(w)
}

func (p *UnbondPayload) Decode(r io.Reader) error {
	return p.Validator.Decode(r)
}

func (p *UnbondPayload) String() string {
	return fmt.Sprintf("{Unbond 🔓 %s",
		p.Validator.ShortString(),
	)
}
