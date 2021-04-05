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

type initState struct {
}

func (s *initState) enter()                             {}
func (s *initState) onSetProposal(p *proposal.Proposal) {}
func (s *initState) onTimedout(t *ticker)               {}
func (s *initState) onAddVote(v *vote.Vote)             {}
func (s *initState) name() string {
	return "initializing"
}
