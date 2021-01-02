package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

type ConsensusReader interface {
	RoundVotes(round int) []*vote.Vote
	RoundVotesHash(round int) []crypto.Hash
	LastProposal() *vote.Proposal
	HRS() hrs.HRS
	Fingerprint() string
}

type Consensus interface {
	ConsensusReader

	MoveToNewHeight()
	Stop()
	AddVote(v *vote.Vote)
	SetProposal(proposal *vote.Proposal)
}
