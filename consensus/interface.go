package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

type Consensus interface {
	MoveToNewHeight()
	Stop()
	AddVote(v *vote.Vote)
	RoundVotes(round int) []*vote.Vote
	RoundVotesHash(round int) []crypto.Hash
	SetProposal(proposal *vote.Proposal)
	LastProposal() *vote.Proposal
	HRS() hrs.HRS
	Fingerprint() string
}
