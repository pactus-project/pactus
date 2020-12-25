package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

//Imagine we have four nodes: (Nx, Ny, Nb, Np) which:
// Nb is a byzantine node and Nx, Ny, Np are honest nodes,
// however Np is partitioned and see the network through Nb (Byzantine node).
// In Height H, B sends its pre-votes to all the nodes
// but only sends valid pre-commit to P and nil pre-commit to X,Y.
// For should not hapens
func TestByzantineVote(t *testing.T) {
	cons := newTestConsensus(t, VAL4)
	cons.enterNewHeight(1)

	p := makeTestProposal(t, VAL1, 1, 0)
	cons.SetProposal(p)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), VAL1, false)
	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL1, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL2, false) // Byzantine vote

	cons.enterNewRound(1, 1)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommitWait)
}

func TestConsensusGotoNextRound(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	cons.enterNewHeight(1)

	// Validator_1 is offline
	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, VAL3, false)
	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, VAL4, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL3, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL4, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrepare)

	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 1, p.Block().Hash(), VAL1, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrepare)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 1, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL1, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
}

func TestConsensusGotoNextRound2(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	cons.enterNewHeight(1)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, crypto.GenerateTestHash(), VAL1, false)
	checkHRSWait(t, cons, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), VAL1, false)
	checkHRSWait(t, cons, 1, 0, hrs.StepTypePrepare)

	shouldPublishProposalReqquest(t, cons)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL2, false)
	checkHRSWait(t, cons, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, VAL3, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL3, false)
	checkHRSWait(t, cons, 1, 0, hrs.StepTypePrecommitWait)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, VAL4, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL4, false)
	checkHRSWait(t, cons, 1, 1, hrs.StepTypePropose)
}
