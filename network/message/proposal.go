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
		Type:   PayloadTypeProposal,
		Height: proposal.Height(),
		Payload: &ProposalPayload{
			Proposal: proposal,
		},
	}
}

func (m *ProposalPayload) SanityCheck() error {
	if err := m.Proposal.SanityCheck(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}
	return nil
}

func (m *ProposalPayload) Type() PayloadType {
	return PayloadTypeProposal
}

func (p *ProposalPayload) Fingerprint() string {
	return p.Proposal.Fingerprint()
}
