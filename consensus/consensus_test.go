package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tIndexX = 0
	tIndexY = 1
	tIndexB = 2
	tIndexP = 3
)

type testData struct {
	*testsuite.TestSuite

	signers []crypto.Signer
	txPool  *txpool.MockTxPool
	genDoc  *genesis.Genesis
	consX   *consensus // Good peer
	consY   *consensus // Good peer
	consB   *consensus // Byzantine or offline peer
	consP   *consensus // Partitioned peer
}

type OverrideFingerprint struct {
	cons Consensus
	name string
}

func testConfig() *Config {
	return &Config{
		ChangeProposerTimeout: 1 * time.Second,
		ChangeProposerDelta:   200 * time.Millisecond,
	}
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.cons.Fingerprint()
}

func setup(t *testing.T) *testData {
	ts := testsuite.NewTestSuite(t)

	_, signers := ts.GenerateTestCommittee(4)
	txPool := txpool.MockingTxPool()

	vals := make([]*validator.Validator, 4)
	for i, s := range signers {
		val := validator.NewValidator(s.PublicKey().(*bls.PublicKey), int32(i))
		vals[i] = val
	}

	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	accs := map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc}
	params := param.DefaultParams()
	params.CommitteeSize = 4
	params.BlockTimeInSecond = 1

	// to prevent triggering timers before starting the tests to avoid double entries for new heights in some tests.
	getTime := util.RoundNow(params.BlockTimeInSecond).Add(time.Duration(params.BlockTimeInSecond) * time.Second)
	genDoc := genesis.MakeGenesis(getTime, accs, vals, params)
	stX, err := state.LoadOrNewState(genDoc, []crypto.Signer{signers[tIndexX]},
		store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)
	stY, err := state.LoadOrNewState(genDoc, []crypto.Signer{signers[tIndexY]},
		store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)
	stB, err := state.LoadOrNewState(genDoc, []crypto.Signer{signers[tIndexB]},
		store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)
	stP, err := state.LoadOrNewState(genDoc, []crypto.Signer{signers[tIndexP]},
		store.MockingStore(ts), txPool, nil)
	require.NoError(t, err)

	ConsX := NewConsensus(testConfig(), stX, signers[tIndexX], signers[tIndexX].Address(),
		make(chan message.Message, 100), newMediator())
	ConsY := NewConsensus(testConfig(), stY, signers[tIndexY], signers[tIndexY].Address(),
		make(chan message.Message, 100), newMediator())
	ConsB := NewConsensus(testConfig(), stB, signers[tIndexB], signers[tIndexB].Address(),
		make(chan message.Message, 100), newMediator())
	ConsP := NewConsensus(testConfig(), stP, signers[tIndexP], signers[tIndexP].Address(),
		make(chan message.Message, 100), newMediator())

	consX := ConsX.(*consensus)
	consY := ConsY.(*consensus)
	consB := ConsB.(*consensus)
	consP := ConsP.(*consensus)

	// -------------------------------
	// For better logging when testing
	overrideLogger := func(cons *consensus, name string) {
		cons.logger = logger.NewLogger("_consensus",
			&OverrideFingerprint{name: fmt.Sprintf("%s - %s: ", name, t.Name()), cons: cons})
	}

	overrideLogger(consX, "consX")
	overrideLogger(consY, "consY")
	overrideLogger(consB, "consB")
	overrideLogger(consP, "consP")
	// -------------------------------

	logger.Info("setup finished, start running the test", "name", t.Name())

	return &testData{
		TestSuite: ts,
		signers:   signers,
		txPool:    txPool,
		genDoc:    genDoc,
		consX:     consX,
		consY:     consY,
		consB:     consB,
		consP:     consP,
	}
}

func (td *testData) shouldPublishBlockAnnounce(t *testing.T, cons *consensus, hash hash.Hash) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishBlockAnnounce", "message", msg)

			if msg.Type() == message.MessageTypeBlockAnnounce {
				m := msg.(*message.BlockAnnounceMessage)
				assert.Equal(t, m.Block.Hash(), hash)
				return
			}
		}
	}
}

func (td *testData) shouldPublishProposal(t *testing.T, cons *consensus,
	height uint32, round int16) *proposal.Proposal {
	return shouldPublishProposal(t, cons, height, round)
}

func shouldPublishProposal(t *testing.T, cons *consensus,
	height uint32, round int16) *proposal.Proposal {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return nil
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishProposal", "message", msg)

			if msg.Type() == message.MessageTypeProposal {
				m := msg.(*message.ProposalMessage)
				require.Equal(t, m.Proposal.Height(), height)
				require.Equal(t, m.Proposal.Round(), round)
				return m.Proposal
			}
		}
	}
}

