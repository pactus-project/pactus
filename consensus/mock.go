package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

var _ Consensus = &MockConsensus{}

type MockConsensus struct {
}

func NewMockConsensus() *MockConsensus {
	return &MockConsensus{}
}

func (m *MockConsensus) MoveToNewHeight()     {}
func (m *MockConsensus) AddVote(v *vote.Vote) {}
func (m *MockConsensus) AllVotes() []*vote.Vote {
	return nil
}
func (m *MockConsensus) AllVotesHashes() []crypto.Hash {
	return nil
}
func (m *MockConsensus) SetProposal(proposal *vote.Proposal) {

}
func (m *MockConsensus) LastProposal() *vote.Proposal {
	return nil
}
func (m *MockConsensus) HRS() hrs.HRS {
	return hrs.NewHRS(0, 0, 0)
}
func (m *MockConsensus) Fingerprint() string {
	return ""
}
