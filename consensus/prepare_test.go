package consensus

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestPrepareQueryProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	// After receiving one vote, it should query for proposal (if don't have it yet)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tIndexX)

	shouldPublishQueryProposal(t, tConsP, 2, 0)
}
