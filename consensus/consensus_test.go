package consensus

import (
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
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var (
	pvals []*validator.PrivValidator
)

const (
	VAL1 = 0
	VAL2 = 1
	VAL3 = 2
	VAL4 = 3
)

func init() {
	_, keys := validator.GenerateTestValidatorSet()

	pvals = make([]*validator.PrivValidator, 4)
	for i, k := range keys {
		pval := validator.NewPrivValidator(k)
		pvals[i] = pval
	}
}

func newTestConsensus(t *testing.T, valID int) *Consensus {
	consConf := TestConfig()
	stateConf := state.TestConfig()
	txPoolConf := txpool.TestConfig()
	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)

	vals := make([]*validator.Validator, 4)
	for i, pval := range pvals {
		val := validator.NewValidator(pval.PublicKey(), 0)
		val.AddToStake(100)
		vals[i] = val
	}

	acc := account.NewAccount(crypto.MintbaseAddress)
	acc.SetBalance(21000000000000)

	ch := make(chan *message.Message, 10)
	go func() {
		for {
			select {
			case <-ch:
			default:
			}
		}
	}()

	genDoc := genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, vals)
	txpool, _ := txpool.NewTxPool(txPoolConf, ch)
	st, _ := state.LoadOrNewState(stateConf, genDoc, pvals[valID].Address(), txpool)

	cons, _ := NewConsensus(consConf, st, pvals[valID], ch)
	assert.Equal(t, cons.votes.height, 0)
	assert.Equal(t, hrs.NewHRS(0, 0, hrs.StepTypeNewHeight), cons.hrs)

	return cons
}

func checkHRS(t *testing.T, cons *Consensus, height, round int, step hrs.StepType) {
	assert.Equal(t, hrs.NewHRS(height, round, step), cons.hrs)
}

func checkHRSWait(t *testing.T, cons *Consensus, height, round int, step hrs.StepType) {
	expected := hrs.NewHRS(height, round, step)
	for i := 0; i < 100; i++ {
		if expected.EqualsTo(cons.hrs) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	assert.Equal(t, expected, cons.hrs)
}

func testAddVote(t *testing.T,
	cons *Consensus,
	voteType vote.VoteType,
	height int,
	round int,
	blockHash crypto.Hash,
	pvalID int,
	expectError bool) {

	v := vote.NewVote(voteType, height, round, blockHash, pvals[pvalID].Address())
	pvals[pvalID].SignMsg(v)

	if expectError {
		assert.Error(t, cons.AddVote(v))
	} else {
		assert.NoError(t, cons.AddVote(v))
	}
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

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL4, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, cons.votes.Precommits(0).Len(), 4) // Votes from validator 1,2,3,4
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

	// Validator_1 is online, but the proposal is not accepted by other nodes
	// Validator_4 is offline
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.GenerateTestHash(), VAL1, false)
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrevoteWait)
	checkHRSWait(t, cons, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), VAL1, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL3, false)
	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommitWait)
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
		assert.Error(t, cons.AddVote(v))
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

	addr := pvals[VAL1].Address()
	b, _ := block.GenerateTestBlock(&addr)
	p := vote.NewProposal(1, 0, b)

	cons.SetProposal(p)
	assert.Nil(t, cons.LastProposal())

	pvals[VAL2].SignMsg(p)
	cons.SetProposal(p)
	assert.Nil(t, cons.LastProposal())

}
