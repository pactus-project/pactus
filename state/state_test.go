package state

import (
	"testing"
	"time"

	"github.com/zarbchain/zarb-go/tx"

	"github.com/zarbchain/zarb-go/tx/payload"

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

var tTxPool *txpool.MockTxPool
var tGenDoc *genesis.Genesis
var tValSigner crypto.Signer

func setup(t *testing.T) {
	tTxPool = txpool.NewMockTxPool()

	_, pb, priv := crypto.GenerateTestKeyPair()
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21000000000000)
	val := validator.NewValidator(pb, 0, 0)
	val.AddToStake(util.RandInt64(10000000))
	tGenDoc = genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val}, 1)
	tValSigner = crypto.NewSigner(priv)

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
	st, err := LoadOrNewState(stateConfig, tGenDoc, *signer, tTxPool)
	require.NoError(t, err)
	s, _ := st.(*state)

	return s, signer.Address()
}

func TestProposeBlockAndValidation(t *testing.T) {
	setup(t)

	st, _ := mockState(t, &tValSigner)
	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}

func proposeAndSignBlock(t *testing.T, st *state) (block.Block, block.Commit) {
	addr := tValSigner.Address()
	b := st.ProposeBlock()
	v := vote.NewPrecommit(1, 0, b.Hash(), addr)
	sig := tValSigner.Sign(v.SignBytes())
	c := block.NewCommit(0, []block.Committer{{Status: 1, Address: addr}}, *sig)

	return b, *c
}

func TestLoadState(t *testing.T) {
	setup(t)

	st1, _ := mockState(t, &tValSigner)
	st2, _ := mockState(t, &tValSigner)

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
		b, c := proposeAndSignBlock(t, st1)

		assert.NoError(t, st1.ApplyBlock(i+1, b, c))
		assert.NoError(t, st2.ApplyBlock(i+1, b, c))
	}

	assert.NoError(t, st2.Close())

	// Load last state info
	st3, err := LoadOrNewState(st2.config, tGenDoc, tValSigner, tTxPool)
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
	setup(t)

	interval := 210000
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy(1, 210000))
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy((1*interval)-1, 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((1*interval), 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((2*interval)-1, 210000))
	assert.Equal(t, int64(1.25*1e8), calcBlockSubsidy((2*interval), 210000))
}

func TestBlockSubsidyTx(t *testing.T) {
	setup(t)
	st, _ := mockState(t, &tValSigner)

	trx := st.createSubsidyTx(7)
	assert.True(t, trx.IsSubsidyTx())
	assert.Equal(t, trx.Payload().Value(), calcBlockSubsidy(1, st.params.SubsidyReductionInterval)+7)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, tValSigner.Address())

	addr, _, _ := crypto.GenerateTestKeyPair()
	st.config.MintbaseAddress = &addr
	trx = st.createSubsidyTx(0)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, addr)
}

func TestApplyBlocks(t *testing.T) {
	setup(t)

	st, _ := mockState(t, &tValSigner)
	b1, c1 := proposeAndSignBlock(t, st)
	invBlock, _ := block.GenerateTestBlock(nil)
	assert.Error(t, st.ApplyBlock(1, invBlock, c1))
	assert.Error(t, st.ApplyBlock(2, b1, c1))
}

func TestProposeBlock(t *testing.T) {
	setup(t)

	st1, _ := mockState(t, &tValSigner)
	st2, _ := mockState(t, nil)

	b1, c1 := proposeAndSignBlock(t, st1)
	assert.NoError(t, st1.ApplyBlock(1, b1, c1))

	invSubsidyTx := st2.createSubsidyTx(100)
	invSendTx, _ := tx.GenerateTestSendTx()
	invBondTx, _ := tx.GenerateTestBondTx()
	invSortitionTx, _ := tx.GenerateTestSortitionTx()

	pub := tValSigner.PublicKey()
	trx1 := tx.NewSendTx(b1.Hash(), 1, tValSigner.Address(), tValSigner.Address(), 1, 1000, "", &pub, nil)
	tValSigner.SignMsg(trx1)

	trx2 := tx.NewBondTx(b1.Hash(), 2, tValSigner.Address(), pub, 1, "", &pub, nil)
	tValSigner.SignMsg(trx2)

	tTxPool.AppendTx(invSendTx)
	tTxPool.AppendTx(invBondTx)
	tTxPool.AppendTx(invSortitionTx)
	tTxPool.AppendTx(invSubsidyTx)
	tTxPool.AppendTx(trx1)
	tTxPool.AppendTx(trx2)

	b2 := st1.ProposeBlock()
	assert.Equal(t, b2.Header().LastBlockHash(), b1.Hash())
	assert.Equal(t, b2.TxIDs().IDs()[1:], []crypto.Hash{trx1.ID(), trx2.ID()})
}
