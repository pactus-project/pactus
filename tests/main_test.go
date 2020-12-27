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

var tSigners map[string]crypto.Signer
var tConfigs map[string]*config.Config
var tNodes map[string]*node.Node
var tCurlAddress = "0.0.0.0:1337"

func TestMain(m *testing.M) {
	tSigners = make(map[string]crypto.Signer)
	tConfigs = make(map[string]*config.Config)
	tNodes = make(map[string]*node.Node)

	_, _, priv1 := crypto.GenerateTestKeyPair()
	_, _, priv2 := crypto.GenerateTestKeyPair()
	_, _, priv3 := crypto.GenerateTestKeyPair()
	_, _, priv4 := crypto.GenerateTestKeyPair()
	tSigners["node_1"] = crypto.NewSigner(priv1)
	tSigners["node_2"] = crypto.NewSigner(priv2)
	tSigners["node_3"] = crypto.NewSigner(priv3)
	tSigners["node_4"] = crypto.NewSigner(priv4)

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

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21000000000000)

	vals := make([]*validator.Validator, 4)
	i := 0
	for _, s := range tSigners {
		val := validator.NewValidator(s.PublicKey(), 0, i)
		vals[i] = val
		i++
	}
	genDoc := genesis.MakeGenesis("test", util.Now(), []*account.Account{acc}, vals, 1)

	tNodes["node_1"], _ = node.NewNode(genDoc, tConfigs["node_1"], tSigners["node_1"])
	tNodes["node_2"], _ = node.NewNode(genDoc, tConfigs["node_2"], tSigners["node_2"])
	tNodes["node_3"], _ = node.NewNode(genDoc, tConfigs["node_3"], tSigners["node_3"])
	tNodes["node_4"], _ = node.NewNode(genDoc, tConfigs["node_4"], tSigners["node_4"])

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

	exitCode := m.Run()

	tNodes["node_1"].Stop()
	tNodes["node_2"].Stop()
	tNodes["node_3"].Stop()
	tNodes["node_4"].Stop()

	os.Exit(exitCode)
}
