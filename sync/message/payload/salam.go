package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type SalamPayload struct {
	Agent       string         `cbor:"1,keyasint"`
	Moniker     string         `cbor:"2,keyasint"`
	PublicKey   *bls.PublicKey `cbor:"3,keyasint"`
	Signature   *bls.Signature `cbor:"4,keyasint"`
	Height      int            `cbor:"5,keyasint"`
	Flags       int            `cbor:"6,keyasint"`
	GenesisHash hash.Hash      `cbor:"7,keyasint"`
}

func NewSalamPayload(moniker string, pub crypto.PublicKey, sig crypto.Signature,
	height int, flags int, genesisHash hash.Hash) Payload {
	return &SalamPayload{
		Agent:       version.Agent(),
		Moniker:     moniker,
		PublicKey:   pub.(*bls.PublicKey),
		Signature:   sig.(*bls.Signature),
		GenesisHash: genesisHash,
		Height:      height,
		Flags:       flags,
	}
}

func (p *SalamPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid height")
	}
	if !p.PublicKey.Verify(p.PublicKey.RawBytes(), p.Signature) {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid signature")
	}
	return nil
}

func (p *SalamPayload) Type() Type {
	return PayloadTypeSalam
}

func (p *SalamPayload) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", p.Moniker, p.Height)
}
