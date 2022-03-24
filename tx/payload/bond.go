package payload

import (
	"fmt"
	"io"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/encoding"
	"github.com/zarbchain/zarb-go/errors"
)

type BondPayload struct {
	Sender    crypto.Address
	PublicKey *bls.PublicKey
	Stake     int64
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
		return errors.Errorf(errors.ErrInvalidAmount, "invalid sender address")
	}
	if err := p.PublicKey.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidPublicKey, "invalid receiver public key")
	}

	return nil
}

func (p *BondPayload) Encode(w io.Writer) error {
	err := encoding.WriteElement(w, &p.Sender)
	if err != nil {
		return err
	}
	err = p.PublicKey.Encode(w)
	if err != nil {
		return err
	}
	return encoding.WriteVarInt(w, uint64(p.Stake))
}

func (p *BondPayload) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &p.Sender)
	if err != nil {
		return err
	}
	p.PublicKey = new(bls.PublicKey)
	err = p.PublicKey.Decode(r)
	if err != nil {
		return err
	}
	stake, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	p.Stake = int64(stake)
	return nil
}

func (p *BondPayload) Fingerprint() string {
	return fmt.Sprintf("{Bond 🔐 %v->%v %v",
		p.Sender.Fingerprint(),
		p.PublicKey.Address().Fingerprint(),
		p.Stake)
}
