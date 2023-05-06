package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/committee"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tIndexX = 0
	tIndexY = 1
	tIndexB = 2
	tIndexP = 3
)

var (
	tSigners []crypto.Signer
	tTxPool  *txpool.MockTxPool
	tGenDoc  *genesis.Genesis
	tConsX   *consensus // Good peer
	tConsY   *consensus // Good peer
	tConsB   *consensus // Byzantine or offline peer
	tConsP   *consensus // Partitioned peer
)

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

func setup(t *testing.T) {
	_, tSigners = committee.GenerateTestCommittee(4)
	tTxPool = txpool.MockingTxPool()

	vals := make([]*validator.Validator, 4)
	for i, s := range tSigners {
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
	tGenDoc = genesis.MakeGenesis(getTime, accs, vals, params)
	stX, err := state.LoadOrNewState(tGenDoc, []crypto.Signer{tSigners[tIndexX]}, store.MockingStore(), tTxPool, nil)
	require.NoError(t, err)
	stY, err := state.LoadOrNewState(tGenDoc, []crypto.Signer{tSigners[tIndexY]}, store.MockingStore(), tTxPool, nil)
	require.NoError(t, err)
	stB, err := state.LoadOrNewState(tGenDoc, []crypto.Signer{tSigners[tIndexB]}, store.MockingStore(), tTxPool, nil)
	require.NoError(t, err)
	stP, err := state.LoadOrNewState(tGenDoc, []crypto.Signer{tSigners[tIndexP]}, store.MockingStore(), tTxPool, nil)
	require.NoError(t, err)

	consX := NewConsensus(testConfig(), stX, tSigners[tIndexX], tSigners[tIndexX].Address(),
		make(chan message.Message, 100), newMediator())
	consY := NewConsensus(testConfig(), stY, tSigners[tIndexY], tSigners[tIndexY].Address(),
		make(chan message.Message, 100), newMediator())
	consB := NewConsensus(testConfig(), stB, tSigners[tIndexB], tSigners[tIndexB].Address(),
		make(chan message.Message, 100), newMediator())
	consP := NewConsensus(testConfig(), stP, tSigners[tIndexP], tSigners[tIndexP].Address(),
		make(chan message.Message, 100), newMediator())

	tConsX = consX.(*consensus)
	tConsY = consY.(*consensus)
	tConsB = consB.(*consensus)
	tConsP = consP.(*consensus)

	// -------------------------------
	// For better logging when testing
	overrideLogger := func(cons *consensus, name string) {
		cons.logger = logger.NewLogger("_consensus",
			&OverrideFingerprint{name: fmt.Sprintf("%s - %s: ", name, t.Name()), cons: cons})
	}

	overrideLogger(tConsX, "consX")
	overrideLogger(tConsY, "consY")
	overrideLogger(tConsB, "consB")
	overrideLogger(tConsP, "consP")
	// -------------------------------

	logger.Info("setup finished, start running the test", "name", t.Name())
}

func shouldPublishBlockAnnounce(t *testing.T, cons *consensus, hash hash.Hash) {
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

func shouldPublishProposal(t *testing.T, cons *consensus, height uint32, round int16) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishProposal", "message", msg)

			if msg.Type() == message.MessageTypeProposal {
				m := msg.(*message.ProposalMessage)
				assert.Equal(t, m.Proposal.Height(), height)
				assert.Equal(t, m.Proposal.Round(), round)
				return
			}
		}
	}
}

func shouldPublishQueryProposal(t *testing.T, cons *consensus, height uint32, round int16) {
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

func shouldPublishVote(t *testing.T, cons *consensus, voteType vote.Type, hash hash.Hash) {
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

func testAddVote(cons *consensus, voteType vote.Type, height uint32, round int16,
	blockHash hash.Hash, valID int) *vote.Vote {
	v := vote.NewVote(voteType, height, round, blockHash, tSigners[valID].Address())
	tSigners[valID].SignMsg(v)

	cons.AddVote(v)

	return v
}

// testEnterNewHeight helps tests to enter new height safely
// without scheduling new height. It boosts the test speed.
func testEnterNewHeight(cons *consensus) {
	cons.lk.Lock()
	cons.enterNewState(cons.newHeightState)
	cons.currentState.onTimeout(&ticker{0, cons.height, cons.round, tickerTargetNewHeight})
	cons.lk.Unlock()
}

// testEnterNextRound helps tests to enter next round safely.
func testEnterNextRound(cons *consensus) {
	cons.lk.Lock()
	cons.round++
	cons.enterNewState(cons.proposeState)
	cons.lk.Unlock()
}

func commitBlockForAllStates(t *testing.T) (*block.Block, *block.Certificate) {
	height := tConsX.state.LastBlockHeight()
	var err error
	p := makeProposal(t, height+1, 0)

	sb := block.CertificateSignBytes(p.Block().Hash(), 0)
	sig1 := tSigners[0].SignData(sb).(*bls.Signature)
	sig2 := tSigners[1].SignData(sb).(*bls.Signature)
	sig4 := tSigners[3].SignData(sb).(*bls.Signature)

	sig := bls.Aggregate([]*bls.Signature{sig1, sig2, sig4})
	cert := block.NewCertificate(0, []int32{0, 1, 2, 3}, []int32{2}, sig)
	block := p.Block()

	err = tConsX.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)
	err = tConsY.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)
	err = tConsB.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)
	err = tConsP.state.CommitBlock(height+1, block, cert)
	assert.NoError(t, err)

	return block, cert
}

