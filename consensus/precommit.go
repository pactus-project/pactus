package consensus

type precommitState struct {
	*consensus
	hasVoted bool
}

func (s *precommitState) enter() {
	s.hasVoted = false

	s.decide()
}

func (s *precommitState) decide() {
	s.vote()

	precommits := s.log.PrecommitVoteSet(s.round)
	precommitQH := precommits.QuorumHash()
	if precommitQH != nil {
		s.logger.Debug("pre-commit has quorum", "hash", precommitQH.ShortString())

		roundProposal := s.log.RoundProposal(s.round)
		if roundProposal == nil {
			// There is a consensus about a proposal that we don't have yet.
			// Ask peers for this proposal.
			s.logger.Info("query for a decided proposal", "hash", precommitQH.ShortString())
			s.queryProposal()
		} else {
			// To ensure we have voted and won't be absent from the certificate
			if s.hasVoted {
				s.enterNewState(s.commitState)
			}
		}
	} else {
		//
		// If a replica receives a set of f+1 valid change-proposer votes for this round,
		// it goes to the change-proposer state, even if its timer has not expired;
		// This prevents it from starting the change-proposer state too late.
		//
		cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
		if cpPreVotes.HasOneThirdOfTotalPower(0) {
			s.startChangingProposer()
		}
	}
}

func (s *precommitState) vote() {
	if s.hasVoted {
		return
	}

	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		s.queryProposal()
		s.logger.Debug("no proposal yet")
		return
	}

	prepares := s.log.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	if !roundProposal.IsForBlock(*prepareQH) {
		s.log.SetRoundProposal(s.round, nil)
		s.queryProposal()
		s.logger.Warn("We don't have the quorum proposal", "our proposal", roundProposal, "quorum hash", *prepareQH)
		return
	}

	// Everything is good
	s.signAddPrecommitVote(*prepareQH)
	s.hasVoted = true
}

func (s *precommitState) timeout(t *ticker) {
	s.logger.Debug("timer expired")

	if t.Target == tickerTargetChangeProposer {
		s.startChangingProposer()
	}
}

func (s *precommitState) name() string {
	return "precommit"
}
