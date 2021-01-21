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
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
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
	if tConsX != nil {
		tConsX.state.Close()
		tConsY.state.Close()
		tConsB.state.Close()
		tConsP.state.Close()
	}
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
	acc.AddToBalance(2100000000000000)
	params := param.MainnetParams()
	params.BlockTimeInSecond = 1
	params.MaximumPower = 4

	tGenDoc = genesis.MakeGenesis("test", util.Now(), []*account.Account{acc}, vals, params)
	stX, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexX], tTxPool)
	require.NoError(t, err)
	stY, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexY], tTxPool)
	require.NoError(t, err)
	stB, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexB], tTxPool)
	require.NoError(t, err)
	stP, err := state.LoadOrNewState(state.TestConfig(), tGenDoc, tSigners[tIndexP], tTxPool)
	require.NoError(t, err)

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

func shouldPublishBlockAnnounce(t *testing.T, cons *consensus, hash crypto.Hash) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishBlockAnnounce", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeBlockAnnounce {
				pld := msg.Payload.(*payload.BlockAnnouncePayload)
				assert.Equal(t, pld.Block.Hash(), hash)
				return
			}
		}
	}
}

func shouldPublishProposal(t *testing.T, cons *consensus, hash crypto.Hash) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishProposal", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeProposal {
				pld := msg.Payload.(*payload.ProposalPayload)
				assert.Equal(t, pld.Proposal.Hash(), hash)
				return
			}
		}
	}
}

func shouldPublishQueryProposal(t *testing.T, cons *consensus, height, round int) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishQueryProposal", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeQueryProposal {
				pld := msg.Payload.(*payload.QueryProposalPayload)
				assert.Equal(t, pld.Height, height)
				assert.Equal(t, pld.Round, round)
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
		time.Sleep(200 * time.Millisecond)
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

func commitBlockForAllStates(t *testing.T) {
	height := tConsX.state.LastBlockHeight()
	var err error
	p := makeProposal(t, height+1, 0)

	sb := block.CommitSignBytes(p.Block().Hash(), 0)
	sig1 := tSigners[0].Sign(sb)
	sig2 := tSigners[1].Sign(sb)
	sig3 := tSigners[2].Sign(sb)
	sig4 := tSigners[3].Sign(sb)

	sig := crypto.Aggregate([]*crypto.Signature{sig1, sig2, sig3, sig4})
	c := block.NewCommit(p.Block().Hash(), 0, []block.Committer{
		{Number: 0, Status: 1},
		{Number: 1, Status: 1},
		{Number: 2, Status: 1},
		{Number: 3, Status: 1},
	}, sig)

	require.NotNil(t, c)
	err = tConsX.state.ApplyBlock(height+1, p.Block(), *c)
	assert.NoError(t, err)
	err = tConsY.state.ApplyBlock(height+1, p.Block(), *c)
	assert.NoError(t, err)
	err = tConsB.state.ApplyBlock(height+1, p.Block(), *c)
	assert.NoError(t, err)
	err = tConsP.state.ApplyBlock(height+1, p.Block(), *c)
	assert.NoError(t, err)
}

func makeProposal(t *testing.T, height, round int) *vote.Proposal {
	var p *vote.Proposal
	switch (height % 4) + round {
	case 1:
		pb, err := tConsX.state.ProposeBlock(round)
		require.NoError(t, err)
		p = vote.NewProposal(height, round, *pb)
		tConsX.signer.SignMsg(p)
	case 2:
		pb, err := tConsY.state.ProposeBlock(round)
		require.NoError(t, err)
		p = vote.NewProposal(height, round, *pb)
		tConsY.signer.SignMsg(p)
	case 3:
		pb, err := tConsB.state.ProposeBlock(round)
		require.NoError(t, err)
		p = vote.NewProposal(height, round, *pb)
		tConsB.signer.SignMsg(p)
	case 0:
		pb, err := tConsP.state.ProposeBlock(round)
		require.NoError(t, err)
		p = vote.NewProposal(height, round, *pb)
		tConsP.signer.SignMsg(p)
	}

	return p
}

func TestHandleTimeout(t *testing.T) {
	setup(t)

	tConsX.handleTimeout(timeout{Height: 1})
	checkHRS(t, tConsX, 0, 0, 0)
}

func TestNotInValidatorSet(t *testing.T) {
	setup(t)

	_, _, priv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv)
	st, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, signer, tTxPool)
	cons, err := NewConsensus(TestConfig(), st, signer, make(chan *message.Message, 100))
	assert.NoError(t, err)

	cons.(*consensus).enterNewHeight()

	cons.(*consensus).signAddVote(vote.VoteTypePrepare, crypto.GenerateTestHash())
	assert.Zero(t, len(cons.RoundVotes(0)))
}

func TestRoundVotes(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t) // height 1
	tConsP.enterNewHeight()

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
	shouldPublishBlockAnnounce(t, tConsX, p.Block().Hash())
}

func TestConsensusUpdateVote(t *testing.T) {
	setup(t)

	tConsY.enterNewHeight()

	h1 := crypto.GenerateTestHash()
	assert.Nil(t, tConsY.LastProposal())

	// Ignore votes from invalid height
	testAddVote(t, tConsY, vote.VoteTypePrepare, 2, 0, h1, tIndexB, false)

	v1 := testAddVote(t, tConsY, vote.VoteTypePrepare, 1, 0, crypto.UndefHash, tIndexX, false)
	v2 := testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, tIndexP, false)
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

	tConsB.SetProposal(p)

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

func TestConsensusInvalidVote(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()

	v, _ := vote.GenerateTestPrecommitVote(1, 0)
	assert.Error(t, tConsX.addVote(v))
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

func TestPickRandomVote(t *testing.T) {
	setup(t)

	tConsY.enterNewHeight()
	assert.Nil(t, tConsY.PickRandomVote())

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), tIndexY, false)
	assert.NotNil(t, tConsY.PickRandomVote())
}