func (td *testData) shouldPublishQueryProposal(t *testing.T, cons *consensus, height uint32, round int16) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishQueryProposal", "message", msg)

			if msg.Type() == message.MessageTypeQueryProposal {
				m := msg.(*message.QueryProposalMessage)
				assert.Equal(t, m.Height, height)
				assert.Equal(t, m.Round, round)
				return
			}
		}
	}
}

func (td *testData) shouldPublishVote(t *testing.T, cons *consensus, voteType vote.Type, hash hash.Hash) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishVote", "message", msg)

			if msg.Type() == message.MessageTypeVote {
				m := msg.(*message.VoteMessage)
				if m.Vote.Type() == voteType &&
					m.Vote.BlockHash().EqualsTo(hash) {
					return
				}
			}
		}
	}
}

func checkHeightRound(t *testing.T, cons *consensus, height uint32, round int16) {
	h, r := cons.HeightRound()
	assert.Equal(t, h, height)
	assert.Equal(t, r, round)
}

func (td *testData) checkHeightRound(t *testing.T, cons *consensus, height uint32, round int16) {
	checkHeightRound(t, cons, height, round)
}

func checkHeightRoundWait(t *testing.T, cons *consensus, height uint32, round int16) {
	for i := 0; i < 20; i++ {
		h, r := cons.HeightRound()
		if h == height && r == round {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	checkHeightRound(t, cons, height, round)
}

func (td *testData) checkHeightRoundWait(t *testing.T, cons *consensus, height uint32, round int16) {
	checkHeightRoundWait(t, cons, height, round)
}

func (td *testData) addVote(cons *consensus, voteType vote.Type, height uint32, round int16,
	blockHash hash.Hash, valID int) *vote.Vote {
	v := vote.NewVote(voteType, height, round, blockHash, td.signers[valID].Address())
	td.signers[valID].SignMsg(v)

	cons.AddVote(v)

	return v
}

// enterNewHeight helps tests to enter new height safely
// without scheduling new height. It boosts the test speed.
func (td *testData) enterNewHeight(cons *consensus) {
	cons.lk.Lock()
	cons.enterNewState(cons.newHeightState)
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetNewHeight})
	cons.lk.Unlock()
}

// enterNextRound helps tests to enter next round safely.
func (td *testData) enterNextRound(cons *consensus) {
	cons.lk.Lock()
	cons.round++
	cons.enterNewState(cons.proposeState)
	cons.lk.Unlock()
}

func (td *testData) commitBlockForAllStates(t *testing.T) (*block.Block, *block.Certificate) {
	height := td.consX.state.LastBlockHeight()
	var err error
	p := td.makeProposal(t, height+1, 0)

	sb := block.CertificateSignBytes(p.Block().Hash(), 0)
	sig1 := td.signers[0].SignData(sb).(*bls.Signature)
	sig2 := td.signers[1].SignData(sb).(*bls.Signature)
	sig4 := td.signers[3].SignData(sb).(*bls.Signature)

	sig := bls.SignatureAggregate([]*bls.Signature{sig1, sig2, sig4})
	cert := block.NewCertificate(0, []int32{0, 1, 2, 3}, []int32{2}, sig)
	block := p.Block()

	err = td.consX.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)
	err = td.consY.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)
	err = td.consB.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)
	err = td.consP.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)

	return block, cert
}

