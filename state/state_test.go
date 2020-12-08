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
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

var mockTxPool *txpool.MockTxPool
var gen *genesis.Genesis
var valSigner crypto.Signer

func init() {
	mockTxPool = txpool.NewMockTxPool()

	_, pb, priv := crypto.GenerateTestKeyPair()
	acc := account.NewAccount(crypto.MintbaseAddress)
	acc.SetBalance(21000000000000)
	val := validator.NewValidator(pb, 1)
	gen = genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val})
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
	st, _ := mockState(t, &valSigner)
	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}

func propsoeAndSignBlock(t *testing.T, st *state) (block.Block, block.Commit) {
	addr := valSigner.Address()
	b := st.ProposeBlock()
	v := vote.NewPrecommit(1, 0, b.Hash(), addr)
	sig := valSigner.Sign(v.SignBytes())
	c := block.NewCommit(0, []block.Commiter{block.Commiter{Signed: true, Address: addr}}, *sig)

	return b, *c
}

func TestLoadState(t *testing.T) {
	st1, _ := mockState(t, &valSigner)
	st2, _ := mockState(t, &valSigner)

	for i := 0; i < 10; i++ {
		b, c := propsoeAndSignBlock(t, st1)

		assert.NoError(t, st1.ApplyBlock(i+1, b, c))
		assert.NoError(t, st2.ApplyBlock(i+1, b, c))
	}

	// Propose second block
	b2 := st1.ProposeBlock()

	// Load state and propose second block
	st2.tryLoadLastInfo()
	b22 := st2.ProposeBlock()

	assert.Equal(t, b2.Hash(), b22.Hash())
}

func TestBlockSubsidy(t *testing.T) {
	interval := 210000
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy(1, 210000))
	assert.Equal(t, int64(5*1e8), calcBlockSubsidy((1*interval)-1, 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((1*interval), 210000))
	assert.Equal(t, int64(2.5*1e8), calcBlockSubsidy((2*interval)-1, 210000))
	assert.Equal(t, int64(1.25*1e8), calcBlockSubsidy((2*interval), 210000))
}

func TestApplyBlocks(t *testing.T) {
	st, _ := mockState(t, &valSigner)
	b1, c1 := propsoeAndSignBlock(t, st)
	invBlock, _ := block.GenerateTestBlock(nil)
	assert.Error(t, st.ApplyBlock(1, invBlock, c1))
	assert.Error(t, st.ApplyBlock(2, b1, c1))

}
