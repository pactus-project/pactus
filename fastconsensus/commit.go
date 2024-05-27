package fastconsensus

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type commitState struct {
	*consensus
}

func (s *commitState) enter() {
	s.decide()
}

func (s *commitState) decide() {
	roundProposal := s.log.RoundProposal(s.round)
	certBlock := roundProposal.Block()
	err := s.bcState.CommitBlock(certBlock, s.blockCert)
	if err != nil {
		s.logger.Error("committing block failed", "block", certBlock, "error", err)
	} else {
		s.logger.Info("block committed, schedule new height", "hash", certBlock.Hash())

		// Now we can announce the committed block and certificate
		s.announceNewBlock(certBlock, s.blockCert)
	}

	s.enterNewState(s.newHeightState)
}

func (*commitState) onAddVote(_ *vote.Vote) {
	panic("Unreachable")
}

func (*commitState) onSetProposal(_ *proposal.Proposal) {
	panic("Unreachable")
}

func (*commitState) onTimeout(_ *ticker) {
	panic("Unreachable")
}

func (*commitState) name() string {
	return "commit"
}
