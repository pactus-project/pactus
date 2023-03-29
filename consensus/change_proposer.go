package consensus

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type changeProposerState struct {
	*consensus
}

func (s *changeProposerState) enter() {
	s.logger.Info("requesting for changing proposer", "proposer", s.proposer(s.round).Address())
	s.signAddVote(vote.VoteTypeChangeProposer, hash.UndefHash)
	s.log.SetRoundProposal(s.round, nil)

	s.decide()
}

func (s *changeProposerState) decide() {
	voteset := s.log.ChangeProposerVoteSet(s.round)
	if voteset.QuorumHash() != nil {
		s.logger.Debug("change proposer has quorum", "proposer", s.proposer(s.round).Address())
		s.round++

		s.enterNewState(s.proposeState)
	}
}

func (s *changeProposerState) onAddVote(v *vote.Vote) {
	// Only accept change proposer votes
	if v.Type() == vote.VoteTypeChangeProposer {
		s.doAddVote(v)
		s.decide()
	}
}

func (s *changeProposerState) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposals
}

func (s *changeProposerState) onTimeout(_ *ticker) {
	// Ignore timeouts
}

func (s *changeProposerState) name() string {
	return "change-proposer"
}
