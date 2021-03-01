package state

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var tState1 *state
var tState2 *state
var tState3 *state
var tState4 *state
var tValSigner1 crypto.Signer
var tValSigner2 crypto.Signer
var tValSigner3 crypto.Signer
var tValSigner4 crypto.Signer
var tGenTime time.Time
var tCommonTxPool *txpool.MockTxPool

func setup(t *testing.T) {
	if tState1 != nil {
		tState1.Close()
		tState2.Close()
		tState3.Close()
		tState4.Close()
	}
	logger.InitLogger(logger.TestConfig())

	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()

	tValSigner1 = crypto.NewSigner(priv1)
	tValSigner2 = crypto.NewSigner(priv2)
	tValSigner3 = crypto.NewSigner(priv3)
	tValSigner4 = crypto.NewSigner(priv4)

	tGenTime = util.RoundNow(10)
	tCommonTxPool = txpool.MockingTxPool()

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14) // 2,100,000,000,000,000
	val1 := validator.NewValidator(tValSigner1.PublicKey(), 0, 0)
	val2 := validator.NewValidator(tValSigner2.PublicKey(), 1, 0)
	val3 := validator.NewValidator(tValSigner3.PublicKey(), 2, 0)
	val4 := validator.NewValidator(tValSigner4.PublicKey(), 3, 0)
	params := param.DefaultParams()
	params.CommitteeSize = 4
	gnDoc := genesis.MakeGenesis(tGenTime, []*account.Account{acc}, []*validator.Validator{val1, val2, val3, val4}, params)

	st1, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner1, tCommonTxPool)
	require.NoError(t, err)
	st2, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner2, tCommonTxPool)
	require.NoError(t, err)
	st3, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner3, tCommonTxPool)
	require.NoError(t, err)
	st4, err := LoadOrNewState(TestConfig(), gnDoc, tValSigner4, tCommonTxPool)
	require.NoError(t, err)

	tState1, _ = st1.(*state)
	tState2, _ = st2.(*state)
	tState3, _ = st3.(*state)
	tState4, _ = st4.(*state)
}

func makeBlockAndCommit(t *testing.T, round int, signers ...crypto.Signer) (block.Block, block.Commit) {
	var st *state
	if tState1.validatorSet.IsProposer(tState1.signer.Address(), round) {
		st = tState1
	} else if tState1.validatorSet.IsProposer(tState2.signer.Address(), round) {
		st = tState2
	} else if tState1.validatorSet.IsProposer(tState3.signer.Address(), round) {
		st = tState3
	} else {
		st = tState4
	}

	b, err := st.ProposeBlock(round)
	require.NoError(t, err)
	c := makeCommitAndSign(t, b.Hash(), round, signers...)

	return *b, c
}

func makeCommitAndSign(t *testing.T, blockHash crypto.Hash, round int, signers ...crypto.Signer) block.Commit {
	sigs := make([]crypto.Signature, len(signers))
	sb := block.CommitSignBytes(blockHash, round)
	committers := make([]block.Committer, 4)
	committers[0] = block.Committer{Status: 0, Number: 0}
	committers[1] = block.Committer{Status: 0, Number: 1}
	committers[2] = block.Committer{Status: 0, Number: 2}
	committers[3] = block.Committer{Status: 0, Number: 3}

	for i, s := range signers {
		if s.Address().EqualsTo(tValSigner1.Address()) {
			committers[0] = block.Committer{Status: 1, Number: 0}
		}

		if s.Address().EqualsTo(tValSigner2.Address()) {
			committers[1] = block.Committer{Status: 1, Number: 1}
		}

		if s.Address().EqualsTo(tValSigner3.Address()) {
			committers[2] = block.Committer{Status: 1, Number: 2}
		}

		if s.Address().EqualsTo(tValSigner4.Address()) {
			committers[3] = block.Committer{Status: 1, Number: 3}
		}
		sigs[i] = s.SignData(sb)
	}
	return *block.NewCommit(blockHash, round, committers, crypto.Aggregate(sigs))
}

func CommitBlockAndCommitForAllStates(t *testing.T, b block.Block, c block.Commit) {
	assert.NoError(t, tState1.CommitBlock(tState1.lastBlockHeight+1, b, c))
	assert.NoError(t, tState2.CommitBlock(tState2.lastBlockHeight+1, b, c))
	assert.NoError(t, tState3.CommitBlock(tState3.lastBlockHeight+1, b, c))
	assert.NoError(t, tState4.CommitBlock(tState4.lastBlockHeight+1, b, c))
}

