package tests

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/node"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/www/capnp"
	"zombiezen.com/go/capnproto2/rpc"
)

var tSigners []crypto.Signer
var tConfigs []*config.Config
var tNodes []*node.Node
var tCapnpAddress = "0.0.0.0:1337"
var tGenDoc *genesis.Genesis
var tCapnpServer capnp.ZarbServer
var tCtx context.Context
var tSequences map[crypto.Address]int

const tNodeIdx1 = 0
const tNodeIdx2 = 1
const tNodeIdx3 = 2
const tNodeIdx4 = 3

func incSequence(t *testing.T, addr crypto.Address) {
	tSequences[addr] = tSequences[addr] + 1
}

func getSequence(t *testing.T, addr crypto.Address) int {
	return tSequences[addr]
}

func TestMain(m *testing.M) {
	nodeCount := 8
	committeeSize := 4
	blockTime := 2

	tSigners = make([]crypto.Signer, nodeCount)
	tConfigs = make([]*config.Config, nodeCount)
	tNodes = make([]*node.Node, nodeCount)
	tSequences = make(map[crypto.Address]int)

	for i := 0; i < nodeCount; i++ {
		addr, _, priv := crypto.GenerateTestKeyPair()
		tSigners[i] = crypto.NewSigner(priv)
		tConfigs[i] = config.DefaultConfig()
		tConfigs[i].Sync.StartingTimeout = 0
		tConfigs[i].State.Store.Path = util.TempDirPath()
		tConfigs[i].Network.NodeKeyFile = util.TempFilePath()
		if i == 0 {
			tConfigs[i].Capnp.Address = tCapnpAddress
			f, _ := os.Create(tConfigs[i].Network.NodeKeyFile)
			_, err := f.WriteString("08011240f22591817d8803e32525db7fc5cb9949d77c402e20867a6cac6b3ffb3dc643fb2521ef3c844a12eee79c275f19958999aeebb173496b67ea4a40f5d34b0a1355")
			if err != nil {
				panic(err)
			}
			f.Close()
		} else {
			tConfigs[i].Capnp.Enable = false
		}
		tConfigs[i].Http.Enable = false
		tConfigs[i].GRPC.Enable = false

		tConfigs[i].Consensus.ChangeProposerTimeout = 1 * time.Second

		tConfigs[i].Logger.Levels["default"] = "error"
		tConfigs[i].Logger.Levels["_state"] = "info"
		tConfigs[i].Logger.Levels["_sync"] = "error"
		tConfigs[i].Logger.Levels["_consensus"] = "debug"
		tConfigs[i].Logger.Levels["_pool"] = "error"

		tConfigs[i].TxPool.WaitingTimeout = 500 * time.Millisecond
		tConfigs[i].Sync.CacheSize = 1000
		tConfigs[i].Network.ListenAddress = []string{fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", 32125+i)}
		tConfigs[i].Network.Bootstrap.Addresses = []string{"/ip4/127.0.0.1/tcp/32125/p2p/12D3KooWCKKGMMGDhqRUZh6MnH2to6XUN9N2YPof4LrNNMe5Mbek"}

		fmt.Printf("Node %d address: %s\n", i+1, addr)
	}

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)

	vals := make([]*validator.Validator, 4)
	vals[0] = validator.NewValidator(tSigners[tNodeIdx1].PublicKey(), 0, 0)
	vals[1] = validator.NewValidator(tSigners[tNodeIdx2].PublicKey(), 1, 0)
	vals[2] = validator.NewValidator(tSigners[tNodeIdx3].PublicKey(), 2, 0)
	vals[3] = validator.NewValidator(tSigners[tNodeIdx4].PublicKey(), 3, 0)
	params := param.DefaultParams()
	params.BlockTimeInSecond = blockTime
	params.CommitteeSize = committeeSize
	params.TransactionToLiveInterval = 8
	tGenDoc = genesis.MakeGenesis(util.Now(), []*account.Account{acc}, vals, params)

	t := &testing.T{}
	for i := 0; i < nodeCount; i++ {
		tNodes[i], _ = node.NewNode(tGenDoc, tConfigs[i], tSigners[i])
		if err := tNodes[i].Start(); err != nil {
			panic(fmt.Sprintf("Error on starting the node: %v", err.Error()))
		}
		time.Sleep(500 * time.Millisecond)
	}

	c, _ := net.Dial("tcp", tCapnpAddress)

	tCtx = context.Background()
	conn := rpc.NewConn(rpc.StreamTransport(c))
	tCapnpServer = capnp.ZarbServer{Client: conn.Bootstrap(tCtx)}

	waitForNewBlock(t)
	waitForNewBlock(t)
	waitForNewBlock(t)
	waitForNewBlock(t)

	// These validators are not in the committee now.
	// Bond transactions are valid and they can enter the committee soon
	for i := committeeSize; i < nodeCount; i++ {
		amt := util.RandInt64(1000000 - 1) // fee is always 1000
		err := broadcastBondTransaction(t, tSigners[tNodeIdx2], tSigners[i].PublicKey(), amt, 1000)
		if err != nil {
			panic(fmt.Sprintf("Error on broadcasting transaction: %v", err))
		}
		fmt.Printf("Staking %v to %v\n", amt, tSigners[i].Address())
		incSequence(t, tSigners[tNodeIdx2].Address())
	}

	fmt.Println("Running tests")

	exitCode := m.Run()

	// Some more blocks
	for i := 0; i < 20; i++ {
		waitForNewBlock(t)
	}

	tCtx.Done()
	for i := 0; i < nodeCount; i++ {
		tNodes[i].Stop()
	}

	s, _ := store.NewStore(tConfigs[tNodeIdx1].State.Store)
	total := int64(0)
	s.IterateAccounts(func(a *account.Account) bool {
		total += a.Balance()
		return false
	})

	s.IterateValidators(func(v *validator.Validator) bool {
		total += v.Stake()
		return false
	})
	if total != int64(21*1e14) {
		panic(fmt.Sprintf("Some coins missed: %v", total-21*1e14))
	}

	os.Exit(exitCode)
}
