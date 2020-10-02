package consensus

import (
	"gitlab.com/zarb-chain/zarb-go/crypto"
	"gitlab.com/zarb-chain/zarb-go/errors"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/validator"
	"gitlab.com/zarb-chain/zarb-go/vote"
)

type RoundVoteSet struct {
	Prevotes   *vote.VoteSet
	Precommits *vote.VoteSet
	proposal   *vote.Proposal
}

type HeightVoteSet struct {
	height         int
	valSet         *validator.ValidatorSet
	roundVoteSets  map[int]*RoundVoteSet
	votes          map[crypto.Hash]*vote.Vote
	lockedProposal *vote.Proposal
}

func NewHeightVoteSet(height int, valSet *validator.ValidatorSet) *HeightVoteSet {
	hvs := &HeightVoteSet{
		height:        height,
		valSet:        valSet,
		roundVoteSets: make(map[int]*RoundVoteSet),
		votes:         make(map[crypto.Hash]*vote.Vote),
	}
	hvs.addRound(0)
	return hvs
}

func (hvs *HeightVoteSet) addRound(round int) *RoundVoteSet {
	if rvs, ok := hvs.roundVoteSets[round]; ok {
		logger.Error("addRound() for an existing round")
		return rvs
	}
	prevotes := vote.NewVoteSet(hvs.height, round, vote.VoteTypePrevote, hvs.valSet)
	precommits := vote.NewVoteSet(hvs.height, round, vote.VoteTypePrecommit, hvs.valSet)
	rvs := &RoundVoteSet{
		Prevotes:   prevotes,
		Precommits: precommits,
	}

	hvs.roundVoteSets[round] = rvs
	return rvs
}

func (hvs *HeightVoteSet) AddVote(vote *vote.Vote) (bool, error) {
	if err := vote.SanityCheck(); err != nil {
		return false, errors.Errorf(errors.ErrInvalidVote, "%v", err)
	}
	voteSet := hvs.voteSet(vote.Round(), vote.Type())
	if voteSet == nil {
		hvs.addRound(vote.Round())
		voteSet = hvs.voteSet(vote.Round(), vote.Type())
	}
	added, err := voteSet.AddVote(vote)
	if added {
		hvs.votes[vote.Hash()] = vote

	}
	return added, err
}

func (hvs *HeightVoteSet) Prevotes(round int) *vote.VoteSet {
	return hvs.voteSet(round, vote.VoteTypePrevote)
}

func (hvs *HeightVoteSet) Precommits(round int) *vote.VoteSet {
	return hvs.voteSet(round, vote.VoteTypePrecommit)
}

func (hvs *HeightVoteSet) RoundProposal(round int) *vote.Proposal {
	rvs, ok := hvs.roundVoteSets[round]
	if !ok {
		return nil
	}
	return rvs.proposal
}

func (hvs *HeightVoteSet) SetRoundProposal(round int, proposal *vote.Proposal) {
	rvs, ok := hvs.roundVoteSets[round]
	if !ok {
		rvs = hvs.addRound(round)
	}
	rvs.proposal = proposal
}

func (hvs *HeightVoteSet) voteSet(round int, voteType vote.VoteType) *vote.VoteSet {
	rvs, ok := hvs.roundVoteSets[round]
	if !ok {
		return nil
	}
	switch voteType {
	case vote.VoteTypePrevote:
		return rvs.Prevotes
	case vote.VoteTypePrecommit:
		return rvs.Precommits
	}

	logger.Panic("Unexpected vote type %d", voteType)
	return nil
}
