package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tValKeys     [][]*bls.ValidatorKey
	tConfigs     []*config.Config
	tNodes       []*node.Node
	tGRPCAddress = "127.0.0.1:1337"
	tGenDoc      *genesis.Genesis
	tGRPC        *grpc.ClientConn
	tBlockchain  pactus.BlockchainClient
	tTransaction pactus.TransactionClient
	tNetwork     pactus.NetworkClient
	tCtx         context.Context
)

const (
	tNodeIdx1      = 0
	tNodeIdx2      = 1
	tNodeIdx3      = 2
	tNodeIdx4      = 3
	tNodeIdx5      = 4
	tNodeIdx6      = 5
	tTotalNodes    = 6 // each node has 3 validators
	tCommitteeSize = 11
)

func TestMain(m *testing.M) {
	// Prevent log from messing the workspace
	logger.LogFilename = util.TempFilePath()

	tValKeys = make([][]*bls.ValidatorKey, tTotalNodes)
	tConfigs = make([]*config.Config, tTotalNodes)
	tNodes = make([]*node.Node, tTotalNodes)

	ikm := hash.CalcHash([]byte{})
	for i := 0; i < tTotalNodes; i++ {
		ikm = hash.CalcHash(ikm.Bytes())
		key0, _ := bls.KeyGen(ikm.Bytes(), nil)
		ikm = hash.CalcHash(ikm.Bytes())
		key1, _ := bls.KeyGen(ikm.Bytes(), nil)
		ikm = hash.CalcHash(ikm.Bytes())
		key2, _ := bls.KeyGen(ikm.Bytes(), nil)

		tValKeys[i] = make([]*bls.ValidatorKey, 3)
		tValKeys[i][0] = bls.NewValidatorKey(key0)
		tValKeys[i][1] = bls.NewValidatorKey(key1)
		tValKeys[i][2] = bls.NewValidatorKey(key2)
		tConfigs[i] = config.DefaultConfigMainnet()

		tConfigs[i].TxPool.MinFeePAC = 0.000001
		tConfigs[i].Store.Path = util.TempDirPath()
		tConfigs[i].Consensus.ChangeProposerTimeout = 4 * time.Second
		tConfigs[i].Consensus.ChangeProposerDelta = 0
		tConfigs[i].Logger.Levels["default"] = "warn"
		tConfigs[i].Logger.Levels["_state"] = "warn"
		tConfigs[i].Logger.Levels["_sync"] = "debug"
		tConfigs[i].Logger.Levels["_consensus"] = "debug"
		tConfigs[i].Logger.Levels["_network"] = "debug"
		tConfigs[i].Logger.Levels["_pool"] = "warn"
		tConfigs[i].Sync.NodeNetwork = false
		tConfigs[i].Sync.Firewall.Enabled = false
		tConfigs[i].Sync.LatestBlockInterval = 10
		tConfigs[i].Network.EnableMdns = true
		tConfigs[i].Network.EnableRelay = false
		tConfigs[i].Network.DefaultBootstrapAddrStrings = []string{}
		tConfigs[i].Network.BootstrapAddrStrings = []string{}
		tConfigs[i].Network.ForcePrivateNetwork = true
		tConfigs[i].Network.NetworkKey = util.TempFilePath()
		tConfigs[i].Network.NetworkName = "test"
		tConfigs[i].Network.ListenAddrStrings = []string{"/ip4/127.0.0.1/tcp/0", "/ip4/127.0.0.1/udp/0/quic-v1"}
		tConfigs[i].Network.MaxConns = 32
		tConfigs[i].HTTP.Enable = false
		tConfigs[i].GRPC.Enable = false

		if i == 0 {
			tConfigs[i].Sync.NodeNetwork = true
			tConfigs[i].GRPC.Enable = true
			tConfigs[i].GRPC.Listen = tGRPCAddress
		}
		fmt.Printf("Node %d created.\n", i+1)
	}

	acc1 := account.NewAccount(0)
	acc1.AddToBalance(21 * 1e14)
	key, _ := bls.KeyGen(ikm.Bytes(), nil)
	acc2 := account.NewAccount(1)
	acc2.AddToBalance(21 * 1e14)

	accs := map[crypto.Address]*account.Account{
		crypto.TreasuryAddress:                 acc1,
		key.PublicKeyNative().AccountAddress(): acc2,
	}

	vals := make([]*validator.Validator, 6)
	vals[0] = validator.NewValidator(tValKeys[tNodeIdx1][0].PublicKey(), 0)
	vals[1] = validator.NewValidator(tValKeys[tNodeIdx2][0].PublicKey(), 1)
	vals[2] = validator.NewValidator(tValKeys[tNodeIdx3][0].PublicKey(), 2)
	vals[3] = validator.NewValidator(tValKeys[tNodeIdx4][0].PublicKey(), 3)
	vals[4] = validator.NewValidator(tValKeys[tNodeIdx5][0].PublicKey(), 4)
	vals[5] = validator.NewValidator(tValKeys[tNodeIdx6][0].PublicKey(), 5)
	params := param.DefaultParams()
	params.MinimumStake = 1000
	params.BlockIntervalInSecond = 2
	params.BondInterval = 8
	params.CommitteeSize = tCommitteeSize
	params.TransactionToLiveInterval = 8
	tGenDoc = genesis.MakeGenesis(util.Now(), accs, vals, params)

	for i := 0; i < tTotalNodes; i++ {
		tNodes[i], _ = node.NewNode(
			tGenDoc, tConfigs[i],
			tValKeys[i],
			[]crypto.Address{
				tValKeys[i][0].PublicKey().AccountAddress(),
				tValKeys[i][1].PublicKey().AccountAddress(),
				tValKeys[i][2].PublicKey().AccountAddress(),
			})

		if err := tNodes[i].Start(); err != nil {
			panic(fmt.Sprintf("Error on starting the node: %v", err))
		}

		if i == 0 {
			// Set bootstrap address for better connectivity
			bootstrapAddr := fmt.Sprintf("%v/p2p/%v",
				tNodes[i].Network().HostAddrs()[0], tNodes[i].Network().SelfID())
			fmt.Println("Bootstrap address is: " + bootstrapAddr)

			tConfigs[tNodeIdx2].Network.BootstrapAddrStrings = []string{bootstrapAddr}
			tConfigs[tNodeIdx3].Network.BootstrapAddrStrings = []string{bootstrapAddr}
			tConfigs[tNodeIdx4].Network.BootstrapAddrStrings = []string{bootstrapAddr}
		}

		time.Sleep(1 * time.Second)
	}

	tCtx = context.Background()
	conn, err := grpc.NewClient(
		tGRPCAddress,
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
	fmt.Println("Waiting to commit some blocks...")
	waitForNewBlocks(8)

	fmt.Println("Running tests...")
	exitCode := m.Run()
	// Commit more blocks, then new nodes can catch up and send sortition transactions

	fmt.Println("Waiting to commit some blocks...")
	waitForNewBlocks(20)

	// Check if sortition worked or not?
	block := lastBlock()
	cert := block.PrevCert
	if block.Height == 1 {
		panic("block height should be greater than 1")
	}
	if len(cert.Committers) == 7 {
		panic("Sortition didn't work")
	}

	// Lets shutdown the nodes
	tCtx.Done()
	for i := 0; i < tTotalNodes; i++ {
		tNodes[i].Stop()
	}

	s, _ := store.NewStore(tConfigs[tNodeIdx1].Store)
	total := amount.Amount(0)
	s.IterateAccounts(func(_ crypto.Address, acc *account.Account) bool {
		total += acc.Balance()

		return false
	})

	s.IterateValidators(func(v *validator.Validator) bool {
		total += v.Stake()

		return false
	})
	if total != tGenDoc.TotalSupply() {
		panic(fmt.Sprintf("Some coins missed: %v", tGenDoc.TotalSupply()-total))
	}

	os.Exit(exitCode)
}
