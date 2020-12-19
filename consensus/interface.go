package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

type Consensus interface {
	MoveToNewHeight()
	AddVote(v *vote.Vote)
	AllVotes() []*vote.Vote
	AllVotesHashes() []crypto.Hash
	SetProposal(proposal *vote.Proposal)
	LastProposal() *vote.Proposal
	HRS() hrs.HRS
	Fingerprint() string
}
