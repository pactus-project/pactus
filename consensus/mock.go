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

func NewMockConsensus() *MockConsensus {
	return &MockConsensus{}
}

func (m *MockConsensus) MoveToNewHeight() {
	m.Started = true
}
func (m *MockConsensus) AddVote(v *vote.Vote) {
	m.Votes = append(m.Votes, v)
}
func (m *MockConsensus) RoundVotes(round int) []*vote.Vote {
	return m.Votes
}
func (m *MockConsensus) RoundVotesHash(round int) []crypto.Hash {
	hashes := make([]crypto.Hash, len(m.Votes))
	for i, v := range m.Votes {
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
