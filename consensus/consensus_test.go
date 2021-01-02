package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var (
	tSigners []crypto.Signer
	tTxPool  *txpool.MockTxPool
	tGenDoc  *genesis.Genesis
	tConsX   *consensus
	tConsY   *consensus
	tConsB   *consensus // Byzantine
	tConsP   *consensus // partitioned
)

const (
	tIndexX = 0
	tIndexY = 1
	tIndexB = 2
	tIndexP = 3
)

func setup(t *testing.T) {
	conf := logger.TestConfig()
	conf.Levels["_state"] = "debug"
	logger.InitLogger(conf)

	_, keys := validator.GenerateTestValidatorSet()
	tTxPool = txpool.MockingTxPool()

	tSigners = make([]crypto.Signer, 4)
	for i, k := range keys {
		tSigners[i] = crypto.NewSigner(k)
	}

	vals := make([]*validator.Validator, 4)
	for i, s := range tSigners {
		val := validator.NewValidator(s.PublicKey(), 0, i)
		vals[i] = val
	}

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21000000000000)

	tGenDoc = genesis.MakeGenesis("test", util.Now(), []*account.Account{acc}, vals, 1)
	stX, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexX], tTxPool)
	stY, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexY], tTxPool)
	stB, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexB], tTxPool)
	stP, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexP], tTxPool)

	consX, err := NewConsensus(TestConfig(), stX, tSigners[tIndexX], make(chan *message.Message, 100))
	assert.NoError(t, err)
	consY, err := NewConsensus(TestConfig(), stY, tSigners[tIndexY], make(chan *message.Message, 100))
	assert.NoError(t, err)
	consB, err := NewConsensus(TestConfig(), stB, tSigners[tIndexB], make(chan *message.Message, 100))
	assert.NoError(t, err)
	consP, err := NewConsensus(TestConfig(), stP, tSigners[tIndexP], make(chan *message.Message, 100))
	assert.NoError(t, err)
	tConsX = consX.(*consensus)
	tConsY = consY.(*consensus)
	tConsB = consB.(*consensus)
	tConsP = consP.(*consensus)
}

func shouldPublishProposalBlock(t *testing.T, cons *consensus) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishProposalBlock", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeLatestBlocks {
				return
			}
		}
	}
}
func shouldPublishProposalReqquest(t *testing.T, cons *consensus) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishProposalReqquest", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeProposalRequest {
				return
			}
		}
	}
}

func shouldPublishVote(t *testing.T, cons *consensus, voteType vote.VoteType, hash crypto.Hash) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishUndefVote", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeVote {
				pld := msg.Payload.(*payload.VotePayload)
				assert.Equal(t, pld.Vote.VoteType(), voteType)
				assert.Equal(t, pld.Vote.BlockHash(), hash)
				return
			}
		}
	}
}

func checkHRS(t *testing.T, cons *consensus, height, round int, step hrs.StepType) {
	assert.Equal(t, hrs.NewHRS(height, round, step), cons.hrs)
}

func checkHRSWait(t *testing.T, cons *consensus, height, round int, step hrs.StepType) {
	expected := hrs.NewHRS(height, round, step)
	for i := 0; i < 20; i++ {
		if expected.EqualsTo(cons.HRS()) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	assert.Equal(t, expected, cons.hrs)
}

func testAddVote(t *testing.T,
	cons *consensus,
	voteType vote.VoteType,
	height int,
	round int,
	blockHash crypto.Hash,
	valID int,
	expectError bool) *vote.Vote {

	v := vote.NewVote(voteType, height, round, blockHash, tSigners[valID].Address())
	tSigners[valID].SignMsg(v)

	if expectError {
		assert.Error(t, cons.addVote(v))
	} else {
		assert.NoError(t, cons.addVote(v))
	}
	return v
}

func TestNotInValidatorSet(t *testing.T) {
	setup(t)

	_, _, priv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv)
	st, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, signer, tTxPool)
	cons, err := NewConsensus(TestConfig(), st, signer, make(chan *message.Message, 100))
	assert.NoError(t, err)

	cons.MoveToNewHeight()

	cons.(*consensus).signAddVote(vote.VoteTypePrepare, crypto.GenerateTestHash())
	assert.Nil(t, cons.(*consensus).pendingVotes.GetRoundVotes(0))
}

func TestRoundVotes(t *testing.T) {
	setup(t)

	tConsY.enterNewHeight()
	checkHRSWait(t, tConsY, 1, 0, hrs.StepTypePrepare)

	v1 := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.GenerateTestHash(), tSigners[tIndexY].Address())
	tSigners[tIndexY].SignMsg(v1)
	v2 := vote.NewVote(vote.VoteTypePrepare, 1, 1, crypto.GenerateTestHash(), tSigners[tIndexB].Address())
	tSigners[tIndexB].SignMsg(v2)
	v3 := vote.NewVote(vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tSigners[tIndexP].Address())
	tSigners[tIndexP].SignMsg(v3)
	tConsY.AddVote(v1)
	tConsY.AddVote(v2)
	tConsY.AddVote(v3)

	require.True(t, tConsY.HasVote(v1.Hash()))
}

