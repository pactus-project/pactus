package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

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
		s.logger.Debug("pre-commit has quorum", "hash", precommitQH)

		roundProposal := s.log.RoundProposal(s.round)
		if roundProposal == nil {
			// There is a consensus about a proposal that we don't have yet.
			// Ask peers for this proposal.
			s.logger.Info("query for a decided proposal", "hash", precommitQH)
			s.queryProposal()
		} else if s.hasVoted {
			// To ensure we have voted and won't be absent from the certificate
			s.enterNewState(s.commitState)
		}
	} else {
		//
		// If a validator receives a set of f+1 valid cp:PRE-VOTE votes for this round,
		// it starts the changing proposer phase, even if its timer has not expired;
		// This prevents it from starting the change-proposer phase too late.
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
		s.logger.Warn("double proposal detected", "roundProposal", roundProposal, "prepared", *prepareQH)

		return
	}

	// Everything is good
	s.signAddPrecommitVote(*prepareQH)
	s.hasVoted = true
}

func (s *precommitState) onAddVote(v *vote.Vote) {
	if v.Type() == vote.VoteTypePrecommit ||
		v.Type() == vote.VoteTypeCPPreVote {
		s.decide()
	}
}

func (s *precommitState) onSetProposal(_ *proposal.Proposal) {
	s.decide()
}

func (s *precommitState) onTimeout(t *ticker) {
	if t.Target == tickerTargetChangeProposer {
		s.startChangingProposer()
	}
}

func (*precommitState) name() string {
	return "precommit"
}
