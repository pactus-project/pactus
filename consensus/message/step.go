package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/errors"
)

type StepPayload struct {
	Round int          `cbor:"1,keyasint"`
	Step  hrs.StepType `cbor:"2,keyasint"`
}

func NewStepMessage(hrs hrs.HRS) *Message {
	return &Message{
		Type:   PayloadTypeStep,
		Height: hrs.Height(),
		Payload: &StepPayload{
			Round: hrs.Round(),
			Step:  hrs.Step(),
		},
	}
}

func (p *StepPayload) SanityCheck() error {
	if !p.Step.IsValid() {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid step")
	}
	return nil
}

func (p *StepPayload) Type() PayloadType {
	return PayloadTypeStep
}

func (p *StepPayload) Fingerprint() string {
	return fmt.Sprintf("{%s}", p.Step)
}
