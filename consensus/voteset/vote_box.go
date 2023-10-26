package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/vote"
)

type voteBox struct {
	votes      map[crypto.Address]*vote.Vote
	votedPower int64
}

func newVoteBox() *voteBox {
	return &voteBox{
		votes:      make(map[crypto.Address]*vote.Vote),
		votedPower: 0,
	}
}

func (vs *voteBox) addVote(vote *vote.Vote, power int64) {
	if vs.votes[vote.Signer()] == nil {
		vs.votes[vote.Signer()] = vote
		vs.votedPower += power
	}
}
