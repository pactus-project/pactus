package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
)

type changeProposerState struct {
	*consensus
}

func (s *changeProposerState) enter() {
	s.logger.Info("Requesting for chaning proposer")
	s.signAddVote(vote.VoteTypeChangeProposer, crypto.UndefHash)

	s.decide()
}

func (s *changeProposerState) decide() {
	votes := s.pendingVotes.ChangeProposerVoteSet(s.round)
	if votes.QuorumHash() != nil {
		s.logger.Debug("change proposer has quorum")
		s.round++

		s.enterNewState(s.newRoundState)
	}
}

func (s *changeProposerState) onAddVote(v *vote.Vote) {
	// Only accept change propser votes
	if v.VoteType() == vote.VoteTypeChangeProposer {
		s.doAddVote(v)
		s.decide()
	}
}

func (s *changeProposerState) onSetProposal(p *proposal.Proposal) {
	// Ignore proposals
}

func (s *changeProposerState) onTimedout(t *ticker) {
	// Ignore timeouts
}

func (s *changeProposerState) name() string {
	return "change-proposer"
}
