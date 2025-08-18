package consensusv2

// type prepareState struct {
// 	*fastConsensus
// 	hasVoted bool
// }

// func (s *prepareState) enter() {
// 	s.hasVoted = false

// 	changeProperTimeout := s.config.CalculateChangeProposerTimeout(s.round)
// 	queryProposalTimeout := changeProperTimeout / 2
// 	s.scheduleTimeout(queryProposalTimeout, s.height, s.round, tickerTargetQueryProposal)
// 	s.scheduleTimeout(changeProperTimeout, s.height, s.round, tickerTargetChangeProposer)

// 	s.decide()
// }

// func (s *prepareState) decide() {
// 	s.vote()
// 	s.strongCommit()

// 	//
// 	// If a validator receives a set of f+1 valid cp:PRE-VOTE votes for this round,
// 	// it starts changing the proposer phase, even if its timer has not expired;
// 	// This prevents it from starting the change-proposer phase too late.
// 	//
// 	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
// 	if cpPreVotes.HasFPlusOneVotesFor(0, vote.CPValueYes) {
// 		s.startChangingProposer()
// 	}
// }

// func (s *prepareState) vote() {
// 	if s.hasVoted {
// 		return
// 	}

// 	roundProposal := s.log.RoundProposal(s.round)
// 	if roundProposal == nil {
// 		s.logger.Debug("no proposal yet")

// 		return
// 	}

// 	// Everything is good
// 	s.signAddPrepareVote(roundProposal.Block().Hash())
// 	s.hasVoted = true
// }

// func (s *prepareState) onTimeout(t *ticker) {
// 	if t.Target == tickerTargetQueryProposal {
// 		roundProposal := s.log.RoundProposal(s.round)
// 		if roundProposal == nil {
// 			s.queryProposal()
// 		}
// 		if s.isProposer() {
// 			s.queryVotes()
// 		}
// 	} else if t.Target == tickerTargetChangeProposer {
// 		s.startChangingProposer()
// 	}
// }

// func (s *prepareState) onAddVote(v *vote.Vote) {
// 	if v.Type() == vote.VoteTypePrepare ||
// 		v.Type() == vote.VoteTypeCPPreVote {
// 		s.decide()
// 	}
// }

// func (s *prepareState) onSetProposal(_ *proposal.Proposal) {
// 	s.decide()
// }

// func (*prepareState) name() string {
// 	return "prepare"
// }
