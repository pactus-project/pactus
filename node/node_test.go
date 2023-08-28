package node

import (
	"testing"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunningNode(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(util.Now(),
		map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc},
		[]*validator.Validator{val}, param.DefaultParams())
	conf := config.DefaultConfig()
	conf.Network.Listens = []string{"/ip4/0.0.0.0/tcp/0"}
	conf.GRPC.Enable = false
	conf.HTTP.Enable = false
	conf.Store.Path = util.TempDirPath()
	conf.Network.EnableRelay = false
	conf.Network.NetworkKey = util.TempFilePath()

	signers := []crypto.Signer{ts.RandSigner(), ts.RandSigner()}
	rewardAddrs := []crypto.Address{ts.RandAddress(), ts.RandAddress()}
	n, err := NewNode(gen, conf, signers, rewardAddrs)

	require.NoError(t, err)
	assert.Equal(t, n.state.LastBlockHash(), hash.UndefHash)

	err = n.Start()
	require.NoError(t, err)
	n.Stop()
}
