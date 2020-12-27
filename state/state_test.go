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
	logger.InitLogger(logger.TestConfig())

	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()

	tValSigner1 = crypto.NewSigner(priv1)
	tValSigner2 = crypto.NewSigner(priv2)
	tValSigner3 = crypto.NewSigner(priv3)
	tValSigner4 = crypto.NewSigner(priv4)

	tGenTime = util.Now()
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

func proposeAndSignBlock(t *testing.T, st *state) (block.Block, block.Commit) {
	b := st.ProposeBlock()
	committers := make([]block.Committer, 1)
	sb := vote.CommitSignBytes(b.Hash(), 0)
	committers[0] = block.Committer{Status: 1, Address: tValSigner1.Address()}
	sig := tValSigner1.Sign(sb)

	c := block.NewCommit(0, committers, *sig)
	return b, *c
}

func makeCommitAndSign(t *testing.T, blockHash crypto.Hash, round int, signers ...crypto.Signer) block.Commit {
	committers := make([]block.Committer, 4)
	sigs := make([]*crypto.Signature, len(signers))
	sb := vote.CommitSignBytes(blockHash, round)
	committers[0] = block.Committer{Status: 0, Address: tValSigner1.Address()}
	committers[1] = block.Committer{Status: 0, Address: tValSigner2.Address()}
	committers[2] = block.Committer{Status: 0, Address: tValSigner3.Address()}
	committers[3] = block.Committer{Status: 0, Address: tValSigner4.Address()}

	for i, s := range signers {
		if s.Address().EqualsTo(tValSigner1.Address()) {
			committers[0] = block.Committer{Status: 1, Address: s.Address()}
		}

		if s.Address().EqualsTo(tValSigner2.Address()) {
			committers[1] = block.Committer{Status: 1, Address: s.Address()}
		}

		if s.Address().EqualsTo(tValSigner3.Address()) {
			committers[2] = block.Committer{Status: 1, Address: s.Address()}
		}

		if s.Address().EqualsTo(tValSigner4.Address()) {
			committers[3] = block.Committer{Status: 1, Address: s.Address()}
		}

		sigs[i] = s.Sign(sb)
	}
	return *block.NewCommit(round, committers, crypto.Aggregate(sigs))
}

func TestProposeBlockAndValidation(t *testing.T) {
	st := setupStatewithOneValidator(t)

	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
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
	invBlock, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, st.ApplyBlock(1, *invBlock, c1))
	assert.Error(t, st.ApplyBlock(2, b1, c1))
	assert.NoError(t, st.ApplyBlock(1, b1, c1))
}

func TestCommitSandbox(t *testing.T) {

	t.Run("Commit new account", func(t *testing.T) {
		st := setupStatewithFourValidators(t, tValSigner1)

		addr, _, _ := crypto.GenerateTestKeyPair()
		newAcc := st.executionSandbox.MakeNewAccount(addr)
		newAcc.AddToBalance(1)
		st.commitSandbox(0)

		assert.True(t, st.store.HasAccount(addr))
	})

	t.Run("Commit new validator", func(t *testing.T) {
		st := setupStatewithFourValidators(t, tValSigner1)

		addr, pub, _ := crypto.GenerateTestKeyPair()
		newVal := st.executionSandbox.MakeNewValidator(pub)
		newVal.AddToStake(1)
		st.commitSandbox(0)

		assert.True(t, st.store.HasValidator(addr))
	})

	t.Run("Modify account", func(t *testing.T) {
		st := setupStatewithFourValidators(t, tValSigner1)

		acc := st.executionSandbox.Account(crypto.TreasuryAddress)
		acc.SubtractFromBalance(1)
		st.executionSandbox.UpdateAccount(acc)
		st.commitSandbox(0)

		acc1, _ := st.store.Account(crypto.TreasuryAddress)
		assert.Equal(t, acc1.Balance(), acc.Balance())
	})

	t.Run("Modify validator", func(t *testing.T) {
		st := setupStatewithFourValidators(t, tValSigner1)

		val := st.executionSandbox.Validator(tValSigner2.Address())
		val.AddToStake(1)
		st.executionSandbox.UpdateValidator(val)
		st.commitSandbox(0)

		val1, _ := st.store.Validator(tValSigner2.Address())
		assert.Equal(t, val1.Stake(), val.Stake())
	})

	t.Run("Move valset", func(t *testing.T) {
		st := setupStatewithFourValidators(t, tValSigner1)

		nextProposer := st.validatorSet.Proposer(1)

		st.commitSandbox(0)

		assert.Equal(t, st.validatorSet.Proposer(0).Address(), nextProposer.Address())
	})

	t.Run("Move valset next round", func(t *testing.T) {
		st := setupStatewithFourValidators(t, tValSigner1)

		nextNextProposer := st.validatorSet.Proposer(2)

		st.commitSandbox(1)

		assert.Equal(t, st.validatorSet.Proposer(0).Address(), nextNextProposer.Address())
	})
}

func TestUpdateLastCommit(t *testing.T) {
	st := setupStatewithFourValidators(t, tValSigner1)
	b := st.ProposeBlock()
	c1 := makeCommitAndSign(t, b.Hash(), 0, tValSigner1, tValSigner3, tValSigner4)
	c2 := makeCommitAndSign(t, b.Hash(), 0, tValSigner1, tValSigner2, tValSigner3, tValSigner4)

	st.lastCommit = &c1
	st.lastBlockHash = b.Hash()
	assert.NoError(t, st.UpdateLastCommit(&c1))
	assert.Equal(t, st.lastCommit.Hash(), c1.Hash())
	assert.NoError(t, st.UpdateLastCommit(&c2))
	assert.NoError(t, st.UpdateLastCommit(&c1))
	assert.Equal(t, st.lastCommit.Hash(), c2.Hash())
}
