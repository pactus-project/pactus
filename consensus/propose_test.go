package consensus

/*

func TestSetProposalInvalidProposer(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsY)
	assert.Nil(t, tConsY.RoundProposal(0))

	addr := tSigners[tIndexB].Address()
	b, _ := block.GenerateTestBlock(&addr, nil)
	p := proposal.NewProposal(1, 0, *b)

	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(0))

	tSigners[tIndexB].SignMsg(p) // Invalid signature
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(0))
}

func TestSetProposalInvalidBlock(t *testing.T) {
	setup(t)

	a := tSigners[tIndexB].Address()
	invBlock, _ := block.GenerateTestBlock(&a, nil)
	p := proposal.NewProposal(1, 2, *invBlock)
	tSigners[tIndexB].SignMsg(p)

	testEnterNewHeight(tConsY)
	// MMMM tConsY.enterNewRound(2)
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(2))
}

func TestSetProposalInvalidHeight(t *testing.T) {
	setup(t)

	a := tSigners[tIndexB].Address()
	invBlock, _ := block.GenerateTestBlock(&a, nil)
	p := proposal.NewProposal(2, 0, *invBlock)
	tSigners[tIndexB].SignMsg(p)

	testEnterNewHeight(tConsY)
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(2))
}

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)

	testEnterNewHeight(tConsP)
	commitBlockForAllStates(t)
	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(0))
}

func TestGotoNextRoundWithoutProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexB)

	checkState(t, tConsP, 3, 1, hrs.StepTypePrepare)
}

func TestSecondProposalCommitted(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	// Now it's turn for Byzantine node to propose a block
	// Other nodes are going to not accept its proposal, even it is valid
	p1 := makeProposal(t, 3, 0) // valid proposal for round 0, byzantine proposer
	p2 := makeProposal(t, 3, 1) // valid proposal for round 1, partitioned proposer

	// Probably we have blocked Byzantine node
	//tConsX.SetProposal(p1)

	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, p1.Block().Hash(), tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexP)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, p1.Block().Hash(), tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexP)

	tConsX.SetProposal(p2)

	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p2.Block().Hash())
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, crypto.UndefHash, tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexP)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p2.Block().Hash())
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, crypto.UndefHash, tIndexB)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexP)

	shouldPublishBlockAnnounce(t, tConsX, p2.Block().Hash())
}

func TestNetworkLagging1(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsP)

	h := 1
	r := 0
	p := makeProposal(t, h, r)
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p)

	checkStateWait(t, tConsP, h, r, hrs.StepTypePrepare)
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexB)

	// Now let's set the proposal
	tConsP.SetProposal(p)
	checkState(t, tConsP, h, r, hrs.StepTypePrecommit)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestNetworkLagging2(t *testing.T) {
	setup(t)

	h := 1
	r := 0
	p1 := makeProposal(t, h, r)

	testEnterNewHeight(tConsP)
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p1)

	// Networks lags and we don't receive prepare from val_1 and pre-commit from val_4
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexY)

	checkState(t, tConsP, h, r, hrs.StepTypePropose)

	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Now let's set the proposal
	tConsP.SetProposal(p1)

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())
	checkState(t, tConsP, h, r, hrs.StepTypePrepare)

	// We can't go to precommit stage, because we haven't prepared yet
	// But if we receive another vote we go to commit phase directly
	// Let's do it
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexB)
	shouldPublishBlockAnnounce(t, tConsP, p1.Block().Hash())
}

func TestLateProposal(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsP)

	h := 1
	r := 0
	p := makeProposal(t, h, r)

	// tConsP is partitioned, so tConsP doesn't have the proposal
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexB)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexB)

	// Now partition healed.
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	/// MMMMM assert.True(t, tConsP.isCommitted)
}

func TestLateProposal2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := 3
	p := makeProposal(t, h, 0) // tConsP should propose for this round

	testEnterNewHeight(tConsX)

	// tConsP is partitioned, so tConsX doesn't have the proposal
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, 0, crypto.UndefHash, tIndexB)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, 0, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, 0, crypto.UndefHash, tIndexB)

	checkStateWait(t, tConsX, h, 1, hrs.StepTypePrepare)

	tConsX.SetProposal(p)

	checkState(t, tConsX, h, 1, hrs.StepTypePrepare)
}

func TestLateUndefVote(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	h := 3
	r := 0
	p := makeProposal(t, h, r) // Other nodes doesn't accept byzantine proposal

	// tConsP is partitioned
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB)

	//  tConsP doesn't accept tConsB's proposal
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Now partition healed
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexX)

	// Enough prepare votes, tConsP can vote for undef precommit
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, crypto.UndefHash)
}

func TestSetProposalForNextRoundWithoutFinishingTheFirstRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	// Byzantine node sends proposal for second round (his turn)
	b, err := tConsB.state.ProposeBlock(1)
	assert.NoError(t, err)
	p := proposal.NewProposal(2, 1, *b)
	tSigners[tIndexB].SignMsg(p)

	tConsX.SetProposal(p)
	// tConsX doesn't accept the proposal for next rounds
	assert.Nil(t, tConsX.RoundProposal(1))

	// But doesn't move to prepare phase
	checkState(t, tConsX, 2, 0, hrs.StepTypePropose)
}

func TestEnterPrepareAfterPrecommit(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := 4
	r := 0
	p := makeProposal(t, h, r)

	// tConsP is partitioned, so tConsX doesn't have the proposal
	testEnterNewHeight(tConsX)
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)

	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexB)
	checkState(t, tConsX, h, r, hrs.StepTypePrecommit)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, crypto.GenerateTestHash(), tIndexB)

	// Now partition healed
	tConsX.SetProposal(p)
	/// MMMMM tConsX.enterPrepare(0)
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p.Block().Hash())

	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexP)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestProposeIvalidArgs(t *testing.T) {
	setup(t)

	// MMMM tConsP.hrs = hrs.NewHRS(1, 0, hrs.StepTypeNewHeight)
	// Invalid args for propose phase
	/// MMMMM tConsP.enterPropose(1)
	checkState(t, tConsP, 1, 0, hrs.StepTypeNewHeight)
}

func TestCreateProposal(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsX)
	testEnterNewHeight(tConsY)

	/// MMMMM tConsX.createProposal(1, 0)
	assert.NotNil(t, tConsX.RoundProposal(0))

	/// MMMMM tConsY.createProposal(1, 0)
	assert.Nil(t, tConsY.RoundProposal(0))
}
*/