func makeProposal(t *testing.T, height uint32, round int16) *proposal.Proposal {
	var p *proposal.Proposal
	switch (height % 4) + uint32(round) {
	case 1:
		blk, err := tConsX.state.ProposeBlock(tConsX.signer, tConsX.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		tConsX.signer.SignMsg(p)
	case 2:
		blk, err := tConsY.state.ProposeBlock(tConsY.signer, tConsY.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		tConsY.signer.SignMsg(p)
	case 3:
		blk, err := tConsB.state.ProposeBlock(tConsB.signer, tConsB.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		tConsB.signer.SignMsg(p)
	case 0, 4:
		blk, err := tConsP.state.ProposeBlock(tConsP.signer, tConsP.rewardAddr, round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, blk)
		tConsP.signer.SignMsg(p)
	}

	return p
}

func TestNotInCommittee(t *testing.T) {
	setup(t)

	_, prv := bls.GenerateTestKeyPair()
	signer := crypto.NewSigner(prv)
	store := store.MockingStore()

	st, _ := state.LoadOrNewState(tGenDoc, []crypto.Signer{signer}, store, tTxPool, nil)
	Cons := NewConsensus(testConfig(), st, signer, signer.Address(), make(chan message.Message, 100),
		newMediator())
	cons := Cons.(*consensus)

	testEnterNewHeight(cons)

	cons.signAddVote(vote.VoteTypePrepare, hash.GenerateTestHash())
}

func TestRoundVotes(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t) // height 1
	testEnterNewHeight(tConsP)

	t.Run("Ignore votes from invalid height", func(t *testing.T) {
		v1 := testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, hash.GenerateTestHash(), tIndexX)
		v2 := testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.GenerateTestHash(), tIndexX)
		v3 := testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.GenerateTestHash(), tIndexY)
		v4 := testAddVote(tConsP, vote.VoteTypeChangeProposer, 3, 0, hash.GenerateTestHash(), tIndexX)

		require.False(t, tConsP.log.HasVote(v1.Hash()))
		require.True(t, tConsP.log.HasVote(v2.Hash()))
		require.True(t, tConsP.log.HasVote(v3.Hash()))
		require.False(t, tConsP.log.HasVote(v4.Hash()))
	})
}

func TestConsensusAddVotesNormal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t) // height 1

	testEnterNewHeight(tConsX)
	checkHeightRound(t, tConsX, 2, 0)

	p := makeProposal(t, 2, 0)
	tConsX.SetProposal(p)

	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsX, vote.VoteTypePrepare, 2, 0, p.Block().Hash(), tIndexP)
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p.Block().Hash())

	testAddVote(tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexY)
	testAddVote(tConsX, vote.VoteTypePrecommit, 2, 0, p.Block().Hash(), tIndexP)
	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p.Block().Hash())
	shouldPublishBlockAnnounce(t, tConsX, p.Block().Hash())
}

func TestConsensusAddVote(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsP)

	v1 := testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, hash.GenerateTestHash(), tIndexX)
	v2 := testAddVote(tConsP, vote.VoteTypePrepare, 1, 0, hash.GenerateTestHash(), tIndexX)
	v3 := testAddVote(tConsP, vote.VoteTypePrecommit, 1, 0, hash.GenerateTestHash(), tIndexX)
	v4 := testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, hash.GenerateTestHash(), tIndexX)
	v5 := testAddVote(tConsP, vote.VoteTypePrepare, 1, 2, hash.GenerateTestHash(), tIndexX)

	assert.False(t, tConsP.log.HasVote(v1.Hash())) // invalid height
	assert.True(t, tConsP.log.HasVote(v2.Hash()))
	assert.True(t, tConsP.log.HasVote(v3.Hash()))
	assert.True(t, tConsP.log.HasVote(v4.Hash()))
	assert.True(t, tConsP.log.HasVote(v5.Hash())) // next round
}

