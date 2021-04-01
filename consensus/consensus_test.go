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
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/proposal"
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

type OverrideFingerprint struct {
	cons Consensus
	name string
}

func (o *OverrideFingerprint) Fingerprint() string {
	return o.name + o.cons.Fingerprint()
}

func setup(t *testing.T) {
	if tConsX != nil {
		tConsX.state.Close()
		tConsY.state.Close()
		tConsB.state.Close()
		tConsP.state.Close()
	}
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

	tGenDoc = genesis.MakeGenesis(util.Now(), []*account.Account{acc}, vals, params)
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

	tConsX.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: "consX: ", cons: tConsX})
	tConsY.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: "consY: ", cons: tConsY})
	tConsB.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: "consB: ", cons: tConsB})
	tConsP.logger = logger.NewLogger("_consensus", &OverrideFingerprint{name: "consP: ", cons: tConsP})
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

func shouldPublishVote(t *testing.T, cons *consensus, voteType vote.VoteType, hash crypto.Hash) *vote.Vote {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishVote", "msg", msg)

			if msg.PayloadType() == payload.PayloadTypeVote {
				pld := msg.Payload.(*payload.VotePayload)
				if pld.Vote.BlockHash().EqualsTo(hash) {
					assert.Equal(t, pld.Vote.VoteType(), voteType)
					return pld.Vote
				}
			}
		}
	}
}

func checkHRS(t *testing.T, cons *consensus, height, round int, step hrs.StepType) {
	expected := hrs.NewHRS(height, round, step)
	assert.True(t, expected.EqualsTo(cons.HRS()))
}

func checkHRSWait(t *testing.T, cons *consensus, height, round int, step hrs.StepType) {
	expected := hrs.NewHRS(height, round, step)
	for i := 0; i < 20; i++ {
		if expected.EqualsTo(cons.HRS()) {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	assert.True(t, expected.EqualsTo(cons.HRS()))
}

func testAddVote(t *testing.T,
	cons *consensus,
	voteType vote.VoteType,
	height int,
	round int,
	blockHash crypto.Hash,
	valID int) *vote.Vote {

	v := vote.NewVote(voteType, height, round, blockHash, tSigners[valID].Address())
	tSigners[valID].SignMsg(v)

	cons.AddVote(v)

	return v
}

func commitBlockForAllStates(t *testing.T) {
	height := tConsX.state.LastBlockHeight()
	var err error
	p := makeProposal(t, height+1, 0)

	sb := block.CertificateSignBytes(p.Block().Hash(), 0)
	sig1 := tSigners[0].SignData(sb)
	sig2 := tSigners[1].SignData(sb)
	sig3 := tSigners[2].SignData(sb)
	sig4 := tSigners[3].SignData(sb)

	sig := crypto.Aggregate([]crypto.Signature{sig1, sig2, sig3, sig4})
	cert := block.NewCertificate(p.Block().Hash(), 0, []int{0, 1, 2, 3}, []int{}, sig)

	require.NotNil(t, cert)
	err = tConsX.state.CommitBlock(height+1, p.Block(), *cert)
	assert.NoError(t, err)
	err = tConsY.state.CommitBlock(height+1, p.Block(), *cert)
	assert.NoError(t, err)
	err = tConsB.state.CommitBlock(height+1, p.Block(), *cert)
	assert.NoError(t, err)
	err = tConsP.state.CommitBlock(height+1, p.Block(), *cert)
	assert.NoError(t, err)
}

func makeProposal(t *testing.T, height, round int) *proposal.Proposal {
	var p *proposal.Proposal
	switch (height % 4) + round {
	case 1:
		pb, err := tConsX.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, *pb)
		tConsX.signer.SignMsg(p)
	case 2:
		pb, err := tConsY.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, *pb)
		tConsY.signer.SignMsg(p)
	case 3:
		pb, err := tConsB.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, *pb)
		tConsB.signer.SignMsg(p)
	case 0, 4:
		pb, err := tConsP.state.ProposeBlock(round)
		require.NoError(t, err)
		p = proposal.NewProposal(height, round, *pb)
		tConsP.signer.SignMsg(p)
	}

	return p
}

