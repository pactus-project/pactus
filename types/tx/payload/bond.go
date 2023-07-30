package payload

import (
	"errors"
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/encoding"
)

type BondPayload struct {
	Sender    crypto.Address
	Receiver  crypto.Address
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
	if err := p.Sender.SanityCheck(); err != nil {
		return err
	}
	if err := p.Receiver.SanityCheck(); err != nil {
		return err
	}
	if p.PublicKey != nil {
		if err := p.PublicKey.VerifyAddress(p.Receiver); err != nil {
			return err
		}
	}

	return nil
}

func (p *BondPayload) SerializeSize() int {
	if p.PublicKey != nil {
		return 139 + encoding.VarIntSerializeSize(uint64(p.Stake))
	}
	return 43 + encoding.VarIntSerializeSize(uint64(p.Stake))
}

func (p *BondPayload) Encode(w io.Writer) error {
	err := encoding.WriteElements(w, &p.Sender, &p.Receiver)
	if err != nil {
		return err
	}
	if p.PublicKey != nil {
		err := encoding.WriteElements(w, uint8(bls.PublicKeySize))
		if err != nil {
			return err
		}
		err = p.PublicKey.Encode(w)
		if err != nil {
			return err
		}
	} else {
		err := encoding.WriteElements(w, uint8(0))
		if err != nil {
			return err
		}
	}

	return encoding.WriteVarInt(w, uint64(p.Stake))
}

func (p *BondPayload) Decode(r io.Reader) error {
	err := encoding.ReadElements(r, &p.Sender, &p.Receiver)
	if err != nil {
		return err
	}
	pubKeySize, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	if pubKeySize == bls.PublicKeySize {
		p.PublicKey = new(bls.PublicKey)
		err = p.PublicKey.Decode(r)
		if err != nil {
			return err
		}
	} else if pubKeySize != 0 {
		return errors.New("invalid public key size")
	}

	stake, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	p.Stake = int64(stake)
	return nil
}

func (p *BondPayload) String() string {
	return fmt.Sprintf("{Bond 🔐 %v->%v %v",
		p.Sender.ShortString(),
		p.Receiver.ShortString(),
		p.Stake)
}
