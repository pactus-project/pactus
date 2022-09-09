package voteset

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/logger"
)

type blockVotes struct {
	votes map[crypto.Address]*vote.Vote
	power int64
}

func newBlockVotes() *blockVotes {
	return &blockVotes{
		votes: make(map[crypto.Address]*vote.Vote),
		power: 0,
	}
}

func (vs *blockVotes) addVote(vote *vote.Vote) {
	signer := vote.Signer()
	if existing, ok := vs.votes[signer]; ok {
		if !existing.Signature().EqualsTo(vote.Signature()) {
			// Signature malleability?
			logger.Panic("invalid vote", "sig1", existing.Signature().Bytes(), "sig2", vote.Signature().Bytes())
		}
	}

	vs.votes[signer] = vote
}
