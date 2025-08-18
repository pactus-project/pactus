package consensusv2

import (
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

// The `mediator“ interface defines a mechanism for setting proposals and votes
// between independent consensus instances.
type mediator interface {
	OnPublishProposal(from consensus.Consensus, prop *proposal.Proposal)
	OnPublishVote(from consensus.Consensus, vte *vote.Vote)
	OnBlockAnnounce(from consensus.Consensus)
	Register(cons consensus.Consensus)
}

// ConcreteMediator struct.
type ConcreteMediator struct {
	instances []consensus.Consensus
}

func newConcreteMediator() mediator {
	return &ConcreteMediator{}
}

func (m *ConcreteMediator) OnPublishProposal(from consensus.Consensus, prop *proposal.Proposal) {
	for _, cons := range m.instances {
		if cons != from {
			cons.SetProposal(prop)
		}
	}
}

func (m *ConcreteMediator) OnPublishVote(from consensus.Consensus, vte *vote.Vote) {
	for _, cons := range m.instances {
		if cons != from {
			cons.AddVote(vte)
		}
	}
}

func (m *ConcreteMediator) OnBlockAnnounce(from consensus.Consensus) {
	for _, cons := range m.instances {
		if cons != from {
			cons.MoveToNewHeight()
		}
	}
}

// Register a new Consensus instance to the mediator.
func (m *ConcreteMediator) Register(cons consensus.Consensus) {
	m.instances = append(m.instances, cons)
}
