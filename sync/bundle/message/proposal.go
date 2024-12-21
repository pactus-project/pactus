package message

import (
	"github.com/pactus-project/pactus/network"
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

func (*ProposalMessage) BasicCheck() error {
	// Basic checks for the proposal are deferred to the consensus phase
	// to avoid unnecessary validation for validators outside the committee.
	return nil
}

func (*ProposalMessage) Type() Type {
	return TypeProposal
}

func (*ProposalMessage) TopicID() network.TopicID {
	return network.TopicIDConsensus
}

func (*ProposalMessage) ShouldBroadcast() bool {
	return true
}

func (m *ProposalMessage) ConsensusHeight() uint32 {
	return m.Height()
}

func (m *ProposalMessage) Height() uint32 {
	return m.Proposal.Height()
}

func (m *ProposalMessage) String() string {
	return m.Proposal.String()
}
