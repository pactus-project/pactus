package consensus

import (
	"github.com/zarbchain/zarb-go/types/proposal"
	"github.com/zarbchain/zarb-go/types/vote"
)

type Reader interface {
	PickRandomVote() *vote.Vote
	AllVotes() []*vote.Vote
	RoundVotes(round int16) []*vote.Vote
	RoundProposal(round int16) *proposal.Proposal
	HeightRound() (int32, int16)
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
