package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

// The `mediator“ interface defines a mechanism for setting proposals and votes
// between independent consensus instances.
type mediator interface {
	OnPublishProposal(from Consensus, prop *proposal.Proposal)
	OnPublishVote(from Consensus, vote *vote.Vote)
	OnBlockAnnounce(from Consensus)
	Register(cons Consensus)
}

// ConcreteMediator struct.
type ConcreteMediator struct {
	instances []Consensus
}

func newConcreteMediator() mediator {
	return &ConcreteMediator{}
}

func (m *ConcreteMediator) OnPublishProposal(from Consensus, prop *proposal.Proposal) {
	for _, cons := range m.instances {
		if cons != from {
			cons.SetProposal(prop)
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

func (m *ConcreteMediator) OnBlockAnnounce(from Consensus) {
	for _, cons := range m.instances {
		if cons != from {
			cons.MoveToNewHeight()
		}
	}
}

// Register a new Consensus instance to the mediator.
func (m *ConcreteMediator) Register(cons Consensus) {
	m.instances = append(m.instances, cons)
}
