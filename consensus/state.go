package consensus

import (
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type consState interface {
	enter()
	onSetProposal(p *proposal.Proposal)
	onAddVote(v *vote.Vote)
	onTimedout(t *ticker)
	name() string
}
