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
	cons  *Consensus
	st    state.State
	pvals []*validator.PrivValidator
)

const (
	VAL_1 = 0
	VAL_2 = 1
	VAL_3 = 2
	VAL_4 = 3
)

func newTestConsensus(t *testing.T, val_id int) (*Consensus, []*validator.PrivValidator) {
	_, keys := validator.GenerateTestValidatorSet()
	consConf := TestConfig()
	stateConf := state.TestConfig()
	txPoolConf := txpool.TestConfig()
	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)

	vals := make([]*validator.Validator, 4)
	pvals := make([]*validator.PrivValidator, 4)
	for i, k := range keys {
		val := validator.NewValidator(k.PublicKey(), 0)
		val.AddToStake(100)
		vals[i] = val

		pval := validator.NewPrivValidator(k)
		pvals[i] = pval
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
	st, _ = state.LoadOrNewState(stateConf, genDoc, pvals[val_id].Address(), txpool)

	cons, _ := NewConsensus(consConf, st, pvals[val_id], ch)
	assert.Equal(t, cons.votes.height, 0)
	assert.Equal(t, hrs.NewHRS(0, 0, hrs.StepTypeNewHeight), cons.hrs)
	cons.ScheduleNewHeight()
	assert.Equal(t, hrs.NewHRS(0, 0, hrs.StepTypeNewHeight), cons.hrs)

	// Calling ScheduleNewHeight for the second time
	cons.ScheduleNewHeight()
	assert.Equal(t, hrs.NewHRS(0, 0, hrs.StepTypeNewHeight), cons.hrs)

	return cons, pvals
}

func checkHRS(t *testing.T, height, round int, step hrs.StepType) {
	assert.Equal(t, hrs.NewHRS(height, round, step), cons.hrs)
}

func testAddVote(t *testing.T,
	voteType vote.VoteType,
	height int,
	round int,
	blockHash crypto.Hash,
	pval_id int,
	expectError bool) {

	v := vote.NewVote(voteType, height, round, blockHash, pvals[pval_id].Address())
	pvals[pval_id].SignMsg(v)

	if expectError {
		assert.Error(t, cons.AddVote(v))
	} else {
		assert.NoError(t, cons.AddVote(v))
	}
}

func TestConsensusAddVotesNormal(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_1)

	cons.enterNewHeight(1)

	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL_2, false)
	checkHRS(t, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 0, hrs.StepTypePrecommit)

	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL_2, false)
	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
	assert.Equal(t, cons.votes.Precommits(0).Len(), 3) // Votes from validator 1,2,3
}
func TestConsensusUpdateVote(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_1)

	cons.enterNewHeight(1)

	p := cons.LastProposal()
	assert.NotNil(t, p)

	testAddVote(t, vote.VoteTypePrevote, 2, 0, p.Block().Hash(), VAL_2, true)
	// Validator_2 doesn't have proposal now vote for nil
	testAddVote(t, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL_2, false)
	checkHRS(t, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 0, hrs.StepTypePrevoteWait)

	// Validator_2 have proposal now and vote for that
	testAddVote(t, vote.VoteTypePrevote, 1, 0, p.Block().Hash(), VAL_2, false)
	checkHRS(t, 1, 0, hrs.StepTypePrecommit)

	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL_2, false)
	assert.Equal(t, cons.isCommitted, false)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
	assert.Equal(t, cons.votes.Precommits(0).Len(), 3)
}

func TestConsensusNoPrevotes(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_1)

	cons.enterNewHeight(1)
	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL_2, false)
	checkHRS(t, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 0, hrs.StepTypePrevote)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, p.Block().Hash(), VAL_4, false)
	checkHRS(t, 1, 0, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
	assert.Equal(t, cons.votes.Precommits(0).Len(), 3)
}

func TestConsensusGotoNextRound(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_2)

	cons.enterNewHeight(1)

	// Validator_1 is offline
	testAddVote(t, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL_2, false)
	testAddVote(t, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL_3, false)
	testAddVote(t, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL_4, false)
	checkHRS(t, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL_2, false)
	testAddVote(t, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL_3, false)
	testAddVote(t, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL_4, false)
	checkHRS(t, 1, 1, hrs.StepTypePrevote)

	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL_1, false)
	checkHRS(t, 1, 1, hrs.StepTypePrevote)

	testAddVote(t, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL_1, false)
	checkHRS(t, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 1, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
}

func TestConsensusGotoNextRound2(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_2)

	cons.enterNewHeight(1)

	// Validator_1 is online, but the proposal is not accepted by other nodes
	// Validator_4 is offline
	testAddVote(t, vote.VoteTypePrevote, 1, 0, crypto.GenerateTestHash(), VAL_1, false)
	testAddVote(t, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL_2, false)
	testAddVote(t, vote.VoteTypePrevote, 1, 0, crypto.UndefHash, VAL_3, false)
	checkHRS(t, 1, 0, hrs.StepTypePrevoteWait)
	time.Sleep(200 * time.Millisecond)
	checkHRS(t, 1, 0, hrs.StepTypePrecommit)

	testAddVote(t, vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), VAL_1, false)
	testAddVote(t, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL_2, false)
	testAddVote(t, vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, VAL_3, false)
	checkHRS(t, 1, 0, hrs.StepTypePrecommitWait)
	time.Sleep(200 * time.Millisecond)
	checkHRS(t, 1, 1, hrs.StepTypePrevote)

	p := cons.LastProposal()
	require.NotNil(t, p)

	testAddVote(t, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL_1, false)
	checkHRS(t, 1, 1, hrs.StepTypePrevote)

	testAddVote(t, vote.VoteTypePrevote, 1, 1, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL_1, false)
	checkHRS(t, 1, 1, hrs.StepTypePrecommit)

	testAddVote(t, vote.VoteTypePrecommit, 1, 1, p.Block().Hash(), VAL_3, false)
	checkHRS(t, 1, 1, hrs.StepTypeCommit)
	assert.Equal(t, cons.isCommitted, true)
}
func TestConsensusSpamming(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_1)

	cons.enterNewHeight(1)

	for i := 0; i < 100; i++ {
		v, _ := vote.GenerateTestPrecommitVote(1, 0)
		assert.Error(t, cons.AddVote(v))
	}
}

func TestConsensusSpammingProposal(t *testing.T) {
	cons, pvals = newTestConsensus(t, VAL_2)

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
	cons, pvals = newTestConsensus(t, VAL_2)

	cons.enterNewHeight(1)
	assert.Nil(t, cons.LastProposal())

	addr := pvals[VAL_1].Address()
	b, _ := block.GenerateTestBlock(&addr)
	p := vote.NewProposal(1, 0, b)

	cons.SetProposal(p)
	assert.Nil(t, cons.LastProposal())

	pvals[VAL_2].SignMsg(p)
	cons.SetProposal(p)
	assert.Nil(t, cons.LastProposal())

}