func (td *testData) makeProposal(t *testing.T, height uint32, round int16) *proposal.Proposal {
	var p *proposal.Proposal
	switch (height % 4) + uint32(round) {
	case 1:
		blk, err := td.consX.state.ProposeBlock(td.consX.signer, td.consX.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.consX.signer.SignMsg(p)
	case 2:
		blk, err := td.consY.state.ProposeBlock(td.consY.signer, td.consY.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.consY.signer.SignMsg(p)
	case 3:
		blk, err := td.consB.state.ProposeBlock(td.consB.signer, td.consB.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.consB.signer.SignMsg(p)
	case 0, 4:
		blk, err := td.consP.state.ProposeBlock(td.consP.signer, td.consP.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		td.consP.signer.SignMsg(p)
	}

	return p
}

func TestNotInCommittee(t *testing.T) {
	td := setup(t)

	_, prv := td.RandomBLSKeyPair()
	signer := crypto.NewSigner(prv)
	store := store.MockingStore(td.TestSuite)

	st, _ := state.LoadOrNewState(td.genDoc, []crypto.Signer{signer}, store, td.txPool, nil)
	Cons := NewConsensus(testConfig(), st, signer, signer.Address(), make(chan message.Message, 100),
		newMediator())
	cons := Cons.(*consensus)

	td.enterNewHeight(cons)

	cons.signAddVote(vote.VoteTypePrepare, td.RandomHash())
}

func TestRoundVotes(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1
	td.enterNewHeight(td.consP)

	t.Run("Ignore votes from invalid height", func(t *testing.T) {
		v1 := td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 0, td.RandomHash(), tIndexX)
		v2 := td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, td.RandomHash(), tIndexX)
		v3 := td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, td.RandomHash(), tIndexY)
		v4 := td.addVote(td.consP, vote.VoteTypeChangeProposer, 3, 0, td.RandomHash(), tIndexX)

		require.False(t, td.consP.HasVote(v1.Hash()))
		require.True(t, td.consP.HasVote(v2.Hash()))
		require.True(t, td.consP.HasVote(v3.Hash()))
		require.False(t, td.consP.HasVote(v4.Hash()))
	})
}

func TestConsensusAddVotesNormal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consX)
	td.checkHeightRound(t, td.consX, 2, 0)

	p := td.makeProposal(t, 2, 0)
	td.consX.SetProposal(p)

	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	td.addVote(td.consX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexP)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrepare, p.Block().Hash())

	td.addVote(td.consX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexY)
	td.addVote(td.consX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexP)
	td.shouldPublishVote(t, td.consX, vote.VoteTypePrecommit, p.Block().Hash())
	td.shouldPublishBlockAnnounce(t, td.consX, p.Block().Hash())
}

func TestConsensusAddVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	v1 := td.addVote(td.consP, vote.VoteTypePrepare, 2, 0, td.RandomHash(), tIndexX)
	v2 := td.addVote(td.consP, vote.VoteTypePrepare, 1, 0, td.RandomHash(), tIndexX)
	v3 := td.addVote(td.consP, vote.VoteTypePrecommit, 1, 0, td.RandomHash(), tIndexX)
	v4 := td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 0, td.RandomHash(), tIndexX)
	v5 := td.addVote(td.consP, vote.VoteTypePrepare, 1, 2, td.RandomHash(), tIndexX)

	assert.False(t, td.consP.HasVote(v1.Hash())) // invalid height
	assert.True(t, td.consP.HasVote(v2.Hash()))
	assert.True(t, td.consP.HasVote(v3.Hash()))
	assert.True(t, td.consP.HasVote(v4.Hash()))
	assert.True(t, td.consP.HasVote(v5.Hash())) // next round
}

// TestConsensusLateProposal tests the scenario where a slow node receives a proposal
// after votes have been cast.
func TestConsensusLateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	h := uint32(2)
	r := int16(0)
	p := td.makeProposal(t, h, r)
	require.NotNil(t, p)

	// The partitioned node receives all the votes first
	td.addVote(td.consP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB)

	td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

	// Partitioned node receives proposal now
	td.consP.SetProposal(p)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, p.Block().Hash())
	td.shouldPublishBlockAnnounce(t, td.consP, p.Block().Hash())
}

// TestConsensusVeryLateProposal tests the scenario where a slow node receives a proposal
// after a block is committed.
func TestConsensusVeryLateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	h := uint32(2)
	r := int16(0)
	p := td.makeProposal(t, h, r)
	require.NotNil(t, p)

	td.commitBlockForAllStates(t) // height 2

	// Partitioned node doesn't receive all the votes
	td.addVote(td.consP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB)

	td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB)

	// Partitioned node receives proposal now
	td.consP.SetProposal(p)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, p.Block().Hash())
	td.shouldPublishBlockAnnounce(t, td.consP, p.Block().Hash())
}

func TestConsensusInvalidVote(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)

	v1, _ := td.GenerateTestPrecommitVote(1, 0)
	v2 := vote.NewVote(vote.VoteTypePrepare, 2, 0, td.RandomHash(),
		td.signers[tIndexB].Address())
	td.signers[tIndexB].SignMsg(v2)

	td.consX.AddVote(v1)
	td.consX.AddVote(v2)
	assert.False(t, td.consX.HasVote(v1.Hash()))
	assert.False(t, td.consX.HasVote(v2.Hash()))
}

