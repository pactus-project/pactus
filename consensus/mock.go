package consensus

import (
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
}

func MockingConsensus(state *state.MockState) *MockConsensus {
	return &MockConsensus{State: state}
}

func (m *MockConsensus) MoveToNewHeight() {
	m.Scheduled = true
}
func (m *MockConsensus) Stop() {}

func (m *MockConsensus) AddVote(v *vote.Vote) {
	m.Votes = append(m.Votes, v)
}
func (m *MockConsensus) RoundVotes(round int) []*vote.Vote {
	votes := make([]*vote.Vote, 0)
	for _, v := range m.Votes {
		if v.Round() == round {
			votes = append(votes, v)
		}
	}
	return votes
}
func (m *MockConsensus) SetProposal(p *proposal.Proposal) {
	m.Proposal = p
}
func (m *MockConsensus) RoundProposal(round int) *proposal.Proposal {
	if m.Proposal == nil || m.Proposal.Round() != round {
		return nil
	}
	return m.Proposal
}
func (m *MockConsensus) HRS() hrs.HRS {
	return hrs.NewHRS(m.State.LastBlockHeight()+1, m.Round, hrs.StepTypeNewHeight)
}
func (m *MockConsensus) Fingerprint() string {
	return ""
}
func (m *MockConsensus) PickRandomVote(round int) *vote.Vote {
	if len(m.Votes) == 0 {
		return nil
	}
	r := util.RandInt(len(m.Votes))
	return m.Votes[r]
}
