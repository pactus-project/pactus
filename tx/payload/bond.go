package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/errors"
)

type BondPayload struct {
	Sender    crypto.Address `cbor:"1,keyasint"`
	PublicKey *bls.PublicKey `cbor:"2,keyasint"`
	Stake     int64          `cbor:"3,keyasint"`
}

func (p *BondPayload) Type() Type {
	return PayloadTypeBond
}

func (p *BondPayload) Signer() crypto.Address {
	return p.Sender
}

func (p *BondPayload) Value() int64 {
	return p.Stake
}

func (p *BondPayload) SanityCheck() error {
	if p.Stake < 0 {
		return errors.Errorf(errors.ErrInvalidTx, "invalid amount")
	}
	if err := p.Sender.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "invalid sender address")
	}
	if err := p.PublicKey.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidTx, "invalid receiver address")
	}

	return nil
}

func (p *BondPayload) Fingerprint() string {
	return fmt.Sprintf("{Bond 🔐 %v->%v %v",
		p.Sender.Fingerprint(),
		p.PublicKey.Address().Fingerprint(),
		p.Stake)
}