func TestPickRandomVote(t *testing.T) {
	td := setup(t)

	assert.Nil(t, td.consP.PickRandomVote())

	td.enterNewHeight(td.consP)
	assert.Nil(t, td.consP.PickRandomVote())

	p1 := td.makeProposal(t, 1, 0)
	p2 := td.makeProposal(t, 1, 1)

	// round 0
	td.addVote(td.consP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexP)
	td.addVote(td.consP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexY)

	assert.NotNil(t, td.consP.PickRandomVote())

	td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexX)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexY)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexP)

	// Round 1
	td.addVote(td.consP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexX)
	td.addVote(td.consP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexY)
	td.addVote(td.consP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexP)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 1, hash.UndefHash, tIndexX)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 1, hash.UndefHash, tIndexY)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 1, hash.UndefHash, tIndexP)

	// Round 2
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 1, 2, hash.UndefHash, tIndexP)

	for i := 0; i < 10; i++ {
		rndVote := td.consP.PickRandomVote()
		assert.NotNil(t, rndVote)
		assert.Equal(t, rndVote.Type(), vote.VoteTypeChangeProposer,
			"Should only pick Change Proposer votes")
	}
}

func TestSetProposalFromPreviousRound(t *testing.T) {
	td := setup(t)

	p := td.makeProposal(t, 1, 0)
	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	// It should ignore proposal for previous rounds
	td.consP.SetProposal(p)

	assert.Nil(t, td.consP.RoundProposal(0))
	td.checkHeightRoundWait(t, td.consP, 1, 1)
}

func TestSetProposalFromPreviousHeight(t *testing.T) {
	td := setup(t)

	p := td.makeProposal(t, 1, 0)
	td.commitBlockForAllStates(t) // height 1

	td.enterNewHeight(td.consP)

	td.consP.SetProposal(p)
	assert.Nil(t, td.consP.RoundProposal(0))
	td.checkHeightRoundWait(t, td.consP, 2, 0)
}

func TestDuplicateProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	h := uint32(4)
	r := int16(0)
	p1 := td.makeProposal(t, h, r)
	trx := tx.NewTransferTx(hash.UndefHash.Stamp(), 1, td.signers[0].Address(),
		td.signers[1].Address(), 1000, 1000, "proposal changer")
	td.signers[0].SignMsg(trx)
	assert.NoError(t, td.txPool.AppendTx(trx))
	p2 := td.makeProposal(t, h, r)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	td.consX.SetProposal(p1)
	td.consX.SetProposal(p2)

	assert.Equal(t, td.consX.RoundProposal(0).Hash(), p1.Hash())
}

func TestNonActiveValidator(t *testing.T) {
	td := setup(t)

	signer := td.RandomSigner()
	Cons := NewConsensus(testConfig(), state.MockingState(td.TestSuite),
		signer, signer.Address(), make(chan message.Message, 100), newMediator())
	nonActiveCons := Cons.(*consensus)

	t.Run("non-active instances should be in new-height state", func(t *testing.T) {
		nonActiveCons.MoveToNewHeight()
		td.checkHeightRoundWait(t, nonActiveCons, 1, 0)

		assert.False(t, nonActiveCons.IsActive())
		assert.Equal(t, nonActiveCons.currentState.name(), "new-height")
	})

	t.Run("non-active instances should ignore proposals", func(t *testing.T) {
		p := td.makeProposal(t, 1, 0)
		nonActiveCons.SetProposal(p)

		assert.Nil(t, nonActiveCons.RoundProposal(0))
	})

	t.Run("non-active instances should ignore votes", func(t *testing.T) {
		v := td.addVote(nonActiveCons, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexX)

		assert.False(t, nonActiveCons.HasVote(v.Hash()))
	})

	t.Run("non-active instances should move to new height", func(t *testing.T) {
		b1, cert1 := td.commitBlockForAllStates(t)
		b2, cert2 := td.commitBlockForAllStates(t)

		nonActiveCons.MoveToNewHeight()
		td.checkHeightRoundWait(t, nonActiveCons, 1, 0)

		assert.NoError(t, nonActiveCons.state.CommitBlock(1, b1, cert1))
		assert.NoError(t, nonActiveCons.state.CommitBlock(2, b2, cert2))

		nonActiveCons.MoveToNewHeight()
		td.checkHeightRoundWait(t, nonActiveCons, 3, 0)
	})
}

func TestValidVoteWithBigRound(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)

	v := vote.NewVote(vote.VoteTypeChangeProposer, 1, util.MaxInt16, hash.UndefHash,
		td.signers[tIndexB].Address())
	td.signers[tIndexB].SignMsg(v)

	td.consX.AddVote(v)
	assert.False(t, td.consX.HasVote(v.Hash()))
}
