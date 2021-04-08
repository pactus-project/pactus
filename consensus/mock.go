package consensus

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
)

var _ Consensus = &MockConsensus{}

type MockConsensus struct {
	Votes     []*vote.Vote
	Proposal  *proposal.Proposal
	Scheduled bool
	State     *state.MockState
	Round     int
	Lock      deadlock.RWMutex
}

func MockingConsensus(state *state.MockState) *MockConsensus {
	return &MockConsensus{State: state}
}

func (m *MockConsensus) MoveToNewHeight() {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Scheduled = true
}
func (m *MockConsensus) Stop() {}

func (m *MockConsensus) AddVote(v *vote.Vote) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Votes = append(m.Votes, v)
}
func (m *MockConsensus) RoundVotes(round int) []*vote.Vote {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	votes := make([]*vote.Vote, 0)
	for _, v := range m.Votes {
		if v.Round() == round {
			votes = append(votes, v)
		}
	}
	return votes
}
func (m *MockConsensus) SetProposal(p *proposal.Proposal) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Proposal = p
}
func (m *MockConsensus) RoundProposal(round int) *proposal.Proposal {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if m.Proposal == nil || m.Proposal.Round() != round {
		return nil
	}
	return m.Proposal
}
func (m *MockConsensus) HRS() *hrs.HRS {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	return hrs.NewHRS(m.State.LastBlockHeight()+1, m.Round, hrs.StepTypeNewHeight)
}
func (m *MockConsensus) Fingerprint() string {
	return ""
}
func (m *MockConsensus) PickRandomVote() *vote.Vote {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if len(m.Votes) == 0 {
		return nil
	}
	r := util.RandInt(len(m.Votes))
	return m.Votes[r]
}
