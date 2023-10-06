package consensus

import (
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type Reader interface {
	ConsensusKey() *bls.PublicKey
	AllVotes() []*vote.Vote
	PickRandomVote(round int16) *vote.Vote
	Proposal() *proposal.Proposal
	HasVote(hash hash.Hash) bool
	HeightRound() (uint32, int16)
	IsActive() bool
}

type Consensus interface {
	Reader

	Start()
	MoveToNewHeight()
	AddVote(vote *vote.Vote)
	SetProposal(proposal *proposal.Proposal)
}

type ManagerReader interface {
	Instances() []Reader
	PickRandomVote(round int16) *vote.Vote
	Proposal() *proposal.Proposal
	HeightRound() (uint32, int16)
	HasActiveInstance() bool
}

type Manager interface {
	ManagerReader

	Start() error
	Stop()
	MoveToNewHeight()
	AddVote(vote *vote.Vote)
	SetProposal(proposal *proposal.Proposal)
}
