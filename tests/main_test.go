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
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/node"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync"
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
var tSequences map[crypto.Address]int32

const tNodeIdx1 = 0
const tNodeIdx2 = 1
const tNodeIdx3 = 2
const tNodeIdx4 = 3
const tTotalNodes = 8
const tCommitteeSize = 7

func incSequence(addr crypto.Address) {
	tSequences[addr] = tSequences[addr] + 1
}

func getSequence(addr crypto.Address) int32 {
	return tSequences[addr]
}

func TestMain(m *testing.M) {
	tSigners = make([]crypto.Signer, tTotalNodes)
	tConfigs = make([]*config.Config, tTotalNodes)
	tNodes = make([]*node.Node, tTotalNodes)
	tSequences = make(map[crypto.Address]int32)

	for i := 0; i < tTotalNodes; i++ {
		pub, prv := bls.GenerateTestKeyPair()
		tSigners[i] = crypto.NewSigner(prv)
		tConfigs[i] = config.DefaultConfig()

		tConfigs[i].Store.Path = util.TempDirPath()
		tConfigs[i].Consensus.ChangeProposerTimeout = 4 * time.Second
		tConfigs[i].Logger.Levels["default"] = "warning"
		tConfigs[i].Logger.Levels["_state"] = "info"
		tConfigs[i].Logger.Levels["_sync"] = "error"
		tConfigs[i].Logger.Levels["_consensus"] = "error"
		tConfigs[i].Logger.Levels["_network"] = "error"
		tConfigs[i].Logger.Levels["_pool"] = "error"
		tConfigs[i].Sync.CacheSize = 1000
		tConfigs[i].Sync.StartingTimeout = 0
		tConfigs[i].Sync.NodeNetwork = false
		tConfigs[i].Sync.Firewall.Enabled = false
		tConfigs[i].Network.NodeKeyFile = util.TempFilePath()
		tConfigs[i].Network.ListenAddress = []string{fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", 32125+i)}
		tConfigs[i].Network.Bootstrap.Addresses = []string{"/ip4/127.0.0.1/tcp/32125/p2p/12D3KooWCKKGMMGDhqRUZh6MnH2to6XUN9N2YPof4LrNNMe5Mbek"}
		tConfigs[i].Network.Bootstrap.Period = 10 * time.Second
		tConfigs[i].Network.Bootstrap.MinThreshold = 3
		tConfigs[i].HTTP.Enable = false
		tConfigs[i].GRPC.Enable = false
		tConfigs[i].Capnp.Enable = false

		sync.LatestBlockInterval = 10

		if i == 0 {
			tConfigs[i].Sync.NodeNetwork = true
			tConfigs[i].Capnp.Enable = true
			tConfigs[i].Capnp.Address = tCapnpAddress

			f, _ := os.Create(tConfigs[i].Network.NodeKeyFile)
			_, err := f.WriteString("08011240f22591817d8803e32525db7fc5cb9949d77c402e20867a6cac6b3ffb3dc643fb2521ef3c844a12eee79c275f19958999aeebb173496b67ea4a40f5d34b0a1355")
			if err != nil {
				panic(err)
			}
			f.Close()
		}
		fmt.Printf("Node %d address: %s\n", i+1, pub.Address())
	}

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)

	vals := make([]*validator.Validator, 4)
	vals[0] = validator.NewValidator(tSigners[tNodeIdx1].PublicKey().(*bls.PublicKey), 0)
	vals[1] = validator.NewValidator(tSigners[tNodeIdx2].PublicKey().(*bls.PublicKey), 1)
	vals[2] = validator.NewValidator(tSigners[tNodeIdx3].PublicKey().(*bls.PublicKey), 2)
	vals[3] = validator.NewValidator(tSigners[tNodeIdx4].PublicKey().(*bls.PublicKey), 3)
	params := param.DefaultParams()
	params.BlockTimeInSecond = 2
	params.BondInterval = 8
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
	waitForNewBlocks(8)

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
	waitForNewBlocks(20)

	// Check if sortition worked or not?
	b, _ := lastBlock().Block()
	cert, _ := b.PrevCert()
	committers, _ := cert.Committers()
	if committers.Len() == 4 {
		panic("Sortition didn't work")
	}

	// Let's shutdown the nodes
	tCtx.Done()
	for i := 0; i < tTotalNodes; i++ {
		tNodes[i].Stop()
	}

	s, _ := store.NewStore(tConfigs[tNodeIdx1].Store, 0)
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
