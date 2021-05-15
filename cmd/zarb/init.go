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

		mainnetOpt := c.Bool(cli.BoolOpt{
			Name:  "mainnet",
			Desc:  "Initialize working directory for joining the mainnet",
			Value: false,
		})

		c.LongDesc = "Initializing the working directory by new validator's private key and genesis file."
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			path, _ := filepath.Abs(*workingDirOpt)

			if !util.IsDirNotExistsOrEmpty(path) {
				cmd.PrintErrorMsg("The workspace directory is not empty: %v", path)
				return
			}

			// TODO: Show Mnemonics for validator key to user
			// Generate key for the validator and save it to file system
			valKey := key.GenerateRandomKey()
			if err := key.EncryptKeyToFile(valKey, path+"/validator_key.json", "", ""); err != nil {
				cmd.PrintErrorMsg("Failed to crate validator key: %v", err)
				return
			}

			var gen *genesis.Genesis
			conf := config.DefaultConfig()

			if *testnetOpt {
				gen = genesis.Testnet()

				conf.Network.Name = "zarb-testnet"
				conf.Network.Bootstrap.Addresses = []string{"/ip4/139.162.135.180/tcp/31887/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq"}
				conf.Network.Bootstrap.MinThreshold = 4
				conf.Network.Bootstrap.MaxThreshold = 8
			} else if *mainnetOpt {
				gen = genesis.Mainnet()

				conf.Network.Name = "zarb"
				conf.Network.Bootstrap.Addresses = []string{"/ip4/172.104.186.100/tcp/8421/p2p/12D3KooWLB7zCZ2VV1AtqHwYy2RBgpxdtwYRYt1ZU7iECFNfpks6", "/ip4/139.177.199.21/tcp/8421/p2p/12D3KooWMjGbsP2XbR11RmjevPvTsC33qT48sbqiBhB9ekoFiedx"}
			} else {
				gen, _ = makeLocalGenesis(*workingDirOpt, valKey.PublicKey())
				conf.Network.Name = "zarb-local"
			}

			// Save genesis file to file system
			genFile := path + "/genesis.json"
			if err := gen.SaveToFile(genFile); err != nil {
				cmd.PrintErrorMsg("Failed to write genesis file: %v", err)
				return
			}

			// Save config file to file system
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

// makeLocalGenesis makes genisis file for the local network
func makeLocalGenesis(workingDir string, pub crypto.PublicKey) (*genesis.Genesis, error) {

	// Treasury account
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)
	accs := []*account.Account{acc}

	val := validator.NewValidator(pub, 0)
	vals := []*validator.Validator{val}

	// create genesis
	gen := genesis.MakeGenesis(util.RoundNow(60), accs, vals, param.DefaultParams())
	return gen, nil
}
