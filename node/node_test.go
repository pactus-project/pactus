package node

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet"
	wltmgr "github.com/pactus-project/pactus/wallet/manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunningNode(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(time.Now(),
		map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc},
		[]*validator.Validator{val}, genesis.DefaultGenesisParams())
	conf := config.DefaultConfigForChain(genesis.Mainnet)
	// Prevent log from messing the workspace
	conf.Logger.Filename = util.TempFilePath()
	conf.GRPC.Enable = true
	conf.GRPC.Listen = "0.0.0.0:0"
	conf.HTML.Enable = true
	conf.HTML.Listen = "0.0.0.0:0"
	conf.HTTP.Enable = true
	conf.HTTP.Listen = "0.0.0.0:0"
	conf.JSONRPC.Enable = true
	conf.JSONRPC.Listen = "0.0.0.0:0"
	conf.Store.Path = util.TempDirPath()
	conf.Network.EnableRelay = false
	conf.Network.NetworkKey = util.TempFilePath()
	conf.Network.PeerStorePath = util.TempFilePath()
	conf.WalletManager.WalletsDir = util.TempDirPath()

	walletPath := filepath.Join(conf.WalletManager.WalletsDir, "default_wallet")
	mnemonic, _ := wallet.GenerateMnemonic(128)
	password := "password"
	wlt, err := wallet.Create(t.Context(), walletPath, mnemonic, password, genesis.Mainnet)
	require.NoError(t, err)
	wlt.Close()

	valKeys := []*bls.ValidatorKey{ts.RandValKey(), ts.RandValKey()}
	rewardAddrs := []crypto.Address{ts.RandAccAddress(), ts.RandAccAddress()}
	node, err := NewNode(t.Context(), gen, conf, valKeys, rewardAddrs)
	require.NoError(t, err)

	assert.True(t, conf.Sync.Services.IsFullNode())
	assert.True(t, conf.Sync.Services.IsPrunedNode())
	assert.Equal(t, hash.UndefHash, node.state.LastBlockHash())

	err = node.Start()
	require.NoError(t, err)

	consHeight, _ := node.ConsManager().HeightRound()
	assert.Equal(t, types.Height(1), consHeight)

	lastBlockTime := node.State().LastBlockTime()
	assert.Equal(t, gen.GenesisTime(), lastBlockTime)

	syncSelfID := node.Sync().SelfID()
	netSelfID := node.Network().SelfID()
	assert.Equal(t, syncSelfID, netSelfID)

	assert.NotEmpty(t, node.GRPC().Address())

	wallets, err := node.WalletManager().ListWallets()
	require.NoError(t, err)
	assert.NotEmpty(t, wallets)

	node.Stop()
}

func TestNewNodeReleasesStoreOnError(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(time.Now(),
		map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc},
		[]*validator.Validator{val}, genesis.DefaultGenesisParams())
	conf := config.DefaultConfigForChain(genesis.Mainnet)
	conf.Logger.Filename = util.TempFilePath()
	conf.Store.Path = util.TempDirPath()
	conf.Network.EnableRelay = false
	// The network key can't be read, because its path is a directory.
	conf.Network.NetworkKey = util.TempDirPath()
	conf.Network.PeerStorePath = util.TempFilePath()
	conf.WalletManager.WalletsDir = util.TempDirPath()

	valKeys := []*bls.ValidatorKey{ts.RandValKey()}
	rewardAddrs := []crypto.Address{ts.RandAccAddress()}

	_, err := NewNode(t.Context(), gen, conf, valKeys, rewardAddrs)
	require.Error(t, err)

	// The store must be released, so it can be opened again.
	str, err := store.NewStore(conf.Store)
	require.NoError(t, err, "the store is still locked after NewNode has failed")
	str.Close()
}

func TestStartStopsWaitingForGenesisOnCancel(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(time.Now().Add(time.Hour),
		map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc},
		[]*validator.Validator{val}, genesis.DefaultGenesisParams())
	conf := config.DefaultConfigForChain(genesis.Mainnet)
	conf.Logger.Filename = util.TempFilePath()
	conf.Store.Path = util.TempDirPath()
	conf.Network.EnableRelay = false
	conf.Network.NetworkKey = util.TempFilePath()
	conf.Network.PeerStorePath = util.TempFilePath()
	conf.WalletManager.WalletsDir = util.TempDirPath()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	valKeys := []*bls.ValidatorKey{ts.RandValKey()}
	rewardAddrs := []crypto.Address{ts.RandAccAddress()}
	node, err := NewNode(ctx, gen, conf, valKeys, rewardAddrs)
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() { done <- node.Start() }()

	cancel()

	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(10 * time.Second):
		t.Fatal("Start is still waiting for the genesis time after the context was canceled")
	}
}

