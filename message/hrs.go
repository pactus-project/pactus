package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/errors"
)

type HRSPayload struct {
	HRS hrs.HRS `cbor:"1,keyasint"`
}

func NewHRSMessage(hrs hrs.HRS) Message {
	return Message{
		Type: PayloadTypeHRS,
		Payload: &HRSPayload{
			HRS: hrs,
		},
	}
}

func (p *HRSPayload) SanityCheck() error {
	if !p.HRS.IsValid() {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid step")
	}
	return nil
}

func (p *HRSPayload) Type() PayloadType {
	return PayloadTypeHRS
}

func (p *HRSPayload) Fingerprint() string {
	return fmt.Sprintf("{%s}", p.HRS.Fingerprint())
}
