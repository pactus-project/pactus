package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestNewRound(t *testing.T) {
	cons1 := newTestConsensus(t, VAL1)
	cons4 := newTestConsensus(t, VAL4)

	cons1.enterNewHeight(1)
	cons4.enterNewHeight(1)

	//
	// 1- Move to round 0
	// 2- PreCommits  for round 0 => missed
	// 3- PreCommits  for round 1 => missed
	// 4- PreCommits  for round 2 => received
	// 5- Moved to round 3
	// 6- PreCommits  for round 0 => received
	// 7- Should ignore moving to round 1

	voteRound0Val1 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL1, false)
	voteRound0Val2 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL2, false)
	voteRound0Val3 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL3, false)

	voteRound2Val1 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 2, crypto.UndefHash, VAL1, false)
	voteRound2Val2 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 2, crypto.UndefHash, VAL2, false)
	voteRound2Val3 := testAddVote(t, cons1, vote.VoteTypePrecommit, 1, 2, crypto.UndefHash, VAL3, false)

	assert.NoError(t, cons4.addVote(voteRound2Val1))
	assert.NoError(t, cons4.addVote(voteRound2Val2))
	assert.NoError(t, cons4.addVote(voteRound2Val3))

	checkHRSWait(t, cons4, 1, 3, hrs.StepTypePrepare)

	assert.NoError(t, cons4.addVote(voteRound0Val1))
	assert.NoError(t, cons4.addVote(voteRound0Val2))
	assert.NoError(t, cons4.addVote(voteRound0Val3))

	checkHRS(t, cons4, 1, 3, hrs.StepTypePrepare)
}
