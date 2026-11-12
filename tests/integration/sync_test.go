package integration

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/stretchr/testify/require"
)

// A node that joins late downloads the blocks from the other nodes. Its sessions are small,
// so the window of sessions moves many times before the node catches up.
func TestLateNodeCatchesUp(t *testing.T) {
	target := lastHeight()

	conf := config.DefaultConfigForChain(genesis.Mainnet)
	conf.Logger.Filename = util.TempFilePath()
	conf.Logger.Levels["_sync"] = "info"
	conf.Store.Path = util.TempDirPath()
	conf.Sync.Firewall.BannedNets = make([]string, 0)
	conf.Sync.BlockPerSession = 3
	conf.Sync.MaxSessions = 2
	conf.Network.EnableMdns = true
	conf.Network.EnableRelay = false
	conf.Network.DefaultBootstrapAddrStrings = []string{}
	conf.Network.BootstrapAddrStrings = []string{
		fmt.Sprintf("%v/p2p/%v", tNodes[tNodeIdx1].Network().HostAddrs()[0], tNodes[tNodeIdx1].Network().SelfID()),
	}
	conf.Network.ForcePrivateNetwork = true
	conf.Network.NetworkKey = util.TempFilePath()
	conf.Network.NetworkName = "test"
	conf.Network.ListenAddrStrings = []string{"/ip4/127.0.0.1/tcp/0"}
	conf.Network.MaxConns = 32
	conf.Network.PeerStorePath = util.TempFilePath()
	conf.HTML.Enable = false
	conf.GRPC.Enable = false
	conf.WalletManager.WalletsDir = util.TempDirPath()

	walletPath := filepath.Join(conf.WalletManager.WalletsDir, "default_wallet")
	mnemonic, err := wallet.GenerateMnemonic(128)
	require.NoError(t, err)
	wlt, err := wallet.Create(t.Context(), walletPath, mnemonic, "password", genesis.Mainnet)
	require.NoError(t, err)
	wlt.Close()

	// The validator of this node is not in the committee.
	key, err := bls.KeyGen(hash.CalcHash([]byte("late node")).Bytes(), nil)
	require.NoError(t, err)
	valKey := bls.NewValidatorKey(key)

	lateNode, err := node.NewNode(tCtx, tGenDoc, conf, []*bls.ValidatorKey{valKey},
		[]crypto.Address{valKey.PublicKey().AccountAddress()})
	require.NoError(t, err)
	require.NoError(t, lateNode.Start())
	defer lateNode.Stop()

	deadline := time.Now().Add(2 * time.Minute)
	for lateNode.State().LastBlockHeight() < target {
		require.True(t, time.Now().Before(deadline),
			"the late node is at height %d, the network was at %d", lateNode.State().LastBlockHeight(), target)
		time.Sleep(time.Second)
	}

	// The blocks are the ones of the network.
	ours, err := lateNode.State().CommittedBlock(target)
	require.NoError(t, err)
	theirs, err := tNodes[tNodeIdx1].State().CommittedBlock(target)
	require.NoError(t, err)
	require.Equal(t, theirs.BlockHash, ours.BlockHash)
}
