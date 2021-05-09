package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var (
	tSigners []crypto.Signer
	tTxPool  *txpool.MockTxPool
	tGenDoc  *genesis.Genesis
	tConsX   *consensus
	tConsY   *consensus
	tConsB   *consensus // Byzantine of offline
	tConsP   *consensus // partitioned
)

const (
	tIndexX = 0
	tIndexY = 1
	tIndexB = 2
	tIndexP = 3
)

type OverrideFingerprint struct {
	cons Consensus
	name string
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.cons.Fingerprint()
}

func setup(t *testing.T) {
	conf := logger.TestConfig()
	conf.Levels["_consensus"] = "debug"
	logger.InitLogger(conf)

	_, tSigners = committee.GenerateTestCommittee()
	tTxPool = txpool.MockingTxPool()

	vals := make([]*validator.Validator, 4)
	for i, s := range tSigners {
		val := validator.NewValidator(s.PublicKey(), 0, i)
		vals[i] = val
	}

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)
	params := param.DefaultParams()
	params.CommitteeSize = 4
	params.BlockTimeInSecond = 2

	store1 := store.MockingStore()
	store2 := store.MockingStore()
	store3 := store.MockingStore()
	store4 := store.MockingStore()

	// To prevent trigging timers before starting the tests, otherwise some tests will have double entry for new height.
	getTime := util.RoundNow(params.BlockTimeInSecond).Add(time.Duration(params.BlockTimeInSecond) * time.Second)
	tGenDoc = genesis.MakeGenesis(getTime, []*account.Account{acc}, vals, params)
	stX, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexX], store1, tTxPool)
	require.NoError(t, err)
	stY, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexY], store2, tTxPool)
	require.NoError(t, err)
	stB, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexB], store3, tTxPool)
	require.NoError(t, err)
	stP, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexP], store4, tTxPool)
	require.NoError(t, err)

	consX, err := NewConsensus(TestConfig(), stX, tSigners[tIndexX], make(chan payload.Payload, 100))
	assert.NoError(t, err)
	consY, err := NewConsensus(TestConfig(), stY, tSigners[tIndexY], make(chan payload.Payload, 100))
	assert.NoError(t, err)
	consB, err := NewConsensus(TestConfig(), stB, tSigners[tIndexB], make(chan payload.Payload, 100))
	assert.NoError(t, err)
	consP, err := NewConsensus(TestConfig(), stP, tSigners[tIndexP], make(chan payload.Payload, 100))
	assert.NoError(t, err)
	tConsX = consX.(*consensus)
	tConsY = consY.(*consensus)
	tConsB = consB.(*consensus)
	tConsP = consP.(*consensus)

	tConsX.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: fmt.Sprintf("consX - %s: ", t.Name()), cons: tConsX})
	tConsY.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: fmt.Sprintf("consY - %s: ", t.Name()), cons: tConsY})
	tConsB.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: fmt.Sprintf("consB - %s: ", t.Name()), cons: tConsB})
	tConsP.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: fmt.Sprintf("consP - %s: ", t.Name()), cons: tConsP})

	logger.Info("Setup finished, start running the test", "name", t.Name())
}

func shouldPublishBlockAnnounce(t *testing.T, cons *consensus, hash crypto.Hash) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case pld := <-cons.broadcastCh:
			logger.Info("shouldPublishBlockAnnounce", "pld", pld)

			if pld.Type() == payload.PayloadTypeBlockAnnounce {
				p := pld.(*payload.BlockAnnouncePayload)
				assert.Equal(t, p.Block.Hash(), hash)
				return
			}
		}
	}
}

func shouldPublishProposal(t *testing.T, cons *consensus, height, round int) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case pld := <-cons.broadcastCh:
			logger.Info("shouldPublishProposal", "pld", pld)

			if pld.Type() == payload.PayloadTypeProposal {
				p := pld.(*payload.ProposalPayload)
				assert.Equal(t, p.Proposal.Height(), height)
				assert.Equal(t, p.Proposal.Round(), round)
				return
			}
		}
	}
}

