package node

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/validator"
)

func TestRunningNode(t *testing.T) {
	_, pb, pv := crypto.RandomKeyPair()
	acc := account.NewAccount(crypto.MintbaseAddress, 0)
	acc.SetBalance(21000000000000)
	val := validator.NewValidator(pb, 0, 0)
	gen := genesis.MakeGenesis("test", time.Now(), []*account.Account{acc}, []*validator.Validator{val}, 1)
	conf := config.TestConfig()
	signer := crypto.NewSigner(pv)
	n, err := NewNode(gen, conf, signer)

	assert.Equal(t, n.state.LastBlockHash(), crypto.UndefHash)

	require.NoError(t, err)
	err = n.Start()
	require.NoError(t, err)
	n.Stop()
}
