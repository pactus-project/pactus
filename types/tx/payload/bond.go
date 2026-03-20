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

	// Delegation fields
	DelegateOwner  crypto.Address
	DelegateShare  amount.Amount
	DelegateExpiry uint32
}

func (*BondPayload) Type() Type {
	return TypeBond
}

func (p *BondPayload) Signer() crypto.Address {
	return p.From
}

func (p *BondPayload) Value() amount.Amount {
	return p.Stake
}

func (p *BondPayload) IsDelegated() bool {
	return p.DelegateOwner != crypto.TreasuryAddress
}

// BasicCheck performs basic checks on the Bond payload.
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

	if p.IsDelegated() {
		if !p.DelegateOwner.IsAccountAddress() {
			return BasicCheckError{
				Reason: "delegate owner is not an account address: " + p.DelegateOwner.String(),
			}
		}

		if p.DelegateShare < 0 || p.DelegateShare > 7e8 {
			return BasicCheckError{
				Reason: "delegate share must be between 0 and 0.7 PAC",
			}
		}
	}

	return nil
}

func (p *BondPayload) SerializeSize() int {
	size := 43 + encoding.VarIntSerializeSize(uint64(p.Stake))
	if p.PublicKey != nil {
		size += 96 // pubkey size
	}
	if p.IsDelegated() {
		// delegate owner size (21) + delegate share size (var) + delegate expiry size (4)
		size += 21 + encoding.VarIntSerializeSize(uint64(p.DelegateShare)) + 4
	}

	return size
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

	if err := encoding.WriteVarInt(w, uint64(p.Stake)); err != nil {
		return err
	}

	if p.IsDelegated() {
		if err := p.DelegateOwner.Encode(w); err != nil {
			return err
		}
		if err := encoding.WriteVarInt(w, uint64(p.DelegateShare)); err != nil {
			return err
		}
		if err := encoding.WriteElement(w, p.DelegateExpiry); err != nil {
			return err
		}
	}

	return nil
}

func (p *BondPayload) Decode(ctx DecodeContext, r io.Reader) error {
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

	if ctx.WithDelegation {
		if err := p.DelegateOwner.Decode(r); err != nil {
			return err
		}

		share, err := encoding.ReadVarInt(r)
		if err != nil {
			return err
		}
		p.DelegateShare = amount.Amount(share)

		if err := encoding.ReadElement(r, &p.DelegateExpiry); err != nil {
			return err
		}
	}

	return nil
}

// LogString returns a concise string representation intended for use in logs.
func (p *BondPayload) LogString() string {
	return fmt.Sprintf("{Bond 🔐 %s->%s %s",
		p.From.LogString(),
		p.To.LogString(),
		p.Stake)
}
