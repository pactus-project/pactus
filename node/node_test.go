package node

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunningNode(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// Prevent log from messing the workspace
	logger.LogFilename = util.TempFilePath()
	pub, _ := ts.RandBLSKeyPair()
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(time.Now(),
		map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc},
		[]*validator.Validator{val}, genesis.DefaultGenesisParams())
	conf := config.DefaultConfigMainnet()
	conf.GRPC.Enable = true
	conf.GRPC.Listen = "0.0.0.0:0"
	conf.Rest.Enable = true
	conf.Rest.Listen = "0.0.0.0:0"
	conf.HTTP.Enable = true
	conf.HTTP.Listen = "0.0.0.0:0"
	conf.JSONRPC.Enable = true
	conf.JSONRPC.Listen = "0.0.0.0:0"
	conf.Store.Path = util.TempDirPath()
	conf.Network.EnableRelay = false
	conf.Network.NetworkKey = util.TempFilePath()
	conf.Network.PeerStorePath = util.TempFilePath()

	valKeys := []*bls.ValidatorKey{ts.RandValKey(), ts.RandValKey()}
	rewardAddrs := []crypto.Address{ts.RandAccAddress(), ts.RandAccAddress()}
	node, err := NewNode(gen, conf, valKeys, rewardAddrs)
	assert.True(t, conf.Sync.Services.IsFullNode())
	assert.True(t, conf.Sync.Services.IsPrunedNode())

	require.NoError(t, err)
	assert.Equal(t, hash.UndefHash, node.state.LastBlockHash())

	err = node.Start()
	require.NoError(t, err)

	consHeight, _ := node.ConsManager().HeightRound()
	assert.Equal(t, uint32(1), consHeight)

	lastBlockTime := node.State().LastBlockTime()
	assert.Equal(t, gen.GenesisTime(), lastBlockTime)

	syncSelfID := node.Sync().SelfID()
	netSelfID := node.Network().SelfID()
	assert.Equal(t, syncSelfID, netSelfID)

	assert.NotEmpty(t, node.GRPC().Address())

	node.Stop()
}
