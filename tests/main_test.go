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
const tTotalNodes = 8
const tCommitteeSize = 7

func incSequence(t *testing.T, addr crypto.Address) {
	tSequences[addr] = tSequences[addr] + 1
}

func getSequence(t *testing.T, addr crypto.Address) int {
	return tSequences[addr]
}

func TestMain(m *testing.M) {
	tSigners = make([]crypto.Signer, tTotalNodes)
	tConfigs = make([]*config.Config, tTotalNodes)
	tNodes = make([]*node.Node, tTotalNodes)
	tSequences = make(map[crypto.Address]int)

	for i := 0; i < tTotalNodes; i++ {
		addr, _, priv := crypto.GenerateTestKeyPair()
		tSigners[i] = crypto.NewSigner(priv)
		tConfigs[i] = config.DefaultConfig()

		tConfigs[i].Store.Path = util.TempDirPath()
		tConfigs[i].Consensus.ChangeProposerTimeout = 4 * time.Second
		tConfigs[i].Logger.Levels["default"] = "warning"
		tConfigs[i].Logger.Levels["_state"] = "info"
		tConfigs[i].Logger.Levels["_sync"] = "error"
		tConfigs[i].Logger.Levels["_consensus"] = "error"
		tConfigs[i].Logger.Levels["_pool"] = "error"
		tConfigs[i].TxPool.WaitingTimeout = 500 * time.Millisecond
		tConfigs[i].Sync.CacheSize = 1000
		tConfigs[i].Sync.RequestBlockInterval = 10
		tConfigs[i].Sync.StartingTimeout = 0
		tConfigs[i].Sync.InitialBlockDownload = false
		tConfigs[i].Network.NodeKeyFile = util.TempFilePath()
		tConfigs[i].Network.ListenAddress = []string{fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", 32125+i)}
		tConfigs[i].Network.Bootstrap.Addresses = []string{"/ip4/127.0.0.1/tcp/32125/p2p/12D3KooWCKKGMMGDhqRUZh6MnH2to6XUN9N2YPof4LrNNMe5Mbek"}
		tConfigs[i].Network.Bootstrap.Period = 10 * time.Second
		tConfigs[i].Network.Bootstrap.MinThreshold = 3
		tConfigs[i].Http.Enable = false
		tConfigs[i].GRPC.Enable = false
		tConfigs[i].Capnp.Enable = false

		if i == 0 {
			tConfigs[i].Sync.InitialBlockDownload = true
			tConfigs[i].Capnp.Enable = true
			tConfigs[i].Capnp.Address = tCapnpAddress

			f, _ := os.Create(tConfigs[i].Network.NodeKeyFile)
			_, err := f.WriteString("08011240f22591817d8803e32525db7fc5cb9949d77c402e20867a6cac6b3ffb3dc643fb2521ef3c844a12eee79c275f19958999aeebb173496b67ea4a40f5d34b0a1355")
			if err != nil {
				panic(err)
			}
			f.Close()
		}
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
	params.BlockTimeInSecond = 2
	params.CommitteeSize = tCommitteeSize
	params.TransactionToLiveInterval = 8
	tGenDoc = genesis.MakeGenesis(util.Now(), []*account.Account{acc}, vals, params)

	for i := 0; i < tCommitteeSize; i++ {
		tNodes[i], _ = node.NewNode(tGenDoc, tConfigs[i], tSigners[i])
		if err := tNodes[i].Start(); err != nil {
			panic(fmt.Sprintf("Error on starting the node: %v", err.Error()))
		}
	}

	c, _ := net.Dial("tcp", tCapnpAddress)
	tCtx = context.Background()
	conn := rpc.NewConn(rpc.StreamTransport(c))
	tCapnpServer = capnp.ZarbServer{Client: conn.Bootstrap(tCtx)}

	// Wait for some blocks
	for i := 0; i < 10; i++ {
		waitForNewBlock()
	}

	fmt.Println("Running tests")
	exitCode := m.Run()

	// Running other nodes
	for i := tCommitteeSize; i < tTotalNodes; i++ {
		tNodes[i], _ = node.NewNode(tGenDoc, tConfigs[i], tSigners[i])
		if err := tNodes[i].Start(); err != nil {
			panic(fmt.Sprintf("Error on starting the node: %v", err.Error()))
		}
	}

	// Commit more blocks, then new nodes can catch up and send sortition transactions
	for i := 0; i < 40; i++ {
		waitForNewBlock()
	}

	// Check if sortition worked or not?
	b := lastBlock()
	committers := b.LastCertificate().Committers()
	if len(committers) == 4 {
		panic("Sortition didn't work")
	}

	// Let's shutdown the nodes
	tCtx.Done()
	for i := 0; i < tTotalNodes; i++ {
		tNodes[i].Stop()
	}

	s, _ := store.NewStore(tConfigs[tNodeIdx1].Store)
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
