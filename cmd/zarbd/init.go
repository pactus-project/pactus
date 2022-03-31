package main

import (
	"fmt"
	"path/filepath"

	cli "github.com/jawher/mow.cli"
	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

// Init initializes a node for zarb blockchain
func Init() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {

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

			if !util.IsDirNotExistsOrEmpty(path) {
				cmd.PrintErrorMsg("The workspace directory is not empty: %v", path)
				return
			}

			// Generate key for the validator and save it to file system
			entropy, _ := bip39.NewEntropy(128)
			mnemonic, _ := bip39.NewMnemonic(entropy)
			seed := bip39.NewSeed(mnemonic, "")
			prv, err := bls.PrivateKeyFromSeed(seed,nil)
			if err != nil {
				cmd.PrintErrorMsg("Failed to create key from the seed: %v", err)
				return
			}
			err = util.WriteFile(path+"/validator_key", []byte(prv.String()))
			if err != nil {
				cmd.PrintErrorMsg("Failed to write key file: %v", err)
				return
			}

			var gen *genesis.Genesis
			conf := config.DefaultConfig()

			if *testnetOpt {
				gen = genesis.Testnet()

				conf.Network.Name = "perdana-testnet"
				conf.Network.Bootstrap.Addresses = []string{"/ip4/172.104.169.94/tcp/21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq"}
				conf.Network.Bootstrap.MinThreshold = 4
				conf.Network.Bootstrap.MaxThreshold = 8
			} else {
				pub := prv.PublicKey()
				gen = makeLocalGenesis(pub.(*bls.PublicKey))
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
			cmd.PrintInfoMsg("You validator address is: %v", prv.PublicKey().Address())
			cmd.PrintInfoMsg("mnemonic: \"" + mnemonic + "\"")
			cmd.PrintWarnMsg("Write down your 12 word mnemonic on a piece of paper to recover your validator key in future.")
		}
	}
}

// makeLocalGenesis makes genisis file for the local network
func makeLocalGenesis(pub *bls.PublicKey) *genesis.Genesis {
	// Treasury account
	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)
	accs := []*account.Account{acc}

	val := validator.NewValidator(pub, 0)
	vals := []*validator.Validator{val}

	// create genesis
	gen := genesis.MakeGenesis(util.RoundNow(60), accs, vals, param.DefaultParams())
	return gen
}
