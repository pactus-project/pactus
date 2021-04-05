package pending_votes

import (
	"github.com/zarbchain/zarb-go/consensus/vote_set"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

type RoundVotes struct {
	prepareVotes        *vote_set.VoteSet
	precommitVotes      *vote_set.VoteSet
	changeProposerVotes *vote_set.VoteSet
	proposal            *proposal.Proposal
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
	votes = append(votes, rv.prepareVotes.AllVotes()...)
	votes = append(votes, rv.precommitVotes.AllVotes()...)
	votes = append(votes, rv.changeProposerVotes.AllVotes()...)

	return votes
}

func (rv *RoundVotes) voteSet(voteType vote.VoteType) *vote_set.VoteSet {
	switch voteType {
	case vote.VoteTypePrepare:
		return rv.prepareVotes
	case vote.VoteTypePrecommit:
		return rv.precommitVotes
	case vote.VoteTypeChangeProposer:
		return rv.changeProposerVotes
	}

	logger.Panic("Unexpected vote type %d", voteType)
	return nil
}
