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

func TestProposeBlockAndValidation(t *testing.T) {
	st := setupStatewithOneValidator(t)

	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}

func proposeAndSignBlock(t *testing.T, st *state, signer crypto.Signer) (block.Block, block.Commit) {
	b := st.ProposeBlock()
	c := makeCommitAndSign(t, b.Hash(), signer)

	return b, c
}

func makeCommitAndSign(t *testing.T, blockHash crypto.Hash, signers ...crypto.Signer) block.Commit {
	committers := make([]block.Committer, len(signers))
	sigs := make([]*crypto.Signature, len(signers))
	for i, s := range signers {
		v := vote.NewPrecommit(-1, 0, blockHash, s.Address())

		committers[i] = block.Committer{Status: 1, Address: s.Address()}
		sigs[i] = s.Sign(v.SignBytes())
	}
	return *block.NewCommit(0, committers, crypto.Aggregate(sigs))
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

	b1, c1 := proposeAndSignBlock(t, st, tValSigner1)
	invBlock, _ := block.GenerateTestBlock(nil, nil)
	assert.Error(t, st.ApplyBlock(1, *invBlock, c1))
	assert.Error(t, st.ApplyBlock(2, b1, c1))
	assert.NoError(t, st.ApplyBlock(1, b1, c1))
}
