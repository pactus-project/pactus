package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

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
		s.logger.Debug("prepare has quorum", "hash", prepareQH)
		s.enterNewState(s.precommitState)
	} else {
		//
		// If a validator receives a set of f+1 valid cp:PRE-VOTE votes for this round,
		// it starts changing the proposer phase, even if its timer has not expired;
		// This prevents it from starting the change-proposer phase too late.
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

func (s *prepareState) onTimeout(ticker *ticker) {
	switch ticker.Target {
	case tickerTargetQueryProposal:
		roundProposal := s.log.RoundProposal(s.round)
		if roundProposal == nil {
			s.queryProposal()
		}
		if s.isProposer() {
			s.queryVote()
		}

		// Schedule another timeout to retry querying for the proposal or votes.
		// This ensures that delayed or missing data doesn't cause the process to stall.
		s.scheduleTimeout(ticker.Duration*2, s.height, s.round, tickerTargetQueryProposal)

	case tickerTargetChangeProposer:
		s.startChangingProposer()

	case tickerTargetNewHeight, tickerTargetQueryVote:
		// These targets are not used in the prepare state
	}
}

func (s *prepareState) onAddVote(v *vote.Vote) {
	if v.Type() == vote.VoteTypePrepare ||
		v.Type() == vote.VoteTypeCPPreVote {
		s.decide()
	}
}

func (s *prepareState) onSetProposal(_ *proposal.Proposal) {
	s.decide()
}

func (*prepareState) name() string {
	return "prepare"
}