func moveToNextHeightForAllStates(t *testing.T) {
	b, c := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockAndCommitForAllStates(t, b, c)
}

func TestProposeBlockAndValidation(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	b, err := tState1.ProposeBlock(0)
	assert.Error(t, err)
	assert.Nil(t, b)

	trx := tx.NewSendTx(crypto.UndefHash, 1, tValSigner1.Address(), tValSigner2.Address(), 1000, 1000, "")
	tValSigner1.SignMsg(trx)
	assert.NoError(t, tCommonTxPool.AppendTx(trx))

	b, err = tState2.ProposeBlock(0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.TxIDs().Len(), 2)

	err = tState1.ValidateBlock(*b)
	require.NoError(t, err)

	// Propose and validate again
	b, err = tState2.ProposeBlock(0)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.TxIDs().Len(), 2)

	err = tState1.ValidateBlock(*b)
	require.NoError(t, err)
}

func TestBlockSubsidy(t *testing.T) {
	interval := 2100000
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy(1, interval))
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy((1*interval)-1, interval))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((1*interval), interval))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((2*interval)-1, interval))
	assert.Equal(t, int64(1.25*1e8), calcBlockSubsidy((2*interval), interval))
}

func TestBlockSubsidyTx(t *testing.T) {
	setup(t)

	// Without mintbase address in config
	trx := tState1.createSubsidyTx(7)
	assert.True(t, trx.IsMintbaseTx())
	assert.Equal(t, trx.Payload().Value(), calcBlockSubsidy(1, tState1.params.SubsidyReductionInterval)+7)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, tValSigner1.Address())

	// With ivalid mintbase address in config
	tState1.config.MintbaseAddress = "invalid"
	tState1.Close()
	_, err := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, tCommonTxPool)
	assert.Error(t, err)

	// With mintbase address in config
	addr, _, _ := crypto.GenerateTestKeyPair()
	tState1.config.MintbaseAddress = addr.String()
	tState1.Close()
	st, err := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, tCommonTxPool)
	assert.NoError(t, err)
	trx = st.(*state).createSubsidyTx(0)
	assert.Equal(t, trx.Payload().(*payload.SendPayload).Receiver, addr)
}

func TestCommitBlocks(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCommit(t, 1, tValSigner1, tValSigner2, tValSigner3)
	invBlock, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.CommitBlock(1, *invBlock, c1))
	// No error here but block is ignored, because the height is invalid
	assert.NoError(t, tState1.CommitBlock(2, b1, c1))
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))

	assert.Equal(t, tState1.LastBlockHash(), b1.Hash())
	assert.Equal(t, tState1.LastBlockTime(), b1.Header().Time())
	assert.Equal(t, tState1.LastCommit().Hash(), c1.Hash())
	assert.Equal(t, tState1.LastBlockHeight(), 1)
	assert.Equal(t, tState1.GenesisHash(), tState2.GenesisHash())
}

func TestCommitSandbox(t *testing.T) {

	t.Run("Commit new account", func(t *testing.T) {
		setup(t)

		addr, _, _ := crypto.GenerateTestKeyPair()
		newAcc := tState1.executionSandbox.MakeNewAccount(addr)
		newAcc.AddToBalance(1)
		tState1.commitSandbox(0)

		assert.True(t, tState1.store.HasAccount(addr))
	})

	t.Run("Commit new validator", func(t *testing.T) {
		setup(t)

		addr, pub, _ := crypto.GenerateTestKeyPair()
		newVal := tState1.executionSandbox.MakeNewValidator(pub)
		newVal.AddToStake(1)
		tState1.executionSandbox.UpdateValidator(newVal)
		tState1.commitSandbox(0)

		assert.True(t, tState1.store.HasValidator(addr))
		assert.Equal(t, tState1.executionSandbox.TotalStakeChange(), int64(1))
		assert.Equal(t, tState1.sortition.TotalStake(), int64(1))
	})

	t.Run("Modify account", func(t *testing.T) {
		setup(t)

		acc := tState1.executionSandbox.Account(crypto.TreasuryAddress)
		acc.SubtractFromBalance(1)
		tState1.executionSandbox.UpdateAccount(acc)
		tState1.commitSandbox(0)

		acc1, _ := tState1.store.Account(crypto.TreasuryAddress)
		assert.Equal(t, acc1.Balance(), acc.Balance())
	})

	t.Run("Modify validator", func(t *testing.T) {
		setup(t)

		val := tState1.executionSandbox.Validator(tValSigner2.Address())
		val.AddToStake(2)
		tState1.executionSandbox.UpdateValidator(val)
		tState1.commitSandbox(0)

		val1, _ := tState1.store.Validator(tValSigner2.Address())
		assert.Equal(t, val1.Stake(), val.Stake())
		assert.Equal(t, tState1.executionSandbox.TotalStakeChange(), int64(2))
	})

	t.Run("Move valset", func(t *testing.T) {
		setup(t)

		nextProposer := tState1.validatorSet.Proposer(1)

		tState1.commitSandbox(0)

		assert.Equal(t, tState1.validatorSet.Proposer(0).Address(), nextProposer.Address())
	})

	t.Run("Move valset next round", func(t *testing.T) {
		setup(t)

		nextNextProposer := tState1.validatorSet.Proposer(2)

		tState1.commitSandbox(1)

		assert.Equal(t, tState1.validatorSet.Proposer(0).Address(), nextNextProposer.Address())
	})
}

