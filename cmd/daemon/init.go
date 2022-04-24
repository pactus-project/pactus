package main

import (
	"fmt"
	"path/filepath"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/wallet"
)

// Init initializes a node for zarb blockchain
func Init() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {

		workingDirOpt := c.String(cli.StringOpt{
			Name:  "w working-dir",
			Desc:  "Working directory to save node configuration and genesis files.",
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
			workingDir, _ := filepath.Abs(*workingDirOpt)
			if !util.IsDirNotExistsOrEmpty(workingDir) {
				cmd.PrintErrorMsg("The working directory is not empty: %s", workingDir)
				return
			}

			cmd.PrintInfoMsg("Creating wallet...")
			mnemonic := wallet.GenerateMnemonic()
			cmd.PrintLine()
			cmd.PrintInfoMsg("Your wallet seed:")
			cmd.PrintInfoMsg("\"" + mnemonic + "\"")
			cmd.PrintWarnMsg("Write down your 12 word mnemonic on a piece of paper to recover your validator key in future.")
			cmd.PrintLine()
			confirmed := cmd.PromptConfirm("Do you want to continue?")
			if !confirmed {
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Please enter a password for wallet")
			password := cmd.PromptPassword("Password: ", true)
			walletPath := cmd.ZarbDefaultWalletPath(workingDir)

			// To make process faster, we update the password
			// after creating the addresses
			wallet, err := wallet.FromMnemonic(walletPath, mnemonic, "", 0)
			if err != nil {
				cmd.PrintErrorMsg("Failed to create wallet: %v", err)
				return
			}
			cmd.PrintInfoMsg("Wallet created successfully")
			valAddrStr, err := wallet.MakeNewAddress("", "Validator address")
			if err != nil {
				cmd.PrintErrorMsg("Failed to create validator address: %v", err)
				return
			}
			rewardAddrStr, err := wallet.MakeNewAddress("", "Reward address")
			if err != nil {
				cmd.PrintErrorMsg("Failed to create reward address: %v", err)
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
				conf.State.RewardAddress = rewardAddrStr
			} else {
				valPubStr, err := wallet.PublicKey("", valAddrStr)
				if err != nil {
					cmd.PrintErrorMsg("Failed to get validator public key: %v", err)
					return
				}
				valPub, err := bls.PublicKeyFromString(valPubStr)
				if err != nil {
					cmd.PrintErrorMsg("Failed to create validator public key: %v", err)
					return
				}

				gen = makeLocalGenesis(valPub)
				conf.Network.Name = "local-test"
			}

			// Save genesis file to file system
			genFile := cmd.ZarbGenesisPath(workingDir)
			if err := gen.SaveToFile(genFile); err != nil {
				cmd.PrintErrorMsg("Failed to write genesis file: %v", err)
				return
			}

			// Save config file to file system
			confFile := cmd.ZarbConfigPath(workingDir)
			if err := conf.SaveToFile(confFile); err != nil {
				cmd.PrintErrorMsg("Failed to write config file: %v", err)
				return
			}

			err = wallet.UpdatePassword("", password)
			if err != nil {
				cmd.PrintErrorMsg("Failed to update wallet password: %v", err)
				return
			}
			err = wallet.Save()
			if err != nil {
				cmd.PrintErrorMsg("Failed to save wallet: %v", err)
				return
			}
			cmd.PrintLine()

			cmd.PrintSuccessMsg("A zarb node is successfully initialized at %v", workingDir)
			cmd.PrintInfoMsg("You validator address is: %v", valAddrStr)
			cmd.PrintInfoMsg("You reward address is: %v", rewardAddrStr)
			cmd.PrintLine()
			cmd.PrintInfoMsg("You can start the node by running this command:")
			cmd.PrintInfoMsg("./zarb-daemon start -w %v", workingDir)
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
