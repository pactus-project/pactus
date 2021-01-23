package main

import (
	"fmt"
	"path/filepath"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

// Init initializes a node for zarb blockchain
func Init() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Hidden = true

		workingDirOpt := c.String(cli.StringOpt{
			Name:  "w working-dir",
			Desc:  "Working directory to save configuration and genesis files.",
			Value: cmd.ZarbHomeDir(),
		})
		testnetOpt := c.Bool(cli.BoolOpt{
			Name:  "testnet",
			Desc:  "Initialize working directory for joining the testnet",
			Value: false,
		})

		c.LongDesc = "Initializing the working directory by new validator's private key and genesis file."
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			path, _ := filepath.Abs(*workingDirOpt)

			var gen *genesis.Genesis
			conf := config.DefaultConfig()

			if *testnetOpt {
				gen = genesis.Testnet()

				conf.Network.Name = "zarb-testnet"
				conf.Network.Bootstrap.Addresses = []string{""}
				conf.Network.Bootstrap.MinPeerThreshold = 1
			} else {
				name := fmt.Sprintf("zarb-local-%v", cmd.RandomHex(2))
				gen, _ = makeGenesis(*workingDirOpt, name)
			}

			// save genesis file to file system
			genFile := path + "/genesis.json"
			if err := gen.SaveToFile(genFile); err != nil {
				cmd.PrintErrorMsg("Failed to write genesis file: %v", err)
				return
			}

			// save config file to file system
			confFile := path + "/config.toml"
			if err := conf.SaveToFile(confFile); err != nil {
				cmd.PrintErrorMsg("Failed to write config file: %v", err)
				return
			}

			fmt.Println()
			cmd.PrintSuccessMsg("A zarb node is successfully initialized at %v", path)
		}
	}
}

// makeGenesis makes genisis file while on initialize
func makeGenesis(workingDir string, chainName string) (*genesis.Genesis, error) {

	// create  accounts for genesis
	accs := make([]*account.Account, 5)
	// Treasury account
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(2100000000000000)

	accs[0] = acc

	for i := 1; i < len(accs); i++ {
		k := key.GenKey()
		if err := key.EncryptKeyToFile(k, workingDir+"/keys/"+k.Address().String()+".json", "", ""); err != nil {
			return nil, err
		}
		acc := account.NewAccount(k.Address(), i+1)
		acc.AddToBalance(1000000)

		accs[i] = acc
	}

	// create validator account for genesis
	k := key.GenKey()
	if err := key.EncryptKeyToFile(k, workingDir+"/validator_key.json", "", ""); err != nil {
		return nil, err
	}
	val := validator.NewValidator(k.PublicKey(), 0, 0)
	vals := []*validator.Validator{val}

	// create genesis
	gen := genesis.MakeGenesis(chainName, util.RoundNow(60), accs, vals, param.MainnetParams())
	return gen, nil
}
