package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/version"
)

type StatusReqPayload struct {
	Version version.Version `cbor:"1,keyasint"`
	Height  int             `cbor:"2,keyasint"`
}

func NewStatusReqMessage(height int) Message {
	return Message{
		Type: PayloadTypeStatusReq,
		Payload: &StatusReqPayload{
			Version: version.NodeVersion,
			Height:  height,
		},
	}

}
func (p *StatusReqPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	return nil
}

func (p *StatusReqPayload) Type() PayloadType {
	return PayloadTypeStatusReq
}

func (p *StatusReqPayload) Fingerprint() string {
	return fmt.Sprintf("{%v}", p.Height)
}
