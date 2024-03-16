package payload

import (
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/encoding"
)

type BondPayload struct {
	From      crypto.Address
	To        crypto.Address
	PublicKey *bls.PublicKey
	Stake     amount.Amount
}

func (p *BondPayload) Type() Type {
	return TypeBond
}

func (p *BondPayload) Signer() crypto.Address {
	return p.From
}

func (p *BondPayload) Value() amount.Amount {
	return p.Stake
}

func (p *BondPayload) BasicCheck() error {
	if !p.From.IsAccountAddress() {
		return BasicCheckError{
			Reason: "sender is not an account address: " + p.From.String(),
		}
	}

	if !p.To.IsValidatorAddress() {
		return BasicCheckError{
			Reason: "receiver is not a validator address: " + p.To.String(),
		}
	}

	if p.PublicKey != nil {
		if err := p.PublicKey.VerifyAddress(p.To); err != nil {
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
	err := p.From.Encode(w)
	if err != nil {
		return err
	}

	err = p.To.Encode(w)
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
	err := p.From.Decode(r)
	if err != nil {
		return err
	}

	err = p.To.Decode(r)
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
		return ErrInvalidPublicKeySize
	}

	stake, err := encoding.ReadVarInt(r)
	if err != nil {
		return err
	}
	p.Stake = amount.Amount(stake)

	return nil
}

func (p *BondPayload) String() string {
	return fmt.Sprintf("{Bond 🔐 %s->%s %s",
		p.From.ShortString(),
		p.To.ShortString(),
		p.Stake)
}

func (p *BondPayload) Receiver() *crypto.Address {
	return &p.To
}