func TestUpdateLastCommit(t *testing.T) {
	setup(t)
	b1, c1 := makeBlockAndCommit(t, 0, tValSigner1, tValSigner3, tValSigner4)
	b11, c11 := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	assert.Equal(t, b1.Hash(), b11.Hash())

	CommitBlockAndCommitForAllStates(t, b1, c1)

	b2, c2 := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	assert.NotEqual(t, b1.Hash(), b2.Hash())

	assert.Equal(t, tState1.lastCommit.Hash(), c1.Hash())
	assert.Error(t, tState1.UpdateLastCommit(&c2))
	assert.NoError(t, tState1.UpdateLastCommit(&c1))
	assert.Equal(t, tState1.lastCommit.Hash(), c1.Hash())
	assert.NoError(t, tState1.UpdateLastCommit(&c11))
	assert.NoError(t, tState1.UpdateLastCommit(&c1))
	assert.Equal(t, tState1.lastCommit.Hash(), c11.Hash())
}

func TestInvalidProposerProposeBlock(t *testing.T) {
	setup(t)

	_, err := tState2.ProposeBlock(0)
	assert.Error(t, err)
	_, err = tState2.ProposeBlock(1)
	assert.NoError(t, err)
}

func TestBlockProposal(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	t.Run("validity of proposed block", func(t *testing.T) {
		b, err := tState2.ProposeBlock(0)
		assert.NoError(t, err)
		assert.NoError(t, tState1.ValidateBlock(*b)) // State1 check state2's proposed block
	})

	t.Run("Tx pool has two subsidy transactions", func(t *testing.T) {
		trx := tState2.createSubsidyTx(0)
		assert.NoError(t, tState2.txPool.AppendTx(trx))

		// Moving to the next round
		b, err := tState3.ProposeBlock(1)
		assert.NoError(t, err)
		assert.NoError(t, tState1.ValidateBlock(*b))
	})
}

func TestInvalidBlock(t *testing.T) {
	setup(t)

	b, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.ValidateBlock(*b))
}

func TestForkDetection(t *testing.T) {
	setup(t)

	b1, c1 := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3)
	b2, c2 := makeBlockAndCommit(t, 1, tValSigner1, tValSigner2, tValSigner3)
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))
	assert.NoError(t, tState1.CommitBlock(1, b1, c1))
	assert.Error(t, tState1.CommitBlock(1, b2, c2))
}

func TestNodeShutdown(t *testing.T) {
	setup(t)
	b1, c1 := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3)

	// Should not panic or crash
	tState1.Close()
	assert.Error(t, tState1.CommitBlock(1, b1, c1))
	b, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, tState1.ValidateBlock(*b))
	_, err := tState1.ProposeBlock(0)
	assert.Error(t, err)
}

