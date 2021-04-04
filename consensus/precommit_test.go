package consensus

/*
func TestPrecommitNoProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := 2
	r := 0
	p := makeProposal(t, h, r)

	testEnterNewHeight(tConsP)
	checkStateWait(t, tConsP, h, r, hrs.StepTypePrepare)
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Still no proposal
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

	checkStateWait(t, tConsP, h, r, hrs.StepTypePrecommit)
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

	testEnterNewHeight(tConsP)
	checkState(t, tConsP, h, r, hrs.StepTypePropose)
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Still no proposal
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB)

	checkState(t, tConsP, h, r, hrs.StepTypeCommit)

	// Set proposal now
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

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

	testEnterNewHeight(tConsP)
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	// Validator_1 is offline
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.GenerateTestHash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.GenerateTestHash(), tIndexY)

	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
}

func TestSuspiciousPrepare2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := 3
	r := 0
	p := makeProposal(t, h, r) // Byzantine node send different proposal for every node, all valid

	testEnterNewHeight(tConsP)
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	// Validator_1 is offline
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexY)

	shouldPublishProposal(t, tConsP, p.Hash())
}

func TestPrecommitTimeout(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsP)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexY)

	checkStateWait(t, tConsP, 1, 0, hrs.StepTypePrecommit)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
}

func TestPrecommitIvalidArgs(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsP)

	// Invalid args for propose phase
	// MMMMM tConsP.enterPrecommit(1)
	checkState(t, tConsP, 1, 0, hrs.StepTypePropose)
}
*/
