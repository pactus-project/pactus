package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
)

func TestPrepareQueryProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)
	shouldPublishQueryProposal(t, tConsP, 2, 0)
}

func TestEnterPrecommitWithoutProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsB)
	shouldPublishQueryProposal(t, tConsB, 2, 0)

	p := makeProposal(t, 2, 0)
	testAddVote(t, tConsB, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsB, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsB, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexP)

	assert.Equal(t, tConsB.currentState.name(), "precommit")
}