func TestSortition(t *testing.T) {
	setup(t)

	addr, pub, priv := crypto.GenerateTestKeyPair()
	signer := crypto.NewSigner(priv)

	st, err := LoadOrNewState(TestConfig(), tState1.genDoc, signer, tCommonTxPool)
	assert.NoError(t, err)
	st1 := st.(*state)

	assert.False(t, st1.evaluateSortition()) //  not a validator

	// Commit 10 blocks, bonding tx is in block 2
	for i := 0; i < 10; i++ {
		if i == 2 {
			trx := tx.NewBondTx(crypto.UndefHash, 1, tValSigner1.Address(), pub, 1000, 1000, "")
			tValSigner1.SignMsg(trx)
			assert.NoError(t, tCommonTxPool.AppendTx(trx))
		}

		b, c := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
		CommitBlockAndCommitForAllStates(t, b, c)
		require.NoError(t, st1.CommitBlock(i+1, b, c))
	}

	assert.False(t, st1.evaluateSortition()) //  bonding period

	// Commit another block
	b, c := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockAndCommitForAllStates(t, b, c)
	require.NoError(t, st1.CommitBlock(11, b, c))

	assert.True(t, st1.evaluateSortition())                //  ok
	assert.False(t, tState1.ValidatorSet().Contains(addr)) // still not in the set

	// ---------------------------------------------
	// Commit another block, new validator should be in the set now
	b, c = makeBlockAndCommit(t, 1, tValSigner1, tValSigner2, tValSigner3, tValSigner4)
	CommitBlockAndCommitForAllStates(t, b, c)
	require.NoError(t, st1.CommitBlock(12, b, c))

	assert.False(t, st1.evaluateSortition()) // already in the set
	assert.False(t, tState1.ValidatorSet().Contains(tValSigner1.Address()))
	assert.True(t, tState1.ValidatorSet().Contains(addr))

	// ---------------------------------------------
	// Let's save and load tState1
	committeeHash := tState1.validatorSet.CommitteeHash()
	tState1.Close()
	state1, _ := LoadOrNewState(tState1.config, tState1.genDoc, tValSigner1, tCommonTxPool)

	assert.Equal(t, state1.(*state).validatorSet.CommitteeHash(), committeeHash)

	// ---------------------------------------------
	// Let's commit another block with new Validator set
	b1, err := st1.ProposeBlock(3)
	require.NoError(t, err)
	require.NotNil(t, b1)

	sigs := make([]crypto.Signature, 4)
	sb := block.CommitSignBytes(b1.Hash(), 3)
	committers := make([]block.Committer, 4)
	committers[0] = block.Committer{Status: 1, Number: 1}
	committers[1] = block.Committer{Status: 1, Number: 2}
	committers[2] = block.Committer{Status: 1, Number: 3}
	committers[3] = block.Committer{Status: 1, Number: 4}

	sigs[0] = tValSigner2.SignData(sb)
	sigs[1] = tValSigner3.SignData(sb)
	sigs[2] = tValSigner4.SignData(sb)
	sigs[3] = signer.SignData(sb)
	c1 := block.NewCommit(b1.Hash(), 3, committers, crypto.Aggregate(sigs))

	require.NoError(t, st1.CommitBlock(13, *b1, *c1))
	require.NoError(t, tState2.CommitBlock(13, *b1, *c1))
}

func TestValidateBlockTime(t *testing.T) {
	setup(t)

	fmt.Printf("BlockTimeInSecond: %d\n", tState1.params.BlockTimeInSecond)
	roundedNow := util.RoundNow(10)
	tState1.lastBlockTime = roundedNow.Add(-1 * time.Minute)

	// Time not rounded
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-15*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(-5*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(5*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(15*time.Second)))

	// Too early
	assert.Error(t, tState1.validateBlockTime(tState1.lastBlockTime.Add(-20*time.Second)))
	assert.Error(t, tState1.validateBlockTime(tState1.lastBlockTime.Add(-10*time.Second)))
	assert.Error(t, tState1.validateBlockTime(tState1.lastBlockTime))

	// Ok
	assert.NoError(t, tState1.validateBlockTime(roundedNow.Add(10*time.Second)))
	assert.NoError(t, tState1.validateBlockTime(roundedNow.Add(20*time.Second)))

	// Too late
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(30*time.Second)))
	assert.Error(t, tState1.validateBlockTime(roundedNow.Add(40*time.Second)))
}

func TestInvalidBlockVersion(t *testing.T) {
	setup(t)

	tState1.params.BlockVersion = 2
	b, _ := tState1.ProposeBlock(0)
	assert.Error(t, tState2.ValidateBlock(*b))
}

func TestInvalidBlockTime(t *testing.T) {
	setup(t)
	moveToNextHeightForAllStates(t)

	validBlock, _ := makeBlockAndCommit(t, 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	invalidBlock := block.MakeBlock(
		validBlock.Header().Version(),
		validBlock.Header().Time().Add(30*time.Second),
		validBlock.TxIDs(),
		validBlock.Header().LastBlockHash(),
		validBlock.Header().CommitteeHash(),
		validBlock.Header().StateHash(),
		validBlock.Header().LastReceiptsHash(),
		validBlock.LastCommit(),
		validBlock.Header().SortitionSeed(),
		validBlock.Header().ProposerAddress())

	assert.NoError(t, tState1.ValidateBlock(validBlock))
	assert.Error(t, tState1.ValidateBlock(invalidBlock))

}
