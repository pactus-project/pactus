package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

// Init initializes a node for the Pactus blockchain.
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
			mnemonic := wallet.GenerateMnemonic(128)
			cmd.PrintLine()
			cmd.PrintInfoMsg("Your wallet seed is:")
			cmd.PrintInfoMsgBold("   " + mnemonic)
			cmd.PrintLine()
			cmd.PrintWarnMsg("Write down this seed on a piece of paper to recover your validator key in future.")
			cmd.PrintLine()
			confirmed := cmd.PromptConfirm("Do you want to continue")
			if !confirmed {
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Enter a password for wallet")
			password := cmd.PromptPassword("Password", true)
			walletPath := cmd.PactusDefaultWalletPath(workingDir)

			cmd.PrintLine()
			cmd.PrintInfoMsg("How many validators do you want to create?")
			cmd.PrintInfoMsg("Each node can run up to 32 validators, and each validator can hold up to 1000 staked coins.")
			cmd.PrintInfoMsg("You can define validators based on the amount of coins you want to stake.")
			numValidatorsStr := cmd.PromptInputWithSuggestion("Number of Validators", "7")
			numValidators, err := strconv.Atoi(numValidatorsStr)
			cmd.FatalErrorCheck(err)


			

			if numValidators < 1 || numValidators > 32 {
				cmd.PrintErrorMsg("Invalid validator number")
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Creating wallet...")

			// To make process faster, we update the password
			// after creating the addresses
			network := wallet.NetworkMainNet
			if *testnetOpt {
				network = wallet.NetworkTestNet
			}
			wallet, err := wallet.Create(walletPath, mnemonic, "", network)
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintSuccessMsg("Wallet created successfully at %s.", walletPath)
			cmd.PrintLine()
			cmd.PrintInfoMsg("Generating keys...")
			cmd.PrintLine()
			cmd.PrintInfoMsgBold("Validator addresses:")
			for i := 0; i < numValidators; i++ {
				valAddrStr, err := wallet.DeriveNewAddress(fmt.Sprintf("Validator address %v", i+1))
				cmd.FatalErrorCheck(err)

				cmd.PrintInfoMsg("%v- %s", i+1, valAddrStr)
			}
			cmd.PrintLine()

			cmd.PrintInfoMsgBold("Reward addresses:")
			for i := 0; i < numValidators; i++ {
				rewardAddrStr, err := wallet.DeriveNewAddress(fmt.Sprintf("Reward address %v", i+1))
				cmd.FatalErrorCheck(err)

				cmd.PrintInfoMsg("%v- %s", i+1, rewardAddrStr)
			}
			cmd.PrintLine()
			cmd.PrintInfoMsg("Initializing node...")
			cmd.PrintLine()

			var gen *genesis.Genesis
			confFile := cmd.PactusConfigPath(workingDir)

			if *testnetOpt {
				cmd.PrintInfoMsg("Network: Testnet")

				gen = genesis.Testnet()

				// Save config for testnet
				err := config.SaveTestnetConfig(confFile, numValidators)
				cmd.FatalErrorCheck(err)
			} else if *localnetOpt {
				cmd.PrintInfoMsg("Network: Localnet")

				info := wallet.AddressInfo(wallet.AddressLabels()[0].Address)
				valPub := info.Pub.(*bls.PublicKey)
				gen = makeLocalGenesis(valPub)

				// Save config for localnet
				err := config.SaveLocalnetConfig(confFile)
				cmd.FatalErrorCheck(err)
			} else {
				panic("not yet!")
				// gen = genesis.Mainnet()

				// // Save config for mainnet
				// if err := config.SaveMainnetConfig(confFile, rewardAddrStr); err != nil {
				// 	cmd.PrintErrorMsg("Failed to write config file: %v", err)
				// 	return
				// }
			}

			genFile := cmd.PactusGenesisPath(workingDir)
			err = gen.SaveToFile(genFile)
			cmd.FatalErrorCheck(err)

			err = wallet.UpdatePassword("", password)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintSuccessMsg("A pactus node is successfully initialized at %v", workingDir)
			cmd.PrintLine()
			cmd.PrintInfoMsg("You can start the node by running this command:")
			cmd.PrintInfoMsg("./pactus-daemon start -w %v", workingDir)
		}
	}
}

// makeLocalGenesis makes genesis file for the local network.
func makeLocalGenesis(pub *bls.PublicKey) *genesis.Genesis {
	// Treasury account
	acc := account.NewAccount(0)
	acc.AddToBalance(21 * 1e14)
	accs := map[crypto.Address]*account.Account{}

	val := validator.NewValidator(pub, 0)
	vals := []*validator.Validator{val}

	// create genesis
	params := param.DefaultParams()
	params.BlockVersion = 63
	gen := genesis.MakeGenesis(util.RoundNow(60), accs, vals, params)
	return gen
}
