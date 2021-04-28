package consensus

import (
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/util"
)

type newHeightState struct {
	*consensus
}

func (s *newHeightState) enter() {
	sleep := s.state.LastBlockTime().Add(s.state.BlockTime()).Sub(util.Now())
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
}

func (s *newHeightState) decide() {
	sateHeight := s.state.LastBlockHeight()
	if s.height == sateHeight+1 {
		s.logger.Warn("Duplicated entry")
		return
	}

	// Apply last certificate. We may have more votes now
	if s.height == sateHeight && s.round >= 0 {
		vs := s.pendingVotes.PrecommitVoteSet(s.round)
		if vs == nil {
			s.logger.Warn("Entering new height without certificate")
		} else {
			// Update last certificate here, consensus had enough time to populate more votes
			lastCert := vs.ToCertificate()
			if lastCert != nil {
				if err := s.state.UpdateLastCertificate(lastCert); err != nil {
					s.logger.Warn("Updating last certificate failed", "err", err)
				}
			}
		}
	}

	vals := s.state.CommitteeValidators()
	s.pendingVotes.MoveToNewHeight(sateHeight+1, vals)

	s.height = sateHeight + 1
	s.round = 0
	s.logger.Info("Entering new height", "height", s.height)

	s.enterNewState(s.proposeState)
}

func (s *newHeightState) onAddVote(v *vote.Vote) {
	s.doAddVote(v)
}

func (s *newHeightState) onSetProposal(p *proposal.Proposal) {
}

func (s *newHeightState) onTimedout(t *ticker) {
	if t.Target != tickerTargetNewHeight {
		s.logger.Debug("Invalid ticker", "ticker", t)
		return
	}
	s.decide()
}

func (s *newHeightState) name() string {
	return "new-height"
}
