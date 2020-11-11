package message

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"
)

type ProposalReqPayload struct {
	Height int `cbor:"1,keyasint"`
	Round  int `cbor:"2,keyasint"`
}

func NewProposalReqMessage(height, round int) *Message {
	return &Message{
		Type: PayloadTypeProposalReq,
		Payload: &ProposalReqPayload{
			Height: height,
			Round:  round,
		},
	}
}

func (p *ProposalReqPayload) SanityCheck() error {
	if p.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid height")
	}
	if p.Round < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid round")
	}

	return nil
}

func (p *ProposalReqPayload) Type() PayloadType {
	return PayloadTypeProposalReq
}

func (p *ProposalReqPayload) Fingerprint() string {
	return fmt.Sprintf("%v/%v", p.Height, p.Round)
}
