package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
)

func TestSaveLoadLastInfo(t *testing.T) {
	setup(t)
	b, _ := block.GenerateTestBlock(nil, nil)

	tState1.saveLastInfo(125, *b.LastCommit(), b.Header().LastReceiptsHash())
	li, err := tState1.loadLastInfo()
	assert.NoError(t, err)
	assert.Equal(t, li.LastHeight, 125)
	assert.Equal(t, li.LastCommit.Hash(), b.LastCommit().Hash())
	assert.Equal(t, li.LastReceiptHash, b.Header().LastReceiptsHash())
}

func TestLoadState(t *testing.T) {
	setup(t)

	addr, pub, _ := crypto.GenerateTestKeyPair()

	tx1 := tx.NewSendTx(crypto.UndefHash, 1, tValSigner1.Address(), addr, 8888000, 8888, "")
	tValSigner1.SignMsg((tx1))

	tx2 := tx.NewBondTx(crypto.UndefHash, 2, tValSigner1.Address(), pub, 8888000, 8888, "")
	tValSigner1.SignMsg((tx2))

	assert.NoError(t, tCommonTxPool.AppendTx(tx1))
	assert.NoError(t, tCommonTxPool.AppendTx(tx2))

	i := 0
	for ; i < 8; i++ {
		moveToNextHeightForAllStates(t)
	}

	newBlock, newCommit := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.NoError(t, tState1.Close())

	// Load last state info
	st2, err := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, txpool.MockingTxPool())
	require.NoError(t, err)

	assert.Equal(t, tState1.store.TotalAccounts(), st2.(*state).store.TotalAccounts())
	assert.Equal(t, tState1.store.TotalValidators(), st2.(*state).store.TotalValidators())
	assert.Equal(t, tState1.sortition.TotalStake(), st2.(*state).sortition.TotalStake())
	assert.Equal(t, tState1.executionSandbox.LastBlockHeight(), st2.(*state).executionSandbox.LastBlockHeight())
	assert.Equal(t, tState1.executionSandbox.LastBlockHash(), st2.(*state).executionSandbox.LastBlockHash())
	assert.Equal(t, tState1.store.TotalAccounts(), 6)
	assert.Equal(t, tState1.sortition.TotalStake(), int64(8888000))

	proposed, err := st2.ProposeBlock(0)
	assert.NoError(t, err)
	assert.Equal(t, newBlock.Hash(), proposed.Hash())
	require.NoError(t, st2.CommitBlock(i+1, newBlock, newCommit))
}

func TestLoadStateAfterChangingGenesis(t *testing.T) {
	setup(t)

	// Let's commit some blocks
	i := 0
	for ; i < 10; i++ {
		moveToNextHeightForAllStates(t)
	}

	assert.NoError(t, tState1.Close())

	_, err := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, txpool.MockingTxPool())
	require.NoError(t, err)

	// Load last state info after modifying genesis
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21*1e14 + 1) // manipulating genesis
	val := validator.NewValidator(tValSigner1.PublicKey(), 0, 0)
	genDoc := genesis.MakeGenesis("test", tGenTime, []*account.Account{acc}, []*validator.Validator{val}, param.MainnetParams())

	_, err = LoadOrNewState(tState1.config, genDoc, tValSigner1, txpool.MockingTxPool())
	require.Error(t, err)
}