func TestHandleTimeout(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	tConsX.hrs = hrs.NewHRS(2, 0, hrs.StepTypeNewHeight)

	tConsX.handleTimeout(timeout{Height: 1, Step: hrs.StepTypePrepare})
	checkHRS(t, tConsX, 2, 0, hrs.StepTypeNewHeight)

	tConsX.handleTimeout(timeout{Height: 2, Step: hrs.StepTypePrepare})
	checkHRS(t, tConsX, 2, 0, hrs.StepTypePrepare)
}

func TestNotInCommittee(t *testing.T) {
	setup(t)

	_, _, priv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv)
	st, _ := state.LoadOrNewState(state.TestConfig(), tGenDoc, signer, tTxPool)
	cons, err := NewConsensus(TestConfig(), st, signer, make(chan *message.Message, 100))
	assert.NoError(t, err)

	cons.(*consensus).enterNewHeight()

	cons.(*consensus).signAddVote(vote.VoteTypePrepare, 0, crypto.GenerateTestHash())
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

	tConsX.enterNewHeight()
	checkHRSWait(t, tConsX, 1, 0, hrs.StepTypePrepare)

	p := tConsX.RoundProposal(0)
	require.NotNil(t, p)

	testAddVote(t, tConsX, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), tIndexY)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrepare)

	testAddVote(t, tConsX, vote.VoteTypePrepare, 1, 0, p.Block().Hash(), tIndexP)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexY)
	checkHRS(t, tConsX, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), tIndexP)
	shouldPublishBlockAnnounce(t, tConsX, p.Block().Hash())
}

func TestConsensusAddVote(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()

	v1 := testAddVote(t, tConsP, vote.VoteTypePrepare, 2, 0, crypto.GenerateTestHash(), tIndexX)
	v2 := testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, crypto.GenerateTestHash(), tIndexX)
	v3 := testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), tIndexX)
	v4 := testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 1, crypto.GenerateTestHash(), tIndexX)
	v5 := testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 2, crypto.GenerateTestHash(), tIndexX)

	assert.False(t, tConsP.HasVote(v1.Hash())) // invalid height
	assert.True(t, tConsP.HasVote(v2.Hash()))
	assert.True(t, tConsP.HasVote(v3.Hash()))
	assert.True(t, tConsP.HasVote(v4.Hash()))
	assert.True(t, tConsP.HasVote(v5.Hash()))
}

func TestConsensusNoPrepares(t *testing.T) {
	setup(t)

	tConsB.enterNewHeight()

	h := 1
	r := 0
	p := makeProposal(t, h, r)
	require.NotNil(t, p)

	tConsB.SetProposal(p)

	testAddVote(t, tConsB, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX)
	testAddVote(t, tConsB, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY)
	checkHRS(t, tConsB, h, r, hrs.StepTypePrepare)

	testAddVote(t, tConsB, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexP)
	checkHRS(t, tConsB, h, r, hrs.StepTypeCommit)

	shouldPublishBlockAnnounce(t, tConsB, p.Block().Hash())
	assert.Equal(t, tConsB.pendingVotes.PrecommitVoteSet(0).Len(), 3)
}

func TestConsensusInvalidVote(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()

	v, _ := vote.GenerateTestPrecommitVote(1, 0)
	assert.Error(t, tConsX.addVote(v))
}

func TestPickRandomVote(t *testing.T) {
	setup(t)

	tConsY.enterNewHeight()
	assert.Nil(t, tConsY.PickRandomVote())

	testAddVote(t, tConsY, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), tIndexY)
	assert.NotNil(t, tConsY.PickRandomVote())
}

func TestSignProposalFromPreviousRound(t *testing.T) {
	setup(t)

	p0 := makeProposal(t, 1, 0)
	tConsP.enterNewHeight()
	tConsP.enterNewRound(1)

	tConsP.SetProposal(p0)

	v := shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p0.Block().Hash())
	assert.Equal(t, v.Round(), 0)
}
