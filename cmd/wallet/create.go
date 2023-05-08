package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
)

// Generate creates a new wallet.
func Generate() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		testnetOpt := c.Bool(cli.BoolOpt{
			Name:  "testnet",
			Desc:  "creating wallet for testnet",
			Value: false,
		})

		entropyOpt := c.Int(cli.IntOpt{
			Name:  "entropy",
			Desc:  "Entropy bit length",
			Value: 128,
		})

		c.Before = func() {}
		c.Action = func() {
			password := cmd.PromptPassword("Password", true)
			mnemonic := wallet.GenerateMnemonic(*entropyOpt)

			network := wallet.NetworkMainNet
			if *testnetOpt {
				network = wallet.NetworkTestNet
			}
			wallet, err := wallet.Create(*pathArg, mnemonic, password, network)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintSuccessMsg("Wallet created successfully at: %s", wallet.Path())
			cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
			cmd.PrintWarnMsg("Please keep your seed in a safe place; if you lose it, you will not be able to restore your wallet.")
		}
	}
}

// ChangePassword updates the wallet password.
func ChangePassword() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			wallet, err := openWallet()
			cmd.FatalErrorCheck(err)

			oldPassword := getPassword(wallet, *passOpt)
			newPassword := cmd.PromptPassword("New Password", true)

			err = wallet.UpdatePassword(oldPassword, newPassword)
			cmd.FatalErrorCheck(err)

			err = wallet.Save()
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintWarnMsg("Wallet password updated")
		}
	}
}
