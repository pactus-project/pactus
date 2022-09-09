package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type Reader interface {
	PickRandomVote() *vote.Vote
	AllVotes() []*vote.Vote
	RoundVotes(round int16) []*vote.Vote
	RoundProposal(round int16) *proposal.Proposal
	HeightRound() (uint32, int16)
	Fingerprint() string
}

type Consensus interface {
	Reader

	MoveToNewHeight()
	Start() error
	Stop()
	AddVote(v *vote.Vote)
	SetProposal(proposal *proposal.Proposal)
}
