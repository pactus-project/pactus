package consensus

import (
	"sync"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type manager struct {
	lk sync.RWMutex

	instances []Consensus
}

// NewManager creates a new manager instance that manages a set of consensus instances,
// each associated with a signer and reward address.
func NewManager(
	conf *Config,
	state state.Facade,
	signers []crypto.Signer,
	rewardAddrs []crypto.Address,
	broadcastCh chan message.Message) Manager {
	mgr := &manager{
		instances: make([]Consensus, len(signers)),
	}
	mediator := newMediator(&mgr.lk)

	for i, signer := range signers {
		cons := NewConsensus(conf, state, signer, rewardAddrs[i], broadcastCh, mediator)

		mgr.instances[i] = cons
	}

	return mgr
}

// Start starts the manager.
func (mgr *manager) Start() error {
	return nil
}

// Stop stops the manager.
func (mgr *manager) Stop() {
}

// Instances returns all consensus instances that are read-only and
// can be safely accessed without modifying their state.
func (mgr *manager) Instances() []Reader {
	readers := make([]Reader, len(mgr.instances))
	for i, cons := range mgr.instances {
		readers[i] = cons
	}
	return readers
}

// PickRandomVote retrieves a random vote from a random consensus instance.
func (mgr *manager) PickRandomVote() *vote.Vote {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	cons := mgr.getBestInstance()
	return cons.PickRandomVote()
}

// RoundProposal retrieves the proposal for a specific round from a random consensus instance.
func (mgr *manager) RoundProposal(round int16) *proposal.Proposal {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	cons := mgr.getBestInstance()
	return cons.RoundProposal(round)
}

// HeightRound retrieves the current height and round from a random consensus instance.
func (mgr *manager) HeightRound() (uint32, int16) {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	cons := mgr.getBestInstance()
	return cons.HeightRound()
}

// HasActiveInstance checks if any of the consensus instances are currently active.
func (mgr *manager) HasActiveInstance() bool {
	mgr.lk.RLock()
	defer mgr.lk.RUnlock()

	for _, cons := range mgr.instances {
		if cons.IsActive() {
			return true
		}
	}

	return false
}

// MoveToNewHeight moves all consensus instances to a new height.
func (mgr *manager) MoveToNewHeight() {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	for _, cons := range mgr.instances {
		cons.MoveToNewHeight()
	}
}

// AddVote adds a vote to all consensus instances.
func (mgr *manager) AddVote(v *vote.Vote) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	for _, cons := range mgr.instances {
		cons.AddVote(v)
	}
}

// SetProposal sets the proposal for all consensus instances.
func (mgr *manager) SetProposal(proposal *proposal.Proposal) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	for _, cons := range mgr.instances {
		cons.SetProposal(proposal)
	}
}

// getBestInstance iterates through all consensus instances and returns the instance
// that is currently active, if there is one.
// If there are no active instances, it returns the first instance.
//
// Note that all active instances are assumed to be in the same state, and all inactive
// instances are assumed to be in the same state as well.
func (mgr *manager) getBestInstance() Consensus {
	for _, cons := range mgr.instances {
		if cons.IsActive() {
			return cons
		}
	}

	return mgr.instances[0]
}

func (mgr *manager) OnPublishProposal(from Consensus, proposal *proposal.Proposal) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	for _, cons := range mgr.instances {
		if cons != from {
			cons.SetProposal(proposal)
		}
	}
}

func (mgr *manager) OnPublishVote(from Consensus, vote *vote.Vote) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	for _, cons := range mgr.instances {
		if cons != from {
			cons.AddVote(vote)
		}
	}
}

func (mgr *manager) OnBlockAnnounce(from Consensus) {
	mgr.lk.Lock()
	defer mgr.lk.Unlock()

	for _, cons := range mgr.instances {
		if cons != from {
			cons.MoveToNewHeight()
		}
	}
}
