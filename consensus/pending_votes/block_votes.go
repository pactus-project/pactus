package pending_votes

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/vote"
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

func (vs *blockVotes) addVote(vote *vote.Vote) bool {
	signer := vote.Signer()
	if existing, ok := vs.votes[signer]; ok {
		if !existing.Signature().EqualsTo(*vote.Signature()) {
			// Signature malleability?
			logger.Panic("Invalid vote")
		} else {
			//
			return false
		}
	}

	vs.votes[signer] = vote

	return true
}
