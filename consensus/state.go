package consensus

import (
	"github.com/zarbchain/zarb-go/types/proposal"
	"github.com/zarbchain/zarb-go/types/vote"
)

type consState interface {
	enter()
	onSetProposal(p *proposal.Proposal)
	onAddVote(v *vote.Vote)
	onTimeout(t *ticker)
	name() string
}
