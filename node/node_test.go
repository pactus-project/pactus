package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/node/config"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/types/genesis"
	"github.com/zarbchain/zarb-go/types/param"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util"
)

func TestRunningNode(t *testing.T) {
	pub, pv := bls.GenerateTestKeyPair()
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(util.Now(), []*account.Account{acc}, []*validator.Validator{val}, param.DefaultParams())
	conf := config.DefaultConfig()
	conf.Network.Listens = []string{"/ip4/0.0.0.0/tcp/0"}
	conf.GRPC.Enable = false
	conf.Capnp.Enable = false
	conf.HTTP.Enable = false
	conf.Store.Path = util.TempDirPath()
	conf.Network.NodeKey = util.TempFilePath()

	signer := crypto.NewSigner(pv)
	n, err := NewNode(gen, conf, signer)

	require.NoError(t, err)
	assert.Equal(t, n.state.LastBlockHash(), hash.UndefHash)

	err = n.Start()
	require.NoError(t, err)
	n.Stop()
}
