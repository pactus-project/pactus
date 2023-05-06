package consensus

import (
	"sync"

	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

// The `mediator“ interface defines a mechanism for setting proposals and votes
// between independent consensus instances.
type mediator interface {
	OnPublishProposal(from Consensus, proposal *proposal.Proposal)
	OnPublishVote(from Consensus, vote *vote.Vote)
	OnBlockAnnounce(from Consensus)
	Register(cons Consensus)
}

// ConcreteMediator struct
type ConcreteMediator struct {
	lk *sync.RWMutex

	instances []Consensus
}

func newMediator(lk *sync.RWMutex) mediator {
	return &ConcreteMediator{lk: lk}
}

func (m *ConcreteMediator) OnPublishProposal(from Consensus, proposal *proposal.Proposal) {
	m.lk.Lock()
	defer m.lk.Unlock()

	for _, cons := range m.instances {
		if cons != from {
			cons.SetProposal(proposal)
		}
	}
}

func (m *ConcreteMediator) OnPublishVote(from Consensus, vote *vote.Vote) {
	m.lk.Lock()
	defer m.lk.Unlock()

	for _, cons := range m.instances {
		if cons != from {
			cons.AddVote(vote)
		}
	}
}

func (m *ConcreteMediator) OnBlockAnnounce(from Consensus) {
	m.lk.Lock()
	defer m.lk.Unlock()

	for _, cons := range m.instances {
		if cons != from {
			cons.MoveToNewHeight()
		}
	}
}

// Register a new Consensus instance to the mediator
func (m *ConcreteMediator) Register(cons Consensus) {
	m.instances = append(m.instances, cons)
}
