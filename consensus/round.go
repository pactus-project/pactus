package consensus

import (
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type newRoundState struct {
	*consensus
}

func (s *newRoundState) enter() {
	sleep := s.config.ChangeProposerTimeout
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
	s.logger.Debug("Change proposer timer started...", "timeout", sleep.Seconds())

	s.decide()
}

func (s *newRoundState) decide() {
	// make sure we have quorum votes for previous round
	if s.round > 0 {
		prepares := s.pendingVotes.PrepareVoteSet(s.round - 1)
		precommits := s.pendingVotes.PrecommitVoteSet(s.round - 1)
		// Normally when there is no proposal for this round, every one should vote for nil
		prepareQH := prepares.QuorumHash()
		precommitQH := precommits.QuorumHash()
		if prepareQH == nil || !prepareQH.IsUndef() {
			s.logger.Warn("Suspicious prepares", "prepareQH", prepareQH)
		}
		if precommitQH == nil || !precommitQH.IsUndef() {
			s.logger.Warn("Suspicious precommits", "precommitQH", precommitQH)
		}
	}

	s.logger.Info("Entering new round", "round", s.round)
	s.enterNewState(s.proposeState)
}

func (s *newRoundState) onAddVote(v *vote.Vote) {
	panic("Unreachable")
}

func (s *newRoundState) onSetProposal(p *proposal.Proposal) {
	panic("Unreachable")
}

func (s *newRoundState) onTimedout(t *ticker) {
	panic("Unreachable")
}

func (s *newRoundState) name() string {
	return "new-round"
}
