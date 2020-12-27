package state

import (
	"testing"

	"github.com/zarbchain/zarb-go/block"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
)

func TestSaveLoadLastInfo(t *testing.T) {
	st := setupStatewithOneValidator(t)
	b, _ := block.GenerateTestBlock(nil, nil)

	st.saveLastInfo(125, *b.LastCommit(), b.Header().LastReceiptsHash())
	li, err := st.loadLastInfo()
	assert.NoError(t, err)
	assert.Equal(t, li.LastHeight, 125)
	assert.Equal(t, li.LastCommit.Hash(), b.LastCommit().Hash())
	assert.Equal(t, li.LastReceiptHash, b.Header().LastReceiptsHash())
}

func TestLoadState(t *testing.T) {
	st1 := setupStatewithFourValidators(t, tValSigner1)

	i := 0
	for ; i < 10; i++ {
		b := st1.ProposeBlock()
		c := makeCommitAndSign(t, b.Hash(), 2, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

		require.NoError(t, st1.ApplyBlock(i+1, b, c))
	}

	newBlock := st1.ProposeBlock()
	newCommit := makeCommitAndSign(t, newBlock.Hash(), 1, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.NoError(t, st1.Close())

	// Load last state info
	st2, err := LoadOrNewState(st1.config, st1.genDoc, tValSigner1, txpool.NewMockTxPool())
	require.NoError(t, err)

	assert.Equal(t, st1.store.TotalAccounts(), st2.(*state).store.TotalAccounts())
	assert.Equal(t, st1.store.TotalValidators(), st2.(*state).store.TotalValidators())
	assert.Equal(t, st1.sortition.TotalStake(), st2.(*state).sortition.TotalStake())
	assert.Equal(t, st1.executionSandbox.LastBlockHeight(), st2.(*state).executionSandbox.LastBlockHeight())
	assert.Equal(t, st1.executionSandbox.LastBlockHash(), st2.(*state).executionSandbox.LastBlockHash())

	assert.Equal(t, newBlock.Hash(), st2.ProposeBlock().Hash())
	require.NoError(t, st2.ApplyBlock(i+1, newBlock, newCommit))
}

func TestLoadStateAfterChangingGenesis(t *testing.T) {
	st1 := setupStatewithOneValidator(t)

	// Let's commit some blocks
	i := 0
	for ; i < 10; i++ {
		b1, c1 := proposeAndSignBlock(t, st1)
		assert.NoError(t, st1.ApplyBlock(i+1, b1, c1))
	}

	assert.NoError(t, st1.Close())

	_, err := LoadOrNewState(st1.config, st1.genDoc, tValSigner1, txpool.NewMockTxPool())
	require.NoError(t, err)

	// Load last state info after modifying genesis
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21*1e14 + 1) // manipulating genesis
	val := validator.NewValidator(tValSigner1.PublicKey(), 0, 0)
	genDoc := genesis.MakeGenesis("test", tGenTime, []*account.Account{acc}, []*validator.Validator{val}, 1)

	_, err = LoadOrNewState(st1.config, genDoc, tValSigner1, txpool.NewMockTxPool())
	require.Error(t, err)
}
