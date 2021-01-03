package tests

import (
	"os"
	"testing"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/node"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var tSigners map[string]*crypto.Signer
var tConfigs map[string]*config.Config
var tNodes map[string]*node.Node
var tCurlAddress = "0.0.0.0:1337"
var tGenDoc *genesis.Genesis

func TestMain(m *testing.M) {
	tSigners = make(map[string]*crypto.Signer)
	tConfigs = make(map[string]*config.Config)
	tNodes = make(map[string]*node.Node)

	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()
	signer1 := crypto.NewSigner(priv1)
	signer2 := crypto.NewSigner(priv2)
	signer3 := crypto.NewSigner(priv3)
	signer4 := crypto.NewSigner(priv4)

	tSigners["node_1"] = &signer1
	tSigners["node_2"] = &signer2
	tSigners["node_3"] = &signer3
	tSigners["node_4"] = &signer4

	tConfigs["node_1"] = config.DefaultConfig()
	tConfigs["node_2"] = config.DefaultConfig()
	tConfigs["node_3"] = config.DefaultConfig()
	tConfigs["node_4"] = config.DefaultConfig()

	tConfigs["node_1"].Sync.StartingTimeout = 0
	tConfigs["node_2"].Sync.StartingTimeout = 0
	tConfigs["node_3"].Sync.StartingTimeout = 0
	tConfigs["node_4"].Sync.StartingTimeout = 0

	tConfigs["node_1"].State.Store.Path = util.TempDirPath()
	tConfigs["node_2"].State.Store.Path = util.TempDirPath()
	tConfigs["node_3"].State.Store.Path = util.TempDirPath()
	tConfigs["node_4"].State.Store.Path = util.TempDirPath()

	tConfigs["node_1"].Network.NodeKeyFile = util.TempFilePath()
	tConfigs["node_2"].Network.NodeKeyFile = util.TempFilePath()
	tConfigs["node_3"].Network.NodeKeyFile = util.TempFilePath()
	tConfigs["node_4"].Network.NodeKeyFile = util.TempFilePath()

	tConfigs["node_1"].Http.Address = tCurlAddress
	tConfigs["node_2"].Http.Enable = false
	tConfigs["node_3"].Http.Enable = false
	tConfigs["node_4"].Http.Enable = false

	tConfigs["node_1"].Capnp.Address = "0.0.0.0:0"
	tConfigs["node_2"].Capnp.Enable = false
	tConfigs["node_3"].Capnp.Enable = false
	tConfigs["node_4"].Capnp.Enable = false

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21000000000000)

	vals := make([]*validator.Validator, 4)
	vals[0] = validator.NewValidator(tSigners["node_1"].PublicKey(), 0, 0)
	vals[1] = validator.NewValidator(tSigners["node_2"].PublicKey(), 1, 0)
	vals[2] = validator.NewValidator(tSigners["node_3"].PublicKey(), 2, 0)
	vals[3] = validator.NewValidator(tSigners["node_4"].PublicKey(), 3, 0)
	tGenDoc = genesis.MakeGenesis("test", util.Now(), []*account.Account{acc}, vals, 1)

	tNodes["node_1"], _ = node.NewNode(tGenDoc, tConfigs["node_1"], *tSigners["node_1"])
	tNodes["node_2"], _ = node.NewNode(tGenDoc, tConfigs["node_2"], *tSigners["node_2"])
	tNodes["node_3"], _ = node.NewNode(tGenDoc, tConfigs["node_3"], *tSigners["node_3"])
	tNodes["node_4"], _ = node.NewNode(tGenDoc, tConfigs["node_4"], *tSigners["node_4"])

	err := tNodes["node_1"].Start()
	if err != nil {
		panic(err)
	}
	err = tNodes["node_2"].Start()
	if err != nil {
		panic(err)
	}
	err = tNodes["node_3"].Start()
	if err != nil {
		panic(err)
	}
	err = tNodes["node_4"].Start()
	if err != nil {
		panic(err)
	}

	t := testing.T{}
	getBlockAt(&t, 1)

	exitCode := m.Run()

	// Random crash here
	// TODO: fix ma later
	// tNodes["node_1"].Stop()
	// tNodes["node_2"].Stop()
	// tNodes["node_3"].Stop()
	// tNodes["node_4"].Stop()

	os.Exit(exitCode)
}
