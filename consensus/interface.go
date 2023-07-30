package consensus

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type Reader interface {
	SignerKey() crypto.PublicKey
	AllVotes() []*vote.Vote
	PickRandomVote(round int16) *vote.Vote
	RoundProposal(round int16) *proposal.Proposal
	HasVote(hash hash.Hash) bool
	HeightRound() (uint32, int16)
	IsActive() bool
}

type Consensus interface {
	Reader

	MoveToNewHeight()
	AddVote(v *vote.Vote)
	SetProposal(proposal *proposal.Proposal)
}

type ManagerReader interface {
	Instances() []Reader
	PickRandomVote(round int16) *vote.Vote
	RoundProposal(round int16) *proposal.Proposal
	HeightRound() (uint32, int16)
	HasActiveInstance() bool
}

type Manager interface {
	ManagerReader

	Start() error
	Stop()
	MoveToNewHeight()
	AddVote(v *vote.Vote)
	SetProposal(proposal *proposal.Proposal)
}
