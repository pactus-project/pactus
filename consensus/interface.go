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
	HandleQueryVote(height uint32, round int16) *vote.Vote
	HandleQueryProposal(height uint32, round int16) *proposal.Proposal
	Proposal() *proposal.Proposal
	HasVote(h hash.Hash) bool
	HeightRound() (uint32, int16)
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
