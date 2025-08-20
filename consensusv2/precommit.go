package consensusv2

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type precommitState struct {
	*consensusV2
	hasVoted bool
}

func (s *precommitState) enter() {
	s.hasVoted = false

	changeProperTimeout := s.config.CalculateChangeProposerTimeout(s.round)
	queryProposalTimeout := changeProperTimeout / 2
	s.scheduleTimeout(queryProposalTimeout, s.height, s.round, tickerTargetQueryProposal)
	s.scheduleTimeout(changeProperTimeout, s.height, s.round, tickerTargetChangeProposer)

	s.decide()
}

func (s *precommitState) decide() {
	s.vote()
	s.absoluteCommit()

	//
	// The block can be committed by `2f+1` votes from the committee and
	// the proof of the change-proposer phase.
	//
	if s.cpDecidedCert != nil {
		roundProposal := s.log.RoundProposal(s.round)
		if roundProposal == nil {
			s.queryProposal()

			return
		}

		precommits := s.log.PrecommitVoteSet(s.round)
		precommitQH := precommits.QuorumHash()
		if precommitQH == nil {
			s.queryVote()

			return
		}

		s.enterNewState(s.commitState)
	}

	//
	// If a validator receives a set of `f+1` valid `cp:PRE-VOTE` votes for this round,
	// it starts changing the proposer phase, even if its timer has not expired;
	// This prevents it from starting the change-proposer phase too late.
	//
	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
	if cpPreVotes.HasFPlusOneVotesFor(0, vote.CPValueYes) {
		s.startChangingProposer()
	}
}

func (s *precommitState) vote() {
	if s.hasVoted {
		return
	}

	roundProposal := s.log.RoundProposal(s.round)
	if roundProposal == nil {
		s.logger.Debug("no proposal yet")

		return
	}

	// Everything is good
	s.signAddPrecommitVote(roundProposal.Block().Hash())
	s.hasVoted = true
}

func (s *precommitState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (s *precommitState) onSetProposal(_ *proposal.Proposal) {
	s.decide()
}

func (s *precommitState) onTimeout(ticker *ticker) {
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
		// Ignore it
	}
}

func (*precommitState) name() string {
	return "precommit"
}
