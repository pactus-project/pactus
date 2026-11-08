package manager

import (
	"context"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/consensusv2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"golang.org/x/exp/slices"
)

type manager struct {
	instances []consensus.Consensus

	// Caching future votes and proposals due to potential server time misalignments.
	// Votes and proposals for upcoming blocks may be received before
	// the current block's consensus is complete.
	upcomingVotes     []*vote.Vote         // Map to cache votes for future block heights
	upcomingProposals []*proposal.Proposal // Map to cache proposals for future block heights
	state             state.State

	// observer keeps the votes and proposals of the current height while no
	// local validator is in the committee. It is off unless WithObserver is set.
	observer *observer
}

// Option configures a consensus manager.
type Option func(*manager)

// WithObserver turns on the observer mode: while none of the local validators
// is in the committee, the manager verifies and keeps the votes and proposals
// of the current height, without taking part in consensus. See Observation.
func WithObserver(enabled bool) Option {
	return func(mgr *manager) {
		mgr.observer = newObserver(enabled)
	}
}

// maxUpcomingVotes and maxUpcomingProposals bound upcomingVotes and
// upcomingProposals. Neither cache checks a signature, a committee
// membership, or even a duplicate before caching an entry for height+1
// (unlike the current-height path, which forwards into each consensus
// instance's own, checked AddVote/SetProposal), so an unauthenticated peer
// could otherwise grow them without limit just by sending future-height
// votes or proposals. Comfortably above what every real committee member
// sending every vote type once for the round could produce.
const (
	maxUpcomingVotes     = 512
	maxUpcomingProposals = 100
)

// NewManagerV1 creates a new manager instance that manages a set of consensus instances,
// each associated with a validator key and a reward address.
// It is not thread-safe.
func NewManagerV1(
	ctx context.Context,
	conf *consensus.Config,
	state state.State,
	valKeys []*bls.ValidatorKey,
	rewardAddrs []crypto.Address,
	broadcastPipe pipeline.Pipeline[message.Message],
	opts ...Option,
) ConsensusManager {
	instances := make([]consensus.Consensus, len(valKeys))
	mediatorConcrete := consensus.NewConcreteMediator()
	for i, key := range valKeys {
		instances[i] = consensus.NewConsensus(
			ctx, conf, state, key, rewardAddrs[i], broadcastPipe, mediatorConcrete,
		)
	}

	return NewManager(state, instances, opts...)
}

func NewManagerV2(
	ctx context.Context,
	conf *consensusv2.Config,
	state state.State,
	valKeys []*bls.ValidatorKey,
	rewardAddrs []crypto.Address,
	broadcastPipe pipeline.Pipeline[message.Message],
	opts ...Option,
) ConsensusManager {
	instances := make([]consensus.Consensus, len(valKeys))
	mediatorConcrete := consensus.NewConcreteMediator()

	for i, key := range valKeys {
		instances[i] = consensusv2.NewConsensus(
			ctx, conf, state, key, rewardAddrs[i], broadcastPipe, mediatorConcrete,
		)
	}

	return NewManager(state, instances, opts...)
}

func NewManager(state state.State, instances []consensus.Consensus, opts ...Option) ConsensusManager {
	mgr := &manager{
		instances:         instances,
		upcomingVotes:     make([]*vote.Vote, 0),
		upcomingProposals: make([]*proposal.Proposal, 0),
		state:             state,
		observer:          newObserver(false),
	}
	for _, opt := range opts {
		opt(mgr)
	}

	return mgr
}

// Instances return all consensus instances that are read-only and
// can be safely accessed without modifying their state.
func (mgr *manager) Instances() []consensus.ConsensusReader {
	readers := make([]consensus.ConsensusReader, len(mgr.instances))
	for i, cons := range mgr.instances {
		readers[i] = cons
	}

	return readers
}

// Proposal returns the current proposal for the active round from a random consensus instance.
func (mgr *manager) Proposal() *proposal.Proposal {
	cons := mgr.getBestInstance()

	return cons.Proposal()
}

