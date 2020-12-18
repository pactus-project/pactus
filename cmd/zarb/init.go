package main

import (
	"fmt"
	"path/filepath"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/keystore/key"
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
		chainNameOpt := c.String(cli.StringOpt{
			Name: "n chain-name",
			Desc: "A name for the blockchain",
		})

		c.Spec = "[-w=<path>] [-n=<name>]"
		c.LongDesc = "Initializing the working directory"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			// Check chain-name for genesis
			if *chainNameOpt == "" {
				*chainNameOpt = fmt.Sprintf("local-chain-%v", cmd.RandomHex(2))
			}

			path, _ := filepath.Abs(*workingDirOpt)
			gen := makeGenesis(*workingDirOpt, *chainNameOpt)
			conf := makeConfigfile()

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
func makeGenesis(workingDir string, chainName string) *genesis.Genesis {

	// create  accounts for genesis
	accs := make([]*account.Account, 5)
	// Treasury account
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21000000000000)

	accs[0] = acc

	for i := 1; i < len(accs); i++ {
		k := key.GenKey()
		if err := key.EncryptKeyFile(k, workingDir+"/keys/"+k.Address().String()+".json", "", ""); err != nil {
			return nil
		}
		acc := account.NewAccount(k.Address(), i+1)
		acc.AddToBalance(1000000)

		accs[i] = acc
	}

	// create validator account for genesis
	k := key.GenKey()
	if err := key.EncryptKeyFile(k, workingDir+"/validator_key.json", "", ""); err != nil {
		return nil
	}
	val := validator.NewValidator(k.PublicKey(), 0, 0)
	vals := []*validator.Validator{val}

	tm := time.Now().Truncate(0).UTC()

	// create genesis
	gen := genesis.MakeGenesis(chainName, tm, accs, vals, 10)
	return gen

}

// makeConfigfile makes configuration file
func makeConfigfile() *config.Config {
	conf := config.DefaultConfig()
	return conf

}
