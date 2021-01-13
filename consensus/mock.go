package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

var _ Consensus = &MockConsensus{}

type MockConsensus struct {
	Votes    []*vote.Vote
	Proposal *vote.Proposal
	HRS_     hrs.HRS
	Started  bool
}

func MockingConsensus() *MockConsensus {
	return &MockConsensus{}
}

func (m *MockConsensus) MoveToNewHeight() {
	m.Started = true
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
func (m *MockConsensus) RoundVotesHash(round int) []crypto.Hash {
	hashes := make([]crypto.Hash, len(m.Votes))
	for i, v := range m.RoundVotes(round) {
		hashes[i] = v.Hash()
	}
	return hashes
}
func (m *MockConsensus) SetProposal(p *vote.Proposal) {
	m.Proposal = p
}
func (m *MockConsensus) LastProposal() *vote.Proposal {
	return m.Proposal
}
func (m *MockConsensus) HRS() hrs.HRS {
	return m.HRS_
}
func (m *MockConsensus) Fingerprint() string {
	return ""
}
