package consensus

import (
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
)

type newHeightState struct {
	*consensus
}

func (s *newHeightState) enter() {
	sleep := s.state.LastBlockTime().Add(s.state.BlockTime()).Sub(util.Now())
	s.scheduleTimeout(sleep, s.height, s.round, tickerTargetNewHeight)
	s.logger.Debug("NewHeight is scheduled", "timeout", sleep.Seconds())
}

func (s *newHeightState) execute() {
	sateHeight := s.state.LastBlockHeight()
	if s.height == sateHeight+1 {
		s.logger.Trace("Duplicated entry")
		return
	}

	// Apply last certificate. We may have more votes now
	if s.height == sateHeight && s.round >= 0 {
		vs := s.pendingVotes.PrecommitVoteSet(s.round)
		if vs == nil {
			s.logger.Warn("Entering new height without last commit")
		} else {
			// Update last commit here, consensus had enough time to populate more votes
			lastCert := vs.ToCertificate()
			if lastCert != nil {
				if err := s.state.UpdateLastCertificate(lastCert); err != nil {
					s.logger.Warn("Updating last commit failed", "err", err)
				}
			}
		}
	}

	vals := s.state.CommitteeValidators()
	s.pendingVotes.MoveToNewHeight(sateHeight+1, vals)

	s.height = sateHeight + 1
	s.round = -1
	s.logger.Info("Entering new height", "height", s.height)

	s.enterNewState(s.newRoundState)
}

func (s *newHeightState) timedout(t *ticker) {
	if t.Target != tickerTargetNewHeight {
		s.logger.Debug("Invalid ticker", "ticker", t)
		return
	}
	s.execute()
}

func (s *newHeightState) voteAdded(v *vote.Vote) {
}

func (s *newHeightState) name() string {
	return newHeightName
}
