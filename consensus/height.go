package consensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
)

type newHeightState struct {
	*consensus
}

func (s *newHeightState) enter() {
	sleep := s.state.LastBlockTime().Add(s.state.Params().BlockInterval()).Sub(util.Now())
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
}

func (s *newHeightState) decide() {
	sateHeight := s.state.LastBlockHeight()

	// Try to update the last certificate. We may have more votes now.
	if s.height == sateHeight {
		precommits := s.log.PrecommitVoteSet(s.round)
		if precommits != nil {
			roundProposal := s.log.RoundProposal(s.round)
			if roundProposal != nil {
				// The last certificate is updated at this point since consensus has
				// had sufficient time to populate additional votes.
				votes := precommits.BlockVotes(roundProposal.Block().Hash())
				lastCert := s.makeCertificate(votes)
				if lastCert != nil {
					if err := s.state.UpdateLastCertificate(lastCert); err != nil {
						s.logger.Warn("updating last certificate failed", "err", err)
					}
				}
			}
		}
	}

	validators := s.state.CommitteeValidators()
	s.log.MoveToNewHeight(validators)

	s.validators = validators
	s.height = sateHeight + 1
	s.round = 0
	s.active = s.state.IsInCommittee(s.signer.Address())
	s.logger.Info("entering new height", "height", s.height, "active", s.active)

	if s.active {
		s.enterNewState(s.proposeState)
	}
}

func (s *newHeightState) onAddVote(_ *vote.Vote) {
	// Ignore votes
}

func (s *newHeightState) onSetProposal(_ *proposal.Proposal) {
	// Ignore proposal
}

func (s *newHeightState) onTimeout(t *ticker) {
	if t.Target == tickerTargetNewHeight {
		s.decide()
	}
}

func (s *newHeightState) name() string {
	return "new-height"
}
