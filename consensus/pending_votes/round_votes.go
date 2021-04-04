package pending_votes

import (
	"github.com/zarbchain/zarb-go/consensus/vote_set"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type RoundVotes struct {
	prepares   *vote_set.VoteSet
	precommits *vote_set.VoteSet
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
	votes = append(votes, rv.prepares.AllVotes()...)
	votes = append(votes, rv.precommits.AllVotes()...)

	return votes
}

func (rv *RoundVotes) voteSet(voteType vote.VoteType) *vote_set.VoteSet {
	switch voteType {
	case vote.VoteTypePrepare:
		return rv.prepares
	case vote.VoteTypePrecommit:
		return rv.precommits
	}

	logger.Panic("Unexpected vote type %d", voteType)
	return nil
}
