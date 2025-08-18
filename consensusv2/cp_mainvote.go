package consensusv2

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

type cpMainVoteState struct {
	*changeProposer
}

func (s *cpMainVoteState) enter() {
	s.decide()
}

func (s *cpMainVoteState) decide() {
	s.strongCommit()
	s.cpStrongTermination()
	s.checkForWeakValidity()
	s.detectByzantineProposal()

	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
	if cpPreVotes.HasTwoFPlusOneVotes(s.cpRound) {
		if cpPreVotes.HasTwoFPlusOneVotesFor(s.cpRound, vote.CPValueNo) {
			s.logger.Debug("cp: quorum for pre-votes", "value", "no")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueNo)
			s.cpDecidedCert = s.makeVoteCertificate(votes)
			s.enterNewState(s.precommitState)
		} else if cpPreVotes.HasTwoFPlusOneVotesFor(s.cpRound, vote.CPValueYes) {
			s.logger.Debug("cp: quorum for pre-votes", "value", "yes")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeVoteCertificate(votes)
			just := &vote.JustMainVoteNoConflict{
				QCert: cert,
			}
			s.signAddCPMainVote(hash.UndefHash, s.cpRound, vote.CPValueYes, just)
			s.enterNewState(s.cpDecideState)
		} else {
			s.logger.Debug("cp: no-quorum for pre-votes", "value", "abstain")

			vote0 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueNo)
			vote1 := cpPreVotes.GetRandomVote(s.cpRound, vote.CPValueYes)

			just := &vote.JustMainVoteConflict{
				JustNo:  vote0.CPJust(),
				JustYes: vote1.CPJust(),
			}

			s.signAddCPMainVote(*s.cpWeakValidity, s.cpRound, vote.CPValueAbstain, just)
			s.enterNewState(s.cpDecideState)
		}
	}
}

func (s *cpMainVoteState) checkForWeakValidity() {
	if s.cpWeakValidity != nil {
		return
	}

	preVotes := s.log.CPPreVoteVoteSet(s.round)
	randVote := preVotes.GetRandomVote(s.cpRound, vote.CPValueNo)
	if randVote != nil {
		bh := randVote.BlockHash()
		s.cpWeakValidity = &bh
	}
}

func (s *cpMainVoteState) detectByzantineProposal() {
	if s.cpWeakValidity == nil {
		return
	}

	roundProposal := s.log.RoundProposal(s.round)

	if roundProposal != nil &&
		roundProposal.Block().Hash() != *s.cpWeakValidity {
		s.logger.Warn("double proposal detected",
			"proposal_1", s.cpWeakValidity,
			"proposal_2", roundProposal.Block().Hash())

		s.log.SetRoundProposal(s.round, nil)
		s.queryProposal()
	}
}

func (s *cpMainVoteState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (*cpMainVoteState) name() string {
	return "cp:main-vote"
}
