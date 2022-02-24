package message

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/errors"
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
		return errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	return nil
}

func (m *ProposalMessage) Type() Type {
	return MessageTypeProposal
}

func (m *ProposalMessage) Fingerprint() string {
	return m.Proposal.Fingerprint()
}
