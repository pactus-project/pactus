package consensusv2

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type consState interface {
	enter()
	decide()
	onAddVote(v *vote.Vote)
	onSetProposal(p *proposal.Proposal)
	onTimeout(t *ticker)
	name() string
}
