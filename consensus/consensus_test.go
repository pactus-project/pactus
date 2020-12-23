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
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var (
	tSigners []crypto.Signer
	tTxPool  *txpool.MockTxPool
)

const (
	VAL1 = 0
	VAL2 = 1
	VAL3 = 2
	VAL4 = 3
)

func init() {
	logger.InitLogger(logger.TestConfig())

	_, keys := validator.GenerateTestValidatorSet()
	tTxPool = txpool.NewMockTxPool()

	tSigners = make([]crypto.Signer, 4)
	for i, k := range keys {
		tSigners[i] = crypto.NewSigner(k)
	}
}

func newTestConsensus(t *testing.T, valID int) *consensus {
	vals := make([]*validator.Validator, 4)
	for i, s := range tSigners {
		val := validator.NewValidator(s.PublicKey(), 0, i)
		vals[i] = val
	}

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21000000000000)

	ch := make(chan *message.Message, 100)

	genDoc := genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, vals, 1)
	st, _ := state.LoadOrNewState(state.TestConfig(), genDoc, tSigners[valID], tTxPool)

	cons1, err := NewConsensus(TestConfig(), st, tSigners[valID], ch)
	assert.NoError(t, err)
	cons := cons1.(*consensus)
	assert.Equal(t, cons.votes.height, 0)
	assert.Equal(t, hrs.NewHRS(0, 0, hrs.StepTypeNewHeight), cons.hrs)

	return cons
}

func shouldPublishMessageWithThisType(t *testing.T, cons *consensus, payloadType payload.PayloadType) *message.Message {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return nil
		case msg := <-cons.broadcastCh:
			logger.Info("shouldPublishMessageWithThisType", "msg", msg, "type", payloadType.String())

			if msg.PayloadType() == payloadType {
				return msg
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

func TestAllVoteHashes(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	cons.enterNewHeight(1)

	v1 := vote.NewVote(vote.VoteTypePrevote, 1, 0, crypto.GenerateTestHash(), tSigners[VAL2].Address())
	tSigners[VAL2].SignMsg(v1)
	v2 := vote.NewVote(vote.VoteTypePrevote, 1, 1, crypto.GenerateTestHash(), tSigners[VAL3].Address())
	tSigners[VAL3].SignMsg(v2)
	v3 := vote.NewVote(vote.VoteTypePrevote, 2, 0, crypto.GenerateTestHash(), tSigners[VAL4].Address())
	tSigners[VAL4].SignMsg(v3)
	cons.AddVote(v1)
	cons.AddVote(v2)
	cons.AddVote(v3)

	votes := cons.AllVotes()
	hashes := []crypto.Hash{}
	for _, v := range votes {
		hashes = append(hashes, v.Hash())
	}
	assert.ElementsMatch(t, hashes, cons.AllVotesHashes())
	assert.NotContains(t, hashes, v3.Hash())
	assert.NotNil(t, cons.Vote(v2.Hash()))
}

func TestConsensusAddVotesNormal(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	cons.enterNewHeight(1)

	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL2, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommit)

	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL2, false)
	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
	assert.Equal(t, cons.votes.Precommits(0).Len(), 3) // Votes from validator 1,2,3
}

func TestConsensusUpdateVote(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	cons.enterNewHeight(1)

	p := cons.LastProposal()
	assert.NotNil(t, p)

	assert.Equal(t, cons.votes.Prevotes(0).Len(), 1)
	// Ignore votes from invalid height
	testAddVote(t, cons, vote.VoteTypePrevote, 2, 0, p.Block().Hash(), VAL2, false)
	assert.Equal(t, cons.votes.Prevotes(0).Len(), 1)

	// Validator_2 doesn't have proposal now vote for nil
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL2, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrevoteWait)

	// Validator_2 have proposal now and vote for that
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL2, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommit)

	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL2, false)
	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
	assert.Equal(t, cons.votes.Precommits(0).Len(), 3)
}

func TestConsensusNoPrevotes(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	cons.enterNewHeight(1)
	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL2, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL4, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
	precommits := cons.votes.Precommits(0)
	assert.Equal(t, precommits.Len(), 3)

	// Commit block again
	assert.NoError(t, cons.state.ApplyBlock(1, p.Block(), *precommits.ToCommit()))

	// Commit a block for wrong height
	assert.Error(t, cons.state.ApplyBlock(5, p.Block(), *precommits.ToCommit()))
}

func TestConsensusGotoNextRound(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	cons.enterNewHeight(1)

	// Validator_1 is offline
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL3, false)
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL4, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL3, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL4, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrevote)

	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL1, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrevote)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL1, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
}

func TestConsensusGotoNextRound2(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	cons.enterNewHeight(1)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.GenerateTestHash(), VAL1, false)
	shouldPublishMessageWithThisType(t, cons, payload.PayloadTypeProposalReq)
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL3, false)
	checkHRSWait(t, cons, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), VAL1, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL3, false)
	checkHRSWait(t, cons, 1, 1, hrs.StepTypePrevote)

	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL1, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrevote)

	testAddVote(t, cons, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL1, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL3, false)
	checkHRS(t, cons, 1, 1, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
}
func TestConsensusSpamming(t *testing.T) {
	cons := newTestConsensus(t, VAL1)

	cons.enterNewHeight(1)

	for i := 0; i < 100; i++ {
		v, _ := vote.GenerateTestPrecommitVote(1, 0)
		assert.Error(t, cons.addVote(v))
	}
}

func TestConsensusSpammingProposal(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	cons.enterNewHeight(1)
	p := cons.LastProposal()
	assert.Nil(t, p)

	for i := 0; i < 100; i++ {
		proposal, _ := vote.GenerateTestProposal(1, 0)
		cons.SetProposal(proposal)
	}
	p = cons.LastProposal()
	assert.Nil(t, p)
}

func TestConsensusInvalidProposal(t *testing.T) {
	cons := newTestConsensus(t, VAL2)

	cons.enterNewHeight(1)
	assert.Nil(t, cons.LastProposal())

	addr := tSigners[VAL1].Address()
	b, _ := block.GenerateTestBlock(&addr, nil)
	p := vote.NewProposal(1, 0, *b)

	cons.SetProposal(p)
	assert.Nil(t, cons.LastProposal())

	tSigners[VAL2].SignMsg(p)
	cons.SetProposal(p)
	assert.Nil(t, cons.LastProposal())

}

func TestConsensusFingerprint(t *testing.T) {
	cons := newTestConsensus(t, VAL2)
	assert.Contains(t, cons.Fingerprint(), cons.hrs.Fingerprint())
}
