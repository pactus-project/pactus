package consensus

// func TestNewHeightTimeout(t *testing.T) {
// 	td := setup(t)

// 	td.enterNewHeight(td.consY)
// 	td.commitBlockForAllStates(t)

// 	consState := &newHeightState{td.consY}
// 	consState.enter()

// 	// Invalid target
// 	consState.onTimeout(&ticker{Height: 2, Target: -1})
// 	td.checkHeightRound(t, td.consY, 2, 0)

// 	consState.onTimeout(&ticker{Height: 2, Target: tickerTargetNewHeight})
// 	td.checkHeightRound(t, td.consY, 2, 0)
// 	td.shouldPublishProposal(t, td.consY, 2, 0)
// }

// func TestNewHeightDoubleEntry(t *testing.T) {
// 	td := setup(t)

// 	td.commitBlockForAllStates(t)

// 	td.consX.MoveToNewHeight()
// 	td.newHeightTimeout(td.consX)

// 	// Double entry to new height state
// 	td.consX.enterNewState(td.consX.newHeightState)

// 	td.checkHeightRound(t, td.consX, 2, 0)
// 	assert.True(t, td.consX.active)
// 	assert.Equal(t, types.Round(0), td.consX.round)
// 	assert.Equal(t, types.Height(2), td.consX.height)
// }

// func TestNewHeightTimeBehindNetwork(t *testing.T) {
// 	td := setup(t)

// 	td.commitBlockForAllStates(t)
// 	td.consP.MoveToNewHeight()

// 	height := types.Height(2)
// 	round := types.Round(0)
// 	prop := td.makeProposal(t, height, round)

// 	td.consP.SetProposal(prop)
// 	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
// 	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
// 	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

// 	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, prop.Block().Hash())
// 	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop.Block().Hash())
// }