func shouldPublishQueryProposal(t *testing.T, cons *consensus, height, round int) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case pld := <-cons.broadcastCh:
			logger.Info("shouldPublishQueryProposal", "pld", pld)

			if pld.Type() == payload.PayloadTypeQueryProposal {
				p := pld.(*payload.QueryProposalPayload)
				assert.Equal(t, p.Height, height)
				assert.Equal(t, p.Round, round)
				return
			}
		}
	}
}

func shouldPublishVote(t *testing.T, cons *consensus, voteType vote.Type, hash crypto.Hash) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
		case pld := <-cons.broadcastCh:
			logger.Info("shouldPublishVote", "pld", pld)

			if pld.Type() == payload.PayloadTypeVote {
				p := pld.(*payload.VotePayload)
				if p.Vote.Type() == voteType &&
					p.Vote.BlockHash().EqualsTo(hash) {
					return
				}
			}
		}
	}
}

func checkHeightRound(t *testing.T, cons *consensus, height, round int) {
	assert.Equal(t, cons.Height(), height)
	assert.Equal(t, cons.Round(), round)
}

func checkHeightRoundWait(t *testing.T, cons *consensus, height, round int) {
	for i := 0; i < 20; i++ {
		if cons.Height() == height && cons.Round() == round {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	checkHeightRound(t, cons, height, round)
}

func testAddVote(cons *consensus,
	voteType vote.Type,
	height int,
	round int,
	blockHash crypto.Hash,
	valID int) *vote.Vote {

	v := vote.NewVote(voteType, height, round, blockHash, tSigners[valID].Address())
	tSigners[valID].SignMsg(v)

	cons.AddVote(v)

	return v
}

// testEnterNewHeight helps tests to enter new height safely
// without scheduling new height. It boosts the test speed
func testEnterNewHeight(cons *consensus) {
	cons.lk.Lock()
	cons.enterNewState(cons.newHeightState)
	cons.currentState.onTimedout(&ticker{0, cons.height, cons.round, tickerTargetNewHeight})
	cons.lk.Unlock()
}

// testEnterNextRound helps tests to enter next round safely
func testEnterNextRound(cons *consensus) {
	cons.lk.Lock()
	cons.round++
	cons.enterNewState(cons.proposeState)
	cons.lk.Unlock()
}

func commitBlockForAllStates(t *testing.T) {
	height := tConsX.state.LastBlockHeight()
	var err error
	p := makeProposal(t, height+1, 0)

	sb := block.CertificateSignBytes(p.Block().Hash(), 0)
	sig1 := tSigners[0].SignData(sb)
	sig2 := tSigners[1].SignData(sb)
	sig4 := tSigners[3].SignData(sb)

	sig := crypto.Aggregate([]crypto.Signature{sig1, sig2, sig4})
	cert := block.NewCertificate(p.Block().Hash(), 0, []int{0, 1, 2, 3}, []int{2}, sig)

	require.NotNil(t, cert)
	err = tConsX.state.CommitBlock(height+1, p.Block(), cert)
	assert.NoError(t, err)
	err = tConsY.state.CommitBlock(height+1, p.Block(), cert)
	assert.NoError(t, err)
	err = tConsB.state.CommitBlock(height+1, p.Block(), cert)
	assert.NoError(t, err)
	err = tConsP.state.CommitBlock(height+1, p.Block(), cert)
	assert.NoError(t, err)
}

func makeProposal(t *testing.T, height, round int) *proposal.Proposal {
	var p *proposal.Proposal
	switch (height % 4) + round {
	case 1:
		pb, err := tConsX.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, pb)
		tConsX.signer.SignMsg(p)
	case 2:
		pb, err := tConsY.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, pb)
		tConsY.signer.SignMsg(p)
	case 3:
		pb, err := tConsB.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, pb)
		tConsB.signer.SignMsg(p)
	case 0, 4:
		pb, err := tConsP.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, pb)
		tConsP.signer.SignMsg(p)
	}

	return p
}

func TestNotInCommittee(t *testing.T) {
	setup(t)

	_, _, priv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv)
	store := store.MockingStore()

	st, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, signer, store, tTxPool)
	cons, err := NewConsensus(TestConfig(), st, signer, make(chan payload.Payload, 100))
	assert.NoError(t, err)

	testEnterNewHeight(cons.(*consensus))

	cons.(*consensus).signAddVote(vote.VoteTypePrepare, crypto.GenerateTestHash())
	assert.Zero(t, len(cons.RoundVotes(0)))
}

