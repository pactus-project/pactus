package consensusv2

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
)

type proposeState struct {
	*consensusV2
}

func (s *proposeState) enter() {
	s.decide()
}

func (s *proposeState) decide() {
	proposer := s.proposer(s.round)
	s.cpRound = 0
	s.cpDecidedCert = nil
	s.cpWeakValidity = hash.UndefHash

	// Based on PIP-19, if the Availability Score is less than the Minimum threshold,
	// we initiate the Change-Proposer phase.
	// TODO: write test for me
	score := s.bcState.AvailabilityScore(proposer.Number())
	if score < s.config.MinimumAvailabilityScore {
		s.logger.Info("availability score of proposer is low",
			"score", score, "proposer", proposer.Address())
		s.startChangingProposer()

		return
	}

	if proposer.Address() == s.valKey.Address() {
		s.logger.Info("our turn to propose", "proposer", proposer.Address())
		s.createProposal(s.height, s.round)
	} else {
		s.logger.Debug("not our turn to propose", "proposer", proposer.Address())
	}

	s.enterNewState(s.precommitState)
}

func (s *proposeState) createProposal(height uint32, round int16) {
	block, err := s.bcState.ProposeBlock(s.valKey, s.rewardAddr)
	if err != nil {
		s.logger.Error("unable to propose a block!", "error", err)

		return
	}

	prop := proposal.NewProposal(height, round, block)
	sig := s.valKey.Sign(prop.SignBytes())
	prop.SetSignature(sig)

	s.log.SetRoundProposal(round, prop)

	s.broadcastProposal(prop)

	s.logger.Info("proposal signed and broadcasted", "proposal", prop)
}

func (*proposeState) onAddVote(_ *vote.Vote) {
	panic("Unreachable")
}

func (*proposeState) onSetProposal(_ *proposal.Proposal) {
	panic("Unreachable")
}

func (*proposeState) onTimeout(_ *ticker) {
	panic("Unreachable")
}

func (*proposeState) name() string {
	return "propose"
}
