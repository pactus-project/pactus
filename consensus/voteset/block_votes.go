package voteset

import (
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
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
		if !existing.Signature().EqualsTo(*vote.Signature()) {
			// Signature malleability?
			logger.Panic("Invalid vote", "sig1", existing.Signature().RawBytes(), "sig2", vote.Signature().RawBytes())
		}
	}

	vs.votes[signer] = vote
}