func TestConsensusAddVotesNormal(t *testing.T) {
	setup(t)

	tConsX.MoveToNewHeight()
	checkHRSWait(t, tConsX, 1, 0, hrs.StepTypePrepare)

	p := tConsX.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, tConsX, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), tIndexY, false)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, tConsX, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), tIndexP, false)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY, false)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP, false)
	checkHRSWait(t, tConsX, 2, 0, hrs.StepTypePropose)
}

func TestConsensusUpdateVote(t *testing.T) {
	setup(t)

	tConsY.enterNewHeight()

	h1 := crypto.GenerateTestHash()
	assert.Nil(t, tConsY.LastProposal())

	// Ignore votes from invalid height
	testAddVote(t, tConsY, vote.VoteTypePrepare, 2, 0, h1, tIndexB, false)

	v1 := testAddVote(t, tConsY, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexX, false)
	v2 := testAddVote(t, tConsY, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexP, false)
	tConsY.enterNewRound(1)
	v3 := testAddVote(t, tConsY, vote.VoteTypePrepare, 1, 1, crypto.UndefHash, tIndexX, false)
	tConsY.enterNewRound(2)
	v4 := testAddVote(t, tConsY, vote.VoteTypePrepare, 1, 2, crypto.UndefHash, tIndexX, false)

	assert.Contains(t, tConsY.RoundVotesHash(0), v1.Hash())
	assert.Contains(t, tConsY.RoundVotesHash(0), v2.Hash())
	assert.Contains(t, tConsY.RoundVotesHash(1), v3.Hash())
	assert.Contains(t, tConsY.RoundVotesHash(2), v4.Hash())
	assert.NotContains(t, tConsY.RoundVotesHash(2), v1.Hash())

	assert.Contains(t, tConsY.RoundVotes(0), v1)
	assert.Contains(t, tConsY.RoundVotes(0), v2)
	assert.Contains(t, tConsY.RoundVotes(1), v3)
	assert.Contains(t, tConsY.RoundVotes(2), v4)
}

func TestConsensusNoPrepares(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	tConsB.enterNewHeight()

	p := tConsX.LastProposal()
	require.NotNil(t, p)

	tConsB.pendingVotes.SetRoundProposal(0, p)

	testAddVote(t, tConsB, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexX, false)
	checkHRSWait(t, tConsB, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, tConsB, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY, false)
	checkHRS(t, tConsB, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, tConsB, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP, false)
	checkHRS(t, tConsB, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, tConsB.isCommitted, true)
	precommits := tConsB.pendingVotes.PrecommitVoteSet(0)
	assert.Equal(t, precommits.Len(), 3)

	// Commit block again
	assert.NoError(t, tConsB.state.ApplyBlock(1, p.Block(), *precommits.ToCommit()))

	// Commit a block for wrong height
	assert.Error(t, tConsB.state.ApplyBlock(5, p.Block(), *precommits.ToCommit()))
}

func TestConsensusSpamming(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()

	for i := 0; i < 100; i++ {
		v, _ := vote.GenerateTestPrecommitVote(1, 0)
		assert.Error(t, tConsX.addVote(v))
	}
}

func TestConsensusSpammingProposal(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()
	assert.Nil(t, tConsP.LastProposal())

	for i := 0; i < 100; i++ {
		proposal, _ := vote.GenerateTestProposal(1, 0)
		tConsP.SetProposal(proposal)
	}
	assert.Nil(t, tConsP.LastProposal())
}

func TestConsensusInvalidProposal(t *testing.T) {
	setup(t)

	tConsY.enterNewHeight()
	assert.Nil(t, tConsY.LastProposal())

	addr := tSigners[tIndexX].Address()
	b, _ := block.GenerateTestBlock(&addr, nil)
	p := vote.NewProposal(1, 0, *b)

	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.LastProposal())

	tSigners[tIndexY].SignMsg(p)
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.LastProposal())

}

func TestConsensusFingerprint(t *testing.T) {
	setup(t)

	assert.Contains(t, tConsX.Fingerprint(), tConsX.hrs.String())
}

func TestStop(t *testing.T) {
	setup(t)

	tConsX.Stop()
	tConsX.handleTimeout(timeout{1 * time.Second, 1, 0, 0})
	assert.Equal(t, tConsX.hrs.Height(), -1)
}