func TestConsensusLateProposal1(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t) // height 1

	testEnterNewHeight(tConsB)

	h := uint32(2)
	r := int16(0)
	p := makeProposal(t, h, r)
	require.NotNil(t, p)

	// Partitioned node doesn't receive all the votes
	testAddVote(tConsB, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	testAddVote(tConsB, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	testAddVote(tConsB, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexP)

	testAddVote(tConsB, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(tConsB, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	testAddVote(tConsB, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexP)

	// Partitioned node receives proposal now
	tConsB.SetProposal(p)
	shouldPublishBlockAnnounce(t, tConsB, p.Block().Hash())
}

func TestConsensusInvalidVote(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsX)

	v1, _ := vote.GenerateTestPrecommitVote(1, 0)
	v2 := vote.NewVote(vote.VoteTypePrepare, 2, 0, hash.GenerateTestHash(),
		tSigners[tIndexB].Address())
	tSigners[tIndexB].SignMsg(v2)

	tConsX.AddVote(v1)
	tConsX.AddVote(v2)
	assert.False(t, tConsX.log.HasVote(v1.Hash()))
	assert.False(t, tConsX.log.HasVote(v2.Hash()))
}

func TestPickRandomVote(t *testing.T) {
	setup(t)

	assert.Nil(t, tConsP.PickRandomVote())

	testEnterNewHeight(tConsP)
	assert.Nil(t, tConsP.PickRandomVote())

	p1 := makeProposal(t, 1, 0)
	p2 := makeProposal(t, 1, 1)

	// round 0
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexY)
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexP)
	testAddVote(tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexY)

	assert.NotNil(t, tConsP.PickRandomVote())

	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexY)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexP)

	// Round 1
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexY)
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexP)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 1, hash.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 1, hash.UndefHash, tIndexY)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 1, hash.UndefHash, tIndexP)

	// Round 2
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 2, hash.UndefHash, tIndexP)

	for i := 0; i < 10; i++ {
		rndVote := tConsP.PickRandomVote()
		assert.NotNil(t, rndVote)
		assert.Equal(t, rndVote.Type(), vote.VoteTypeChangeProposer,
			"Should only pick Change Proposer votes")
	}
}

func TestSetProposalFromPreviousRound(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)
	testEnterNewHeight(tConsP)
	testEnterNextRound(tConsP)

	// It should ignore proposal for previous rounds
	tConsP.SetProposal(p)

	assert.Nil(t, tConsP.RoundProposal(0))
	checkHeightRoundWait(t, tConsP, 1, 1)
}

func TestSetProposalFromPreviousHeight(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)
	commitBlockForAllStates(t) // height 1

	testEnterNewHeight(tConsP)

	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(0))
	checkHeightRoundWait(t, tConsP, 2, 0)
}

func TestDuplicateProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	h := uint32(4)
	r := int16(0)
	p1 := makeProposal(t, h, r)
	trx := tx.NewSendTx(hash.UndefHash.Stamp(), 1, tSigners[0].Address(),
		tSigners[1].Address(), 1000, 1000, "proposal changer")
	tSigners[0].SignMsg(trx)
	assert.NoError(t, tTxPool.AppendTx(trx))
	p2 := makeProposal(t, h, r)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	tConsX.SetProposal(p1)
	tConsX.SetProposal(p2)

	assert.Equal(t, tConsX.RoundProposal(0).Hash(), p1.Hash())
}

func TestNonActiveValidator(t *testing.T) {
	setup(t)

	signer := bls.GenerateTestSigner()
	Cons := NewConsensus(testConfig(), state.MockingState(), signer, signer.Address(), make(chan message.Message, 100),
		newMediator())
	nonActiveCons := Cons.(*consensus)

	t.Run("non-active instances should be in new-height state", func(t *testing.T) {
		nonActiveCons.MoveToNewHeight()
		checkHeightRoundWait(t, nonActiveCons, 1, 0)

		assert.False(t, nonActiveCons.IsActive())
		assert.Equal(t, nonActiveCons.currentState.name(), "new-height")
	})

	t.Run("non-active instances should ignore proposals", func(t *testing.T) {
		p := makeProposal(t, 1, 0)
		nonActiveCons.SetProposal(p)

		assert.False(t, nonActiveCons.log.HasRoundProposal(0))
	})

	t.Run("non-active instances should ignore votes", func(t *testing.T) {
		v := testAddVote(nonActiveCons, vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, tIndexX)

		assert.False(t, nonActiveCons.log.HasVote(v.Hash()))
	})

	t.Run("non-active instances should move to new height", func(t *testing.T) {
		b1, cert1 := commitBlockForAllStates(t)
		b2, cert2 := commitBlockForAllStates(t)

		nonActiveCons.MoveToNewHeight()
		checkHeightRoundWait(t, nonActiveCons, 1, 0)

		assert.NoError(t, nonActiveCons.state.CommitBlock(1, b1, cert1))
		assert.NoError(t, nonActiveCons.state.CommitBlock(2, b2, cert2))

		nonActiveCons.MoveToNewHeight()
		checkHeightRoundWait(t, nonActiveCons, 3, 0)
	})
}
