package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type AleykPayload struct {
	Version     version.Version `cbor:"1,keyasint"`
	GenesisHash crypto.Hash     `cbor:"2,keyasint"`
	Height      int             `cbor:"3,keyasint"`
}

func (p *AleykPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *AleykPayload) Type() PayloadType {
	return PayloadTypeAleyk
}

func (p *AleykPayload) Fingerprint() string {
	return fmt.Sprintf("{%v}", p.Height)
}
