package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrecommitNoProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := 2
	r := 0
	p := makeProposal(t, h, r)

	tConsP.enterNewHeight()
	checkHRSWait(t, tConsP, h, r, hrs.StepTypePrepare)
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Still no proposal
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB, false)

	checkHRSWait(t, tConsP, h, r, hrs.StepTypePrecommit)
	shouldPublishQueryProposal(t, tConsP, h, r)

	// Set proposal now
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

// This is a worse case scenario
func TestPrecommitNoProposalWithPrecommitQuorom(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := 2
	r := 0
	p := makeProposal(t, h, r)

	tConsP.enterNewHeight()
	checkHRS(t, tConsP, h, r, hrs.StepTypePropose)
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Still no proposal
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB, false)

	checkHRS(t, tConsP, h, r, hrs.StepTypeCommit)

	// Set proposal now
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB, false)

	shouldPublishBlockAnnounce(t, tConsP, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestSuspiciousPrepare1(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := 3
	r := 0
	p := makeProposal(t, h, r) // Byzantine node send different proposal for every node, all valid

	tConsP.enterNewHeight()
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	// Validator_1 is offline
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.GenerateTestHash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.GenerateTestHash(), tIndexY, false)

	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
}

func TestSuspiciousPrepare2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := 3
	r := 0
	p := makeProposal(t, h, r) // Byzantine node send different proposal for every node, all valid

	tConsP.enterNewHeight()
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	// Validator_1 is offline
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexY, false)

	shouldPublishProposal(t, tConsP, p.Hash())
}

func TestPrecommitTimeout(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexY, false)

	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrecommit)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
}

func TestPrecommitIvalidArgs(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()

	// Invalid args for propose phase
	tConsP.enterPrecommit(1)
	checkHRS(t, tConsP, 1, 0, hrs.StepTypePropose)
}

func TestUpdatePrecommitFromPreviousRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	// Byzantine turn to propose a block
	h := 3
	p0 := makeProposal(t, h, 0)

	tConsX.enterNewHeight()
	prepareXRound0Null := shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)

	tConsY.enterNewHeight()
	prepareYRound0Null := shouldPublishVote(t, tConsY, vote.VoteTypePrepare, crypto.UndefHash)

	// Byzantine node set proposal for Partitioned node, but not for others
	tConsP.enterNewHeight()
	tConsP.SetProposal(p0)
	preparePRound0Block := shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p0.Block().Hash())

	assert.NoError(t, tConsX.addVote(prepareYRound0Null))
	assert.NoError(t, tConsX.addVote(preparePRound0Block))
	precommitXRound0Null := shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)

	assert.NoError(t, tConsY.addVote(prepareXRound0Null))
	assert.NoError(t, tConsY.addVote(preparePRound0Block))
	precommitYRound0Null := shouldPublishVote(t, tConsY, vote.VoteTypePrecommit, crypto.UndefHash)

	assert.NoError(t, tConsP.addVote(prepareXRound0Null))
	assert.NoError(t, tConsP.addVote(prepareYRound0Null))
	shouldPublishProposal(t, tConsP, p0.Hash())

	// Byzantine node send its Null votes to partitioned node
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, 0, crypto.UndefHash, tIndexB, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, 0, crypto.UndefHash, tIndexB, false)
	precommitPRound0Null := shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)

	assert.NoError(t, tConsX.addVote(precommitYRound0Null))
	assert.NoError(t, tConsX.addVote(precommitPRound0Null))

	assert.NoError(t, tConsY.addVote(precommitXRound0Null))
	assert.NoError(t, tConsY.addVote(precommitPRound0Null))

	assert.NoError(t, tConsP.addVote(precommitXRound0Null))
	assert.NoError(t, tConsP.addVote(precommitYRound0Null))

	// ConsP can't see others votes
	// It goes to the next round and publish its proposal.
	checkHRSWait(t, tConsP, h, 1, hrs.StepTypePrepare)
	p1 := tConsP.RoundProposal(1)
	assert.NotNil(t, p1)
	preparePRound1Block := shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())
	assert.NotNil(t, preparePRound1Block)

	// Now partitoned heals
	tConsX.SetProposal(p0)
	prepareXRound0Block := shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p0.Block().Hash())
	assert.NotNil(t, prepareXRound0Block)

	tConsY.SetProposal(p0)
	prepareYRound0Block := shouldPublishVote(t, tConsY, vote.VoteTypePrepare, p0.Block().Hash())
	assert.NotNil(t, prepareYRound0Block)

	assert.NoError(t, tConsX.addVote(prepareYRound0Block))
	precommitXRound0Block := shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p0.Block().Hash())
	assert.NotNil(t, precommitXRound0Block)

	assert.NoError(t, tConsY.addVote(prepareXRound0Block))
	precommitYRound0Block := shouldPublishVote(t, tConsY, vote.VoteTypePrecommit, p0.Block().Hash())
	assert.NotNil(t, precommitYRound0Block)

	tConsP.AddVote(prepareXRound0Block)
	tConsP.AddVote(prepareYRound0Block)
	tConsP.AddVote(precommitXRound0Block)
	tConsP.AddVote(precommitYRound0Block)

	shouldPublishBlockAnnounce(t, tConsP, p0.Block().Hash())
}
