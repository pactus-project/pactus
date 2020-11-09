package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type SalamPayload struct {
	Version     version.Version `cbor:"1,keyasint"`
	GenesisHash crypto.Hash     `cbor:"2,keyasint"`
	Height      int             `cbor:"3,keyasint"`
}

func NewSalamMessage(genesisHash crypto.Hash, height int) *Message {
	return &Message{
		Type: PayloadTypeSalam,
		Payload: &SalamPayload{
			Version:     version.NodeVersion,
			GenesisHash: genesisHash,
			Height:      height,
		},
	}

}
func (p *SalamPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *SalamPayload) Type() PayloadType {
	return PayloadTypeSalam
}

func (p *SalamPayload) Fingerprint() string {
	return fmt.Sprintf("{%v}", p.Height)
}