func TestRoundVotes(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t) // height 1
	testEnterNewHeight(tConsP)

	t.Run("Ignore votes from invalid height", func(t *testing.T) {
		v1 := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.GenerateTestHash(), tSigners[tIndexX].Address())
		tSigners[tIndexX].SignMsg(v1)

		v2 := vote.NewVote(vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tSigners[tIndexX].Address())
		tSigners[tIndexX].SignMsg(v2)

		v3 := vote.NewVote(vote.VoteTypePrepare, 3, 0, crypto.GenerateTestHash(), tSigners[tIndexX].Address())
		tSigners[tIndexX].SignMsg(v3)

		tConsP.AddVote(v1)
		tConsP.AddVote(v2)
		tConsP.AddVote(v3)

		require.False(t, tConsP.HasVote(v1.Hash()))
		require.True(t, tConsP.HasVote(v2.Hash()))
		require.False(t, tConsP.HasVote(v3.Hash()))
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

	v1 := testAddVote(tConsP, vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tIndexX)
	v2 := testAddVote(tConsP, vote.VoteTypePrepare, 1, 0, crypto.GenerateTestHash(), tIndexX)
	v3 := testAddVote(tConsP, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), tIndexX)
	v4 := testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, crypto.GenerateTestHash(), tIndexX)
	v5 := testAddVote(tConsP, vote.VoteTypePrepare, 1, 2, crypto.GenerateTestHash(), tIndexX)

	assert.False(t, tConsP.HasVote(v1.Hash())) // invalid height
	assert.True(t, tConsP.HasVote(v2.Hash()))
	assert.True(t, tConsP.HasVote(v3.Hash()))
	assert.True(t, tConsP.HasVote(v4.Hash()))
	assert.True(t, tConsP.HasVote(v5.Hash())) // next round
}

func TestConsensusLateProposal1(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t) // height 1

	testEnterNewHeight(tConsB)

	h := 2
	r := 0
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

	v, _ := vote.GenerateTestPrecommitVote(1, 0)
	tConsX.AddVote(v)
	assert.False(t, tConsX.HasVote(v.Hash()))
}

func TestPickRandomVote(t *testing.T) {
	setup(t)

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

	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, crypto.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, crypto.UndefHash, tIndexY)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 0, crypto.UndefHash, tIndexP)

	// Round 1
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexY)
	testAddVote(tConsP, vote.VoteTypePrepare, 1, 1, p2.Block().Hash(), tIndexP)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 1, crypto.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 1, crypto.UndefHash, tIndexY)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 1, crypto.UndefHash, tIndexP)

	// Round 2
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 1, 2, crypto.UndefHash, tIndexP)

	for i := 0; i < 10; i++ {
		rndVote := tConsP.PickRandomVote()
		assert.NotNil(t, rndVote)
		assert.Equal(t, rndVote.Type(), vote.VoteTypeChangeProposer, "Should only pick Change Proposer votes")
	}
}

func TestSetProposalFromPreviousRound(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)
	testEnterNewHeight(tConsP)
	testEnterNextRound(tConsP)

	// Keep proposal for previous round, but don't change the state
	tConsP.SetProposal(p)

	assert.NotNil(t, tConsP.RoundProposal(0), 0)
	checkHeightRoundWait(t, tConsP, 1, 1)
}

func TestSetProposalFromPreviousHeight(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)
	commitBlockForAllStates(t) // height 1

	testEnterNewHeight(tConsP)

	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(0), 0)
	checkHeightRoundWait(t, tConsP, 2, 0)
}

func TestDuplicateProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	h := 4
	r := 0
	p1 := makeProposal(t, h, r)
	trx := tx.NewSendTx(crypto.UndefHash, 1, tSigners[0].Address(), tSigners[1].Address(), 1000, 1000, "proposal changer")
	tSigners[0].SignMsg(trx)
	assert.NoError(t, tTxPool.AppendTx(trx))
	p2 := makeProposal(t, h, r)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	tConsX.SetProposal(p1)
	tConsX.SetProposal(p2)

	assert.Equal(t, tConsX.RoundProposal(0).Hash(), p1.Hash())
}
