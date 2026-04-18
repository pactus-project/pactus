package manager

import (
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type Reader interface {
	ConsensusKey() *bls.PublicKey
	AllVotes() []*vote.Vote
	HandleQueryVote(height types.Height, round types.Round) *vote.Vote
	HandleQueryProposal(height types.Height, round types.Round) *proposal.Proposal
	Proposal() *proposal.Proposal
	HasVote(h hash.Hash) bool
	HeightRound() (types.Height, types.Round)
	IsActive() bool
	IsProposer() bool
}

type Consensus interface {
	Reader

	MoveToNewHeight()
	AddVote(vote *vote.Vote)
	SetProposal(prop *proposal.Proposal)
	IsDeprecated() bool
}

type ManagerReader interface {
	Instances() []Reader
	HandleQueryVote(height types.Height, round types.Round) *vote.Vote
	HandleQueryProposal(height types.Height, round types.Round) *proposal.Proposal
	Proposal() *proposal.Proposal
	HeightRound() (types.Height, types.Round)
	HasActiveInstance() bool
}

type Manager interface {
	ManagerReader

	MoveToNewHeight()
	AddVote(vote *vote.Vote)
	SetProposal(prop *proposal.Proposal)
	IsDeprecated() bool
}
