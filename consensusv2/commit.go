package consensusv2

import (
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type commitState struct {
	*consensusV2
}

func (s *commitState) enter() {
	s.decide()
}

func (s *commitState) decide() {
	roundProposal := s.log.RoundProposal(s.round)
	block := roundProposal.Block()
	precommits := s.log.PrecommitVoteSet(s.round)
	votes := precommits.BlockVotes(block.Hash())
	cert := s.makeCertificate(votes)

	err := s.bcState.CommitBlock(block, cert)
	if err != nil {
		s.logger.Error("committing block failed", "block", block, "error", err)
	} else {
		s.logger.Info("block committed, schedule new height", "hash", block.Hash())

		// Now we can announce the committed block and certificate
		s.announceNewBlock(block, cert, s.cpDecidedCert)
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
