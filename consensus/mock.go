package consensus

import (
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
)

var _ Consensus = &MockConsensus{}

type MockConsensus struct {
	ts *testsuite.TestSuite

	State       *state.MockState
	ValKey      *bls.ValidatorKey
	Votes       []*vote.Vote
	CurProposal *proposal.Proposal
	Active      bool
	Proposer    bool
	Height      uint32
	Round       int16
}

func MockingConsensus(ts *testsuite.TestSuite, state *state.MockState, valKey *bls.ValidatorKey) *MockConsensus {
	return &MockConsensus{
		ts:     ts,
		State:  state,
		ValKey: valKey,
	}
}

func (m *MockConsensus) ConsensusKey() *bls.PublicKey {
	return m.ValKey.PublicKey()
}

func (m *MockConsensus) MoveToNewHeight() {
	m.Height = m.State.LastBlockHeight() + 1
}

func (m *MockConsensus) AddVote(v *vote.Vote) {
	m.Votes = append(m.Votes, v)
}

func (m *MockConsensus) AllVotes() []*vote.Vote {
	return m.Votes
}

func (m *MockConsensus) SetProposal(p *proposal.Proposal) {
	m.CurProposal = p
}

func (m *MockConsensus) HasVote(h hash.Hash) bool {
	for _, v := range m.Votes {
		if v.Hash() == h {
			return true
		}
	}

	return false
}

func (m *MockConsensus) Proposal() *proposal.Proposal {
	return m.CurProposal
}

func (m *MockConsensus) HandleQueryProposal(_ uint32, _ int16) *proposal.Proposal {
	return m.CurProposal
}

func (m *MockConsensus) HeightRound() (uint32, int16) {
	return m.Height, m.Round
}

func (*MockConsensus) String() string {
	return ""
}

func (m *MockConsensus) HandleQueryVote(_ uint32, _ int16) *vote.Vote {
	if len(m.Votes) == 0 {
		return nil
	}
	r := m.ts.RandInt32(int32(len(m.Votes)))

	return m.Votes[r]
}

func (m *MockConsensus) IsActive() bool {
	return m.Active
}

func (m *MockConsensus) IsProposer() bool {
	return m.Proposer
}

func (m *MockConsensus) SetActive(active bool) {
	m.Active = active
}
