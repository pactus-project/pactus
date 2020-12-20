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
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var tValSigner1 crypto.Signer
var tValSigner2 crypto.Signer
var tValSigner3 crypto.Signer
var tValSigner4 crypto.Signer
var tGenTime time.Time
var tCommonTxPool *txpool.MockTxPool

func init() {
	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()

	tValSigner1 = crypto.NewSigner(priv1)
	tValSigner2 = crypto.NewSigner(priv2)
	tValSigner3 = crypto.NewSigner(priv3)
	tValSigner4 = crypto.NewSigner(priv4)

	tGenTime = util.Now()

	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)

	tCommonTxPool = txpool.NewMockTxPool()
}

func setupStatewithFourValidators(t *testing.T, signer crypto.Signer) *state {

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	// 2,100,000,000,000,000
	acc.AddToBalance(21 * 1e14)
	val1 := validator.NewValidator(tValSigner1.PublicKey(), 0, 0)
	val2 := validator.NewValidator(tValSigner2.PublicKey(), 1, 0)
	val3 := validator.NewValidator(tValSigner3.PublicKey(), 2, 0)
	val4 := validator.NewValidator(tValSigner4.PublicKey(), 3, 0)
	gnDoc := genesis.MakeGenesis("test", tGenTime, []*account.Account{acc}, []*validator.Validator{val1, val2, val3, val4}, 1)

	st, err := LoadOrNewState(TestConfig(), gnDoc, signer, tCommonTxPool)
	require.NoError(t, err)
	s, _ := st.(*state)

	return s
}

func setupStatewithOneValidator(t *testing.T) *state {
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(tValSigner1.PublicKey(), 0, 0)
	genDoc := genesis.MakeGenesis("test", tGenTime, []*account.Account{acc}, []*validator.Validator{val}, 1)

	st, err := LoadOrNewState(TestConfig(), genDoc, tValSigner1, txpool.NewMockTxPool())
	require.NoError(t, err)
	s, _ := st.(*state)

	return s
}

func TestProposeBlockAndValidation(t *testing.T) {
	st := setupStatewithOneValidator(t)

	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}

func proposeAndSignBlock(t *testing.T, st *state) (block.Block, block.Commit) {
	addr := tValSigner1.Address()
	b := st.ProposeBlock()
	v := vote.NewPrecommit(1, 0, b.Hash(), addr)
	sig := tValSigner1.Sign(v.SignBytes())
	c := block.NewCommit(0, []block.Committer{{Status: 1, Address: addr}}, *sig)

	return b, *c
}

func TestLoadState(t *testing.T) {
	st1 := setupStatewithOneValidator(t)
	st2 := setupStatewithOneValidator(t)

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
	for ; i < 10; i++ {
		b1, c1 := proposeAndSignBlock(t, st1)
		b2, c2 := proposeAndSignBlock(t, st2)

		assert.NoError(t, st1.ApplyBlock(i+1, b1, c1))
		assert.NoError(t, st2.ApplyBlock(i+1, b2, c2))
	}

	assert.Equal(t, st1.stateHash(), st2.stateHash())

	assert.NoError(t, st2.Close())

	// Load last state info
	st3, err := LoadOrNewState(st2.config, st2.genDoc, tValSigner1, txpool.NewMockTxPool())
	require.NoError(t, err)

	b, c := proposeAndSignBlock(t, st1)
	assert.Equal(t, b.Hash(), st3.ProposeBlock().Hash())
	require.NoError(t, st1.ApplyBlock(i+1, b, c))
	require.NoError(t, st3.ApplyBlock(i+1, b, c))
	assert.Equal(t, st1.store.TotalAccounts(), st3.(*state).store.TotalAccounts())
	assert.Equal(t, st1.store.TotalValidators(), st3.(*state).store.TotalValidators())
	assert.Equal(t, st1.sortition.TotalStake(), st3.(*state).sortition.TotalStake())
	assert.Equal(t, st1.executionSandbox.LastBlockHeight(), st3.(*state).executionSandbox.LastBlockHeight())
	assert.Equal(t, st1.executionSandbox.LastBlockHash(), st3.(*state).executionSandbox.LastBlockHash())
}

func TestBlockSubsidy(t *testing.T) {
	interval := 210000
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy(1, 210000))
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy((1*interval)-1, 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((1*interval), 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((2*interval)-1, 210000))
	assert.Equal(t, int64(1.25*1e8), calcBlockSubsidy((2*interval), 210000))
}

func TestBlockSubsidyTx(t *testing.T) {
	st := setupStatewithOneValidator(t)

	trx := st.createSubsidyTx(7)
	assert.True(t, trx.IsSubsidyTx())
	assert.Equal(t, trx.Payload().Value(), calcBlockSubsidy(1, st.params.SubsidyReductionInterval)+7)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, tValSigner1.Address())

	addr, _, _ := crypto.GenerateTestKeyPair()
	st.config.MintbaseAddress = &addr
	trx = st.createSubsidyTx(0)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, addr)
}

func TestApplyBlocks(t *testing.T) {
	st := setupStatewithOneValidator(t)

	b1, c1 := proposeAndSignBlock(t, st)
	invBlock, _ := block.GenerateTestBlock(nil)
	assert.Error(t, st.ApplyBlock(1, *invBlock, c1))
	assert.Error(t, st.ApplyBlock(2, b1, c1))
}

func TestProposeBlock(t *testing.T) {
	st1 := setupStatewithOneValidator(t)

	b1, c1 := proposeAndSignBlock(t, st1)
	assert.NoError(t, st1.ApplyBlock(1, b1, c1))

	subsidy := calcBlockSubsidy(st1.LastBlockHeight(), st1.params.SubsidyReductionInterval)
	invSubsidyTx := tx.NewSubsidyTx(st1.LastBlockHash(), 1, tValSigner2.Address(), subsidy, "")
	invSendTx, _ := tx.GenerateTestSendTx()
	invBondTx, _ := tx.GenerateTestBondTx()
	invSortitionTx, _ := tx.GenerateTestSortitionTx()

	pub := tValSigner1.PublicKey()
	trx1 := tx.NewSendTx(b1.Hash(), 1, tValSigner1.Address(), tValSigner1.Address(), 1, 1000, "", &pub, nil)
	tValSigner1.SignMsg(trx1)

	trx2 := tx.NewBondTx(b1.Hash(), 2, tValSigner1.Address(), pub, 1, "", &pub, nil)
	tValSigner1.SignMsg(trx2)

	assert.NoError(t, st1.txPool.AppendTx(invSendTx))
	assert.NoError(t, st1.txPool.AppendTx(invBondTx))
	assert.NoError(t, st1.txPool.AppendTx(invSortitionTx))
	assert.NoError(t, st1.txPool.AppendTx(invSubsidyTx))
	assert.NoError(t, st1.txPool.AppendTx(trx1))
	assert.NoError(t, st1.txPool.AppendTx(trx2))

	b2 := st1.ProposeBlock()
	assert.Equal(t, b2.Header().LastBlockHash(), b1.Hash())
	assert.Equal(t, b2.TxIDs().IDs()[1:], []crypto.Hash{trx1.ID(), trx2.ID()})
}
