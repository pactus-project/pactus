package consensus

import "github.com/zarbchain/zarb-go/vote"

type newRoundState struct {
	*consensus
}

func (s *newRoundState) enter() {
	s.execute()
}

func (s *newRoundState) execute() {
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

func (s *newRoundState) timedout(t *ticker) {
}

func (s *newRoundState) voteAdded(v *vote.Vote) {
}

func (s *newRoundState) name() string {
	return newRoundName
}
