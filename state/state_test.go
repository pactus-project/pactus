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
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

func mockState(t *testing.T, pb crypto.PublicKey) (*state, crypto.Address) {
	addr := pb.Address()
	acc := account.NewAccount(crypto.MintbaseAddress)
	acc.SetBalance(21000000000000)
	val := validator.NewValidator(pb, 1)
	gen := genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val})
	loggerConfig := logger.TestConfig()
	logger.InitLogger(loggerConfig)
	stateConfig := TestConfig()
	txPoolConfig := txpool.TestConfig()
	txPool, err := txpool.NewTxPool(txPoolConfig, make(chan *message.Message, 10))
	require.NoError(t, err)
	st, err := LoadOrNewState(stateConfig, gen, val.Address(), txPool)
	require.NoError(t, err)
	s, _ := st.(*state)
	return s, addr
}

func TestBlockValidate(t *testing.T) {
	_, pb, _ := crypto.RandomKeyPair()

	st, _ := mockState(t, pb)
	block := st.ProposeBlock()
	err := st.ValidateBlock(block)
	require.NoError(t, err)
}

func TestReplayBlock(t *testing.T) {
	a, pb, pv := crypto.RandomKeyPair()

	st1, _ := mockState(t, pb)
	st2, _ := mockState(t, pb)

	// apply first block
	b1 := st1.ProposeBlock()
	v := vote.NewPrecommit(1, 0, b1.Hash(), a)
	sig1 := pv.Sign(v.SignBytes())
	c1 := block.NewCommit(0, []block.Commiter{block.Commiter{Signed: true, Address: a}}, *sig1)

	h := b1.TxHashes().Hashes()[0]
	st2.txPool.AppendTx(*st1.txPool.PendingTx(h))
	st1.ApplyBlock(1, b1, *c1)
	st2.ApplyBlock(1, b1, *c1)

	// Propose second block
	b2 := st1.ProposeBlock()

	// Load state and propose second block
	st2.tryLoadLastInfo()
	b22 := st2.ProposeBlock()

	assert.Equal(t, b2.Hash(), b22.Hash())
}
