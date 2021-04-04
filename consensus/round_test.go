package consensus

import (
	"testing"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestGotoNextRoundWithoutProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexY)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexY)

	checkHeightRoundWait(t, tConsP, 3, 1)
}
