package consensus

import (
	"sync"

	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/types/proposal"
	"github.com/zarbchain/zarb-go/types/vote"
	"github.com/zarbchain/zarb-go/util"
)

var _ Consensus = &MockConsensus{}

type MockConsensus struct {
	Lock      sync.RWMutex
	Votes     []*vote.Vote
	Proposal  *proposal.Proposal
	Scheduled bool
	State     *state.MockState
	Round     int16
}

func MockingConsensus(state *state.MockState) *MockConsensus {
	return &MockConsensus{State: state}
}
func (m *MockConsensus) MoveToNewHeight() {
	//m.Lock.Lock()
	//defer m.Lock.Unlock()
	m.Scheduled = true
}
func (m *MockConsensus) Start() error {
	return nil
}
func (m *MockConsensus) Stop() {}

func (m *MockConsensus) AddVote(v *vote.Vote) {
	//m.Lock.Lock()
	//defer m.Lock.Unlock()

	m.Votes = append(m.Votes, v)
}
func (m *MockConsensus) AllVotes() []*vote.Vote {
	//m.Lock.RLock()
	//defer m.Lock.RUnlock()
	return m.Votes
}
func (m *MockConsensus) RoundVotes(round int16) []*vote.Vote {
	//m.Lock.RLock()
	//defer m.Lock.RUnlock()

	votes := make([]*vote.Vote, 0)
	for _, v := range m.Votes {
		if v.Round() == round {
			votes = append(votes, v)
		}
	}
	return votes
}
func (m *MockConsensus) SetProposal(p *proposal.Proposal) {
	//m.Lock.RLock()
	//defer m.Lock.RUnlock()

	m.Proposal = p
}
func (m *MockConsensus) RoundProposal(round int16) *proposal.Proposal {
	//m.Lock.RLock()
	//defer m.Lock.RUnlock()

	if m.Proposal == nil || m.Proposal.Round() != round {
		return nil
	}
	return m.Proposal
}
func (m *MockConsensus) QueryProposal(round int16) *proposal.Proposal {
	//m.Lock.RLock()
	//defer m.Lock.RUnlock()

	if m.Proposal == nil || m.Proposal.Round() != round {
		return nil
	}
	return m.Proposal
}
func (m *MockConsensus) HeightRound() (int32, int16) {
	//m.Lock.RLock()
	//defer m.Lock.RUnlock()

	return m.State.LastBlockHeight() + 1, m.Round
}
func (m *MockConsensus) Fingerprint() string {
	return ""
}
func (m *MockConsensus) PickRandomVote() *vote.Vote {
	//m.Lock.RLock()
	//defer m.Lock.RUnlock()

	if len(m.Votes) == 0 {
		return nil
	}
	r := util.RandInt32(int32(len(m.Votes)))
	return m.Votes[r]
}
