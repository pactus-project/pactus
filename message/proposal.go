package message

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/vote"
)

type ProposalPayload struct {
	Proposal *vote.Proposal `cbor:"1,keyasint"`
}

func NewProposalMessage(proposal *vote.Proposal) *Message {
	return &Message{
		Type: PayloadTypeProposal,
		Payload: &ProposalPayload{
			Proposal: proposal,
		},
	}
}

func (p *ProposalPayload) SanityCheck() error {
	if err := p.Proposal.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	return nil
}

func (p *ProposalPayload) Type() PayloadType {
	return PayloadTypeProposal
}

func (p *ProposalPayload) Fingerprint() string {
	return p.Proposal.Fingerprint()
}
