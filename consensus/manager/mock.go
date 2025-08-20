package manager

import (
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
)

func MockingManager(ts *testsuite.TestSuite, state *state.MockState,
	valKeys []*bls.ValidatorKey,
) (Manager, []*consensus.MockConsensus) {
	mocks := make([]*consensus.MockConsensus, len(valKeys))
	instances := make([]Consensus, len(valKeys))
	for i, key := range valKeys {
		cons := consensus.MockingConsensus(ts, state, key)
		mocks[i] = cons
		instances[i] = cons
	}

	return &manager{
		instances:         instances,
		upcomingVotes:     make([]*vote.Vote, 0),
		upcomingProposals: make([]*proposal.Proposal, 0),
	}, mocks
}
