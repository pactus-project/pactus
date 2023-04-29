package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/node/config"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/genesis"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tSigners []crypto.Signer
var tConfigs []*config.Config
var tNodes []*node.Node
var tGRPCAddress = "0.0.0.0:1337"
var tGenDoc *genesis.Genesis
var tGRPC *grpc.ClientConn
var tBlockchain pactus.BlockchainClient
var tTransaction pactus.TransactionClient
var tNetwork pactus.NetworkClient
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
		tConfigs[i].Sync.NodeNetwork = false
		tConfigs[i].Sync.Firewall.Enabled = false
		tConfigs[i].Network.EnableMdns = true
		tConfigs[i].Network.NodeKey = util.TempFilePath()
		tConfigs[i].Network.Listens = []string{"/ip4/127.0.0.1/tcp/0"}
		tConfigs[i].Network.Bootstrap.Period = 10 * time.Second
		tConfigs[i].Network.Bootstrap.MinThreshold = 3
		tConfigs[i].HTTP.Enable = false
		tConfigs[i].GRPC.Enable = false

		sync.LatestBlockInterval = 10

		if i == 0 {
			tConfigs[i].Sync.NodeNetwork = true
			tConfigs[i].GRPC.Enable = true
			tConfigs[i].GRPC.Listen = tGRPCAddress
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
		tNodes[i], _ = node.NewNode(tGenDoc, tConfigs[i],
			[]crypto.Signer{tSigners[i]}, []crypto.Address{tSigners[i].Address()})
		if err := tNodes[i].Start(); err != nil {
			panic(fmt.Sprintf("Error on starting the node: %v", err))
		}
	}

	tCtx = context.Background()
	conn, err := grpc.DialContext(
		tCtx,
		tGRPCAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Errorf("failed to dial server: %w", err))
	}

	tGRPC = conn
	tBlockchain = pactus.NewBlockchainClient(conn)
	tTransaction = pactus.NewTransactionClient(conn)
	tNetwork = pactus.NewNetworkClient(conn)

	// Wait for some blocks
	waitForNewBlocks(8)

	fmt.Println("Running tests")
	exitCode := m.Run()

	// Running other nodes
	for i := tCommitteeSize; i < tTotalNodes; i++ {
		tNodes[i], _ = node.NewNode(tGenDoc, tConfigs[i],
			[]crypto.Signer{tSigners[i]}, []crypto.Address{tSigners[i].Address()})
		if err := tNodes[i].Start(); err != nil {
			panic(fmt.Sprintf("Error on starting the node: %v", err))
		}
	}

	// Commit more blocks, then new nodes can catch up and send sortition transactions
	waitForNewBlocks(20)

	// Check if sortition worked or not?
	block := lastBlock()
	cert := block.PrevCert
	// TODO: determine why the certificate can be nil?
	if cert != nil && len(cert.Committers) == 4 {
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
