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
	NodeVersion string         `cbor:"1,keyasint"`
	Moniker     string         `cbor:"2,keyasint"`
	PublicKey   *bls.PublicKey `cbor:"3,keyasint"`
	GenesisHash hash.Hash      `cbor:"4,keyasint"`
	Height      int            `cbor:"5,keyasint"`
	Flags       int            `cbor:"6,keyasint"`
}

func NewSalamPayload(moniker string,
	publicKey crypto.PublicKey, genesisHash hash.Hash,
	height int, flags int) Payload {
	return &SalamPayload{
		NodeVersion: version.Version(),
		Moniker:     moniker,
		PublicKey:   publicKey.(*bls.PublicKey),
		GenesisHash: genesisHash,
		Height:      height,
		Flags:       flags,
	}
}

func (p *SalamPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := p.PublicKey.SanityCheck(); err != nil {
		return err
	}
	return nil
}

func (p *SalamPayload) Type() Type {
	return PayloadTypeSalam
}

func (p *SalamPayload) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", p.Moniker, p.Height)
}
