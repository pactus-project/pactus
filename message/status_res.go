package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type StatusResPayload struct {
	Version version.Version `cbor:"1,keyasint"`
	Height  int             `cbor:"2,keyasint"`
}

func NewStatusResMessage(height int) Message {
	return Message{
		Type: PayloadTypeStatusRes,
		Payload: &StatusResPayload{
			Version: version.NodeVersion,
			Height:  height,
		},
	}

}
func (p *StatusResPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *StatusResPayload) Type() PayloadType {
	return PayloadTypeStatusRes
}

func (p *StatusResPayload) Fingerprint() string {
	return fmt.Sprintf("{%v}", p.Height)
}
