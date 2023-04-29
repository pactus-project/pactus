package consensus

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
)

var _ Consensus = &MockConsensus{}

type MockConsensus struct {
	Signer   crypto.Signer
	Votes    []*vote.Vote
	Proposal *proposal.Proposal
	Active   bool
	Height   uint32
	Round    int16
}

func MockingManager(signers []crypto.Signer) (Manager, []*MockConsensus) {
	mocks := make([]*MockConsensus, len(signers))
	instances := make([]Consensus, len(signers))
	for i, s := range signers {
		cons := MockingConsensus(s)
		mocks[i] = cons
		instances[i] = cons
	}

	return &manager{
		instances: instances,
	}, mocks
}

func MockingConsensus(signer crypto.Signer) *MockConsensus {
	return &MockConsensus{Signer: signer}
}
func (m *MockConsensus) SignerKey() crypto.PublicKey {
	return m.Signer.PublicKey()
}
func (m *MockConsensus) MoveToNewHeight() {
	m.Height++
}
func (m *MockConsensus) Start() error {
	return nil
}
func (m *MockConsensus) Stop() {}

func (m *MockConsensus) AddVote(v *vote.Vote) {
	m.Votes = append(m.Votes, v)
}
func (m *MockConsensus) AllVotes() []*vote.Vote {
	return m.Votes
}
func (m *MockConsensus) SetProposal(p *proposal.Proposal) {
	m.Proposal = p
}
func (m *MockConsensus) RoundProposal(round int16) *proposal.Proposal {
	if m.Proposal == nil || m.Proposal.Round() != round {
		return nil
	}
	return m.Proposal
}
func (m *MockConsensus) HeightRound() (uint32, int16) {
	return m.Height, m.Round
}
func (m *MockConsensus) Fingerprint() string {
	return ""
}
func (m *MockConsensus) PickRandomVote() *vote.Vote {
	if len(m.Votes) == 0 {
		return nil
	}
	r := util.RandInt32(int32(len(m.Votes)))
	return m.Votes[r]
}
func (m *MockConsensus) IsActive() bool {
	return m.Active
}
