package manager

import (
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type ConsensusManagerReader interface {
	Instances() []consensus.ConsensusReader
	HandleQueryVote(height types.Height, round types.Round) *vote.Vote
	HandleQueryProposal(height types.Height, round types.Round) *proposal.Proposal
	Proposal() *proposal.Proposal
	HeightRound() (types.Height, types.Round)
	HasActiveInstance() bool
}

type ConsensusManager interface {
	ConsensusManagerReader

	MoveToNewHeight()
	AddVote(vote *vote.Vote)
	SetProposal(prop *proposal.Proposal)
	IsDeprecated() bool
}
