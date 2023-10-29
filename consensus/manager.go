package consensus

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/logger"
	"golang.org/x/exp/slices"
)

type manager struct {
	instances []Consensus

	// Caching future votes and proposals due to potential server time misalignments.
	// Votes and proposals for upcoming blocks may be received before
	// the current block's consensus is complete.
	upcomingVotes     []*vote.Vote         // Map to cache votes for future block heights
	upcomingProposals []*proposal.Proposal // Map to cache proposals for future block heights
	state             state.Facade
}

// NewManager creates a new manager instance that manages a set of consensus instances,
// each associated with a validator key and a reward address.
// It is not thread-safe.
func NewManager(
	conf *Config,
	state state.Facade,
	valKeys []*bls.ValidatorKey,
	rewardAddrs []crypto.Address,
	broadcastCh chan message.Message,
) Manager {
	mgr := &manager{
		instances:         make([]Consensus, len(valKeys)),
		upcomingVotes:     make([]*vote.Vote, 0),
		upcomingProposals: make([]*proposal.Proposal, 0),
		state:             state,
	}
	mediatorConcrete := newConcreteMediator()

	for i, key := range valKeys {
		cons := NewConsensus(conf, state, key, rewardAddrs[i], broadcastCh, mediatorConcrete)

		mgr.instances[i] = cons
	}

	return mgr
}

// Start starts the manager.
func (mgr *manager) Start() error {
	logger.Debug("starting consensus instances")

	for _, cons := range mgr.instances {
		cons.Start()
	}

	return nil
}

// Stop stops the manager.
func (mgr *manager) Stop() {
}

// Instances return all consensus instances that are read-only and
// can be safely accessed without modifying their state.
func (mgr *manager) Instances() []Reader {
	readers := make([]Reader, len(mgr.instances))
	for i, cons := range mgr.instances {
		readers[i] = cons
	}

	return readers
}

// PickRandomVote returns a random vote from a random consensus instance.
func (mgr *manager) PickRandomVote(round int16) *vote.Vote {
	cons := mgr.getBestInstance()
	return cons.PickRandomVote(round)
}

// Proposal returns the proposal for a specific round from a random consensus instance.
func (mgr *manager) Proposal() *proposal.Proposal {
	cons := mgr.getBestInstance()
	return cons.Proposal()
}

// HeightRound retrieves the current height and round from a random consensus instance.
func (mgr *manager) HeightRound() (uint32, int16) {
	cons := mgr.getBestInstance()
	return cons.HeightRound()
}

// HasActiveInstance checks if any of the consensus instances are currently active.
func (mgr *manager) HasActiveInstance() bool {
	for _, cons := range mgr.instances {
		if cons.IsActive() {
			return true
		}
	}

	return false
}

// MoveToNewHeight moves all consensus instances to a new height.
func (mgr *manager) MoveToNewHeight() {
	for _, cons := range mgr.instances {
		cons.MoveToNewHeight()
	}

	inst := mgr.getBestInstance()
	curHeight, _ := inst.HeightRound()

	for i := len(mgr.upcomingProposals) - 1; i >= 0; i-- {
		p := mgr.upcomingProposals[i]

		switch {
		case p.Height() < curHeight:
			// Ignore old proposals

		case p.Height() > curHeight:
			// keep this vote
			continue

		case p.Height() == curHeight:
			logger.Debug("upcoming proposal processed", "height", curHeight)

			for _, cons := range mgr.instances {
				cons.SetProposal(p)
			}
		}

		mgr.upcomingProposals = slices.Delete(mgr.upcomingProposals, i, i+1)
	}

	for i := len(mgr.upcomingVotes) - 1; i >= 0; i-- {
		v := mgr.upcomingVotes[i]

		switch {
		case v.Height() < curHeight:
			// Ignore old votes

		case v.Height() > curHeight:
			// keep this vote
			continue

		case v.Height() == curHeight:
			logger.Debug("upcoming votes processed", "height", curHeight)

			for _, cons := range mgr.instances {
				cons.AddVote(v)
			}
		}

		mgr.upcomingVotes = slices.Delete(mgr.upcomingVotes, i, i+1)
	}
}

// AddVote adds a vote to all consensus instances.
func (mgr *manager) AddVote(v *vote.Vote) {
	inst := mgr.getBestInstance()
	curHeight, _ := inst.HeightRound()

	switch {
	case v.Height() < curHeight:
		_ = mgr.state.UpdateLastCertificate(v)

	case v.Height() > curHeight:
		mgr.upcomingVotes = append(mgr.upcomingVotes, v)

	case v.Height() == curHeight:
		for _, cons := range mgr.instances {
			cons.AddVote(v)
		}
	}
}

// SetProposal sets the proposal for all consensus instances.
func (mgr *manager) SetProposal(p *proposal.Proposal) {
	inst := mgr.getBestInstance()
	curHeight, _ := inst.HeightRound()

	switch {
	case p.Height() < curHeight:
		// discard the old proposal

	case p.Height() > curHeight:
		mgr.upcomingProposals = append(mgr.upcomingProposals, p)

	case p.Height() == curHeight:
		for _, cons := range mgr.instances {
			cons.SetProposal(p)
		}
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
