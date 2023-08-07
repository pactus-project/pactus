package message

import (
	"github.com/pactus-project/pactus/types/proposal"
)

type ProposalMessage struct {
	Proposal *proposal.Proposal `cbor:"1,keyasint"`
}

func NewProposalMessage(p *proposal.Proposal) *ProposalMessage {
	return &ProposalMessage{
		Proposal: p,
	}
}

func (m *ProposalMessage) SanityCheck() error {
	return m.Proposal.SanityCheck()
}

func (m *ProposalMessage) Type() Type {
	return TypeProposal
}

func (m *ProposalMessage) String() string {
	return m.Proposal.String()
}
