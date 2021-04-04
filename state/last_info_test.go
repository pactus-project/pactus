package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
)

func TestSaveLoadLastInfo(t *testing.T) {
	setup(t)
	for i := 0; i < 4; i++ {
		moveToNextHeightForAllStates(t)
	}

	tState1.saveLastInfo(
		tState1.lastInfo.BlockHeight(),
		*tState1.lastInfo.Certificate(),
		tState1.lastInfo.ReceiptsHash(),
		tState1.committee.Committers(),
		tState1.committee.Proposer(0).Address())

	li, err := tState1.loadLastInfo()
	assert.NoError(t, err)
	assert.Equal(t, li.LastHeight, tState1.lastInfo.BlockHeight())
	assert.Equal(t, li.LastCertificate.Hash(), tState1.lastInfo.Certificate().Hash())
	assert.Equal(t, li.LastReceiptHash, tState1.lastInfo.ReceiptsHash())
	assert.Equal(t, li.Committee, tState1.committee.Committers())
	assert.Equal(t, li.NextProposer, tState1.committee.Proposer(0).Address())
}

func TestLoadState(t *testing.T) {
	setup(t)

	// Add a bond transactions to change total stake
	_, pub, _ := crypto.GenerateTestKeyPair()
	tx2 := tx.NewBondTx(crypto.UndefHash, 1, tValSigner1.Address(), pub, 8888000, 8888, "")
	tValSigner1.SignMsg((tx2))

	assert.NoError(t, tCommonTxPool.AppendTx(tx2))

	for i := 0; i < 4; i++ {
		moveToNextHeightForAllStates(t)
	}
	b5, c5 := makeBlockAndCertificate(t, 1, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockForAllStates(t, b5, c5)

	b6, c6 := makeBlockAndCertificate(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.NoError(t, tState1.Close())

	// Load last state info
	st2, err := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, tCommonTxPool)
	require.NoError(t, err)

	assert.Equal(t, tState1.store.TotalAccounts(), st2.(*state).store.TotalAccounts())
	assert.Equal(t, tState1.store.TotalValidators(), st2.(*state).store.TotalValidators())
	assert.Equal(t, tState1.sortition.TotalStake(), st2.(*state).sortition.TotalStake()-4) //4 validators are yet in genesis state
	assert.Equal(t, tState1.store.TotalAccounts(), 5)
	assert.Equal(t, tState1.sortition.TotalStake(), int64(8888000))

	require.NoError(t, st2.CommitBlock(6, b6, c6))
	require.NoError(t, tState2.CommitBlock(6, b6, c6))
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
	genDoc := genesis.MakeGenesis(tGenTime, []*account.Account{acc}, []*validator.Validator{val}, param.DefaultParams())

	_, err = LoadOrNewState(tState1.config, genDoc, tValSigner1, txpool.MockingTxPool())
	require.Error(t, err)
}
