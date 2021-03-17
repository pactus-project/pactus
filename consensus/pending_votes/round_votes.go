package pending_votes

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type RoundVotes struct {
	Prepares   *vote.VoteSet
	Precommits *vote.VoteSet
	proposal   *proposal.Proposal
}

func (rv *RoundVotes) addVote(v *vote.Vote) (bool, error) {
	vs := rv.voteSet(v.VoteType())
	return vs.AddVote(v)
}

func (rv *RoundVotes) HasVote(hash crypto.Hash) bool {
	votes := rv.AllVotes()
	for _, v := range votes {
		if v.Hash().EqualsTo(hash) {
			return true
		}
	}
	return false
}

func (rv *RoundVotes) AllVotes() []*vote.Vote {
	votes := []*vote.Vote{}
	votes = append(votes, rv.Prepares.AllVotes()...)
	votes = append(votes, rv.Precommits.AllVotes()...)

	return votes
}

func (rv *RoundVotes) voteSet(voteType vote.VoteType) *vote.VoteSet {
	switch voteType {
	case vote.VoteTypePrepare:
		return rv.Prepares
	case vote.VoteTypePrecommit:
		return rv.Precommits
	}

	logger.Panic("Unexpected vote type %d", voteType)
	return nil
}
