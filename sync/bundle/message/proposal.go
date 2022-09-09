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
	if err := m.Proposal.SanityCheck(); err != nil {
		return err
	}

	return nil
}

func (m *ProposalMessage) Type() Type {
	return MessageTypeProposal
}

func (m *ProposalMessage) Fingerprint() string {
	return m.Proposal.Fingerprint()
}
