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
	s.checkForWeakValidity()
	s.detectDoubleProposal()

	cpPreVotes := s.log.CPPreVoteVoteSet(s.round)
	if cpPreVotes.Has2FP1Votes(s.cpRound) {
		if cpPreVotes.Has2FP1VotesFor(s.cpRound, vote.CPValueNo) {
			s.logger.Debug("cp: quorum for pre-votes", "value", "no")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueNo)
			s.cpDecidedCert = s.makeCertificate(votes)
			s.enterNewState(s.precommitState)
		} else if cpPreVotes.Has2FP1VotesFor(s.cpRound, vote.CPValueYes) {
			s.logger.Debug("cp: quorum for pre-votes", "value", "yes")

			votes := cpPreVotes.BinaryVotes(s.cpRound, vote.CPValueYes)
			cert := s.makeCertificate(votes)
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

			s.signAddCPMainVote(s.cpWeakValidity, s.cpRound, vote.CPValueAbstain, just)
			s.enterNewState(s.cpDecideState)
		}
	}

	s.cpStrongTermination()
	s.absoluteCommit()
}

func (s *cpMainVoteState) checkForWeakValidity() {
	if s.cpWeakValidity != hash.UndefHash {
		return
	}

	preVotes := s.log.CPPreVoteVoteSet(s.round)
	randVote := preVotes.GetRandomVote(s.cpRound, vote.CPValueNo)
	if randVote != nil {
		s.cpWeakValidity = randVote.BlockHash()
	}
}

func (s *cpMainVoteState) detectDoubleProposal() {
	if s.cpWeakValidity == hash.UndefHash {
		return
	}

	preVotes := s.log.CPPreVoteVoteSet(s.round)
	votesForNo := preVotes.BinaryVotes(0, vote.CPValueNo)
	for _, vte := range votesForNo {
		if vte.BlockHash() != s.cpWeakValidity {
			s.logger.Warn("double proposal detected",
				"proposal_1", s.cpWeakValidity,
				"proposal_2", vte.BlockHash())

			s.log.SetRoundProposal(s.round, nil)
			s.cpWeakValidity = hash.UndefHash

			return
		}
	}
}

func (s *cpMainVoteState) onAddVote(_ *vote.Vote) {
	s.decide()
}

func (*cpMainVoteState) name() string {
	return "cp:main-vote"
}