// HandleQueryProposal returns the proposal for a specific round from a random consensus instance.
func (mgr *manager) HandleQueryProposal(height types.Height, round types.Round) *proposal.Proposal {
	cons := mgr.getBestInstance()

	return cons.HandleQueryProposal(height, round)
}

// HandleQueryVote returns a random vote from a random consensus instance.
func (mgr *manager) HandleQueryVote(height types.Height, round types.Round) *vote.Vote {
	cons := mgr.getBestInstance()

	return cons.HandleQueryVote(height, round)
}

// HeightRound retrieves the current height and round from a random consensus instance.
func (mgr *manager) HeightRound() (types.Height, types.Round) {
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
	stateHeight := mgr.state.LastBlockHeight()
	consHeight, _ := mgr.getBestInstance().HeightRound()
	if stateHeight < consHeight {
		return
	}

	for _, cons := range mgr.instances {
		cons.MoveToNewHeight()
	}

	if mgr.observer.enabled {
		newHeight, _ := mgr.getBestInstance().HeightRound()
		mgr.observer.moveToHeight(newHeight, mgr.state.CommitteeValidators(),
			mgr.state.Proposer, mgr.HasActiveInstance())
	}

	for index := len(mgr.upcomingProposals) - 1; index >= 0; index-- {
		prop := mgr.upcomingProposals[index]
		for _, cons := range mgr.instances {
			cons.SetProposal(prop)
		}
		mgr.observer.setProposal(prop)

		mgr.upcomingProposals = slices.Delete(mgr.upcomingProposals, index, index+1)
	}

	for index := len(mgr.upcomingVotes) - 1; index >= 0; index-- {
		vote := mgr.upcomingVotes[index]
		for _, cons := range mgr.instances {
			cons.AddVote(vote)
		}
		mgr.observer.addVote(vote)

		mgr.upcomingVotes = slices.Delete(mgr.upcomingVotes, index, index+1)
	}
}

// AddVote adds a vote to all consensus instances.
func (mgr *manager) AddVote(vote *vote.Vote) {
	cons := mgr.getBestInstance()
	curHeight, _ := cons.HeightRound()
	switch {
	case vote.Height() == curHeight-1:
		_ = mgr.state.UpdateLastCertificate(vote)

	case vote.Height() == curHeight+1:
		if len(mgr.upcomingVotes) < maxUpcomingVotes {
			mgr.upcomingVotes = append(mgr.upcomingVotes, vote)
		}

	case vote.Height() == curHeight:
		for _, cons := range mgr.instances {
			cons.AddVote(vote)
		}
		mgr.observer.addVote(vote)
	}
}

// SetProposal sets the proposal for all consensus instances.
func (mgr *manager) SetProposal(prop *proposal.Proposal) {
	cons := mgr.getBestInstance()
	curHeight, _ := cons.HeightRound()
	switch {
	case prop.Height() == curHeight+1:
		if len(mgr.upcomingProposals) < maxUpcomingProposals {
			mgr.upcomingProposals = append(mgr.upcomingProposals, prop)
		}

	case prop.Height() == curHeight:
		for _, cons := range mgr.instances {
			cons.SetProposal(prop)
		}
		mgr.observer.setProposal(prop)
	}
}

// Observation returns what the observer saw for the current height, or nil
// when the observer is off or a local validator is in the committee.
// It is safe to call from any goroutine.
func (mgr *manager) Observation() *Observation {
	return mgr.observer.observation()
}

// getBestInstance iterates through all consensus instances and returns the instance
// that is currently active, if there is one.
// If there are no active instances, it returns the first instance.
//
// Note that all active instances are assumed to be in the same state, and all inactive
// instances are assumed to be in the same state as well.
func (mgr *manager) getBestInstance() consensus.ConsensusReader {
	for _, cons := range mgr.instances {
		if cons.IsActive() {
			return cons
		}
	}

	return mgr.instances[0]
}

// IsDeprecated checks if any of the consensus instances are deprecated.
func (mgr *manager) IsDeprecated() bool {
	return mgr.instances[0].IsDeprecated()
}
