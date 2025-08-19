package consensusv2

// func TestByzantineProposal(t *testing.T) {
// 	td := setup(t)

// 	td.commitBlockForAllStates(t)
// 	td.commitBlockForAllStates(t)
// 	h := uint32(3)
// 	r := int16(0)
// 	prop := td.makeProposal(t, h, r)
// 	propBlockHash := prop.Block().Hash()

// 	td.enterNewHeight(td.consP)

// 	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexX)
// 	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexY)
// 	td.addPrecommitVote(td.consP, propBlockHash, h, r, tIndexB)

// 	assert.Nil(t, td.consP.Proposal())

// 	// Byzantine node sends second proposal to Partitioned node.
// 	trx := tx.NewTransferTx(h, td.consX.rewardAddr, td.RandAccAddress(), 1000, 1000)
// 	td.HelperSignTransaction(td.consX.valKey.PrivateKey(), trx)
// 	assert.NoError(t, td.txPool.AppendTx(trx))
// 	byzProp := td.makeProposal(t, h, r)
// 	assert.NotEqual(t, prop.Hash(), byzProp.Hash())

// 	td.consP.SetProposal(byzProp)
// 	assert.Nil(t, td.consP.Proposal())
// 	td.shouldPublishQueryProposal(t, td.consP, h, r)
// 	td.checkHeightRound(t, td.consP, h, r)
// }
