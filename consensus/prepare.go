package consensus

type prepareState struct {
	*consensus
	hasVoted bool
}

func (s *prepareState) enter() {
	s.hasVoted = false

	changeProperTimeout := s.config.CalculateChangeProposerTimeout(s.round)
	queryProposalTimeout := changeProperTimeout / 2
	s.scheduleTimeout(queryProposalTimeout, s.height, s.round, tickerTargetQueryProposal)
	s.scheduleTimeout(changeProperTimeout, s.height, s.round, tickerTargetChangeProposer)

	s.decide()
}

func (s *prepareState) decide() {
	s.vote()

	prepares := s.log.PrepareVoteSet(s.round)
	prepareQH := prepares.QuorumHash()
	if prepareQH != nil {
		s.logger.Debug("prepare has quorum", "hash", prepareQH.ShortString())
		s.enterNewState(s.precommitState)
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

func (s *prepareState) vote() {
	if s.hasVoted {
		return
	}

	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		s.logger.Debug("no proposal yet")
		return
	}

	// Everything is good
	s.signAddPrepareVote(roundProposal.Block().Hash())
	s.hasVoted = true
}

func (s *prepareState) timeout(t *ticker) {
	s.logger.Debug("timer expired", "ticker", t)

	if t.Target == tickerTargetQueryProposal {
		s.queryProposal()
		s.decide()
	} else if t.Target == tickerTargetChangeProposer {
		s.startChangingProposer()
	}
}

func (s *prepareState) name() string {
	return "prepare"
}
