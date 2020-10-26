package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type SalamPayload struct {
	Version version.Version `cbor:"1,keyasint"`
	Height  int             `cbor:"2,keyasint"`
}

func NewSalamMessage(height int) Message {
	return Message{
		Type: PayloadTypeSalam,
		Payload: &SalamPayload{
			Version: version.NodeVersion,
			Height:  height,
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
