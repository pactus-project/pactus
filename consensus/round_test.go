package consensus

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestNewRound(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()
	checkHRS(t, tConsP, 1, 0, hrs.StepTypePropose)

	//
	// 1- Move to round 0
	// 2- PreCommits  for round 0 => missed
	// 3- PreCommits  for round 1 => missed
	// 4- PreCommits  for round 2 => received
	// 5- Moved to round 3
	// 6- PreCommits  for round 0 => received
	// 7- Should ignore them

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 2, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 2, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 2, crypto.UndefHash, tIndexB, false)

	checkHRS(t, tConsP, 1, 3, hrs.StepTypePrepare)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexB, false)

	checkHRS(t, tConsP, 1, 3, hrs.StepTypePrepare)
}

func TestConsensusGotoNextRound(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()

	// Validator_1 is offline
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexY, false)
	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexY, false)
	checkHRSWait(t, tConsP, 1, 1, hrs.StepTypePrepare)
}

func TestConsensusGotoNextRound2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsP.enterNewHeight()

	// Byzantine node sends different valid proposals for every node
	h := 3
	r := 0
	p := makeProposal(t, h, r)
	tConsP.SetProposal(p)

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.GenerateTestHash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.GenerateTestHash(), tIndexY, false)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexY, false)
	checkHRSWait(t, tConsP, h, r+1, hrs.StepTypePrepare)
}

func TestEnterNewRound(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()

	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexY, false)
	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrecommit)
	tConsP.enterNewRound(0)
	checkHRS(t, tConsP, 1, 0, hrs.StepTypePrecommit)
}
