package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type ConsensusReader interface {
	PickRandomVote() *vote.Vote
	RoundVotes(round int) []*vote.Vote
	RoundProposal(round int) *proposal.Proposal
	HRS() hrs.HRS
	Fingerprint() string
}

type Consensus interface {
	ConsensusReader

	MoveToNewHeight()
	Stop()
	AddVote(v *vote.Vote)
	SetProposal(proposal *proposal.Proposal)
}
