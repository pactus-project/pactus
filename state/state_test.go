package state

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var mockTxPool *txpool.MockTxPool
var gen *genesis.Genesis
var valSigner crypto.Signer

func setup(t *testing.T) {
	mockTxPool = txpool.NewMockTxPool()

	_, pb, priv := crypto.GenerateTestKeyPair()
	acc := account.NewAccount(crypto.MintbaseAddress, 0)
	acc.AddToBalance(21000000000000)
	val := validator.NewValidator(pb, 0, 0)
	val.AddToStake(util.RandInt64(10000000))
	gen = genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val}, 1)
	valSigner = crypto.NewSigner(priv)

	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)
}

func mockState(t *testing.T, signer *crypto.Signer) (*state, crypto.Address) {
	if signer == nil {
		_, _, priv := crypto.GenerateTestKeyPair()
		s := crypto.NewSigner(priv)
		signer = &s
	}

	stateConfig := TestConfig()
	st, err := LoadOrNewState(stateConfig, gen, *signer, mockTxPool)
	require.NoError(t, err)
	s, _ := st.(*state)

	return s, signer.Address()
}

func TestProposeBlockValidation(t *testing.T) {
	setup(t)

	st, _ := mockState(t, &valSigner)
	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}

func proposeAndSignBlock(t *testing.T, st *state) (block.Block, block.Commit) {
	addr := valSigner.Address()
	b := st.ProposeBlock()
	v := vote.NewPrecommit(1, 0, b.Hash(), addr)
	sig := valSigner.Sign(v.SignBytes())
	c := block.NewCommit(0, []block.Committer{{Status: 1, Address: addr}}, *sig)

	return b, *c
}

func TestLoadState(t *testing.T) {
	setup(t)

	st1, _ := mockState(t, &valSigner)
	st2, _ := mockState(t, &valSigner)

	// Add this dummy acc and val for testing purpose
	dummyAcc, _ := account.GenerateTestAccount(1)
	st1.store.UpdateAccount(dummyAcc)
	st2.store.UpdateAccount(dummyAcc)
	dummyVal, _ := validator.GenerateTestValidator(1)
	st1.store.UpdateValidator(dummyVal)
	st2.store.UpdateValidator(dummyVal)
	st1.sortition.AddToTotalStake(dummyVal.Stake())
	st2.sortition.AddToTotalStake(dummyVal.Stake())

	i := 0
	for ; i < st1.params.TransactionToLiveInterval+10; i++ {
		b, c := proposeAndSignBlock(t, st1)

		assert.NoError(t, st1.ApplyBlock(i+1, b, c))
		assert.NoError(t, st2.ApplyBlock(i+1, b, c))
	}

	assert.NoError(t, st2.Close())

	// Load last state info
	st3, err := LoadOrNewState(st2.config, gen, valSigner, mockTxPool)
	require.NoError(t, err)

	b, c := proposeAndSignBlock(t, st1)
	assert.Equal(t, b.Hash(), st3.ProposeBlock().Hash())
	require.NoError(t, st1.ApplyBlock(i+1, b, c))
	require.NoError(t, st3.ApplyBlock(i+1, b, c))

	t.Run("Check sandbox after loading blockchain", func(t *testing.T) {
		assert.Equal(t, st1.executionSandbox.recentBlocks, st3.(*state).executionSandbox.recentBlocks)
		assert.Equal(t, st1.executionSandbox.totalAccounts, st3.(*state).executionSandbox.totalAccounts)
		assert.Equal(t, st1.executionSandbox.totalValidators, st3.(*state).executionSandbox.totalValidators)
		assert.Equal(t, st1.sortition.TotalStake(), st3.(*state).sortition.TotalStake())
	})
}

func TestBlockSubsidy(t *testing.T) {
	setup(t)

	interval := 210000
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy(1, 210000))
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy((1*interval)-1, 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((1*interval), 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((2*interval)-1, 210000))
	assert.Equal(t, int64(1.25*1e8), calcBlockSubsidy((2*interval), 210000))
}

func TestApplyBlocks(t *testing.T) {
	setup(t)

	st, _ := mockState(t, &valSigner)
	b1, c1 := proposeAndSignBlock(t, st)
	invBlock, _ := block.GenerateTestBlock(nil)
	assert.Error(t, st.ApplyBlock(1, invBlock, c1))
	assert.Error(t, st.ApplyBlock(2, b1, c1))

}
