package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type consState interface {
	enter()
	onSetProposal(p *proposal.Proposal)
	onAddVote(v *vote.Vote)
	onTimeout(t *ticker)
	name() string
}
