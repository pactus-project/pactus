package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type UnbondPayload struct {
	Validator crypto.Address `cbor:"1,keyasint"`
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
		return errors.Errorf(errors.ErrInvalidTx, "invalid validator address")
	}

	return nil
}

func (p *UnbondPayload) Fingerprint() string {
	return fmt.Sprintf("{Unbond 🔓 %v",
		p.Validator.Fingerprint(),
	)
}
