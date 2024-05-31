package consensus

import (
	"sync"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
)

var _ Consensus = &MockConsensus{}

type MockConsensus struct {
	// This locks prevents the Data Race in tests
	lk sync.RWMutex
	ts *testsuite.TestSuite

	ValKey      *bls.ValidatorKey
	Votes       []*vote.Vote
	CurProposal *proposal.Proposal
	Active      bool
	Proposer    bool
	Height      uint32
	Round       int16
}

func MockingManager(ts *testsuite.TestSuite, valKeys []*bls.ValidatorKey) (Manager, []*MockConsensus) {
	mocks := make([]*MockConsensus, len(valKeys))
	instances := make([]Consensus, len(valKeys))
	for i, s := range valKeys {
		cons := MockingConsensus(ts, s)
		mocks[i] = cons
		instances[i] = cons
	}

	return &manager{
		instances:         instances,
		upcomingVotes:     make([]*vote.Vote, 0),
		upcomingProposals: make([]*proposal.Proposal, 0),
	}, mocks
}

func MockingConsensus(ts *testsuite.TestSuite, valKey *bls.ValidatorKey) *MockConsensus {
	return &MockConsensus{
		ts:     ts,
		ValKey: valKey,
	}
}

func (m *MockConsensus) ConsensusKey() *bls.PublicKey {
	return m.ValKey.PublicKey()
}

func (m *MockConsensus) MoveToNewHeight() {
	m.lk.Lock()
	defer m.lk.Unlock()

	m.Height++
}

func (m *MockConsensus) Start() {}

func (m *MockConsensus) AddVote(v *vote.Vote) {
	m.lk.Lock()
	defer m.lk.Unlock()

	m.Votes = append(m.Votes, v)
}

func (m *MockConsensus) AllVotes() []*vote.Vote {
	m.lk.Lock()
	defer m.lk.Unlock()

	return m.Votes
}

func (m *MockConsensus) SetProposal(p *proposal.Proposal) {
	m.lk.Lock()
	defer m.lk.Unlock()

	m.CurProposal = p
}

func (m *MockConsensus) HasVote(h hash.Hash) bool {
	m.lk.Lock()
	defer m.lk.Unlock()

	for _, v := range m.Votes {
		if v.Hash() == h {
			return true
		}
	}

	return false
}

func (m *MockConsensus) Proposal() *proposal.Proposal {
	m.lk.Lock()
	defer m.lk.Unlock()

	return m.CurProposal
}

func (m *MockConsensus) HeightRound() (uint32, int16) {
	m.lk.Lock()
	defer m.lk.Unlock()

	return m.Height, m.Round
}

func (m *MockConsensus) String() string {
	return ""
}

func (m *MockConsensus) PickRandomVote(_ int16) *vote.Vote {
	m.lk.Lock()
	defer m.lk.Unlock()

	if len(m.Votes) == 0 {
		return nil
	}
	r := m.ts.RandInt32(int32(len(m.Votes)))

	return m.Votes[r]
}

func (m *MockConsensus) IsActive() bool {
	m.lk.Lock()
	defer m.lk.Unlock()

	return m.Active
}

func (m *MockConsensus) IsProposer() bool {
	m.lk.Lock()
	defer m.lk.Unlock()

	return m.Proposer
}

func (m *MockConsensus) SetActive(active bool) {
	m.lk.Lock()
	defer m.lk.Unlock()

	m.Active = active
}
