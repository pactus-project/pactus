package consensus

import (
	"github.com/zarbchain/zarb-go/types/proposal"
	"github.com/zarbchain/zarb-go/types/vote"
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
		s.logger.Warn("duplicated entry")
		return
	}

	// Apply last certificate. We may have more votes now
	if s.height == sateHeight && s.round >= 0 {
		vs := s.log.PrecommitVoteSet(s.round)
		if vs == nil {
			s.logger.Warn("entering new height without certificate")
		} else {
			// Update last certificate here, consensus had enough time to populate more votes
			lastCert := vs.ToCertificate()
			if lastCert != nil {
				if err := s.state.UpdateLastCertificate(lastCert); err != nil {
					s.logger.Warn("updating last certificate failed", "err", err)
				}
			}
		}
	}

	vals := s.state.CommitteeValidators()
	s.log.MoveToNewHeight(vals)

	s.height = sateHeight + 1
	s.round = 0
	s.logger.Info("entering new height", "height", s.height)

	s.enterNewState(s.proposeState)
}

func (s *newHeightState) onAddVote(v *vote.Vote) {
	s.doAddVote(v)
}

func (s *newHeightState) onSetProposal(p *proposal.Proposal) {
}

func (s *newHeightState) onTimeout(t *ticker) {
	if t.Target != tickerTargetNewHeight {
		s.logger.Debug("invalid ticker", "ticker", t)
		return
	}
	s.decide()
}

func (s *newHeightState) name() string {
	return "new-height"
}
