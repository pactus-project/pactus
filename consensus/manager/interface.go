package manager

import (
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type ManagerReader interface {
	Instances() []consensus.Reader
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
