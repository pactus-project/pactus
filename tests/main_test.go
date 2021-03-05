package tests

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

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
var tCurlAddress = "0.0.0.0:1337"
var tCapnpAddress = "0.0.0.0:31337"
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
	nodeCount := 7
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
			tConfigs[i].Http.Address = tCurlAddress
			tConfigs[i].Capnp.Address = tCapnpAddress
		} else {
			tConfigs[i].Http.Enable = false
			tConfigs[i].Capnp.Enable = false
		}
		tConfigs[i].GRPC.Enable = false

		tConfigs[i].Logger.Levels["default"] = "info"
		tConfigs[i].Logger.Levels["_state"] = "info"
		tConfigs[i].Logger.Levels["_sync"] = "error"
		tConfigs[i].Logger.Levels["_consensus"] = "error"
		tConfigs[i].Logger.Levels["_txpool"] = "error"

		tConfigs[i].Sync.CacheSize = 1000
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
	}

	c, _ := net.Dial("tcp", tCapnpAddress)

	tCtx = context.Background()
	conn := rpc.NewConn(rpc.StreamTransport(c))
	tCapnpServer = capnp.ZarbServer{Client: conn.Bootstrap(tCtx)}

	waitForNewBlock(t)
	waitForNewBlock(t)
	waitForNewBlock(t)
	waitForNewBlock(t)

	totalStake := int64(0)
	for i := 0; i < nodeCount; i++ {
		amt := util.RandInt64(1000000 - 1) // fee is always 1000
		err := broadcastBondTransaction(t, tSigners[tNodeIdx1], tSigners[i].PublicKey(), amt, 1000)
		if err != nil {
			panic(fmt.Sprintf("Error on broadcasting transaction: %v", err))
		}
		fmt.Printf("Staking %v to %v\n", amt, tSigners[i].Address())
		incSequence(t, tSigners[tNodeIdx1].Address())
		totalStake += amt
	}

	// Uncomment this for debugging the test
	// go func() {
	// 	file, _ := os.OpenFile("./debug.log", os.O_CREATE|os.O_WRONLY, 0666)
	// 	fmt.Fprintf(file, "total stake: %d\n\n", totalStake)

	// 	for {
	// 		for i := 0; i < nodeCount; i++ {
	// 			tNodes[i].Consensus().RoundProposal(tNodes[i].Consensus().HRS().Round())

	// 			fmt.Fprintf(file, "node %d: %s %s %s proposal: %v ",
	// 				i,
	// 				tNodes[i].Sync().Fingerprint(),
	// 				tNodes[i].State().Fingerprint(),
	// 				tNodes[i].Consensus().Fingerprint(),
	// 				tNodes[i].Consensus().RoundProposal(tNodes[i].Consensus().HRS().Round()) != nil)

	// 			votes := tNodes[i].Consensus().RoundVotes(tNodes[i].Consensus().HRS().Round())

	// 			for _, v := range votes {
	// 				fmt.Fprintf(file, "%s:%s,%s ", v.VoteType(), v.BlockHash().Fingerprint(), v.Signer().Fingerprint())
	// 			}
	// 			fmt.Fprintf(file, "\n")
	// 		}
	// 		fmt.Fprintf(file, "================================================================================\n")

	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	for i := 0; i < 20; i++ {
		waitForNewBlock(t)
	}

	fmt.Printf("Running the tests")

	exitCode := m.Run()

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
