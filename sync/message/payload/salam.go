package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type SalamPayload struct {
	NodeVersion version.Version  `cbor:"1,keyasint"`
	Moniker     string           `cbor:"2,keyasint"`
	PublicKey   crypto.PublicKey `cbor:"3,keyasint"`
	GenesisHash crypto.Hash      `cbor:"4,keyasint"`
	Height      int              `cbor:"5,keyasint"`
	Flags       int              `cbor:"6,keyasint"`
}

func NewSalamPayload(moniker string,
	publicKey crypto.PublicKey, genesisHash crypto.Hash,
	height int, flags int) Payload {
	return &SalamPayload{
		NodeVersion: version.NodeVersion,
		Moniker:     moniker,
		PublicKey:   publicKey,
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

func (p *SalamPayload) Type() PayloadType {
	return PayloadTypeSalam
}

func (p *SalamPayload) Fingerprint() string {
	return fmt.Sprintf("{%s %v}", p.Moniker, p.Height)
}
