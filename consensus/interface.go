package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
)

type ConsensusReader interface {
	PickRandomVote() *vote.Vote
	AllVotes() []*vote.Vote
	RoundVotes(round int) []*vote.Vote
	RoundProposal(round int) *proposal.Proposal
	HeightRound() (int, int)
	Fingerprint() string
}

type Consensus interface {
	ConsensusReader

	MoveToNewHeight()
	Start() error
	Stop()
	AddVote(v *vote.Vote)
	SetProposal(proposal *proposal.Proposal)
}
