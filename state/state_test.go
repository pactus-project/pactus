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

func TestBlockValidate(t *testing.T) {

	st, _ := mockState(t, &valSigner)
	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}

func TestReplayBlock(t *testing.T) {
	a, _, priv := crypto.RandomKeyPair()

	st1, _ := mockState(t, &valSigner)
	st2, _ := mockState(t, &valSigner)

	// apply first block
	b1 := st1.ProposeBlock()
	v := vote.NewPrecommit(1, 0, b1.Hash(), a)
	sig1 := priv.Sign(v.SignBytes())
	c1 := block.NewCommit(0, []block.Commiter{block.Commiter{Signed: true, Address: a}}, *sig1)

	st1.ApplyBlock(1, b1, *c1)
	st2.ApplyBlock(1, b1, *c1)

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
