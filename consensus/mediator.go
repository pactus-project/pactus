package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

// mediator interface to set proposal and votes between the consensus instances
type mediator interface {
	OnPublishProposal(from Consensus, proposal *proposal.Proposal)
	OnPublishVote(from Consensus, vote *vote.Vote)
	Register(cons Consensus)
}

// ConcreteMediator struct
type ConcreteMediator struct {
	instances []Consensus
}

func newMediator() mediator {
	return &ConcreteMediator{}
}

func (m *ConcreteMediator) OnPublishProposal(from Consensus, proposal *proposal.Proposal) {
	for _, cons := range m.instances {
		if cons != from {
			cons.SetProposal(proposal)
		}
	}
}

func (m *ConcreteMediator) OnPublishVote(from Consensus, vote *vote.Vote) {
	for _, cons := range m.instances {
		if cons != from {
			cons.AddVote(vote)
		}
	}
}

// Register a new Consensus instance to the mediator
func (m *ConcreteMediator) Register(cons Consensus) {
	m.instances = append(m.instances, cons)
}
