package main

import (
	"fmt"
	"path/filepath"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/node/config"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/genesis"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

// Init initializes a node for pactus blockchain.
func Init() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		workingDirOpt := c.String(cli.StringOpt{
			Name:  "w working-dir",
			Desc:  "Working directory to save node configuration and genesis files.",
			Value: cmd.PactusHomeDir(),
		})
		testnetOpt := c.Bool(cli.BoolOpt{
			Name:  "testnet",
			Desc:  "Initialize working directory for joining the testnet",
			Value: true, // TODO: make it false after mainnet launch
		})
		localnetOpt := c.Bool(cli.BoolOpt{
			Name:  "localnet",
			Desc:  "Initialize working directory for localnet (for developers)",
			Value: false,
		})

		c.LongDesc = "Initializing the working directory by new validator's private key and genesis file."
		c.Before = func() { fmt.Println(cmd.Pactus) }
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
			confirmed := cmd.PromptConfirm("Do you want to continue")
			if !confirmed {
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Please enter a password for wallet")
			password := cmd.PromptPassword("Password", true)
			walletPath := cmd.PactusDefaultWalletPath(workingDir)

			// To make process faster, we update the password
			// after creating the addresses
			network := wallet.NetworkMainNet
			if *testnetOpt {
				network = wallet.NetworkTestNet
			}
			wallet, err := wallet.FromMnemonic(walletPath, mnemonic, "", network)
			if err != nil {
				cmd.PrintErrorMsg("Failed to create wallet: %v", err)
				return
			}
			cmd.PrintInfoMsg("Wallet created successfully")
			valAddrStr, err := wallet.DeriveNewAddress("Validator address")
			if err != nil {
				cmd.PrintErrorMsg("Failed to create validator address: %v", err)
				return
			}
			rewardAddrStr, err := wallet.DeriveNewAddress("Reward address")
			if err != nil {
				cmd.PrintErrorMsg("Failed to create reward address: %v", err)
				return
			}

			var gen *genesis.Genesis
			confFile := cmd.PactusConfigPath(workingDir)

			if *testnetOpt {
				gen = genesis.Testnet()

				// Save config for testnet
				if err := config.SaveTestnetConfig(confFile, rewardAddrStr); err != nil {
					cmd.PrintErrorMsg("Failed to write config file: %v", err)
					return
				}
			} else if *localnetOpt {
				info := wallet.AddressInfo(valAddrStr)
				if info == nil {
					cmd.PrintErrorMsg("Failed to get validator public key")
					return
				}
				valPub := info.Pub.(*bls.PublicKey)
				gen = makeLocalGenesis(valPub)

				// Save config for localnet
				if err := config.SaveLocalnetConfig(confFile, rewardAddrStr); err != nil {
					cmd.PrintErrorMsg("Failed to write config file: %v", err)
					return
				}
			} else {
				panic("not yet!")
				// gen = genesis.Mainnet()

				// // Save config for mainnet
				// if err := config.SaveMainnetConfig(confFile, rewardAddrStr); err != nil {
				// 	cmd.PrintErrorMsg("Failed to write config file: %v", err)
				// 	return
				// }
			}

			// Save genesis file
			genFile := cmd.PactusGenesisPath(workingDir)
			if err := gen.SaveToFile(genFile); err != nil {
				cmd.PrintErrorMsg("Failed to write genesis file: %v", err)
				return
			}

			err = wallet.UpdatePassword("", password)
			if err != nil {
				cmd.PrintErrorMsg("Failed to update wallet password: %v", err)
				return
			}

			// Save wallet
			err = wallet.Save()
			if err != nil {
				cmd.PrintErrorMsg("Failed to save wallet: %v", err)
				return
			}
			cmd.PrintLine()

			cmd.PrintSuccessMsg("A pactus node is successfully initialized at %v", workingDir)
			cmd.PrintInfoMsg("You validator address is: %v", valAddrStr)
			cmd.PrintInfoMsg("You reward address is: %v", rewardAddrStr)
			cmd.PrintLine()
			cmd.PrintInfoMsg("You can start the node by running this command:")
			cmd.PrintInfoMsg("./pactus-daemon start -w %v", workingDir)
		}
	}
}

// makeLocalGenesis makes genisis file for the local network.
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