// freeTCPPort returns a TCP port on the loopback interface that is not in use.
func freeTCPPort(t *testing.T) int {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = lis.Close() }()

	tcpAddr, ok := lis.Addr().(*net.TCPAddr)
	require.True(t, ok)

	return tcpAddr.Port
}

func TestNewNodeReleasesNetworkOnError(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(time.Now(),
		map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc},
		[]*validator.Validator{val}, genesis.DefaultGenesisParams())

	netPort := freeTCPPort(t)
	zmqPort := freeTCPPort(t)

	conf := config.DefaultConfigForChain(genesis.Mainnet)
	conf.Logger.Filename = util.TempFilePath()
	conf.Store.Path = util.TempDirPath()
	conf.Network.EnableRelay = false
	conf.Network.ListenAddrStrings = []string{fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", netPort)}
	conf.Network.NetworkKey = util.TempFilePath()
	conf.Network.PeerStorePath = util.TempFilePath()
	conf.ZeroMq.ZmqPubBlockInfo = fmt.Sprintf("tcp://127.0.0.1:%d", zmqPort)

	// The wallet manager is the last component to be created, and it can't
	// work on a file.
	walletsFile := util.TempFilePath()
	require.NoError(t, os.WriteFile(walletsFile, []byte("not a directory"), 0o600))
	conf.WalletManager.WalletsDir = walletsFile

	valKeys := []*bls.ValidatorKey{ts.RandValKey()}
	rewardAddrs := []crypto.Address{ts.RandAccAddress()}

	_, err := NewNode(t.Context(), gen, conf, valKeys, rewardAddrs)
	require.Error(t, err)

	// The network and the ZeroMQ sockets must be released: nothing is
	// listening on their ports anymore.
	for _, port := range []int{netPort, zmqPort} {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		require.NoError(t, err, "port %d is still in use after NewNode has failed", port)
		_ = lis.Close()
	}
}

type closeHookState struct {
	state.State

	onClose func()
}

func (s *closeHookState) Close() {
	s.onClose()
	s.State.Close()
}

type stopHookWalletManager struct {
	wltmgr.WalletManager

	onStop func()
}

func (m *stopHookWalletManager) Stop() {
	m.onStop()
	m.WalletManager.Stop()
}

// loopbackAddr returns the loopback address with the port of the given listen address.
func loopbackAddr(t *testing.T, listenAddr string) string {
	t.Helper()

	_, port, err := net.SplitHostPort(listenAddr)
	require.NoError(t, err)

	return net.JoinHostPort("127.0.0.1", port)
}

func TestStopShutsDownTheAPIBeforeClosingTheState(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	val := validator.NewValidator(pub, 0)
	gen := genesis.MakeGenesis(time.Now(),
		map[crypto.Address]*account.Account{crypto.TreasuryAddress: acc},
		[]*validator.Validator{val}, genesis.DefaultGenesisParams())
	conf := config.DefaultConfigForChain(genesis.Mainnet)
	conf.Logger.Filename = util.TempFilePath()
	conf.GRPC.Enable = true
	conf.GRPC.Listen = "127.0.0.1:0"
	conf.JSONRPC.Enable = true
	conf.JSONRPC.Listen = "127.0.0.1:0"
	conf.Store.Path = util.TempDirPath()
	conf.Network.EnableRelay = false
	conf.Network.NetworkKey = util.TempFilePath()
	conf.Network.PeerStorePath = util.TempFilePath()
	conf.WalletManager.WalletsDir = util.TempDirPath()

	valKeys := []*bls.ValidatorKey{ts.RandValKey()}
	rewardAddrs := []crypto.Address{ts.RandAccAddress()}
	node, err := NewNode(t.Context(), gen, conf, valKeys, rewardAddrs)
	require.NoError(t, err)
	require.NoError(t, node.Start())

	servers := map[string]string{
		"gRPC":     loopbackAddr(t, node.GRPC().Address()),
		"JSON-RPC": loopbackAddr(t, node.jsonrpc.Address()),
	}
	for name, addr := range servers {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		require.NoError(t, err, "%s server is not running", name)
		_ = conn.Close()
	}

	events := make([]string, 0)
	stillServing := make([]string, 0)
	node.walletMgr = &stopHookWalletManager{
		WalletManager: node.walletMgr,
		onStop:        func() { events = append(events, "wallet manager") },
	}
	node.state = &closeHookState{
		State: node.state,
		onClose: func() {
			events = append(events, "state")
			// Requests must not reach a state that is being closed.
			for name, addr := range servers {
				conn, err := net.DialTimeout("tcp", addr, time.Second)
				if err == nil {
					_ = conn.Close()
					stillServing = append(stillServing, name)
				}
			}
		},
	}

	node.Stop()

	assert.Empty(t, stillServing, "these servers were still running when the state was closed")
	assert.Equal(t, []string{"wallet manager", "state"}, events)
}
